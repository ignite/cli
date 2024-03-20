package gomodule_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/gomodule"
)

func TestSplitPath(t *testing.T) {
	cases := []struct {
		name        string
		path        string
		wantPath    string
		wantVersion string
	}{
		{
			name:        "path with version",
			path:        "foo@v0.1.0",
			wantPath:    "foo",
			wantVersion: "v0.1.0",
		},
		{
			name:     "path without version",
			path:     "foo",
			wantPath: "foo",
		},
		{
			name: "invalid path",
			path: "@v0.1.0",
		},
		{
			name: "empty path",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			p, v := gomodule.SplitPath(tt.path)

			// Assert
			require.Equal(t, tt.wantPath, p)
			require.Equal(t, tt.wantVersion, v)
		})
	}
}

func TestJoinPath(t *testing.T) {
	require.Equal(t, "foo@v0.1.0", gomodule.JoinPath("foo", "v0.1.0"))
	require.Equal(t, "", gomodule.JoinPath("", "v0.1.0"))
	require.Equal(t, "foo", gomodule.JoinPath("foo", ""))
}

func TestFindModule(t *testing.T) {
	cases := []struct {
		name       string
		importPath string
		version    string
		wantErr    error
	}{
		{
			name:       "module exists",
			importPath: "github.com/gorilla/mux",
			version:    "v1.8.0",
		},
		{
			name:       "module exists with local replace",
			importPath: "../local-module-fork",
			version:    "",
		},
		{
			name:       "module missing",
			importPath: "github.com/foo/bar",
			version:    "v0.1.0",
			wantErr:    gomodule.ErrModuleNotFound,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			path := gomodule.JoinPath(tt.importPath, tt.version)

			// Act
			m, err := gomodule.FindModule(ctx, "testdata/module", path)

			// Assert
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.importPath, m.Path)
				require.Equal(t, tt.version, m.Version)
				require.True(t, strings.HasSuffix(m.Dir, path))
			}
		})
	}
}
