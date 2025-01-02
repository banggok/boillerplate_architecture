package entity

import (
	"appointment_management_system/internal/config/validator"
	"appointment_management_system/internal/pkg/custom_errors"
	"appointment_management_system/internal/pkg/password"
)

type Account interface {
	Entity
	GetTenantId() uint
	GetName() string
	GetPassword() string
	GetEmail() string
	GetPhone() string
	GetTenant() Tenant
}

type accountImpl struct {
	entity
	name     string `validate:"required,alpha_space" example:"John Doe"`
	email    string `validate:"required,email" example:"admin@example.com"`
	phone    string `validate:"omitempty,e164" example:"+1234567890"`
	password string
	tenantId uint
	tenant   Tenant
}

type accountIdentity struct {
	name  string `validate:"required,alpha_space" example:"John Doe"`
	email string `validate:"required,email" example:"admin@example.com"`
	phone string `validate:"omitempty,e164" example:"+1234567890"`
}

func NewAccountIdentity(name, email, phone string) accountIdentity {
	return accountIdentity{
		name:  name,
		email: email,
		phone: phone,
	}
}

type accountTenant struct {
	tenantId uint
	tenant   Tenant
}

func NewAccountTenant(tenantId uint, tenant Tenant) accountTenant {
	return accountTenant{
		tenantId: tenantId,
		tenant:   tenant,
	}
}

type newAccountParams struct {
	accountIdentity
	tenant Tenant
}

type makeAccountParams struct {
	metadata
	accountIdentity
	accountTenant
	password string `validate:"required"`
}

func MakeAccount(
	metadata metadata,
	accountIdentity accountIdentity,
	accountTenant accountTenant,
	password string) (Account, error) {
	params := makeAccountParams{
		metadata:        metadata,
		accountIdentity: accountIdentity,
		accountTenant:   accountTenant,
		password:        password,
	}
	if err := validator.Validate.Struct(params); err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.AccountUnprocessEntity,
			"failed to validate make account entity")
	}

	return &accountImpl{
		entity: entity{
			id:        params.id,
			createdAt: params.createdAt,
			updatedAt: params.updatedAt,
		},
		name:     params.name,
		email:    params.email,
		phone:    params.phone,
		tenantId: params.tenantId,
		tenant:   params.tenant,
		password: params.password,
	}, nil
}

func NewAccount(accountIdentity accountIdentity, tenant Tenant) (Account, *string, error) {
	params := newAccountParams{
		accountIdentity: accountIdentity,
		tenant:          tenant,
	}
	if err := validator.Validate.Struct(params); err != nil {
		return nil, nil, custom_errors.New(
			err,
			custom_errors.AccountUnprocessEntity,
			"failed to validate new account entity")
	}

	plain, hashed, err := password.GeneratePassword(8)
	if err != nil || plain == nil || hashed == nil {
		return nil, nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to generate password")
	}

	account := accountImpl{
		name:     params.name,
		email:    params.email,
		phone:    params.phone,
		password: *hashed,
	}

	if params.tenant != nil {
		account.tenantId = params.tenant.GetID()
		account.tenant = params.tenant
	}

	return &account, plain, nil
}

// GetTenantId retrieves the id of the tenant account.
func (a *accountImpl) GetTenantId() uint {
	return a.tenantId
}

// GetName retrieves the name of the account.
func (a *accountImpl) GetName() string {
	return a.name
}

// GetPassword retrieves the hash password of the account.
func (a *accountImpl) GetPassword() string {
	return a.password
}

// GetEmail retrieves the email of the account.
func (a *accountImpl) GetEmail() string {
	return a.email
}

// GetPhone retrieves the phone of the account.
func (a *accountImpl) GetPhone() string {
	return a.phone
}

// GetTenant retrieves the tenant of the account.
func (a *accountImpl) GetTenant() Tenant {
	return a.tenant
}
