package common

import "fmt"

type Version int

// Account holds the options related to setting up Cosmos wallets.
type Account struct {
	Name     string   `yaml:"name"`
	Coins    []string `yaml:"coins,omitempty"`
	Mnemonic string   `yaml:"mnemonic,omitempty"`
	Address  string   `yaml:"address,omitempty"`
	CoinType string   `yaml:"cointype,omitempty"`

	// The RPCAddress off the chain that account is issued at.
	RPCAddress string `yaml:"rpc_address,omitempty"`
}

// Build holds build configs.
type Build struct {
	Main    string   `yaml:"main"`
	Binary  string   `yaml:"binary"`
	LDFlags []string `yaml:"ldflags"`
	Proto   Proto    `yaml:"proto"`
}

// Proto holds proto build configs.
type Proto struct {
	// Path is the relative path of where app's proto files are located at.
	Path string `yaml:"path"`

	// ThirdPartyPath is the relative path of where the third party proto files are
	// located that used by the app.
	ThirdPartyPaths []string `yaml:"third_party_paths"`
}

// Client configures code generation for clients.
type Client struct {
	// Vuex configures code generation for Vuex.
	Vuex Vuex `yaml:"vuex"`

	// Dart configures client code generation for Dart.
	Dart Dart `yaml:"dart"`

	// OpenAPI configures OpenAPI spec generation for API.
	OpenAPI OpenAPI `yaml:"openapi"`
}

// Vuex configures code generation for Vuex.
type Vuex struct {
	// Path configures out location for generated Vuex code.
	Path string `yaml:"path"`
}

// Dart configures client code generation for Dart.
type Dart struct {
	// Path configures out location for generated Dart code.
	Path string `yaml:"path"`
}

// OpenAPI configures OpenAPI spec generation for API.
type OpenAPI struct {
	Path string `yaml:"path"`
}

// Faucet configuration.
type Faucet struct {
	// Name is faucet account's name.
	Name *string `yaml:"name"`

	// Coins holds type of coin denoms and amounts to distribute.
	Coins []string `yaml:"coins"`

	// CoinsMax holds of chain denoms and their max amounts that can be transferred
	// to single user.
	CoinsMax []string `yaml:"coins_max"`

	// LimitRefreshTime sets the timeframe at the end of which the limit will be refreshed
	RateLimitWindow string `yaml:"rate_limit_window"`

	// Host is the host of the faucet server
	Host string `yaml:"host"`

	// Port number for faucet server to listen at.
	Port int `yaml:"port"`
}

// Init overwrites sdk configurations with given values.
type Init struct {
	// App overwrites appd's config/app.toml configs.
	App map[string]interface{} `yaml:"app"`

	// Client overwrites appd's config/client.toml configs.
	Client map[string]interface{} `yaml:"client"`

	// Config overwrites appd's config/config.toml configs.
	Config map[string]interface{} `yaml:"config"`

	// Home overwrites default home directory used for the app
	Home string `yaml:"home"`

	// KeyringBackend is the default keyring backend to use for blockchain initialization
	KeyringBackend string `yaml:"keyring-backend"`
}

// Host keeps configuration related to started servers.
type Host struct {
	RPC     string `yaml:"rpc"`
	P2P     string `yaml:"p2p"`
	Prof    string `yaml:"prof"`
	GRPC    string `yaml:"grpc"`
	GRPCWeb string `yaml:"grpc-web"`
	API     string `yaml:"api"`
}

// BaseConfig is the struct containing all the common fields for the config across all the versions.
type BaseConfig struct {
	Version  Version                `yaml:"version"`
	Build    Build                  `yaml:"build"`
	Accounts []Account              `yaml:"accounts"`
	Faucet   Faucet                 `yaml:"faucet"`
	Client   Client                 `yaml:"client"`
	Genesis  map[string]interface{} `yaml:"genesis"`
}

// GetVersion returns the version of the config.yaml file.
func (c *BaseConfig) GetVersion() Version {
	return c.Version
}

// GetGenesis returns the Genesis.
func (c *BaseConfig) GetGenesis() map[string]interface{} {
	return c.Genesis
}

// AccountByName finds account by name.
func (c *BaseConfig) AccountByName(name string) (acc Account, found bool) {
	for _, acc := range c.Accounts {
		if acc.Name == name {
			return acc, true
		}
	}
	return Account{}, false
}

// GetBuild returns the Build.
func (c *BaseConfig) GetBuild() Build {
	return c.Build
}

// GetFaucet returns the Faucet.
func (c *BaseConfig) GetFaucet() Faucet {
	return c.Faucet
}

// GetClient returns the Client.
func (c *BaseConfig) GetClient() Client {
	return c.Client
}

// FaucetHost returns the faucet host to use
func FaucetHost(conf Config) string {
	// We keep supporting Port option for backward compatibility
	// TODO: drop this option in the future
	host := conf.GetFaucet().Host
	if conf.GetFaucet().Port != 0 {
		host = fmt.Sprintf(":%d", conf.GetFaucet().Port)
	}

	return host
}
