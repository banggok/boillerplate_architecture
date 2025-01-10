package register

import (
	"net/http"
	"testing"

	asserthelper "github.com/banggok/boillerplate_architecture/tests/helper/assert_helper"
	"github.com/banggok/boillerplate_architecture/tests/helper/data"
	test_register "github.com/banggok/boillerplate_architecture/tests/helper/request/tenant/register"

	valueobject "github.com/banggok/boillerplate_architecture/internal/data/entity/value_object"
	"github.com/banggok/boillerplate_architecture/internal/data/model"
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/v1/tenant/register"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/services/tenant"
	"github.com/banggok/boillerplate_architecture/tests/helper/response"
	"github.com/banggok/boillerplate_architecture/tests/helper/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestTenantRegister(t *testing.T) {
	tenantData := data.TenantData
	accountData := data.AccountData
	reqBody := register.Request{
		Name:         tenantData.Name,
		Address:      tenantData.Address,
		Email:        tenantData.Email,
		Phone:        tenantData.Phone,
		Timezone:     tenantData.Timezone,
		OpeningHours: tenantData.OpeningHours,
		ClosingHours: tenantData.ClosingHours,
		Account: register.Account{
			Name:  accountData.Name,
			Email: accountData.Email,
			Phone: accountData.Phone,
		},
	}

	t.Run("Tenant Registration Success", func(t *testing.T) {
		server, cleanUp, mysqlCfg := setup.TestingEnv(t, true)
		defer cleanUp(mysqlCfg)

		resp := test_register.CallPostTenants(t, reqBody, server)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Read and parse the response body
		var responseBody struct {
			Status  string            `json:"status"`
			Message string            `json:"message"`
			Data    register.Response `json:"data"`
		}

		response.GetBodyResponse(t, resp, &responseBody)

		assert.Equal(t, "success", responseBody.Status)
		assert.Equal(t, "Tenant registered successfully", responseBody.Message)

		expectedResponse := register.Response{
			ID:        1,
			CreatedAt: responseBody.Data.CreatedAt,
			UpdatedAt: responseBody.Data.UpdatedAt,
		}

		assert.Equal(t, responseBody.Data, expectedResponse)
		var tenantDB model.Tenant

		require.NoError(t, mysqlCfg.Slave.Session(&gorm.Session{}).
			Preload("Accounts").Preload("Accounts.AccountVerifications").First(&tenantDB).Error)

		expectedTenantDB := model.Tenant{
			Metadata: model.Metadata{
				ID:        tenantDB.ID,
				CreatedAt: tenantDB.CreatedAt,
				UpdatedAt: tenantDB.UpdatedAt,
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
					Metadata: model.Metadata{
						ID:        (*tenantDB.Accounts)[0].ID,
						CreatedAt: (*tenantDB.Accounts)[0].CreatedAt,
						UpdatedAt: (*tenantDB.Accounts)[0].UpdatedAt,
					},
					Name:     reqBody.Account.Name,
					Email:    reqBody.Account.Email,
					Phone:    reqBody.Account.Phone,
					Password: (*tenantDB.Accounts)[0].Password,
					TenantID: tenantDB.ID,
					AccountVerifications: &[]model.AccountVerification{
						{
							Metadata: model.Metadata{
								ID:        (*(*tenantDB.Accounts)[0].AccountVerifications)[0].ID,
								CreatedAt: (*(*tenantDB.Accounts)[0].AccountVerifications)[0].CreatedAt,
								UpdatedAt: (*(*tenantDB.Accounts)[0].AccountVerifications)[0].UpdatedAt,
							},
							AccountID: (*tenantDB.Accounts)[0].ID,
							Type:      string(valueobject.EMAIL_VERIFICATION.String()),
							Token:     (*(*tenantDB.Accounts)[0].AccountVerifications)[0].Token,
							ExpiresAt: (*(*tenantDB.Accounts)[0].AccountVerifications)[0].ExpiresAt,
							Verified:  false,
						},
						{
							Metadata: model.Metadata{
								ID:        (*(*tenantDB.Accounts)[0].AccountVerifications)[1].ID,
								CreatedAt: (*(*tenantDB.Accounts)[0].AccountVerifications)[1].CreatedAt,
								UpdatedAt: (*(*tenantDB.Accounts)[0].AccountVerifications)[1].UpdatedAt,
							},
							AccountID: (*tenantDB.Accounts)[0].ID,
							Type:      string(valueobject.RESET_PASSWORD.String()),
							ExpiresAt: (*(*tenantDB.Accounts)[0].AccountVerifications)[1].ExpiresAt,
							Verified:  false,
						},
					},
				},
			},
		}

		asserthelper.AssertStruct(t, tenantDB, expectedTenantDB)

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
		var errorResponse response.ErrorResponse
		response.GetBodyResponse(t, resp, &errorResponse)

		expectedResponse := response.ErrorResponse{
			RequestID: errorResponse.RequestID,
			Status:    "error",
			Code:      int(custom_errors.TenantConflictEntity),
			Message:   "tenant exists",
			Details: map[string]string{
				"error": tenant.TenantDuplicate,
			},
		}

		// Assert response body structure and values
		assert.Equal(t, errorResponse, expectedResponse)
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
		var errorResponse response.ErrorResponse
		response.GetBodyResponse(t, resp, &errorResponse)

		expectedResponse := response.ErrorResponse{
			RequestID: errorResponse.RequestID,
			Status:    "error",
			Code:      int(custom_errors.TenantBadRequest),
			Message:   "invalid tenant request",
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
		var errorResponse response.ErrorResponse
		response.GetBodyResponse(t, resp, &errorResponse)

		expectedResponse := response.ErrorResponse{
			RequestID: errorResponse.RequestID,
			Status:    "error",
			Code:      int(custom_errors.TenantUnprocessEntity),
			Message:   "can not process tenant",
			Details: map[string]string{
				"name":         "Name is required.",
				"account.name": "Account name is required and can only contain letters and spaces.",
			},
		}

		// Assert response body structure and values
		assert.Equal(t, expectedResponse, errorResponse)
	})
}
