package chain

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"

	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/ignite/cli/v29/ignite/config/chain/defaults"
	"github.com/ignite/cli/v29/ignite/config/chain/version"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// Parse reads a config file.
// When the version of the file being read is not the latest
// it is automatically migrated to the latest version.
func Parse(configFile io.Reader) (*Config, error) {
	cfg, err := parse(configFile)
	if err != nil {
		return cfg, errors.Errorf("error parsing config file: %w", err)
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
	v, err := ReadConfigVersion(io.TeeReader(configFile, &buf))
	if err != nil {
		return DefaultChainConfig(), err
	}

	// Decode the current config file version and assign default
	// values for the fields that are empty
	c, err := decodeConfig(&buf, v)
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

	// Handle includes
	if err := handleIncludes(cfg); err != nil {
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

// ReadProtoPath reads the proto path.
func ReadProtoPath(configFile io.Reader) (string, error) {
	c := struct {
		Build struct {
			Proto struct {
				Path string `yaml:"path"`
			} `yaml:"proto"`
		} `yaml:"build"`
	}{}

	c.Build.Proto.Path = defaults.ProtoDir
	err := yaml.NewDecoder(configFile).Decode(&c)

	return c.Build.Proto.Path, err
}

// decodeConfig decodes a config from an io.Reader using the specified version.
// It returns a version.Converter interface or an error if version is not supported
// or if decoding fails.
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

// validateConfig validates a chain configuration by checking that at least one
// account exists and that all validators have required name and bonded fields.
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

// validateNetworkConfig validates a network genesis configuration by ensuring
// no validators exist and that all accounts have valid addresses, coins and no mnemonics.
func validateNetworkConfig(c *Config) error {
	if len(c.Validators) != 0 {
		return &ValidationError{"no validators can be used in config for network genesis"}
	}

	for _, account := range c.Accounts {
		// must have valid bech32 addr.
		if _, _, err := bech32.DecodeAndConvert(account.Address); err != nil {
			return errors.Errorf("invalid address %s: %w", account.Address, err)
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

// handleIncludes processes included configuration files referenced in the main config.
// It supports both local files and remote URLs, merging their contents with the main config.
func handleIncludes(cfg *Config) error {
	if len(cfg.Include) == 0 {
		return nil
	}

	for _, includePath := range cfg.Include {
		if u, err := url.ParseRequestURI(includePath); err == nil && u.Scheme != "" {
			includePath, err = fetchConfigFile(includePath)
			if err != nil {
				return errors.Wrapf(err, "failed to fetch included config file '%s'", includePath)
			}
			defer os.Remove(includePath)
		}

		// Resolve path - if relative, use the base directory.
		absPath, err := filepath.Abs(includePath)
		if err != nil {
			return errors.Wrapf(err, "failed to resolve included path '%s'", includePath)
		}

		includeFile, err := os.Open(absPath)
		if err != nil {
			return errors.Errorf("failed to open included file '%s'", includePath)
		}
		defer includeFile.Close()

		// Parse the included config.
		includeCfg, err := parse(includeFile)
		if err != nil {
			return errors.Wrapf(err, "failed to parse included config file '%s'", includePath)
		}

		if cfg.Version != includeCfg.Version {
			return errors.Errorf("included config version '%d' does not match with chain config version '%d'", includeCfg.Version, cfg.Version)
		}

		// Merge the included config with the primary config.
		if err = mergo.Merge(cfg, includeCfg, mergo.WithAppendSlice, mergo.WithOverride); err != nil {
			return errors.Wrapf(err, "failed to merge included file '%s'", includePath)
		}
	}

	return nil
}

// fetchConfigFile downloads a configuration file from a URL and saves it to a temporary file.
// Returns the path to the temporary file or an error if the download fails.
func fetchConfigFile(url string) (string, error) {
	// Download the file from URL to a temporary file.
	tmpFile, err := os.CreateTemp("", "config-*.yml")
	if err != nil {
		return "", errors.Wrapf(err, "failed to create temp file for URL")
	}
	defer tmpFile.Close()

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "failed to download from URL '%s'", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("failed to download file, status code: %d", resp.StatusCode)
	}

	if _, err = io.Copy(tmpFile, resp.Body); err != nil {
		return "", errors.Wrapf(err, "failed to save downloaded file from '%s'", url)
	}

	if _, err = tmpFile.Seek(0, io.SeekStart); err != nil {
		return "", errors.Wrapf(err, "failed to rewind temp file from '%s'", url)
	}

	return tmpFile.Name(), nil
}
