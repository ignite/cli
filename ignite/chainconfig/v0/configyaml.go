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

// GetHost returns the Host.
func (c *Config) GetHost() common.Host {
	return c.Host
}

// GetGenesis returns the Genesis.
func (c *Config) GetGenesis() map[string]interface{} {
	return c.Genesis
}

// GetInit returns the Init.
func (c *Config) GetInit() common.Init {
	return c.Init
}

// ListAccounts returns the list of all the accounts.
func (c *Config) ListAccounts() []common.Account {
	return c.Accounts
}

// ListValidators returns the list of all the validators.
func (c *Config) ListValidators() []common.Validator {
	return []common.Validator{&c.Validator}
}

// Clone returns an identical copy of the instance
func (c *Config) Clone() common.Config {
	copy := *c
	return &copy
}

// Default returns the instance with the default value
func (c *Config) Default() common.Config {
	return &Config{
		Host: common.Host{
			// when in Docker on MacOS, it only works with 0.0.0.0.
			RPC:     "0.0.0.0:26657",
			P2P:     "0.0.0.0:26656",
			Prof:    "0.0.0.0:6060",
			GRPC:    "0.0.0.0:9090",
			GRPCWeb: "0.0.0.0:9091",
			API:     "0.0.0.0:1317",
		},
		BaseConfig: common.BaseConfig{
			Build: common.Build{
				Proto: common.Proto{
					Path: "proto",
					ThirdPartyPaths: []string{
						"third_party/proto",
						"proto_vendor",
					},
				},
			},
			Faucet: common.Faucet{
				Host: "0.0.0.0:4500",
			},
		},
	}
}

// Validator holds info related to validator settings.
type Validator struct {
	Name   string `yaml:"name"`
	Staked string `yaml:"staked"`
}

// GetName returns the name of the validator.
func (v *Validator) GetName() string {
	return v.Name
}

// GetBonded returns the bonded value.
func (v *Validator) GetBonded() string {
	return v.Staked
}
