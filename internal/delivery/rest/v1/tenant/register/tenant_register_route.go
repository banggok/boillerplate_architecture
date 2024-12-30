package register

import (
	"appointment_management_system/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	tenantV1RouterGroup *gin.RouterGroup,
	serviceConfig services.Config,
) {
	tenantCreateService := serviceConfig.Tenant()
	email := serviceConfig.Email()
	uc := newTenantRegisterUsecase(tenantCreateService, email)
	h := newTenantRegisterHandler(uc)
	tenantV1RouterGroup.POST("", h.Register)
}
