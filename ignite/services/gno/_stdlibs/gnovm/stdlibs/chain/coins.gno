package chain

import (
	"math/overflow"
	"sort"
	"strconv"
)

// Coin holds some amount of one currency.
// A negative amount is invalid.
type Coin struct {
	Denom  string `json:"denom"`
	Amount int64  `json:"amount"`
}

// NewCoin returns a new coin with a denomination and amount
func NewCoin(denom string, amount int64) Coin {
	return Coin{
		Denom:  denom,
		Amount: amount,
	}
}

// String provides a human-readable representation of a coin
func (c Coin) String() string {
	return strconv.Itoa(int(c.Amount)) + c.Denom
}

// IsGTE returns true if they are the same type and the receiver is
// an equal or greater value
func (c Coin) IsGTE(other Coin) bool {
	mustMatchDenominations(c.Denom, other.Denom)

	return c.Amount >= other.Amount
}

// IsLT returns true if they are the same type and the receiver is
// a smaller value
func (c Coin) IsLT(other Coin) bool {
	mustMatchDenominations(c.Denom, other.Denom)

	return c.Amount < other.Amount
}

// IsEqual returns true if the two sets of Coins have the same value
func (c Coin) IsEqual(other Coin) bool {
	mustMatchDenominations(c.Denom, other.Denom)

	return c.Amount == other.Amount
}

// Add adds amounts of two coins with same denom.
// If the coins differ in denom then it panics.
// An overflow or underflow panics.
//
// The sign of the result is not checked. Adding two negative coins returns a
// negative one, so test the sign yourself where the value is a balance.
func (c Coin) Add(other Coin) Coin {
	mustMatchDenominations(c.Denom, other.Denom)

	sum, ok := overflow.Add64(c.Amount, other.Amount)
	if !ok {
		panic("coin add overflow/underflow: " + strconv.Itoa(int(c.Amount)) + " + " + strconv.Itoa(int(other.Amount)))
	}

	c.Amount = sum
	return c
}

// Sub subtracts amounts of two coins with same denom.
// If the coins differ in denom then it panics.
// An overflow or underflow panics.
//
// The sign of the result is not checked. Taking more than the receiver holds
// returns a negative coin, so test the sign yourself where the value is a
// balance.
func (c Coin) Sub(other Coin) Coin {
	mustMatchDenominations(c.Denom, other.Denom)

	dff, ok := overflow.Sub64(c.Amount, other.Amount)
	if !ok {
		panic("coin sub overflow/underflow: " + strconv.Itoa(int(c.Amount)) + " - " + strconv.Itoa(int(other.Amount)))
	}
	c.Amount = dff

	return c
}

// IsPositive returns true if coin amount is positive.
func (c Coin) IsPositive() bool {
	return c.Amount > 0
}

// IsNegative returns true if the coin amount is negative and false otherwise.
func (c Coin) IsNegative() bool {
	return c.Amount < 0
}

// IsZero returns true if the amount of given coin is zero
func (c Coin) IsZero() bool {
	return c.Amount == 0
}

func mustMatchDenominations(denomA, denomB string) {
	if denomA != denomB {
		panic("incompatible coin denominations: " + denomA + ", " + denomB)
	}
}

// Coins is a set of Coin, one per currency
type Coins []Coin

// NewCoins returns a new set of Coins given one or more Coins.
// Internally, it joins coins with the same mechanism as [Coins.Add], so the
// resulting [Coins] will be sorted and one per coin, with the caveats mentioned
// in [Coins.Add]'s godoc.
//
// The argument is never modified, including when it is a slice passed with
// `NewCoins(myVar...)`. Joining two coins of one denom shifts the rest of the
// set down a place, so working in the caller's array would leave the caller
// holding a repeat of its own last coin — value it never had. [Coins.Add]
// copies for the same reason.
func NewCoins(coins ...Coin) Coins {
	cz := make(Coins, len(coins))
	copy(cz, coins)
	return cz.sortAndJoin()
}

func (cz Coins) sortAndJoin() Coins {
	sort.Sort(coinsSort(cz))

	// Join each coin in cz[i] with cz[i+1] if possible.
	for i := 0; i+1 < len(cz); {
		current, next := cz[i], cz[i+1]
		if current.Denom == next.Denom {
			// Consolidate next into current.
			sum, ok := overflow.Add64(current.Amount, next.Amount)
			if !ok {
				panic("coin add overflow/underflow: " + strconv.Itoa(int(current.Amount)) +
					" + " + strconv.Itoa(int(next.Amount)) + " " + next.Denom)
			}
			if sum == 0 {
				// current + next cancel each other out. Remove both.
				cz = append(cz[:i], cz[i+2:]...)
			} else {
				cz[i].Amount = sum
				// Remove next from the slice.
				cz = append(cz[:i+1], cz[i+2:]...)
			}
		} else {
			// Move onto next coin.
			i++
		}
	}
	return cz[:len(cz):len(cz)]
}

// String returns the string representation of Coins
func (cz Coins) String() string {
	if len(cz) == 0 {
		return ""
	}

	res := ""
	for i, c := range cz {
		if i > 0 {
			res += ","
		}
		res += c.String()
	}

	return res
}

// AmountOf returns the amount of a specific coin from the Coins set.
// Returns 0 if the denom is not found. Panics if there are duplicate coins with
// the given denom.
func (cz Coins) AmountOf(denom string) int64 {
	amt, set := int64(0), false
	for _, c := range cz {
		if c.Denom == denom {
			if set {
				panic("duplicate denom " + denom + " in coin set")
			}
			amt, set = c.Amount, true
		}
	}

	return amt
}

// Add adds the given Coins to the set.
//
// The coins of the receiver and the argument are added together into a new
// Coins slice. Coins are then sorted and all coins with the same denom are
// joined together in one single denom. Consequently:
//
//   - If cz or b contain multiple coins with the same Denom (though generally
//     obtained as a result of an incorrect operation), they will be joined
//     together.
//   - Operations may unexpectedly overflow in some edge cases depending on the
//     ordering of values after sorting (only if there are coins with duplicate
//     denoms).
//   - The resulting Coins slice will be sorted by ascending denoms.
//   - If any coin has an amount of 0 after it is joined with another coin with
//     the same denom, it will be removed.
//   - The receiver and the parameter will not be modified.
func (cz Coins) Add(b Coins) Coins {
	// Create dst then sort to have O(n*log(n)) complexity.
	dst := make(Coins, len(cz)+len(b))
	copy(dst[:len(cz)], cz)
	copy(dst[len(cz):], b)
	return dst.sortAndJoin()
}

// IsZero returns true if there are no coins or all coins are zero.
func (coins Coins) IsZero() bool {
	for _, coin := range coins {
		if !coin.IsZero() {
			return false
		}
	}
	return true
}

type coinsSort []Coin

func (coins coinsSort) Len() int { return len(coins) }

func (coins coinsSort) Less(i, j int) bool {
	return coins[i].Denom < coins[j].Denom
}

func (coins coinsSort) Swap(i, j int) {
	coins[i], coins[j] = coins[j], coins[i]
}

func CoinDenom(pkgPath, coinName string) string {
	// TODO: Possibly remove after https://github.com/gnolang/gno/issues/3164
	// Similar to ibc spec
	// ibc_denom := 'ibc/' + hash('path' + 'base_denom')
	// gno_qualified_denom := '/' + 'pkg_path' + ':' + 'base_denom'
	return "/" + pkgPath + ":" + coinName
}
