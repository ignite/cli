package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "masterclass/testutil/keeper"
	"masterclass/x/masterclass/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.MasterclassKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
