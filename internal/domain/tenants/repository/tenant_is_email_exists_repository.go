package repository

import "github.com/gin-gonic/gin"

type TenantIsEmailExistsRepository interface {
	Execute(ctx *gin.Context, phoneNumber string) (*bool, error)
}

type tenantIsEmailExistsRepository struct {
}

// Execute implements TenantIsPhoneExists.
func (t *tenantIsEmailExistsRepository) Execute(ctx *gin.Context, email string) (*bool, error) {
	panic("unimplemented")
}

func NewTenantIsEmailExistsRepository() TenantIsEmailExistsRepository {
	return &tenantIsEmailExistsRepository{}
}
