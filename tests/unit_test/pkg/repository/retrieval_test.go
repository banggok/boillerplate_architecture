package repository_test

import (
	"appointment_management_system/internal/data/entity"
	"appointment_management_system/internal/data/model"
	"appointment_management_system/internal/pkg/custom_errors"
	"appointment_management_system/internal/pkg/repository"
	"appointment_management_system/tests/e2e/helper/setup"
	"encoding/json"
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetOne(t *testing.T) {
	_, cleanUp, db := setup.TestingEnv(t, false)
	defer cleanUp(db)

	ctx := setupTransaction(db.Slave, db.Master)

	t.Run("account not found", func(t *testing.T) {
		accountRepo := repository.NewGenericRepository[entity.Account, model.Account](model.NewAccountModel)

		account, err := accountRepo.GetOne(ctx, &repository.Condition{})

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

		tenant, err := tenantRepo.GetOne(ctx, &repository.Condition{})

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

	res, count, err := repo.GetAllWithPagination(ctx, nil, repository.Pagination{
		Page: 1,
		Size: 1,
	})
	assert.NoError(t, err)

	expectedRes := []entity.Tenant{
		tenant[0],
	}
	expectedResJson, err := json.Marshal(expectedRes)
	assert.NoError(t, err)

	resJson, err := json.Marshal(res)
	assert.NoError(t, err)

	assert.JSONEq(t, string(expectedResJson), string(resJson))
	assert.Equal(t, int64(2), count)
	assert.NoError(t, err)
}

func TestGetAll(t *testing.T) {
	_, cleanUp, db := setup.TestingEnv(t, false)
	defer cleanUp(db)
	ctx := setupTransaction(db.Slave, db.Master)

	repo := repository.NewGenericRepository[entity.Tenant, model.Tenant](model.NewTenantModel)

	tenant := insertTenant(t, repo, ctx)

	res, err := repo.GetAll(ctx, nil)
	assert.NoError(t, err)

	expectedResJson, err := json.Marshal(tenant)
	assert.NoError(t, err)

	resJson, err := json.Marshal(res)
	assert.NoError(t, err)

	assert.JSONEq(t, string(expectedResJson), string(resJson))
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
