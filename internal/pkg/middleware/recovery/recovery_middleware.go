package recovery

import (
	"appointment_management_system/internal/pkg/custom_errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// CustomRecoveryMiddleware handles both panic recovery and registered errors.
func CustomRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer recoverFromPanic(c) // Recover from panics

		c.Next() // Process the request

		// Handle errors from the context
		if len(c.Errors) > 0 {
			handleErrors(c)
		}
	}
}

// Handles panics and returns a structured error response
func recoverFromPanic(c *gin.Context) {
	if err := recover(); err != nil {
		log.Errorf("Panic recovered: %v", err)

		c.Error(custom_errors.New(nil, custom_errors.InternalServerError, fmt.Sprintf("%v", err)))

		// Send the response immediately and stop further processing
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error",
			"details": "An unexpected error occurred. Please try again later.",
		})
	}
}

// Handles errors from the context and ensures consistent responses
func handleErrors(c *gin.Context) {
	// Handle only the first error to avoid duplication
	if err, ok := c.Errors[0].Err.(custom_errors.CustomError); ok {
		c.JSON(err.HTTPCode, gin.H{
			"status":  "error",
			"code":    err.Code,
			"message": err.Message,
			"details": err.Details,
		})
		return
	}

	// Fallback for generic errors
	var errorMessages []string
	for _, err := range c.Errors {
		errorMessages = append(errorMessages, err.Error())
	}
	log.Errorf("Unhandled errors: %v", errorMessages)
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  "error",
		"message": "An unexpected error occurred",
		"details": errorMessages,
	})
}
