package plugins

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// ParseFile parses a plugins config.
func ParseFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return DefaultConfig(), err
	}

	defer file.Close()

	return Parse(file)
}

// Parse reads a config file for ignite binary plugins
func Parse(configFile io.Reader) (*Config, error) {
	return parse(configFile)
}

func parse(configFile io.Reader) (*Config, error) {
	var c Config

	err := yaml.NewDecoder(configFile).Decode(&c)
	// if the error is end of file we have an empty configuration, return the default config
	if err == io.EOF {
		return DefaultConfig(), nil
	}

	return &c, err
}
