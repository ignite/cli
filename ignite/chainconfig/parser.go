package chainconfig

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"

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

	DefaultVersion = "v0"

	// Migration defines the version as the key and the config instance as the value
	Migration = map[int]common.Config{0: &v0.Config{}, 1: &v1.Config{}}

	DefaultConfig0 = v0.Config{
		Host: common.Host{
			// when in Docker on MacOS, it only works with 0.0.0.0.
			RPC:     "0.0.0.0:26657",
			P2P:     "0.0.0.0:26656",
			Prof:    "0.0.0.0:6060",
			GRPC:    "0.0.0.0:9090",
			GRPCWeb: "0.0.0.0:9091",
			API:     "0.0.0.0:1317",
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
)

var (
	// ErrCouldntLocateConfig returned when config.yml cannot be found in the source code.
	ErrCouldntLocateConfig = errors.New(
		"could not locate a config.yml in your chain. please follow the link for" +
			"how-to: https://github.com/ignite-hq/cli/blob/develop/docs/configure/index.md")
)

// Parse parses config.yml into UserConfig based on the version.
func Parse(r io.Reader) (common.Config, error) {
	// The io.Reader can only be read once, so we need to keep the content for further usage.
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	// Read the version field
	version, err := getConfigVersion(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	fmt.Println(string(content))
	fmt.Println(version)
	conf, err := GetConfigInstance(version)
	if err != nil {
		return nil, err
	}
	if err = yaml.NewDecoder(bytes.NewReader(content)).Decode(conf); err != nil {
		return conf, err
	}
	if err = mergo.Merge(conf, DefaultConfig0); err != nil {
		return nil, err
	}

	return conf, validate(conf)
}

// getConfigVersion returns the version in the io.Reader based on the field version.
func getConfigVersion(r io.Reader) (int, error) {
	var baseConf common.BaseConfig
	if err := yaml.NewDecoder(r).Decode(&baseConf); err != nil {
		return 0, err
	}
	return baseConf.Version, nil
}

// GetConfigInstance retrieves correct config instance based on the version.
func GetConfigInstance(version int) (common.Config, error) {
	var config common.Config
	var ok bool
	if config, ok = Migration[version]; !ok {
		// If there is no matching instance, return the config with the v0 version.
		return nil, &UnsupportedVersionError{"the version is not available in the supported list"}
	}
	// If we find the matching instance, clone the instance and return it.
	return config.Clone(), nil
}

// ParseFile parses config.yml from the path.
func ParseFile(path string) (common.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Parse(file)
}

// validate validates user config.
func validate(conf common.Config) error {
	if len(conf.ListAccounts()) == 0 {
		return &ValidationError{"at least 1 account is needed"}
	}
	if conf.ListValidators()[0].GetName() == "" {
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

// UnsupportedVersionError is returned when the version of the config is not supported.
type UnsupportedVersionError struct {
	Message string
}

func (e *UnsupportedVersionError) Error() string {
	return fmt.Sprintf("the version of the config is unsupported: %s", e.Message)
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
func FaucetHost(conf common.Config) string {
	// We keep supporting Port option for backward compatibility
	// TODO: drop this option in the future
	host := conf.GetFaucet().Host
	if conf.GetFaucet().Port != 0 {
		host = fmt.Sprintf(":%d", conf.GetFaucet().Port)
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
