package v0

import (
	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
)

// ConvertNext implements the conversion of the current config to the next version.
func (c *Config) ConvertNext() (config.Converter, error) {
	targetConfig := &v1.Config{}

	// All the fields in the base config remain the same
	targetConfig.BaseConfig = c.BaseConfig

	// Change the version to 1
	targetConfig.Version = config.Version(1)

	// There is only one validator in v0. Set it as the only one validator in v1.
	validators := make([]v1.Validator, 0, 1)
	validator := v1.Validator{}

	sourceValidator := c.Validator
	validator.Name = sourceValidator.Name
	validator.Bonded = sourceValidator.Staked

	// Set the fields in Init to the target validator
	validator.Home = c.Init.Home
	validator.KeyringBackend = c.Init.KeyringBackend
	validator.Client = c.Init.Client

	validator.App = c.Init.App
	validator.Config = c.Init.Config

	// If the fields in Host is not empty, we need to merge them into validator.app and validator.config.
	if c.Host.P2P != "" {
		if validator.Config == nil {
			// Create an empty map for Config
			validator.Config = make(map[string]interface{})
		}
		p2p := map[string]interface{}{"laddr": c.Host.P2P}
		validator.Config["p2p"] = p2p
	}

	if c.Host.Prof != "" {
		if validator.Config == nil {
			// Create an empty map for Config
			validator.Config = make(map[string]interface{})
		}
		validator.Config["pprof_laddr"] = c.Host.Prof
	}

	if c.Host.RPC != "" {
		if validator.Config == nil {
			// Create an empty map for Config
			validator.Config = make(map[string]interface{})
		}
		rpc := map[string]interface{}{"laddr": c.Host.RPC}
		validator.Config["rpc"] = rpc
	}

	if c.Host.GRPCWeb != "" {
		if validator.App == nil {
			// Create an empty map for App
			validator.App = make(map[string]interface{})
		}
		grpcweb := map[string]interface{}{"address": c.Host.GRPCWeb}
		validator.App["grpc-web"] = grpcweb
	}

	if c.Host.GRPC != "" {
		if validator.App == nil {
			// Create an empty map for App
			validator.App = make(map[string]interface{})
		}
		grpc := map[string]interface{}{"address": c.Host.GRPC}
		validator.App["grpc"] = grpc
	}

	if c.Host.API != "" {
		if validator.App == nil {
			// Create an empty map for App
			validator.App = make(map[string]interface{})
		}
		api := map[string]interface{}{"address": c.Host.API}
		validator.App["api"] = api
	}

	validators = append(validators, validator)
	targetConfig.Validators = validators

	return targetConfig, nil
}
