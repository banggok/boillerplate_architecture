package model

import (
	"appointment_management_system/internal/data/entity"
	"appointment_management_system/internal/pkg/custom_errors"
)

type Tenant struct {
	Model
	Name         string `gorm:"type:varchar(255);not null"`              // Tenant's name, cannot be null
	Address      string `gorm:"type:varchar(255)"`                       // Optional, maximum 255 characters
	Email        string `gorm:"type:varchar(255);not null;unique;index"` // Unique and indexed email
	Phone        string `gorm:"type:varchar(20);not null;unique;index"`  // unique and indexed
	Timezone     string `gorm:"type:varchar(100);not null"`              // Required timezone
	OpeningHours string `gorm:"type:time"`                               // Optional opening hours
	ClosingHours string `gorm:"type:time"`                               // Optional closing hours
	Accounts     *[]Account
}

func (Tenant) TableName() string {
	return "tenants"
}

func (Tenant) NotFoundError() custom_errors.ErrorCode {
	return custom_errors.TenantNotFound
}

func NewTenantModel(entity entity.Tenant) Tenant {
	tenant := Tenant{
		Model: Model{
			ID:        entity.GetID(),
			CreatedAt: entity.GetCreatedAt(),
			UpdatedAt: entity.GetUpdatedAt(),
		},
		Name:         entity.GetName(),
		Address:      entity.GetAddress(),
		Email:        entity.GetEmail(),
		Phone:        entity.GetPhone(),
		Timezone:     entity.GetTimezone(),
		OpeningHours: entity.GetOpeningHours(),
		ClosingHours: entity.GetClosingHours(),
	}
	if entity.GetAccounts() != nil {
		accountModels := make([]Account, 0)
		for _, account := range *entity.GetAccounts() {

			accountModels = append(accountModels, NewAccountModel(account))
		}
		tenant.Accounts = &accountModels
	}
	return tenant
}

func (m *Tenant) ToEntity() (*entity.Tenant, error) {
	param := entity.MakeTenantParams{
		Metadata: entity.Metadata{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		Name:         m.Name,
		Address:      m.Address,
		Email:        m.Email,
		Phone:        m.Phone,
		Timezone:     m.Timezone,
		OpeningHours: m.OpeningHours,
		ClosingHours: m.ClosingHours,
	}

	if m.Accounts != nil {
		accountEntities := make([]entity.Account, 0)
		for _, accountModel := range *m.Accounts {
			accountEntity, err := entity.MakeAccount(
				entity.MakeAccountParams{
					ID:        accountModel.ID,
					Name:      accountModel.Name,
					Email:     accountModel.Email,
					Phone:     accountModel.Phone,
					TenantId:  accountModel.TenantID,
					CreatedAt: accountModel.CreatedAt,
					UpdatedAt: accountModel.UpdatedAt,
				},
			)
			if err != nil {
				return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to convert account model to entity")
			}
			accountEntities = append(accountEntities, *accountEntity)
		}
		param.Account = &accountEntities
	}

	return entity.MakeTenant(param)
}
