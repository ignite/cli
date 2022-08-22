package xurl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPEnsurePort(t *testing.T) {
	cases := []struct {
		name string
		addr string
		want string
	}{
		{
			name: "http",
			addr: "http://localhost",
			want: "http://localhost:80",
		},
		{
			name: "https",
			addr: "https://localhost",
			want: "https://localhost:443",
		},
		{
			name: "custom",
			addr: "http://localhost:4000",
			want: "http://localhost:4000",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			addr := HTTPEnsurePort(tt.addr)
			require.Equal(t, tt.want, addr)
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
			addr: "tcp://github.com/ignite/cli",
			want: "tcp://github.com/ignite/cli",
		},
		{
			name: "without scheme",
			addr: "github.com/ignite/cli",
			want: "tcp://github.com/ignite/cli",
		},
		{
			name: "with invalid scheme",
			addr: "ftp://github.com/ignite/cli",
			want: "tcp://github.com/ignite/cli",
		},
		{
			name: "with ip and port",
			addr: "0.0.0.0:4500",
			want: "tcp://0.0.0.0:4500",
		},
		{
			name: "with localhost and port",
			addr: "localhost:4500",
			want: "tcp://localhost:4500",
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
			addr: "http://github.com/ignite/cli",
			want: "http://github.com/ignite/cli",
		},
		{
			name: "without scheme",
			addr: "github.com/ignite/cli",
			want: "http://github.com/ignite/cli",
		},
		{
			name: "with invalid scheme",
			addr: "ftp://github.com/ignite/cli",
			want: "http://github.com/ignite/cli",
		},
		{
			name: "with ip and port",
			addr: "0.0.0.0:4500",
			want: "http://0.0.0.0:4500",
		},
		{
			name: "with localhost and port",
			addr: "localhost:4500",
			want: "http://localhost:4500",
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
			addr: "https://github.com/ignite/cli",
			want: "https://github.com/ignite/cli",
		},
		{
			name: "without scheme",
			addr: "github.com/ignite/cli",
			want: "https://github.com/ignite/cli",
		},
		{
			name: "with invalid scheme",
			addr: "ftp://github.com/ignite/cli",
			want: "https://github.com/ignite/cli",
		},
		{
			name: "with ip and port",
			addr: "0.0.0.0:4500",
			want: "https://0.0.0.0:4500",
		},
		{
			name: "with localhost and port",
			addr: "localhost:4500",
			want: "https://localhost:4500",
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
			addr: "ws://github.com/ignite/cli",
			want: "ws://github.com/ignite/cli",
		},
		{
			name: "without scheme",
			addr: "github.com/ignite/cli",
			want: "ws://github.com/ignite/cli",
		},
		{
			name: "with invalid scheme",
			addr: "ftp://github.com/ignite/cli",
			want: "ws://github.com/ignite/cli",
		},
		{
			name: "with ip and port",
			addr: "0.0.0.0:4500",
			want: "ws://0.0.0.0:4500",
		},
		{
			name: "with localhost and port",
			addr: "localhost:4500",
			want: "ws://localhost:4500",
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

func TestMightHTTPS(t *testing.T) {
	cases := []struct {
		name  string
		addr  string
		want  string
		error bool
	}{
		{
			name: "with http scheme",
			addr: "http://github.com/ignite/cli",
			want: "http://github.com/ignite/cli",
		},
		{
			name: "with https scheme",
			addr: "https://github.com/ignite/cli",
			want: "https://github.com/ignite/cli",
		},
		{
			name: "without scheme",
			addr: "github.com/ignite/cli",
			want: "https://github.com/ignite/cli",
		},
		{
			name: "with invalid scheme",
			addr: "ftp://github.com/ignite/cli",
			want: "https://github.com/ignite/cli",
		},
		{
			name: "with ip and port",
			addr: "0.0.0.0:4500",
			want: "https://0.0.0.0:4500",
		},
		{
			name: "with localhost and port",
			addr: "localhost:4500",
			want: "https://localhost:4500",
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
			addr, err := MightHTTPS(tt.addr)
			if tt.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, addr)
			}
		})
	}
}

func Test_addressPort(t *testing.T) {
	tests := []struct {
		name     string
		arg      string
		wantHost string
		want     bool
	}{
		{
			name: "URI path",
			arg:  "/test/false",
			want: false,
		},
		{
			name: "invalid address",
			arg:  "aeihf3/aef/f..//",
			want: false,
		},
		{
			name:     "host and port",
			arg:      "102.33.3.43:10000",
			wantHost: "102.33.3.43:10000",
			want:     true,
		},
		{
			name:     "local port",
			arg:      "0.0.0.0:10000",
			wantHost: "0.0.0.0:10000",
			want:     true,
		},
		{
			name:     "only port",
			arg:      ":10000",
			wantHost: "0.0.0.0:10000",
			want:     true,
		},
		{
			name: "only host",
			arg:  "102.33.3.43",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, got := addressPort(tt.arg)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantHost, gotHost)
		})
	}
}
