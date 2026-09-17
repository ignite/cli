package banker

import (
	"fmt"

	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/stdlibs/internal/execctx"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/std"
)

// BankerInterface is the interface through which Gno is capable of accessing
// the blockchain's banker.
//
// The name is what it is to avoid a collision with Gno's Banker, when
// transpiling.
type BankerInterface interface {
	GetCoins(addr crypto.Bech32Address) (dst std.Coins)
	GetCoin(addr crypto.Bech32Address, denom string) int64
	SendCoins(from, to crypto.Bech32Address, amt std.Coins)
	TotalCoin(denom string) int64
	IssueCoin(addr crypto.Bech32Address, denom string, amount int64)
	RemoveCoin(addr crypto.Bech32Address, denom string, amount int64)
}

const (
	// Can only read state.
	btReadonly uint8 = iota //nolint
	// Can only send from tx send: at most the coins the current message's
	// send-envelope credited to the realm's own address, and only if that
	// realm is the envelope's recipient.
	//
	// The recipient half is checked here at use time, in X_bankerSendCoins,
	// rather than only at construction in banker.gno. It is a second layer:
	// the primary rule is that an OriginSend banker cannot outlive its
	// message at all, which banker.gno enforces structurally by pinning it
	// to a realm value the persistence walk refuses to store. This check
	// does not depend on that reasoning holding.
	btOriginSend
	// Can send from all realm coins.
	btRealmSend
	// Can issue and remove realm coins.
	btRealmIssue
)

func X_bankerGetCoins(m *gno.Machine, bt uint8, addr string) (denoms []string, amounts []int64) {
	coins := execctx.GetContext(m).Banker.GetCoins(crypto.Bech32Address(addr))
	return ExpandCoins(coins)
}

func X_bankerGetCoin(m *gno.Machine, bt uint8, addr string, denom string) int64 {
	return execctx.GetContext(m).Banker.GetCoin(crypto.Bech32Address(addr), denom)
}

func X_bankerSendCoins(m *gno.Machine, bt uint8, fromS, toS string, denoms []string, amounts []int64) {
	// bt != BankerTypeReadonly (checked in gno)

	ctx := execctx.GetContext(m)
	amt := CompactCoins(denoms, amounts)
	from, to := crypto.Bech32Address(fromS), crypto.Bech32Address(toS)

	switch bt {
	case btOriginSend:
		// Second layer. The primary rule is that an OriginSend banker
		// cannot outlive the message that created it — enforced
		// structurally in banker.gno, which pins such a banker to a
		// realm value that the persistence walk refuses to store.
		//
		// This gate is independent of that: the envelope of this message
		// was credited to exactly one address (ctx.OriginSendRecipient),
		// and `from` is already pinned to the banker's bound pkgAddr by
		// banker.gno's SendCoins, so it asserts the spending realm is
		// the realm that actually received the envelope. It holds even
		// for a banker built and used entirely within one message, and
		// it does not depend on any Gno-side lifetime reasoning.
		//
		// NOTE: this check alone is NOT sufficient. It only stops a
		// revived banker whose realm is not the current entry realm; a
		// banker delegated to another realm and stashed there revives at
		// full envelope value in any later message where the grantor
		// happens to be the entry realm. That is why the lifetime pin,
		// not this gate, is the primary control.
		if from != ctx.OriginSendRecipient {
			m.PanicString(
				fmt.Sprintf(
					`cannot send from "%v": it did not receive the origin send of this message`,
					from),
			)
			return
		}
		// Past the gate above, `from` is the address the envelope was
		// credited to, so this banker belongs to the realm that was paid.
		// Forwarding the envelope through a banker counts as noticing the
		// payment just as much as reading it does — the limit check below
		// consults ctx.OriginSend. Mark here rather than after the limit
		// check: a realm that gets this far is payable even if the send is
		// then rejected for being too large, and failing it twice over
		// would be wrong.
		//
		// Unconditional by design: the owner may have handed this banker to
		// another realm to spend within this message, so the realm running
		// right now is not necessarily the one that was paid.
		ctx.MarkOriginSendObserved()
		// indirection allows us to "commit" in a second phase
		spent := (*ctx.OriginSendSpent).Add(amt)
		if !ctx.OriginSend.IsAllGTE(spent) {
			m.PanicString(
				fmt.Sprintf(
					`cannot send "%v", limit "%v" exceeded with "%v" already spent`,
					amt, ctx.OriginSend, *ctx.OriginSendSpent),
			)
			return
		}
		ctx.Banker.SendCoins(from, to, amt)
		*ctx.OriginSendSpent = spent
	case btRealmSend, btRealmIssue:
		ctx.Banker.SendCoins(from, to, amt)
	default:
		panic(fmt.Sprintf("invalid banker type %d in bankerSendCoins", bt))
	}
}

func X_bankerTotalCoin(m *gno.Machine, bt uint8, denom string) int64 {
	return execctx.GetContext(m).Banker.TotalCoin(denom)
}

func X_bankerIssueCoin(m *gno.Machine, bt uint8, addr string, denom string, amount int64) {
	execctx.GetContext(m).Banker.IssueCoin(crypto.Bech32Address(addr), denom, amount)
}

func X_bankerRemoveCoin(m *gno.Machine, bt uint8, addr string, denom string, amount int64) {
	execctx.GetContext(m).Banker.RemoveCoin(crypto.Bech32Address(addr), denom, amount)
}

func ExpandCoins(c std.Coins) (denoms []string, amounts []int64) {
	denoms = make([]string, len(c))
	amounts = make([]int64, len(c))
	for i, coin := range c {
		denoms[i] = coin.Denom
		amounts[i] = coin.Amount
	}
	return denoms, amounts
}

func CompactCoins(denoms []string, amounts []int64) std.Coins {
	coins := make(std.Coins, len(denoms))
	for i := range coins {
		coins[i] = std.Coin{Denom: denoms[i], Amount: amounts[i]}
	}
	return coins
}
