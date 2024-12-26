package register

import (
	"appointment_management_system/internal/data/entity"
	"appointment_management_system/internal/pkg/custom_errors"
	"appointment_management_system/internal/services/tenant"

	"github.com/gin-gonic/gin"
)

// TenantRegisterUsecase defines the interface for registering a tenant
type TenantRegisterUsecase interface {
	Execute(ctx *gin.Context, request RegisterTenantRequest) (*entity.Tenant, error)
}

type tenantRegisterUsecase struct {
	createTenantService tenant.Service
}

// Execute implements TenantRegisterUsecase.
func (t *tenantRegisterUsecase) Execute(ctx *gin.Context, request RegisterTenantRequest) (*entity.Tenant, error) {

	account, _, err := entity.NewAccount(
		entity.NewAccountIdentity(request.Name, request.Email, request.Phone),
		nil)
	if err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.InternalServerError,
			"failed to validate new account entity")
	}

	tenant, err := entity.NewTenant(
		entity.NewTenantIdentity(request.Name, request.Email, request.Phone),
		entity.NewTenantStoreInfo(request.Address, request.Timezone, request.OpeningHours, request.ClosingHours),
		&[]entity.Account{*account})

	if err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.InternalServerError,
			"failed to validate new tenant entity")
	}

	if err := t.createTenantService.Create(ctx, tenant); err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.InternalServerError,
			"failed to new tenant")
	}

	return tenant, nil
}

func newTenantRegisterUsecase(createTenantService tenant.Service) TenantRegisterUsecase {
	return &tenantRegisterUsecase{
		createTenantService: createTenantService,
	}
}
