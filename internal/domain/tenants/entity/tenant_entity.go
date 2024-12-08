package entity

import (
	"appointment_management_system/internal/config"
	"fmt"
	"time"
)

type Tenant struct {
	id           uint      // id from database
	name         string    `validate:"required"`              // Tenant's name (required, no strict restrictions)
	address      string    `validate:"omitempty,max=255"`     // Optional, with a maximum of 255 characters
	email        string    `validate:"required,email"`        // Tenant's contact email (required)
	phone        string    `validate:"omitempty,e164"`        // Tenant's phone number (optional)
	timezone     string    `validate:"required,iana_tz"`      // Tenant's timezone (required, IANA format)
	openingHours string    `validate:"omitempty,time_format"` // Tenant's opening hours in HH:mm format (optional)
	closingHours string    `validate:"omitempty,time_format"` // Tenant's closing hours in HH:mm format (optional)
	createdAt    time.Time // timestamp from database
	updatedAt    time.Time // timestamp from database
}

type NewTenantParams struct {
	Name         string `validate:"required"`              // Tenant's name (required, no strict restrictions)
	Address      string `validate:"omitempty,max=255"`     // Optional, with a maximum of 255 characters
	Email        string `validate:"required,email"`        // Tenant's contact email (required)
	Phone        string `validate:"omitempty,e164"`        // Tenant's phone number (optional)
	Timezone     string `validate:"required,iana_tz"`      // Tenant's timezone (required, IANA format)
	OpeningHours string `validate:"omitempty,time_format"` // Tenant's opening hours in HH:mm format (optional)
	ClosingHours string `validate:"omitempty,time_format"` // Tenant's closing hours in HH:mm format (optional)
}

type MakeTenantParams struct {
	ID           uint      `validate:"required"`
	Name         string    `validate:"required"`
	Address      string    `validate:"required,max=255"`
	Email        string    `validate:"required,email"`
	Phone        string    `validate:"required,e164"`
	Timezone     string    `validate:"required,iana_tz"`
	OpeningHours string    `validate:"required,time_format"`
	ClosingHours string    `validate:"required,time_format"`
	CreatedAt    time.Time `validate:"required"`
	UpdatedAt    time.Time `validate:"required"`
}

var validate = config.SetupValidator()

func NewTenant(params NewTenantParams) (*Tenant, error) {
	if err := validate.Struct(params); err != nil {
		return nil, fmt.Errorf("failed NewTenant validate: %w", err)
	}

	return &Tenant{
		name:         params.Name,
		address:      params.Address,
		email:        params.Email,
		phone:        params.Phone,
		timezone:     params.Timezone,
		openingHours: params.OpeningHours,
		closingHours: params.ClosingHours,
	}, nil
}

func MakeTenant(params MakeTenantParams) (*Tenant, error) {
	// Validate the input parameters

	if err := validate.Struct(params); err != nil {
		return nil, fmt.Errorf("failed MakeTenant validate: %w", err)
	}

	return &Tenant{
		id:           params.ID,
		name:         params.Name,
		address:      params.Address,
		email:        params.Email,
		phone:        params.Phone,
		timezone:     params.Timezone,
		openingHours: params.OpeningHours,
		closingHours: params.ClosingHours,
		createdAt:    params.CreatedAt,
		updatedAt:    params.UpdatedAt,
	}, nil
}

// Getter and Setter for ID
func (t *Tenant) GetID() uint {
	return t.id
}

func (t *Tenant) SetID(id uint) {
	t.id = id
}

// Getter and Setter for Name
func (t *Tenant) GetName() string {
	return t.name
}

func (t *Tenant) SetName(name string) {
	t.name = name
}

// Getter and Setter for Address
func (t *Tenant) GetAddress() string {
	return t.address
}

func (t *Tenant) SetAddress(address string) {
	t.address = address
}

// Getter and Setter for Email
func (t *Tenant) GetEmail() string {
	return t.email
}

func (t *Tenant) SetEmail(email string) {
	t.email = email
}

// Getter and Setter for Phone
func (t *Tenant) GetPhone() string {
	return t.phone
}

func (t *Tenant) SetPhone(phone string) {
	t.phone = phone
}

// Getter and Setter for Timezone
func (t *Tenant) GetTimezone() string {
	return t.timezone
}

func (t *Tenant) SetTimezone(timezone string) {
	t.timezone = timezone
}

// Getter and Setter for OpeningHours
func (t *Tenant) GetOpeningHours() string {
	return t.openingHours
}

func (t *Tenant) SetOpeningHours(openingHours string) {
	t.openingHours = openingHours
}

// Getter and Setter for ClosingHours
func (t *Tenant) GetClosingHours() string {
	return t.closingHours
}

func (t *Tenant) SetClosingHours(closingHours string) {
	t.closingHours = closingHours
}

// Getter and Setter for CreatedAt
func (t *Tenant) GetCreatedAt() time.Time {
	return t.createdAt
}

func (t *Tenant) SetCreatedAt(createdAt time.Time) {
	t.createdAt = createdAt
}

// Getter and Setter for UpdatedAt
func (t *Tenant) GetUpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Tenant) SetUpdatedAt(updatedAt time.Time) {
	t.updatedAt = updatedAt
}
