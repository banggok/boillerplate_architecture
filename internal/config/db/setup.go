package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/banggok/boillerplate_architecture/internal/config/app"
	log "github.com/sirupsen/logrus"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConnection struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

var dialector map[app.SupportedDriver]func(string) gorm.Dialector = map[app.SupportedDriver]func(string) gorm.Dialector{
	app.MySQLDriver:  mysql.Open,
	app.SQLiteDriver: sqlite.Open,
}

var logLevel map[app.Environment]logger.LogLevel = map[app.Environment]logger.LogLevel{
	app.ENV_DEV:     logger.Info,
	app.ENV_PROD:    logger.Silent,
	app.ENV_TESTING: logger.Info,
}

func Setup() (*DBConnection, func(*DBConnection), error) {
	// Connect to the master database
	dbCfg := app.AppConfig.DBConfig
	dbLogger := logger.New(
		log.New(),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel[app.AppConfig.Environment],
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	driver := dialector[dbCfg.Driver]
	master, err := gorm.Open(driver(dbCfg.MasterDSN), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		return nil, nil, err // coverage:ignore-line
	}

	// Configure the master connection pool
	sqlDB, err := master.DB()
	if err != nil {
		return nil, nil, err // coverage:ignore-line
	}
	configureConnectionPool(sqlDB, dbCfg.PoolConfig.Master, "master")

	slave := master

	// Connect to the slave database
	slave, err = gorm.Open(driver(dbCfg.SlaveDSN), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to slave database: %w", err) // coverage:ignore-line
	}

	// Configure the slave connection pool
	sqlDB, err = slave.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get slave database connection pool: %w", err) // coverage:ignore-line
	}
	configureConnectionPool(sqlDB, dbCfg.PoolConfig.Slave, "slave")

	return &DBConnection{Master: master, Slave: slave}, cleanUp, nil
}

func configureConnectionPool(sqlDB *sql.DB, cfg app.PoolConfig, label string) {
	// Apply the connection pool configurations
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)                                // Maximum number of open connections
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)                                // Maximum number of idle connections
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Minute) // Connection max lifetime

	log.Infof("Database connection pool configured for %s: maxOpenConns=%d, maxIdleConns=%d, connMaxLifetime=%dm",
		label, cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.MaxLifetime)
}

func cleanUp(conn *DBConnection) {
	closeConnection := func(db *gorm.DB, label string) {
		if db == nil {
			return
		}
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Errorf("Failed to close %s database: %v", label, err) // coverage:ignore-line
			} else {
				log.Infof("%s database closed successfully", label)
			}
		} else {
			log.Errorf("Error retrieving %s database connection: %v", label, err) // coverage:ignore-line
		}
	}
	closeConnection(conn.Master, "master")
	closeConnection(conn.Slave, "slave")
}
