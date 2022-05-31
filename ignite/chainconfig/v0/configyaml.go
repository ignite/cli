package v0

import (
	"github.com/ignite/cli/ignite/chainconfig/common"
)

// ConfigYaml is the user given configuration to do additional setup
// during serve.
type Config struct {
	Validator         Validator   `yaml:"validator"`
	Init              common.Init `yaml:"init"`
	Host              common.Host `yaml:"host"`
	common.BaseConfig `yaml:",inline"`
}

// ListAccounts returns the list of all the accounts.
func (c *Config) ListAccounts() []common.Account {
	return c.Accounts
}

// Clone returns an identical copy of the instance
func (c *Config) Clone() common.Config {
	copy := *c
	return &copy
}

// Validator holds info related to validator settings.
type Validator struct {
	Name   string `yaml:"name"`
	Staked string `yaml:"staked"`
}
