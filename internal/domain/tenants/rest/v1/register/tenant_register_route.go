package register

import (
	"appointment_management_system/internal/domain/tenants/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func NewTenantV1RegisterRoutes(
	tenantV1RouterGroup *gin.RouterGroup,
	validator *validator.Validate,
	createTenantRepository repository.TenantCreateRepository,
) {
	uc := newTenantRegisterUsecase(createTenantRepository)
	h := newTenantRegisterHandler(validator, uc)
	tenantV1RouterGroup.POST("/register", h.Register)
}
