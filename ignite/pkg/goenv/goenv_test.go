package goenv_test

import (
	"go/build"
	"path/filepath"
	"testing"

	"github.com/ignite/cli/v29/ignite/pkg/goenv"
	"github.com/stretchr/testify/require"
)

func TestGoModCache(t *testing.T) {
	cases := []struct {
		name, envKey, envValue, want string
	}{
		{
			name:     "from go module cache",
			envKey:   "GOMODCACHE",
			envValue: "/foo/cache/pkg/mod",
			want:     "/foo/cache/pkg/mod",
		},
		{
			name:     "from go path",
			envKey:   "GOPATH",
			envValue: "/foo/go",
			want:     "/foo/go/pkg/mod",
		},
		{
			name: "from default path",
			want: filepath.Join(build.Default.GOPATH, "pkg/mod"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			if tt.envKey != "" {
				t.Setenv(tt.envKey, tt.envValue)
			}

			// Act
			path := goenv.GoModCache()

			// Assert
			require.Equal(t, tt.want, path)
		})
	}
}
