package login

import (
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/request"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Request struct {
	request.Base
	Email    string `json:"email" validate:"required,email" example:"admin@example.com"` // Admin's email
	Password string `json:"password" validate:"required"`                                // Admin's email
}

var fieldMessages = map[string]string{
	"Request.Email":    "Email is required and must be in a valid format.",
	"Request.Password": "Password is required",
}

var fieldKeys = map[string]string{
	"Request.Email":    "email",
	"Request.Password": "password",
}

// ParseAndValidateRequest implements request.IRequest.
// Subtle: this method shadows the method (Base).ParseAndValidateRequest of Request.Base.
func (r *Request) ParseAndValidateRequest(c *gin.Context) error {
	return r.Base.ParseAndValidateRequest(c, r)
}

func newRequest() request.IRequest {
	return &Request{
		Base: request.Base{
			Binding:                  binding.JSON,
			FieldMessages:            fieldMessages,
			FieldKeys:                fieldKeys,
			BadRequestErrorCode:      custom_errors.AccountBadRequest,
			UnprocessEntityErrorCode: custom_errors.AccountUnprocessEntity,
		},
	}
}
