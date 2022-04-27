package cosmosutil

import (
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var rePercentageRequired = regexp.MustCompile(`[0-9]+%`)

// ParseCoinsNormalizedWithPercentageRequired parses coins by requiring percentages.
// format: 20%foo,50%staking
func ParseCoinsNormalizedWithPercentageRequired(coins string) (sdk.Coins, error) {
	s := strings.Split(coins, ",")
	for _, ss := range s {
		if len(rePercentageRequired.FindStringIndex(ss)) == 0 {
			return nil, fmt.Errorf("amount for %s has to have a %% after the number", ss)
		}
	}
	trimmedCoins := strings.ReplaceAll(coins, "%", "")
	return sdk.ParseCoinsNormalized(trimmedCoins)
}
