package repository

import (
	"appointment_management_system/internal/domain/tenants/entity"
	"fmt"

	"github.com/gin-gonic/gin"
)

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

	if err := t.checkPhoneExistence(ctx, tenant, &phoneExists); err != nil {
		return phoneExists, emailExists, err
	}
	if err := t.checkEmailExistence(ctx, tenant, &emailExists); err != nil {
		return phoneExists, emailExists, err
	}

	return phoneExists, emailExists, nil
}

// Helper to check phone existence
func (t *tenantCreateRepository) checkPhoneExistence(ctx *gin.Context, tenant *entity.Tenant, phoneExists *bool) error {

	res, err := t.isPhoneExists.Execute(ctx, tenant.GetPhone())
	if err != nil {
		return err
	}
	*phoneExists = res != nil && *res
	return nil
}

// Helper to check email existence
func (t *tenantCreateRepository) checkEmailExistence(ctx *gin.Context, tenant *entity.Tenant, emailExists *bool) error {
	res, err := t.isEmailExists.Execute(ctx, tenant.GetEmail())
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
