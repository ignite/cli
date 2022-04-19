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
			addr := HTTPEnsurePort(tt.addr)
			require.Equal(t, tt.ensured, addr)
		})
	}
}

func TestHTTPS(t *testing.T) {
	cases := []struct {
		name string
		addr string
		want string
	}{
		{
			name: "with scheme",
			addr: "https://github.com/ignite-hq/cli",
			want: "https://github.com/ignite-hq/cli",
		},
		{
			name: "without scheme",
			addr: "github.com/ignite-hq/cli",
			want: "https://github.com/ignite-hq/cli",
		},
		{
			name: "empty",
			addr: "",
			want: "https://",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, HTTPS(tt.addr))
		})
	}
}
