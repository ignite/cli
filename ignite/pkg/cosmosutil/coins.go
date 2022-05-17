package cosmosutil

import (
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var rePercentageRequired = regexp.MustCompile(`^[0-9]+%`)

// ParseCoinsNormalizedWithPercentageRequired parses coins by requiring percentages.
// format: 20%foo,50%staking
func ParseCoinsNormalizedWithPercentageRequired(coins string) (sdk.Coins, error) {
	trimPercentage := func(s string) string {
		return strings.ReplaceAll(s, "%", "")
	}

	s := strings.Split(coins, ",")
	for _, ss := range s {
		if len(rePercentageRequired.FindStringIndex(ss)) == 0 {
			return nil, fmt.Errorf("amount for %s has to have a %% after the number", trimPercentage(ss))
		}
	}
	c, err := sdk.ParseCoinsNormalized(trimPercentage(coins))
	if err != nil {
		return nil, err
	}
	for _, coin := range c {
		if coin.Amount.Int64() > 100 {
			return nil, fmt.Errorf("%q can not be bigger than 100", coin.Denom)
		}
	}
	return c, nil
}
