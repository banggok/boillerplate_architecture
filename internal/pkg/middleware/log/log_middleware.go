package log

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger = logrus.StandardLogger()

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture the start time
		startTime := time.Now()

		// Retrieve query parameters
		queryParams := c.Request.URL.RawQuery

		// Retrieve path parameters
		pathParams := c.Params

		// Retrieve the request body
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// Restore the request body for further processing
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Retrieve request headers
		headers := make(map[string]string)
		for key, values := range c.Request.Header {
			headers[key] = values[0] // Assuming single-value headers
		}

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
			"request_id":   reqID,
			"status_code":  statusCode,
			"latency":      fmt.Sprintf("%dms", latency.Milliseconds()),
			"client_ip":    clientIP,
			"method":       method,
			"path":         path,
			"query_params": queryParams,
			"path_params":  pathParams,
			"headers":      headers,
			"body":         requestBody,
		}).Info("Request completed")
	}
}
