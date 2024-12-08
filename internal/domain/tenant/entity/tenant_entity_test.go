package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMakeTenant(t *testing.T) {
	// Test data
	now := time.Now()
	createdAt := now.Add(-24 * time.Hour)
	updatedAt := now

	tenant, _ := MakeTenant(
		MakeTenantParams{
			ID:           uint(1),
			Name:         "John Doe",
			Address:      "123 Elm Street",
			Email:        "johndoe@example.com",
			Phone:        "+1234567890",
			Timezone:     "America/New_York",
			OpeningHours: "09:00",
			ClosingHours: "17:00",
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		},
	)

	// Table-driven test cases
	testCases := []struct {
		field       string
		expected    interface{}
		actual      func() interface{}
		description string
	}{
		{"ID", uint(1), func() interface{} { return tenant.GetID() }, "Tenant ID mismatch"},
		{"Name", "John Doe", func() interface{} { return tenant.GetName() }, "Tenant name mismatch"},
		{"Address", "123 Elm Street", func() interface{} { return tenant.GetAddress() }, "Tenant address mismatch"},
		{"Email", "johndoe@example.com", func() interface{} { return tenant.GetEmail() }, "Tenant email mismatch"},
		{"Phone", "+1234567890", func() interface{} { return tenant.GetPhone() }, "Tenant phone mismatch"},
		{"Timezone", "America/New_York", func() interface{} { return tenant.GetTimezone() }, "Tenant timezone mismatch"},
		{"OpeningHours", "09:00", func() interface{} { return tenant.GetOpeningHours() }, "Tenant opening hours mismatch"},
		{"ClosingHours", "17:00", func() interface{} { return tenant.GetClosingHours() }, "Tenant closing hours mismatch"},
		{"CreatedAt", createdAt, func() interface{} { return tenant.GetCreatedAt() }, "Tenant createdAt mismatch"},
		{"UpdatedAt", updatedAt, func() interface{} { return tenant.GetUpdatedAt() }, "Tenant updatedAt mismatch"},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.field, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.actual(), tc.description)
		})
	}
}

func TestTenantGettersAndSetters(t *testing.T) {
	tenant := &Tenant{}

	// Test ID
	tenant.SetID(10)
	assert.Equal(t, uint(10), tenant.GetID(), "Tenant ID mismatch")

	// Test Name
	tenant.SetName("Acme Corp")
	assert.Equal(t, "Acme Corp", tenant.GetName(), "Tenant name mismatch")

	// Test Address
	tenant.SetAddress("456 Oak Avenue")
	assert.Equal(t, "456 Oak Avenue", tenant.GetAddress(), "Tenant address mismatch")

	// Test Email
	tenant.SetEmail("contact@acme.com")
	assert.Equal(t, "contact@acme.com", tenant.GetEmail(), "Tenant email mismatch")

	// Test Phone
	tenant.SetPhone("+1987654321")
	assert.Equal(t, "+1987654321", tenant.GetPhone(), "Tenant phone mismatch")

	// Test Timezone
	tenant.SetTimezone("Asia/Jakarta")
	assert.Equal(t, "Asia/Jakarta", tenant.GetTimezone(), "Tenant timezone mismatch")

	// Test OpeningHours
	tenant.SetOpeningHours("08:00")
	assert.Equal(t, "08:00", tenant.GetOpeningHours(), "Tenant opening hours mismatch")

	// Test ClosingHours
	tenant.SetClosingHours("18:00")
	assert.Equal(t, "18:00", tenant.GetClosingHours(), "Tenant closing hours mismatch")

	// Test CreatedAt
	now := time.Now()
	tenant.SetCreatedAt(now)
	assert.Equal(t, now, tenant.GetCreatedAt(), "Tenant createdAt mismatch")

	// Test UpdatedAt
	updated := now.Add(1 * time.Hour)
	tenant.SetUpdatedAt(updated)
	assert.Equal(t, updated, tenant.GetUpdatedAt(), "Tenant updatedAt mismatch")
}

func TestMakeTenantValidation(t *testing.T) {
	// Common test data
	now := time.Now()

	t.Run("Validations", func(t *testing.T) {
		testValidTenant(t, now)
		testMissingRequiredFields(t, now)
		testInvalidFields(t, now)
	})
}

func testValidTenant(t *testing.T, now time.Time) {
	params := MakeTenantParams{
		ID:           1,
		Name:         "Acme Corporation",
		Address:      "123 Main St",
		Email:        "contact@acme.com",
		Phone:        "+1234567890",
		Timezone:     "America/New_York",
		OpeningHours: "09:00",
		ClosingHours: "18:00",
		CreatedAt:    now.Add(-24 * time.Hour),
		UpdatedAt:    now,
	}

	tenant, err := MakeTenant(params)
	assert.NoError(t, err, "Expected no error for valid tenant data")
	assert.NotNil(t, tenant, "Tenant should be created successfully")
}

