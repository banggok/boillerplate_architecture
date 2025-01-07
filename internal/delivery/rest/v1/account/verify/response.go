package verify

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
func transform(account entity.Account) Response {
	return Response{
		ID:        account.ID(),
		CreatedAt: account.CreatedAt(),
		UpdatedAt: account.UpdatedAt(),
	}
}
