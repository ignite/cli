package v0_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	v0testdata "github.com/ignite/cli/ignite/config/chain/v0/testdata"
	v1 "github.com/ignite/cli/ignite/config/chain/v1"
	"github.com/ignite/cli/ignite/config/chain/version"
)

func TestV0ToV1(t *testing.T) {
	// Arrange
	cfgV0 := v0testdata.GetConfig(t)

	// Act
	c, err := cfgV0.ConvertNext()
	cfgV1, _ := c.(*v1.Config)

	// Assert
	require.NoError(t, err)
	require.NotNilf(t, cfgV1, "expected *v1.Config, got %T", c)
	require.Equal(t, version.Version(1), cfgV1.GetVersion())
	require.Equal(t, cfgV0.Build, cfgV1.Build)
	require.Equal(t, cfgV0.Accounts, cfgV1.Accounts)
	require.Equal(t, cfgV0.Faucet, cfgV1.Faucet)
	require.Equal(t, cfgV0.Client, cfgV1.Client)
	require.Equal(t, cfgV0.Genesis, cfgV1.Genesis)
	require.Len(t, cfgV1.Validators, 1)
}

func TestV0ToV1Validator(t *testing.T) {
	// Arrange
	cfgV0 := v0testdata.GetConfig(t)
	cfgV0.Host.RPC = "127.0.0.0:1"
	cfgV0.Host.P2P = "127.0.0.0:2"
	cfgV0.Host.GRPC = "127.0.0.0:3"
	cfgV0.Host.GRPCWeb = "127.0.0.0:4"
	cfgV0.Host.Prof = "127.0.0.0:5"
	cfgV0.Host.API = "127.0.0.0:6"

	// Act
	c, _ := cfgV0.ConvertNext()
	cfgV1, _ := c.(*v1.Config)
	validator := cfgV1.Validators[0]
	servers, _ := validator.GetServers()

	// Assert
	require.Equal(t, cfgV0.Validator.Name, validator.Name)
	require.Equal(t, cfgV0.Validator.Staked, validator.Bonded)
	require.Equal(t, cfgV0.Init.Home, validator.Home)
	require.Equal(t, cfgV0.Init.Client, validator.Client)
	require.Equal(t, cfgV0.Host.RPC, servers.RPC.Address)
	require.Equal(t, cfgV0.Host.P2P, servers.P2P.Address)
	require.Equal(t, cfgV0.Host.GRPC, servers.GRPC.Address)
	require.Equal(t, cfgV0.Host.GRPCWeb, servers.GRPCWeb.Address)
	require.Equal(t, cfgV0.Host.Prof, servers.RPC.PProfAddress)
	require.Equal(t, cfgV0.Host.API, servers.API.Address)
}
