package register

import (
	"time"

	"github.com/banggok/boillerplate_architecture/internal/data/entity"
)

// Response represents the response payload for tenant registration.
type Response struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// transform transforms a Tenant entity into a RegisterTenantResponse DTO.
func transform(tenant entity.Tenant) Response {
	return Response{
		ID:        tenant.ID(),
		CreatedAt: tenant.CreatedAt(),
		UpdatedAt: tenant.UpdatedAt(),
	}
}
