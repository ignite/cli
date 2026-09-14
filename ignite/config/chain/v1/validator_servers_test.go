package v1_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	v1 "github.com/ignite/cli/v29/ignite/config/chain/v1"
	"github.com/ignite/cli/v29/ignite/pkg/xyaml"
)

func TestValidatorGetServers(t *testing.T) {
	// Arrange
	want := v1.DefaultServers()
	want.RPC.Address = "127.0.0.0:1"
	want.P2P.Address = "127.0.0.0:2"
	want.GRPC.Address = "127.0.0.0:3"
	want.GRPCWeb.Address = "127.0.0.0:4"
	want.RPC.PProfAddress = "127.0.0.0:5"
	want.API.Address = "127.0.0.0:6"

	v := v1.Validator{
		App: map[string]any{
			"grpc":     map[string]any{"address": want.GRPC.Address},
			"grpc-web": map[string]any{"address": want.GRPCWeb.Address},
			"api":      map[string]any{"address": want.API.Address},
		},
		Config: map[string]any{
			"p2p": map[string]any{"laddr": want.P2P.Address},
			"rpc": map[string]any{
				"laddr":       want.RPC.Address,
				"pprof_laddr": want.RPC.PProfAddress,
			},
		},
	}

	// Act
	s, err := v.GetServers()

	// Assert
	require.NoError(t, err)
	require.Equal(t, want, s)
}

func TestValidatorSetServers(t *testing.T) {
	// Arrange
	v := v1.Validator{}
	s := v1.DefaultServers()
	wantApp := xyaml.Map{
		"grpc":     map[string]any{"address": s.GRPC.Address},
		"grpc-web": map[string]any{"address": s.GRPCWeb.Address},
		"api":      map[string]any{"address": s.API.Address},
	}
	wantConfig := xyaml.Map{
		"p2p": map[string]any{"laddr": s.P2P.Address},
		"rpc": map[string]any{
			"laddr":       s.RPC.Address,
			"pprof_laddr": s.RPC.PProfAddress,
		},
	}

	// Act
	err := v.SetServers(s)

	// Assert
	require.NoError(t, err)
	require.Equal(t, wantApp, v.App, "cosmos app config is not equal")
	require.Equal(t, wantConfig, v.Config, "tendermint config is not equal")
}
