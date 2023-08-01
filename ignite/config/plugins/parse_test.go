package plugins_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	pluginsconfig "github.com/ignite/cli/ignite/config/plugins"
)

func TestParseDir(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		expectedError   string
		expectedPlugins []pluginsconfig.Plugin
		expectedPath    string
	}{
		{
			name:          "fail: path is not a dir",
			path:          "testdata/plugins.yml",
			expectedError: "plugin config parse: path testdata/plugins.yml is not a dir",
		},
		{
			name:          "fail: path doesn't exists",
			path:          "testdata/xxx/yyy",
			expectedError: "plugin config parse: stat testdata/xxx/yyy: no such file or directory",
		},
		{
			name:            "ok: path doesn't contain any config",
			path:            "testdata/noconfig",
			expectedPlugins: nil,
			expectedPath:    "testdata/noconfig/plugins.yml",
		},
		{
			name:          "fail: path contains an invalid yml file",
			path:          "testdata/invalid",
			expectedError: "plugin config parse: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `not yaml !` into plugins.Config",
		},
		{
			name: "ok: path contains a plugin.yml file",
			path: "testdata",
			expectedPlugins: []pluginsconfig.Plugin{
				{
					Path: "/path/to/plugin1",
				},
				{
					Path: "/path/to/plugin2",
					With: map[string]string{
						"bar": "baz",
						"foo": "bar",
					},
				},
			},
			expectedPath: "testdata/plugins.yml",
		},
		{
			name: "ok: path contains a plugin.yaml file",
			path: "testdata/other",
			expectedPlugins: []pluginsconfig.Plugin{
				{
					Path: "/path/to/plugin1",
				},
				{
					Path: "/path/to/plugin2",
					With: map[string]string{
						"bar": "baz",
						"foo": "bar",
					},
				},
			},
			expectedPath: "testdata/other/plugins.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			cfg, err := pluginsconfig.ParseDir(tt.path)

			if tt.expectedError != "" {
				require.EqualError(err, tt.expectedError)
				return
			}
			require.NoError(err)
			require.Equal(tt.expectedPlugins, cfg.Plugins)
			require.Equal(tt.expectedPath, cfg.Path())
		})
	}
}
