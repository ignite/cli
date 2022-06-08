package network

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SharePercents []SharePercent

func (sp SharePercents) Empty() bool {
	return len(sp) == 0
}

var rePercentageRequired = regexp.MustCompile(`^[0-9]+.[0-9]*%`)

// SharePercent represent percent of total share
type SharePercent struct {
	denom string
	// in order to avoid using numbers with floating point
	// fractional representation is used: 297/10000 instead of 2.97%
	nominator, denominator uint64
}

// NewSharePercent creates new share percent representation
func NewSharePercent(denom string, nominator, denominator uint64) (SharePercent, error) {
	if denominator < nominator {
		return SharePercent{}, fmt.Errorf("%q can not be bigger than 100", denom)
	}
	return SharePercent{
		denom:       denom,
		nominator:   nominator,
		denominator: denominator,
	}, nil
}

// Share returns coin share of total according to underlying percent
func (p SharePercent) Share(total uint64) (sdk.Coin, error) {
	resultNominator := total * p.nominator
	if resultNominator%p.denominator != 0 {
		err := fmt.Errorf("%s share from total %d is not integer: %f",
			p.denom,
			total,
			float64(resultNominator)/float64(p.denominator),
		)
		return sdk.Coin{}, err
	}
	return sdk.NewInt64Coin(p.denom, int64(resultNominator/p.denominator)), nil
}

// SharePercentFromString parses share percent from string
// format: 11.87%foo
func SharePercentFromString(str string) (SharePercent, error) {
	// validate raw percentage format
	if len(rePercentageRequired.FindStringIndex(str)) == 0 {
		return SharePercent{}, newInvalidPercentageFormat(str)
	}
	var (
		foo        = strings.Split(str, "%")
		fractional = strings.Split(foo[0], ".")
		denom      = foo[1]
	)

	switch len(fractional) {
	case 1:
		nominator, err := strconv.ParseUint(fractional[0], 10, 64)
		if err != nil {
			return SharePercent{}, newInvalidPercentageFormat(str)
		}
		return NewSharePercent(denom, nominator, 100)
	case 2:
		trimmedFractionalPart := strings.TrimRight(fractional[1], "0")
		nominator, err := strconv.ParseUint(fractional[0]+trimmedFractionalPart, 10, 64)
		if err != nil {
			return SharePercent{}, newInvalidPercentageFormat(str)
		}
		return NewSharePercent(denom, nominator, uintPow(10, uint64(len(trimmedFractionalPart)+2)))

	default:
		return SharePercent{}, newInvalidPercentageFormat(str)
	}
}

// ParseSharePercents parses SharePercentage list from string
// format: 12.4%foo,10%bar,0.133%baz
func ParseSharePercents(percents string) (SharePercents, error) {
	rawPercentages := strings.Split(percents, ",")
	ps := make([]SharePercent, len(rawPercentages))
	for i, percentage := range rawPercentages {
		sp, err := SharePercentFromString(percentage)
		if err != nil {
			return nil, err
		}
		ps[i] = sp

	}

	return ps, nil
}

func uintPow(x, y uint64) uint64 {
	var result = x
	for i := 1; uint64(i) < y; i++ {
		result *= x
	}
	return result
}

func newInvalidPercentageFormat(s string) error {
	return fmt.Errorf("invalid percentage format %s", s)
}
