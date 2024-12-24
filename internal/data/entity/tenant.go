package entity

import (
	"appointment_management_system/internal/pkg/custom_errors"
	"time"
)

type Tenant struct {
	entity
	name         string `validate:"required"`              // Tenant's name (required, no strict restrictions)
	address      string `validate:"omitempty,max=255"`     // Optional, with a maximum of 255 characters
	email        string `validate:"required,email"`        // Tenant's contact email (required)
	phone        string `validate:"omitempty,e164"`        // Tenant's phone number (optional)
	timezone     string `validate:"required,iana_tz"`      // Tenant's timezone (required, IANA format)
	openingHours string `validate:"omitempty,time_format"` // Tenant's opening hours in HH:mm format (optional)
	closingHours string `validate:"omitempty,time_format"` // Tenant's closing hours in HH:mm format (optional)
	account      *[]Account
}

type Metadata struct {
	ID        uint      `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
}

type NewTenantParams struct {
	Name         string `validate:"required"`              // Tenant's name (required, no strict restrictions)
	Address      string `validate:"omitempty,max=255"`     // Optional, with a maximum of 255 characters
	Email        string `validate:"required,email"`        // Tenant's contact email (required)
	Phone        string `validate:"omitempty,e164"`        // Tenant's phone number (optional)
	Timezone     string `validate:"required,iana_tz"`      // Tenant's timezone (required, IANA format)
	OpeningHours string `validate:"omitempty,time_format"` // Tenant's opening hours in HH:mm format (optional)
	ClosingHours string `validate:"omitempty,time_format"` // Tenant's closing hours in HH:mm format (optional)
	Account      *[]Account
}

type MakeTenantParams struct {
	Metadata
	Name         string `validate:"required"`
	Address      string `validate:"required,max=255"`
	Email        string `validate:"required,email"`
	Phone        string `validate:"required,e164"`
	Timezone     string `validate:"required,iana_tz"`
	OpeningHours string `validate:"required,time_format"`
	ClosingHours string `validate:"required,time_format"`
	Account      *[]Account
}

func NewTenant(params NewTenantParams) (*Tenant, error) {
	if err := validate.Struct(params); err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.TenantUnprocessEntity,
			"failed to validate new tenant entity")
	}

	return &Tenant{
		name:         params.Name,
		address:      params.Address,
		email:        params.Email,
		phone:        params.Phone,
		timezone:     params.Timezone,
		openingHours: params.OpeningHours,
		closingHours: params.ClosingHours,
		account:      params.Account,
	}, nil
}

func MakeTenant(params MakeTenantParams) (*Tenant, error) {
	// Validate the input parameters

	if err := validate.Struct(params); err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.TenantUnprocessEntity,
			"failed to validate make tenant")
	}

	return &Tenant{
		entity: entity{
			id:        params.Metadata.ID,
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
		account:      params.Account,
	}, nil
}

// Getter and Setter for Name
func (t *Tenant) GetName() string {
	return t.name
}

// Getter and Setter for Address
func (t *Tenant) GetAddress() string {
	return t.address
}

// Getter and Setter for Email
func (t *Tenant) GetEmail() string {
	return t.email
}

// Getter and Setter for Phone
func (t *Tenant) GetPhone() string {
	return t.phone
}

// Getter and Setter for Timezone
func (t *Tenant) GetTimezone() string {
	return t.timezone
}

// Getter and Setter for OpeningHours
func (t *Tenant) GetOpeningHours() string {
	return t.openingHours
}

// Getter and Setter for ClosingHours
func (t *Tenant) GetClosingHours() string {
	return t.closingHours
}

func (t *Tenant) GetAccounts() *[]Account {
	return t.account
}
