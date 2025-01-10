package login

import (
	"net/http"
	"testing"

	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/v1/account/login"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/tests/helper/data"
	test_login "github.com/banggok/boillerplate_architecture/tests/helper/request/account/login"
	"github.com/banggok/boillerplate_architecture/tests/helper/response"
	"github.com/banggok/boillerplate_architecture/tests/helper/setup"
	"github.com/stretchr/testify/assert"
)

func TestAccountLogin(t *testing.T) {
	requestBody := login.Request{
		Email:    data.AccountData.Email,
		Password: data.PlainPassword,
	}
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

		requestBody.Email = ""
		requestBody.Password = ""

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
}
