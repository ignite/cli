package common

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

// Validator holds info related to validator settings.
type Validator struct {
	Name   string `yaml:"name"`
	Staked string `yaml:"staked"`
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
