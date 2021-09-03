package gomodulepath

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
			err: fmt.Errorf("app name is an invalid go module name: %w",
				errors.New(`malformed module path "github.com/a/b/c@": invalid char '@'`)),
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
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			path, err := Parse(tt.rawpath)
			if err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.path, path)
		})
	}
}
