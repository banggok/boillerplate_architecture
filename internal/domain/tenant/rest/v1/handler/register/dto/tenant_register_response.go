package dto

import (
	"appointment_management_system/internal/domain/tenant/entity"
	"time"
)

// RegisterTenantResponse represents the response payload for tenant registration.
type RegisterTenantResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Address      string    `json:"address,omitempty"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone,omitempty"`
	Timezone     string    `json:"timezone"`
	OpeningHours string    `json:"opening_hours,omitempty"`
	ClosingHours string    `json:"closing_hours,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewRegisterTenantResponse transforms a Tenant entity into a RegisterTenantResponse DTO.
func NewRegisterTenantResponse(tenant *entity.Tenant) RegisterTenantResponse {
	return RegisterTenantResponse{
		ID:           tenant.GetID(),
		Name:         tenant.GetName(),
		Address:      tenant.GetAddress(),
		Email:        tenant.GetEmail(),
		Phone:        tenant.GetPhone(),
		Timezone:     tenant.GetTimezone(),
		OpeningHours: tenant.GetOpeningHours(),
		ClosingHours: tenant.GetClosingHours(),
		CreatedAt:    tenant.GetCreatedAt(),
		UpdatedAt:    tenant.GetUpdatedAt(),
	}
}
