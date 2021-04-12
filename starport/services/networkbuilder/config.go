package networkbuilder

import (
	"encoding/json"
	"os"

	"github.com/tendermint/starport/starport/pkg/xfilepath"
	"github.com/tendermint/starport/starport/services"
)

var (
	confDir = "networkbuilder"

	// confPath returns the path of Starport Network configuration
	confPath = xfilepath.Join(
		xfilepath.PathWithError(services.StarportConfPath()),
		xfilepath.Path(confDir),
	)
)

// Config holds configuration about network builder's state.
type Config struct {
	// SPNAccount is the default spn account in use.
	SPNAccount string

	// FinalizedChains can be started without any further preparation.
	FinalizedChains []string
}

// IsChainMarkedFinalized checks if chain marked as finalized.
func (c *Config) IsChainMarkedFinalized(chainID string) bool {
	for _, c := range c.FinalizedChains {
		if c == chainID {
			return true
		}
	}
	return false
}

// MarkFinalized marks chain as finalized.
func (c *Config) MarkFinalized(chainID string) {
	c.FinalizedChains = append(c.FinalizedChains, chainID)
}

// ConfigGet retrieves the current state of Config.
func ConfigGet() (*Config, error) {
	conf, err := confPath()
	if err != nil {
		return &Config{}, nil
	}

	cf, err := os.Open(conf)
	if err != nil {
		return &Config{}, nil
	}
	defer cf.Close()
	var c Config
	return &c, json.NewDecoder(cf).Decode(&c)
}

// ConfigSave saves the current state of Config.
func ConfigSave(c *Config) error {
	conf, err := confPath()
	if err != nil {
		return nil
	}

	cf, err := os.OpenFile(conf, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer cf.Close()
	return json.NewEncoder(cf).Encode(c)
}
