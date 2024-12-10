package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestCustomLogger(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new logrus test hook to capture logs
	hook := test.NewGlobal()

	// Create a Gin router with the CustomLogger middleware
	router := gin.New()
	router.Use(requestid.New()) // Add request ID middleware
	router.Use(CustomLogger())

	// Add a sample route
	router.GET("/test", func(c *gin.Context) {
		time.Sleep(50 * time.Millisecond) // Simulate some latency
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create a request to the sample route with a simulated client IP
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345" // Inject a client IP directly

	// Serve the request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")
	assert.JSONEq(t, `{"message":"success"}`, w.Body.String(), "Response body does not match")

	// Check that the log entry was created
	assert.Len(t, hook.AllEntries(), 1, "Expected one log entry")
	entry := hook.LastEntry()
	assert.Equal(t, logrus.InfoLevel, entry.Level, "Expected log level to be info")
	assert.Equal(t, "Request completed", entry.Message, "Expected log message to match")

	// Check log fields
	fields := entry.Data
	assert.Equal(t, http.StatusOK, fields["status_code"], "Log field 'status_code' does not match")
	assert.Equal(t, "GET", fields["method"], "Log field 'method' does not match")
	assert.Equal(t, "/test", fields["path"], "Log field 'path' does not match")
	assert.NotEmpty(t, fields["latency"], "Log field 'latency' should not be empty")
	assert.Equal(t, "192.168.1.1", fields["client_ip"], "Log field 'client_ip' does not match expected IP")
	assert.NotEmpty(t, fields["request_id"], "Log field 'request_id' should not be empty")
}
