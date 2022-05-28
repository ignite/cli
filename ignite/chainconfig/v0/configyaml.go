package v0

import (
	"github.com/ignite-hq/cli/ignite/chainconfig/common"
	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
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

// ConvertNext implements the conversion of the current config to the next version.
func (c *Config) ConvertNext() (common.Config, error) {
	targetConfig := &v1.Config{}

	// All the fields in the base config remain the same
	targetConfig.BaseConfig = c.BaseConfig

	// Change the version to 1
	targetConfig.Version = 1

	// There is only one validator in v0. Set it as the only one validator in v1.
	validators := make([]v1.Validator, 0, 1)
	validator := v1.Validator{}

	sourceValidator := c.Validator
	validator.Name = sourceValidator.Name
	validator.Bonded = sourceValidator.Staked

	// Set the fields in Init to the target validator
	validator.Home = c.Init.Home
	validator.KeyringBackend = c.Init.KeyringBackend
	validator.Client = c.Init.Client

	validator.App = c.Init.App
	validator.Config = c.Init.Config

	// If the fields in Host is not empty, we need to merge them into validator.app and validator.config.
	if c.Host.P2P != "" {
		if validator.Config == nil {
			// Create an empty map for Config
			validator.Config = make(map[string]interface{})
		}
		p2p := map[string]interface{}{"laddr": c.Host.P2P}
		validator.Config["p2p"] = p2p
	}

	if c.Host.Prof != "" {
		if validator.Config == nil {
			// Create an empty map for Config
			validator.Config = make(map[string]interface{})
		}
		validator.Config["pprof_laddr"] = c.Host.Prof
	}

	if c.Host.RPC != "" {
		if validator.Config == nil {
			// Create an empty map for Config
			validator.Config = make(map[string]interface{})
		}
		rpc := map[string]interface{}{"laddr": c.Host.RPC}
		validator.Config["rpc"] = rpc
	}

	if c.Host.GRPCWeb != "" {
		if validator.App == nil {
			// Create an empty map for App
			validator.App = make(map[string]interface{})
		}
		grpcweb := map[string]interface{}{"address": c.Host.GRPCWeb}
		validator.App["grpc-web"] = grpcweb
	}

	if c.Host.GRPC != "" {
		if validator.App == nil {
			// Create an empty map for App
			validator.App = make(map[string]interface{})
		}
		grpc := map[string]interface{}{"address": c.Host.GRPC}
		validator.App["grpc"] = grpc
	}

	if c.Host.API != "" {
		if validator.App == nil {
			// Create an empty map for App
			validator.App = make(map[string]interface{})
		}
		api := map[string]interface{}{"address": c.Host.API}
		validator.App["api"] = api
	}

	validators = append(validators, validator)
	targetConfig.Validators = validators

	return targetConfig, nil
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
