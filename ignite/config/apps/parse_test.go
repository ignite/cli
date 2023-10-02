package apps_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	appsconfig "github.com/ignite/cli/ignite/config/apps"
)

func TestParseDir(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedError string
		expectedApps  []appsconfig.App
		expectedPath  string
	}{
		{
			name:          "fail: path is not a dir",
			path:          "testdata/igniteapps.yml",
			expectedError: "app config parse: path testdata/igniteapps.yml is not a dir",
		},
		{
			name:          "fail: path doesn't exists",
			path:          "testdata/xxx/yyy",
			expectedError: "app config parse: stat testdata/xxx/yyy: no such file or directory",
		},
		{
			name:         "ok: path doesn't contain any config",
			path:         "testdata/noconfig",
			expectedApps: nil,
			expectedPath: "testdata/noconfig/igniteapps.yml",
		},
		{
			name:          "fail: path contains an invalid yml file",
			path:          "testdata/invalid",
			expectedError: "app config parse: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `not yaml !` into apps.Config",
		},
		{
			name: "ok: path contains a igniteapps.yml file",
			path: "testdata",
			expectedApps: []appsconfig.App{
				{
					Path: "/path/to/app1",
				},
				{
					Path: "/path/to/app2",
					With: map[string]string{
						"bar": "baz",
						"foo": "bar",
					},
				},
			},
			expectedPath: "testdata/igniteapps.yml",
		},
		{
			name: "ok: path contains a igniteapps.yaml file",
			path: "testdata/other",
			expectedApps: []appsconfig.App{
				{
					Path: "/path/to/app1",
				},
				{
					Path: "/path/to/app2",
					With: map[string]string{
						"bar": "baz",
						"foo": "bar",
					},
				},
			},
			expectedPath: "testdata/other/igniteapps.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			cfg, err := appsconfig.ParseDir(tt.path)

			if tt.expectedError != "" {
				require.EqualError(err, tt.expectedError)
				return
			}
			require.NoError(err)
			require.Equal(tt.expectedApps, cfg.Apps)
			require.Equal(tt.expectedPath, cfg.Path())
		})
	}
}
