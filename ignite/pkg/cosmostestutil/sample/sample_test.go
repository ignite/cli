package cosmostestutilsample

import (
	"testing"

	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/stretchr/testify/require"
)

func TestAccAddress(t *testing.T) {
	got := AccAddress()
	require.NotEmpty(t, got)
	exampleAccountAddress := addresscodec.NewBech32Codec("cosmos")
	_, err := exampleAccountAddress.StringToBytes(got)
	require.NoError(t, err)
}
