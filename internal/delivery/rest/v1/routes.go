package v1

import (
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/v1/tenant"
	"github.com/banggok/boillerplate_architecture/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, serviceCfg services.Config) {
	v1 := api.Group("/v1")
	tenant.RegisterRoutes(v1, serviceCfg)
}
