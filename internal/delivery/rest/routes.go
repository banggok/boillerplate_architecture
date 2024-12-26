package rest

import (
	"appointment_management_system/internal/delivery/rest/v1/tenant"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup) {
	v1 := api.Group("v1/")
	tenant.RegisterRoutes(v1)
}
