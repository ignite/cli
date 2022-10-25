package masterclass_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "masterclass/testutil/keeper"
	"masterclass/testutil/nullify"
	"masterclass/x/masterclass"
	"masterclass/x/masterclass/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MasterclassKeeper(t)
	masterclass.InitGenesis(ctx, *k, genesisState)
	got := masterclass.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
