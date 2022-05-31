package chainconfig

import (
	"testing"

	"github.com/ignite/cli/ignite/chainconfig/common"
	v0 "github.com/ignite/cli/ignite/chainconfig/v0"
	v1 "github.com/ignite/cli/ignite/chainconfig/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertLatest(t *testing.T) {
	origin := v0.GetInitialV0Config()
	result, err := ConvertLatest(origin)
	assert.Nil(t, err)
	expected := v0.GetConvertedLatestConfig()

	require.Equal(t, common.Version(0), origin.Version())
	require.Equal(t, common.Version(1), result.Version())
	require.Equal(t, origin.Faucet, result.(*v1.Config).Faucet)
	require.Equal(t, origin.Client, result.(*v1.Config).Client)
	require.Equal(t, origin.Build, result.(*v1.Config).Build)
	require.Equal(t, origin.Host.RPC, result.(*v1.Config).Validators[0].GetRPC())
	require.Equal(t, origin.Host.P2P, result.(*v1.Config).Validators[0].GetP2P())
	require.Equal(t, origin.Host.GRPC, result.(*v1.Config).Validators[0].GetGRPC())
	require.Equal(t, origin.Host.GRPCWeb, result.(*v1.Config).Validators[0].GetGRPCWeb())
	require.Equal(t, origin.Host.Prof, result.(*v1.Config).Validators[0].GetProf())
	require.Equal(t, origin.Host.API, result.(*v1.Config).Validators[0].GetAPI())
	require.Equal(t, origin.Genesis, result.(*v1.Config).Genesis)
	require.Equal(t, origin.ListAccounts(), result.(*v1.Config).ListAccounts())
	require.Equal(t, origin.Init.KeyringBackend, result.(*v1.Config).Validators[0].KeyringBackend)
	require.Equal(t, origin.Init.Client, result.(*v1.Config).Validators[0].Client)
	require.Equal(t, origin.Init.Home, result.(*v1.Config).Validators[0].Home)
	require.Equal(t, expected, result)
}
