package login

import (
	"net/http"
	"testing"

	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	valueobject "github.com/banggok/boillerplate_architecture/internal/data/entity/value_object"
	"github.com/banggok/boillerplate_architecture/internal/data/model"
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/v1/account/login"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/middleware/authentication"
	"github.com/banggok/boillerplate_architecture/tests/helper/data"
	test_login "github.com/banggok/boillerplate_architecture/tests/helper/request/account/login"
	"github.com/banggok/boillerplate_architecture/tests/helper/response"
	"github.com/banggok/boillerplate_architecture/tests/helper/setup"
	"github.com/stretchr/testify/assert"
)

func TestAccountLogin(t *testing.T) {
	t.Run("Login Bad Request (Invalid JSON)", func(t *testing.T) {
		server, cleanUp, mysqlCfg := setup.TestingEnv(t, true)
		defer cleanUp(mysqlCfg)

		// Invalid JSON payload
		invalidPayload := `{"email": "invalid-email",}` // Malformed JSON

		resp := test_login.CallRequestWithRawPayload(t, invalidPayload, server)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Read and parse the response body
		var errorResponse response.ErrorResponse
		response.GetBodyResponse(t, resp, &errorResponse)

		expectedResponse := response.ErrorResponse{
			RequestID: errorResponse.RequestID,
			Status:    "error",
			Code:      int(custom_errors.AccountBadRequest),
			Message:   "invalid account request",
			Details: map[string]string{
				"error": "invalid character '}' looking for beginning of object key string",
			},
		}

		// Assert response body structure and values
		assert.Equal(t, errorResponse, expectedResponse)
	})

	t.Run("login Unprocessable Entity (Validation Error)", func(t *testing.T) {
		server, cleanUp, mysqlCfg := setup.TestingEnv(t, true)
		defer cleanUp(mysqlCfg)

		requestBody := login.Request{}

		resp := test_login.CallRequest(t, requestBody, server)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		// Read and parse the response body
		var errorResponse response.ErrorResponse
		response.GetBodyResponse(t, resp, &errorResponse)

		expectedResponse := response.ErrorResponse{
			RequestID: errorResponse.RequestID,
			Status:    "error",
			Code:      int(custom_errors.AccountUnprocessEntity),
			Message:   "can not process account",
			Details: map[string]string{
				"email":    "Email is required and must be in a valid format.",
				"password": "Password is required",
			},
		}

		// Assert response body structure and values
		assert.Equal(t, expectedResponse, errorResponse)
	})

	t.Run("invalid password", func(t *testing.T) {
		server, cleanUp, mysqlCfg := setup.TestingEnv(t, true)
		defer cleanUp(mysqlCfg)

		accountData := data.AccountData
		tenantData := data.TenantData
		tenantData.Accounts = &[]model.Account{accountData}

		data.PrepareData(t, mysqlCfg, tenantData)

		requestBody := login.Request{
			Email:    accountData.Email,
			Password: "invalid",
		}

		resp := test_login.CallRequest(t, requestBody, server)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Read and parse the response body
		var errorResponse response.ErrorResponse
		response.GetBodyResponse(t, resp, &errorResponse)

		expectedResponse := response.ErrorResponse{
			RequestID: errorResponse.RequestID,
			Status:    "error",
			Code:      int(custom_errors.Unauthorized),
			Message:   "Unauthorized",
			Details: map[string]string{
				"error": "wrong password",
			},
		}

		// Assert response body structure and values
		assert.Equal(t, expectedResponse, errorResponse)
	})

	t.Run("success with email unverified", func(t *testing.T) {
		verifiedAction(t, string(entity.EMAIL))
	})

	t.Run("success with change password unverified", func(t *testing.T) {
		verifiedAction(t, string(entity.CHANGE_PASSWORD))
	})

	t.Run("success verified", func(t *testing.T) {
		verifiedAction(t, string(entity.VERIFIED))
	})
}

func verifiedAction(t *testing.T, action string) {
	server, cleanUp, mysqlCfg := setup.TestingEnv(t, true)
	defer cleanUp(mysqlCfg)

	emailVerificationData := data.AccountVerificationData
	if action == string(entity.CHANGE_PASSWORD) || action == string(entity.VERIFIED) {
		emailVerificationData.Verified = true
	}
	changePasswordVerificationData := data.AccountVerificationData
	changePasswordVerificationData.Type = valueobject.CHANGE_PASSWORD.String()
	changePasswordVerificationData.Token = nil
	if action == string(entity.VERIFIED) {
		changePasswordVerificationData.Verified = true
	}
	accountData := data.AccountData
	accountData.AccountVerifications = &[]model.AccountVerification{emailVerificationData, changePasswordVerificationData}
	tenantData := data.TenantData
	tenantData.Accounts = &[]model.Account{accountData}

	data.PrepareData(t, mysqlCfg, &tenantData)

	accountData = (*tenantData.Accounts)[0]

	requestBody := login.Request{
		Email:    accountData.Email,
		Password: data.PlainPassword,
	}

	resp := test_login.CallRequest(t, requestBody, server)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read and parse the response body
	var responseBody struct {
		Status  string         `json:"status"`
		Message string         `json:"message"`
		Data    login.Response `json:"data"`
	}
	response.GetBodyResponse(t, resp, &responseBody)

	assert.Equal(t, "success", responseBody.Status)
	assert.Equal(t, "Account logged in successfully", responseBody.Message)

	assert.Equal(t, accountData.ID, responseBody.Data.ID)
	assert.Equal(t, action, responseBody.Data.Action)

	accessToken, err := authentication.ValidateToken(responseBody.Data.AccessToken, true)
	assert.NoError(t, err)
	assert.Equal(t, accountData.ID, accessToken.UserID)

	refreshToken, err := authentication.ValidateToken(responseBody.Data.RefreshToken, false)
	assert.NoError(t, err)
	assert.Equal(t, accountData.ID, refreshToken.UserID)

	assert.Equal(t, accountData.CreatedAt.UTC(), responseBody.Data.CreatedAt.UTC())
	assert.Equal(t, accountData.UpdatedAt.UTC(), responseBody.Data.UpdatedAt.UTC())
}
