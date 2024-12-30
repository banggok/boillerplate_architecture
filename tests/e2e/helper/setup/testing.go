package setup

import (
	"appointment_management_system/internal/config/app"
	"appointment_management_system/internal/config/db"
	"appointment_management_system/internal/config/server"
	"appointment_management_system/internal/data/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestingEnv(t *testing.T) (*gin.Engine, func(*db.DBConnection), *db.DBConnection) {
	// Use a unique in-memory SQLite database for this test
	dsn := "file::memory:?cache=shared"

	// Configure the application to use this DSN
	appConfig := app.Setup()
	appConfig.DBConfig.MasterDSN = dsn
	appConfig.DBConfig.SlaveDSN = dsn
	appConfig.DBConfig.Driver = app.SQLiteDriver
	appConfig.Environment = app.ENV_TESTING

	// Set up the database
	mysqlCfg, cleanUp, err := db.Setup(appConfig)
	require.NoError(t, err)

	mysqlCfg.Master.AutoMigrate(model.Tenant{}, model.Account{})
	mysqlCfg.Slave.AutoMigrate(model.Tenant{}, model.Account{})

	// Set up the Gin server
	server := server.Setup(appConfig, mysqlCfg)
	return server, cleanUp, mysqlCfg
}