func testMissingRequiredFields(t *testing.T, now time.Time) {
	params := MakeTenantParams{
		ID:           1,
		Address:      "123 Main St",
		Email:        "contact@acme.com",
		Phone:        "+1234567890",
		Timezone:     "America/New_York",
		OpeningHours: "09:00",
		ClosingHours: "18:00",
		CreatedAt:    now.Add(-24 * time.Hour),
		UpdatedAt:    now,
	}

	_, err := MakeTenant(params)
	assert.Error(t, err, "Expected error for missing required field")
	assert.Contains(t, err.Error(), "Name", "Error should be related to the Name field")
}

func testInvalidFields(t *testing.T, now time.Time) {
	tests := getInvalidFieldTests(now)

	for _, test := range tests {
		runInvalidFieldTest(t, test.name, test.params, test.errorField)
	}
}

func getInvalidFieldTests(now time.Time) []struct {
	name       string
	params     MakeTenantParams
	errorField string
} {
	return []struct {
		name       string
		params     MakeTenantParams
		errorField string
	}{
		{
			name:       "Invalid Email Format",
			params:     createBaseParams(now, "invalid-email", "+1234567890", "123 Main St"),
			errorField: "Email",
		},
		{
			name:       "Invalid Phone Format",
			params:     createBaseParams(now, "contact@acme.com", "1234567890", "123 Main St"),
			errorField: "Phone",
		},
		{
			name:       "Address Exceeds Max Length",
			params:     createBaseParams(now, "contact@acme.com", "+1234567890", string(make([]byte, 256))),
			errorField: "Address",
		},
	}
}

func createBaseParams(now time.Time, email, phone, address string) MakeTenantParams {
	return MakeTenantParams{
		ID:           1,
		Name:         "Acme Corporation",
		Address:      address,
		Email:        email,
		Phone:        phone,
		Timezone:     "America/New_York",
		OpeningHours: "09:00",
		ClosingHours: "18:00",
		CreatedAt:    now.Add(-24 * time.Hour),
		UpdatedAt:    now,
	}
}

func runInvalidFieldTest(t *testing.T, name string, params MakeTenantParams, errorField string) {
	t.Run(name, func(t *testing.T) {
		_, err := MakeTenant(params)
		assert.Error(t, err, "Expected error for invalid field")
		assert.Contains(t, err.Error(), errorField, "Error should be related to the "+errorField+" field")
	})
}

func TestNewTenant_Success(t *testing.T) {
	params := NewTenantParams{
		Name:         "John Doe",
		Address:      "123 Main St",
		Email:        "johndoe@example.com",
		Phone:        "+1234567890",
		Timezone:     "America/New_York",
		OpeningHours: "09:00",
		ClosingHours: "18:00",
	}

	tenant, err := NewTenant(params)

	// Assertions
	assert.NoError(t, err, "Expected no error when creating a new tenant")
	assert.NotNil(t, tenant, "Expected a non-nil tenant")
	assert.Equal(t, params.Name, tenant.GetName(), "Tenant name mismatch")
	assert.Equal(t, params.Address, tenant.GetAddress(), "Tenant address mismatch")
	assert.Equal(t, params.Email, tenant.GetEmail(), "Tenant email mismatch")
	assert.Equal(t, params.Phone, tenant.GetPhone(), "Tenant phone mismatch")
	assert.Equal(t, params.Timezone, tenant.GetTimezone(), "Tenant timezone mismatch")
	assert.Equal(t, params.OpeningHours, tenant.GetOpeningHours(), "Tenant opening hours mismatch")
	assert.Equal(t, params.ClosingHours, tenant.GetClosingHours(), "Tenant closing hours mismatch")
}

func TestNewTenant_ValidationFailure(t *testing.T) {
	params := NewTenantParams{
		Name:         "", // Missing required field
		Address:      "123 Main St",
		Email:        "invalid-email", // Invalid email format
		Phone:        "1234567890",    // Invalid phone format
		Timezone:     "Invalid/Timezone",
		OpeningHours: "25:00", // Invalid time format
		ClosingHours: "26:00", // Invalid time format
	}

	tenant, err := NewTenant(params)

	// Assertions
	assert.Error(t, err, "Expected validation error when creating a new tenant with invalid params")
	assert.Nil(t, tenant, "Expected tenant to be nil on validation failure")
	assert.Contains(t, err.Error(), "failed NewTenant validate", "Error message should contain validation failure")
}
