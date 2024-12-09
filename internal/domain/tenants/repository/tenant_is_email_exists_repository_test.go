package repository

import (
	"appointment_management_system/internal/pkg/constants"
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create memory database: %v", err)
	}
	if err := db.AutoMigrate(&tenantModel{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}
	return db
}

func seedDatabase(t *testing.T, db *gorm.DB, email string) {
	if err := db.Create(&tenantModel{Email: email}).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}
}

func setupGinContext(db *gorm.DB) *gin.Context {
	ctx := &gin.Context{}
	ctx.Set(constants.DBTRANSACTION, db)
	return ctx
}

func runExecuteTest(t *testing.T, repo TenantIsEmailExistsRepository, ctx *gin.Context, email string, expectedExists *bool, expectedErr string) {
	exists, err := repo.Execute(ctx, email)

	if expectedErr != "" {
		assert.Error(t, err)
		assert.Nil(t, exists)
		assert.EqualError(t, err, expectedErr)
	} else {
		assert.NoError(t, err)
		assert.NotNil(t, exists)
		assert.Equal(t, *expectedExists, *exists)
	}
}

func TestTenantIsEmailExistsRepository_Execute(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTenantIsEmailExistsRepository()

	t.Run("email exists", func(t *testing.T) {
		seedDatabase(t, db, "existing@example.com")
		ctx := setupGinContext(db)
		x := true
		runExecuteTest(t, repo, ctx, "existing@example.com", &x, "")
	})

	t.Run("email does not exist", func(t *testing.T) {
		ctx := setupGinContext(db)
		x := false
		runExecuteTest(t, repo, ctx, "nonexistent@example.com", &x, "")
	})

	t.Run("email is empty", func(t *testing.T) {
		ctx := setupGinContext(nil)
		runExecuteTest(t, repo, ctx, "", nil, "email cannot be empty")
	})

	t.Run("database transaction error", func(t *testing.T) {
		originalHelper := helperGetDB
		defer func() { helperGetDB = originalHelper }()
		helperGetDB = func(_ *gin.Context) (*gorm.DB, error) {
			return nil, errors.New("database transaction error")
		}

		ctx := &gin.Context{}
		runExecuteTest(t, repo, ctx, "test@example.com", nil, "failed to get database transaction: database transaction error")
	})

	t.Run("database query error", func(t *testing.T) {
		sqlDB, _ := db.DB()
		sqlDB.Close()

		ctx := setupGinContext(db)
		runExecuteTest(t, repo, ctx, "test@example.com", nil, "failed to query database: sql: database is closed")
	})
}
