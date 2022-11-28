package v0

import (
	"io"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"

	chainconfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/config/chain/base"
)

// Config is the user given configuration to do additional setup during serve.
type Config struct {
	base.Config `yaml:",inline"`

	Validator Validator `yaml:"validator"`
	Init      base.Init `yaml:"init"`
	Host      base.Host `yaml:"host"`
}

// Clone returns an identical copy of the instance.
func (c *Config) Clone() (chainconfig.Converter, error) {
	copy := Config{}
	if err := mergo.Merge(&copy, c, mergo.WithAppendSlice); err != nil {
		return nil, err
	}

	return &copy, nil
}

// Decode decodes the config file values from YAML.
func (c *Config) Decode(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(c)
}

// Validator holds info related to validator settings.
type Validator struct {
	Name   string `yaml:"name"`
	Staked string `yaml:"staked"`
}
