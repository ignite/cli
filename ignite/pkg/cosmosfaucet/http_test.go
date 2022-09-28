package cosmosfaucet_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cosmosfaucet"
)

func TestServeHTTPCORS(t *testing.T) {
	f := cosmosfaucet.Faucet{}
	cases := []struct {
		name, method string
	}{
		{
			name:   "root endpoint",
			method: "POST",
		},
		{
			name:   "info endpoint",
			method: "GET",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			res := httptest.NewRecorder()
			req, _ := http.NewRequest("OPTIONS", "/", nil)
			req.Header.Set("Access-Control-Request-Method", tt.method)

			// Act
			f.ServeHTTP(res, req)

			// Assert
			require.Equal(t, http.StatusNoContent, res.Result().StatusCode)
		})
	}
}
