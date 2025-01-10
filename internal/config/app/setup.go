package app

import (
	"fmt"
	"time"

	"github.com/banggok/boillerplate_architecture/internal/config"
)

var AppConfig appConfig

var getConfigValue = config.GetConfigValue
var getConfigValueAsInt = config.GetConfigValueAsInt

type Environment string

const (
	ENV_DEV     Environment = "development"
	ENV_PROD    Environment = "production"
	ENV_TESTING Environment = "testing"
)

type secretKey struct {
	Access         string
	Refresh        string
	AccessExpired  time.Duration
	RefreshExpired time.Duration
}

type appConfig struct {
	BaseUrl          string
	Port             string
	CORSAllowOrigins string
	RateLimit        int
	GracefulShutdown time.Duration
	DBConfig         DBConfig
	Environment      Environment
	ExpiredDuration  ExpiredDuration
	SecretKey        secretKey
}

type ExpiredDuration struct {
	EmailVerification         time.Duration
	ResetPasswordVerification time.Duration
}

type EnvPoolConfig struct {
	Master PoolConfig
	Slave  PoolConfig
}

type PoolConfig struct {
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
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
	environment := getConfigValue("ENVIRONMENT", string(ENV_DEV))
	AppConfig = appConfig{
		BaseUrl:          getConfigValue("BASE_URL", "http://localhost"),
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
		ExpiredDuration: ExpiredDuration{
			EmailVerification:         time.Duration(getConfigValueAsInt("EMAIL_VERIFICATION_EXPIRED", 24)) * time.Hour,
			ResetPasswordVerification: time.Duration(getConfigValueAsInt("RESET_PASSWORD_EXPIRED", 90)) * 24 * time.Hour,
		},
		SecretKey: secretKey{
			Access:         getConfigValue("ACCESS_SECRET_KEY", ""),
			Refresh:        getConfigValue("REFRESH_SECRET_KEY", ""),
			AccessExpired:  time.Duration(getConfigValueAsInt("ACCESS_EXPIRED", 15)) * time.Minute,
			RefreshExpired: time.Duration(getConfigValueAsInt("REFRESH_EXPIRED", 7)) * 24 * time.Hour,
		},
	}

}

func getPoolConfig(label string) PoolConfig {
	// Fetch connection pool settings from environment variables
	return PoolConfig{
		MaxOpenConns: getConfigValueAsInt(fmt.Sprintf("DB_%s_MAX_OPEN_CONNS", label), 100),
		MaxIdleConns: getConfigValueAsInt(fmt.Sprintf("DB_%s_MAX_IDLE_CONNS", label), 10),
		MaxLifetime:  time.Duration(getConfigValueAsInt(fmt.Sprintf("DB_%s_CONN_MAX_LIFETIME", label), 30)) * time.Minute,
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
