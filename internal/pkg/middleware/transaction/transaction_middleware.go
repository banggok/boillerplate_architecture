package transaction

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const DBTRANSACTION = "db_tx"
const DBREAD = "db_read"

func CustomTransaction(master, slave *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		read := slave.Session(&gorm.Session{AllowGlobalUpdate: false,
			PrepareStmt: true,
		})

		c.Set(DBREAD, read)

		isWrite := isWriteMethod(c.Request.Method)

		var tx *gorm.DB
		if isWrite {
			tx = master.Session(&gorm.Session{
				PrepareStmt: true,
			}).Begin()

			if tx.Error != nil {
				log.WithFields(log.Fields{"error": tx.Error}).Error("Failed to start transaction")
			}

			log.Info("Transaction started")

			c.Set(DBTRANSACTION, tx)
		}

		c.Next()

		finalizeTransaction(tx, isWrite, c)
	}
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
func GetTransaction(c *gin.Context, isWrite bool) (*gorm.DB, error) {
	txName := DBREAD

	if isWrite {
		txName = DBTRANSACTION
	}
	// Retrieve the value from the context
	tx, exists := c.Get(txName)
	if !exists {
		return nil, fmt.Errorf("failed to get database transaction from context: key not found")
	}

	// Validate the type of the value
	db, ok := tx.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("failed to get database transaction from context: invalid type")
	}

	if !isWrite {
		db = db.Session(&gorm.Session{})
	}

	return db, nil
}
