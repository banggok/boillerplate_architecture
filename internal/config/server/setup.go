package server

import (
	"net/http"
	"time"

	"github.com/banggok/boillerplate_architecture/internal/config/app"
	"github.com/banggok/boillerplate_architecture/internal/config/db"
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest"
	"github.com/banggok/boillerplate_architecture/internal/pkg/middleware/recovery"
	"github.com/banggok/boillerplate_architecture/internal/pkg/middleware/transaction"
	"github.com/banggok/boillerplate_architecture/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	log_middleware "github.com/banggok/boillerplate_architecture/internal/pkg/middleware/log"

	"github.com/ulule/limiter/v3"
	ginmiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func Setup(appCfg app.Config, mysqlCfg *db.DBConnection) *gin.Engine {
	if appCfg.Environment == app.ENV_PROD {
		gin.SetMode(gin.ReleaseMode)
	}
	server := gin.Default()

	setupMiddleware(server, appCfg, mysqlCfg)

	setupRoutes(server, appCfg)

	return server
}

func setupRoutes(server *gin.Engine, cfg app.Config) {
	// Health route
	server.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	serviceCfg := services.Setup(cfg)
	rest.RegisterRoutes(server, serviceCfg)
}

func setupMiddleware(router *gin.Engine, appCfg app.Config, db *db.DBConnection) {
	router.Use(requestid.New())
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{appCfg.CORSAllowOrigins},
		AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))
	router.Use(log_middleware.CustomLogger())

	router.Use(ginmiddleware.NewMiddleware(setupRateLimiter(appCfg.RateLimit)))
	router.Use(transaction.CustomTransaction(db.Master, db.Slave))
	router.Use(recovery.CustomRecoveryMiddleware())
}

func setupRateLimiter(rateLimit int) *limiter.Limiter {
	rate := limiter.Rate{
		Period: time.Minute,
		Limit:  int64(rateLimit),
	}
	store := memory.NewStore()
	return limiter.New(store, rate)
}
