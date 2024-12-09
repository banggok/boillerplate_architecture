package repository

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type TenantIsPhoneExistsRepository interface {
	Execute(ctx *gin.Context, phoneNumber string) (*bool, error)
}

type tenantIsPhoneExistsRepository struct {
}

// Execute implements TenantIsPhoneExists.
func (t *tenantIsPhoneExistsRepository) Execute(ctx *gin.Context, phoneNumber string) (*bool, error) {
	if phoneNumber == "" {
		return nil, fmt.Errorf("phone number cannot be empty")
	}

	// Get database transaction from context
	tx, err := helperGetDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get database transaction: %w", err)
	}

	// Query the database
	var count int64
	if err := tx.Model(&tenantModel{}).Where("phone = ?", phoneNumber).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}

	// Determine if the phone number exists
	exists := count > 0
	return &exists, nil
}

func NewTenantIsPhoneExistsRepository() TenantIsPhoneExistsRepository {
	return &tenantIsPhoneExistsRepository{}
}
