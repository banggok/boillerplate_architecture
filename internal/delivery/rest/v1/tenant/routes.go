package tenant

import (
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/v1/tenant/register"
	"github.com/banggok/boillerplate_architecture/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(v1 *gin.RouterGroup, serviceConfig services.Config) {
	tenants := v1.Group("/tenants")
	register.RegisterRoute(tenants, serviceConfig)
}
