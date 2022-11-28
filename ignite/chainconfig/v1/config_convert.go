package v1

// ConvertNext implements the conversion of the current config to the next version.
func (c *Config) ConvertNext() (base.Converter, error) {
	// v1 is the latest version, there is no need to convert.
	return c, nil
}
