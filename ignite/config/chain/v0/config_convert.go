package v0

import (
	v1 "github.com/ignite/cli/ignite/config/chain/v1"
	"github.com/ignite/cli/ignite/config/chain/version"
)

// ConvertNext converts the current config version to the next one.
func (c *Config) ConvertNext() (version.Converter, error) {
	targetCfg := v1.DefaultConfig()

	// All the fields in the base config remain the same
	targetCfg.Config = c.Config
	targetCfg.Version = 1

	// There is always only one validator in version 0
	validator := v1.Validator{}
	validator.Name = c.Validator.Name
	validator.Bonded = c.Validator.Staked
	validator.Home = c.Init.Home
	validator.Client = c.Init.Client

	if c.Init.App != nil {
		validator.App = c.Init.App
	}

	if c.Init.Config != nil {
		validator.Config = c.Init.Config
	}

	// The host configuration must be defined in the validators for version 1
	servers := v1.Servers{}

	if c.Host.P2P != "" {
		servers.P2P.Address = c.Host.P2P
	}

	if c.Host.RPC != "" {
		servers.RPC.Address = c.Host.RPC
	}

	if c.Host.Prof != "" {
		servers.RPC.PProfAddress = c.Host.Prof
	}

	if c.Host.GRPCWeb != "" {
		servers.GRPCWeb.Address = c.Host.GRPCWeb
	}

	if c.Host.GRPC != "" {
		servers.GRPC.Address = c.Host.GRPC
	}

	if c.Host.API != "" {
		servers.API.Address = c.Host.API
	}

	if err := validator.SetServers(servers); err != nil {
		return nil, err
	}

	targetCfg.Validators = append(targetCfg.Validators, validator)

	return targetCfg, nil
}
