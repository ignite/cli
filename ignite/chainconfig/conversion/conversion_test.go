package conversion

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v0 "github.com/ignite/cli/ignite/chainconfig/v0"
	v1 "github.com/ignite/cli/ignite/chainconfig/v1"
)

func TestConvertLatest(t *testing.T) {
	origin := v0.GetInitialV0Config()
	result, err := ConvertLatest(origin)
	assert.Nil(t, err)
	expected := v0.GetConvertedLatestConfig()

	require.Equal(t, 0, origin.GetVersion())
	require.Equal(t, 1, result.GetVersion())
	require.Equal(t, origin.GetFaucet(), result.GetFaucet())
	require.Equal(t, origin.GetClient(), result.GetClient())
	require.Equal(t, origin.GetBuild(), result.GetBuild())
	require.Equal(t, origin.GetHost(), result.GetHost())
	require.Equal(t, origin.GetGenesis(), result.GetGenesis())
	require.Equal(t, origin.ListAccounts(), result.ListAccounts())
	require.Equal(t, origin.GetInit(), result.GetInit())
	require.Equal(t, expected.(*v1.Config), result.(*v1.Config))
}
