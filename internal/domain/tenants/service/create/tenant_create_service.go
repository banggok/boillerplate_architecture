package create

import (
	"appointment_management_system/internal/data/entity"
	"appointment_management_system/internal/data/model"
	"appointment_management_system/internal/pkg/custom_errors"
	"appointment_management_system/internal/pkg/repository"
	"fmt"

	"github.com/gin-gonic/gin"
)

type TenantCreateService interface {
	Execute(*gin.Context, *entity.Tenant) error
}

type tenantCreateService struct {
	tenantRepo  repository.GenericRepository[entity.Tenant, model.Tenant]
	accountRepo repository.GenericRepository[entity.Account, model.Account]
}

func (t *tenantCreateService) Execute(ctx *gin.Context, tenant *entity.Tenant) error {
	if tenant == nil {
		return fmt.Errorf("tenant can not be empty in TenantCreateService")
	}

	if tenant.GetAccounts() == nil {
		return custom_errors.New(nil, custom_errors.InternalServerError, "account can not empty when create tenant")
	}

	// Validate tenant existence (phone and email)
	countTenant, err := t.tenantRepo.Count(ctx, repository.Condition{
		Query: "phone = ? or email = ?",
		Args: []interface{}{
			tenant.GetPhone(),
			tenant.GetEmail(),
		},
	})

	if err != nil {
		return custom_errors.New(err, custom_errors.InternalServerError, "failed to count tenant")
	}

	if countTenant > 0 {
		return custom_errors.New(nil, custom_errors.TenantConflictEntity, "tenant phone number and email already exists")
	}

	accounts := tenant.GetAccounts()
	account := (*accounts)[0]
	countAccount, err := t.accountRepo.Count(ctx, repository.Condition{
		Query: "phone = ? or email = ?",
		Args: []interface{}{
			account.GetPhone(),
			account.GetEmail(),
		},
	})

	if err != nil {
		return custom_errors.New(err, custom_errors.InternalServerError, "failed to count account")
	}

	if countAccount > 0 {
		return custom_errors.New(nil, custom_errors.TenantConflictEntity, "account phone number and email already exists")
	}

	return t.tenantRepo.Persist(ctx, tenant)
}

func NewTenantCreateService() TenantCreateService {
	tenantRepo := repository.NewGenericRepository[entity.Tenant, model.Tenant](
		model.NewTenantModel,
	)
	accountRepo := repository.NewGenericRepository[entity.Account, model.Account](
		model.NewAccountModel,
	)
	return &tenantCreateService{
		tenantRepo:  tenantRepo,
		accountRepo: accountRepo,
	}
}
