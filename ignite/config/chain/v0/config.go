package v0

import (
	"io"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"

	"github.com/ignite/cli/v29/ignite/config/chain/base"
	"github.com/ignite/cli/v29/ignite/config/chain/version"
)

// Config is the user given configuration to do additional setup during serve.
type Config struct {
	base.Config `yaml:",inline"`

	Validator Validator `yaml:"validator" doc:"Contains information related to the validator and settings."`
	Init      base.Init `yaml:"init" doc:"Overwrites the appd's config/config.toml configurations."`
	Host      base.Host `yaml:"host" doc:"Keeps configuration related to started servers."`
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
	Name   string `yaml:"name" doc:"Name of the validator."`
	Staked string `yaml:"staked" doc:"Amount staked by the validator."`
}
