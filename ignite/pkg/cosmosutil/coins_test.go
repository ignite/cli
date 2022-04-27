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
			"20%foo,50%staking",
			sdk.NewCoins(sdk.NewInt64Coin("foo", 20), sdk.NewInt64Coin("staking", 50)),
			nil,
		},
		{
			"wrong format",
			"20nova,50baz",
			sdk.NewCoins(sdk.NewInt64Coin("nova", 20), sdk.NewInt64Coin("baz", 50)),
			errors.New("amount for 20nova has to have a % after the number"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coins, err := cosmosutil.ParseCoinsNormalizedWithPercentageRequired(tt.coins)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.parsed, coins)
			}
		})
	}
}
