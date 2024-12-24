package transaction

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const DBTRANSACTION = "db_tx"

func CustomTransaction(master, slave *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx, isWrite := startTransaction(master, slave, c.Request.Method)
		if tx == nil {
			return
		}

		c.Set(DBTRANSACTION, tx)
		c.Next()

		finalizeTransaction(tx, isWrite, c)
	}
}

func startTransaction(master, slave *gorm.DB, method string) (*gorm.DB, bool) {
	isWrite := isWriteMethod(method)
	var tx *gorm.DB
	if isWrite {
		tx = master.Session(&gorm.Session{
			PrepareStmt: true,
		}).Begin()
	} else {
		tx = slave.Session(&gorm.Session{AllowGlobalUpdate: false,
			PrepareStmt: true,
		})
	}

	if tx.Error != nil {
		log.WithFields(log.Fields{"error": tx.Error}).Error("Failed to start transaction")
		return nil, false
	}

	log.Info("Transaction started")
	return tx, isWrite
}

func finalizeTransaction(tx *gorm.DB, isWrite bool, c *gin.Context) {
	if isWrite {
		handleWriteTransaction(tx, c)
	} else {
		log.Info("Read-only transaction completed successfully")
	}
}

func handleWriteTransaction(tx *gorm.DB, c *gin.Context) {
	if len(c.Errors) > 0 {
		if err := tx.Rollback().Error; err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed to rollback transaction")
			c.Abort()
			return
		}
		log.Info("Transaction rolled back due to errors")
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed to commit transaction")
		c.Abort()
		return
	}

	log.Info("Transaction committed successfully")
}

func isWriteMethod(method string) bool {
	writeMethods := map[string]struct{}{
		"POST":   {},
		"PUT":    {},
		"PATCH":  {},
		"DELETE": {},
	}

	_, exists := writeMethods[strings.ToUpper(method)]
	return exists
}

// Utility function to get transaction from context
func GetTransaction(c *gin.Context) (*gorm.DB, error) {
	// Retrieve the value from the context
	tx, exists := c.Get(DBTRANSACTION)
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
