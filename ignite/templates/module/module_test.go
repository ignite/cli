package module

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProtoPackageName(t *testing.T) {
	cases := []struct {
		name   string
		app    string
		module string
		want   string
	}{
		{
			name:   "name",
			app:    "ignite",
			module: "test",
			want:   "ignite.test",
		},
		{
			name:   "path",
			app:    "ignite/cli",
			module: "test",
			want:   "cli.test",
		},
		{
			name:   "path with dash",
			app:    "ignite/c-li",
			module: "test",
			want:   "cli.test",
		},
		{
			name:   "path with number prefix",
			app:    "0ignite/cli",
			module: "test",
			want:   "cli.test",
		},
		{
			name:   "app with number prefix",
			app:    "ignite/0cli",
			module: "test",
			want:   "_0cli.test",
		},
		{
			name:   "path with number prefix and dash",
			app:    "0ignite/cli",
			module: "test",
			want:   "cli.test",
		},
		{
			name:   "module with dash",
			app:    "ignite",
			module: "test-mod",
			want:   "ignite.testmod",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, ProtoPackageName(tt.app, tt.module))
		})
	}
}
