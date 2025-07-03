package base

import (
	"dario.cat/mergo"

	"github.com/ignite/cli/v29/ignite/config/chain/defaults"
	"github.com/ignite/cli/v29/ignite/config/chain/version"
	"github.com/ignite/cli/v29/ignite/pkg/xyaml"
)

// Account holds the options related to setting up Cosmos wallets.
type Account struct {
	Name          string   `yaml:"name" doc:"Local name associated with the Account's key pair."`
	Coins         []string `yaml:"coins,omitempty" doc:"List of token balances for the account."`
	Mnemonic      string   `yaml:"mnemonic,omitempty" doc:"Mnemonic phrase for the account."`
	Address       string   `yaml:"address,omitempty" doc:"Address of the account."`
	CoinType      string   `yaml:"cointype,omitempty" doc:"Coin type number for HD derivation (default is 118)."`
	AccountNumber string   `yaml:"account_number,omitempty" doc:"Account number for HD derivation (must be ≤ 2147483647)."`
	AddressIndex  string   `yaml:"address_index,omitempty" doc:"Address index number for HD derivation (must be ≤ 2147483647)."`
}

// Build holds build configs.
type Build struct {
	Main    string   `yaml:"main,omitempty" doc:"Path to the main build file."`
	Binary  string   `yaml:"binary,omitempty" doc:"Path to the binary file."`
	LDFlags []string `yaml:"ldflags,omitempty" doc:"List of custom linker flags for building the binary."`
	Proto   Proto    `yaml:"proto" doc:"Contains proto build configuration options."`
}

// Proto holds proto build configs.
type Proto struct {
	// Path is the relative path of where app's proto files are located at.
	Path string `yaml:"path" doc:"Relative path where the application's proto files are located."`
}

// Client configures code generation for clients.
type Client struct {
	// TSClient configures code generation for Typescript Client.
	Typescript Typescript `yaml:"typescript,omitempty" doc:"Relative path where the application's Typescript files are located."`

	// Composables configures code generation for Vue 3 composables.
	Composables Composables `yaml:"composables,omitempty" doc:"Configures Vue 3 composables code generation."`

	// OpenAPI configures OpenAPI spec generation for API.
	OpenAPI OpenAPI `yaml:"openapi,omitempty" doc:"Configures OpenAPI spec generation for the API."`
}

// Typescript configures code generation for Typescript Client.
type Typescript struct {
	// Path configures out location for generated Typescript Client code.
	Path string `yaml:"path" doc:"Relative path where the application's Typescript files are located."`
}

// Composables configures code generation for vue-query hooks.
type Composables struct {
	// Path configures out location for generated vue-query hooks.
	Path string `yaml:"path" doc:"Relative path where the application's composable files are located."`
}

// OpenAPI configures OpenAPI spec generation for API.
type OpenAPI struct {
	Path string `yaml:"path" doc:"Relative path where the application's OpenAPI files are located."`
}

// Faucet configuration.
type Faucet struct {
	// Name is faucet account's name.
	Name *string `yaml:"name" doc:"Name of the faucet account."`

	// Coins holds type of coin denoms and amounts to distribute.
	Coins []string `yaml:"coins" doc:"Types and amounts of coins the faucet distributes."`

	// CoinsMax holds of chain denoms and their max amounts that can be transferred to single user.
	CoinsMax []string `yaml:"coins_max,omitempty" doc:"Maximum amounts of coins that can be transferred to a single user."`

	// LimitRefreshTime sets the timeframe at the end of which the limit will be refreshed.
	RateLimitWindow string `yaml:"rate_limit_window,omitempty" doc:"Timeframe after which the limit will be refreshed."`

	// Host is the host of the faucet server.
	Host string `yaml:"host,omitempty" doc:"Host address of the faucet server."`

	// Port number for faucet server to listen at.
	Port uint `yaml:"port,omitempty" doc:"Port number for the faucet server."`

	// TxFee is the tx fee the faucet needs to pay for each transaction.
	TxFee string `yaml:"tx_fee,omitempty" doc:"Tx fee the faucet needs to pay for each transaction."`
}

// Init overwrites sdk configurations with given values.
// Deprecated: Used in config v0 only.
type Init struct {
	// App overwrites appd's config/app.toml configs.
	App xyaml.Map `yaml:"app" doc:"Overwrites the appd's config/app.toml configurations."`

	// Client overwrites appd's config/client.toml configs.
	Client xyaml.Map `yaml:"client" doc:"Overwrites the appd's config/client.toml configurations."`

	// Config overwrites appd's config/config.toml configs.
	Config xyaml.Map `yaml:"config" doc:"Overwrites the appd's config/config.toml configurations."`

	// Home overwrites default home directory used for the app.
	Home string `yaml:"home" doc:"Overwrites the default home directory used for the application."`
}

// Host keeps configuration related to started servers.
// Deprecated: Used in config v0 only.
type Host struct {
	RPC     string `yaml:"rpc" doc:"RPC server address."`
	P2P     string `yaml:"p2p" doc:"P2P server address."`
	Prof    string `yaml:"prof" doc:"Profiling server address."`
	GRPC    string `yaml:"grpc" doc:"GRPC server address."`
	GRPCWeb string `yaml:"grpc-web" doc:"GRPC Web server address."`
	API     string `yaml:"api" doc:"API server address."`
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
	// This is a special case for ICS chains, used by the consumer ignite app (https://github.com/ignite/apps/issues/101).
	ValidationConsumer = "consumer"
)

// Config defines a struct with the fields that are common to all config versions.
type Config struct {
	Include      []string        `yaml:"include,omitempty" doc:"Include incorporate a separate config.yml file directly in your current config file."`
	Validation   Validation      `yaml:"validation,omitempty" doc:"Specifies the type of validation the blockchain uses (e.g., sovereign)."`
	Version      version.Version `yaml:"version" doc:"Defines the configuration version number."`
	Build        Build           `yaml:"build,omitempty" doc:"Contains build configuration options."`
	Accounts     []Account       `yaml:"accounts" doc:"Lists the options for setting up Cosmos Accounts."`
	Faucet       Faucet          `yaml:"faucet,omitempty" doc:"Configuration for the faucet."`
	Client       Client          `yaml:"client,omitempty" doc:"Configures client code generation."`
	Genesis      xyaml.Map       `yaml:"genesis,omitempty" doc:"Custom genesis block modifications. Follow the nesting of the genesis file here to access all the parameters."`
	DefaultDenom string          `yaml:"default_denom,omitempty" doc:"Default staking denom (default is stake)."`
}

// GetVersion returns the config version.
func (c Config) GetVersion() version.Version {
	return c.Version
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
				Path: defaults.ProtoDir,
			},
		},
		Faucet: Faucet{
			Host: defaults.FaucetHost,
		},
	}
}
