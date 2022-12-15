package plugins_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pluginsconfig "github.com/ignite/cli/ignite/config/plugins"
)

func TestSplitPathRef(t *testing.T) {
	tests := []struct {
		name                  string
		plugin                pluginsconfig.Plugin
		expectedCanonicalPath string
		expectedRef           string
	}{
		{
			name:   "empty path",
			plugin: pluginsconfig.Plugin{},
		},
		{
			name: "simple path",
			plugin: pluginsconfig.Plugin{
				Path: "github.com/ignite/example",
			},
			expectedCanonicalPath: "github.com/ignite/example",
		},
		{
			name: "path with ref",
			plugin: pluginsconfig.Plugin{
				Path: "github.com/ignite/example@v1",
			},
			expectedCanonicalPath: "github.com/ignite/example",
			expectedRef:           "v1",
		},
		{
			name: "path with empty ref",
			plugin: pluginsconfig.Plugin{
				Path: "github.com/ignite/example@",
			},
			expectedCanonicalPath: "github.com/ignite/example",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			canonicalPath, ref := tt.plugin.SplitPathRef()

			assert.Equal(tt.expectedCanonicalPath, canonicalPath)
			assert.Equal(tt.expectedRef, ref)
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
