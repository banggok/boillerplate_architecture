package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// CustomRecoveryMiddleware handles both global panics and registered errors.
func CustomRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer recoverFromPanic(c)

		c.Next() // Continue to the next middleware or handler

		if len(c.Errors) > 0 {
			handleErrors(c)
		}
	}
}

func recoverFromPanic(c *gin.Context) {
	if err := recover(); err != nil {
		// Log the panic details
		log.Errorf("Panic recovered: %v", err)

		// Return a structured error response for the panic
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error",
			"details": "An unexpected error occurred. Please try again later.",
		})

		// Abort further request processing
		c.Abort()
	}
}

func handleErrors(c *gin.Context) {
	var errorMessages []string
	statusCode := http.StatusInternalServerError // Default to 500 if no specific code

	// Collect errors from the context
	for _, err := range c.Errors {
		switch typedErr := err.Err.(type) {
		case *CustomError:
			// Handle CustomError
			statusCode = typedErr.Status
			c.JSON(statusCode, gin.H{
				"status":  "error",
				"message": typedErr.Message,
				"details": typedErr.Details,
			})
			return
		default:
			// Log generic errors
			errorMessages = append(errorMessages, err.Error())
		}
	}

	// Fallback response for generic errors
	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": "An unexpected error occurred",
		"details": errorMessages,
	})
}
