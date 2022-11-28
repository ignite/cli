package testdata

import (
	"github.com/ignite/cli/ignite/config/chain/base"
	v0testdata "github.com/ignite/cli/ignite/config/chain/v0/testdata"
	v1testdata "github.com/ignite/cli/ignite/config/chain/v1/testdata"
	"testing"

	"github.com/ignite/cli/ignite/config"
	networkconfigTestData "github.com/ignite/cli/ignite/config/networkconfig/testdata"
)

var Versions = map[base.Version][]byte{
	0: v0testdata.ConfigYAML,
	1: v1testdata.ConfigYAML,
}

var NetworkConfig = networkconfigTestData.ConfigYAML

func GetLatestConfig(t *testing.T) *config.Config {
	return v1testdata.GetConfig(t)
}

func GetLatestNetworkConfig(t *testing.T) *config.Config {
	return networkconfigTestData.GetConfig(t)
}
