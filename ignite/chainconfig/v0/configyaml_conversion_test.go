package v0

import (
	"testing"

	v1 "github.com/ignite-hq/cli/ignite/chainconfig/v1"

	"github.com/ignite-hq/cli/ignite/chainconfig/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertNext(t *testing.T) {
	origin := GetInitialV0Config()
	result, err := origin.ConvertNext()
	assert.Nil(t, err)
	expected := GetConvertedLatestConfig()

	require.Equal(t, common.Version(0), origin.Version())
	require.Equal(t, common.Version(1), result.Version())
	require.Equal(t, origin.GetFaucet(), result.(*v1.Config).GetFaucet())
	require.Equal(t, origin.GetClient(), result.(*v1.Config).GetClient())
	require.Equal(t, origin.GetBuild(), result.(*v1.Config).GetBuild())
	require.Equal(t, origin.GetHost(), result.(*v1.Config).GetHost())
	require.Equal(t, origin.GetGenesis(), result.(*v1.Config).GetGenesis())
	require.Equal(t, origin.ListAccounts(), result.(*v1.Config).ListAccounts())
	require.Equal(t, origin.GetInit(), result.(*v1.Config).GetInit())
	require.Equal(t, expected, result)
}
