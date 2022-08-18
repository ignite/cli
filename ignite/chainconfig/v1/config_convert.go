package v1

import (
	"github.com/ignite-hq/cli/ignite/chainconfig/common"
)

// ConvertNext implements the conversion of the current config to the next version.
func (c *Config) ConvertNext() (common.Config, error) {
	// v1 is the latest version, there is no need to convert.
	return c, nil
}
