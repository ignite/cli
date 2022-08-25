package v0

import (
	"github.com/imdario/mergo"

	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
)

// ConvertNext convets the current config version to the next one.
func (c *Config) ConvertNext() (config.Converter, error) {
	targetCfg := v1.DefaultConfig()

	// All the fields in the base config remain the same
	targetCfg.BaseConfig = c.BaseConfig

	// There is always only one validator in version 0
	validator := v1.Validator{}
	validator.Name = c.Validator.Name
	validator.Bonded = c.Validator.Staked
	validator.Home = c.Init.Home
	validator.KeyringBackend = c.Init.KeyringBackend
	validator.Client = c.Init.Client

	if c.Init.App != nil {
		validator.App = c.Init.App
	}

	if c.Init.Config != nil {
		validator.Config = c.Init.Config
	}

	// The host configuration must be defined in the validators for version 1
	configValues := make(map[string]interface{})
	appValues := make(map[string]interface{})

	if c.Host.P2P != "" {
		configValues["p2p"] = map[string]interface{}{"laddr": c.Host.P2P}
	}

	if c.Host.RPC != "" {
		configValues["rpc"] = map[string]interface{}{"laddr": c.Host.RPC}
	}

	if c.Host.Prof != "" {
		configValues["pprof_laddr"] = c.Host.Prof
	}

	if c.Host.GRPCWeb != "" {
		appValues["grpc-web"] = map[string]interface{}{"address": c.Host.GRPCWeb}
	}

	if c.Host.GRPC != "" {
		appValues["grpc"] = map[string]interface{}{"address": c.Host.GRPC}
	}

	if c.Host.API != "" {
		appValues["api"] = map[string]interface{}{"address": c.Host.API}
	}

	if err := mergo.Merge(&validator.Config, configValues, mergo.WithOverride); err != nil {
		return nil, err
	}

	if err := mergo.Merge(&validator.App, appValues, mergo.WithOverride); err != nil {
		return nil, err
	}

	targetCfg.Validators = append(targetCfg.Validators, validator)

	return targetCfg, nil
}
