package register

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// RegisterTenantRequest represents the payload for tenant registration
// @description Tenant registration request body
type RegisterTenantRequest struct {
	Name         string                       `json:"name" validate:"required" example:"Example Tenant"`               // Tenant's name
	Address      string                       `json:"address" validate:"omitempty,max=255" example:"123 Main Street"`  // Tenant's address
	Email        string                       `json:"email" validate:"required,email" example:"tenant@example.com"`    // Tenant's email
	Phone        string                       `json:"phone" validate:"required,e164" example:"+1987654321"`            // Tenant's phone
	Timezone     string                       `json:"timezone" validate:"required,iana_tz" example:"America/New_York"` // Tenant's timezone
	OpeningHours string                       `json:"opening_hours" validate:"required,time_format" example:"09:00"`   // Opening hours
	ClosingHours string                       `json:"closing_hours" validate:"required,time_format" example:"18:00"`   // Closing hours
	Account      RegisterTenantAccountRequest `json:"account" validate:"required"`                                     // Admin user details
}

type RegisterTenantAccountRequest struct {
	Name  string `json:"name" validate:"required,alpha_space" example:"John Doe"`     // Admin's name
	Email string `json:"email" validate:"required,email" example:"admin@example.com"` // Admin's email
	Phone string `json:"phone" validate:"required,e164" example:"+1234567890"`        // Admin's phone
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
		"RegisterTenantRequest.Phone":         "Phone is required and must be in a valid international format (E.164).",
		"RegisterTenantRequest.Timezone":      "Timezone is required and must be a valid IANA timezone.",
		"RegisterTenantRequest.OpeningHours":  "Opening hours is required and must be in the format HH:mm.",
		"RegisterTenantRequest.ClosingHours":  "Closing hours is required and must be in the format HH:mm.",
		"RegisterTenantRequest.Account.Name":  "Account name is required and can only contain letters and spaces.",
		"RegisterTenantRequest.Account.Email": "Account email is required and must be in a valid format.",
		"RegisterTenantRequest.Account.Phone": "Account phone is required and must be in a valid international format (E.164).",
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
