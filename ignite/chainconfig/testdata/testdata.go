package testdata

import (
	"testing"

	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	v0testdata "github.com/ignite-hq/cli/ignite/chainconfig/v0/testdata"
	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
	v1testdata "github.com/ignite-hq/cli/ignite/chainconfig/v1/testdata"
)

var Versions = map[config.Version][]byte{
	0: v0testdata.ConfigYAML,
	1: v1testdata.ConfigYAML,
}

func GetLatestConfig(t *testing.T) *v1.Config {
	return v1testdata.GetConfig(t)
}
