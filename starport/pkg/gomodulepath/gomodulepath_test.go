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
		{"standard",
			"github.com/a/b", Path{"github.com/a/b", "b", "b"}, nil,
		},
		{"with dash",
			"github.com/a/b-c", Path{"github.com/a/b-c", "b-c", "bc"}, nil,
		},
		{"long",
			"github.com/a/b/c", Path{"github.com/a/b/c", "c", "c"}, nil,
		},
		{"invalid as go.mod module name",
			"github.com/a/b/c@", Path{}, fmt.Errorf("app name is an invalid go module name: %w",
				errors.New(`malformed module path "github.com/a/b/c@": invalid char '@'`)),
		},
		{"name starting with the letter v",
			"github.com/a/vote", Path{"github.com/a/vote", "vote", "vote"}, nil,
		},
		{"with version",
			"github.com/a/b/v2", Path{"github.com/a/b/v2", "b", "b"}, nil,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			path, err := Parse(tt.rawpath)
			require.Equal(t, tt.err, err)
			if err != nil {
				return
			}
			require.Equal(t, tt.path, path)
		})
	}
}
