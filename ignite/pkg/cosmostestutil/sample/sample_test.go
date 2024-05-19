package cosmostestutilsample

import (
	"testing"

	"github.com/stretchr/testify/require"

	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
)

func TestAccAddress(t *testing.T) {
	got := AccAddress()
	require.NotEmpty(t, got)
	exampleAccountAddress := addresscodec.NewBech32Codec("cosmos")
	_, err := exampleAccountAddress.StringToBytes(got)
	require.NoError(t, err)
}
