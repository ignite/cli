package chainconfig

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"

	"github.com/ignite/cli/ignite/chainconfig/common"
	v0 "github.com/ignite/cli/ignite/chainconfig/v0"
	"github.com/ignite/cli/ignite/pkg/xfilepath"
)

var (
	// ConfigDirPath returns the path of configuration directory of Ignite.
	ConfigDirPath = xfilepath.JoinFromHome(xfilepath.Path(".ignite"))

	// ConfigFileNames is a list of recognized names as for Ignite's config file.
	ConfigFileNames = []string{"config.yml", "config.yaml"}

	DefaultVersion = "v0"

	// Migration defines the version as the key and the config instance as the value
	Migration = map[string]common.Config{"v0": &v0.ConfigYaml{}}
)

var (
	// ErrCouldntLocateConfig returned when config.yml cannot be found in the source code.
	ErrCouldntLocateConfig = errors.New(
		"could not locate a config.yml in your chain. please follow the link for" +
			"how-to: https://github.com/ignite-hq/cli/blob/develop/docs/configure/index.md")
)

// Parse parses config.yml into UserConfig based on the version.
func Parse(content []byte) (common.Config, error) {
	// Read the version field
	version, err := getConfigVersion(bytes.NewReader(content))
	if err != nil {
		return GetDefaultConfig(), err
	}

	conf := GetConfigInstance(version)
	if err = yaml.NewDecoder(bytes.NewReader(content)).Decode(conf); err != nil {
		return conf, err
	}
	if err = mergo.Merge(conf, conf.Default()); err != nil {
		return GetDefaultConfig(), err
	}
	return conf, validate(conf)
}

// getConfigVersion returns the version in the io.Reader based on the field version.
func getConfigVersion(r io.Reader) (string, error) {
	var baseConf common.BaseConfigYaml
	if err := yaml.NewDecoder(r).Decode(&baseConf); err != nil {
		return DefaultVersion, err
	}
	if baseConf.Version != "" {
		return baseConf.Version, nil
	}
	return DefaultVersion, nil
}

// GetConfigInstance retrieves correct config instance based on the version.
func GetConfigInstance(version string) common.Config {
	var config common.Config
	var ok bool
	if config, ok = Migration[version]; !ok {
		// If there is no matching instance, return the config with the v0 version.
		return GetDefaultConfig()
	}
	// If we find the matching instance, clone the instance and return it.
	return config.Clone()
}

// ParseFile parses config.yml from the path.
func ParseFile(path string) (common.Config, error) {
	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		return GetDefaultConfig(), err
	}

	return Parse(yfile)
}

// GetDefaultConfig returns the default instance of the config.
func GetDefaultConfig() common.Config {
	return &v0.ConfigYaml{}
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
