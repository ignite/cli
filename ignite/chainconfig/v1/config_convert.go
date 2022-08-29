package v1

import "github.com/ignite/cli/ignite/chainconfig/config"

// ConvertNext implements the conversion of the current config to the next version.
func (c *Config) ConvertNext() (config.Converter, error) {
	// v1 is the latest version, there is no need to convert.
	return c, nil
}
