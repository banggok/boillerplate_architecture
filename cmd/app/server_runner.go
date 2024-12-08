package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServerRunner interface for running and shutting down a server
type ServerRunner interface {
	Run() error
	Shutdown(ctx context.Context) error
}

// HTTPServerRunner implements ServerRunner for HTTP servers
type HTTPServerRunner struct {
	server *http.Server
}

// NewHTTPServerRunner creates a new HTTPServerRunner
func NewHTTPServerRunner(router *gin.Engine, port string) *HTTPServerRunner {
	return &HTTPServerRunner{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: router,
		},
	}
}

// Run starts the HTTP server
func (r *HTTPServerRunner) Run() error {
	return r.server.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (r *HTTPServerRunner) Shutdown(ctx context.Context) error {
	return r.server.Shutdown(ctx)
}
