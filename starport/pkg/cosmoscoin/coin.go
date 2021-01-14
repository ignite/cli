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

// Parse parses a coin into amount and denom.
func Parse(c string) (amount uint64, denom string, err error) {
	parsed := parseRe.FindStringSubmatch(c)

	if len(parsed) != 3 {
		return 0, "", errInvalidCoin
	}

	amountStr := parsed[1]
	denom = parsed[2]

	amount, err = strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return 0, "", errInvalidCoin
	}

	return
}
