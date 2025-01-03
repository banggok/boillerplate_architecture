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
	app.AppConfig.DBConfig.MasterDSN = dsn
	app.AppConfig.DBConfig.SlaveDSN = dsn
	app.AppConfig.DBConfig.Driver = app.SQLiteDriver
	app.AppConfig.DBConfig.PoolConfig.Master.MaxIdleConns = 1
	app.AppConfig.DBConfig.PoolConfig.Master.MaxOpenConns = 1
	app.AppConfig.DBConfig.PoolConfig.Slave.MaxIdleConns = 1
	app.AppConfig.DBConfig.PoolConfig.Slave.MaxOpenConns = 1
	app.AppConfig.Environment = app.ENV_TESTING

	// Set up the database
	mysqlCfg, cleanUp, err := db.Setup()
	require.NoError(t, err)

	mysqlCfg.Master.AutoMigrate(
		model.Tenant{},
		model.Account{},
		model.AccountVerification{},
	)

	// Set up the Gin server
	if useServer {
		serverEngine = server.Setup(mysqlCfg)
	}

	return serverEngine, cleanUp, mysqlCfg
}
