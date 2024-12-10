package repository

import (
	"appointment_management_system/internal/domain/tenants/entity"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Mock for TenantIsPhoneExistsRepository
type mockTenantIsPhoneExistsRepository struct {
	mock.Mock
}

func (m *mockTenantIsPhoneExistsRepository) Execute(ctx *gin.Context, phoneNumber string) (*bool, error) {
	args := m.Called(ctx, phoneNumber)
	return args.Get(0).(*bool), args.Error(1)
}

// Mock for TenantIsEmailExistsRepository
type mockTenantIsEmailExistsRepository struct {
	mock.Mock
}

func (m *mockTenantIsEmailExistsRepository) Execute(ctx *gin.Context, email string) (*bool, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*bool), args.Error(1)
}

// Unit Test for Execute
func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	return ctx, w
}

func setupMockDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Create schema
	if err := db.AutoMigrate(&tenantModel{}); err != nil {
		return nil, err
	}

	return db, nil
}

func mockRepositories(mockPhoneRepo *mockTenantIsPhoneExistsRepository, mockEmailRepo *mockTenantIsEmailExistsRepository, phoneExists, emailExists bool, ctx *gin.Context, tenant *entity.Tenant) {
	mockPhoneRepo.On("Execute", ctx, tenant.GetPhone()).Return(&phoneExists, nil).Once()
	mockEmailRepo.On("Execute", ctx, tenant.GetEmail()).Return(&emailExists, nil).Once()
}

type testCase struct {
	name          string
	phoneExists   bool
	emailExists   bool
	mockDBError   error
	expectedError string
}

func runTestCase(t *testing.T, repo TenantCreateRepository, ctx *gin.Context, tenant *entity.Tenant, mockPhoneRepo *mockTenantIsPhoneExistsRepository, mockEmailRepo *mockTenantIsEmailExistsRepository, tc testCase) {
	helperGetDBOri := helperGetDB
	defer func() {
		helperGetDB = helperGetDBOri

	}()

	// Setup mocks
	mockRepositories(mockPhoneRepo, mockEmailRepo, tc.phoneExists, tc.emailExists, ctx, tenant)

	// Override DB helper for specific test cases
	helperGetDB = func(_ *gin.Context) (*gorm.DB, error) {
		if tc.mockDBError != nil {
			return nil, tc.mockDBError
		}
		return setupMockDB()
	}

	// Run the test
	err := repo.Execute(ctx, tenant)

	// Validate results
	if tc.expectedError != "" {
		assert.Error(t, err)
		assert.EqualError(t, err, tc.expectedError)
	} else {
		assert.NoError(t, err)
	}

	// Validate mock expectations
	mockPhoneRepo.AssertExpectations(t)
	mockEmailRepo.AssertExpectations(t)
}

func defineTestCase() []testCase {
	testCases := []testCase{
		{
			name:        "successfully create tenant",
			phoneExists: false,
			emailExists: false,
		},
		{
			name:          "phone already exists",
			phoneExists:   true,
			emailExists:   false,
			expectedError: "phone number already exists",
		},
		{
			name:          "email already exists",
			phoneExists:   false,
			emailExists:   true,
			expectedError: "email already exists",
		},
		{
			name:          "database transaction error",
			phoneExists:   false,
			emailExists:   false,
			mockDBError:   errors.New("database connection error"),
			expectedError: "failed to get transaction: database connection error",
		},
	}
	return testCases
}

