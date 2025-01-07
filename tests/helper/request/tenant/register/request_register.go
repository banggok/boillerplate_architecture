package test_register

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/v1/tenant/register"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func CallPostTenants(t *testing.T, reqBody register.Request,
	server *gin.Engine) http.Response {
	payload, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewBuffer(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, req)

	resp := recorder.Result()
	return *resp
}

func CallPostTenantsWithRawPayload(t *testing.T, rawPayload string, server *gin.Engine) http.Response {
	req, err := http.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewBuffer([]byte(rawPayload)))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, req)

	resp := recorder.Result()
	return *resp
}
