package config

import (
	"io"

	"github.com/imdario/mergo"

	xyaml "github.com/ignite/cli/ignite/pkg/yaml"
)

// Version defines the type for the config version number.
type Version uint

// Converter defines the interface required to migrate configurations to newer versions.
type Converter interface {
	// Clone clones the config by returning a new copy of the current one.
	Clone() Converter

	// SetDefaults assigns default values to empty config fields.
	SetDefaults() error

	// GetVersion returns the config version.
	GetVersion() Version

	// ConvertNext converts the config to the next version.
	ConvertNext() (Converter, error)

	// Decode decodes the config file from YAML and updates it's values.
	Decode(io.Reader) error
}

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
	Main    string   `yaml:"main,omitempty"`
	Binary  string   `yaml:"binary,omitempty"`
	LDFlags []string `yaml:"ldflags,omitempty"`
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
	Vuex Vuex `yaml:"vuex,omitempty"`

	// Dart configures client code generation for Dart.
	Dart Dart `yaml:"dart,omitempty"`

	// OpenAPI configures OpenAPI spec generation for API.
	OpenAPI OpenAPI `yaml:"openapi,omitempty"`
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

	// CoinsMax holds of chain denoms and their max amounts that can be transferred to single user.
	CoinsMax []string `yaml:"coins_max,omitempty"`

	// LimitRefreshTime sets the timeframe at the end of which the limit will be refreshed
	RateLimitWindow string `yaml:"rate_limit_window,omitempty"`

	// Host is the host of the faucet server
	Host string `yaml:"host,omitempty"`

	// Port number for faucet server to listen at.
	Port int `yaml:"port,omitempty"`
}

// Init overwrites sdk configurations with given values.
type Init struct {
	// App overwrites appd's config/app.toml configs.
	App xyaml.Map `yaml:"app"`

	// Client overwrites appd's config/client.toml configs.
	Client xyaml.Map `yaml:"client"`

	// Config overwrites appd's config/config.toml configs.
	Config xyaml.Map `yaml:"config"`

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

// BaseConfig defines a struct with the fields that are common to all config versions.
type BaseConfig struct {
	Version  Version                `yaml:"version"`
	Build    Build                  `yaml:"build"`
	Accounts []Account              `yaml:"accounts"`
	Faucet   Faucet                 `yaml:"faucet,omitempty"`
	Client   Client                 `yaml:"client,omitempty"`
	Genesis  map[string]interface{} `yaml:"genesis,omitempty"`
}

// GetVersion returns the config version.
func (c BaseConfig) GetVersion() Version {
	return c.Version
}

// SetDefaults assigns default values to empty config fields.
func (c *BaseConfig) SetDefaults() error {
	if err := mergo.Merge(c, DefaultBaseConfig()); err != nil {
		return err
	}

	return nil
}

// DefaultBaseConfig returns a base config with default values.
func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		Build: Build{
			Proto: Proto{
				Path:            "proto",
				ThirdPartyPaths: []string{"third_party/proto", "proto_vendor"},
			},
		},
		Faucet: Faucet{
			Host: "0.0.0.0:4500",
		},
	}
}
