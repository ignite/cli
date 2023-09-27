package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/ignite/mars/testutil/keeper"
	"github.com/ignite/mars/x/mars/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.MarsKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
