package delivery

import (
	"appointment_management_system/internal/domain/tenants/delivery/v1/register"

	"github.com/gin-gonic/gin"
)

func RegisterTenantRoutes(api *gin.RouterGroup) {
	v1 := api.Group("v1/tenants/")
	register.NewTenantV1RegisterRoutes(v1)
}
