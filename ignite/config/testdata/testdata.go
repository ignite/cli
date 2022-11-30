package testdata

import (
	"testing"

	"github.com/ignite/cli/ignite/config"
	chainconfig "github.com/ignite/cli/ignite/config/chain"
	networkconfigTestdata "github.com/ignite/cli/ignite/config/chain/network/testdata"
	v0testdata "github.com/ignite/cli/ignite/config/chain/v0/testdata"
	v1testdata "github.com/ignite/cli/ignite/config/chain/v1/testdata"
)

var Versions = map[chainconfig.Version][]byte{
	0: v0testdata.ConfigYAML,
	1: v1testdata.ConfigYAML,
}

var NetworkConfig = networkconfigTestdata.ConfigYAML

func GetLatestConfig(t *testing.T) *config.ChainConfig {
	return v1testdata.GetConfig(t)
}

func GetLatestNetworkConfig(t *testing.T) *config.ChainConfig {
	return networkconfigTestdata.GetConfig(t)
}
