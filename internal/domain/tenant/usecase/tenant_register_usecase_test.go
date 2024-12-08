package usecase

import (
	"appointment_management_system/internal/domain/tenant/entity"
	"appointment_management_system/internal/domain/tenant/rest/v1/handler/register/dto"
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTenantCreateService is a mock implementation of the TenantCreateService interface
type MockTenantCreateService struct {
	mock.Mock
}

func (m *MockTenantCreateService) Execute(ctx *gin.Context, tenant *entity.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func TestTenantRegisterUsecase_Success(t *testing.T) {
	// Mock service
	mockService := new(MockTenantCreateService)
	mockService.On("Execute", mock.Anything, mock.Anything).Return(nil)

	// Create the usecase with the mock service
	usecase := NewTenantRegisterUsecase(mockService)

	// Prepare test data
	request := dto.RegisterTenantRequest{
		Name:         "Test Tenant",
		Address:      "123 Main St",
		Email:        "test@tenant.com",
		Phone:        "+1234567890",
		Timezone:     "America/New_York",
		OpeningHours: "09:00",
		ClosingHours: "17:00",
	}

	// Call the usecase
	ctx := &gin.Context{}
	tenant, err := usecase.Execute(ctx, request)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, tenant)
	assert.Equal(t, request.Name, tenant.GetName())
	assert.Equal(t, request.Address, tenant.GetAddress())
	assert.Equal(t, request.Email, tenant.GetEmail())
	assert.Equal(t, request.Phone, tenant.GetPhone())
	assert.Equal(t, request.Timezone, tenant.GetTimezone())
	assert.Equal(t, request.OpeningHours, tenant.GetOpeningHours())
	assert.Equal(t, request.ClosingHours, tenant.GetClosingHours())

	// Verify mock
	mockService.AssertCalled(t, "Execute", mock.Anything, mock.Anything)
}

func TestTenantRegisterUsecase_Failure_NewTenant(t *testing.T) {
	// Mock service (not called in this case)
	mockService := new(MockTenantCreateService)

	// Create the usecase with the mock service
	usecase := NewTenantRegisterUsecase(mockService)

	// Prepare invalid test data
	request := dto.RegisterTenantRequest{
		Name:         "", // Missing required field
		Address:      "123 Main St",
		Email:        "test@tenant.com",
		Phone:        "+1234567890",
		Timezone:     "America/New_York",
		OpeningHours: "09:00",
		ClosingHours: "17:00",
	}

	// Call the usecase
	ctx := &gin.Context{}
	tenant, err := usecase.Execute(ctx, request)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, tenant)
	assert.Contains(t, err.Error(), "failed TenantRegisterUsecase.Execute entity.NewTenant")

	// Verify that the service was never called
	mockService.AssertNotCalled(t, "Execute", mock.Anything, mock.Anything)
}

func TestTenantRegisterUsecase_Failure_CreateService(t *testing.T) {
	// Mock service
	mockService := new(MockTenantCreateService)
	mockService.On("Execute", mock.Anything, mock.Anything).Return(errors.New("service error"))

	// Create the usecase with the mock service
	usecase := NewTenantRegisterUsecase(mockService)

	// Prepare test data
	request := dto.RegisterTenantRequest{
		Name:         "Test Tenant",
		Address:      "123 Main St",
		Email:        "test@tenant.com",
		Phone:        "+1234567890",
		Timezone:     "America/New_York",
		OpeningHours: "09:00",
		ClosingHours: "17:00",
	}

	// Call the usecase
	ctx := &gin.Context{}
	tenant, err := usecase.Execute(ctx, request)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, tenant)
	assert.Contains(t, err.Error(), "failed TenantRegisterUsecase.Execute createTenantService.Execute")

	// Verify mock
	mockService.AssertCalled(t, "Execute", mock.Anything, mock.Anything)
}
