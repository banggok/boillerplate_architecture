package entity

import (
	"github.com/banggok/boillerplate_architecture/internal/config/validator"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
)

type Tenant interface {
	iMetadata
	Name() string
	Address() string
	Email() string
	Accounts() *[]Account
	Phone() string
	Timezone() string
	OpeningHours() string
	ClosingHours() string
}

type tenantImpl struct {
	metadataImpl
	name         string // Tenant's name (required, no strict restrictions)
	address      string // Optional, with a maximum of 255 characters
	email        string // Tenant's contact email (required)
	phone        string // Tenant's phone number (optional)
	timezone     string // Tenant's timezone (required, IANA format)
	openingHours string // Tenant's opening hours in HH:mm format (optional)
	closingHours string // Tenant's closing hours in HH:mm format (optional)
	accounts     *[]Account
}

type tenantIdentity struct {
	Name  string `validate:"required"`       // Tenant's name (required, no strict restrictions)
	Email string `validate:"required,email"` // Tenant's contact email (required)
	Phone string `validate:"omitempty,e164"` // Tenant's phone number (optional)
}

func NewTenantIdentity(name, email, phone string) tenantIdentity {
	return tenantIdentity{
		Name:  name,
		Email: email,
		Phone: phone,
	}
}

type tenantStoreInfo struct {
	Address      string `validate:"omitempty,max=255"`     // Optional, with a maximum of 255 characters
	Timezone     string `validate:"required,iana_tz"`      // Tenant's timezone (required, IANA format)
	OpeningHours string `validate:"omitempty,time_format"` // Tenant's opening hours in HH:mm format (optional)
	ClosingHours string `validate:"omitempty,time_format"` // Tenant's closing hours in HH:mm format (optional)
}

func NewTenantStoreInfo(address, timezone, openingHours, closingHours string) tenantStoreInfo {
	return tenantStoreInfo{
		Address:      address,
		Timezone:     timezone,
		OpeningHours: openingHours,
		ClosingHours: closingHours,
	}
}

type newTenantParams struct {
	tenantIdentity
	tenantStoreInfo
	Accounts *[]Account
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
		Accounts:        accounts,
	}
	if err := validator.Validate.Struct(params); err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.TenantUnprocessEntity,
			"failed to validate new tenant entity")
	}

	return &tenantImpl{
		name:         params.Name,
		address:      params.Address,
		email:        params.Email,
		phone:        params.Phone,
		timezone:     params.Timezone,
		openingHours: params.OpeningHours,
		closingHours: params.ClosingHours,
		accounts:     params.Accounts,
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
		metadataImpl: metadataImpl{
			id:        params.ID,
			createdAt: params.CreatedAt,
			updatedAt: params.UpdatedAt,
		},
		name:         params.Name,
		address:      params.Address,
		email:        params.Email,
		phone:        params.Phone,
		timezone:     params.Timezone,
		openingHours: params.OpeningHours,
		closingHours: params.ClosingHours,
		accounts:     params.accounts,
	}, nil
}

// Getter and Setter for Name
func (t *tenantImpl) Name() string {
	return t.name
}

// Getter and Setter for Address
func (t *tenantImpl) Address() string {
	return t.address
}

// Getter and Setter for Email
func (t *tenantImpl) Email() string {
	return t.email
}

// Getter and Setter for Phone
func (t *tenantImpl) Phone() string {
	return t.phone
}

// Getter and Setter for Timezone
func (t *tenantImpl) Timezone() string {
	return t.timezone
}

// Getter and Setter for OpeningHours
func (t *tenantImpl) OpeningHours() string {
	return t.openingHours
}

// Getter and Setter for ClosingHours
func (t *tenantImpl) ClosingHours() string {
	return t.closingHours
}

func (t *tenantImpl) Accounts() *[]Account {
	return t.accounts
}
