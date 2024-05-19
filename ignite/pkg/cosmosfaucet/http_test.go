package cosmosfaucet_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosfaucet"
	"github.com/stretchr/testify/require"
)

func TestServeHTTPCORS(t *testing.T) {
	f := cosmosfaucet.Faucet{}
	cases := []struct {
		name, method, path string
	}{
		{
			name:   "root endpoint",
			method: "POST",
			path:   "/",
		},
		{
			name:   "info endpoint",
			method: "GET",
			path:   "/info",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			res := httptest.NewRecorder()
			req, _ := http.NewRequest("OPTIONS", tt.path, nil)
			req.Header.Set("Access-Control-Request-Method", tt.method)

			// Act
			f.ServeHTTP(res, req)

			// Assert
			require.Equal(t, http.StatusNoContent, res.Result().StatusCode)
		})
	}
}
