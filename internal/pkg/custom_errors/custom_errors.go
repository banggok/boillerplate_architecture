package custom_errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type CustomError struct {
	Code     int               `json:"code"`
	Message  string            `json:"message"`
	HTTPCode int               `json:"-"` // HTTP status code for the error
	Err      error             `json:"-"` // Underlying error
	Details  map[string]string `json:"details"`
}

func (e CustomError) Error() string {
	return e.Message
}

func New(err error, code ErrorCode, message string, detail ...map[string]string) CustomError {
	logger := logrus.StandardLogger()

	// If err is nil, use message as error message
	if err == nil {
		err = errors.New(message)
	}
	logger.WithFields(logrus.Fields{
		message: fmt.Sprintf("%s: %s", message, err.Error()),
	})

	// If err is already a CustomError, avoid wrapping it again
	if customErr, ok := err.(CustomError); ok {
		return customErr
	}
	// Fetch the error definition from the dictionary
	def, exists := errorDict[code]
	if !exists {
		// Fallback to a default error if the key is not found
		def = errorDict[InternalServerError]
	}

	// Use the first detail map if provided
	details := map[string]string{
		"error": err.Error(),
	}
	if len(detail) > 0 {
		details = detail[0]
	}

	// Create a new CustomError instance
	return CustomError{
		Code:     int(code),
		Message:  message,
		HTTPCode: def.HTTPCode,
		Err:      err,
		Details:  details,
	}
}

// ErrorCode represents the type for error keys
type ErrorCode int

// Enum-like constants for error keys
const (
	InternalServerError ErrorCode = 10500

	TenantUnprocessEntity ErrorCode = 20422
	TenantBadRequest      ErrorCode = 20400
	TenantConflictEntity  ErrorCode = 20409
	TenantNotFound        ErrorCode = 20404

	AccountUnprocessEntity ErrorCode = 30422
	AccountBadRequest      ErrorCode = 30400
	AccountConflictEntity  ErrorCode = 30409
	AccountNotFound        ErrorCode = 30404

	AccountVerificationNotFound        ErrorCode = 40404
	AccountVerificationUnprocessEntity ErrorCode = 40422

	// Add more keys as needed
)

// errorDict is a global dictionary of errors
var errorDict = map[ErrorCode]CustomError{
	InternalServerError: {
		HTTPCode: http.StatusInternalServerError,
		Message:  "Internal Server Error",
	},
	TenantUnprocessEntity: {
		HTTPCode: http.StatusUnprocessableEntity,
	},
	TenantBadRequest: {
		HTTPCode: http.StatusBadRequest,
	},
	TenantConflictEntity: {
		HTTPCode: http.StatusConflict,
	},
	TenantNotFound: {
		HTTPCode: http.StatusNotFound,
	},
}
