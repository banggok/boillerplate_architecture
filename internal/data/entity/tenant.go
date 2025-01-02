package entity

import (
	"appointment_management_system/internal/config/validator"
	"appointment_management_system/internal/pkg/custom_errors"
)

type Tenant interface {
	Entity
	GetName() string
	GetAddress() string
	GetEmail() string
	GetAccounts() *[]Account
	GetPhone() string
	GetTimezone() string
	GetOpeningHours() string
	GetClosingHours() string
}

type tenantImpl struct {
	entity
	name         string `validate:"required"`              // Tenant's name (required, no strict restrictions)
	address      string `validate:"omitempty,max=255"`     // Optional, with a maximum of 255 characters
	email        string `validate:"required,email"`        // Tenant's contact email (required)
	phone        string `validate:"omitempty,e164"`        // Tenant's phone number (optional)
	timezone     string `validate:"required,iana_tz"`      // Tenant's timezone (required, IANA format)
	openingHours string `validate:"omitempty,time_format"` // Tenant's opening hours in HH:mm format (optional)
	closingHours string `validate:"omitempty,time_format"` // Tenant's closing hours in HH:mm format (optional)
	accounts     *[]Account
}

type tenantIdentity struct {
	name  string `validate:"required"`       // Tenant's name (required, no strict restrictions)
	email string `validate:"required,email"` // Tenant's contact email (required)
	phone string `validate:"omitempty,e164"` // Tenant's phone number (optional)
}

func NewTenantIdentity(name, email, phone string) tenantIdentity {
	return tenantIdentity{
		name:  name,
		email: email,
		phone: phone,
	}
}

type tenantStoreInfo struct {
	address      string `validate:"omitempty,max=255"`     // Optional, with a maximum of 255 characters
	timezone     string `validate:"required,iana_tz"`      // Tenant's timezone (required, IANA format)
	openingHours string `validate:"omitempty,time_format"` // Tenant's opening hours in HH:mm format (optional)
	closingHours string `validate:"omitempty,time_format"` // Tenant's closing hours in HH:mm format (optional)
}

func NewTenantStoreInfo(address, timezone, openingHours, closingHours string) tenantStoreInfo {
	return tenantStoreInfo{
		address:      address,
		timezone:     timezone,
		openingHours: openingHours,
		closingHours: closingHours,
	}
}

type newTenantParams struct {
	tenantIdentity
	tenantStoreInfo
	accounts *[]Account
}

type makeTenantParams struct {
	metadata
	tenantIdentity
	tenantStoreInfo
	accounts *[]Account
}

func NewTenant(tenantIdentity tenantIdentity, tenantStoreInfo tenantStoreInfo, accounts *[]Account) (Tenant, error) {
	params := newTenantParams{
		tenantIdentity:  tenantIdentity,
		tenantStoreInfo: tenantStoreInfo,
		accounts:        accounts,
	}
	if err := validator.Validate.Struct(params); err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.TenantUnprocessEntity,
			"failed to validate new tenant entity")
	}

	return &tenantImpl{
		name:         params.name,
		address:      params.address,
		email:        params.email,
		phone:        params.phone,
		timezone:     params.timezone,
		openingHours: params.openingHours,
		closingHours: params.closingHours,
		accounts:     params.accounts,
	}, nil
}

func MakeTenant(metadata metadata,
	tenantIdentity tenantIdentity,
	tenantStoreInfo tenantStoreInfo,
	accounts *[]Account) (Tenant, error) {
	params := makeTenantParams{
		metadata:        metadata,
		tenantIdentity:  tenantIdentity,
		tenantStoreInfo: tenantStoreInfo,
		accounts:        accounts,
	}
	// Validate the input parameters
	if err := validator.Validate.Struct(params); err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.TenantUnprocessEntity,
			"failed to validate make tenant")
	}

	return &tenantImpl{
		entity: entity{
			id:        params.id,
			createdAt: params.createdAt,
			updatedAt: params.updatedAt,
		},
		name:         params.name,
		address:      params.address,
		email:        params.email,
		phone:        params.phone,
		timezone:     params.timezone,
		openingHours: params.openingHours,
		closingHours: params.closingHours,
		accounts:     params.accounts,
	}, nil
}

// Getter and Setter for Name
func (t *tenantImpl) GetName() string {
	return t.name
}

// Getter and Setter for Address
func (t *tenantImpl) GetAddress() string {
	return t.address
}

// Getter and Setter for Email
func (t *tenantImpl) GetEmail() string {
	return t.email
}

// Getter and Setter for Phone
func (t *tenantImpl) GetPhone() string {
	return t.phone
}

// Getter and Setter for Timezone
func (t *tenantImpl) GetTimezone() string {
	return t.timezone
}

// Getter and Setter for OpeningHours
func (t *tenantImpl) GetOpeningHours() string {
	return t.openingHours
}

// Getter and Setter for ClosingHours
func (t *tenantImpl) GetClosingHours() string {
	return t.closingHours
}

func (t *tenantImpl) GetAccounts() *[]Account {
	return t.accounts
}
