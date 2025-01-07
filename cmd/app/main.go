// @title Appointment Management System API
// @version 1.0

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/banggok/boillerplate_architecture/internal/config/app"
	"github.com/banggok/boillerplate_architecture/internal/config/db"
	eventconfig "github.com/banggok/boillerplate_architecture/internal/config/event"
	"github.com/banggok/boillerplate_architecture/internal/config/server"
	"github.com/banggok/boillerplate_architecture/internal/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	app.Setup()
	services.Setup()

	mysqlCfg, cleanUp, err := db.Setup()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer cleanUp(mysqlCfg)

	eventconfig.Setup(mysqlCfg.Master)

	server := server.Setup(mysqlCfg)
	runServer(server)
}

func runServer(router *gin.Engine) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.AppConfig.Port),
		Handler: router,
	}

	// Signal handling
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
		log.Info("server runnig")
	}()

	<-signalCh
	log.Info("Signal received, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), app.AppConfig.GracefulShutdown)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("Server shutdown error: %v", err)
	} else {
		log.Info("Server shut down gracefully")
	}
}
