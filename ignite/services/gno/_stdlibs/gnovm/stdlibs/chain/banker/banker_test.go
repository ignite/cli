package banker

import (
	"testing"

	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/stdlibs/internal/execctx"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/std"
	"github.com/stretchr/testify/require"
)

// countingBanker records whether a send actually reached the bank.
type countingBanker struct{ sends int }

func (b *countingBanker) GetCoins(crypto.Bech32Address) std.Coins          { return nil }
func (b *countingBanker) GetCoin(crypto.Bech32Address, string) int64       { return 0 }
func (b *countingBanker) SendCoins(_, _ crypto.Bech32Address, _ std.Coins) { b.sends++ }
func (b *countingBanker) TotalCoin(string) int64                           { return 0 }
func (b *countingBanker) IssueCoin(crypto.Bech32Address, string, int64)    {}
func (b *countingBanker) RemoveCoin(crypto.Bech32Address, string, int64)   {}

// An origin-send banker may only spend the envelope of the realm the envelope
// was credited to. The primary control is the lifetime pin in banker.gno, which
// makes such a banker unpersistable; this gate is the layer behind it, and only
// fires for a banker that reached a later message anyway. Nothing else can
// exercise it, because Gno code cannot defeat the pin -- so it is pinned here,
// by calling the native directly.
func TestOriginSendIsSpendableOnlyByThePaidRealm(t *testing.T) {
	paid := gno.DerivePkgBech32Addr("gno.land/r/paid")
	other := gno.DerivePkgBech32Addr("gno.land/r/other")
	to := gno.DerivePkgBech32Addr("gno.land/r/dest")
	envelope := std.Coins{{Denom: "ugnot", Amount: 100}}

	newMachine := func(bk *countingBanker) *gno.Machine {
		observed := false
		return &gno.Machine{Context: execctx.ExecContext{
			OriginSend:          envelope,
			OriginSendSpent:     new(std.Coins),
			OriginSendRecipient: paid,
			OriginSendObserved:  &observed,
			Banker:              bk,
		}}
	}

	send := func(m *gno.Machine, from crypto.Bech32Address) (panicked bool) {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		X_bankerSendCoins(m, btOriginSend, string(from), string(to),
			[]string{"ugnot"}, []int64{10})
		return false
	}

	// A realm that was not paid must not move the envelope.
	bk := &countingBanker{}
	send(newMachine(bk), other)
	require.Zero(t, bk.sends,
		"a realm that did not receive the envelope must not spend it")

	// The realm that was paid may, so the check above is not just refusing
	// everything.
	bk = &countingBanker{}
	m := newMachine(bk)
	require.False(t, send(m, paid), "the paid realm must be able to spend")
	require.Equal(t, 1, bk.sends, "the paid realm's send must reach the bank")
}
