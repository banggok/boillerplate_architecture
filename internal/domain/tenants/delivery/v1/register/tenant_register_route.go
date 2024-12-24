package register

import (
	"appointment_management_system/internal/domain/tenants/service/create"

	"github.com/gin-gonic/gin"
)

func NewTenantV1RegisterRoutes(
	tenantV1RouterGroup *gin.RouterGroup,
) {
	tenantCreateService := create.NewTenantCreateService()
	uc := newTenantRegisterUsecase(tenantCreateService)
	h := newTenantRegisterHandler(uc)
	tenantV1RouterGroup.POST("/register", h.Register)
}
