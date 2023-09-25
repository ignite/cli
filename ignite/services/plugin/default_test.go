package plugin

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDefaultPlugins(t *testing.T) {
	// Act
	plugins := GetDefaultPlugins()

	// Assert
	require.Greater(t, len(plugins), 0)
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
