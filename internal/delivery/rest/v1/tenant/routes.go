package tenant

import (
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/v1/tenant/register"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(v1 *gin.RouterGroup) {
	tenants := v1.Group("/tenants")
	register.Register(tenants)
}
