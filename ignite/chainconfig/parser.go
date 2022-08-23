package chainconfig

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
	"github.com/imdario/mergo"
)

// Parse parses config.yml into UserConfig based on the version.
// TODO parse to the given config
func Parse(configFile io.Reader, out config.Converter) error {
	// Read the version field
	version, err := readConfigVersion(configFile)
	if err != nil {
		return err
	}

	conf, err := GetConfigInstance(version)
	if err != nil {
		return err
	}

	// Go back to the beginning of the file.
	_, err = configFile.(io.Seeker).Seek(0, 0)
	if err != nil {
		return err
	}

	// Decode the file by parsing the content again.
	if err = yaml.NewDecoder(configFile).Decode(conf); err != nil {
		return err
	}

	conf, err = ConvertLatest(conf)
	if err != nil {
		return err
	}

	if err = mergo.Merge(conf, DefaultConfig); err != nil {
		return err
	}

	// As the lib does not support the merge of the array, we fill in the default values for the list of validators.
	latestConfig := conf.(*v1.Config)
	if err = latestConfig.FillValidatorsDefaults(v1.DefaultValidator); err != nil {
		return err
	}

	// return latestConfig, validate(latestConfig)
	return nil
}

// IsConfigLatest checks if the version of the config file is the latest.
func IsConfigLatest(cfgPath string) (config.Version, bool, error) {
	file, err := os.Open(cfgPath)
	if err != nil {
		return 0, false, err
	}

	defer file.Close()

	ver, err := readConfigVersion(file)
	if err != nil {
		return 0, false, err
	}

	return ver, ver == LatestVersion, nil
}

// MigrateLatest migrates a config file to the latest version.
func MigrateLatest(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer file.Close()

	cfg, err := Parse(file)
	if err != nil {
		return err
	}

	if err = file.Truncate(0); err != nil {
		return err
	}

	if _, err = file.Seek(0, 0); err != nil {
		return err
	}

	return yaml.NewEncoder(file).Encode(cfg)
}

// GetConfigInstance retrieves correct config instance based on the version.
func GetConfigInstance(version config.Version) (config.Converter, error) {
	cfg, ok := Migration[version]
	if !ok {
		// If there is no matching instance, return the config with the v0 version.
		return nil, &UnsupportedVersionError{"the version is not available in the supported list"}
	}

	// If we find the matching instance, clone the instance and return it.
	return cfg.Clone(), nil
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

func validate(cfg *v1.Config) error {
	if len(cfg.ListAccounts()) == 0 {
		return &ValidationError{"at least 1 account is needed"}
	}

	for _, validator := range cfg.Validators {
		if validator.Name == "" {
			return &ValidationError{"validator is required"}
		}
	}

	return nil
}

// LocateDefault locates the default path for the config file.
// Returns ErrConfigNotFound when no config file found.
func LocateDefault(root string) (path string, err error) {
	for _, name := range ConfigFileNames {
		path = filepath.Join(root, name)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}
	}

	return "", ErrConfigNotFound
}

// FaucetHost returns the faucet host to use
func FaucetHost(cfg *v1.Config) string {
	// We keep supporting Port option for backward compatibility
	// TODO: drop this option in the future
	host := cfg.Faucet.Host
	if cfg.Faucet.Port != 0 {
		host = fmt.Sprintf(":%d", cfg.Faucet.Port)
	}

	return host
}

// CreateConfigDir creates config directory if it is not created yet.
func CreateConfigDir() error {
	path, err := ConfigDirPath()
	if err != nil {
		return err
	}

	return os.MkdirAll(path, 0755)
}

func readConfigVersion(r io.Reader) (config.Version, error) {
	var cfg config.BaseConfig
	if err := yaml.NewDecoder(r).Decode(&cfg); err != nil {
		return 0, err
	}

	return cfg.Version, nil
}
