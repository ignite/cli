package base

import (
	"github.com/imdario/mergo"

	"github.com/ignite/cli/v28/ignite/config/chain/version"
	xyaml "github.com/ignite/cli/v28/ignite/pkg/yaml"
)

var (
	// DefaultGRPCAddress is the default GRPC address.
	DefaultGRPCAddress = "0.0.0.0:9090"

	// DefaultGRPCWebAddress is the default GRPC-Web address.
	DefaultGRPCWebAddress = "0.0.0.0:9091"

	// DefaultAPIAddress is the default API address.
	DefaultAPIAddress = "0.0.0.0:1317"

	// DefaultRPCAddress is the default RPC address.
	DefaultRPCAddress = "0.0.0.0:26657"

	// DefaultP2PAddress is the default P2P address.
	DefaultP2PAddress = "0.0.0.0:26656"

	// DefaultPProfAddress is the default Prof address.
	DefaultPProfAddress = "0.0.0.0:6060"
)

// Account holds the options related to setting up Cosmos wallets.
type Account struct {
	Name     string   `yaml:"name"`
	Coins    []string `yaml:"coins,omitempty"`
	Mnemonic string   `yaml:"mnemonic,omitempty"`
	Address  string   `yaml:"address,omitempty"`
	CoinType string   `yaml:"cointype,omitempty"`
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
}

// Client configures code generation for clients.
type Client struct {
	// TSClient configures code generation for Typescript Client.
	Typescript Typescript `yaml:"typescript,omitempty"`

	// Vuex configures code generation for Vuex stores.
	//
	// Deprecated: Will be removed eventually.
	Vuex Vuex `yaml:"vuex,omitempty"`

	// Composables configures code generation for Vue 3 composables.
	Composables Composables `yaml:"composables,omitempty"`

<<<<<<< HEAD
	// Hooks configures code generation for React hooks.
	Hooks Hooks `yaml:"hooks,omitempty"`

=======
>>>>>>> d1bf508a (refactor!: remove react frontend + re-enable disabled integration tests (#4744))
	// OpenAPI configures OpenAPI spec generation for API.
	OpenAPI OpenAPI `yaml:"openapi,omitempty"`
}

// TSClient configures code generation for Typescript Client.
type Typescript struct {
	// Path configures out location for generated Typescript Client code.
	Path string `yaml:"path"`
}

// Vuex configures code generation for Vuex stores.
//
// Deprecated: Will be removed eventually.
type Vuex struct {
	// Path configures out location for generated Vuex stores code.
	Path string `yaml:"path"`
}

// Composables configures code generation for vue-query hooks.
type Composables struct {
	// Path configures out location for generated vue-query hooks.
	Path string `yaml:"path"`
}

<<<<<<< HEAD
// Hooks configures code generation for react-query hooks.
type Hooks struct {
	// Path configures out location for generated vue-query hooks.
	Path string `yaml:"path"`
}

=======
>>>>>>> d1bf508a (refactor!: remove react frontend + re-enable disabled integration tests (#4744))
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
	Port uint `yaml:"port,omitempty"`

	// TxFee is the tx fee the faucet needs to pay for each transaction.
	TxFee string `yaml:"tx_fee,omitempty"`
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

// Validation describes the kind of validation the chain has.
type Validation string

const (
	// ValidationSovereign is when the chain has his own validator set.
	// Note that an empty string is also considered as a sovereign validation,
	// because this is the default value.
	ValidationSovereign = "sovereign"
	// ValidationConsumer is when the chain is validated by a provider chain.
	// Such chain is called a consumer chain.
	ValidationConsumer = "consumer"
)

// Config defines a struct with the fields that are common to all config versions.
type Config struct {
	Include    []string        `yaml:"include,omitempty"`
	Validation Validation      `yaml:"validation,omitempty"`
	Version    version.Version `yaml:"version"`
	Build      Build           `yaml:"build,omitempty"`
	Accounts   []Account       `yaml:"accounts"`
	Faucet     Faucet          `yaml:"faucet,omitempty"`
	Client     Client          `yaml:"client,omitempty"`
	Genesis    xyaml.Map       `yaml:"genesis,omitempty"`
	Minimal    bool            `yaml:"minimal,omitempty"`
}

// GetVersion returns the config version.
func (c Config) GetVersion() version.Version {
	return c.Version
}

// IsChainMinimal returns true if the chain is minimally scaffolded.
func (c Config) IsChainMinimal() bool {
	return c.Minimal
}

func (c Config) IsSovereignChain() bool {
	return c.Validation == "" || c.Validation == ValidationSovereign
}

func (c Config) IsConsumerChain() bool {
	return c.Validation == ValidationConsumer
}

// SetDefaults assigns default values to empty config fields.
func (c *Config) SetDefaults() error {
	return mergo.Merge(c, DefaultConfig())
}

// DefaultConfig returns a base config with default values.
func DefaultConfig() Config {
	return Config{
		Build: Build{
			Proto: Proto{
				Path: "proto",
			},
		},
		Faucet: Faucet{
			Host: "0.0.0.0:4500",
		},
	}
}
