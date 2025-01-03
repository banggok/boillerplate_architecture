package model

import (
	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
)

// Account struct for the accounts table
type Account struct {
	Metadata
	Name                 string  `gorm:"type:varchar(100);not null"`
	Email                string  `gorm:"type:varchar(255);unique;not null"`
	Phone                string  `gorm:"type:varchar(20);unique;not null"`
	Password             string  `gorm:"type:text;not null"`                            // Hashed password
	TenantID             uint    `gorm:"not null;index"`                                // Foreign key
	Tenant               *Tenant `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // Relationship
	AccountVerifications *[]AccountVerification
}

func (Account) TableName() string {
	return "accounts"
}

func (Account) NotFoundError() custom_errors.ErrorCode {
	return custom_errors.AccountNotFound
}

func NewAccountModel(entity entity.Account) Account {
	account := Account{
		Metadata: Metadata{
			ID:        entity.ID(),
			CreatedAt: entity.CreatedAt(),
			UpdatedAt: entity.UpdatedAt(),
		},
		Name:     entity.Name(),
		Email:    entity.Email(),
		Phone:    entity.Phone(),
		TenantID: entity.TenantId(),
		Password: entity.Password(),
	}
	if entity.Tenant() != nil {
		tenant := NewTenantModel(entity.Tenant())
		account.Tenant = &tenant
	}

	if entity.AccountVerifications() != nil {
		accountVerificationModels := make([]AccountVerification, 0)

		for _, accountVerification := range *entity.AccountVerifications() {
			accountVerificationModels = append(accountVerificationModels, NewAccountVerification(accountVerification))
		}
		account.AccountVerifications = &accountVerificationModels
	}

	return account
}

func (m *Account) ToEntity() (entity.Account, error) {
	var tenantParam entity.Tenant
	if m.Tenant != nil {
		var err error
		tenantParam, err = m.Tenant.ToEntity()
		if err != nil {
			return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed convert account model to entity")
		}
	}

	var accountVerifications *[]entity.AccountVerification

	if m.AccountVerifications != nil {
		accountVerificationEntities := make([]entity.AccountVerification, 0)
		for _, accountVerification := range *m.AccountVerifications {
			accountVerificationEntity, err := accountVerification.ToEntity()
			if err != nil {
				return nil, custom_errors.New(err, custom_errors.InternalServerError,
					"failed to convert account verification model to entity")
			}
			accountVerificationEntities = append(accountVerificationEntities, accountVerificationEntity)
		}
		accountVerifications = &accountVerificationEntities
	}

	return entity.MakeAccount(
		entity.NewMetadata(m.ID, m.CreatedAt, m.UpdatedAt),
		entity.NewAccountIdentity(m.Name, m.Email, m.Phone, &m.Password),
		entity.NewAccountTenant(m.TenantID, tenantParam),
		accountVerifications,
	)
}
