package register

import (
	"appointment_management_system/internal/services/tenant"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	tenantV1RouterGroup *gin.RouterGroup,
) {
	tenantCreateService := tenant.NewTenantCreateService()
	uc := newTenantRegisterUsecase(tenantCreateService)
	h := newTenantRegisterHandler(uc)
	tenantV1RouterGroup.POST("/", h.Register)
}
