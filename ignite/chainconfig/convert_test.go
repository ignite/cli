package chainconfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite-hq/cli/ignite/chainconfig"
	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	v0testdata "github.com/ignite-hq/cli/ignite/chainconfig/v0/testdata"
	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"
)

func TestConvertLatest(t *testing.T) {
	origin := v0testdata.GetInitialV0Config()
	result, err := chainconfig.ConvertLatest(origin)
	assert.Nil(t, err)
	expected := v0testdata.GetConvertedLatestConfig()

	cfg, ok := result.(*v1.Config)
	require.Truef(t, ok, "expected v1 config, got %T", result)

	require.Equal(t, config.Version(0), origin.Version())
	require.Equal(t, config.Version(1), result.Version())
	require.Equal(t, origin.Faucet, cfg.Faucet)
	require.Equal(t, origin.Client, cfg.Client)
	require.Equal(t, origin.Build, cfg.Build)
	require.Equal(t, origin.Host.RPC, cfg.Validators[0].GetRPC())
	require.Equal(t, origin.Host.P2P, cfg.Validators[0].GetP2P())
	require.Equal(t, origin.Host.GRPC, cfg.Validators[0].GetGRPC())
	require.Equal(t, origin.Host.GRPCWeb, cfg.Validators[0].GetGRPCWeb())
	require.Equal(t, origin.Host.Prof, cfg.Validators[0].GetProf())
	require.Equal(t, origin.Host.API, cfg.Validators[0].GetAPI())
	require.Equal(t, origin.Genesis, cfg.Genesis)
	require.Equal(t, origin.ListAccounts(), cfg.ListAccounts())
	require.Equal(t, origin.Init.KeyringBackend, cfg.Validators[0].KeyringBackend)
	require.Equal(t, origin.Init.Client, cfg.Validators[0].Client)
	require.Equal(t, origin.Init.Home, cfg.Validators[0].Home)
	require.Equal(t, expected, result)
}
