package chainconfig

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/imdario/mergo"

	"github.com/ignite-hq/cli/ignite/pkg/xfilepath"
	"github.com/ignite-hq/cli/ignite/pkg/xurl"
)

// internal type alias for convienence
type configmap = map[string]interface{}

var (
	// ConfigDirPath returns the path of configuration directory of Ignite.
	ConfigDirPath = xfilepath.JoinFromHome(xfilepath.Path(".ignite"))

	// ConfigFileNames is a list of recognized names as for Ignite's config file.
	ConfigFileNames = []string{"config.yml", "config.yaml"}
)

// ErrCouldntLocateConfig returned when config.yml cannot be found in the source code.
var ErrCouldntLocateConfig = errors.New(
	"could not locate a config.yml in your chain. please follow the link for" +
		"how-to: https://github.com/ignite/cli/blob/develop/docs/configure/index.md")

// DefaultConf holds default configuration.
var DefaultConf = Config{
	Validators: []Validator{
		{
			// when in Docker on MacOS, it only works with 0.0.0.0.
			App: configmap{
				"grpc": configmap{
					"address": "0.0.0.0:9090",
				},

				"grpc-web": configmap{
					"address": "0.0.0.0:9091",
				},

				"api": configmap{
					"address": "0.0.0.0:1317",
				},
			},

			Config: configmap{
				"rpc": configmap{
					"laddr": "0.0.0.0:26657",
				},

				"p2p": configmap{
					"laddr": "0.0.0.0:26656",
				},

				"pprof_laddr": "0.0.0.0:6060",
			},
		},
	},
	Build: Build{
		Proto: Proto{
			Path: "proto",
			ThirdPartyPaths: []string{
				"third_party/proto",
				"proto_vendor",
			},
		},
	},
	Faucet: Faucet{
		Host: "0.0.0.0:45001",
	},
}

// Config is the user given configuration to do additional setup
// during serve.
type Config struct {
	Accounts   []Account              `yaml:"accounts"`
	Validators []Validator            `yaml:"validators"`
	Faucet     Faucet                 `yaml:"faucet"`
	Build      Build                  `yaml:"build"`
	Genesis    map[string]interface{} `yaml:"genesis"`
	Client     Client                 `yaml:"client"`
}

