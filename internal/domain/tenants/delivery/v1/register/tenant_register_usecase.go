package register

import (
	"appointment_management_system/internal/data/entity"
	"appointment_management_system/internal/domain/tenants/service/create"
	"appointment_management_system/internal/pkg/custom_errors"

	"github.com/gin-gonic/gin"
)

// TenantRegisterUsecase defines the interface for registering a tenant
type TenantRegisterUsecase interface {
	Execute(ctx *gin.Context, request RegisterTenantRequest) (*entity.Tenant, error)
}

type tenantRegisterUsecase struct {
	createTenantService create.TenantCreateService
}

// Execute implements TenantRegisterUsecase.
func (t *tenantRegisterUsecase) Execute(ctx *gin.Context, request RegisterTenantRequest) (*entity.Tenant, error) {
	account, err := entity.NewAccount(
		entity.NewAccountParams{
			Name:  request.Account.Name,
			Email: request.Account.Email,
			Phone: request.Account.Phone,
		},
	)
	if err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.InternalServerError,
			"failed to validate new account entity")
	}

	tenant, err := entity.NewTenant(entity.NewTenantParams{
		Name:         request.Name,
		Address:      request.Address,
		Email:        request.Email,
		Phone:        request.Phone,
		Timezone:     request.Timezone,
		OpeningHours: request.OpeningHours,
		ClosingHours: request.ClosingHours,
		Account: &[]entity.Account{
			*account,
		},
	})

	if err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.InternalServerError,
			"failed to validate new tenant entity")
	}

	if err := t.createTenantService.Execute(ctx, tenant); err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.InternalServerError,
			"failed to new tenant")
	}

	return tenant, nil
}

func newTenantRegisterUsecase(createTenantService create.TenantCreateService) TenantRegisterUsecase {
	return &tenantRegisterUsecase{
		createTenantService: createTenantService,
	}
}
