package base

import (
	"github.com/imdario/mergo"

	"github.com/ignite/cli/v29/ignite/config/chain/defaults"
	"github.com/ignite/cli/v29/ignite/config/chain/version"
	"github.com/ignite/cli/v29/ignite/pkg/xyaml"
)

// Account holds the options related to setting up Cosmos wallets.
type Account struct {
	Name          string   `yaml:"name" doc:"local name of a key pair associated with an account"`
	Coins         []string `yaml:"coins,omitempty" doc:"list of token balances for the account"`
	Mnemonic      string   `yaml:"mnemonic,omitempty" doc:"account mnemonic"`
	Address       string   `yaml:"address,omitempty" doc:"account address"`
	CoinType      string   `yaml:"cointype,omitempty" doc:"coin type number for HD derivation (default 118)"`
	AccountNumber string   `yaml:"account_number,omitempty" doc:"account number for HD derivation (less than equal 2147483647)"`
	AddressIndex  string   `yaml:"address_index,omitempty" doc:"address index number for HD derivation (less than equal 2147483647)"`
}

// Build holds build configs.
type Build struct {
	Main    string   `yaml:"main,omitempty" doc:"build main path"`
	Binary  string   `yaml:"binary,omitempty" doc:"binary path"`
	LDFlags []string `yaml:"ldflags,omitempty" doc:"custom build ld flags"`
	Proto   Proto    `yaml:"proto" doc:"holds proto build configs"`
}

// Proto holds proto build configs.
type Proto struct {
	// Path is the relative path of where app's proto files are located at.
	Path string `yaml:"path" doc:"relative path of where app's proto files are located at"`
}

// Client configures code generation for clients.
type Client struct {
	// TSClient configures code generation for Typescript Client.
	Typescript Typescript `yaml:"typescript,omitempty" doc:"configures code generation for Typescript Client"`

	// Composables configures code generation for Vue 3 composables.
	Composables Composables `yaml:"composables,omitempty" doc:"configures code generation for Vue 3 composables"`

	// Hooks configures code generation for React hooks.
	Hooks Hooks `yaml:"hooks,omitempty" doc:"configures code generation for React hooks"`

	// OpenAPI configures OpenAPI spec generation for API.
	OpenAPI OpenAPI `yaml:"openapi,omitempty" doc:"configures OpenAPI spec generation for API"`
}

// Typescript configures code generation for Typescript Client.
type Typescript struct {
	// Path configures out location for generated Typescript Client code.
	Path string `yaml:"path" doc:"relative path of where app's typescript files are located at"`
}

// Vuex configures code generation for Vuex stores.
//
// Deprecated: Will be removed eventually.
type Vuex struct {
	// Path configures out location for generated Vuex stores code.
	Path string `yaml:"path" doc:"relative path of where app's vuex files are located at"`
}

// Composables configures code generation for vue-query hooks.
type Composables struct {
	// Path configures out location for generated vue-query hooks.
	Path string `yaml:"path" doc:"relative path of where app's composable files are located at"`
}

// Hooks configures code generation for react-query hooks.
type Hooks struct {
	// Path configures out location for generated vue-query hooks.
	Path string `yaml:"path" doc:"relative path of where app's hooks are located at"`
}

// OpenAPI configures OpenAPI spec generation for API.
type OpenAPI struct {
	Path string `yaml:"path" doc:"relative path of where app's openapi files are located at"`
}

// Faucet configuration.
type Faucet struct {
	// Name is faucet account's name.
	Name *string `yaml:"name" doc:"faucet account's name"`

	// Coins holds type of coin denoms and amounts to distribute.
	Coins []string `yaml:"coins" doc:"holds type of coin denoms and amounts to distribute"`

	// CoinsMax holds of chain denoms and their max amounts that can be transferred to single user.
	CoinsMax []string `yaml:"coins_max,omitempty" doc:"holds of chain denoms and their max amounts that can be transferred to single user"`

	// LimitRefreshTime sets the timeframe at the end of which the limit will be refreshed.
	RateLimitWindow string `yaml:"rate_limit_window,omitempty" doc:"sets the timeframe at the end of which the limit will be refreshed"`

	// Host is the host of the faucet server.
	Host string `yaml:"host,omitempty" doc:"the host of the faucet server"`

	// Port number for faucet server to listen at.
	Port uint `yaml:"port,omitempty" doc:"number for faucet server to listen at"`
}

// Init overwrites sdk configurations with given values.
type Init struct {
	// App overwrites appd's config/app.toml configs.
	App xyaml.Map `yaml:"app" doc:"overwrites appd's config/app.toml configs"`

	// Client overwrites appd's config/client.toml configs.
	Client xyaml.Map `yaml:"client" doc:"overwrites appd's config/client.toml configs"`

	// Config overwrites appd's config/config.toml configs.
	Config xyaml.Map `yaml:"config" doc:"overwrites appd's config/config.toml configs"`

	// Home overwrites default home directory used for the app.
	Home string `yaml:"home" doc:"overwrites default home directory used for the app"`
}

// Host keeps configuration related to started servers.
type Host struct {
	RPC     string `yaml:"rpc" doc:"RPC server address"`
	P2P     string `yaml:"p2p" doc:"P2P server address"`
	Prof    string `yaml:"prof" doc:"Profiling server address"`
	GRPC    string `yaml:"grpc" doc:"GRPC server address"`
	GRPCWeb string `yaml:"grpc-web" doc:"GRPC Web server address"`
	API     string `yaml:"api" doc:"API server address"`
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
	Validation Validation      `yaml:"validation,omitempty" doc:"describes the kind of validation the chain has"`
	Version    version.Version `yaml:"version" doc:"defines the type for the config version number"`
	Build      Build           `yaml:"build,omitempty" doc:"holds build configs"`
	Accounts   []Account       `yaml:"accounts" doc:"holds the options related to setting up Cosmos wallets"`
	Faucet     Faucet          `yaml:"faucet,omitempty" doc:"faucet configuration"`
	Client     Client          `yaml:"client,omitempty" doc:"configures code generation for clients"`
	Genesis    xyaml.Map       `yaml:"genesis,omitempty" doc:"custom genesis modifications"`
	Minimal    bool            `yaml:"minimal,omitempty" doc:"minimal blockchain with the minimum required Cosmos SDK modules"`
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
				Path: defaults.ProtoDir,
			},
		},
		Faucet: Faucet{
			Host: defaults.FaucetHost,
		},
	}
}
