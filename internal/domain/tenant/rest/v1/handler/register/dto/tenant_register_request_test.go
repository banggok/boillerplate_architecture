package dto

import (
	"appointment_management_system/internal/config"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var validate *validator.Validate

func TestMain(m *testing.M) {
	// Initialize the validator
	validate = config.SetupValidator()
	m.Run()
}

func TestRegisterTenantRequestValidation(t *testing.T) {

	t.Run("Valid Payload", func(t *testing.T) {
		validRequest := getValidTenantRequest()
		err := validate.Struct(validRequest)
		assert.NoError(t, err, "Valid payload should pass validation")
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		invalidRequest := getMissingFieldsTenantRequest()
		checkValidationErrors(t, validate, invalidRequest, map[string]string{
			"name":         "Name is required.",
			"account.name": "Account name is required and can only contain letters and spaces.",
			"timezone":     "Timezone is required and must be a valid IANA timezone.",
		})
	})

	t.Run("Invalid Field Formats", func(t *testing.T) {
		invalidRequest := getInvalidFormatTenantRequest()
		checkValidationErrors(t, validate, invalidRequest, map[string]string{
			"email":         "Email is required and must be in a valid format.",
			"phone":         "Phone must be in a valid international format (E.164).",
			"timezone":      "Timezone is required and must be a valid IANA timezone.",
			"opening_hours": "Opening hours must be in the format HH:mm.",
			"closing_hours": "Closing hours must be in the format HH:mm.",
			"account.name":  "Account name is required and can only contain letters and spaces.",
			"account.email": "Account email is required and must be in a valid format.",
			"account.phone": "Account phone must be in a valid international format (E.164).",
		})
	})

	t.Run("Invalid Address Length", func(t *testing.T) {
		invalidRequest := getInvalidAddressTenantRequest()
		checkValidationErrors(t, validate, invalidRequest, map[string]string{
			"address": "Address must not exceed 255 characters.",
		})
	})
}

// Helper Functions

func getValidTenantRequest() RegisterTenantRequest {
	return RegisterTenantRequest{
		Name:         "Acme Corporation",
		Address:      "123 Main St",
		Email:        "contact@acme.com",
		Phone:        "+1234567890",
		Timezone:     "America/New_York",
		OpeningHours: "09:00",
		ClosingHours: "18:00",
		Account: Account{
			Name:  "John Doe",
			Email: "admin@acme.com",
			Phone: "+0987654321",
		},
	}
}

func getMissingFieldsTenantRequest() RegisterTenantRequest {
	return RegisterTenantRequest{
		Email: "contact@acme.com", // Missing required fields
	}
}

func getInvalidFormatTenantRequest() RegisterTenantRequest {
	return RegisterTenantRequest{
		Name:         "Acme Corporation",
		Email:        "invalid-email",
		Phone:        "123456",
		Timezone:     "Invalid/Timezone",
		OpeningHours: "25:00",
		ClosingHours: "12:60",
		Account: Account{
			Name:  "John@Doe",
			Email: "admin@invalid-email",
			Phone: "not-a-phone",
		},
	}
}

func getInvalidAddressTenantRequest() RegisterTenantRequest {
	return RegisterTenantRequest{
		Name:    "Acme Corporation",
		Address: strings.Repeat("A", 256), // Exceeds 255 characters
		Email:   "contact@acme.com",
		Account: Account{
			Name:  "John Doe",
			Email: "admin@acme.com",
		},
	}
}

func checkValidationErrors(
	t *testing.T,
	validate *validator.Validate,
	request RegisterTenantRequest,
	expectedErrors map[string]string,
) {
	err := validate.Struct(request)
	assert.Error(t, err, "Validation should fail for invalid input")

	validationErrors := request.CustomValidationMessages(err)
	for field, expectedMessage := range expectedErrors {
		assert.Contains(t, validationErrors, field, "Expected validation error for field: "+field)
		assert.Equal(t, expectedMessage, validationErrors[field], "Validation error message mismatch for field: "+field)
	}
}

func TestMapValidationMessage(t *testing.T) {
	tests := []struct {
		field    string
		expected string
	}{
		{"RegisterTenantRequest.Name", "Name is required."},
		{"RegisterTenantRequest.InvalidField", "Invalid value for RegisterTenantRequest.InvalidField."},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			result := mapValidationMessage(tt.field)
			if result != tt.expected {
				t.Errorf("mapValidationMessage(%q) = %q; want %q", tt.field, result, tt.expected)
			}
		})
	}
}

func TestFormatFieldKey(t *testing.T) {
	tests := []struct {
		field    string
		expected string
	}{
		{"RegisterTenantRequest.Name", "name"},
		{"RegisterTenantRequest.Account.Email", "account.email"},
		{"RegisterTenantRequest.InvalidField", "RegisterTenantRequest.InvalidField"},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			result := formatFieldKey(tt.field)
			if result != tt.expected {
				t.Errorf("formatFieldKey(%q) = %q; want %q", tt.field, result, tt.expected)
			}
		})
	}
}
