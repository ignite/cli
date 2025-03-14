package v1

import (
	"fmt"
	"io"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"

	"github.com/ignite/cli/v29/ignite/config/chain/base"
	"github.com/ignite/cli/v29/ignite/config/chain/defaults"
	"github.com/ignite/cli/v29/ignite/config/chain/version"
	"github.com/ignite/cli/v29/ignite/pkg/xnet"
)

// DefaultConfig returns a config with default values.
func DefaultConfig() *Config {
	c := Config{Config: base.DefaultConfig()}
	c.Version = 1
	return &c
}

// Config is the user given configuration to do additional setup during serve.
type Config struct {
	base.Config `yaml:",inline"`

	Validators []Validator `yaml:"validators" doc:"Contains information related to the list of validators and settings."`
}

func (c *Config) SetDefaults() error {
	if err := c.Config.SetDefaults(); err != nil {
		return err
	}
	return c.updateValidatorAddresses()
}

// Clone returns an identical copy of the instance.
func (c *Config) Clone() (version.Converter, error) {
	cfgCopy := Config{}
	return &cfgCopy, mergo.Merge(&cfgCopy, c, mergo.WithAppendSlice)
}

// Decode decodes the config file values from YAML.
func (c *Config) Decode(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(c)
}

func (c *Config) updateValidatorAddresses() (err error) {
	// Margin to increase port numbers of the default addresses
	margin := 10

	for i := range c.Validators {
		// Use default addresses for the first validator
		if i == 0 {
			continue
		}

		validator := &c.Validators[i]
		servers, err := validator.GetServers()
		if err != nil {
			return err
		}
		portIncrement := margin * i
		if portIncrement < 0 {
			return fmt.Errorf("calculated port increment is negative: %d", portIncrement) //nolint: forbidigo
		}

		servers, err = incrementDefaultServerPortsBy(servers, uint64(portIncrement))
		if err != nil {
			return err
		}

		if err := validator.SetServers(servers); err != nil {
			return err
		}
	}

	return nil
}

// Returns a new server where the default addresses have their ports
// incremented by a margin to avoid port clashing.
func incrementDefaultServerPortsBy(s Servers, inc uint64) (Servers, error) {
	var err error

	if s.GRPC.Address == defaults.GRPCAddress {
		s.GRPC.Address, err = xnet.IncreasePortBy(defaults.GRPCAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	if s.GRPCWeb.Address == defaults.GRPCWebAddress {
		s.GRPCWeb.Address, err = xnet.IncreasePortBy(defaults.GRPCWebAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	if s.API.Address == defaults.APIAddress {
		s.API.Address, err = xnet.IncreasePortBy(defaults.APIAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	if s.P2P.Address == defaults.P2PAddress {
		s.P2P.Address, err = xnet.IncreasePortBy(defaults.P2PAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	if s.RPC.Address == defaults.RPCAddress {
		s.RPC.Address, err = xnet.IncreasePortBy(defaults.RPCAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	if s.RPC.PProfAddress == defaults.PProfAddress {
		s.RPC.PProfAddress, err = xnet.IncreasePortBy(defaults.PProfAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	return s, nil
}
