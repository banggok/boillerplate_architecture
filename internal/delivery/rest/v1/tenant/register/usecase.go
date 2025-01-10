package register

import (
	"fmt"

	"github.com/banggok/boillerplate_architecture/internal/config/app"
	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	valueobject "github.com/banggok/boillerplate_architecture/internal/data/entity/value_object"
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/request"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/services/notification/email"
	"github.com/banggok/boillerplate_architecture/internal/services/tenant"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// usecase defines the interface for registering a tenant
type usecase interface {
	execute(ctx *gin.Context, request request.IRequest) (entity.Tenant, error)
}

type usecaseImpl struct {
	createTenantService tenant.Service
	emailService        email.Service
}

// execute implements TenantRegisterUsecase.
func (t *usecaseImpl) execute(ctx *gin.Context, iRequest request.IRequest) (entity.Tenant, error) {
	request, ok := iRequest.(*Request)
	if !ok {
		return nil, custom_errors.New(nil, custom_errors.InternalServerError, "request invalid")
	}
	emailVerification, err := entity.NewAccountVerification(valueobject.EMAIL_VERIFICATION, nil)
	if err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.InternalServerError,
			"failed to create new account entity")
	}

	resetPassword, err := entity.NewAccountVerification(valueobject.RESET_PASSWORD, nil)
	if err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.InternalServerError,
			"failed to create new account entity")
	}

	account, plainPassword, err := entity.NewAccount(
		entity.NewAccountIdentity(request.Account.Name, request.Account.Email, request.Account.Phone, nil),
		nil, &[]entity.AccountVerification{emailVerification, resetPassword})
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

	if emailVerification.Token() == nil {
		return nil, custom_errors.New(nil, custom_errors.InternalServerError, "token was empty")
	}

	t.sendWelcomeEmail(emailVerification, tenant, plainPassword)

	return tenant, nil
}

func (t *usecaseImpl) sendWelcomeEmail(emailVerification entity.AccountVerification, tenant entity.Tenant, plainPassword *string) {
	if app.AppConfig.Environment == app.ENV_PROD {
		url := fmt.Sprintf("%s/verify/%s", app.AppConfig.BaseUrl, *emailVerification.Token())
		go func() {
			if tenant.Accounts() == nil || len(*tenant.Accounts()) == 0 {
				logrus.Errorf("tenant account was empty")
			}
			account := (*tenant.Accounts())[0]
			if err := t.emailService.SendWelcomeEmail(
				account.Email(),
				email.WelcomeData{
					TenantName:      tenant.Name(),
					Username:        account.Email(),
					Password:        *plainPassword,
					VerificationURL: url,
				},
			); err != nil {
				logrus.Errorf("failed to send email: %v", err.Error())
			}
		}()
	}
}

func newUsecase(createTenantService tenant.Service, email email.Service) usecase {
	return &usecaseImpl{
		createTenantService: createTenantService,
		emailService:        email,
	}
}
