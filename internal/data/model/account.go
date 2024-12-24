package model

import (
	"appointment_management_system/internal/data/entity"
	"appointment_management_system/internal/pkg/custom_errors"
	"time"
)

// Account struct for the accounts table
type Account struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(255);unique;not null"`
	Phone     string    `gorm:"type:varchar(20);unique;not null"`
	TenantID  uint      `gorm:"not null;index"`                                // Foreign key
	Tenant    *Tenant   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // Relationship
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Account) TableName() string {
	return "accounts"
}

func (Account) NotFoundError() custom_errors.ErrorCode {
	return custom_errors.AccountNotFound
}

func NewAccountModel(entity entity.Account) Account {
	account := Account{
		ID:        entity.GetID(),
		Name:      entity.GetName(),
		Email:     entity.GetEmail(),
		Phone:     entity.GetPhone(),
		TenantID:  entity.GetTenantId(),
		CreatedAt: entity.GetCreatedAt(),
		UpdatedAt: entity.GetUpdatedAt(),
	}
	if entity.GetTenant() != nil {
		tenant := NewTenantModel(*entity.GetTenant())
		account.Tenant = &tenant
	}
	return account
}

func (m *Account) ToEntity() (*entity.Account, error) {
	accountParams := entity.MakeAccountParams{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Phone:     m.Phone,
		TenantId:  m.TenantID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.Tenant != nil {
		tenant, err := m.Tenant.ToEntity()
		if err != nil {
			return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed convert account model to entity")
		}
		accountParams.Tenant = tenant
	}
	return entity.MakeAccount(accountParams)
}
