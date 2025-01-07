package register

import (
	"net/http"
	"testing"
	"time"

	asserthelper "github.com/banggok/boillerplate_architecture/tests/helper/assert_helper"
	test_register "github.com/banggok/boillerplate_architecture/tests/helper/request/tenant/register"
	time_testing_helper "github.com/banggok/boillerplate_architecture/tests/helper/time"

	"github.com/banggok/boillerplate_architecture/internal/config/app"
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
	reqBody := register.Request{
		Name:         "Doe's Bakery",
		Address:      "123 Main Street, Springfield, USA",
		Email:        "business@example.com",
		Phone:        "+0987654321",
		Timezone:     "America/New_York",
		OpeningHours: "08:00",
		ClosingHours: "20:00",
		Account: register.Account{
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
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Assert response body structure and values
		time_testing_helper.Sanitize(t, &responseBody.Data.CreatedAt,
			expectedResponse.CreatedAt)

		time_testing_helper.Sanitize(t, &responseBody.Data.UpdatedAt,
			expectedResponse.UpdatedAt)

		assert.Equal(t, responseBody.Data, expectedResponse)
		var tenantDB model.Tenant

		require.NoError(t, mysqlCfg.Slave.Session(&gorm.Session{}).
			Preload("Accounts").Preload("Accounts.AccountVerifications").First(&tenantDB).Error)

		expiredDuration := app.AppConfig.ExpiredDuration.EmailVerification
		expiresAt := time.Now().UTC().Add(expiredDuration)

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
					Metadata: model.Metadata{
						ID:        1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Name:     reqBody.Account.Name,
					Email:    reqBody.Account.Email,
					Phone:    reqBody.Account.Phone,
					Password: (*tenantDB.Accounts)[0].Password,
					TenantID: 1,
					AccountVerifications: &[]model.AccountVerification{
						{
							Metadata: model.Metadata{
								ID:        1,
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
							},
							AccountID: 1,
							Type:      string(valueobject.EMAIL_VERIFICATION.String()),
							Token:     (*(*tenantDB.Accounts)[0].AccountVerifications)[0].Token,
							ExpiresAt: expiresAt,
							Verified:  false,
						},
					},
				},
			},
		}

		time_testing_helper.Sanitize(t, &tenantDB.CreatedAt,
			expectedTenantDB.CreatedAt)

		time_testing_helper.Sanitize(t, &tenantDB.UpdatedAt,
			expectedTenantDB.UpdatedAt)

		for i := range *expectedTenantDB.Accounts {
			accountDb := (*tenantDB.Accounts)[i]

			time_testing_helper.Sanitize(t, &accountDb.CreatedAt,
				(*expectedTenantDB.Accounts)[i].CreatedAt)

			time_testing_helper.Sanitize(t, &accountDb.UpdatedAt,
				(*expectedTenantDB.Accounts)[i].UpdatedAt)

			for j := range *accountDb.AccountVerifications {
				accountVerification := (*accountDb.AccountVerifications)[j]

				time_testing_helper.Sanitize(t, &accountVerification.ExpiresAt,
					(*(*expectedTenantDB.Accounts)[i].AccountVerifications)[j].ExpiresAt)

				time_testing_helper.Sanitize(t, &accountVerification.CreatedAt,
					(*(*expectedTenantDB.Accounts)[i].AccountVerifications)[j].CreatedAt)

				time_testing_helper.Sanitize(t, &accountVerification.UpdatedAt,
					(*(*expectedTenantDB.Accounts)[i].AccountVerifications)[j].UpdatedAt)

				(*accountDb.AccountVerifications)[j] = accountVerification
			}
			(*tenantDB.Accounts)[i] = accountDb
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
