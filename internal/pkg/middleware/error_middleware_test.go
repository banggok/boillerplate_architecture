package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandlingMiddleware_CustomError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandlingMiddleware())

	router.GET("/test-custom-error", func(c *gin.Context) {
		// Add a CustomError to the Gin context
		c.Error(&CustomError{
			Status:  http.StatusBadRequest,
			Code:    "BAD_REQUEST",
			Message: "This is a custom error",
			Details: map[string]string{"field": "value"},
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-custom-error", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{
		"status": "error",
		"code": "BAD_REQUEST",
		"message": "This is a custom error",
		"details": {"field": "value"}
	}`, w.Body.String(), "Response JSON does not match")
}

func TestErrorHandlingMiddleware_GenericError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandlingMiddleware())

	router.GET("/test-generic-error", func(c *gin.Context) {
		// Add a generic error to the Gin context
		c.Error(errors.New("some generic error"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-generic-error", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{
		"status": "error",
		"message": "An unexpected error occurred"
	}`, w.Body.String(), "Response JSON does not match")
}

func TestErrorHandlingMiddleware_NoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ErrorHandlingMiddleware())

	router.GET("/test-no-error", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test-no-error", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "success"}`, w.Body.String(), "Response JSON does not match")
}

func TestCustomError_Error(t *testing.T) {
	// Create a CustomError instance
	customErr := &CustomError{
		Status:  http.StatusBadRequest,
		Code:    "BAD_REQUEST",
		Message: "This is a custom error",
		Details: map[string]string{"field": "value"},
	}

	// Check the error message
	expectedMessage := "This is a custom error"
	assert.Equal(t, expectedMessage, customErr.Error(), "CustomError message does not match")
}
