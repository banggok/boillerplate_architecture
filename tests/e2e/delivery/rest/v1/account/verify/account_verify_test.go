package verify_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/banggok/boillerplate_architecture/internal/config/db"
	eventconfig "github.com/banggok/boillerplate_architecture/internal/config/event"
	valueobject "github.com/banggok/boillerplate_architecture/internal/data/entity/value_object"
	"github.com/banggok/boillerplate_architecture/internal/data/model"
	restverify "github.com/banggok/boillerplate_architecture/internal/delivery/rest/v1/account/verify"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/password"
	asserthelper "github.com/banggok/boillerplate_architecture/tests/helper/assert_helper"
	"github.com/banggok/boillerplate_architecture/tests/helper/request/account/verify"
	"github.com/banggok/boillerplate_architecture/tests/helper/response"
	"github.com/banggok/boillerplate_architecture/tests/helper/setup"
	time_testing_helper "github.com/banggok/boillerplate_architecture/tests/helper/time"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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

		token := prepareAccountData(t, mysqlCfg, true)

		resp := verify.CallGetAccountVerify(t, *token, server)

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

	t.Run("success", func(t *testing.T) {
		server, cleanUp, mysqlCfg := setup.TestingEnv(t, true)
		defer cleanUp(mysqlCfg)
		eventconfig.Setup(mysqlCfg.Master)

		token := prepareAccountData(t, mysqlCfg, false)

		resp := verify.CallGetAccountVerify(t, *token, server)

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
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		time_testing_helper.Sanitize(t, &responseBody.Data.CreatedAt, expected.CreatedAt)
		time_testing_helper.Sanitize(t, &responseBody.Data.UpdatedAt, expected.UpdatedAt)

		asserthelper.AssertStruct(t, responseBody.Data, expected)

		var count int64
		err := mysqlCfg.Slave.Model(model.AccountVerification{}).
			Where("verified = ? and token = ?", true, token).Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})
}

func prepareAccountData(t *testing.T, mysqlCfg *db.DBConnection, verified bool) *string {
	token, err := password.GeneratePassword(16)
	assert.NoError(t, err)

	mysqlCfg.Master.Session(&gorm.Session{FullSaveAssociations: true}).Save(&model.Tenant{
		Name:         "Tenant",
		Address:      "Address",
		Email:        "email@testing.com",
		Phone:        "081234567890",
		Timezone:     "UTC",
		OpeningHours: "08:00",
		ClosingHours: "17:00",
		Accounts: &[]model.Account{
			{
				Name:     "Name",
				Email:    "email@testing.com",
				Phone:    "08123456789",
				Password: "password",
				AccountVerifications: &[]model.AccountVerification{
					{
						Type:      valueobject.EMAIL_VERIFICATION.String(),
						Token:     *token,
						ExpiresAt: time.Now().Add(24 * time.Hour).UTC(),
						Verified:  verified,
					},
				},
			},
		},
	})
	return token
}
