package services

import (
	"github.com/banggok/boillerplate_architecture/internal/config/app"
	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	"github.com/banggok/boillerplate_architecture/internal/data/model"
	"github.com/banggok/boillerplate_architecture/internal/pkg/repository"
	"github.com/banggok/boillerplate_architecture/internal/services/notification/email"
	"github.com/banggok/boillerplate_architecture/internal/services/tenant"
)

type Config interface {
	Tenant() tenant.Service
	Email() email.Service
}

type servicesImplt struct {
	tenant tenant.Service
	email  email.Service
}

// Email implements Services.
func (s *servicesImplt) Email() email.Service {
	return s.email
}

// Tenant implements Services.
func (s *servicesImplt) Tenant() tenant.Service {
	return s.tenant
}

func Setup(cfg app.Config) Config {
	tenantRepo := repository.NewGenericRepository[entity.Tenant, model.Tenant](model.NewTenantModel)
	accountRepo := repository.NewGenericRepository[entity.Account, model.Account](model.NewAccountModel)

	email := email.NewService(cfg.SMTP, nil)
	tenant := tenant.NewTenantCreateService(tenantRepo, accountRepo)
	return &servicesImplt{
		email:  email,
		tenant: tenant,
	}
}
