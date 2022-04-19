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

func TestTCP(t *testing.T) {
	cases := []struct {
		name  string
		addr  string
		want  string
		error bool
	}{
		{
			name: "with scheme",
			addr: "tcp://github.com/ignite-hq/cli",
			want: "tcp://github.com/ignite-hq/cli",
		},
		{
			name: "without scheme",
			addr: "github.com/ignite-hq/cli",
			want: "tcp://github.com/ignite-hq/cli",
		},
		{
			name: "with invalid scheme",
			addr: "ftp://github.com/ignite-hq/cli",
			want: "tcp://github.com/ignite-hq/cli",
		},
		{
			name:  "with invalid url",
			addr:  "tcp://github.com:x",
			error: true,
		},
		{
			name:  "empty",
			addr:  "",
			error: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := TCP(tt.addr)
			if tt.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, addr)
			}
		})
	}
}

func TestHTTP(t *testing.T) {
	cases := []struct {
		name  string
		addr  string
		want  string
		error bool
	}{
		{
			name: "with scheme",
			addr: "http://github.com/ignite-hq/cli",
			want: "http://github.com/ignite-hq/cli",
		},
		{
			name: "without scheme",
			addr: "github.com/ignite-hq/cli",
			want: "http://github.com/ignite-hq/cli",
		},
		{
			name: "with invalid scheme",
			addr: "ftp://github.com/ignite-hq/cli",
			want: "http://github.com/ignite-hq/cli",
		},
		{
			name:  "with invalid url",
			addr:  "http://github.com:x",
			error: true,
		},
		{
			name:  "empty",
			addr:  "",
			error: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := HTTP(tt.addr)
			if tt.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, addr)
			}
		})
	}
}

func TestHTTPS(t *testing.T) {
	cases := []struct {
		name  string
		addr  string
		want  string
		error bool
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
			name: "with invalid scheme",
			addr: "ftp://github.com/ignite-hq/cli",
			want: "https://github.com/ignite-hq/cli",
		},
		{
			name:  "with invalid url",
			addr:  "https://github.com:x",
			error: true,
		},
		{
			name:  "empty",
			addr:  "",
			error: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := HTTPS(tt.addr)
			if tt.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, addr)
			}
		})
	}
}

func TestWS(t *testing.T) {
	cases := []struct {
		name  string
		addr  string
		want  string
		error bool
	}{
		{
			name: "with scheme",
			addr: "ws://github.com/ignite-hq/cli",
			want: "ws://github.com/ignite-hq/cli",
		},
		{
			name: "without scheme",
			addr: "github.com/ignite-hq/cli",
			want: "ws://github.com/ignite-hq/cli",
		},
		{
			name: "with invalid scheme",
			addr: "ftp://github.com/ignite-hq/cli",
			want: "ws://github.com/ignite-hq/cli",
		},
		{
			name:  "with invalid url",
			addr:  "ws://github.com:x",
			error: true,
		},
		{
			name:  "empty",
			addr:  "",
			error: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := WS(tt.addr)
			if tt.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, addr)
			}
		})
	}
}
