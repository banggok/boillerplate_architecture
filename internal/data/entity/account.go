package entity

import (
	"github.com/banggok/boillerplate_architecture/internal/config/validator"

	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/password"
)

type Account interface {
	Entity
	TenantId() uint
	Name() string
	Password() string
	Email() string
	Phone() string
	Tenant() Tenant
	AccountVerifications() *[]AccountVerification
}

type accountImpl struct {
	entity
	name                 string `validate:"required,alpha_space" example:"John Doe"`
	email                string `validate:"required,email" example:"admin@example.com"`
	phone                string `validate:"omitempty,e164" example:"+1234567890"`
	password             string
	tenantId             uint
	tenant               Tenant
	accountVerifications *[]AccountVerification
}

// AccountVerifications implements Account.
func (a *accountImpl) AccountVerifications() *[]AccountVerification {
	return a.accountVerifications
}

type accountIdentity struct {
	name     string  `validate:"required,alpha_space" example:"John Doe"`
	email    string  `validate:"required,email" example:"admin@example.com"`
	phone    string  `validate:"omitempty,e164" example:"+1234567890"`
	password *string `validate:"required"`
}

func NewAccountIdentity(name, email, phone string, password *string) accountIdentity {
	return accountIdentity{
		name:     name,
		email:    email,
		phone:    phone,
		password: password,
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
}

func MakeAccount(
	metadata metadata,
	accountIdentity accountIdentity,
	accountTenant accountTenant,
	accountVerifications *[]AccountVerification) (Account, error) {
	params := makeAccountParams{
		metadata:        metadata,
		accountIdentity: accountIdentity,
		accountTenant:   accountTenant,
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
		name:                 params.name,
		email:                params.email,
		phone:                params.phone,
		tenantId:             params.tenantId,
		tenant:               params.tenant,
		password:             *params.password,
		accountVerifications: accountVerifications,
	}, nil
}

func NewAccount(accountIdentity accountIdentity, tenant Tenant, accountVerifications *[]AccountVerification) (Account, *string, error) {
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

	plain, err := password.GeneratePassword(8)
	if err != nil || plain == nil {
		return nil, nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to generate password")
	}

	hashed, err := password.HashPassword(*plain)
	if err != nil || hashed == nil {
		return nil, nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to hash password")
	}

	account := accountImpl{
		name:                 params.name,
		email:                params.email,
		phone:                params.phone,
		password:             *hashed,
		accountVerifications: accountVerifications,
	}

	if params.tenant != nil {
		account.tenantId = params.tenant.ID()
		account.tenant = params.tenant
	}

	return &account, plain, nil
}

// TenantId retrieves the id of the tenant account.
func (a *accountImpl) TenantId() uint {
	return a.tenantId
}

// Name retrieves the name of the account.
func (a *accountImpl) Name() string {
	return a.name
}

// Password retrieves the hash password of the account.
func (a *accountImpl) Password() string {
	return a.password
}

// Email retrieves the email of the account.
func (a *accountImpl) Email() string {
	return a.email
}

// Phone retrieves the phone of the account.
func (a *accountImpl) Phone() string {
	return a.phone
}

// Tenant retrieves the tenant of the account.
func (a *accountImpl) Tenant() Tenant {
	return a.tenant
}
