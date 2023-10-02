package apps_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appsconfig "github.com/ignite/cli/ignite/config/apps"
)

func TestAppIsGlobal(t *testing.T) {
	assert.False(t, appsconfig.App{}.IsGlobal())
	assert.True(t, appsconfig.App{Global: true}.IsGlobal())
}

func TestAppIsLocalPath(t *testing.T) {
	assert.False(t, appsconfig.App{}.IsLocalPath())
	assert.False(t, appsconfig.App{Path: "github.com/ignite/example"}.IsLocalPath())
	assert.True(t, appsconfig.App{Path: "/home/bob/example"}.IsLocalPath())
}

func TestAppHasPath(t *testing.T) {
	tests := []struct {
		name        string
		app         appsconfig.App
		path        string
		expectedRes bool
	}{
		{
			name:        "empty both path",
			app:         appsconfig.App{},
			expectedRes: false,
		},
		{
			name: "simple path",
			app: appsconfig.App{
				Path: "github.com/ignite/example",
			},
			path:        "github.com/ignite/example",
			expectedRes: true,
		},
		{
			name: "app path with ref",
			app: appsconfig.App{
				Path: "github.com/ignite/example@v1",
			},
			path:        "github.com/ignite/example",
			expectedRes: true,
		},
		{
			name: "app path with empty ref",
			app: appsconfig.App{
				Path: "github.com/ignite/example@",
			},
			path:        "github.com/ignite/example",
			expectedRes: true,
		},
		{
			name: "both path with different ref",
			app: appsconfig.App{
				Path: "github.com/ignite/example@v1",
			},
			path:        "github.com/ignite/example@v2",
			expectedRes: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.app.HasPath(tt.path)

			require.Equal(t, tt.expectedRes, res)
		})
	}
}

func TestAppCanonicalPath(t *testing.T) {
	tests := []struct {
		name         string
		app          appsconfig.App
		expectedPath string
	}{
		{
			name:         "empty both path",
			app:          appsconfig.App{},
			expectedPath: "",
		},
		{
			name: "simple path",
			app: appsconfig.App{
				Path: "github.com/ignite/example",
			},
			expectedPath: "github.com/ignite/example",
		},
		{
			name: "app path with ref",
			app: appsconfig.App{
				Path: "github.com/ignite/example@v1",
			},
			expectedPath: "github.com/ignite/example",
		},
		{
			name: "app path with empty ref",
			app: appsconfig.App{
				Path: "github.com/ignite/example@",
			},
			expectedPath: "github.com/ignite/example",
		},
		{
			name: "app local directory path",
			app: appsconfig.App{
				Path: "/home/user/go/foo/bar",
			},
			expectedPath: "/home/user/go/foo/bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.app.CanonicalPath()
			require.Equal(t, tt.expectedPath, res)
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		configs  []appsconfig.App
		expected []appsconfig.App
	}{
		{
			name:     "do nothing for empty list",
			configs:  []appsconfig.App(nil),
			expected: []appsconfig.App(nil),
		},
		{
			name: "remove duplicates",
			configs: []appsconfig.App{
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
			expected: []appsconfig.App{
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
			configs: []appsconfig.App{
				{
					Path: "foo/bar",
				},
				{
					Path: "bar/foo",
				},
			},
			expected: []appsconfig.App{
				{
					Path: "foo/bar",
				},
				{
					Path: "bar/foo",
				},
			},
		},
		{
			name: "prioritize local apps",
			configs: []appsconfig.App{
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
			expected: []appsconfig.App{
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
			name: "prioritize local apps different versions",
			configs: []appsconfig.App{
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
			expected: []appsconfig.App{
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
			unique := appsconfig.RemoveDuplicates(tt.configs)
			require.EqualValues(t, tt.expected, unique)
		})
	}
}

func TestConfigSave(t *testing.T) {
	tests := []struct {
		name            string
		buildConfig     func(*testing.T) *appsconfig.Config
		expectedError   string
		expectedContent string
	}{
		{
			name: "fail: config path is empty",
			buildConfig: func(t *testing.T) *appsconfig.Config {
				return &appsconfig.Config{}
			},
			expectedError: "app config save: empty path",
		},
		{
			name: "ok: config path is a file that doesn't exist",
			buildConfig: func(t *testing.T) *appsconfig.Config {
				cfg, err := appsconfig.ParseDir(t.TempDir())
				require.NoError(t, err)
				return cfg
			},
			expectedContent: "apps: []\n",
		},
		{
			name: "ok: config path is an existing file",
			buildConfig: func(t *testing.T) *appsconfig.Config {
				// copy testdata/igniteapps.yml to tmp because it will be modified
				dir := t.TempDir()
				bz, err := os.ReadFile("testdata/igniteapps.yml")
				require.NoError(t, err)
				err = os.WriteFile(path.Join(dir, "igniteapps.yml"), bz, 0o666)
				require.NoError(t, err)
				// load from tmp
				cfg, _ := appsconfig.ParseDir(dir)
				// add a new app
				cfg.Apps = append(cfg.Apps, appsconfig.App{
					Path: "/path/to/app3",
					With: map[string]string{"key": "val"},
				})
				// update an app
				cfg.Apps[1].Path = "/path/to/app22"
				cfg.Apps[1].With["key"] = "val"
				return cfg
			},
			expectedContent: `apps:
- path: /path/to/app1
- path: /path/to/app22
  with:
    bar: baz
    foo: bar
    key: val
- path: /path/to/app3
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

func TestConfigHasApp(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	tests := []struct {
		name          string
		cfg           appsconfig.Config
		expectedFound bool
	}{
		{
			name:          "empty config",
			expectedFound: false,
		},
		{
			name: "not found in config",
			cfg: appsconfig.Config{
				Apps: []appsconfig.App{
					{Path: "github.com/ignite/example2"},
				},
			},
			expectedFound: false,
		},
		{
			name: "found in config",
			cfg: appsconfig.Config{
				Apps: []appsconfig.App{
					{Path: "github.com/ignite/example2"},
					{Path: "github.com/ignite/example@master"},
				},
			},
			expectedFound: true,
		},
		{
			name: "found in config but from a local app",
			cfg: appsconfig.Config{
				Apps: []appsconfig.App{
					{Path: "github.com/ignite/example2"},
					{Path: path.Join(wd, "testdata", "localapp", "example")},
				},
			},
			expectedFound: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := tt.cfg.HasApp("github.com/ignite/example@v42")

			assert.Equal(t, tt.expectedFound, found)
		})
	}
}
