package repository

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func seedDatabaseForPhone(t *testing.T, db *gorm.DB, phone string) {
	if err := db.Create(&tenantModel{Phone: phone}).Error; err != nil {
		t.Fatalf("failed to seed database: %v", err)
	}
}

func runExecuteTestForPhone(t *testing.T, repo TenantIsPhoneExistsRepository, ctx *gin.Context, phone string, expectedExists *bool, expectedErr string) {
	exists, err := repo.Execute(ctx, phone)

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

func TestTenantIsPhoneExistsRepository_Execute(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTenantIsPhoneExistsRepository()

	t.Run("phone exists", func(t *testing.T) {
		seedDatabaseForPhone(t, db, "+1234567890")
		ctx := setupGinContext(db)
		x := true
		runExecuteTestForPhone(t, repo, ctx, "+1234567890", &x, "")
	})

	t.Run("phone does not exist", func(t *testing.T) {
		ctx := setupGinContext(db)
		x := false
		runExecuteTestForPhone(t, repo, ctx, "+0987654321", &x, "")
	})

	t.Run("phone is empty", func(t *testing.T) {
		ctx := setupGinContext(nil)
		runExecuteTestForPhone(t, repo, ctx, "", nil, "phone number cannot be empty")
	})

	t.Run("database transaction error", func(t *testing.T) {
		originalHelper := helperGetDB
		defer func() { helperGetDB = originalHelper }()
		helperGetDB = func(_ *gin.Context) (*gorm.DB, error) {
			return nil, errors.New("database transaction error")
		}

		ctx := &gin.Context{}
		runExecuteTestForPhone(t, repo, ctx, "+1234567890", nil, "failed to get database transaction: database transaction error")
	})

	t.Run("database query error", func(t *testing.T) {
		sqlDB, _ := db.DB()
		sqlDB.Close()

		ctx := setupGinContext(db)
		runExecuteTestForPhone(t, repo, ctx, "+1234567890", nil, "failed to query database: sql: database is closed")
	})
}
