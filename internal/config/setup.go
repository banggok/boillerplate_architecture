package config

import (
	"os"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

// SetupLogging configures the logrus logging format and level.
func SetupLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}

// SetTimezone sets the application's default timezone.
func SetTimezone() {
	timezone := os.Getenv("TZ")
	if timezone == "" {
		timezone = "UTC" // Default
	}
	os.Setenv("TZ", timezone)
}

// SetupValidator initializes the validator and registers custom rules.
func SetupValidator() *validator.Validate {
	validate := validator.New()

	// Register custom validation rules
	if err := validate.RegisterValidation("time_format", func(fl validator.FieldLevel) bool {
		_, err := time.Parse("15:04", fl.Field().String())
		return err == nil
	}); err != nil {
		log.Fatalf("Failed to register time_format validation: %v", err)
	}

	if err := validate.RegisterValidation("iana_tz", func(fl validator.FieldLevel) bool {
		_, err := time.LoadLocation(fl.Field().String())
		return err == nil
	}); err != nil {
		log.Fatalf("Failed to register iana_tz validation: %v", err)
	}

	// Register alpha_space validation rule
	if err := validate.RegisterValidation("alpha_space", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		for _, char := range value {
			if !(unicode.IsLetter(char) || unicode.IsSpace(char)) {
				return false
			}
		}
		return true
	}); err != nil {
		log.Fatalf("Failed to register alpha_space validation: %v", err)
	}

	return validate
}
