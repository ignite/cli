package plugin

import (
	"net"
	"testing"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"
)

func TestReadWritePluginConfigCache(t *testing.T) {
	t.Run("Should cache plugin config and read from cache", func(t *testing.T) {
		const path = "/path/to/awesome/plugin"
		unixFD, _ := net.ResolveUnixAddr("unix", "/var/folders/5k/sv4bxrs102n_6rr7430jc7j80000gn/T/plugin193424090")

		rc := hplugin.ReattachConfig{
			Protocol:        hplugin.ProtocolNetRPC,
			ProtocolVersion: hplugin.CoreProtocolVersion,
			Addr:            unixFD,
			Pid:             24464,
		}

		err := WritePluginConfigCache(path, rc)
		require.NoError(t, err)

		c := hplugin.ReattachConfig{}
		err = ReadPluginConfigCache(path, &c)
		require.NoError(t, err)
		require.Equal(t, rc, c)
	})

	t.Run("Should error writing bad plugin config to cache", func(t *testing.T) {
		const path = "/path/to/awesome/plugin"
		rc := hplugin.ReattachConfig{
			Protocol:        hplugin.ProtocolNetRPC,
			ProtocolVersion: hplugin.CoreProtocolVersion,
			Addr:            nil,
			Pid:             24464,
		}

		err := WritePluginConfigCache(path, rc)
		require.Error(t, err)
	})

	t.Run("Should error with invalid plugin path", func(t *testing.T) {
		const path = ""
		rc := hplugin.ReattachConfig{
			Protocol:        hplugin.ProtocolNetRPC,
			ProtocolVersion: hplugin.CoreProtocolVersion,
			Addr:            nil,
			Pid:             24464,
		}

		err := WritePluginConfigCache(path, rc)
		require.Error(t, err)
	})
}

func TestPluginCacheDelete(t *testing.T) {
	t.Run("Delete plugin config after write to cache should remove from cache", func(t *testing.T) {
		const path = "/path/to/awesome/plugin"
		unixFD, _ := net.ResolveUnixAddr("unix", "/var/folders/5k/sv4bxrs102n_6rr7430jc7j80000gn/T/plugin193424090")

		rc := hplugin.ReattachConfig{
			Protocol:        hplugin.ProtocolNetRPC,
			ProtocolVersion: hplugin.CoreProtocolVersion,
			Addr:            unixFD,
			Pid:             24464,
		}

		err := WritePluginConfigCache(path, rc)
		require.NoError(t, err)

		err = DeletePluginConfCache(path)
		require.NoError(t, err)

		c := hplugin.ReattachConfig{}
		// there should be an error after deleting the config from the cache
		err = ReadPluginConfigCache(path, &c)
		require.Error(t, err)
	})

	t.Run("Delete plugin config should return error given empty path", func(t *testing.T) {
		const path = ""
		err := DeletePluginConfCache(path)
		require.Error(t, err)
	})
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
		err := WritePluginConfigCache(path, rc)
		require.NoError(t, err)
		require.Equal(t, true, CheckPluginConfCache(path))
	})

	t.Run("Cache should be empty", func(t *testing.T) {
		_ = DeletePluginConfCache(path)
		require.Equal(t, false, CheckPluginConfCache(path))
	})
}
