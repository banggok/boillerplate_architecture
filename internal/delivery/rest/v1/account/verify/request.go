package verify

import (
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/request"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var _ request.IRequest = &Request{}

type Request struct {
	request.Base
	Token string `json:"token" validate:"required"`
}

func (r *Request) ParseAndValidateRequest(c *gin.Context, _ binding.Binding) error {
	r.Token = c.Param("token")
	return nil
}
