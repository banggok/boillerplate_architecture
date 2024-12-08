package repository

import (
	"appointment_management_system/internal/domain/tenants/entity"
	"appointment_management_system/internal/pkg/helper"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var helperGetDB = helper.GetDB

type contextKey string

const ginContextKey contextKey = "ginContext"

type TenantCreateRepository interface {
	Execute(*gin.Context, *entity.Tenant) error
}

type tenantCreateRepository struct {
	isPhoneExists TenantIsPhoneExistsRepository
	isEmailExists TenantIsEmailExistsRepository
}

func (t *tenantCreateRepository) Execute(ctx *gin.Context, tenant *entity.Tenant) error {
	if tenant == nil {
		return fmt.Errorf("tenant can not be empty in TenantCreateRepository")
	}

	// Validate tenant existence (phone and email)
	phoneExists, emailExists, err := t.validateTenantExistence(ctx, tenant)
	if err != nil {
		return fmt.Errorf("failed to validate tenant in TenantCreateRepository: %w", err)
	}

	if phoneExists {
		return fmt.Errorf("phone number already exists")
	}
	if emailExists {
		return fmt.Errorf("email already exists")
	}

	return t.saveTenant(ctx, tenant)
}

// Helper to validate tenant existence
func (t *tenantCreateRepository) validateTenantExistence(ctx *gin.Context, tenant *entity.Tenant) (bool, bool, error) {
	var phoneExists, emailExists bool
	groupCtx := context.WithValue(
		ctx.Request.Context(),
		ginContextKey,
		ctx)

	g, groupCtx := errgroup.WithContext(groupCtx)

	// Check phone existence
	g.Go(func() error {
		return t.checkPhoneExistence(groupCtx, tenant, &phoneExists)
	})

	// Check email existence
	g.Go(func() error {
		return t.checkEmailExistence(groupCtx, tenant, &emailExists)
	})

	// Wait for goroutines to complete
	if err := g.Wait(); err != nil {
		return false, false, err
	}

	return phoneExists, emailExists, nil
}

// Helper to check phone existence
func (t *tenantCreateRepository) checkPhoneExistence(ctx context.Context, tenant *entity.Tenant, phoneExists *bool) error {
	ginCtx, ok := ctx.Value(ginContextKey).(*gin.Context)
	if !ok {
		return fmt.Errorf("failed to retrieve *gin.Context from context")
	}
	res, err := t.isPhoneExists.Execute(ginCtx, tenant.GetPhone())
	if err != nil {
		return err
	}
	*phoneExists = res != nil && *res
	return nil
}

// Helper to check email existence
func (t *tenantCreateRepository) checkEmailExistence(ctx context.Context, tenant *entity.Tenant, emailExists *bool) error {
	ginCtx, ok := ctx.Value(ginContextKey).(*gin.Context)
	if !ok {
		return fmt.Errorf("failed to retrieve *gin.Context from context")
	}
	res, err := t.isEmailExists.Execute(ginCtx, tenant.GetEmail())
	if err != nil {
		return err
	}
	*emailExists = res != nil && *res
	return nil
}

// Helper to save tenant to the database
func (t *tenantCreateRepository) saveTenant(ctx *gin.Context, tenant *entity.Tenant) error {
	tx, err := helperGetDB(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	tenantModel := newTenantModel(*tenant)

	if err := tx.Create(&tenantModel).Error; err != nil {
		return fmt.Errorf("failed to save tenant to database: %w", err)
	}

	toEntity, err := tenantModel.ToEntity()
	if err != nil {
		return fmt.Errorf("failed to convert tenantModel to entity: %w", err)
	}

	*tenant = *toEntity
	return nil
}

func NewTenantCreateRepository() TenantCreateRepository {
	isPhoneExists := NewTenantIsPhoneExistsRepository()
	isEmailExists := NewTenantIsEmailExistsRepository()
	return &tenantCreateRepository{
		isPhoneExists: isPhoneExists,
		isEmailExists: isEmailExists,
	}
}
