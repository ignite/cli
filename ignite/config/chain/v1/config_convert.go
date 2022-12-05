package v1

import (
	"github.com/ignite/cli/ignite/config/chain/version"
)

// ConvertNext implements the conversion of the current config to the next version.
func (c *Config) ConvertNext() (version.Converter, error) {
	// v1 is the latest version, there is no need to convert.
	return c, nil
}
