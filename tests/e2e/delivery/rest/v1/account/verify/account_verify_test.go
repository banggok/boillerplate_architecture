package verify_test

import (
	"net/http"
	"testing"
	"time"

	eventconfig "github.com/banggok/boillerplate_architecture/internal/config/event"
	"github.com/banggok/boillerplate_architecture/internal/data/model"
	restverify "github.com/banggok/boillerplate_architecture/internal/delivery/rest/v1/account/verify"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	asserthelper "github.com/banggok/boillerplate_architecture/tests/helper/assert_helper"
	"github.com/banggok/boillerplate_architecture/tests/helper/data"
	"github.com/banggok/boillerplate_architecture/tests/helper/request/account/verify"
	"github.com/banggok/boillerplate_architecture/tests/helper/response"
	"github.com/banggok/boillerplate_architecture/tests/helper/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountVerify(t *testing.T) {
	t.Run("failed when token invalid", func(t *testing.T) {
		server, cleanUp, mysqlCfg := setup.TestingEnv(t, true)
		defer cleanUp(mysqlCfg)

		resp := verify.CallGetAccountVerify(t, "invalid_token", server)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		// Read and parse the response body
		var errorResponse response.ErrorResponse

		response.GetBodyResponse(t, resp, &errorResponse)

		expectedResponse := response.ErrorResponse{
			RequestID: errorResponse.RequestID,
			Status:    "error",
			Code:      int(custom_errors.AccountVerificationNotFound),
			Message:   "account verification not found",
			Details: map[string]string{
				"error": "record not found",
			},
		}

		// Assert response body structure and values
		assert.Equal(t, expectedResponse, errorResponse)
	})

	t.Run("failed when token already verified", func(t *testing.T) {
		server, cleanUp, mysqlCfg := setup.TestingEnv(t, true)
		defer cleanUp(mysqlCfg)

		verification := data.AccountVerificationData
		verification.Verified = true

		account := data.AccountData
		account.AccountVerifications = &[]model.AccountVerification{verification}

		tenant := data.TenantData
		tenant.Accounts = &[]model.Account{account}

		data.PrepareData(t, mysqlCfg, tenant)

		token := data.Token

		resp := verify.CallGetAccountVerify(t, token, server)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		var errorResponse response.ErrorResponse
		response.GetBodyResponse(t, resp, &errorResponse)

		expectedError := response.ErrorResponse{
			RequestID: errorResponse.RequestID,
			Status:    "error",
			Code:      int(custom_errors.AccountVerificationConflictEntity),
			Message:   "account verification exists",
			Details: map[string]string{
				"error": "token has been verified",
			},
		}

		assert.Equal(t, expectedError, errorResponse)

	})

	t.Run("failed when token expired", func(t *testing.T) {
		server, cleanUp, mysqlCfg := setup.TestingEnv(t, true)
		defer cleanUp(mysqlCfg)

		verification := data.AccountVerificationData
		verification.ExpiresAt = time.Now().UTC().Add(-1 * time.Second)

		account := data.AccountData
		account.AccountVerifications = &[]model.AccountVerification{verification}

		tenant := data.TenantData
		tenant.Accounts = &[]model.Account{account}

		data.PrepareData(t, mysqlCfg, tenant)

		token := data.Token

		resp := verify.CallGetAccountVerify(t, token, server)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		var errorResponse response.ErrorResponse
		response.GetBodyResponse(t, resp, &errorResponse)

		expectedError := response.ErrorResponse{
			RequestID: errorResponse.RequestID,
			Status:    "error",
			Code:      int(custom_errors.AccountVerificationUnprocessEntity),
			Message:   "can not process account verification",
			Details: map[string]string{
				"error": "token expired",
			},
		}

		assert.Equal(t, expectedError, errorResponse)

	})

	t.Run("success", func(t *testing.T) {
		server, cleanUp, mysqlCfg := setup.TestingEnv(t, true)
		defer cleanUp(mysqlCfg)

		eventconfig.Setup(mysqlCfg.Master)

		verification := data.AccountVerificationData

		account := data.AccountData
		account.AccountVerifications = &[]model.AccountVerification{verification}

		tenant := data.TenantData
		tenant.Accounts = &[]model.Account{account}

		data.PrepareData(t, mysqlCfg, tenant)

		token := data.Token

		resp := verify.CallGetAccountVerify(t, token, server)

		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody struct {
			Status  string              `json:"status"`
			Message string              `json:"message"`
			Data    restverify.Response `json:"data"`
		}
		response.GetBodyResponse(t, resp, &responseBody)

		assert.Equal(t, "success", responseBody.Status)
		assert.Equal(t, "Account verified successfully", responseBody.Message)

		expected := restverify.Response{
			ID:        1,
			CreatedAt: responseBody.Data.CreatedAt,
			UpdatedAt: responseBody.Data.UpdatedAt,
		}

		asserthelper.AssertStruct(t, responseBody.Data, expected)

		var count int64
		err := mysqlCfg.Slave.Model(model.AccountVerification{}).
			Where("verified = ? and token = ?", true, token).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})
}
