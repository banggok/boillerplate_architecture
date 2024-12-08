package service

import (
	"appointment_management_system/internal/domain/tenant/entity"

	"github.com/gin-gonic/gin"
)

type TenantCreateService interface {
	Execute(*gin.Context, *entity.Tenant) error
}

type tenantCreateService struct{}

// Execute implements TenantCreateService.
func (t *tenantCreateService) Execute(*gin.Context, *entity.Tenant) error {
	// check email dan no telpon exists atau tidak
	panic("unimplemented")
}

func NewTenantCreateService() TenantCreateService {
	return &tenantCreateService{}
}
