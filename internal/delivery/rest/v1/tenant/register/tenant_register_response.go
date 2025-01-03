package register

import (
	"time"

	"github.com/banggok/boillerplate_architecture/internal/data/entity"
)

// RegisterTenantResponse represents the response payload for tenant registration.
type RegisterTenantResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewRegisterTenantResponse transforms a Tenant entity into a RegisterTenantResponse DTO.
func NewRegisterTenantResponse(tenant entity.Tenant) RegisterTenantResponse {
	return RegisterTenantResponse{
		ID:        tenant.ID(),
		CreatedAt: tenant.CreatedAt(),
		UpdatedAt: tenant.UpdatedAt(),
	}
}
