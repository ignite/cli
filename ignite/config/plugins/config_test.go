package plugins_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	pluginsconfig "github.com/ignite/cli/ignite/config/plugins"
)

func TestPluginHasPath(t *testing.T) {
	tests := []struct {
		name        string
		plugin      pluginsconfig.Plugin
		path        string
		expectedRes bool
	}{
		{
			name:        "empty both path",
			plugin:      pluginsconfig.Plugin{},
			expectedRes: false,
		},
		{
			name: "simple path",
			plugin: pluginsconfig.Plugin{
				Path: "github.com/ignite/example",
			},
			path:        "github.com/ignite/example",
			expectedRes: true,
		},
		{
			name: "plugin path with ref",
			plugin: pluginsconfig.Plugin{
				Path: "github.com/ignite/example@v1",
			},
			path:        "github.com/ignite/example",
			expectedRes: true,
		},
		{
			name: "plugin path with empty ref",
			plugin: pluginsconfig.Plugin{
				Path: "github.com/ignite/example@",
			},
			path:        "github.com/ignite/example",
			expectedRes: true,
		},
		{
			name: "both path with different ref",
			plugin: pluginsconfig.Plugin{
				Path: "github.com/ignite/example@v1",
			},
			path:        "github.com/ignite/example@v2",
			expectedRes: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.plugin.HasPath(tt.path)
			require.Equal(t, tt.expectedRes, res)
		})
	}
}

func TestPluginCanonicalPath(t *testing.T) {
	tests := []struct {
		name         string
		plugin       pluginsconfig.Plugin
		expectedPath string
	}{
		{
			name:         "empty both path",
			plugin:       pluginsconfig.Plugin{},
			expectedPath: "",
		},
		{
			name: "simple path",
			plugin: pluginsconfig.Plugin{
				Path: "github.com/ignite/example",
			},
			expectedPath: "github.com/ignite/example",
		},
		{
			name: "plugin path with ref",
			plugin: pluginsconfig.Plugin{
				Path: "github.com/ignite/example@v1",
			},
			expectedPath: "github.com/ignite/example",
		},
		{
			name: "plugin path with empty ref",
			plugin: pluginsconfig.Plugin{
				Path: "github.com/ignite/example@",
			},
			expectedPath: "github.com/ignite/example",
		},
		{
			name: "plugin local directory path",
			plugin: pluginsconfig.Plugin{
				Path: "/home/user/go/foo/bar",
			},
			expectedPath: "/home/user/go/foo/bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.plugin.CanonicalPath()
			require.Equal(t, tt.expectedPath, res)
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		configs  []pluginsconfig.Plugin
		expected []pluginsconfig.Plugin
	}{
		{
			name:     "do nothing for empty list",
			configs:  []pluginsconfig.Plugin(nil),
			expected: []pluginsconfig.Plugin(nil),
		},
		{
			name: "remove duplicates",
			configs: []pluginsconfig.Plugin{
				{
					Path: "foo/bar",
				},
				{
					Path: "foo/bar",
				},
				{
					Path: "bar/foo",
				},
			},
			expected: []pluginsconfig.Plugin{
				{
					Path: "foo/bar",
				},
				{
					Path: "bar/foo",
				},
			},
		},
		{
			name: "do nothing for no duplicates",
			configs: []pluginsconfig.Plugin{
				{
					Path: "foo/bar",
				},
				{
					Path: "bar/foo",
				},
			},
			expected: []pluginsconfig.Plugin{
				{
					Path: "foo/bar",
				},
				{
					Path: "bar/foo",
				},
			},
		},
		{
			name: "prioritize local plugins",
			configs: []pluginsconfig.Plugin{
				{
					Path:   "foo/bar",
					Global: true,
				},
				{
					Path:   "bar/foo",
					Global: true,
				},
				{
					Path:   "foo/bar",
					Global: false,
				},
				{
					Path:   "bar/foo",
					Global: false,
				},
			},
			expected: []pluginsconfig.Plugin{
				{
					Path:   "foo/bar",
					Global: false,
				},
				{
					Path:   "bar/foo",
					Global: false,
				},
			},
		},
		{
			name: "prioritize local plugins different versions",
			configs: []pluginsconfig.Plugin{
				{
					Path:   "foo/bar@v1",
					Global: true,
				},
				{
					Path:   "bar/foo",
					Global: true,
				},
				{
					Path:   "foo/bar@v2",
					Global: false,
				},
				{
					Path:   "bar/foo",
					Global: false,
				},
			},
			expected: []pluginsconfig.Plugin{
				{
					Path:   "foo/bar@v2",
					Global: false,
				},
				{
					Path:   "bar/foo",
					Global: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unique := pluginsconfig.RemoveDuplicates(tt.configs)
			require.EqualValues(t, tt.expected, unique)
		})
	}
}

func TestConfigSave(t *testing.T) {
	tests := []struct {
		name            string
		buildConfig     func(*testing.T) *pluginsconfig.Config
		expectedError   string
		expectedContent string
	}{
		{
			name: "fail: config path is empty",
			buildConfig: func(t *testing.T) *pluginsconfig.Config {
				return &pluginsconfig.Config{}
			},
			expectedError: "plugin config save: empty path",
		},
		{
			name: "ok: config path is a file that doesn't exist",
			buildConfig: func(t *testing.T) *pluginsconfig.Config {
				cfg, err := pluginsconfig.ParseDir(t.TempDir())
				require.NoError(t, err)
				return cfg
			},
			expectedContent: "plugins: []\n",
		},
		{
			name: "ok: config path is an existing file",
			buildConfig: func(t *testing.T) *pluginsconfig.Config {
				// copy testdata/plugins.yml to tmp because it will be modified
				dir := t.TempDir()
				bz, err := os.ReadFile("testdata/plugins.yml")
				require.NoError(t, err)
				err = os.WriteFile(path.Join(dir, "plugins.yml"), bz, 0o666)
				require.NoError(t, err)
				// load from tmp
				cfg, _ := pluginsconfig.ParseDir(dir)
				// add a new plugin
				cfg.Plugins = append(cfg.Plugins, pluginsconfig.Plugin{
					Path: "/path/to/plugin3",
					With: map[string]string{"key": "val"},
				})
				// update a plugin
				cfg.Plugins[1].Path = "/path/to/plugin22"
				cfg.Plugins[1].With["key"] = "val"
				return cfg
			},
			expectedContent: `plugins:
- path: /path/to/plugin1
- path: /path/to/plugin22
  with:
    bar: baz
    foo: bar
    key: val
- path: /path/to/plugin3
  with:
    key: val
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			cfg := tt.buildConfig(t)

			err := cfg.Save()

			if tt.expectedError != "" {
				require.EqualError(err, tt.expectedError)
				return
			}
			require.NoError(err)
			bz, err := os.ReadFile(cfg.Path())
			require.NoError(err)
			require.Equal(string(bz), tt.expectedContent)
		})
	}
}
