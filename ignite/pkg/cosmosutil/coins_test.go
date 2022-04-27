package cosmosutil_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosutil"
	"github.com/stretchr/testify/require"
)

func TestValidateCoinsStrWithPercentage(t *testing.T) {
	tests := []struct {
		name   string
		coins  string
		parsed sdk.Coins
		err    error
	}{
		{
			"format is OK",
			"%20foo,%50staking",
			sdk.NewCoins(sdk.NewInt64Coin("foo", 20), sdk.NewInt64Coin("staking", 50)),
			nil,
		},
		{
			"wrong format",
			"20foo,50staking",
			sdk.NewCoins(sdk.NewInt64Coin("foo", 20), sdk.NewInt64Coin("staking", 50)),
			errors.New("amount for 20foo has to have a % prefix"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coins, err := cosmosutil.ParseCoinsNormalizedWithPercentage(tt.coins)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.parsed, coins)
			}
		})
	}
}
