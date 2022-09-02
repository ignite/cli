package v1_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	v1 "github.com/ignite/cli/ignite/chainconfig/v1"
	"github.com/ignite/cli/ignite/pkg/xnet"
)

func TestConfigValidatorDefaultServers(t *testing.T) {
	// Arrange
	c := v1.Config{
		Validators: []v1.Validator{
			{
				Name:   "name-1",
				Bonded: "100ATOM",
			},
		},
	}
	servers := v1.Servers{}

	// Act
	err := c.SetDefaults()
	if err == nil {
		servers, err = c.Validators[0].GetServers()
	}

	// Assert
	require.NoError(t, err)

	// Assert
	require.Equal(t, v1.DefaultGRPCAddress, servers.GRPC.Address)
	require.Equal(t, v1.DefaultGRPCWebAddress, servers.GRPCWeb.Address)
	require.Equal(t, v1.DefaultAPIAddress, servers.API.Address)
	require.Equal(t, v1.DefaultRPCAddress, servers.RPC.Address)
	require.Equal(t, v1.DefaultP2PAddress, servers.P2P.Address)
	require.Equal(t, v1.DefaultPProfAddress, servers.RPC.PProfAddress)
}

func TestConfigValidatorWithExistingServers(t *testing.T) {
	// Arrange
	rpcAddr := "127.0.0.1:1234"
	apiAddr := "127.0.0.1:4321"
	c := v1.Config{
		Validators: []v1.Validator{
			{
				Name:   "name-1",
				Bonded: "100ATOM",
				App: map[string]interface{}{
					// This value should not be ovewritten with the default address
					"api": map[string]interface{}{"address": apiAddr},
				},
				Config: map[string]interface{}{
					// This value should not be ovewritten with the default address
					"rpc": map[string]interface{}{"laddr": rpcAddr},
				},
			},
		},
	}
	servers := v1.Servers{}

	// Act
	err := c.SetDefaults()
	if err == nil {
		servers, err = c.Validators[0].GetServers()
	}

	// Assert
	require.NoError(t, err)

	// Assert
	require.Equal(t, rpcAddr, servers.RPC.Address)
	require.Equal(t, apiAddr, servers.API.Address)
	require.Equal(t, v1.DefaultGRPCAddress, servers.GRPC.Address)
	require.Equal(t, v1.DefaultGRPCWebAddress, servers.GRPCWeb.Address)
	require.Equal(t, v1.DefaultP2PAddress, servers.P2P.Address)
	require.Equal(t, v1.DefaultPProfAddress, servers.RPC.PProfAddress)
}

func TestConfigValidatorsWithExistingServers(t *testing.T) {
	// Arrange
	inc := uint64(10)
	rpcAddr := "127.0.0.1:1234"
	apiAddr := "127.0.0.1:4321"
	c := v1.Config{
		Validators: []v1.Validator{
			{
				Name:   "name-1",
				Bonded: "100ATOM",
			},
			{
				Name:   "name-2",
				Bonded: "200ATOM",
				App: map[string]interface{}{
					// This value should not be ovewritten with the default address
					"api": map[string]interface{}{"address": apiAddr},
				},
				Config: map[string]interface{}{
					// This value should not be ovewritten with the default address
					"rpc": map[string]interface{}{"laddr": rpcAddr},
				},
			},
		},
	}
	servers := v1.Servers{}

	// Act
	err := c.SetDefaults()
	if err == nil {
		servers, err = c.Validators[1].GetServers()
	}

	// Assert
	require.NoError(t, err)

	// Assert: The existing addresses should not be changed
	require.Equal(t, rpcAddr, servers.RPC.Address)
	require.Equal(t, apiAddr, servers.API.Address)

	// Assert: The second validator should have the ports incremented by 10
	require.Equal(t, xnet.MustIncreasePortBy(v1.DefaultGRPCAddress, inc), servers.GRPC.Address)
	require.Equal(t, xnet.MustIncreasePortBy(v1.DefaultGRPCWebAddress, inc), servers.GRPCWeb.Address)
	require.Equal(t, xnet.MustIncreasePortBy(v1.DefaultP2PAddress, inc), servers.P2P.Address)
	require.Equal(t, xnet.MustIncreasePortBy(v1.DefaultPProfAddress, inc), servers.RPC.PProfAddress)
}

func TestConfigValidatorsDefaultServers(t *testing.T) {
	// Arrange
	inc := uint64(10)
	c := v1.Config{
		Validators: []v1.Validator{
			{
				Name:   "name-1",
				Bonded: "100ATOM",
			},
			{
				Name:   "name-2",
				Bonded: "200ATOM",
			},
		},
	}
	servers := v1.Servers{}

	// Act
	err := c.SetDefaults()
	if err == nil {
		servers, err = c.Validators[1].GetServers()
	}

	// Assert
	require.NoError(t, err)

	// Assert: The second validator should have the ports incremented by 10
	require.Equal(t, xnet.MustIncreasePortBy(v1.DefaultGRPCAddress, inc), servers.GRPC.Address)
	require.Equal(t, xnet.MustIncreasePortBy(v1.DefaultGRPCWebAddress, inc), servers.GRPCWeb.Address)
	require.Equal(t, xnet.MustIncreasePortBy(v1.DefaultAPIAddress, inc), servers.API.Address)
	require.Equal(t, xnet.MustIncreasePortBy(v1.DefaultRPCAddress, inc), servers.RPC.Address)
	require.Equal(t, xnet.MustIncreasePortBy(v1.DefaultP2PAddress, inc), servers.P2P.Address)
	require.Equal(t, xnet.MustIncreasePortBy(v1.DefaultPProfAddress, inc), servers.RPC.PProfAddress)
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
	require.Equal(t, c, c2)
}
