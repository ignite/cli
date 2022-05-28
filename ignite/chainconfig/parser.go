package chainconfig

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"

	"github.com/ignite-hq/cli/ignite/chainconfig/common"
	v0 "github.com/ignite-hq/cli/ignite/chainconfig/v0"
	"github.com/ignite-hq/cli/ignite/pkg/xfilepath"
)

var (
	// ConfigDirPath returns the path of configuration directory of Ignite.
	ConfigDirPath = xfilepath.JoinFromHome(xfilepath.Path(".ignite"))

	// ConfigFileNames is a list of recognized names as for Ignite's config file.
	ConfigFileNames = []string{"config.yml", "config.yaml"}
)

var (
	// ErrCouldntLocateConfig returned when config.yml cannot be found in the source code.
	ErrCouldntLocateConfig = errors.New(
		"could not locate a config.yml in your chain. please follow the link for" +
			"how-to: https://github.com/ignite-hq/cli/blob/develop/docs/configure/index.md")
)

// DefaultConf holds default configuration.
var DefaultConf = v0.ConfigYaml{
	Host: common.Host{
		// when in Docker on MacOS, it only works with 0.0.0.0.
		RPC:     "0.0.0.0:26657",
		P2P:     "0.0.0.0:26656",
		Prof:    "0.0.0.0:6060",
		GRPC:    "0.0.0.0:9090",
		GRPCWeb: "0.0.0.0:9091",
		API:     "0.0.0.0:1317",
	},
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
}

// Parse parses config.yml into UserConfig.
func Parse(r io.Reader) (v0.ConfigYaml, error) {
	var conf v0.ConfigYaml
	if err := yaml.NewDecoder(r).Decode(&conf); err != nil {
		return conf, err
	}
	if err := mergo.Merge(&conf, DefaultConf); err != nil {
		return v0.ConfigYaml{}, err
	}
	return conf, validate(conf)
}

// ParseFile parses config.yml from the path.
func ParseFile(path string) (v0.ConfigYaml, error) {
	file, err := os.Open(path)
	if err != nil {
		return v0.ConfigYaml{}, nil
	}
	defer file.Close()
	return Parse(file)
}

// validate validates user config.
func validate(conf v0.ConfigYaml) error {
	if len(conf.Accounts) == 0 {
		return &ValidationError{"at least 1 account is needed"}
	}
	if conf.Validator.Name == "" {
		return &ValidationError{"validator is required"}
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
func FaucetHost(conf v0.ConfigYaml) string {
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

	return os.MkdirAll(confPath, 0755)
}
