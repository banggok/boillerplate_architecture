package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPServerRunner_Run(t *testing.T) {
	// Create a Gin engine with a test route
	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create the HTTPServerRunner
	runner := NewHTTPServerRunner(router, "8081")

	// Run the server in a goroutine
	go func() {
		err := runner.Run()
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Failed to start server: %v", err)
		}
	}()

	// Give the server some time to start
	time.Sleep(100 * time.Millisecond)

	// Make a test HTTP request
	resp, err := http.Get("http://localhost:8081/test")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Shutdown the server
	err = runner.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestHTTPServerRunner_Shutdown(t *testing.T) {
	// Create a Gin engine with a test route
	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create the HTTPServerRunner
	runner := NewHTTPServerRunner(router, "8082")

	// Run the server in a goroutine
	go func() {
		err := runner.Run()
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Failed to start server: %v", err)
		}
	}()

	// Give the server some time to start
	time.Sleep(100 * time.Millisecond)

	// Shutdown the server with a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := runner.Shutdown(ctx)
	assert.NoError(t, err)

	// Verify the server is shut down
	resp, err := http.Get("http://localhost:8082/test")
	assert.Error(t, err)
	assert.Nil(t, resp, "Response should be nil after shutdown")
}
