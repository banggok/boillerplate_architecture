package test_login

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/v1/account/login"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var URL = "/api/v1/accounts/login"

func CallRequest(t *testing.T, reqBody login.Request,
	server *gin.Engine) http.Response {
	payload, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, req)

	resp := recorder.Result()
	return *resp
}

func CallRequestWithRawPayload(t *testing.T, rawPayload string, server *gin.Engine) http.Response {
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer([]byte(rawPayload)))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, req)

	resp := recorder.Result()
	return *resp
}
