package chainconfig

import (
	"testing"

	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"

	"github.com/ignite-hq/cli/ignite/chainconfig/common"
	v0 "github.com/ignite-hq/cli/ignite/chainconfig/v0"
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
	//require.Equal(t, origin.Host, result.(*v1.Config).GetHost())
	require.Equal(t, origin.Genesis, result.(*v1.Config).Genesis)
	require.Equal(t, origin.ListAccounts(), result.(*v1.Config).ListAccounts())
	//require.Equal(t, origin.GetInit(), result.(*v1.Config).GetInit())
	require.Equal(t, expected, result)
}