func makeTenant() *entity.Tenant {
	tenant, _ := entity.MakeTenant(entity.MakeTenantParams{
		ID:           1,
		Name:         "Test Tenant",
		Address:      "123 Test Street",
		Email:        "test@tenant.com",
		Phone:        "+1234567890",
		Timezone:     "UTC",
		OpeningHours: "09:00",
		ClosingHours: "18:00",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	return tenant
}

func TestTenantCreateRepository_Execute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPhoneRepo := new(mockTenantIsPhoneExistsRepository)
	mockEmailRepo := new(mockTenantIsEmailExistsRepository)

	repo := NewTenantCreateRepository()

	repo.(*tenantCreateRepository).isEmailExists = mockEmailRepo
	repo.(*tenantCreateRepository).isPhoneExists = mockPhoneRepo

	tenant := makeTenant()

	ctx, _ := setupTestContext()

	// Define test cases
	testCases := defineTestCase()

	// Execute test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTestCase(t, repo, ctx, tenant, mockPhoneRepo, mockEmailRepo, tc)
		})
	}
}

func TestCheckPhoneExistence_ErrorCases(t *testing.T) {
	mockPhoneRepo := new(mockTenantIsPhoneExistsRepository)

	repo := &tenantCreateRepository{
		isPhoneExists: mockPhoneRepo,
	}

	tenant := &entity.Tenant{}
	var phoneExists bool

	t.Run("phone existence check fails", func(t *testing.T) {
		ginCtx, _ := gin.CreateTestContext(nil)

		mockPhoneRepo.On("Execute", ginCtx, "").Return((*bool)(nil), fmt.Errorf("phone check error"))

		err := repo.checkPhoneExistence(ginCtx, tenant, &phoneExists)

		assert.Error(t, err)
		assert.EqualError(t, err, "phone check error")
		mockPhoneRepo.AssertExpectations(t)
	})
}

func TestCheckEmailExistence_ErrorCases(t *testing.T) {
	mockEmailRepo := new(mockTenantIsEmailExistsRepository)

	repo := &tenantCreateRepository{
		isEmailExists: mockEmailRepo,
	}

	tenant := &entity.Tenant{}
	var emailExists bool

	t.Run("email existence check fails", func(t *testing.T) {
		ginCtx, _ := gin.CreateTestContext(nil)

		// Correctly return (*bool)(nil) instead of nil
		mockEmailRepo.On("Execute", ginCtx, "").Return((*bool)(nil), fmt.Errorf("email check error"))

		err := repo.checkEmailExistence(ginCtx, tenant, &emailExists)

		assert.Error(t, err)
		assert.EqualError(t, err, "email check error")
		mockEmailRepo.AssertExpectations(t)
	})
}

// Unit test for validateTenantExistence with errors in checkPhoneExistence and checkEmailExistence
func TestValidateTenantExistence_ErrorCases(t *testing.T) {
	mockPhoneRepo := new(mockTenantIsPhoneExistsRepository)
	mockEmailRepo := new(mockTenantIsEmailExistsRepository)

	repo := &tenantCreateRepository{
		isPhoneExists: mockPhoneRepo,
		isEmailExists: mockEmailRepo,
	}

	tenant := &entity.Tenant{}
	ginCtx, _ := gin.CreateTestContext(nil)

	runTest := func(t *testing.T, name string, mockSetup func(), expectedError string) {
		t.Run(name, func(t *testing.T) {
			if mockSetup != nil {
				mockSetup()
			}

			phoneExists, emailExists, err := repo.validateTenantExistence(ginCtx, tenant)

			assert.Error(t, err)
			assert.EqualError(t, err, expectedError)
			assert.False(t, phoneExists)
			assert.False(t, emailExists)

			mockPhoneRepo.AssertExpectations(t)
			mockEmailRepo.AssertExpectations(t)
		})
	}

	runTest(t, "checkPhoneExistence returns error", func() {
		mockPhoneRepo.On("Execute", ginCtx, "").Return((*bool)(nil), errors.New("phone check error")).Once()
	}, "phone check error")

	runTest(t, "checkEmailExistence returns error", func() {
		mockPhoneRepo.On("Execute", ginCtx, "").Return((*bool)(nil), nil).Once()
		mockEmailRepo.On("Execute", ginCtx, "").Return((*bool)(nil), errors.New("email check error")).Once()
	}, "email check error")
}
