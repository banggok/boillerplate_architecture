package validator

import (
	"log"
	"time"
	"unicode"

	valueobject "github.com/banggok/boillerplate_architecture/internal/data/entity/value_object"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

// Setup initializes the validator and registers custom rules.
func init() {
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

	// Register verification_type validation rule
	if err := validate.RegisterValidation("verification_type", func(fl validator.FieldLevel) bool {
		// Convert the field value to a VerificationType
		verificationType := valueobject.VerificationType(fl.Field().String())
		// Check if the verification type is valid
		return verificationType.IsValid()
	}); err != nil {
		log.Fatalf("Failed to register verification_type validation: %v", err)
	}

	Validate = validate
}
