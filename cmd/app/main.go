package main

import (
	"appointment_management_system/internal/config"
	tenantHandlerRegister "appointment_management_system/internal/domain/tenant/rest/v1/handler/register"
	"appointment_management_system/internal/domain/tenant/service"
	tenantRegisterUsecase "appointment_management_system/internal/domain/tenant/usecase"
	"appointment_management_system/internal/pkg/middleware"
	"context"
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
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/ulule/limiter/v3"
	ginmiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// AppConfig contains the application configuration
var AppConfig = struct {
	DefaultPort      string
	CORSAllowOrigins string
	RateLimit        int
	GracefulShutdown int
}{
	DefaultPort:      getConfigValue("PORT", "8080"),
	CORSAllowOrigins: getConfigValue("CORS_ALLOW_ORIGINS", "https://example.com"),
	RateLimit:        getConfigValueAsInt("RATE_LIMIT", 10),
	GracefulShutdown: getConfigValueAsInt("SHUTDOWN_TIMEOUT", 5),
}

type DBConnection struct {
	master *gorm.DB
	slave  *gorm.DB
}

var conn DBConnection

func main() {
	logger := log.New()
	// Setup application configurations and dependencies
	validator, rateLimiter := setupConfiguration()

	// Initialize Gin router
	router := gin.Default()

	// Connect to the database
	var err error
	conn, err := connectMySQL()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set up routes and middlewares
	setupRoutes(router, validator, rateLimiter, conn, logger)

	// Create a server runner
	serverRunner := NewHTTPServerRunner(router, AppConfig.DefaultPort)

	// Start the server with graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	runServer(serverRunner, time.Duration(AppConfig.GracefulShutdown)*time.Second, signalCh)
}

func setupConfiguration() (*validator.Validate, *limiter.Limiter) {
	config.SetupLogging()
	config.SetTimezone()

	// Setup validator
	validator := config.SetupValidator()

	// Setup rate limiter
	rate := limiter.Rate{
		Period: time.Minute,
		Limit:  int64(AppConfig.RateLimit),
	}
	store := memory.NewStore()
	rateLimiter := limiter.New(store, rate)

	return validator, rateLimiter
}

func setupRoutes(
	router *gin.Engine,
	validator *validator.Validate,
	rateLimiter *limiter.Limiter,
	db DBConnection,
	logger *log.Logger,
) {
	setupMiddlewares(router, rateLimiter, db, logger)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	api := router.Group("/api")
	v1 := api.Group("/v1")

	setupTenantV1Routes(v1.Group("/tenants"), validator)
}

func setupTenantV1Routes(tenantGroup *gin.RouterGroup, validator *validator.Validate) {
	createTenantService := service.NewTenantCreateService()
	registerUsecase := tenantRegisterUsecase.NewTenantRegisterUsecase(createTenantService)
	registerHandler := tenantHandlerRegister.NewTenantRegisterHandler(validator, registerUsecase)

	tenantGroup.POST("/register", registerHandler.Register)
}

func setupMiddlewares(
	router *gin.Engine,
	rateLimiter *limiter.Limiter,
	db DBConnection,
	logger *log.Logger,
) {
	router.Use(requestid.New())
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{AppConfig.CORSAllowOrigins},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))
	router.Use(middleware.CustomRecoveryMiddleware())
	router.Use(middleware.ErrorHandlingMiddleware())
	router.Use(middleware.CustomLogger(logger))
	router.Use(ginmiddleware.NewMiddleware(rateLimiter))

	router.Use(middleware.TransactionMiddleware(db.master, db.slave))
}

var gormOpen = gorm.Open

func connectMySQL() (DBConnection, error) {
	conn := DBConnection{}

	host := getConfigValue("DB_MASTER_HOST", "localhost")
	port := getConfigValue("DB_MASTER_PORT", "3306")
	user := getConfigValue("DB_MASTER_USER", "user")
	password := getConfigValue("DB_MASTER_PASSWORD", "password")
	dbName := getConfigValue("DB_MASTER_NAME", "appointment_management")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	master, err := gormOpen(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return DBConnection{}, fmt.Errorf("failed to connect to Master MySQL: %w", err)
	}

	log.Infof("Connected to Master MySQL database at %s:%s", host, port)

	conn.master = master

	host = getConfigValue("DB_MASTER_HOST", "localhost")
	port = getConfigValue("DB_MASTER_PORT", "3306")
	user = getConfigValue("DB_MASTER_USER", "user")
	password = getConfigValue("DB_MASTER_PASSWORD", "password")
	dbName = getConfigValue("DB_MASTER_NAME", "appointment_management")

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	slave, err := gormOpen(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return DBConnection{}, fmt.Errorf("failed to connect to Slave MySQL: %w", err)
	}

	log.Infof("Connected to Slave MySQL database at %s:%s", host, port)

	conn.slave = slave

	return conn, nil
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
		log.Warnf("Invalid value for %s. Using fallback: %d", key, fallback)
		return fallback
	}
	return intValue
}

func runServer(runner ServerRunner, shutdownTimeout time.Duration, signals chan os.Signal) {
	serverError := make(chan error, 1)

	go func() {
		log.Infof("Server running on %s", getServerAddress(runner))
		serverError <- runner.Run()
	}()

	select {
	case <-signals:
		log.Info("Signal received, initiating graceful shutdown...")
		gracefulShutdown(runner, shutdownTimeout)
	case err := <-serverError:
		if err != nil && err != http.ErrServerClosed {
			log.Errorf("Server error: %v", err)
		}
	}
}

func gracefulShutdown(runner ServerRunner, shutdownTimeout time.Duration) {
	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := runner.Shutdown(ctx); err != nil {
		log.Panicf("Failed to gracefully shutdown: %v", err)
	}

	if conn.master != nil {
		sqlDB, err := conn.master.DB()
		if err == nil {
			if closeErr := sqlDB.Close(); closeErr != nil {
				log.Errorf("Error closing database connection: %v", closeErr)
			} else {
				log.Info("Database connection closed successfully")
			}
		} else {
			log.Errorf("Error retrieving database connection: %v", err)
		}
	}

	if conn.slave != nil {
		sqlDB, err := conn.slave.DB()
		if err == nil {
			if closeErr := sqlDB.Close(); closeErr != nil {
				log.Errorf("Error closing database connection: %v", closeErr)
			} else {
				log.Info("Database connection closed successfully")
			}
		} else {
			log.Errorf("Error retrieving database connection: %v", err)
		}
	}
	log.Info("Server exited gracefully")
}

func getServerAddress(runner ServerRunner) string {
	if r, ok := runner.(*HTTPServerRunner); ok {
		return r.server.Addr
	}
	return "unknown address"
}
