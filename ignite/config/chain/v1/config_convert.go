package v1

import chainconfig "github.com/ignite/cli/ignite/config/chain"

// ConvertNext implements the conversion of the current config to the next version.
func (c *Config) ConvertNext() (chainconfig.Converter, error) {
	// v1 is the latest version, there is no need to convert.
	return c, nil
}
