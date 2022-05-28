package conversion

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/chainconfig/common"
	v0 "github.com/ignite/cli/ignite/chainconfig/v0"
)

func TestConvertLatest(t *testing.T) {
	origin := v0.GetInitialV0Config()
	result, err := ConvertLatest(origin)
	assert.Nil(t, err)
	expected := v0.GetConvertedLatestConfig()

	require.Equal(t, common.Version(0), origin.Version())
	require.Equal(t, common.Version(1), result.Version())
	require.Equal(t, origin.GetFaucet(), result.GetFaucet())
	require.Equal(t, origin.GetClient(), result.GetClient())
	require.Equal(t, origin.GetBuild(), result.GetBuild())
	require.Equal(t, origin.GetHost(), result.GetHost())
	require.Equal(t, origin.GetGenesis(), result.GetGenesis())
	require.Equal(t, origin.ListAccounts(), result.ListAccounts())
	require.Equal(t, origin.GetInit(), result.GetInit())
	require.Equal(t, expected, result)
}
