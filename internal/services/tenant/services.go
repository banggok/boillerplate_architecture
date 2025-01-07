package tenant

import (
	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	"github.com/banggok/boillerplate_architecture/internal/data/model"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/repository"
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
	if tenant == nil || (*tenant).Accounts() == nil {
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
		key   string
		repo  string
		query string
		args  []interface{}
	}{
		{
			key:   "countPhone",
			repo:  "tenant",
			query: "phone = ?",
			args:  []interface{}{tenant.Phone()},
		},
		{
			key:   "countEmail",
			repo:  "tenant",
			query: "email = ?",
			args:  []interface{}{tenant.Email()},
		},
		{
			key:   "countAccountPhone",
			repo:  "account",
			query: "phone = ?",
			args:  []interface{}{(*tenant.Accounts())[0].Phone()},
		},
		{
			key:   "countAccountEmail",
			repo:  "account",
			query: "email = ?",
			args:  []interface{}{(*tenant.Accounts())[0].Email()},
		},
	}

	for _, query := range queries {
		q := query
		eg.Go(func() error {
			var err error
			counts[q.key], err = t.countEntity(ctx, q.repo, q.query, q.args)
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
	query string,
	args []interface{},
) (*int64, error) {
	if repo == "tenant" {
		return t.tenantRepo.Where(query, args...).Count(ctx)
	}
	return t.accountRepo.Where(query, args...).Count(ctx)
}

func New(
	tenantRepo repository.GenericRepository[entity.Tenant, model.Tenant],
	accountRepo repository.GenericRepository[entity.Account, model.Account]) Service {
	return &tenantCreateService{
		tenantRepo:  tenantRepo,
		accountRepo: accountRepo,
	}
}
