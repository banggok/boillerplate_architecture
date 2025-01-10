package response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func GetBodyResponse(t *testing.T, resp http.Response, container interface{}) {
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	err = json.Unmarshal(body, container)
	require.NoError(t, err)
}
