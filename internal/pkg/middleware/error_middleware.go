package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// CustomError represents a structured application error
type CustomError struct {
	Status  int               `json:"-"`
	Code    string            `json:"code"`    // Custom error code
	Message string            `json:"message"` // Human-readable message
	Details map[string]string `json:"details"` // Additional details
}

// Error implements the error interface
func (e *CustomError) Error() string {
	return e.Message
}

// ErrorHandlingMiddleware handles global errors.
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process the request

		// If there are errors, handle them
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				switch typedErr := err.Err.(type) {
				case *CustomError:
					// Handle CustomError
					c.JSON(typedErr.Status, gin.H{
						"status":  "error",
						"code":    typedErr.Code,
						"message": typedErr.Message,
						"details": typedErr.Details,
					})
					return
				default:
					// Log and handle generic errors
					log.Errorf("Unhandled error: %v", err)
				}
			}

			// Fallback for generic errors
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "An unexpected error occurred",
			})
		}
	}
}
