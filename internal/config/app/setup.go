package app

import (
	"fmt"
	"os"
	"time"

	"github.com/banggok/boillerplate_architecture/internal/config"
	"github.com/banggok/boillerplate_architecture/internal/config/smtp"
	"github.com/banggok/boillerplate_architecture/internal/config/validator"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var AppConfig Config

var getConfigValue = config.GetConfigValue
var getConfigValueAsInt = config.GetConfigValueAsInt

type Environment string

const (
	ENV_DEV     Environment = "development"
	ENV_PROD    Environment = "production"
	ENV_TESTING Environment = "testing"
)

type Config struct {
	Port             string
	CORSAllowOrigins string
	RateLimit        int
	GracefulShutdown time.Duration
	DBConfig         DBConfig
	Environment      Environment
	SMTP             smtp.Config
}

type EnvPoolConfig struct {
	Master PoolConfig
	Slave  PoolConfig
}

type PoolConfig struct {
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  int
}

// SupportedDriver defines the type of database driver
type SupportedDriver string

const (
	MySQLDriver  SupportedDriver = "mysql"
	SQLiteDriver SupportedDriver = "sqlite"
)

type DBConfig struct {
	MasterDSN  string
	SlaveDSN   string
	PoolConfig EnvPoolConfig
	Driver     SupportedDriver
}

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Warn("No .env file found or unable to load")
	}
}

func Setup() {
	setupLogging()
	setTimezone()
	validator.Setup()
	environment := getConfigValue("ENVIRONMENT", string(ENV_DEV))
	AppConfig = Config{
		Port:             getConfigValue("HTTP_PORT", "8080"),
		CORSAllowOrigins: getConfigValue("CORS_ALLOW_ORIGINS", "*"),
		RateLimit:        getConfigValueAsInt("RATE_LIMIT", 100),
		GracefulShutdown: time.Duration(getConfigValueAsInt("SHUTDOWN_TIMEOUT", 5)) * time.Second,
		DBConfig: DBConfig{
			MasterDSN: getDSN("MASTER"),
			SlaveDSN:  getDSN("SLAVE"),
			PoolConfig: EnvPoolConfig{
				Master: getPoolConfig("MASTER"),
				Slave:  getPoolConfig("SLAVE"),
			},
			Driver: MySQLDriver,
		},
		Environment: Environment(environment),
		SMTP:        smtp.Setup(),
	}

}

func getPoolConfig(label string) PoolConfig {
	// Fetch connection pool settings from environment variables
	return PoolConfig{
		MaxOpenConns: getConfigValueAsInt(fmt.Sprintf("DB_%s_MAX_OPEN_CONNS", label), 100),
		MaxIdleConns: getConfigValueAsInt(fmt.Sprintf("DB_%s_MAX_IDLE_CONNS", label), 10),
		MaxLifetime:  getConfigValueAsInt(fmt.Sprintf("DB_%s_CONN_MAX_LIFETIME", label), 30),
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

// SetupLogging configures the logrus logging format and level.
func setupLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}

// SetTimezone sets the application's default timezone.
func setTimezone() {
	os.Setenv("TZ", "UTC")
}
