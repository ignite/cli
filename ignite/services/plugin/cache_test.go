package plugin_test

import (
	"net"
	"testing"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/services/plugin"
)

func TestPluginCacheAdd(t *testing.T) {
	const path = "/path/to/awesome/plugin"
	unixFD, _ := net.ResolveUnixAddr("unix", "/var/folders/5k/sv4bxrs102n_6rr7430jc7j80000gn/T/plugin193424090")

	rc := hplugin.ReattachConfig{
		Protocol:        hplugin.ProtocolNetRPC,
		ProtocolVersion: hplugin.CoreProtocolVersion,
		Addr:            unixFD,
		Pid:             24464,
	}

	err := plugin.WritePluginConfigCache(path, rc)
	require.NoError(t, err)

	c := hplugin.ReattachConfig{}
	err = plugin.ReadPluginConfigCache(path, &c)
	require.NoError(t, err)
}

func TestPluginCacheDelete(t *testing.T) {
	const path = "/path/to/awesome/plugin"
	unixFD, _ := net.ResolveUnixAddr("unix", "/var/folders/5k/sv4bxrs102n_6rr7430jc7j80000gn/T/plugin193424090")

	rc := hplugin.ReattachConfig{
		Protocol:        hplugin.ProtocolNetRPC,
		ProtocolVersion: hplugin.CoreProtocolVersion,
		Addr:            unixFD,
		Pid:             24464,
	}

	err := plugin.WritePluginConfigCache(path, rc)
	require.NoError(t, err)

	err = plugin.DeletePluginConfCache(path)
	require.NoError(t, err)

	c := hplugin.ReattachConfig{}
	// there should be an error after deleting the config from the cache
	err = plugin.ReadPluginConfigCache(path, &c)
	require.Error(t, err)
}

func TestPluginCacheCheck(t *testing.T) {
	const path = "/path/to/awesome/plugin"
	unixFD, _ := net.ResolveUnixAddr("unix", "/var/folders/5k/sv4bxrs102n_6rr7430jc7j80000gn/T/plugin193424090")

	rc := hplugin.ReattachConfig{
		Protocol:        hplugin.ProtocolNetRPC,
		ProtocolVersion: hplugin.CoreProtocolVersion,
		Addr:            unixFD,
		Pid:             24464,
	}

	t.Run("Cache should be hydrated", func(t *testing.T) {
		err := plugin.WritePluginConfigCache(path, rc)
		require.NoError(t, err)
		require.Equal(t, true, plugin.CheckPluginConfCache(path))
	})

	t.Run("Cache should be empty", func(t *testing.T) {
		_ = plugin.DeletePluginConfCache(path)
		require.Equal(t, false, plugin.CheckPluginConfCache(path))
	})
}
