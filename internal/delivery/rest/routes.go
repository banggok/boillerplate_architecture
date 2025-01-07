package rest

import (
	v1 "github.com/banggok/boillerplate_architecture/internal/delivery/rest/v1"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	api := server.Group("/api")
	v1.RegisterRoutes(api)
}
