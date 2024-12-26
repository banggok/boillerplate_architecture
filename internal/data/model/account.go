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
	Password  string    `gorm:"type:text;not null"`                            // Hashed password
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
		Password:  entity.GetPassword(),
	}
	if entity.GetTenant() != nil {
		tenant := NewTenantModel(*entity.GetTenant())
		account.Tenant = &tenant
	}
	return account
}

func (m *Account) ToEntity() (*entity.Account, error) {
	var tenantParam *entity.Tenant
	if m.Tenant != nil {
		var err error
		tenantParam, err = m.Tenant.ToEntity()
		if err != nil {
			return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed convert account model to entity")
		}
	}
	return entity.MakeAccount(
		entity.NewMetadata(m.ID, m.CreatedAt, m.UpdatedAt),
		entity.NewAccountIdentity(m.Name, m.Email, m.Phone),
		entity.NewAccountTenant(m.TenantID, tenantParam),
		m.Password,
	)
}
