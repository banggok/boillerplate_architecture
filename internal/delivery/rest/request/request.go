package request

import (
	"fmt"

	"github.com/banggok/boillerplate_architecture/internal/config/validator"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	govalidator "github.com/go-playground/validator/v10"
)

type IRequest interface {
	ParseAndValidateRequest(c *gin.Context, binding binding.Binding) error
}

type RequestParam struct {
	Context                  *gin.Context
	Request                  interface{}
	Binding                  binding.Binding
	FieldMessages            map[string]string
	FieldKeys                map[string]string
	BadRequestErrorCode      custom_errors.ErrorCode
	UnprocessEntityErrorCode custom_errors.ErrorCode
}

type Base struct{}

// customValidationMessages provides default validation error messages.
func (b *Base) customValidationMessages(err error, fieldMessages map[string]string, fieldKeys map[string]string) map[string]string {
	validationErrors := err.(govalidator.ValidationErrors)
	errorMessages := make(map[string]string)

	for _, fieldError := range validationErrors {
		field := fieldError.StructNamespace()
		message := b.mapValidationMessage(field, fieldMessages)
		if message != "" {
			key := b.formatFieldKey(field, fieldKeys)
			errorMessages[key] = message
		}
	}

	return errorMessages
}

// ParseAndValidateRequest handles parsing and validating a request.
func (b *Base) ParseAndValidateRequest(param RequestParam) error {
	// Parse the incoming JSON request
	if param.Binding != nil {
		if err := param.Context.ShouldBindWith(param.Request, param.Binding); err != nil {
			param.Context.Error(custom_errors.New(
				err,
				param.BadRequestErrorCode,
				"invalid request"))
			return err
		}
	}

	// Validate the request using the validator
	if err := validator.Validate.Struct(param.Request); err != nil {
		validationErrors := b.customValidationMessages(err, param.FieldMessages, param.FieldKeys)
		param.Context.Error(custom_errors.New(
			err,
			param.UnprocessEntityErrorCode,
			"failed to validate request",
			validationErrors))
		return err
	}

	return nil
}

// mapValidationMessage maps field names to error messages.
func (b *Base) mapValidationMessage(field string, fieldMessages map[string]string) string {
	if message, exists := fieldMessages[field]; exists {
		return message
	}
	return fmt.Sprintf("Invalid value for %s.", field)
}

// formatFieldKey maps field names to JSON keys.
func (b *Base) formatFieldKey(field string, fieldKeys map[string]string) string {
	if key, exists := fieldKeys[field]; exists {
		return key
	}
	return field
}
