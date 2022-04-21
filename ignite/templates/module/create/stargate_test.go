package modulecreate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatAPIPath(t *testing.T) {
	cases := []struct {
		name   string
		owner  string
		app    string
		module string
		want   string
	}{
		{
			name:   "app",
			owner:  "ignite",
			app:    "cli",
			module: "cli",
			want:   "/ignite/cli/cli",
		},
		{
			name:   "app with same owner and app names",
			owner:  "ignite",
			app:    "ignite",
			module: "cli",
			want:   "/ignite/cli",
		},
		{
			name:   "module",
			owner:  "ignite",
			app:    "cli",
			module: "bank",
			want:   "/ignite/cli/bank",
		},
		{
			name:   "same names",
			owner:  "cli",
			app:    "cli",
			module: "cli",
			want:   "/cli/cli",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, formatAPIPath(tt.owner, tt.app, tt.module))
		})
	}
}
