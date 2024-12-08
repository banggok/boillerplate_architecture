package middleware

import (
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CustomLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture the start time
		startTime := time.Now()

		// Process the request
		c.Next()

		// Retrieve the request ID from the context
		reqID := requestid.Get(c)

		// Log the request details after the response is finalized
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		method := c.Request.Method
		path := c.Request.URL.Path

		// Use the provided logger for structured logging
		logger.WithFields(logrus.Fields{
			"request_id":  reqID,
			"status_code": statusCode,
			"latency":     latency,
			"client_ip":   clientIP,
			"method":      method,
			"path":        path,
		}).Info("Request completed")
	}
}
