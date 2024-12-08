package dto

import (
	"appointment_management_system/internal/domain/tenant/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRegisterTenantResponse(t *testing.T) {
	t.Run("Transform Tenant to RegisterTenantResponse", func(t *testing.T) {
		// Static datetime values
		createdAt := time.Date(2023, time.April, 10, 14, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2023, time.April, 11, 14, 0, 0, 0, time.UTC)

		// Prepare a tenant entity
		tenant, _ := entity.MakeTenant(
			entity.MakeTenantParams{
				ID:           1,
				Name:         "Acme Corporation",
				Address:      "123 Main St",
				Email:        "contact@acme.com",
				Phone:        "+1234567890",
				Timezone:     "America/New_York",
				OpeningHours: "09:00",
				ClosingHours: "18:00",
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
			},
		)

		// Call the function
		response := NewRegisterTenantResponse(tenant)

		// Assertions
		assert.Equal(t, uint(1), response.ID, "ID should match")
		assert.Equal(t, "Acme Corporation", response.Name, "Name should match")
		assert.Equal(t, "123 Main St", response.Address, "Address should match")
		assert.Equal(t, "contact@acme.com", response.Email, "Email should match")
		assert.Equal(t, "+1234567890", response.Phone, "Phone should match")
		assert.Equal(t, "America/New_York", response.Timezone, "Timezone should match")
		assert.Equal(t, "09:00", response.OpeningHours, "OpeningHours should match")
		assert.Equal(t, "18:00", response.ClosingHours, "ClosingHours should match")
		assert.Equal(t, createdAt, response.CreatedAt, "CreatedAt should match")
		assert.Equal(t, updatedAt, response.UpdatedAt, "UpdatedAt should match")
	})
}
