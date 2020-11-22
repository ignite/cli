package networkbuilder

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var (
	starportConfDir = os.ExpandEnv("$HOME/.starport")
	confPath        = filepath.Join(starportConfDir, "networkbuilder")
)

func init() {
	if err := os.MkdirAll(starportConfDir, 0666); err != nil {
		panic(err)
	}
}

// Config holds configuration about network builder's state.
type Config struct {
	// SPNAccount is the default spn account in use.
	SPNAccount string

	// FinalizedChains can be started without any further preperation.
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
	cf, err := os.Open(confPath)
	if err != nil {
		return &Config{}, nil
	}
	defer cf.Close()
	var c Config
	return &c, json.NewDecoder(cf).Decode(&c)
}

// ConfigSave saves the current state of Config.
func ConfigSave(c *Config) error {
	cf, err := os.OpenFile(confPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer cf.Close()
	return json.NewEncoder(cf).Encode(c)
}
