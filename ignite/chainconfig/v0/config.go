package v0

import (
	"io"

	"gopkg.in/yaml.v2"

	"github.com/ignite-hq/cli/ignite/chainconfig/config"
)

// Config is the user given configuration to do additional setup during serve.
type Config struct {
	config.BaseConfig `yaml:",inline"`

	Validator Validator   `yaml:"validator"`
	Init      config.Init `yaml:"init"`
	Host      config.Host `yaml:"host"`
}

// Clone returns an identical copy of the instance.
func (c *Config) Clone() config.Converter {
	copy := *c
	return &copy
}

// Decode decodes the config file values from YAML.
func (c *Config) Decode(r io.Reader) error {
	if err := yaml.NewDecoder(r).Decode(c); err != nil {
		return err
	}

	return nil
}

// Validator holds info related to validator settings.
type Validator struct {
	Name   string `yaml:"name"`
	Staked string `yaml:"staked"`
}
