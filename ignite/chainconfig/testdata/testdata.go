package testdata

import (
	_ "embed"

	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	v0testdata "github.com/ignite-hq/cli/ignite/chainconfig/v0/testdata"
	v1testdata "github.com/ignite-hq/cli/ignite/chainconfig/v1/testdata"
)

var Versions = map[config.Version][]byte{
	0: v0testdata.ConfigYAML,
	1: v1testdata.ConfigYAML,
}
