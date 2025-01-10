package verify

import (
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/request"
	"github.com/gin-gonic/gin"
)

var _ request.IRequest = &Request{}

type Request struct {
	Token string `json:"token" validate:"required"`
}

func (r *Request) ParseAndValidateRequest(c *gin.Context) error {
	r.Token = c.Param("token")
	return nil
}