// AccountByName finds account by name.
func (c Config) AccountByName(name string) (acc Account, found bool) {
	for _, acc := range c.Accounts {
		if acc.Name == name {
			return acc, true
		}
	}
	return Account{}, false
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

// Validator holds info related to validator settings.
type Validator struct {
	Name string `yaml:"name"`

	Bonded string `yaml:"bonded"`

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

	// GenTx overwrites the default gentx transaction info.
	GenTx GenTx `yaml:"gentx"`
}

func (v Validator) API() string {
	str, _ := xurl.TCP(v.getConfigItemWithNamespace(v.App, "api", "address"))
	return str
}

func (v Validator) GRPC() string {
	return v.getConfigItemWithNamespace(v.App, "grpc", "address")
}

func (v Validator) GRPCWeb() string {
	return v.getConfigItemWithNamespace(v.App, "grpc-web", "address")
}

func (v Validator) RPC() string {
	str, _ := xurl.TCP(v.getConfigItemWithNamespace(v.Config, "rpc", "laddr"))
	return str
}

func (v Validator) P2P() string {
	str, _ := xurl.TCP(v.getConfigItemWithNamespace(v.Config, "p2p", "laddr"))
	return str
}

func (v Validator) Pprof() string {
	return v.getConfigItem(v.Config, "pprof_laddr")
}

// we can generalize these getter into a single recursive map field getter
// func that accepts an abritrary depth, but we only ever seem to either
// do 1 or two deep, so that might be overkill.
func (v Validator) getConfigItemWithNamespace(conf map[string]interface{}, namespace, key string) string {
	// sanity checks since its map[string]interface{} types all the way down
	if confMap, ok := conf[namespace].(configmap); ok {
		if rawval, ok := confMap[key]; ok {
			if val, ok := rawval.(string); ok {
				return val
			}
			return "" // default?
		}
		return "" // default?
	}
	return "" // default?
}

func (v Validator) getConfigItem(conf map[string]interface{}, key string) string {
	// sanity checks since its map[string]interface{} types all the way down
	if rawval, ok := conf[key]; ok {
		if val, ok := rawval.(string); ok {
			return val
		}
		return "" // default?
	}
	return "" // default?
}

// GenTx validator info
type GenTx struct {
	// Amount to self stake. Overwrites Validator.Bonded
	Amount string `yaml:"amount"`

	// Moniker of valudator. Overwrites Validator.Name
	Moniker string `yaml:"moniker"`

	CommisionRate           string `yaml:"rate"`
	CommissionMaxRate       string `yaml:"max-rate"`
	CommissionMaxChangeRate string `yaml:"max-change-rate"`
	MinDelegation           string `yaml:"min-delegation"`
	GasPrices               string `yaml:"gas-prices"`
	Details                 string `yaml:"details"`
	Identity                string `yaml:"identity"`
	Website                 string `yaml:"website"`
	SecurityContact         string `yaml:"securty-contact"`
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

// // Init overwrites sdk configurations with given values.
// type Init struct {
// 	// App overwrites appd's config/app.toml configs.
// 	App map[string]interface{} `yaml:"app"`

// 	// Client overwrites appd's config/client.toml configs.
// 	Client map[string]interface{} `yaml:"client"`

// 	// Config overwrites appd's config/config.toml configs.
// 	Config map[string]interface{} `yaml:"config"`

// 	// Home overwrites default home directory used for the app
// 	Home string `yaml:"home"`

// 	// KeyringBackend is the default keyring backend to use for blockchain initialization
// 	KeyringBackend string `yaml:"keyring-backend"`
// }

// Host keeps configuration related to started servers.
// Keeping for backwards compatability for now, despite no longer being used
// directly on the config object.
type Host struct {
	RPC     string `yaml:"rpc"`
	P2P     string `yaml:"p2p"`
	Prof    string `yaml:"prof"`
	GRPC    string `yaml:"grpc"`
	GRPCWeb string `yaml:"grpc-web"`
	API     string `yaml:"api"`
}

// note(jsimnz): Using github.com/imdario/mergo to merge the config struct with the default struct
// has a problem with the change to Validators being a slice instead of a struct. There is a mergo
// option `mergo.WithSliceDeepCopy` which will correctly handle merging the default validator slice
// with the defined one in the config, but it has a side-effect of setting `config.Overwrite` true
// on the `mergo.Config` struct. So we need to add an additional functional option call to set
// `config.Overwrite` back to false, otherwise other parts of the config don't get merged correctly
// like Faucet.
func mergoWithoutOverwrite(config *mergo.Config) {
	config.Overwrite = false
}

// Parse parses config.yml into UserConfig.
func Parse(r io.Reader) (Config, error) {
	var conf Config
	if err := yaml.NewDecoder(r).Decode(&conf); err != nil {
		return conf, err
	}
	if err := mergo.Merge(&conf, DefaultConf); err != nil {
		return Config{}, err
	}
	conf, err := mergeValidators(conf)
	if err != nil {
		return Config{}, err
	}
	return conf, validate(conf)
}

func mergeValidators(conf Config) (Config, error) {
	//only first validator for now
	if len(conf.Validators) == 0 {
		return Config{}, &ValidationError{"at least one validator is required"}
	}
	val := conf.Validators[0]
	if err := mergo.Merge(&val, DefaultConf.Validators[0]); err != nil {
		return Config{}, err
	}
	conf.Validators[0] = val
	return conf, nil
}

// ParseFile parses config.yml from the path.
func ParseFile(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, nil
	}
	defer file.Close()
	return Parse(file)
}

// validate validates user config.
func validate(conf Config) error {
	if len(conf.Accounts) == 0 {
		return &ValidationError{"at least 1 account is needed"}
	}
	if len(conf.Validators) == 0 || conf.Validators[0].Name == "" {
		return &ValidationError{"at least one validator is required"}
	}
	return nil
}

// ValidationError is returned when a configuration is invalid.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("config is not valid: %s", e.Message)
}

// LocateDefault locates the default path for the config file, if no file found returns ErrCouldntLocateConfig.
func LocateDefault(root string) (path string, err error) {
	for _, name := range ConfigFileNames {
		path = filepath.Join(root, name)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}
	}
	return "", ErrCouldntLocateConfig
}

// FaucetHost returns the faucet host to use
func FaucetHost(conf Config) string {
	// We keep supporting Port option for backward compatibility
	// TODO: drop this option in the future
	host := conf.Faucet.Host
	if conf.Faucet.Port != 0 {
		host = fmt.Sprintf(":%d", conf.Faucet.Port)
	}

	return host
}

// CreateConfigDir creates config directory if it is not created yet.
func CreateConfigDir() error {
	confPath, err := ConfigDirPath()
	if err != nil {
		return err
	}

	return os.MkdirAll(confPath, 0o755)
}
