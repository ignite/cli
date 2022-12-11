package plugin_test

import (
	"net"
	"testing"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/ignite/services/plugin"
	"github.com/stretchr/testify/require"
)

func TestNewPlugin(t *testing.T) {
	const path = "/path/to/awesome/plugin"
	unixFD, _ := net.ResolveUnixAddr("unix", "/var/folders/5k/sv4bxrs102n_6rr7430jc7j80000gn/T/plugin193424090")

	var rc = hplugin.ReattachConfig{
		Protocol:        hplugin.ProtocolNetRPC,
		ProtocolVersion: hplugin.CoreProtocolVersion,
		Addr:            unixFD,
		Pid:             24464,
	}

	err := plugin.WritePluginConfig(path, rc)
	require.NoError(t, err)

	c := hplugin.ReattachConfig{}
	err = plugin.ReadPluginConfig(path, &c)
	require.NoError(t, err)

	err = plugin.DeletePluginConf(path)
	require.NoError(t, err)

	// there should be an error after deleting the config from the cache
	err = plugin.ReadPluginConfig(path, &c)
	require.Error(t, err)
}
