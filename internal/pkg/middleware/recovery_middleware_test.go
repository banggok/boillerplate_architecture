package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCustomRecoveryMiddleware_Panic(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router with the middleware
	router := gin.New()
	router.Use(CustomRecoveryMiddleware())

	// Add a route that triggers a panic
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	// Perform the request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status code 500")
	assert.JSONEq(t, `{
		"status": "error",
		"message": "Internal server error",
		"details": "An unexpected error occurred. Please try again later."
	}`, w.Body.String(), "Expected structured error response")
}

func TestCustomRecoveryMiddleware_CustomError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router with the middleware
	router := gin.New()
	router.Use(CustomRecoveryMiddleware())

	// Add a route that returns a CustomError
	router.GET("/custom-error", func(c *gin.Context) {
		c.Error(&CustomError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid input provided",
			Details: map[string]string{"field": "email"},
		})
	})

	// Perform the request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/custom-error", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400")
	assert.JSONEq(t, `{
		"code": "INVALID_REQUEST",
		"status": "error",
		"message": "Invalid input provided",
		"details": {"field": "email"}
	}`, w.Body.String(), "Expected structured error response")
}

func TestCustomRecoveryMiddleware_GenericError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router with the middleware
	router := gin.New()
	router.Use(CustomRecoveryMiddleware())

	// Add a route that returns a generic error
	router.GET("/generic-error", func(c *gin.Context) {
		c.Error(errors.New("generic error occurred"))
	})

	// Perform the request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/generic-error", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status code 500")
	assert.JSONEq(t, `{
		"status": "error",
		"message": "An unexpected error occurred",
		"details": ["generic error occurred"]
	}`, w.Body.String(), "Expected structured error response")
}
