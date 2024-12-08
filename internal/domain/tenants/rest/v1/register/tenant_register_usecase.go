package register

import (
	"appointment_management_system/internal/domain/tenants/entity"
	"appointment_management_system/internal/domain/tenants/repository"
	"appointment_management_system/internal/domain/tenants/rest/v1/register/dto"
	"fmt"

	"github.com/gin-gonic/gin"
)

// TenantRegisterUsecase defines the interface for registering a tenant
type TenantRegisterUsecase interface {
	Execute(ctx *gin.Context, request dto.RegisterTenantRequest) (*entity.Tenant, error)
}

type tenantRegisterUsecase struct {
	createTenantService repository.TenantCreateRepository
}

// Execute implements TenantRegisterUsecase.
func (t *tenantRegisterUsecase) Execute(ctx *gin.Context, request dto.RegisterTenantRequest) (*entity.Tenant, error) {
	tenant, err := entity.NewTenant(entity.NewTenantParams{
		Name:         request.Name,
		Address:      request.Address,
		Email:        request.Email,
		Phone:        request.Phone,
		Timezone:     request.Timezone,
		OpeningHours: request.OpeningHours,
		ClosingHours: request.ClosingHours,
	})

	if err != nil {
		return nil, fmt.Errorf("failed TenantRegisterUsecase.Execute entity.NewTenant: %w", err)
	}

	if err := t.createTenantService.Execute(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed TenantRegisterUsecase.Execute createTenantService.Execute: %w", err)
	}

	return tenant, nil
}

func newTenantRegisterUsecase(createTenantService repository.TenantCreateRepository) TenantRegisterUsecase {
	return &tenantRegisterUsecase{
		createTenantService: createTenantService,
	}
}
