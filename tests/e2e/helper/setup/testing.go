package setup

import (
	"fmt"
	"testing"
	"time"

	"github.com/banggok/boillerplate_architecture/internal/config/app"
	"github.com/banggok/boillerplate_architecture/internal/config/db"
	"github.com/banggok/boillerplate_architecture/internal/config/server"
	"github.com/banggok/boillerplate_architecture/internal/data/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestingEnv(t *testing.T, useServer bool) (serverEngine *gin.Engine, cleanUp func(*db.DBConnection), mysqlCfg *db.DBConnection) {
	// Use a unique in-memory SQLite database for this test
	uniqueID := time.Now().UnixNano() // Generate a unique ID based on the current timestamp
	dsn := fmt.Sprintf("file:test_%d?cache=shared&mode=memory", uniqueID)

	// Configure the application to use this DSN
	app.Setup()
	appConfig := app.AppConfig
	appConfig.DBConfig.MasterDSN = dsn
	appConfig.DBConfig.SlaveDSN = dsn
	appConfig.DBConfig.Driver = app.SQLiteDriver
	appConfig.DBConfig.PoolConfig.Master.MaxIdleConns = 1
	appConfig.DBConfig.PoolConfig.Master.MaxOpenConns = 1
	appConfig.DBConfig.PoolConfig.Slave.MaxIdleConns = 1
	appConfig.DBConfig.PoolConfig.Slave.MaxOpenConns = 1
	appConfig.Environment = app.ENV_TESTING

	// Set up the database
	mysqlCfg, cleanUp, err := db.Setup(appConfig)
	require.NoError(t, err)

	mysqlCfg.Master.AutoMigrate(
		model.Tenant{},
		model.Account{})

	// Set up the Gin server
	if useServer {
		serverEngine = server.Setup(appConfig, mysqlCfg)
	}

	return serverEngine, cleanUp, mysqlCfg
}
