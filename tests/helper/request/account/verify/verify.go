package verify

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func CallGetAccountVerify(t *testing.T, token string, server *gin.Engine) http.Response {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/accounts/verify/%s", token), nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, req)

	resp := recorder.Result()
	return *resp
}
