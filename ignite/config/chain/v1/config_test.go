package v1_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/config/chain/base"
	v1 "github.com/ignite/cli/ignite/config/chain/v1"
	"github.com/ignite/cli/ignite/pkg/xnet"
)

func TestConfigDecode(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	f, err := os.Open("testdata/config2.yaml")
	require.NoError(err)
	defer f.Close()
	var cfg v1.Config

	err = cfg.Decode(f)

	require.NoError(err)
	expected := v1.Config{
		Config: base.Config{
			Version: 1,
			Build: base.Build{
				Binary: "evmosd",
				Proto: base.Proto{
					Path:            "proto",
					ThirdPartyPaths: []string{"third_party/proto", "proto_vendor"},
				},
			},
			Accounts: []base.Account{
				{
					Name:     "alice",
					Coins:    []string{"100000000uatom", "100000000000000000000aevmos"},
					Mnemonic: "ozone unfold device pave lemon potato omit insect column wise cover hint narrow large provide kidney episode clay notable milk mention dizzy muffin crazy",
				},
				{
					Name:    "bob",
					Coins:   []string{"5000000000000aevmos"},
					Address: "cosmos1adn9gxjmrc3hrsdx5zpc9sj2ra7kgqkmphf8yw",
				},
			},
			Faucet: base.Faucet{
				Name:  &[]string{"bob"}[0],
				Coins: []string{"10aevmos"},
				Host:  "0.0.0.0:4600",
				Port:  4600,
			},
			Genesis: map[string]any{
				"app_state": map[string]any{
					"crisis": map[string]any{
						"constant_fee": map[string]any{
							"denom": "aevmos",
						},
					},
				},
				"chain_id": "evmosd_9000-1",
			},
		},
		Validators: []v1.Validator{{
			Name:   "alice",
			Bonded: "100000000000000000000aevmos",
			Home:   "$HOME/.evmosd",
			App: map[string]any{
				"evm-rpc": map[string]any{
					"address":    "0.0.0.0:8545",
					"ws-address": "0.0.0.0:8546",
				},
			},
		}},
	}
	assert.Equal(expected, cfg)
}

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
	require.Equal(t, base.DefaultGRPCAddress, servers.GRPC.Address)
	require.Equal(t, base.DefaultGRPCWebAddress, servers.GRPCWeb.Address)
	require.Equal(t, base.DefaultAPIAddress, servers.API.Address)
	require.Equal(t, base.DefaultRPCAddress, servers.RPC.Address)
	require.Equal(t, base.DefaultP2PAddress, servers.P2P.Address)
	require.Equal(t, base.DefaultPProfAddress, servers.RPC.PProfAddress)
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
	require.Equal(t, base.DefaultGRPCAddress, servers.GRPC.Address)
	require.Equal(t, base.DefaultGRPCWebAddress, servers.GRPCWeb.Address)
	require.Equal(t, base.DefaultP2PAddress, servers.P2P.Address)
	require.Equal(t, base.DefaultPProfAddress, servers.RPC.PProfAddress)
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
	require.Equal(t, xnet.MustIncreasePortBy(base.DefaultGRPCAddress, inc), servers.GRPC.Address)
	require.Equal(t, xnet.MustIncreasePortBy(base.DefaultGRPCWebAddress, inc), servers.GRPCWeb.Address)
	require.Equal(t, xnet.MustIncreasePortBy(base.DefaultP2PAddress, inc), servers.P2P.Address)
	require.Equal(t, xnet.MustIncreasePortBy(base.DefaultPProfAddress, inc), servers.RPC.PProfAddress)
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
	require.Equal(t, xnet.MustIncreasePortBy(base.DefaultGRPCAddress, inc), servers.GRPC.Address)
	require.Equal(t, xnet.MustIncreasePortBy(base.DefaultGRPCWebAddress, inc), servers.GRPCWeb.Address)
	require.Equal(t, xnet.MustIncreasePortBy(base.DefaultAPIAddress, inc), servers.API.Address)
	require.Equal(t, xnet.MustIncreasePortBy(base.DefaultRPCAddress, inc), servers.RPC.Address)
	require.Equal(t, xnet.MustIncreasePortBy(base.DefaultP2PAddress, inc), servers.P2P.Address)
	require.Equal(t, xnet.MustIncreasePortBy(base.DefaultPProfAddress, inc), servers.RPC.PProfAddress)
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
	c2, err := c.Clone()

	// Assert
	require.NoError(t, err)
	require.Equal(t, c, c2)
}
