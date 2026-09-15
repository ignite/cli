package runtime

import (
	"testing"

	"github.com/gnolang/gno/tm2/pkg/crypto"
	tm2std "github.com/gnolang/gno/tm2/pkg/std"
	"github.com/stretchr/testify/require"
)

// TotalCoin sums the table, so AmountOf — which validates before reading — is
// reached only once per held denom. With nothing held the loop does not run, and
// only TotalCoin's own check keeps `gno test` from reporting a malformed denom as
// zero where SDKBanker.TotalCoin panics. That divergence is the one this branch
// cares about: a realm whose test suite passed would fail in production.
//
// GetCoin needs no such test — AmountOf validates as its first statement there,
// which is why its explicit check was removed rather than pinned.
func TestTotalCoinRejectsAMalformedDenomWithNothingHeld(t *testing.T) {
	t.Parallel()

	empty := &TestBanker{CoinTable: map[crypto.Bech32Address]tm2std.Coins{}}
	require.PanicsWithError(t, "invalid denom: UPPERCASE", func() {
		empty.TotalCoin("UPPERCASE")
	})

	// A well-formed denom nobody holds has zero supply, not an error — so this
	// cannot pass with a TotalCoin that rejects everything.
	require.Zero(t, empty.TotalCoin("ugnot"))
}
