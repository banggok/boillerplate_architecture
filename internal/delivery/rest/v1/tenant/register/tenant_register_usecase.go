package register

import (
	"appointment_management_system/internal/data/entity"
	"appointment_management_system/internal/pkg/custom_errors"
	"appointment_management_system/internal/services/notification/email"
	"appointment_management_system/internal/services/tenant"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// TenantRegisterUsecase defines the interface for registering a tenant
type TenantRegisterUsecase interface {
	Execute(ctx *gin.Context, request RegisterTenantRequest) (*entity.Tenant, error)
}

type tenantRegisterUsecase struct {
	createTenantService tenant.Service
	emailService        email.Service
}

// Execute implements TenantRegisterUsecase.
func (t *tenantRegisterUsecase) Execute(ctx *gin.Context, request RegisterTenantRequest) (*entity.Tenant, error) {

	account, plainPassword, err := entity.NewAccount(
		entity.NewAccountIdentity(request.Account.Name, request.Account.Email, request.Account.Phone),
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

	go func() {
		if tenant.GetAccounts() == nil || len(*tenant.GetAccounts()) == 0 {
			logrus.Errorf("tenant account was empty")
		}
		account := (*tenant.GetAccounts())[0]
		if err := t.emailService.SendWelcomeEmail(account.GetEmail(), tenant.GetName(), account.GetEmail(),
			*plainPassword, ""); err != nil {
			logrus.Errorf("failed to send email: %v", err.Error())
		}
	}()

	return tenant, nil
}

func newTenantRegisterUsecase(createTenantService tenant.Service, email email.Service) TenantRegisterUsecase {
	return &tenantRegisterUsecase{
		createTenantService: createTenantService,
		emailService:        email,
	}
}
