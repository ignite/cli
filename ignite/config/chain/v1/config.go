package v1

import (
	"io"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"

	"github.com/ignite/cli/ignite/config/chain/base"
	"github.com/ignite/cli/ignite/config/chain/version"
	"github.com/ignite/cli/ignite/pkg/xnet"
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

	Validators []Validator `yaml:"validators"`
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

	if s.GRPC.Address == base.DefaultGRPCAddress {
		s.GRPC.Address, err = xnet.IncreasePortBy(base.DefaultGRPCAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	if s.GRPCWeb.Address == base.DefaultGRPCWebAddress {
		s.GRPCWeb.Address, err = xnet.IncreasePortBy(base.DefaultGRPCWebAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	if s.API.Address == base.DefaultAPIAddress {
		s.API.Address, err = xnet.IncreasePortBy(base.DefaultAPIAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	if s.P2P.Address == base.DefaultP2PAddress {
		s.P2P.Address, err = xnet.IncreasePortBy(base.DefaultP2PAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	if s.RPC.Address == base.DefaultRPCAddress {
		s.RPC.Address, err = xnet.IncreasePortBy(base.DefaultRPCAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	if s.RPC.PProfAddress == base.DefaultPProfAddress {
		s.RPC.PProfAddress, err = xnet.IncreasePortBy(base.DefaultPProfAddress, inc)
		if err != nil {
			return Servers{}, err
		}
	}

	return s, nil
}
