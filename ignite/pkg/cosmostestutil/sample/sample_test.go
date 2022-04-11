package cosmostestutilsample

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestAccAddress(t *testing.T) {
	got := AccAddress()
	require.NotEmpty(t, got)
	_, err := sdk.AccAddressFromBech32(got)
	require.NoError(t, err)
}
