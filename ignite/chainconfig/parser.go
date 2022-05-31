package chainconfig

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"

	"github.com/ignite/cli/ignite/chainconfig/common"
	v0 "github.com/ignite/cli/ignite/chainconfig/v0"
	v1 "github.com/ignite/cli/ignite/chainconfig/v1"
	"github.com/ignite/cli/ignite/pkg/xfilepath"
)

var (
	// ErrCouldntLocateConfig returned when config.yml cannot be found in the source code.
	ErrCouldntLocateConfig = errors.New(
		"could not locate a config.yml in your chain. please follow the link for" +
			"how-to: https://github.com/ignite/cli/blob/develop/docs/configure/index.md")

	// ConfigDirPath returns the path of configuration directory of Ignite.
	ConfigDirPath = xfilepath.JoinFromHome(xfilepath.Path(".ignite"))

	// ConfigFileNames is a list of recognized names as for Ignite's config file.
	ConfigFileNames = []string{"config.yml", "config.yaml"}

	DefaultVersion = "v0"

	// Migration defines the version as the key and the config instance as the value
	Migration = map[common.Version]common.Config{0: &v0.Config{}, 1: &v1.Config{}}

	// DefaultConfig defines the default config without the validators.
	DefaultConfig = &v1.Config{
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

// Parse parses config.yml into UserConfig based on the version.
func Parse(r io.Reader) (*v1.Config, error) {
	// Read the version field
	version, err := getConfigVersion(r)
	if err != nil {
		return nil, err
	}

	conf, err := GetConfigInstance(version)
	if err != nil {
		return nil, err
	}

	// Go back to the beginning of the file.
	_, err = r.(io.Seeker).Seek(0, 0)
	if err != nil {
		return nil, err
	}

	// Decode the file by parsing the content again.
	if err = yaml.NewDecoder(r).Decode(conf); err != nil {
		return nil, err
	}

	conf, err = ConvertLatest(conf)
	if err != nil {
		return nil, err
	}

	if err = mergo.Merge(conf, DefaultConfig); err != nil {
		return nil, err
	}

	latestConfig := conf.(*v1.Config)
	// As the lib does not support the merge of the array, we fill in the default values for the list of validators.
	if err = latestConfig.FillValidatorsDefaults(v1.DefaultValidator); err != nil {
		return nil, err
	}

	return latestConfig, validate(latestConfig)
}

// IsConfigLatest checks if the version of the config file is the latest
func IsConfigLatest(path string) (common.Version, bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, false, err
	}
	defer file.Close()
	version, err := getConfigVersion(file)
	if err != nil {
		return 0, false, err
	}
	return version, version == common.LatestVersion, nil
}

// MigrateConfigFile upgrades the config file to the latest version.
func MigrateConfigFile(configFile string) error {
	configyml, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer configyml.Close()
	conf, err := Parse(configyml)
	if err != nil {
		return err
	}

	err = configyml.Truncate(0)
	if err != nil {
		return err
	}

	_, err = configyml.Seek(0, 0)
	if err != nil {
		return err
	}
	return yaml.NewEncoder(configyml).Encode(conf)
}

// getConfigVersion returns the version in the io.Reader based on the field version.
func getConfigVersion(r io.Reader) (common.Version, error) {
	var baseConf common.BaseConfig
	if err := yaml.NewDecoder(r).Decode(&baseConf); err != nil {
		return 0, err
	}
	return baseConf.ConfigVersion, nil
}

// GetConfigInstance retrieves correct config instance based on the version.
func GetConfigInstance(version common.Version) (common.Config, error) {
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
func ParseFile(path string) (*v1.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Parse(file)
}

// validate validates user config.
func validate(conf *v1.Config) error {
	if len(conf.ListAccounts()) == 0 {
		return &ValidationError{"at least 1 account is needed"}
	}

	for _, validator := range conf.Validators {
		if validator.Name == "" {
			return &ValidationError{"validator is required"}
		}
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

// UnknownInputError is returned when the input of Parse is unknown.
type UnknownInputError struct {
	Message string
}

func (e *UnknownInputError) Error() string {
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
func FaucetHost(conf *v1.Config) string {
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
