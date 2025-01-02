package validator

import (
	"log"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

// Setup initializes the validator and registers custom rules.
func Setup() {
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

	Validate = validate
}
