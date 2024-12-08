package helper

import (
	"appointment_management_system/internal/pkg/constants"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Utility function to get transaction from context
func GetDB(c *gin.Context) (*gorm.DB, error) {
	// Retrieve the value from the context
	tx, exists := c.Get(constants.DBTRANSACTION)
	if !exists {
		return nil, fmt.Errorf("failed to get database transaction from context: key not found")
	}

	// Validate the type of the value
	db, ok := tx.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("failed to get database transaction from context: invalid type")
	}

	return db, nil
}
