package repository

import (
	"appointment_management_system/internal/domain/tenants/entity"
	"time"
)

type tenantModel struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`                // Primary key, auto-incrementing
	Name         string    `gorm:"type:varchar(255);not null"`              // Tenant's name, cannot be null
	Address      string    `gorm:"type:varchar(255)"`                       // Optional, maximum 255 characters
	Email        string    `gorm:"type:varchar(255);not null;unique;index"` // Unique and indexed email
	Phone        string    `gorm:"type:varchar(20);not null;unique;index"`  // unique and indexed
	Timezone     string    `gorm:"type:varchar(100);not null"`              // Required timezone
	OpeningHours string    `gorm:"type:time"`                               // Optional opening hours
	ClosingHours string    `gorm:"type:time"`                               // Optional closing hours
	CreatedAt    time.Time `gorm:"autoCreateTime"`                          // Automatically set when the record is created
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`                          // Automatically set when the record is updated
}

func (tenantModel) TableName() string {
	return "tenants"
}

func newTenantModel(entity entity.Tenant) tenantModel {
	return tenantModel{
		ID:           entity.GetID(),
		Name:         entity.GetName(),
		Address:      entity.GetAddress(),
		Email:        entity.GetEmail(),
		Phone:        entity.GetPhone(),
		Timezone:     entity.GetTimezone(),
		OpeningHours: entity.GetOpeningHours(),
		ClosingHours: entity.GetClosingHours(),
		CreatedAt:    entity.GetCreatedAt().UTC(),
		UpdatedAt:    entity.GetUpdatedAt().UTC(),
	}
}

func (m *tenantModel) ToEntity() (*entity.Tenant, error) {
	param := entity.MakeTenantParams{
		ID:           m.ID,
		Name:         m.Name,
		Address:      m.Address,
		Email:        m.Email,
		Phone:        m.Phone,
		Timezone:     m.Timezone,
		OpeningHours: m.OpeningHours,
		ClosingHours: m.ClosingHours,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}

	return entity.MakeTenant(param)
}
