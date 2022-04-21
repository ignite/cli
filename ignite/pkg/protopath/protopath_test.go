package protopath

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatPackageName(t *testing.T) {
	cases := []struct {
		name   string
		owner  string
		app    string
		module string
		want   string
	}{
		{
			name:   "main",
			owner:  "ignite",
			app:    "cli",
			module: "cli",
			want:   "ignite.cli.cli",
		},
		{
			name:   "main with same owner and app names",
			owner:  "ignite",
			app:    "ignite",
			module: "cli",
			want:   "ignite.cli",
		},
		{
			name:   "main with dash",
			owner:  "ignite-hq",
			app:    "cli",
			module: "cli",
			want:   "ignitehq.cli.cli",
		},
		{
			name:   "main with number prefix",
			owner:  "0ignite",
			app:    "cli",
			module: "cli",
			want:   "_0ignite.cli.cli",
		},
		{
			name:   "main with number prefix and dash",
			owner:  "0ignite-hq",
			app:    "cli",
			module: "cli",
			want:   "_0ignitehq.cli.cli",
		},
		{
			name:   "module",
			owner:  "ignite",
			app:    "cli",
			module: "bank",
			want:   "ignite.cli.bank",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, FormatPackageName(tt.owner, tt.app, tt.module))
		})
	}
}
