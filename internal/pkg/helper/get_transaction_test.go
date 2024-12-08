package helper

import (
	"appointment_management_system/internal/pkg/constants"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successfully retrieves database transaction", func(t *testing.T) {
		// Create a test context
		ctx := &gin.Context{}
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		// Set a valid gorm.DB transaction in the context
		ctx.Set(constants.DBTRANSACTION, db)

		// Call GetDB
		retrievedDB, err := GetDB(ctx)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, retrievedDB)
		assert.Equal(t, db, retrievedDB)
	})

	t.Run("fails if transaction key is not found", func(t *testing.T) {
		// Create a test context
		ctx := &gin.Context{}

		// Call GetDB without setting DBTRANSACTION
		retrievedDB, err := GetDB(ctx)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, retrievedDB)
		assert.EqualError(t, err, "failed to get database transaction from context: key not found")
	})

	t.Run("fails if value is not of type *gorm.DB", func(t *testing.T) {
		// Create a test context
		ctx := &gin.Context{}

		// Set an invalid value in the context
		ctx.Set(constants.DBTRANSACTION, "invalid_type")

		// Call GetDB
		retrievedDB, err := GetDB(ctx)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, retrievedDB)
		assert.EqualError(t, err, "failed to get database transaction from context: invalid type")
	})
}
