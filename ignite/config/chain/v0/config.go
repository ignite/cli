package v0

import (
	"io"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"

	"github.com/ignite/cli/v29/ignite/config/chain/base"
	"github.com/ignite/cli/v29/ignite/config/chain/version"
)

// Config is the user given configuration to do additional setup during serve.
type Config struct {
	base.Config `yaml:",inline"`

	Validator Validator `yaml:"validator" doc:"holds info related to validator settings"`
	Init      base.Init `yaml:"init" doc:"overwrites sdk configurations with given values"`
	Host      base.Host `yaml:"host" doc:"keeps configuration related to started servers"`
}

// Clone returns an identical copy of the instance.
func (c *Config) Clone() (version.Converter, error) {
	cfgCopy := Config{}
	return &cfgCopy, mergo.Merge(&cfgCopy, c, mergo.WithAppendSlice)
}

// Decode decodes the config file values from YAML.
func (c *Config) Decode(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(c)
}

// Validator holds info related to validator settings.
type Validator struct {
	Name   string `yaml:"name" doc:"name of the validator"`
	Staked string `yaml:"staked" doc:"how much the validator has staked"`
}
