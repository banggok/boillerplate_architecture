package response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetBodyResponse(t *testing.T, resp http.Response, container interface{}) {
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	err = json.Unmarshal(body, container)
	assert.NoError(t, err)
}
