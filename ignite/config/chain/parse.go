package chain

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"gopkg.in/yaml.v2"

	"github.com/ignite/cli/ignite/config/chain/version"
)

// Parse reads a config file.
// When the version of the file being read is not the latest
// it is automatically migrated to the latest version.
func Parse(configFile io.Reader) (*Config, error) {
	cfg, err := parse(configFile)
	if err != nil {
		return cfg, fmt.Errorf("error parsing config file: %w", err)
	}

	return cfg, validateConfig(cfg)
}

// ParseNetwork reads a config file for Ignite Network genesis.
// When the version of the file being read is not the latest
// it is automatically migrated to the latest version.
func ParseNetwork(configFile io.Reader) (*Config, error) {
	cfg, err := parse(configFile)
	if err != nil {
		return cfg, err
	}

	return cfg, validateNetworkConfig(cfg)
}

func parse(configFile io.Reader) (*Config, error) {
	var buf bytes.Buffer

	// Read the config file version first to know how to decode it
	version, err := ReadConfigVersion(io.TeeReader(configFile, &buf))
	if err != nil {
		return DefaultChainConfig(), err
	}

	// Decode the current config file version and assign default
	// values for the fields that are empty
	c, err := decodeConfig(&buf, version)
	if err != nil {
		return DefaultChainConfig(), err
	}

	// Make sure that the empty fields contain default values
	// after reading the config from the YAML file
	if err = c.SetDefaults(); err != nil {
		return DefaultChainConfig(), err
	}

	// Finally make sure the config is the latest one before validating it
	cfg, err := ConvertLatest(c)
	if err != nil {
		return DefaultChainConfig(), err
	}

	return cfg, nil
}

// ParseFile parses a config from a file path.
func ParseFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return DefaultChainConfig(), err
	}

	defer file.Close()

	return Parse(file)
}

// ParseNetworkFile parses a config for Ignite Network genesis from a file path.
func ParseNetworkFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return DefaultChainConfig(), err
	}

	defer file.Close()

	return ParseNetwork(file)
}

// ReadConfigVersion reads the config version.
func ReadConfigVersion(configFile io.Reader) (version.Version, error) {
	c := struct {
		Version version.Version `yaml:"version"`
	}{}

	err := yaml.NewDecoder(configFile).Decode(&c)

	return c.Version, err
}

func decodeConfig(r io.Reader, version version.Version) (version.Converter, error) {
	c, ok := Versions[version]
	if !ok {
		return nil, &UnsupportedVersionError{version}
	}

	cfg, err := c.Clone()
	if err != nil {
		return nil, err
	}

	if err = cfg.Decode(r); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateConfig(c *Config) error {
	if len(c.Accounts) == 0 {
		return &ValidationError{"at least one account is required"}
	}

	for _, validator := range c.Validators {
		if validator.Name == "" {
			return &ValidationError{"validator 'name' is required"}
		}

		if validator.Bonded == "" {
			return &ValidationError{"validator 'bonded' is required"}
		}
	}

	return nil
}

func validateNetworkConfig(c *Config) error {
	if len(c.Validators) != 0 {
		return &ValidationError{"no validators can be used in config for network genesis"}
	}

	for _, account := range c.Accounts {
		// must have valid bech32 addr
		if _, _, err := bech32.DecodeAndConvert(account.Address); err != nil {
			return fmt.Errorf("invalid address %s: %w", account.Address, err)
		}

		if account.Coins == nil {
			return &ValidationError{"account coins is required"}
		}

		if account.Mnemonic != "" {
			return &ValidationError{"cannot include mnemonic in network config genesis"}
		}
	}

	return nil
}
