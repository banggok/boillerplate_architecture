package v1

import (
	"appointment_management_system/internal/delivery/rest/v1/tenant"
	"appointment_management_system/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, serviceCfg services.Config) {
	v1 := api.Group("/v1")
	tenant.RegisterRoutes(v1, serviceCfg)
}
