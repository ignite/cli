package xhttp

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

func TestResponseJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{"a": 1}
	require.NoError(t, ResponseJSON(w, http.StatusCreated, data))
	resp := w.Result()
	defer resp.Body.Close() // Ensure the response body is closed

	require.Equal(t, http.StatusCreated, resp.StatusCode)
	require.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	body, _ := io.ReadAll(resp.Body)
	dataJSON, err := json.Marshal(data)
	require.NoError(t, err)
	require.Equal(t, dataJSON, body)
}

func TestNewErrorResponse(t *testing.T) {
	require.Equal(t, ErrorResponseBody{
		Error: ErrorResponse{
			Message: "error",
		},
	}, NewErrorResponse(errors.New("error")))
}
