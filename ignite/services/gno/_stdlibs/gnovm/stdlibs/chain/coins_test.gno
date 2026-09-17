package chain_test

import (
	"chain"
	"testing"
)

func TestCoinIsZero(t *testing.T) {
	coin := chain.NewCoin("ugnot", 0)
	res := coin.IsZero()
	if !res {
		t.Fatalf("expected IsZero=true for %s", coin)
	}

	coin = chain.NewCoin("ugnot", 1)
	res = coin.IsZero()
	if res {
		t.Fatalf("expected IsZero=false for %s", coin)
	}
}

func TestCoinsIsZero(t *testing.T) {
	var (
		testDenom1 = "denom1"
		testDenom2 = "denom2"
	)
	tests := []struct {
		name     string
		coins    chain.Coins
		expected bool
	}{
		{
			name:     "empty coins",
			coins:    chain.Coins{},
			expected: true,
		},
		{
			name:     "one zero denom",
			coins:    chain.Coins{chain.NewCoin(testDenom1, 0)},
			expected: true,
		},
		{
			name: "two zero denom",
			coins: chain.Coins{
				chain.NewCoin(testDenom1, 0), chain.NewCoin(testDenom2, 0),
			},
			expected: true,
		},
		{
			name:     "one non-zero denom",
			coins:    chain.Coins{chain.NewCoin(testDenom1, 1)},
			expected: false,
		},
		{
			name: "one zero denom, one non-zero denom",
			coins: chain.Coins{
				chain.NewCoin(testDenom1, 0), chain.NewCoin(testDenom2, 1),
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zero := tt.coins.IsZero()

			if zero != tt.expected {
				t.Fatalf("expected IsZero=%t for coins=%v", tt.expected, tt.coins)
			}
		})
	}
}

func TestCoins_AmountOf(t *testing.T) {
	coins := chain.Coins{
		chain.NewCoin("denom1", 1),
		chain.NewCoin("denom2", 2),
	}

	amount := coins.AmountOf("denom1")
	if amount != 1 {
		t.Fatalf("expected amount=1 for denom1, got %d", amount)
	}

	amount = coins.AmountOf("denom2")
	if amount != 2 {
		t.Fatalf("expected amount=2 for denom2, got %d", amount)
	}

	amount = coins.AmountOf("denom3")
	if amount != 0 {
		t.Fatalf("expected amount=0 for denom3, got %d", amount)
	}

	coins = append(coins, chain.NewCoin("denom1", 1))

	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fatal("expected panic when getting duplicate denom")
			}
		}()
		coins.AmountOf("denom1")
	}()
}

const maxInt64 int64 = (1 << 63) - 1

func shouldPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic but did not panic")
		}
	}()
	f()
}

func TestNewCoins_SortsJoinsAndRemovesZero(t *testing.T) {
	// b cancels out to zero, a sums to 4
	res := chain.NewCoins(
		chain.NewCoin("b", 2),
		chain.NewCoin("a", 1),
		chain.NewCoin("b", -2),
		chain.NewCoin("a", 3),
	)

	if len(res) != 1 {
		t.Fatalf("expected 1 coin after join, got %d: %v", len(res), res)
	}
	if res[0].Denom != "a" || res[0].Amount != 4 {
		t.Fatalf("unexpected result: %v", res)
	}
}

func TestNewCoins_LeavesInputSliceAlone(t *testing.T) {
	// Passing a slice with the variadic expansion (s...) shares the array with
	// NewCoins, so it must copy before sorting.
	s := []chain.Coin{
		chain.NewCoin("b", 2),
		chain.NewCoin("a", 1),
	}
	_ = chain.NewCoins(s...)

	if s[0].Denom != "b" || s[1].Denom != "a" {
		t.Fatalf("NewCoins reordered the caller's slice: %v", s)
	}
}

func TestNewCoins_DoesNotInflateInputSlice(t *testing.T) {
	// Two coins of one denom make NewCoins join them, which shifts the rest of
	// the set down a place. Done in the caller's array that would overwrite the
	// caller's own coins with a repeat of the last one, so the caller would see
	// more value than it put in and a denom listed twice.
	s := []chain.Coin{
		chain.NewCoin("a", 1),
		chain.NewCoin("a", 2),
		chain.NewCoin("b", 3),
	}

	joined := chain.NewCoins(s...)
	if joined.String() != "3a,3b" {
		t.Fatalf("expected the returned set to join the duplicates, got: %v", joined)
	}

	total := int64(0)
	for _, c := range s {
		total += c.Amount
	}
	if total != 6 {
		t.Fatalf("caller's slice total changed from 6 to %d: %v", total, s)
	}
	if s[2].Denom != "b" || s[2].Amount != 3 {
		t.Fatalf("caller's slice was overwritten: %v", s)
	}
}

func TestNewCoins_OrderBasedPanic(t *testing.T) {
	cz := chain.Coins{
		chain.NewCoin("z", 1),
		chain.NewCoin("z", -1),
		chain.NewCoin("z", maxInt64),
	}
	// Should not panic: 1 + (-1) + maxInt should not cause a panic.
	chain.NewCoins(cz...)

	cz = chain.Coins{
		chain.NewCoin("z", 1),
		chain.NewCoin("z", maxInt64),
		chain.NewCoin("z", -1),
	}
	// Should panic: 1 + maxInt causes an overflow, regardless of whether
	// it is followed up by a -1.
	shouldPanic(t, func() { chain.NewCoins(cz...) })
}

func TestCoinsAdd_Basic(t *testing.T) {
	cz := chain.NewCoins(chain.NewCoin("b", 2), chain.NewCoin("a", 1))
	other := chain.NewCoins(chain.NewCoin("a", 3), chain.NewCoin("c", 5))

	res := cz.Add(other)

	// Expect sorted result and proper joining: a:4, b:2, c:5
	if len(res) != 3 {
		t.Fatalf("expected 3 coins in result, got %d: %v", len(res), res)
	}
	if res[0].Denom != "a" || res[0].Amount != 4 {
		t.Fatalf("unexpected a coin: %v", res[0])
	}
	if res[1].Denom != "b" || res[1].Amount != 2 {
		t.Fatalf("unexpected b coin: %v", res[1])
	}
	if res[2].Denom != "c" || res[2].Amount != 5 {
		t.Fatalf("unexpected c coin: %v", res[2])
	}

	// Ensure inputs were not mutated (Add should not modify inputs)
	if cz.AmountOf("a") != 1 || cz.AmountOf("b") != 2 {
		t.Fatalf("cz was modified: %v", cz)
	}
	if other.AmountOf("a") != 3 || other.AmountOf("c") != 5 {
		t.Fatalf("other was modified: %v", other)
	}
}

func TestCoinsAdd_RemovesZero(t *testing.T) {
	// Test that sums that result in zero are removed.
	cz := chain.NewCoins(chain.NewCoin("a", 5))
	other := chain.NewCoins(chain.NewCoin("a", -5), chain.NewCoin("b", 1))

	res := cz.Add(other)
	// 'a' should be removed because 5 + (-5) == 0
	if len(res) != 1 {
		t.Fatalf("expected only 1 coin after cancellation, got %d: %v", len(res), res)
	}
	if res[0].Denom != "b" || res[0].Amount != 1 {
		t.Fatalf("unexpected remaining coin: %v", res[0])
	}
}

func TestCoinsAdd_OverflowPanicsWhenSumming(t *testing.T) {
	cz := chain.Coins{chain.NewCoin("x", maxInt64)}
	other := chain.Coins{chain.NewCoin("x", 1)}
	shouldPanic(t, func() { _ = cz.Add(other) })
}
