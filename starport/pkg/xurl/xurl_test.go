package xurl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPEnsurePort(t *testing.T) {
	cases := []struct {
		addr    string
		ensured string
	}{
		{"http://localhost", "http://localhost:80"},
		{"https://localhost", "https://localhost:443"},
		{"http://localhost:4000", "http://localhost:4000"},
	}
	for _, tt := range cases {
		t.Run(tt.addr, func(t *testing.T) {
			addr, err := HTTPEnsurePort(tt.addr)
			require.NoError(t, err)
			require.Equal(t, tt.ensured, addr)
		})
	}
}
