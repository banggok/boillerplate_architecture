package register

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"appointment_management_system/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewTenantV1RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set up mock repositories and dependencies
	mockRepository := new(MockTenantCreateRepository)
	mockValidator := config.SetupValidator()

	router := gin.Default()
	tenantGroup := router.Group("/tenant/v1")

	// Call the function under test
	NewTenantV1RegisterRoutes(tenantGroup, mockValidator, mockRepository)

	t.Run("RegisterRoute", func(t *testing.T) {
		// Define the payload
		reqBody := `{
			"name": "Acme Corporation",
			"address": "123 Main St",
			"email": "contact@acme.com",
			"phone": "+1234567890",
			"timezone": "America/New_York",
			"openingHours": "09:00",
			"closingHours": "18:00",
			"account": {
				"name": "John Doe",
				"email": "admin@acme.com",
				"phone": "+0987654321"
			}
		}`

		// Mock the repository call
		mockRepository.On("Execute", mock.Anything, mock.Anything).Return(nil)

		// Create a test request
		req, _ := http.NewRequest(http.MethodPost, "/tenant/v1/register", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert the response
		assert.Equal(t, http.StatusCreated, w.Code)
		mockRepository.AssertCalled(t, "Execute", mock.Anything, mock.Anything)
	})
}
