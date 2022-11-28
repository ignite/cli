package v0

import (
	"github.com/ignite/cli/ignite/chainconfig/base"
	"io"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

// Config is the user given configuration to do additional setup during serve.
type Config struct {
	base.BaseConfig `yaml:",inline"`

	Validator Validator `yaml:"validator"`
	Init      base.Init `yaml:"init"`
	Host      base.Host `yaml:"host"`
}

// Clone returns an identical copy of the instance.
func (c *Config) Clone() (base.Converter, error) {
	copy := Config{}
	if err := mergo.Merge(&copy, c, mergo.WithAppendSlice); err != nil {
		return nil, err
	}

	return &copy, nil
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
