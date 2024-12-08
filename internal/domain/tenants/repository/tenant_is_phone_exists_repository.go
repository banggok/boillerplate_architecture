package repository

import "github.com/gin-gonic/gin"

type TenantIsPhoneExistsRepository interface {
	Execute(ctx *gin.Context, phoneNumber string) (*bool, error)
}

type tenantIsPhoneExistsRepository struct {
}

// Execute implements TenantIsPhoneExists.
func (t *tenantIsPhoneExistsRepository) Execute(ctx *gin.Context, phoneNumber string) (*bool, error) {
	panic("unimplemented")
}

func NewTenantIsPhoneExistsRepository() TenantIsPhoneExistsRepository {
	return &tenantIsPhoneExistsRepository{}
}
