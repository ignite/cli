package starportconf

import (
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
)

var (
	// FileNames holds a list of appropriate names for the config file.
	FileNames = []string{"config.yml", "config.yaml"}
)

// Config is the user given configuration to do additional setup
// during serve.
type Config struct {
	Accounts  []Account              `yaml:"accounts"`
	Validator Validator              `yaml:"validator"`
	Genesis   map[string]interface{} `yaml:"genesis"`
}

// Account holds the options related to setting up Cosmos wallets.
type Account struct {
	Name  string   `yaml:"name"`
	Coins []string `yaml:"coins"`
}

// Validator holds info related to validator settings.
type Validator struct {
	Name   string `yaml:"name"`
	Staked string `yaml:"staked"`
}

// Parse parses config.yml into UserConfig.
func Parse(r io.Reader) (Config, error) {
	var conf Config
	if err := yaml.NewDecoder(r).Decode(&conf); err != nil {
		return conf, err
	}
	return conf, validate(conf)
}

// validate validates user config.
func validate(conf Config) error {
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
