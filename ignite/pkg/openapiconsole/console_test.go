package openapiconsole

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlerRendersTitleAndSpecURL(t *testing.T) {
	h := Handler("My API", "https://example.com/openapi.json")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	body := rr.Body.String()
	require.Contains(t, body, "My API")
	require.Contains(t, body, "https://example.com/openapi.json")
}
