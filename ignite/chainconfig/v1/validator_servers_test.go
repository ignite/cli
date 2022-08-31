package v1_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	v1 "github.com/ignite/cli/ignite/chainconfig/v1"
)

func TestValidatorGetServers(t *testing.T) {
	// Arrange
	want := v1.DefaultServers()
	v := v1.Validator{
		App: map[string]interface{}{
			"grpc":     map[string]interface{}{"address": want.GRPC.Address},
			"grpc-web": map[string]interface{}{"address": want.GRPCWeb.Address},
			"api":      map[string]interface{}{"address": want.API.Address},
		},
		Config: map[string]interface{}{
			"p2p": map[string]interface{}{"laddr": want.P2P.Address},
			"rpc": map[string]interface{}{
				"laddr":       want.RPC.Address,
				"pprof_laddr": want.RPC.PProfAddress,
			},
		},
	}

	// Act
	s, err := v.GetServers()

	// Assert
	require.NoError(t, err)
	require.EqualValues(t, want, s)
}

func TestValidatorSetServers(t *testing.T) {
	// Arrange
	v := v1.Validator{}
	s := v1.DefaultServers()
	wantApp := map[string]interface{}{
		"grpc":     map[string]interface{}{"address": s.GRPC.Address},
		"grpc-web": map[string]interface{}{"address": s.GRPCWeb.Address},
		"api":      map[string]interface{}{"address": s.API.Address},
	}
	wantConfig := map[string]interface{}{
		"p2p": map[string]interface{}{"laddr": s.P2P.Address},
		"rpc": map[string]interface{}{
			"laddr":       s.RPC.Address,
			"pprof_laddr": s.RPC.PProfAddress,
		},
	}

	// Act
	err := v.SetServers(s)

	// Assert
	require.NoError(t, err)
	require.EqualValues(t, wantApp, v.App, "cosmos app config is not equal")
	require.EqualValues(t, wantConfig, v.Config, "tendermint config is not equal")
}
