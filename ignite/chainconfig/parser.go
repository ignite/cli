package chainconfig

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
)

// Parse reads a config file.
// The `DefaultVersion` is assumed when there is no version field in the config file.
// When the version of the file beign read is not the latest it is automatically
// migrated to the latest version.
func Parse(configFile io.ReadSeeker) (*v1.Config, error) {
	version, err := ReadConfigVersion(configFile)
	if err != nil {
		return DefaultConfig(), err
	}

	if _, err := configFile.Seek(0, 0); err != nil {
		return DefaultConfig(), err
	}

	c, err := decodeConfig(configFile, version)
	if err != nil {
		return DefaultConfig(), err
	}

	if err = c.SetDefaults(); err != nil {
		return DefaultConfig(), err
	}

	cfg, err := ConvertLatest(c)
	if err != nil {
		return DefaultConfig(), err
	}

	return cfg, validateConfig(cfg)
}

// ParseFile parses a config from a file path.
func ParseFile(path string) (*v1.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return DefaultConfig(), err
	}

	defer file.Close()

	return Parse(file)
}

// ReadConfigVersion reads the config version.
func ReadConfigVersion(configFile io.Reader) (config.Version, error) {
	c := struct {
		Version config.Version `yaml:"version"`
	}{}

	err := yaml.NewDecoder(configFile).Decode(&c)

	return c.Version, err
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

// FaucetHost returns the faucet host to use.
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

func decodeConfig(r io.Reader, version config.Version) (config.Converter, error) {
	c, ok := Versions[version]
	if !ok {
		return nil, &UnsupportedVersionError{version}
	}

	cfg := c.Clone()
	if err := cfg.Decode(r); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateConfig(c *v1.Config) error {
	if len(c.Accounts) == 0 {
		return &ValidationError{"at least 1 account is needed"}
	}

	// TODO: Fix this validation to check that there are validators and all have a name
	for _, validator := range c.Validators {
		if validator.Name == "" {
			return &ValidationError{"validator is required"}
		}
	}

	return nil
}
