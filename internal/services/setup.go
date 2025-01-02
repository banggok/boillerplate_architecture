package services

import (
	"appointment_management_system/internal/config/app"
	"appointment_management_system/internal/data/entity"
	"appointment_management_system/internal/data/model"
	"appointment_management_system/internal/pkg/repository"
	"appointment_management_system/internal/services/notification/email"
	"appointment_management_system/internal/services/tenant"
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

func Setup(cfg app.AppConfig) Config {
	tenantRepo := repository.NewGenericRepository[entity.Tenant, model.Tenant](model.NewTenantModel)
	accountRepo := repository.NewGenericRepository[entity.Account, model.Account](model.NewAccountModel)

	email := email.NewService(cfg.SMTP, nil)
	tenant := tenant.NewTenantCreateService(tenantRepo, accountRepo)
	return &servicesImplt{
		email:  email,
		tenant: tenant,
	}
}
