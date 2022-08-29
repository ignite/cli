package testdata

import (
	_ "embed"

	"github.com/ignite-hq/cli/ignite/chainconfig/config"
)

//go:embed configv0.yaml
var ConfigV0YAML []byte

//go:embed configv1.yaml
var ConfigV1YAML []byte

var Versions = map[config.Version][]byte{
	0: ConfigV0YAML,
	1: ConfigV1YAML,
}
