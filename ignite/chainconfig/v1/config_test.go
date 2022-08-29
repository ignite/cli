package v1_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	v1 "github.com/ignite/cli/ignite/chainconfig/v1"
	"github.com/ignite/cli/ignite/pkg/xnet"
)

func TestValidatorAddresses(t *testing.T) {
	// Arrange
	grpcAddr := xnet.AnyIPv4Address(9090)
	grpcWebAddr := xnet.AnyIPv4Address(9091)
	apiAddr := xnet.AnyIPv4Address(1317)
	rpcAddr := xnet.AnyIPv4Address(26657)
	p2pAddr := xnet.AnyIPv4Address(26656)
	pprofAddr := xnet.AnyIPv4Address(6060)

	// Act
	v := v1.Validator{
		Name:   "test-name-1",
		Bonded: "102ATOM",
		App: map[string]interface{}{
			"grpc":     map[string]interface{}{"address": grpcAddr},
			"grpc-web": map[string]interface{}{"address": grpcWebAddr},
			"api":      map[string]interface{}{"address": apiAddr},
		},
		Config: map[string]interface{}{
			"rpc":         map[string]interface{}{"laddr": rpcAddr},
			"p2p":         map[string]interface{}{"laddr": p2pAddr},
			"pprof_laddr": pprofAddr,
		},
	}

	// Assert
	require.Equal(t, grpcAddr, v.GetGRPC())
	require.Equal(t, grpcWebAddr, v.GetGRPCWeb())
	require.Equal(t, apiAddr, v.GetAPI())
	require.Equal(t, rpcAddr, v.GetRPC())
	require.Equal(t, p2pAddr, v.GetP2P())
	require.Equal(t, pprofAddr, v.GetProf())
}

func TestConfigSetDefault(t *testing.T) {
	// Arrange
	inc := 10
	c := v1.Config{
		Validators: []v1.Validator{
			{
				Name:   "name-1",
				Bonded: "100ATOM",
				Config: map[string]interface{}{
					// This address should be overwritten
					"rpc": map[string]interface{}{"laddr": "127.0.0.1:1234"},
				},
			},
			{
				Name:   "name-2",
				Bonded: "200ATOM",
				App: map[string]interface{}{
					// This address should be overwritten
					"api": map[string]interface{}{"address": "127.0.0.1:8888"},
				},
			},
		},
	}

	// Act
	err := c.SetDefaults()

	// Assert
	require.NoError(t, err)

	// Assert: First validator
	require.Equal(t, xnet.AnyIPv4Address(v1.GRPCPort), c.Validators[0].GetGRPC())
	require.Equal(t, xnet.AnyIPv4Address(v1.GRPCWebPort), c.Validators[0].GetGRPCWeb())
	require.Equal(t, xnet.AnyIPv4Address(v1.APIPort), c.Validators[0].GetAPI())
	require.Equal(t, xnet.AnyIPv4Address(v1.RPCPort), c.Validators[0].GetRPC())
	require.Equal(t, xnet.AnyIPv4Address(v1.P2PPort), c.Validators[0].GetP2P())
	require.Equal(t, xnet.AnyIPv4Address(v1.PProfPort), c.Validators[0].GetProf())

	// Assert: Second validator
	require.Equal(t, xnet.AnyIPv4Address(v1.GRPCPort+inc), c.Validators[1].GetGRPC())
	require.Equal(t, xnet.AnyIPv4Address(v1.GRPCWebPort+inc), c.Validators[1].GetGRPCWeb())
	require.Equal(t, xnet.AnyIPv4Address(v1.APIPort+inc), c.Validators[1].GetAPI())
	require.Equal(t, xnet.AnyIPv4Address(v1.RPCPort+inc), c.Validators[1].GetRPC())
	require.Equal(t, xnet.AnyIPv4Address(v1.P2PPort+inc), c.Validators[1].GetP2P())
	require.Equal(t, xnet.AnyIPv4Address(v1.PProfPort+inc), c.Validators[1].GetProf())
}

func TestClone(t *testing.T) {
	// Arrange
	c := &v1.Config{
		Validators: []v1.Validator{
			{
				Name:   "alice",
				Bonded: "100000000stake",
			},
		},
	}

	// Act
	c2 := c.Clone()

	// Assert
	require.EqualValues(t, c, c2)
}
