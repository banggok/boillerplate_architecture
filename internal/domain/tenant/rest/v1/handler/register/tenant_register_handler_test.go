package register

import (
	"appointment_management_system/internal/config"
	"appointment_management_system/internal/domain/tenant/entity"
	"appointment_management_system/internal/domain/tenant/rest/v1/handler/register/dto"
	"appointment_management_system/internal/pkg/middleware"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for TenantRegisterUsecase
type MockTenantRegisterUsecase struct {
	mock.Mock
}

func (m *MockTenantRegisterUsecase) Execute(ctx *gin.Context, request dto.RegisterTenantRequest) (*entity.Tenant, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*entity.Tenant), args.Error(1)
}

// Setup a fresh test router for each test case
func setupTestRouter(handler TenantRegisterHandler) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.ErrorHandlingMiddleware())
	router.POST("/v1/tenants/register", handler.Register)
	return router
}

func TestTenantHandler_Register(t *testing.T) {
	validate := config.SetupValidator()
	mockUsecase := new(MockTenantRegisterUsecase)

	t.Run("Success Case", func(t *testing.T) {
		testSuccessCase(t, validate, mockUsecase)
	})

	t.Run("Validation Failure", func(t *testing.T) {
		testValidationFailure(t, validate, mockUsecase)
	})

	t.Run("Invalid JSON Payload", func(t *testing.T) {
		testInvalidJSONPayload(t, validate, mockUsecase)
	})

	t.Run("Usecase Error", func(t *testing.T) {
		testUsecaseError(t, validate)
	})
}

func testSuccessCase(t *testing.T, validate *validator.Validate, mockUsecase *MockTenantRegisterUsecase) {
	handler := createTestHandler(validate, mockUsecase)
	router := setupTestRouter(handler)
	payload := getValidTenantRequest()

	tenant := createTestTenant(payload)

	jsonPayload, _ := json.Marshal(payload)
	mockUsecase.On("Execute", mock.Anything, payload).Return(tenant, nil)

	req := createTestRequest(http.MethodPost, "/v1/tenants/register", jsonPayload)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	validateSuccessResponse(t, recorder, tenant)
	mockUsecase.AssertCalled(t, "Execute", mock.Anything, payload)
}

func createTestHandler(validate *validator.Validate, mockUsecase *MockTenantRegisterUsecase) *tenantRegisterHandler {
	return &tenantRegisterHandler{
		validator:       validate,
		registerUsecase: mockUsecase,
	}
}

func getValidTenantRequest() dto.RegisterTenantRequest {
	return dto.RegisterTenantRequest{
		Name:         "Acme Corporation",
		Address:      "123 Main St",
		Email:        "contact@acme.com",
		Phone:        "+1234567890",
		Timezone:     "America/New_York",
		OpeningHours: "09:00",
		ClosingHours: "18:00",
		Account: dto.Account{
			Name:  "John Doe",
			Email: "admin@acme.com",
			Phone: "+0987654321",
		},
	}
}

func createTestTenant(request dto.RegisterTenantRequest) *entity.Tenant {
	tenant, _ := entity.MakeTenant(entity.MakeTenantParams{
		ID:           1,
		Name:         request.Name,
		Address:      request.Address,
		Email:        request.Email,
		Phone:        request.Phone,
		Timezone:     request.Timezone,
		OpeningHours: request.OpeningHours,
		ClosingHours: request.ClosingHours,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	return tenant
}

func createTestRequest(method, url string, body []byte) *http.Request {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func validateSuccessResponse(t *testing.T, recorder *httptest.ResponseRecorder, tenant *entity.Tenant) {
	var responseBody map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, "Tenant registered successfully", responseBody["message"])

	data, ok := responseBody["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, tenant.GetName(), data["name"])
	assert.Equal(t, tenant.GetAddress(), data["address"])
	assert.Equal(t, tenant.GetEmail(), data["email"])
	assert.Equal(t, tenant.GetPhone(), data["phone"])
	assert.Equal(t, tenant.GetTimezone(), data["timezone"])
	assert.Equal(t, tenant.GetOpeningHours(), data["opening_hours"])
	assert.Equal(t, tenant.GetClosingHours(), data["closing_hours"])
}

func testValidationFailure(t *testing.T, validate *validator.Validate, mockUsecase *MockTenantRegisterUsecase) {
	handler := &tenantRegisterHandler{
		validator:       validate,
		registerUsecase: mockUsecase,
	}
	router := setupTestRouter(handler)

	invalidPayload := dto.RegisterTenantRequest{}
	jsonPayload, _ := json.Marshal(invalidPayload)

	req, _ := http.NewRequest(http.MethodPost, "/v1/tenants/register", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "error", responseBody["status"])
	assert.Equal(t, "VALIDATION_FAILED", responseBody["code"])

	details := responseBody["details"].(map[string]interface{})
	assert.Equal(t, "Name is required.", details["name"])
}

func testInvalidJSONPayload(t *testing.T, validate *validator.Validate, mockUsecase *MockTenantRegisterUsecase) {
	handler := &tenantRegisterHandler{
		validator:       validate,
		registerUsecase: mockUsecase,
	}
	router := setupTestRouter(handler)

	invalidJSON := `{"name": "Acme Corporation", "address": "123 Main St", "email": "contact@acme.com",`
	req, _ := http.NewRequest(http.MethodPost, "/v1/tenants/register", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "error", responseBody["status"])
	assert.Equal(t, "INVALID_PAYLOAD", responseBody["code"])
}

func testUsecaseError(t *testing.T, validate *validator.Validate) {
	mockUsecase := new(MockTenantRegisterUsecase)
	handler := &tenantRegisterHandler{
		validator:       validate,
		registerUsecase: mockUsecase,
	}
	router := setupTestRouter(handler)

	payload := dto.RegisterTenantRequest{
		Name:         "Acme Corporation",
		Address:      "123 Main St",
		Email:        "contact@acme.com",
		Phone:        "+1234567890",
		Timezone:     "America/New_York",
		OpeningHours: "09:00",
		ClosingHours: "18:00",
		Account: dto.Account{
			Name:  "John Doe",
			Email: "admin@acme.com",
			Phone: "+0987654321",
		},
	}
	jsonPayload, _ := json.Marshal(payload)

	mockUsecase.On("Execute", mock.Anything, payload).Return((*entity.Tenant)(nil), fmt.Errorf("Database unavailable"))

	req, _ := http.NewRequest(http.MethodPost, "/v1/tenants/register", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "error", responseBody["status"])
	assert.Equal(t, "TENANT_REGISTRATION_FAILED", responseBody["code"])
}

func TestNewTenantRegisterHandler(t *testing.T) {
	v := validator.New()
	mockUsecase := new(MockTenantRegisterUsecase)
	handler := NewTenantRegisterHandler(v, mockUsecase)

	assert.NotNil(t, handler, "Handler should be created successfully")
}
