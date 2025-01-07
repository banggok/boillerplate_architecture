package repository_test

import (
	"errors"
	"testing"

	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	"github.com/banggok/boillerplate_architecture/internal/data/model"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/repository"
	asserthelper "github.com/banggok/boillerplate_architecture/tests/helper/assert_helper"
	"github.com/banggok/boillerplate_architecture/tests/helper/setup"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetOne(t *testing.T) {
	_, cleanUp, db := setup.TestingEnv(t, false)
	defer cleanUp(db)

	ctx := setupTransaction(db.Slave, db.Master)

	t.Run("account not found", func(t *testing.T) {
		accountRepo := repository.NewGenericRepository[entity.Account, model.Account](model.NewAccountModel)

		account, err := accountRepo.GetOne(ctx)

		assert.Error(t, err)
		assert.True(t, errors.As(err, &custom_errors.CustomError{}))
		var customErr custom_errors.CustomError
		if ce, ok := err.(custom_errors.CustomError); ok {
			customErr = ce
		}

		assert.Equal(t, int(custom_errors.AccountNotFound), customErr.Code)
		assert.Nil(t, account)
	})

	t.Run("tenant not found", func(t *testing.T) {
		tenantRepo := repository.NewGenericRepository[entity.Tenant, model.Tenant](model.NewTenantModel)

		tenant, err := tenantRepo.GetOne(ctx)

		assert.Error(t, err)
		assert.True(t, errors.As(err, &custom_errors.CustomError{}))
		var customErr custom_errors.CustomError
		if ce, ok := err.(custom_errors.CustomError); ok {
			customErr = ce
		}

		assert.Equal(t, int(custom_errors.TenantNotFound), customErr.Code)
		assert.Nil(t, tenant)
	})

}

func TestGetAllWithPagination(t *testing.T) {
	_, cleanUp, db := setup.TestingEnv(t, false)
	defer cleanUp(db)

	ctx := setupTransaction(db.Slave, db.Master)

	repo := repository.NewGenericRepository[entity.Tenant, model.Tenant](model.NewTenantModel)

	tenant := insertTenant(t, repo, ctx)

	res, count, err := repo.GetAllWithPagination(ctx, repository.Pagination{
		Page: 1,
		Size: 1,
	})
	assert.NoError(t, err)

	expectedRes := []entity.Tenant{
		tenant[0],
	}
	asserthelper.AssertStruct(t, res, expectedRes)
	assert.Equal(t, int64(2), count)
	assert.NoError(t, err)
}

func TestGetAll(t *testing.T) {
	_, cleanUp, db := setup.TestingEnv(t, false)
	defer cleanUp(db)
	ctx := setupTransaction(db.Slave, db.Master)

	repo := repository.NewGenericRepository[entity.Tenant, model.Tenant](model.NewTenantModel)

	tenant := insertTenant(t, repo, ctx)

	res, err := repo.Where("id = ?", tenant[1].ID()).GetAll(ctx)
	assert.NoError(t, err)

	asserthelper.AssertStruct(t, res, []entity.Tenant{tenant[1]})

	assert.NoError(t, err)
}

func insertTenant(t *testing.T, repo repository.GenericRepository[entity.Tenant, model.Tenant],
	ctx *gin.Context) []entity.Tenant {
	tenant, err := entity.NewTenant(
		entity.NewTenantIdentity("tenant 1", "email1@test.com", "081234567891"),
		entity.NewTenantStoreInfo("address 1", "UTC", "08:00", "21:00"),
		nil,
	)

	assert.NoError(t, err)

	err = repo.Persist(ctx, &tenant)

	assert.NoError(t, err)

	tenant2, err := entity.NewTenant(
		entity.NewTenantIdentity("tenant 2", "email2@test.com", "081234567892"),
		entity.NewTenantStoreInfo("address 2", "UTC", "08:00", "21:00"),
		nil,
	)

	assert.NoError(t, err)

	err = repo.Persist(ctx, &tenant2)

	assert.NoError(t, err)
	return []entity.Tenant{
		tenant,
		tenant2,
	}
}
