// @title Appointment Management System API
// @version 1.0

package main

import (
	"appointment_management_system/internal/config"
	"appointment_management_system/internal/domain/tenants/delivery"
	"appointment_management_system/internal/pkg/custom_errors"
	log_middleware "appointment_management_system/internal/pkg/middleware/log"
	"appointment_management_system/internal/pkg/middleware/recovery"
	"appointment_management_system/internal/pkg/middleware/transaction"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/ulule/limiter/v3"
	ginmiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AppConfig struct {
	Port             string
	CORSAllowOrigins string
	RateLimit        int
	GracefulShutdown time.Duration
	DBConfig         DBConfig
}

type DBConfig struct {
	MasterDSN string
	SlaveDSN  string
}

type DBConnection struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func main() {
	config.SetupLogging()
	config.SetTimezone()
	appConfig := loadAppConfig()
	rateLimiter := setupRateLimiter(appConfig.RateLimit)

	dbConn, err := connectDatabase(appConfig.DBConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer cleanupDatabase(dbConn)

	router := setupRouter(dbConn, rateLimiter)
	runServer(router, appConfig)
}

func loadAppConfig() AppConfig {
	return AppConfig{
		Port:             getConfigValue("PORT", "8080"),
		CORSAllowOrigins: getConfigValue("CORS_ALLOW_ORIGINS", "https://example.com"),
		RateLimit:        getConfigValueAsInt("RATE_LIMIT", 10),
		GracefulShutdown: time.Duration(getConfigValueAsInt("SHUTDOWN_TIMEOUT", 5)) * time.Second,
		DBConfig: DBConfig{
			MasterDSN: getDSN("MASTER"),
			SlaveDSN:  getDSN("SLAVE"),
		},
	}
}

func getDSN(role string) string {
	host := getConfigValue(fmt.Sprintf("DB_%s_HOST", role), "localhost")
	port := getConfigValue(fmt.Sprintf("DB_%s_PORT", role), "3306")
	user := getConfigValue(fmt.Sprintf("DB_%s_USER", role), "user")
	password := getConfigValue(fmt.Sprintf("DB_%s_PASSWORD", role), "password")
	dbName := getConfigValue(fmt.Sprintf("DB_%s_NAME", role), "appointment_management")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
}
func connectDatabase(cfg DBConfig) (*DBConnection, error) {
	// Connect to the master database
	dbLogger := logger.New(
		log.New(),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	master, err := gorm.Open(mysql.Open(cfg.MasterDSN), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to connect database")
	}

	// Configure the master connection pool
	sqlDB, err := master.DB()
	if err != nil {
		return nil, custom_errors.New(
			err, custom_errors.InternalServerError, "failed to get master database connection pool")
	}
	configureConnectionPool(sqlDB, "master")

	// Connect to the slave database
	slave, err := gorm.Open(mysql.Open(cfg.SlaveDSN), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to slave database: %w", err)
	}

	// Configure the slave connection pool
	sqlDB, err = slave.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get slave database connection pool: %w", err)
	}
	configureConnectionPool(sqlDB, "slave")

	return &DBConnection{Master: master, Slave: slave}, nil
}

func configureConnectionPool(sqlDB *sql.DB, label string) {
	sqlDB.SetMaxOpenConns(50)                  // Maximum number of open connections to the database
	sqlDB.SetMaxIdleConns(25)                  // Maximum number of idle connections in the pool
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Maximum amount of time a connection may be reused

	log.Infof("Database connection pool configured for %s", label)
}

func cleanupDatabase(conn *DBConnection) {
	closeConnection := func(db *gorm.DB, label string) {
		if db == nil {
			return
		}
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Errorf("Failed to close %s database: %v", label, err)
			} else {
				log.Infof("%s database closed successfully", label)
			}
		} else {
			log.Errorf("Error retrieving %s database connection: %v", label, err)
		}
	}
	closeConnection(conn.Master, "master")
	closeConnection(conn.Slave, "slave")
}

func setupRouter(
	db *DBConnection,
	rateLimiter *limiter.Limiter,
) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(requestid.New())
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{getConfigValue("CORS_ALLOW_ORIGINS", "*")},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))
	router.Use(log_middleware.CustomLogger())
	router.Use(ginmiddleware.NewMiddleware(rateLimiter))
	router.Use(transaction.CustomTransaction(db.Master, db.Slave))
	router.Use(recovery.CustomRecoveryMiddleware())

	// Health route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// API Routes
	api := router.Group("/api")
	setupRoutes(api)

	return router
}

func setupRoutes(api *gin.RouterGroup) {
	delivery.RegisterTenantRoutes(api)
}

func runServer(router *gin.Engine, cfg AppConfig) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	// Signal handling
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-signalCh
	log.Info("Signal received, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdown)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("Server shutdown error: %v", err)
	} else {
		log.Info("Server shut down gracefully")
	}
}

func getConfigValue(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getConfigValueAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Warnf("Invalid value for %s, using fallback: %d", key, fallback)
		return fallback
	}
	return intValue
}

func setupRateLimiter(rateLimit int) *limiter.Limiter {
	rate := limiter.Rate{
		Period: time.Minute,
		Limit:  int64(rateLimit),
	}
	store := memory.NewStore()
	return limiter.New(store, rate)
}
