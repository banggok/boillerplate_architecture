package register

import (
	"appointment_management_system/internal/data/entity"
	"time"
)

// RegisterTenantResponse represents the response payload for tenant registration.
type RegisterTenantResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewRegisterTenantResponse transforms a Tenant entity into a RegisterTenantResponse DTO.
func NewRegisterTenantResponse(tenant *entity.Tenant) RegisterTenantResponse {
	return RegisterTenantResponse{
		ID:        tenant.GetID(),
		CreatedAt: tenant.GetCreatedAt(),
		UpdatedAt: tenant.GetUpdatedAt(),
	}
}
