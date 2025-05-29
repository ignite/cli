package plugins_test

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pluginsconfig "github.com/ignite/cli/v29/ignite/config/plugins"
)

func TestPluginIsGlobal(t *testing.T) {
	assert.False(t, pluginsconfig.Plugin{}.IsGlobal())
	assert.True(t, pluginsconfig.Plugin{Global: true}.IsGlobal())
}

func TestPluginIsLocalPath(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)

	assert.False(t, pluginsconfig.Plugin{}.IsLocalPath())
	assert.False(t, pluginsconfig.Plugin{Path: "github.com/ignite/example"}.IsLocalPath())
	assert.False(t, pluginsconfig.Plugin{Path: "invalid_path"}.IsLocalPath())
	assert.False(t, pluginsconfig.Plugin{Path: "/testdata"}.IsLocalPath())
	assert.True(t, pluginsconfig.Plugin{Path: "testdata"}.IsLocalPath())
	assert.True(t, pluginsconfig.Plugin{Path: "/"}.IsLocalPath())
	assert.True(t, pluginsconfig.Plugin{Path: "."}.IsLocalPath())
	assert.True(t, pluginsconfig.Plugin{Path: ".."}.IsLocalPath())
	assert.True(t, pluginsconfig.Plugin{Path: pwd}.IsLocalPath())
	assert.True(t, pluginsconfig.Plugin{Path: filepath.Join(pwd, "testdata")}.IsLocalPath())
}

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
				t.Helper()
				return &pluginsconfig.Config{}
			},
			expectedError: "plugin config save: empty path",
		},
		{
			name: "ok: config path is a file that doesn't exist",
			buildConfig: func(t *testing.T) *pluginsconfig.Config {
				t.Helper()
				cfg, err := pluginsconfig.ParseDir(t.TempDir())
				require.NoError(t, err)
				return cfg
			},
			expectedContent: "apps: []\n",
		},
		{
			name: "ok: config path is an existing file",
			buildConfig: func(t *testing.T) *pluginsconfig.Config {
				t.Helper()
				// copy testdata/igniteapps.yml to tmp because it will be modified
				dir := t.TempDir()
				bz, err := os.ReadFile("testdata/igniteapps.yml")
				require.NoError(t, err)
				err = os.WriteFile(path.Join(dir, "igniteapps.yml"), bz, 0o666)
				require.NoError(t, err)
				// load from tmp
				cfg, _ := pluginsconfig.ParseDir(dir)
				// add a new plugin
				cfg.Apps = append(cfg.Apps, pluginsconfig.Plugin{
					Path: "/path/to/plugin3",
					With: map[string]string{"key": "val"},
				})
				// update a plugin
				cfg.Apps[1].Path = "/path/to/plugin22"
				cfg.Apps[1].With["key"] = "val"
				return cfg
			},
			expectedContent: `apps:
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

func TestConfigHasPlugin(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	tests := []struct {
		name          string
		cfg           pluginsconfig.Config
		expectedFound bool
	}{
		{
			name:          "empty config",
			expectedFound: false,
		},
		{
			name: "not found in config",
			cfg: pluginsconfig.Config{
				Apps: []pluginsconfig.Plugin{
					{Path: "github.com/ignite/example2"},
				},
			},
			expectedFound: false,
		},
		{
			name: "found in config",
			cfg: pluginsconfig.Config{
				Apps: []pluginsconfig.Plugin{
					{Path: "github.com/ignite/example2"},
					{Path: "github.com/ignite/example@master"},
				},
			},
			expectedFound: true,
		},
		{
			name: "found in config but from a local plugin",
			cfg: pluginsconfig.Config{
				Apps: []pluginsconfig.Plugin{
					{Path: "github.com/ignite/example2"},
					{Path: path.Join(wd, "testdata", "localplugin", "example")},
				},
			},
			expectedFound: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := tt.cfg.HasPlugin("github.com/ignite/example@v42")

			assert.Equal(t, tt.expectedFound, found)
		})
	}
}
