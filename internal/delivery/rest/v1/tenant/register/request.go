package register

import (
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/request"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Request represents the payload for tenant registration
// @description Tenant registration request body
type Request struct {
	request.Base
	Name         string  `json:"name" validate:"required" example:"Example Tenant"`               // Tenant's name
	Address      string  `json:"address" validate:"omitempty,max=255" example:"123 Main Street"`  // Tenant's address
	Email        string  `json:"email" validate:"required,email" example:"tenant@example.com"`    // Tenant's email
	Phone        string  `json:"phone" validate:"required,e164" example:"+1987654321"`            // Tenant's phone
	Timezone     string  `json:"timezone" validate:"required,iana_tz" example:"America/New_York"` // Tenant's timezone
	OpeningHours string  `json:"opening_hours" validate:"required,time_format" example:"09:00"`   // Opening hours
	ClosingHours string  `json:"closing_hours" validate:"required,time_format" example:"18:00"`   // Closing hours
	Account      Account `json:"account" validate:"required"`                                     // Admin user details
}

var _ request.IRequest = &Request{}

type Account struct {
	Name  string `json:"name" validate:"required,alpha_space" example:"John Doe"`     // Admin's name
	Email string `json:"email" validate:"required,email" example:"admin@example.com"` // Admin's email
	Phone string `json:"phone" validate:"required,e164" example:"+1234567890"`        // Admin's phone
}

var fieldMessages = map[string]string{
	"Request.Name":          "Name is required.",
	"Request.Address":       "Address must not exceed 255 characters.",
	"Request.Email":         "Email is required and must be in a valid format.",
	"Request.Phone":         "Phone is required and must be in a valid international format (E.164).",
	"Request.Timezone":      "Timezone is required and must be a valid IANA timezone.",
	"Request.OpeningHours":  "Opening hours is required and must be in the format HH:mm.",
	"Request.ClosingHours":  "Closing hours is required and must be in the format HH:mm.",
	"Request.Account.Name":  "Account name is required and can only contain letters and spaces.",
	"Request.Account.Email": "Account email is required and must be in a valid format.",
	"Request.Account.Phone": "Account phone is required and must be in a valid international format (E.164).",
}

var fieldKeys = map[string]string{
	"Request.Name":          "name",
	"Request.Address":       "address",
	"Request.Email":         "email",
	"Request.Phone":         "phone",
	"Request.Timezone":      "timezone",
	"Request.OpeningHours":  "opening_hours",
	"Request.ClosingHours":  "closing_hours",
	"Request.Account.Name":  "account.name",
	"Request.Account.Email": "account.email",
	"Request.Account.Phone": "account.phone",
}

// ParseAndValidateRequest implements request.IRequest.
// Subtle: this method shadows the method (BaseRequest).ParseAndValidateRequest of Request.BaseRequest.
func (r *Request) ParseAndValidateRequest(c *gin.Context) error {
	return r.Base.ParseAndValidateRequest(request.RequestParam{
		Context:                  c,
		Request:                  r,
		Binding:                  binding.JSON,
		FieldMessages:            fieldMessages,
		FieldKeys:                fieldKeys,
		BadRequestErrorCode:      custom_errors.TenantBadRequest,
		UnprocessEntityErrorCode: custom_errors.TenantUnprocessEntity,
	})
}
