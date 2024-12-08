package middleware

import (
	"appointment_management_system/internal/pkg/constants"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func TransactionMiddleware(master, slave *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx, isWrite := startTransaction(master, slave, c.Request.Method)
		if tx == nil {
			return
		}

		defer rollbackOnPanic(tx, c)

		c.Set(constants.DBTRANSACTION, tx)
		c.Next()

		finalizeTransaction(tx, isWrite, c)
	}
}

func startTransaction(master, slave *gorm.DB, method string) (*gorm.DB, bool) {
	isWrite := isWriteMethod(method)
	var tx *gorm.DB
	if isWrite {
		tx = master.Begin()
	} else {
		tx = slave.Session(&gorm.Session{AllowGlobalUpdate: false})
	}

	if tx.Error != nil {
		log.WithFields(log.Fields{"error": tx.Error}).Error("Failed to start transaction")
		return nil, false
	}

	log.Info("Transaction started")
	return tx, isWrite
}

func rollbackOnPanic(tx *gorm.DB, c *gin.Context) {
	if r := recover(); r != nil {
		tx.Rollback()
		log.WithField("panic", r).Error("Panic occurred, transaction rolled back")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unexpected server error",
		})
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
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to rollback transaction",
			})
			return
		}
		log.Info("Transaction rolled back due to errors")
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Failed to commit transaction")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction",
		})
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
