package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// RegisterTenantRequest represents the payload for tenant registration.
type RegisterTenantRequest struct {
	Name         string  `json:"name" validate:"required"`                       // Tenant's name (required, no strict restrictions)
	Address      string  `json:"address" validate:"omitempty,max=255"`           // Optional, with a maximum of 255 characters
	Email        string  `json:"email" validate:"required,email"`                // Tenant's contact email (required)
	Phone        string  `json:"phone" validate:"omitempty,e164"`                // Tenant's phone number (optional)
	Timezone     string  `json:"timezone" validate:"required,iana_tz"`           // Tenant's timezone (required, IANA format)
	OpeningHours string  `json:"opening_hours" validate:"omitempty,time_format"` // Tenant's opening hours in HH:mm format (optional)
	ClosingHours string  `json:"closing_hours" validate:"omitempty,time_format"` // Tenant's closing hours in HH:mm format (optional)
	Account      Account `json:"account" validate:"required"`                    // Admin user details (required)
}

// Account represents the admin user's information.
type Account struct {
	Name  string `json:"name" validate:"required,alpha_space"` // Admin user's name (required, only letters and spaces allowed)
	Email string `json:"email" validate:"required,email"`      // Admin user's email (required)
	Phone string `json:"phone" validate:"omitempty,e164"`      // Admin user's phone number (optional)
}

// CustomValidationMessages maps validation errors to user-friendly messages.
func (r *RegisterTenantRequest) CustomValidationMessages(err error) map[string]string {
	validationErrors := err.(validator.ValidationErrors)
	errorMessages := make(map[string]string)

	for _, fieldError := range validationErrors {
		field := fieldError.StructNamespace() // Use StructNamespace for nested fields
		message := mapValidationMessage(field)
		if message != "" {
			key := formatFieldKey(field)
			errorMessages[key] = message
		}
	}

	return errorMessages
}

// mapValidationMessage maps a field name to a user-friendly error message.
func mapValidationMessage(field string) string {
	fieldMessages := map[string]string{
		"RegisterTenantRequest.Name":          "Name is required.",
		"RegisterTenantRequest.Address":       "Address must not exceed 255 characters.",
		"RegisterTenantRequest.Email":         "Email is required and must be in a valid format.",
		"RegisterTenantRequest.Phone":         "Phone must be in a valid international format (E.164).",
		"RegisterTenantRequest.Timezone":      "Timezone is required and must be a valid IANA timezone.",
		"RegisterTenantRequest.OpeningHours":  "Opening hours must be in the format HH:mm.",
		"RegisterTenantRequest.ClosingHours":  "Closing hours must be in the format HH:mm.",
		"RegisterTenantRequest.Account.Name":  "Account name is required and can only contain letters and spaces.",
		"RegisterTenantRequest.Account.Email": "Account email is required and must be in a valid format.",
		"RegisterTenantRequest.Account.Phone": "Account phone must be in a valid international format (E.164).",
	}
	if message, exists := fieldMessages[field]; exists {
		return message
	}
	return fmt.Sprintf("Invalid value for %s.", field)
}

// formatFieldKey formats the field key for use in the errorMessages map.
func formatFieldKey(field string) string {
	fieldKeys := map[string]string{
		"RegisterTenantRequest.Name":          "name",
		"RegisterTenantRequest.Address":       "address",
		"RegisterTenantRequest.Email":         "email",
		"RegisterTenantRequest.Phone":         "phone",
		"RegisterTenantRequest.Timezone":      "timezone",
		"RegisterTenantRequest.OpeningHours":  "opening_hours",
		"RegisterTenantRequest.ClosingHours":  "closing_hours",
		"RegisterTenantRequest.Account.Name":  "account.name",
		"RegisterTenantRequest.Account.Email": "account.email",
		"RegisterTenantRequest.Account.Phone": "account.phone",
	}
	if key, exists := fieldKeys[field]; exists {
		return key
	}
	return field
}
