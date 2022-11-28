package testdata

import (
	"testing"

	"github.com/ignite/cli/ignite/config"
	"github.com/ignite/cli/ignite/config/base"
	networkconfigTestData "github.com/ignite/cli/ignite/config/networkconfig/testdata"
	v0testdata "github.com/ignite/cli/ignite/config/v0/testdata"
	v1testdata "github.com/ignite/cli/ignite/config/v1/testdata"
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
