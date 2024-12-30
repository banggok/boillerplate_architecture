package tenant

import (
	"appointment_management_system/internal/delivery/rest/v1/tenant/register"
	"appointment_management_system/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(v1 *gin.RouterGroup, serviceConfig services.Config) {
	tenants := v1.Group("/tenants")
	register.RegisterRoutes(tenants, serviceConfig)
}
