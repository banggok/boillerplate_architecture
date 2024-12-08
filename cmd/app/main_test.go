package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSetupConfiguration(t *testing.T) {
	// Call the function
	validator, rateLimiter := setupConfiguration()

	// Validate return values
	assert.NotNil(t, validator, "Validator should not be nil")
	assert.NotNil(t, rateLimiter, "RateLimiter should not be nil")
}

func TestSetupRoutes(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Mock validator and rate limiter
	validate := validator.New()
	rate := limiter.Rate{Period: time.Minute, Limit: 10}
	store := memory.NewStore()
	rateLimiter := limiter.New(store, rate)

	// Initialize an SQLite in-memory database
	mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize mock database: %v", err)
	}

	conn = DBConnection{
		master: mockDB,
		slave:  mockDB,
	}

	// Setup routes with mocked dependencies
	setupRoutes(router, validate, rateLimiter, conn, logrus.New())

	// Test health endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Health endpoint should return 200 OK")
	assert.JSONEq(t, `{"status":"healthy"}`, w.Body.String(), "Health endpoint response should match")
}

func TestGetConfigValue(t *testing.T) {
	// Test with environment variable set
	os.Setenv("TEST_ENV", "test_value")
	assert.Equal(t, "test_value", getConfigValue("TEST_ENV", "default_value"))
	os.Unsetenv("TEST_ENV")

	// Test with fallback value
	assert.Equal(t, "default_value", getConfigValue("TEST_ENV", "default_value"))
}

func TestGetConfigValueAsInt(t *testing.T) {
	// Test with valid integer value
	os.Setenv("TEST_INT", "123")
	assert.Equal(t, 123, getConfigValueAsInt("TEST_INT", 0))
	os.Unsetenv("TEST_INT")

	// Test with fallback value
	assert.Equal(t, 0, getConfigValueAsInt("TEST_INT", 0))

	// Test with invalid integer value
	os.Setenv("TEST_INT", "invalid")
	assert.Equal(t, 0, getConfigValueAsInt("TEST_INT", 0))
	os.Unsetenv("TEST_INT")
}

// MockServerRunner is a mock implementation of the ServerRunner interface
type MockServerRunner struct {
	mock.Mock
}

func (m *MockServerRunner) Run() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockServerRunner) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestRunServer(t *testing.T) {
	// Mock server runner
	mockRunner := new(MockServerRunner)

	// Simulate a blocking Run method that waits for a signal
	runBlocked := make(chan struct{}) // Used to block the Run method
	mockRunner.On("Run").Run(func(args mock.Arguments) {
		<-runBlocked // Block until the signal is sent
	}).Return(nil)

	// Mock graceful shutdown success
	mockRunner.On("Shutdown", mock.Anything).Return(nil)

	// Simulate signal channel
	signalCh := make(chan os.Signal, 1)

	// Track when the shutdown is complete
	shutdownComplete := make(chan struct{})

	// Start the server in a goroutine
	go func() {
		runServer(mockRunner, 5*time.Second, signalCh)
		close(shutdownComplete) // Notify that the server has shut down
	}()

	// Ensure the server is running before sending the signal
	time.Sleep(100 * time.Millisecond)

	// Send the signal
	t.Log("Sending signal: os.Interrupt")
	signalCh <- os.Interrupt

	// Unblock the Run method to simulate server exit
	close(runBlocked)

	// Wait for the server to shut down
	select {
	case <-shutdownComplete:
		t.Log("Server shutdown complete")
	case <-time.After(1 * time.Second):
		t.Fatal("Test timed out waiting for server shutdown")
	}

	// Verify the correct calls were made
	mockRunner.AssertCalled(t, "Run")
	mockRunner.AssertCalled(t, "Shutdown", mock.Anything)
}

func TestRunServer_ServerError(t *testing.T) {
	// Mock server runner
	mockRunner := new(MockServerRunner)
	mockRunner.On("Run").Return(errors.New("server error")) // Mock server run failure

	// Simulate signal channel
	signalCh := make(chan os.Signal, 1)

	// Call the runServer function
	runServer(mockRunner, 5*time.Second, signalCh)

	// Assert Run was called, but Shutdown was not
	mockRunner.AssertCalled(t, "Run")
	mockRunner.AssertNotCalled(t, "Shutdown")
}

func TestConnectMySQL_Success(t *testing.T) {
	// Mock gormOpen to use SQLite for testing
	originalGormOpen := gormOpen
	defer func() { gormOpen = originalGormOpen }()
	gormOpen = func(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error) {
		return gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	}

	// Call the function
	db, err := connectMySQL()

	// Assertions
	assert.NoError(t, err, "Expected no error when connecting to mock database")
	assert.NotNil(t, db, "Expected non-nil database connection")
}

func TestConnectMySQL_Failure(t *testing.T) {
	// Mock gormOpen to simulate a failure
	originalGormOpen := gormOpen
	defer func() { gormOpen = originalGormOpen }()
	gormOpen = func(dialector gorm.Dialector, opts ...gorm.Option) (*gorm.DB, error) {
		return nil, errors.New("mock connection failure")
	}

	// Call the function
	db, err := connectMySQL()

	// Assertions
	assert.Error(t, err, "Expected an error when failing to connect to the database")
	assert.Nil(t, db.master, "Expected nil database connection on failure")
	assert.Equal(t, "failed to connect to Master MySQL: mock connection failure", err.Error(), "Error message should match")
}

func TestGracefulShutdown_Success(t *testing.T) {
	// Mock ServerRunner
	mockRunner := new(MockServerRunner)
	mockRunner.On("Shutdown", mock.Anything).Return(nil)

	// Mock database connection
	mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err, "Should be able to open mock database")

	conn = DBConnection{
		master: mockDB,
		slave:  mockDB,
	} // Set the global db for testing

	// Call gracefulShutdown
	gracefulShutdown(mockRunner, 5*time.Second)

	// Assert server shutdown was called
	mockRunner.AssertCalled(t, "Shutdown", mock.Anything)

	// Assert database connection is closed
	sqlDB, err := mockDB.DB()
	assert.NoError(t, err, "Should retrieve underlying sql.DB")
	err = sqlDB.Ping()
	assert.Error(t, err, "Database connection should be closed")
}

func TestGracefulShutdown_ServerError(t *testing.T) {
	// Mock ServerRunner with shutdown error
	mockRunner := new(MockServerRunner)
	mockRunner.On("Shutdown", mock.Anything).Return(errors.New("mock shutdown error"))

	// Call gracefulShutdown
	assert.Panics(t, func() {
		gracefulShutdown(mockRunner, 5*time.Second)
	}, "Expected gracefulShutdown to panic on server shutdown error")

	// Assert server shutdown was called
	mockRunner.AssertCalled(t, "Shutdown", mock.Anything)
}

func TestGetServerAddress_HTTPServerRunner(t *testing.T) {
	// Create an instance of HTTPServerRunner
	runner := &HTTPServerRunner{
		server: &http.Server{
			Addr: ":8080",
		},
	}

	// Call the function
	address := getServerAddress(runner)

	// Assertions
	assert.Equal(t, ":8080", address, "Expected server address to be ':8080'")
}

func TestGetServerAddress_UnknownRunner(t *testing.T) {
	// Create a mock runner that does not implement HTTPServerRunner
	mockRunner := new(MockServerRunner)

	// Call the function
	address := getServerAddress(mockRunner)

	// Assertions
	assert.Equal(t, "unknown address", address, "Expected address to be 'unknown address' for non-HTTPServerRunner")
}
