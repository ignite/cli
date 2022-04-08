package xurl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPEnsurePort(t *testing.T) {
	fmt.Println(HTTPEnsurePort("https://26657-crimson-pheasant-2x3fbpak.ws-eu03.gitpod.io/"))
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
			addr := HTTPEnsurePort(tt.addr)
			require.Equal(t, tt.ensured, addr)
		})
	}
}
