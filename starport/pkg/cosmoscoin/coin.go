// Package cosmoscoin provides utilities to deal with SDK coins.
package cosmoscoin

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var (
	errInvalidCoin = errors.New("coin is invalid")
)

var (
	reDnmString = `[a-zA-Z][a-zA-Z0-9/]{2,127}`
	reDecAmt    = `[[:digit:]]+(?:\.[[:digit:]]+)?|\.[[:digit:]]+`
	reSpc       = `[[:space:]]*`
	pattern     = fmt.Sprintf(`^(%s)%s(%s)$`, reDecAmt, reSpc, reDnmString)
	parseRe     = regexp.MustCompile(pattern)
)

// Coin represents a SDK coin
type Coin struct {
	// Amount of coins
	Amount uint64

	// Denom of the coin
	Denom string
}

// String stringifies the coin
func (c Coin) String() string {
	return fmt.Sprintf("%d%s", c.Amount, c.Denom)
}

// Parse parses a string into a coin containing an amount and a denom.
func Parse(c string) (coin Coin, err error) {
	parsed := parseRe.FindStringSubmatch(c)

	if len(parsed) != 3 {
		return coin, errInvalidCoin
	}

	amountStr := parsed[1]
	coin.Denom = parsed[2]

	coin.Amount, err = strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return coin, errInvalidCoin
	}

	return
}
