package module

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProtoPackageName(t *testing.T) {
	cases := []struct {
		name    string
		app     string
		module  string
		version string
		want    string
	}{
		{
			name:    "name",
			app:     "ignite",
			module:  "test",
			version: "v1",
			want:    "ignite.test.v1",
		},
		{
			name:    "name",
			app:     "ignite",
			module:  "test",
			version: "v2",
			want:    "ignite.test.v2",
		},
		{
			name:    "path",
			app:     "ignite/cli",
			module:  "test",
			version: "v1",
			want:    "cli.test.v1",
		},
		{
			name:    "path with dash",
			app:     "ignite/c-li",
			module:  "test",
			version: "v1",
			want:    "cli.test.v1",
		},
		{
			name:    "path with number prefix",
			app:     "0ignite/cli",
			module:  "test",
			version: "v1",
			want:    "cli.test.v1",
		},
		{
			name:    "app with number prefix",
			app:     "ignite/0cli",
			module:  "test",
			version: "v1",
			want:    "_0cli.test.v1",
		},
		{
			name:    "path with number prefix and dash",
			app:     "0ignite/cli",
			module:  "test",
			version: "v1",
			want:    "cli.test.v1",
		},
		{
			name:    "module with dash",
			app:     "ignite",
			module:  "test-mod",
			version: "v1",
			want:    "ignite.testmod.v1",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, ProtoPackageName(tt.app, tt.module, tt.version))
		})
	}
}

func TestProtoModulePackageName(t *testing.T) {
	cases := []struct {
		name    string
		app     string
		module  string
		version string
		want    string
	}{
		{
			name:    "name",
			app:     "ignite",
			module:  "test",
			version: "v1",
			want:    "ignite.test.module.v1",
		},
		{
			name:    "name",
			app:     "ignite",
			module:  "test",
			version: "v2",
			want:    "ignite.test.module.v2",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, ProtoModulePackageName(tt.app, tt.module, tt.version))
		})
	}
}
