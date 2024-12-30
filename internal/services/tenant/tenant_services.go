package tenant

import (
	"appointment_management_system/internal/data/entity"
	"appointment_management_system/internal/data/model"
	"appointment_management_system/internal/pkg/custom_errors"
	"appointment_management_system/internal/pkg/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var TenantDuplicate = "phone number and email for tenant or account already used"

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

	if err := t.validateExistence(ctx, *tenant); err != nil {
		return custom_errors.New(err, custom_errors.InternalServerError, "failed to validate tenant existence")
	}

	return t.tenantRepo.Persist(ctx, tenant)
}

func (t *tenantCreateService) validateExistence(ctx *gin.Context, tenant entity.Tenant) error {
	eg := new(errgroup.Group)

	counts := make(map[string]*int64)
	queries := []struct {
		key       string
		repo      string
		condition repository.Condition
	}{
		{
			key:  "countPhone",
			repo: "tenant",
			condition: repository.Condition{
				Query: "phone = ?",
				Args:  []interface{}{tenant.GetPhone()},
			},
		},
		{
			key:  "countEmail",
			repo: "tenant",
			condition: repository.Condition{
				Query: "email = ?",
				Args:  []interface{}{tenant.GetEmail()},
			},
		},
		{
			key:  "countAccountPhone",
			repo: "account",
			condition: repository.Condition{
				Query: "phone = ?",
				Args:  []interface{}{(*tenant.GetAccounts())[0].GetPhone()},
			},
		},
		{
			key:  "countAccountEmail",
			repo: "account",
			condition: repository.Condition{
				Query: "email = ?",
				Args:  []interface{}{(*tenant.GetAccounts())[0].GetEmail()},
			},
		},
	}

	for _, query := range queries {
		q := query
		eg.Go(func() error {
			var err error
			counts[q.key], err = t.countEntity(ctx, q.repo, q.condition)
			if err != nil {
				return custom_errors.New(err, custom_errors.InternalServerError, "failed to count "+q.key)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return custom_errors.New(err, custom_errors.InternalServerError, "failed to validate tenant existence")
	}

	if counts["countPhone"] == nil || counts["countEmail"] == nil ||
		counts["countAccountPhone"] == nil || counts["countAccountEmail"] == nil {
		return custom_errors.New(nil, custom_errors.InternalServerError, "failed to validate tenant existence")
	}

	if *counts["countPhone"] > 0 || *counts["countEmail"] > 0 ||
		*counts["countAccountPhone"] > 0 || *counts["countAccountEmail"] > 0 {
		return custom_errors.New(nil, custom_errors.TenantConflictEntity, TenantDuplicate)
	}

	return nil
}

func (t *tenantCreateService) countEntity(
	ctx *gin.Context,
	repo string,
	condition repository.Condition,
) (*int64, error) {
	if repo == "tenant" {
		return t.tenantRepo.Count(ctx, condition)
	}
	return t.accountRepo.Count(ctx, condition)
}

func NewTenantCreateService(
	tenantRepo repository.GenericRepository[entity.Tenant, model.Tenant],
	accountRepo repository.GenericRepository[entity.Account, model.Account]) Service {
	return &tenantCreateService{
		tenantRepo:  tenantRepo,
		accountRepo: accountRepo,
	}
}
