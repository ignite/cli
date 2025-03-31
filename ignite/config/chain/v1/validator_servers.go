package v1

import (
	"github.com/mitchellh/mapstructure"

	baseconfig "github.com/ignite/cli/v29/ignite/config/chain/defaults"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

func DefaultServers() Servers {
	s := Servers{}
	s.GRPC.Address = baseconfig.GRPCAddress
	s.GRPCWeb.Address = baseconfig.GRPCWebAddress
	s.API.Address = baseconfig.APIAddress
	s.P2P.Address = baseconfig.P2PAddress
	s.RPC.Address = baseconfig.RPCAddress
	s.RPC.PProfAddress = baseconfig.PProfAddress

	return s
}

type Servers struct {
	cosmosServers     `mapstructure:",squash"`
	tendermintServers `mapstructure:",squash"`
}

type cosmosServers struct {
	GRPC    CosmosHost `mapstructure:"grpc"`
	GRPCWeb CosmosHost `mapstructure:"grpc-web"`
	API     CosmosHost `mapstructure:"api"`
}

type tendermintServers struct {
	P2P TendermintHost    `mapstructure:"p2p"`
	RPC TendermintRPCHost `mapstructure:"rpc"`
}

type CosmosHost struct {
	Address string `mapstructure:"address,omitempty"`
}

type TendermintHost struct {
	Address string `mapstructure:"laddr,omitempty"`
}

type TendermintRPCHost struct {
	TendermintHost `mapstructure:",squash"`

	PProfAddress string `mapstructure:"pprof_laddr,omitempty"`
}

func (v Validator) GetServers() (Servers, error) {
	// Initialize servers with default addresses
	s := DefaultServers()

	// Overwrite the default Cosmos SDK addresses with the configured ones
	if err := mapstructure.Decode(v.App, &s); err != nil {
		return Servers{}, errors.Errorf("error reading validator app servers: %w", err)
	}

	// Overwrite the default Tendermint addresses with the configured ones
	if err := mapstructure.Decode(v.Config, &s); err != nil {
		return Servers{}, errors.Errorf("error reading tendermint validator config servers: %w", err)
	}

	return s, nil
}

func (v *Validator) SetServers(s Servers) error {
	if err := v.setAppServers(s); err != nil {
		return errors.Errorf("error updating validator app servers: %w", err)
	}

	if err := v.setConfigServers(s); err != nil {
		return errors.Errorf("error updating validator config servers: %w", err)
	}

	return nil
}

func (v *Validator) setAppServers(s Servers) error {
	c, err := decodeServers(s.cosmosServers)
	if err != nil {
		return err
	}

	v.App = mergeMaps(c, v.App)

	return nil
}

func (v *Validator) setConfigServers(s Servers) error {
	m, err := decodeServers(s.tendermintServers)
	if err != nil {
		return errors.Errorf("error updating validator config servers: %w", err)
	}

	v.Config = mergeMaps(m, v.Config)

	return nil
}

func decodeServers(input interface{}) (output map[string]interface{}, err error) {
	// Decode the input structure into a map
	if err := mapstructure.Decode(input, &output); err != nil {
		return nil, err
	}

	// Remove keys with empty server values from the map
	for k := range output {
		if v, _ := output[k].(map[string]interface{}); len(v) == 0 {
			delete(output, k)
		}
	}

	// Don't return an empty map to avoid the generation of empty
	// fields when the validator is saved to a YAML config file.
	if len(output) == 0 {
		return nil, nil
	}

	return
}

func mergeMaps(src, dst map[string]interface{}) map[string]interface{} {
	if len(src) == 0 {
		return dst
	}

	// Allow dst to be nil by initializing it here
	if dst == nil {
		dst = make(map[string]interface{})
	}

	for k, v := range src {
		// When the current value is a map in both merge their values
		if srcValue, ok := v.(map[string]interface{}); ok {
			if dstValue, ok := dst[k].(map[string]interface{}); ok {
				mergeMaps(srcValue, dstValue)

				continue
			}
		}

		// By default overwrite the destination map with the source value
		dst[k] = v
	}

	return dst
}
