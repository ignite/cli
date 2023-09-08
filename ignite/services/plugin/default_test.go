package plugin

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePluginsMarkdown(t *testing.T) {
	// Arrange
	markdown := []byte(`
- [foo](github.com/foo/ignite-app-foo/tree/v0.0.1): Foo description
- [bar](github.com/bar/ignite-app-bar/tree/v0.1.0): Bar description
	`)

	cases := []struct {
		name        string
		markdown    []byte
		plugins     []DefaultPlugin
		expectError bool
	}{
		{
			name:     "ok",
			markdown: markdown,
			plugins: []DefaultPlugin{
				{
					Use:     "foo",
					Short:   "Foo description",
					Aliases: []string{"f"},
					Path:    "github.com/foo/ignite-app-foo@v0.0.1",
				},
				{
					Use:     "bar",
					Short:   "Bar description",
					Aliases: []string{"b"},
					Path:    "github.com/bar/ignite-app-bar@v0.1.0",
				},
			},
		}, {
			name:        "empty",
			expectError: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			plugins, err := parsePluginsMarkdown(tt.markdown)

			// Assert
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.plugins, plugins)
		})
	}
}

func TestParsePluginRepoURL(t *testing.T) {
	// Arrange
	cases := []struct {
		name, url, want string
	}{
		{
			name: "ok",
			url:  "github.com/ignite/cli-plugin-network/tree/v0.1.1",
			want: "github.com/ignite/cli-plugin-network@v0.1.1",
		}, {
			name: "invalid",
			url:  "github.com/ignite/cli-plugin-network",
		}, {
			name: "empty",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			pkg := parsePluginRepoURL(tt.url)

			// Assert
			require.Equal(t, tt.want, pkg)
		})
	}
}

func TestGetDefaultPlugins(t *testing.T) {
	// Act
	plugins, err := GetDefaultPlugins()

	// Assert
	require.NoError(t, err)
	require.Len(t, plugins, 1)
}

func TestGetDefaultNetworkPlugin(t *testing.T) {
	// Act
	p, err := GetDefaultNetworkPlugin()

	// Assert
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(p.Path, "github.com/ignite/cli-plugin-network@"))
	require.Equal(t, "network", p.Use)
	require.Equal(t, []string{"n"}, p.Aliases)
	require.NotEmpty(t, p.Short)
}
