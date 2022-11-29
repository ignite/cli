package testdata

import (
	pluginsconfig "github.com/ignite/cli/ignite/config/plugins"
	"testing"

	"github.com/ignite/cli/ignite/config"
	chainconfig "github.com/ignite/cli/ignite/config/chain"
	networkconfigTestdata "github.com/ignite/cli/ignite/config/chain/network/testdata"
	v0testdata "github.com/ignite/cli/ignite/config/chain/v0/testdata"
	v1testdata "github.com/ignite/cli/ignite/config/chain/v1/testdata"
	pluginsconfigTestdata "github.com/ignite/cli/ignite/config/plugins/testdata"
)

var Versions = map[chainconfig.Version][]byte{
	0: v0testdata.ConfigYAML,
	1: v1testdata.ConfigYAML,
}

var NetworkConfig = networkconfigTestdata.ConfigYAML

var PluginsConfig = pluginsconfigTestdata.ConfigYAML

func GetLatestConfig(t *testing.T) *config.ChainConfig {
	return v1testdata.GetConfig(t)
}

func GetLatestNetworkConfig(t *testing.T) *config.ChainConfig {
	return networkconfigTestdata.GetConfig(t)
}

func GetPluginsConfig(t *testing.T) *pluginsconfig.Config {
	return pluginsconfigTestdata.GetConfig(t)
}
