package v1

import (
	"io"

	"gopkg.in/yaml.v2"

	"github.com/ignite/cli/ignite/chainconfig/config"
	"github.com/ignite/cli/ignite/pkg/xnet"
)

// DefaultConfig returns a config with default values.
func DefaultConfig() *Config {
	c := Config{BaseConfig: config.DefaultBaseConfig()}
	c.Version = 1
	return &c
}

// Config is the user given configuration to do additional setup during serve.
type Config struct {
	config.BaseConfig `yaml:",inline"`

	Validators []Validator `yaml:"validators"`
}

func (c *Config) SetDefaults() error {
	if err := c.BaseConfig.SetDefaults(); err != nil {
		return err
	}

	// Make sure that validator addresses don't chash with each other
	if err := c.updateValidatorAddresses(); err != nil {
		return err
	}

	return nil
}

// Clone returns an identical copy of the instance
func (c *Config) Clone() config.Converter {
	copy := *c
	return &copy
}

// Decode decodes the config file values from YAML.
func (c *Config) Decode(r io.Reader) error {
	if err := yaml.NewDecoder(r).Decode(c); err != nil {
		return err
	}

	return nil
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

		servers, err = incrementDefaultServerPortsBy(servers, uint64(margin*i))
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

	if s.GRPC.Address == DefaultGRPCAddress {
		s.GRPC.Address, err = xnet.IncreasePortBy(DefaultGRPCAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	if s.GRPCWeb.Address == DefaultGRPCWebAddress {
		s.GRPCWeb.Address, err = xnet.IncreasePortBy(DefaultGRPCWebAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	if s.API.Address == DefaultAPIAddress {
		s.API.Address, err = xnet.IncreasePortBy(DefaultAPIAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	if s.P2P.Address == DefaultP2PAddress {
		s.P2P.Address, err = xnet.IncreasePortBy(DefaultP2PAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	if s.RPC.Address == DefaultRPCAddress {
		s.RPC.Address, err = xnet.IncreasePortBy(DefaultRPCAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	if s.RPC.PProfAddress == DefaultPProfAddress {
		s.RPC.PProfAddress, err = xnet.IncreasePortBy(DefaultPProfAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	return s, nil
}
