package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTransactionMiddleware(t *testing.T) {
	// Create master and slave database connections using SQLite in-memory databases for testing
	master, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err, "Failed to create master database connection")

	slave, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err, "Failed to create slave database connection")

	// Create a Gin engine with the middleware
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(TransactionMiddleware(master, slave))

	// Add a test route for write operations
	router.POST("/write", func(c *gin.Context) {
		// Retrieve the transaction from the context
		tx, exists := c.Get("db_tx")
		assert.True(t, exists, "Transaction should exist in context")
		assert.NotNil(t, tx, "Transaction should not be nil")
		c.JSON(http.StatusOK, gin.H{"status": "write successful"})
	})

	// Add a test route for read operations
	router.GET("/read", func(c *gin.Context) {
		// Retrieve the transaction from the context
		tx, exists := c.Get("db_tx")
		assert.True(t, exists, "Transaction should exist in context")
		assert.NotNil(t, tx, "Transaction should not be nil")
		c.JSON(http.StatusOK, gin.H{"status": "read successful"})
	})

	// Test write operation
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/write", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200 for write operation")
	assert.JSONEq(t, `{"status":"write successful"}`, w.Body.String())

	// Test read operation
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/read", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200 for read operation")
	assert.JSONEq(t, `{"status":"read successful"}`, w.Body.String())
}

// MockDB is used to simulate the behavior of a database transaction
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Begin(opts ...*sql.TxOptions) *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Session(config *gorm.Session) *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Commit() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func TestTransactionMiddlewareRollbackAndCommit(t *testing.T) {
	// Create mock master and slave DB connections
	masterMock, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	slave, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// Create Gin engine
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(TransactionMiddleware(masterMock, slave))

	// Add a route to simulate successful processing
	router.POST("/commit", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "commit successful"})
	})

	// Add a route to simulate an error during processing
	router.POST("/rollback", func(c *gin.Context) {
		c.Error(errors.New("simulated processing error")) // Register an error to trigger rollback
		c.JSON(http.StatusInternalServerError, gin.H{"error": "rollback triggered"})
	})

	router.POST("/panic", func(c *gin.Context) {
		panic("mock panic")
	})

	// Test commit case
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/commit", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200 for commit case")
	assert.JSONEq(t, `{"status":"commit successful"}`, w.Body.String())

	// Test rollback case
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/rollback", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status code 500 for rollback case")
	assert.JSONEq(t, `{"error":"rollback triggered"}`, w.Body.String())

	// Test panic case
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/panic", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status code 500 for panic case")
	assert.JSONEq(t, `{"error":"Unexpected server error"}`, w.Body.String())

}

func TestHandleWriteTransaction_RollbackError(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a mock SQL database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, got error: %v", err)
	}
	defer db.Close()

	// Open a GORM DB connection using the mock database
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm db, got error: %v", err)
	}

	// Set expectations
	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(errors.New("rollback error"))

	// Begin a transaction
	tx := gormDB.Begin()
	if tx.Error != nil {
		t.Fatalf("Failed to begin transaction: %v", tx.Error)
	}

	// Create a test Gin context with a response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate an error in the Gin context
	c.Error(errors.New("simulated context error"))

	// Call the function under test
	handleWriteTransaction(tx, c)

	// Assert the response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error": "Failed to rollback transaction"}`, w.Body.String())

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestHandleWriteTransaction_CommitError(t *testing.T) {
	// Initialize Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a mock SQL database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock sql db, got error: %v", err)
	}
	defer db.Close()

	// Open a GORM DB connection using the mock database
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm db, got error: %v", err)
	}

	// Set expectations
	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(errors.New("commit error"))

	// Begin a transaction
	tx := gormDB.Begin()
	if tx.Error != nil {
		t.Fatalf("Failed to begin transaction: %v", tx.Error)
	}

	// Create a test Gin context with a response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Ensure no errors in the Gin context
	assert.Empty(t, c.Errors)

	// Call the function under test
	handleWriteTransaction(tx, c)

	// Assert the response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error": "Failed to commit transaction"}`, w.Body.String())

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}