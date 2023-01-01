package plugin

import (
	"net"
	"testing"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/require"
)

func TestReadWriteConfigCache(t *testing.T) {
	t.Run("Should cache plugin config and read from cache", func(t *testing.T) {
		const path = "/path/to/awesome/plugin"
		unixFD, _ := net.ResolveUnixAddr("unix", "/var/folders/5k/sv4bxrs102n_6rr7430jc7j80000gn/T/plugin193424090")

		rc := hplugin.ReattachConfig{
			Protocol:        hplugin.ProtocolNetRPC,
			ProtocolVersion: hplugin.CoreProtocolVersion,
			Addr:            unixFD,
			Pid:             24464,
		}

		err := writeConfigCache(path, rc)
		require.NoError(t, err)

		c, err := readConfigCache(path)
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

		err := writeConfigCache(path, rc)
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

		err := writeConfigCache(path, rc)
		require.Error(t, err)
	})
}

func TestDeleteConfCache(t *testing.T) {
	t.Run("Delete plugin config after write to cache should remove from cache", func(t *testing.T) {
		const path = "/path/to/awesome/plugin"
		unixFD, _ := net.ResolveUnixAddr("unix", "/var/folders/5k/sv4bxrs102n_6rr7430jc7j80000gn/T/plugin193424090")

		rc := hplugin.ReattachConfig{
			Protocol:        hplugin.ProtocolNetRPC,
			ProtocolVersion: hplugin.CoreProtocolVersion,
			Addr:            unixFD,
			Pid:             24464,
		}

		err := writeConfigCache(path, rc)
		require.NoError(t, err)

		err = deleteConfCache(path)
		require.NoError(t, err)

		// there should be an error after deleting the config from the cache
		_, err = readConfigCache(path)
		require.Error(t, err)
	})

	t.Run("Delete plugin config should return error given empty path", func(t *testing.T) {
		const path = ""
		err := deleteConfCache(path)
		require.Error(t, err)
	})
}

func TestCheckConfCache(t *testing.T) {
	const path = "/path/to/awesome/plugin"
	unixFD, _ := net.ResolveUnixAddr("unix", "/var/folders/5k/sv4bxrs102n_6rr7430jc7j80000gn/T/plugin193424090")

	rc := hplugin.ReattachConfig{
		Protocol:        hplugin.ProtocolNetRPC,
		ProtocolVersion: hplugin.CoreProtocolVersion,
		Addr:            unixFD,
		Pid:             24464,
	}

	t.Run("Cache should be hydrated", func(t *testing.T) {
		err := writeConfigCache(path, rc)
		require.NoError(t, err)
		require.Equal(t, true, checkConfCache(path))
	})

	t.Run("Cache should be empty", func(t *testing.T) {
		_ = deleteConfCache(path)
		require.Equal(t, false, checkConfCache(path))
	})
}
