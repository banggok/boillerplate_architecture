package repository

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type TenantIsEmailExistsRepository interface {
	Execute(ctx *gin.Context, phoneNumber string) (*bool, error)
}

type tenantIsEmailExistsRepository struct {
}

// Execute implements TenantIsPhoneExists.
func (t *tenantIsEmailExistsRepository) Execute(ctx *gin.Context, email string) (*bool, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	// Get database transaction from context
	tx, err := helperGetDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get database transaction: %w", err)
	}

	// Query the database
	var count int64
	if err := tx.Model(&tenantModel{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}

	// Determine if the email exists
	exists := count > 0
	return &exists, nil
}

func NewTenantIsEmailExistsRepository() TenantIsEmailExistsRepository {
	return &tenantIsEmailExistsRepository{}
}
