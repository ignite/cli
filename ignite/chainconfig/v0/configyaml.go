package v0

import "github.com/ignite-hq/cli/ignite/chainconfig/common"

// ConfigYaml is the user given configuration to do additional setup
// during serve.
type ConfigYaml struct {
	Validator             common.Validator       `yaml:"validator"`
	Client                common.Client          `yaml:"client"`
	Init                  common.Init            `yaml:"init"`
	Genesis               map[string]interface{} `yaml:"genesis"`
	Host                  common.Host            `yaml:"host"`
	common.BaseConfigYaml `yaml:",inline"`
}

// GetHost returns the Host.
func (c *ConfigYaml) GetHost() common.Host {
	return c.Host
}

// GetGenesis returns the Genesis.
func (c *ConfigYaml) GetGenesis() map[string]interface{} {
	return c.Genesis
}

// GetInit returns the Init.
func (c *ConfigYaml) GetInit() common.Init {
	return c.Init
}

// GetClient returns the Client.
func (c *ConfigYaml) GetClient() common.Client {
	return c.Client
}

// ListAccounts returns the list of all the accounts.
func (c *ConfigYaml) ListAccounts() []common.Account {
	return c.Accounts
}

// ListValidators returns the list of all the validators.
func (c *ConfigYaml) ListValidators() []common.Validator {
	return []common.Validator{c.Validator}
}

// Clone returns an identical copy of the instance
func (c *ConfigYaml) Clone() common.Config {
	copy := *c
	return &copy
}

// Default returns the instance with the default value
func (c *ConfigYaml) Default() common.Config {
	return &ConfigYaml{
		Host: common.Host{
			// when in Docker on MacOS, it only works with 0.0.0.0.
			RPC:     "0.0.0.0:26657",
			P2P:     "0.0.0.0:26656",
			Prof:    "0.0.0.0:6060",
			GRPC:    "0.0.0.0:9090",
			GRPCWeb: "0.0.0.0:9091",
			API:     "0.0.0.0:1317",
		},
		BaseConfigYaml: common.BaseConfigYaml{
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
