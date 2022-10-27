package gomodulepath

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/mod/module"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name    string
		rawpath string
		path    Path
		err     error
	}{
		{
			name:    "standard",
			rawpath: "github.com/a/b",
			path:    Path{RawPath: "github.com/a/b", Root: "b", Package: "b"},
		},
		{
			name:    "with dash",
			rawpath: "github.com/a/b-c",
			path:    Path{RawPath: "github.com/a/b-c", Root: "b-c", Package: "bc"},
		},
		{
			name:    "short",
			rawpath: "github.com/a",
			path:    Path{RawPath: "github.com/a", Root: "a", Package: "a"},
		},
		{
			name:    "short with dash",
			rawpath: "github.com/a-c",
			path:    Path{RawPath: "github.com/a-c", Root: "a-c", Package: "ac"},
		},
		{
			name:    "short with version",
			rawpath: "github.com/a/v2",
			path:    Path{RawPath: "github.com/a/v2", Root: "a", Package: "a"},
		},
		{
			name:    "long",
			rawpath: "github.com/a/b/c",
			path:    Path{RawPath: "github.com/a/b/c", Root: "c", Package: "c"},
		},
		{
			name:    "invalid as go.mod module name",
			rawpath: "github.com/a/b/c@",
			err: &module.InvalidPathError{
				Kind: "module",
				Path: "github.com/a/b/c@",
				Err:  fmt.Errorf("invalid char '@'"),
			},
		},
		{
			name:    "name starting with the letter v",
			rawpath: "github.com/a/vote",
			path:    Path{RawPath: "github.com/a/vote", Root: "vote", Package: "vote"},
		},
		{
			name:    "with version",
			rawpath: "github.com/a/b/v2",
			path:    Path{RawPath: "github.com/a/b/v2", Root: "b", Package: "b"},
		},
		{
			name:    "with underscore",
			rawpath: "github.com/a/b_c",
			path:    Path{RawPath: "github.com/a/b_c", Root: "b_c", Package: "bc"},
		},
		{
			name:    "with mixed case",
			rawpath: "github.com/a/bC",
			path:    Path{RawPath: "github.com/a/bC", Root: "bC", Package: "bc"},
		},
		{
			name:    "with a name",
			rawpath: "a",
			path:    Path{RawPath: "a", Root: "a", Package: "a"},
		},
		{
			name:    "with a name containing underscore",
			rawpath: "a_b",
			path:    Path{RawPath: "a_b", Root: "a_b", Package: "ab"},
		},
		{
			name:    "with a name containing dash",
			rawpath: "a-b",
			path:    Path{RawPath: "a-b", Root: "a-b", Package: "ab"},
		},
		{
			name:    "with a path",
			rawpath: "a/b/c",
			path:    Path{RawPath: "a/b/c", Root: "c", Package: "c"},
		},
		{
			name:    "with a path containing underscore",
			rawpath: "a/b_c",
			path:    Path{RawPath: "a/b_c", Root: "b_c", Package: "bc"},
		},
		{
			name:    "with a path containing dash",
			rawpath: "a/b-c",
			path:    Path{RawPath: "a/b-c", Root: "b-c", Package: "bc"},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			path, err := Parse(tt.rawpath)
			if err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err, errors.Unwrap(err))
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.path, path)
		})
	}
}

func TestExtractAppPath(t *testing.T) {
	cases := []struct {
		name string
		path string
		want string
	}{
		{
			name: "github uri",
			path: "github.com/ignite/cli",
			want: "ignite/cli",
		},
		{
			name: "short uri",
			path: "domain.com/ignite",
			want: "ignite",
		},
		{
			name: "long uri",
			path: "domain.com/a/b/c/ignite/cli",
			want: "ignite/cli",
		},
		{
			name: "name",
			path: "cli",
			want: "cli",
		},
		{
			name: "path",
			path: "ignite/cli",
			want: "ignite/cli",
		},
		{
			name: "long path",
			path: "a/b/c/ignite/cli",
			want: "ignite/cli",
		},
		{
			name: "empty",
			path: "",
			want: "",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, ExtractAppPath(tt.path))
		})
	}
}

func TestValidateURIPath(t *testing.T) {
	require.NoError(t, validateURIPath("github.com/ignite/cli"))
}

func TestValidateURIPathWithInvalidPath(t *testing.T) {
	require.Error(t, validateURIPath("github/ignite/cli"))
}

func TestValidateNamePath(t *testing.T) {
	require.NoError(t, validateNamePath("cli"))
}

func TestValidateNamePathWithInvalidPath(t *testing.T) {
	require.Error(t, validateNamePath("cli@"))
}
