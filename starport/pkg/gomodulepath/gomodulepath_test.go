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

func TestExtractUserAndRepo(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		want     string
		mustFail bool
	}{
		{
			name: "url",
			path: "github.com/tendermint/starport",
			want: "tendermint/starport",
		},
		{
			name: "name",
			path: "starport",
			want: "starport/starport",
		},
		{
			name:     "invalid url",
			path:     "github.com/tendermint",
			mustFail: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			path, err := ExtractUserAndRepo(tt.path)

			if tt.mustFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, path)
		})
	}
}

func TestValidateURLPath(t *testing.T) {
	require.NoError(t, validateURLPath("github.com/tendermint/starport"))
}

func TestValidateURLPathWithInvalidPath(t *testing.T) {
	require.Error(t, validateURLPath("github/tendermint/starport"))
}

func TestValidateNamePath(t *testing.T) {
	require.NoError(t, validateNamePath("starport"))
}

func TestValidateNamePathWithInvalidPath(t *testing.T) {
	require.Error(t, validateNamePath("starport."))
}
