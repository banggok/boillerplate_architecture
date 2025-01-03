package register

import (
	"github.com/banggok/boillerplate_architecture/internal/config/app"
	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/services/notification/email"
	"github.com/banggok/boillerplate_architecture/internal/services/tenant"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// usecase defines the interface for registering a tenant
type usecase interface {
	execute(ctx *gin.Context, request Request) (entity.Tenant, error)
}

type usecaseImpl struct {
	createTenantService tenant.Service
	emailService        email.Service
}

// execute implements TenantRegisterUsecase.
func (t *usecaseImpl) execute(ctx *gin.Context, request Request) (entity.Tenant, error) {
	accountVerification, err := entity.NewAccountVerification(entity.EMAIL_VERIFICATION, nil)
	if err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.InternalServerError,
			"failed to create new account entity")
	}
	account, plainPassword, err := entity.NewAccount(
		entity.NewAccountIdentity(request.Account.Name, request.Account.Email, request.Account.Phone, nil),
		nil, &[]entity.AccountVerification{accountVerification})
	if err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.InternalServerError,
			"failed to create new account entity")
	}

	tenant, err := entity.NewTenant(
		entity.NewTenantIdentity(request.Name, request.Email, request.Phone),
		entity.NewTenantStoreInfo(request.Address, request.Timezone, request.OpeningHours, request.ClosingHours),
		&[]entity.Account{account})

	if err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.InternalServerError,
			"failed to create new tenant entity")
	}

	if err := t.createTenantService.Create(ctx, &tenant); err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.InternalServerError,
			"failed to new tenant")
	}

	go func() {
		cfg := app.AppConfig
		if cfg.Environment != app.ENV_PROD {
			return
		}
		if tenant.Accounts() == nil || len(*tenant.Accounts()) == 0 {
			logrus.Errorf("tenant account was empty")
		}
		account := (*tenant.Accounts())[0]
		if err := t.emailService.SendWelcomeEmail(
			account.Email(),
			email.WelcomeData{
				TenantName: tenant.Name(),
				Username:   account.Email(),
				Password:   *plainPassword,
				LoginURL:   "",
			},
		); err != nil {
			logrus.Errorf("failed to send email: %v", err.Error())
		}
	}()

	return tenant, nil
}

func newUsecase(createTenantService tenant.Service, email email.Service) usecase {
	return &usecaseImpl{
		createTenantService: createTenantService,
		emailService:        email,
	}
}
