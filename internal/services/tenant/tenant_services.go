package tenant

import (
	"appointment_management_system/internal/data/entity"
	"appointment_management_system/internal/data/model"
	"appointment_management_system/internal/pkg/custom_errors"
	"appointment_management_system/internal/pkg/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type Service interface {
	Create(*gin.Context, *entity.Tenant) error
}

type tenantCreateService struct {
	tenantRepo  repository.GenericRepository[entity.Tenant, model.Tenant]
	accountRepo repository.GenericRepository[entity.Account, model.Account]
}

func (t *tenantCreateService) Create(ctx *gin.Context, tenant *entity.Tenant) error {
	if tenant == nil || tenant.GetAccounts() == nil {
		return custom_errors.New(nil, custom_errors.TenantUnprocessEntity,
			"tenant and account can not empty when create tenant")
	}

	accounts := tenant.GetAccounts()
	account := (*accounts)[0]

	if err := t.validateExistence(ctx, *tenant); err != nil {
		return custom_errors.New(err, custom_errors.InternalServerError, "failed to validate tenant existence")
	}

	countAccount, err := t.accountRepo.Count(ctx, repository.Condition{
		Query: "phone = ? or email = ?",
		Args: []interface{}{
			account.GetPhone(),
			account.GetEmail(),
		},
	})

	if err != nil || countAccount == nil {
		return custom_errors.New(err, custom_errors.InternalServerError, "failed to count account")
	}

	if *countAccount > 0 {
		return custom_errors.New(nil, custom_errors.TenantConflictEntity, "account phone number and email already exists")
	}

	return t.tenantRepo.Persist(ctx, tenant)
}

func (t *tenantCreateService) validateExistence(ctx *gin.Context, tenant entity.Tenant) error {
	eg := new(errgroup.Group)

	var countPhone, countEmail *int64
	var countAccountPhone, countAccountEmail *int64
	account := (*tenant.GetAccounts())[0]

	eg.Go(func() error {
		var err error
		countPhone, err = t.tenantRepo.Count(ctx, repository.Condition{
			Query: "phone = ?",
			Args: []interface{}{
				tenant.GetPhone(),
			},
		})

		if err != nil {
			return custom_errors.New(err, custom_errors.InternalServerError, "failed to count tenant phone")
		}

		return nil
	})

	eg.Go(func() error {
		var err error
		countEmail, err = t.tenantRepo.Count(ctx, repository.Condition{
			Query: "email = ?",
			Args: []interface{}{
				tenant.GetEmail(),
			},
		})

		if err != nil {
			return custom_errors.New(err, custom_errors.InternalServerError, "failed to count tenant email")
		}

		return nil
	})

	eg.Go(func() error {
		var err error
		countAccountEmail, err = t.accountRepo.Count(ctx, repository.Condition{
			Query: "email = ?",
			Args: []interface{}{
				account.GetEmail(),
			},
		})

		if err != nil {
			return custom_errors.New(err, custom_errors.InternalServerError, "failed to count account email")
		}

		return nil
	})

	eg.Go(func() error {
		var err error
		countAccountPhone, err = t.accountRepo.Count(ctx, repository.Condition{
			Query: "email = ?",
			Args: []interface{}{
				account.GetPhone(),
			},
		})

		if err != nil {
			return custom_errors.New(err, custom_errors.InternalServerError, "failed to count account phone")
		}

		return nil
	})

	if err := eg.Wait(); err != nil || countPhone == nil || countEmail == nil || countAccountEmail == nil ||
		countAccountPhone == nil {
		return custom_errors.New(err, custom_errors.InternalServerError, "faild to validate tenant existence")
	}

	if *countPhone > 0 || *countEmail > 0 || *countAccountEmail > 0 || *countAccountPhone > 0 {
		return custom_errors.New(nil, custom_errors.TenantConflictEntity,
			"phone number and email for tenant or account already used")
	}
	return nil
}

func NewTenantCreateService() Service {
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
