package verify

import (
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/request"
	"github.com/gin-gonic/gin"
)

type Request struct {
	request.Base
	Token string `json:"token" validate:"required"`
}

func (r *Request) ParseAndValidateRequest() error {
	r.Token = r.Base.C.Param("token")
	return nil
}

func newRequest(c *gin.Context) request.IRequest {
	return &Request{
		Base: request.Base{
			C: c,
		},
	}
}
