package register

import (
	"appointment_management_system/internal/data/model"
	"appointment_management_system/internal/delivery/rest/v1/tenant/register"
	"appointment_management_system/internal/pkg/custom_errors"
	"appointment_management_system/internal/services/tenant"
	test_register "appointment_management_system/tests/e2e/helper/request/tenant/register"
	"appointment_management_system/tests/e2e/helper/response"
	"appointment_management_system/tests/e2e/helper/setup"
	time_testing_helper "appointment_management_system/tests/e2e/helper/time"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestTenantRegister(t *testing.T) {
	reqBody := register.RegisterTenantRequest{
		Name:         "Doe's Bakery",
		Address:      "123 Main Street, Springfield, USA",
		Email:        "business@example.com",
		Phone:        "+0987654321",
		Timezone:     "America/New_York",
		OpeningHours: "08:00",
		ClosingHours: "20:00",
		Account: register.RegisterTenantAccountRequest{
			Name:  "John Doe",
			Email: "rtriasmono@gmail.com",
			Phone: "+1234567890",
		},
	}

	t.Run("Tenant Registration Success", func(t *testing.T) {
		server, cleanUp, mysqlCfg := setup.TestingEnv(t, true)
		defer cleanUp(mysqlCfg)

		resp := test_register.CallPostTenants(t, reqBody, server)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Read and parse the response body
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var responseBody struct {
			Status  string                          `json:"status"`
			Message string                          `json:"message"`
			Data    register.RegisterTenantResponse `json:"data"`
		}

		err = json.Unmarshal(body, &responseBody)
		require.NoError(t, err)

		assert.Equal(t, "success", responseBody.Status)
		assert.Equal(t, "Tenant registered successfully", responseBody.Message)

		expectedResponse := register.RegisterTenantResponse{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Assert response body structure and values
		err = time_testing_helper.Sanitize(&responseBody.Data.CreatedAt,
			expectedResponse.CreatedAt)
		require.NoError(t, err)

		err = time_testing_helper.Sanitize(&responseBody.Data.UpdatedAt,
			expectedResponse.UpdatedAt)
		require.NoError(t, err)

		assert.Equal(t, responseBody.Data, expectedResponse)
		var tenantDB model.Tenant

		require.NoError(t, mysqlCfg.Slave.Session(&gorm.Session{}).Preload("Accounts").First(&tenantDB).Error)
		expectedTenantDB := model.Tenant{
			Metadata: model.Metadata{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:         reqBody.Name,
			Address:      reqBody.Address,
			Email:        reqBody.Email,
			Phone:        reqBody.Phone,
			Timezone:     reqBody.Timezone,
			OpeningHours: reqBody.OpeningHours,
			ClosingHours: reqBody.ClosingHours,
			Accounts: &[]model.Account{
				{
					ID:        1,
					Name:      reqBody.Account.Name,
					Email:     reqBody.Account.Email,
					Phone:     reqBody.Account.Phone,
					Password:  (*tenantDB.Accounts)[0].Password,
					TenantID:  1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		}

		err = time_testing_helper.Sanitize(&tenantDB.CreatedAt,
			expectedTenantDB.CreatedAt)
		require.NoError(t, err)

		err = time_testing_helper.Sanitize(&tenantDB.UpdatedAt,
			expectedTenantDB.UpdatedAt)
		require.NoError(t, err)

		for i := range *expectedTenantDB.Accounts {
			accountDb := (*tenantDB.Accounts)[i]
			err = time_testing_helper.Sanitize(&accountDb.CreatedAt,
				(*expectedTenantDB.Accounts)[i].CreatedAt)
			require.NoError(t, err)

			err = time_testing_helper.Sanitize(&accountDb.UpdatedAt,
				(*expectedTenantDB.Accounts)[i].UpdatedAt)
			require.NoError(t, err)
			(*tenantDB.Accounts)[i] = accountDb
		}

		tenantDBJSON, err := json.Marshal(tenantDB)
		require.NoError(t, err)

		expectedTenantDBJSON, err := json.Marshal(expectedTenantDB)
		require.NoError(t, err)

		assert.JSONEq(t, string(expectedTenantDBJSON), string(tenantDBJSON), "Tenant data mismatch")

	})

	t.Run("tenant registration should be failed due to duplication", func(t *testing.T) {
		server, cleanUp, mysqlCfg := setup.TestingEnv(t, true)
		defer cleanUp(mysqlCfg)

		resp := test_register.CallPostTenants(t, reqBody, server)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		resp = test_register.CallPostTenants(t, reqBody, server)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		// Read and parse the response body
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var errorResonse response.ErrorResponse

		err = json.Unmarshal(body, &errorResonse)
		require.NoError(t, err)

		expectedResponse := response.ErrorResponse{
			RequestID: errorResonse.RequestID,
			Status:    "error",
			Code:      int(custom_errors.TenantConflictEntity),
			Message:   tenant.TenantDuplicate,
			Details: map[string]string{
				"error": tenant.TenantDuplicate,
			},
		}

		// Assert response body structure and values
		assert.Equal(t, errorResonse, expectedResponse)
	})

	t.Run("Tenant Registration Bad Request (Invalid JSON)", func(t *testing.T) {
		server, cleanUp, mysqlCfg := setup.TestingEnv(t, true)
		defer cleanUp(mysqlCfg)

		// Invalid JSON payload
		invalidPayload := `{"name": "Doe's Bakery", "email": "invalid-email",}` // Malformed JSON

		resp := test_register.CallPostTenantsWithRawPayload(t, invalidPayload, server)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Read and parse the response body
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var errorResponse response.ErrorResponse

		err = json.Unmarshal(body, &errorResponse)
		require.NoError(t, err)

		expectedResponse := response.ErrorResponse{
			RequestID: errorResponse.RequestID,
			Status:    "error",
			Code:      int(custom_errors.TenantBadRequest),
			Message:   "invalid tenant register request payload",
			Details: map[string]string{
				"error": "invalid character '}' looking for beginning of object key string",
			},
		}

		// Assert response body structure and values
		assert.Equal(t, errorResponse, expectedResponse)
	})

	// New Test Case: Validation Error (422 Unprocessable Entity)
	t.Run("Tenant Registration Unprocessable Entity (Validation Error)", func(t *testing.T) {
		server, cleanUp, mysqlCfg := setup.TestingEnv(t, true)
		defer cleanUp(mysqlCfg)

		reqBody.Name = ""
		reqBody.Account.Name = ""

		resp := test_register.CallPostTenants(t, reqBody, server)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		// Read and parse the response body
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var errorResponse response.ErrorResponse

		err = json.Unmarshal(body, &errorResponse)
		require.NoError(t, err)

		expectedResponse := response.ErrorResponse{
			RequestID: errorResponse.RequestID,
			Status:    "error",
			Code:      int(custom_errors.TenantUnprocessEntity),
			Message:   "failed to validate request payload",
			Details: map[string]string{
				"name":         "Name is required.",
				"account.name": "Account name is required and can only contain letters and spaces.",
			},
		}

		// Assert response body structure and values
		assert.Equal(t, errorResponse, expectedResponse)
	})
}
