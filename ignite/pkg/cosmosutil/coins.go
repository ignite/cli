package cosmosutil

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ParseCoinsNormalizedWithPercentage parses coins with percentage prefix.
// format: %20foo,%50staking
func ParseCoinsNormalizedWithPercentage(coins string) (sdk.Coins, error) {
	s := strings.Split(coins, ",")
	for _, ss := range s {
		if !strings.HasPrefix(ss, "%") {
			return nil, fmt.Errorf("amount for %s has to have a %% prefix", ss)
		}
	}
	trimmedCoins := strings.ReplaceAll(coins, "%", "")
	return sdk.ParseCoinsNormalized(trimmedCoins)
}
