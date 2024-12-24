package entity

import (
	"appointment_management_system/internal/pkg/custom_errors"
	"time"
)

type Account struct {
	entity
	name     string `validate:"required,alpha_space" example:"John Doe"`
	email    string `validate:"required,email" example:"admin@example.com"`
	phone    string `validate:"omitempty,e164" example:"+1234567890"`
	tenantId uint
	tenant   *Tenant
}

type NewAccountParams struct {
	Name   string `validate:"required,alpha_space" example:"John Doe"`
	Email  string `validate:"required,email" example:"admin@example.com"`
	Phone  string `validate:"omitempty,e164" example:"+1234567890"`
	Tenant *Tenant
}

type MakeAccountParams struct {
	ID        uint
	Name      string `validate:"required,alpha_space" example:"John Doe"`
	Email     string `validate:"required,email" example:"admin@example.com"`
	Phone     string `validate:"omitempty,e164" example:"+1234567890"`
	TenantId  uint   `validate:"required"`
	Tenant    *Tenant
	CreatedAt time.Time
	UpdatedAt time.Time
}

func MakeAccount(params MakeAccountParams) (*Account, error) {
	if err := validate.Struct(params); err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.AccountUnprocessEntity,
			"failed to validate make account entity")
	}

	return &Account{
		entity: entity{
			id:        params.ID,
			createdAt: params.CreatedAt,
			updatedAt: params.UpdatedAt,
		},
		name:     params.Name,
		email:    params.Email,
		phone:    params.Phone,
		tenantId: params.TenantId,
		tenant:   params.Tenant,
	}, nil
}

func NewAccount(params NewAccountParams) (*Account, error) {
	if err := validate.Struct(params); err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.AccountUnprocessEntity,
			"failed to validate new account entity")
	}

	return &Account{
		entity: entity{},
		name:   params.Name,
		email:  params.Email,
		phone:  params.Phone,
	}, nil
}

// GetTenantId retrieves the id of the tenant account.
func (a *Account) GetTenantId() uint {
	return a.tenantId
}

// GetName retrieves the name of the account.
func (a *Account) GetName() string {
	return a.name
}

// GetEmail retrieves the email of the account.
func (a *Account) GetEmail() string {
	return a.email
}

// GetPhone retrieves the phone of the account.
func (a *Account) GetPhone() string {
	return a.phone
}

// GetTenant retrieves the tenant of the account.
func (a *Account) GetTenant() *Tenant {
	return a.tenant
}
