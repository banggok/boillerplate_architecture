package rest

import (
	v1 "appointment_management_system/internal/delivery/rest/v1"
	"appointment_management_system/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine, serviceCfg services.Config) {
	api := server.Group("/api")
	v1.RegisterRoutes(api, serviceCfg)
}
