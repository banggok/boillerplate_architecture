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
	ParseAndValidateRequest() error
}

type Base struct {
	C                        *gin.Context
	Binding                  binding.Binding
	FieldMessages            map[string]string
	FieldKeys                map[string]string
	BadRequestErrorCode      custom_errors.ErrorCode
	UnprocessEntityErrorCode custom_errors.ErrorCode
}

// customValidationMessages provides default validation error messages.
func (b *Base) customValidationMessages(err error) map[string]string {
	validationErrors := err.(govalidator.ValidationErrors)
	errorMessages := make(map[string]string)

	for _, fieldError := range validationErrors {
		field := fieldError.StructNamespace()
		message := b.mapValidationMessage(field)
		if message != "" {
			key := b.formatFieldKey(field)
			errorMessages[key] = message
		}
	}

	return errorMessages
}

// ParseAndValidateRequest handles parsing and validating a request.
func (b *Base) ParseAndValidateRequest(request interface{}) error {
	// Parse the incoming request
	base := *b
	if b.Binding != nil {
		if err := b.C.ShouldBindWith(request, b.Binding); err != nil {
			b.C.Error(custom_errors.New(
				err,
				b.BadRequestErrorCode,
				"invalid request"))
			return err
		}
	}

	b = &base

	// Validate the request using the validator
	if err := validator.Validate.Struct(request); err != nil {
		validationErrors := b.customValidationMessages(err)
		b.C.Error(custom_errors.New(
			err,
			b.UnprocessEntityErrorCode,
			"failed to validate request",
			validationErrors))
		return err
	}

	return nil
}

// mapValidationMessage maps field names to error messages.
func (b *Base) mapValidationMessage(field string) string {
	if message, exists := b.FieldMessages[field]; exists {
		return message
	}
	return fmt.Sprintf("Invalid value for %s.", field)
}

// formatFieldKey maps field names to JSON keys.
func (b *Base) formatFieldKey(field string) string {
	if key, exists := b.FieldKeys[field]; exists {
		return key
	}
	return field
}
