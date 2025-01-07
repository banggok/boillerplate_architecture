package services

import (
	"github.com/banggok/boillerplate_architecture/internal/config/app"
	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	"github.com/banggok/boillerplate_architecture/internal/data/model"
	"github.com/banggok/boillerplate_architecture/internal/pkg/repository"
	"github.com/banggok/boillerplate_architecture/internal/services/account"
	accountverification "github.com/banggok/boillerplate_architecture/internal/services/account_verification"
	"github.com/banggok/boillerplate_architecture/internal/services/notification/email"
	"github.com/banggok/boillerplate_architecture/internal/services/tenant"
)

var ServiceConFig config

type config interface {
	Tenant() tenant.Service
	Email() email.Service
	AccountVerification() accountverification.Service
	Account() account.Service
}

type servicesImplt struct {
	tenant              tenant.Service
	email               email.Service
	accountVerification accountverification.Service
	account             account.Service
}

// Account implements Config.
func (s *servicesImplt) Account() account.Service {
	return s.account
}

// AccountVerification implements Config.
func (s *servicesImplt) AccountVerification() accountverification.Service {
	return s.accountVerification
}

// Email implements Services.
func (s *servicesImplt) Email() email.Service {
	return s.email
}

// Tenant implements Services.
func (s *servicesImplt) Tenant() tenant.Service {
	return s.tenant
}

func Setup() {
	tenantRepo := repository.NewGenericRepository[entity.Tenant, model.Tenant](model.NewTenantModel)
	accountRepo := repository.NewGenericRepository[entity.Account, model.Account](model.NewAccountModel)
	accountVerificationRepo := repository.NewGenericRepository[entity.AccountVerification, model.AccountVerification](model.NewAccountVerification)

	ServiceConFig = &servicesImplt{
		email:               email.New(app.AppConfig.SMTP, nil),
		tenant:              tenant.New(tenantRepo, accountRepo),
		accountVerification: accountverification.NewService(accountVerificationRepo),
		account:             account.New(accountRepo),
	}
}
