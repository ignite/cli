package v1

import (
	"fmt"

	"github.com/imdario/mergo"

	"github.com/ignite/cli/ignite/chainconfig/common"
)

// Config is the user given configuration to do additional setup
// during serve.
type Config struct {
	Validators        []Validator `yaml:"validators"`
	common.BaseConfig `yaml:",inline"`
}

// GetHost returns the Host.
func (c *Config) GetHost() common.Host {
	if len(c.Validators) == 0 {
		return common.Host{}
	}

	validator := c.Validators[0]

	host := common.Host{}
	rpc := host.RPC
	p2p := host.P2P
	prof := host.Prof
	grpc := host.GRPC
	grpcweb := host.GRPCWeb
	api := host.API
	if validator.Config != nil {
		if val, ok := validator.Config["rpc"]; ok {
			v, yes := val.(map[string]interface{})
			if yes {
				if address, key := v["laddr"]; key {
					rpc = fmt.Sprintf("%v", address)
				}
			}
		}

		if val, ok := validator.Config["p2p"]; ok {
			v, yes := val.(map[string]interface{})
			if yes {
				if address, key := v["laddr"]; key {
					p2p = fmt.Sprintf("%v", address)
				}
			}
		}

		if val, ok := validator.Config["pprof_laddr"]; ok {
			prof = fmt.Sprintf("%v", val)
		}
	}

	if validator.App != nil {
		if val, ok := validator.App["grpc"]; ok {
			v, yes := val.(map[string]interface{})
			if yes {
				if address, key := v["address"]; key {
					grpc = fmt.Sprintf("%v", address)
				}
			}
		}

		if val, ok := validator.App["grpc-web"]; ok {
			v, yes := val.(map[string]interface{})
			if yes {
				if address, key := v["address"]; key {
					grpcweb = fmt.Sprintf("%v", address)
				}
			}
		}

		if val, ok := validator.App["api"]; ok {
			v, yes := val.(map[string]interface{})
			if yes {
				if address, key := v["address"]; key {
					api = fmt.Sprintf("%v", address)
				}
			}
		}
	}
	// Get the information from the first validator.
	return common.Host{
		RPC:     rpc,
		P2P:     p2p,
		Prof:    prof,
		GRPC:    grpc,
		GRPCWeb: grpcweb,
		API:     api,
	}
}

// GetInit returns the Init.
func (c *Config) GetInit() common.Init {
	if len(c.Validators) == 0 {
		return common.Init{}
	}
	validator := c.Validators[0]

	// Get the information from the first validator.
	return common.Init{
		App:            validator.App,
		Client:         validator.Client,
		Config:         validator.Config,
		Home:           validator.Home,
		KeyringBackend: validator.KeyringBackend,
	}
}

// ListAccounts returns the list of all the accounts.
func (c *Config) ListAccounts() []common.Account {
	return c.Accounts
}

// ListValidators returns the list of all the validators.
func (c *Config) ListValidators() []common.Validator {
	validators := make([]common.Validator, len(c.Validators))
	for i := range c.Validators {
		validators[i] = &c.Validators[i]
	}

	return validators
}

// Clone returns an identical copy of the instance
func (c *Config) Clone() common.Config {
	copy := *c
	return &copy
}

// FillValidatorsDefaults fills in the defaults values for the validators if they are missing.
func (c *Config) FillValidatorsDefaults(defaultValidator Validator) error {
	for i := range c.Validators {
		if err := c.Validators[i].FillDefaults(defaultValidator); err != nil {
			return err
		}
	}
	return nil
}

// Default returns the instance with the default value
func (c *Config) Default() common.Config {
	return &Config{
		Validators: []Validator{
			{
				App: map[string]interface{}{"grpc": map[string]interface{}{"address": "0.0.0.0:9090"},
					"grpc-web": map[string]interface{}{"address": "0.0.0.0:9091"}, "api": map[string]interface{}{"address": "0.0.0.0:1317"}},
				Config: map[string]interface{}{"rpc": map[string]interface{}{"laddr": "0.0.0.0:26657"},
					"p2p": map[string]interface{}{"laddr": "0.0.0.0:26656"}, "pprof_laddr": "0.0.0.0:6060"},
			},
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
	Bonded string `yaml:"bonded"`

	// App overwrites appd's config/app.toml configs.
	App map[string]interface{} `yaml:"app"`

	// Config overwrites appd's config/config.toml configs.
	Config map[string]interface{} `yaml:"config"`

	// Client overwrites appd's config/client.toml configs.
	Client map[string]interface{} `yaml:"client"`

	// Home overwrites default home directory used for the app
	Home string `yaml:"home"`

	// KeyringBackend is the default keyring backend to use for blockchain initialization
	KeyringBackend string `yaml:"keyring-backend"`

	// Gentx overwrites appd's config/gentx.toml configs.
	Gentx map[string]interface{} `yaml:"gentx"`
}

// GetName returns the name of the validator.
func (v *Validator) GetName() string {
	return v.Name
}

// GetBonded returns the bonded value.
func (v *Validator) GetBonded() string {
	return v.Bonded
}

// FillDefaults fills in the default values in the parameter defaultValidator.
func (v *Validator) FillDefaults(defaultValidator Validator) error {
	if err := mergo.Merge(v, defaultValidator); err != nil {
		return err
	}
	return nil
}
