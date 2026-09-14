package banker

import (
	"chain"
	"strconv"
	"strings"
)

// Realm functions can call banker.NewBanker(options) to get a banker
// instance. A banker instance can be passed onto other realm functions
// (this allows other realms to spend coins on behalf of the first
// realm) and can be stored in package-level vars across transactions —
// the underlying struct (BankerType + scoped pkgAddr) persists like any
// other Gno value, and the methods continue to dispatch through the
// runtime against live bank state on each call. See
// gno.land/pkg/integration/testdata/banker_persistence.txtar for the
// full cross-realm capability-persistence behavior.
//
// SECURITY: because the struct persists and NewBanker authorizes only at
// construction (SendCoins re-checks nothing but pkgAddr == from), minting
// a banker is an IRREVOCABLE grant of spending authority over that
// address. Passing a realm value to code you do not fully trust is
// therefore equivalent to handing over that authority permanently: the
// callee can mint a banker, keep it in a package-level var, and spend in
// later transactions with no further involvement from you. Unlike the
// realm value itself — which cannot be persisted (refusePersistRealmHIV)
// and cannot be reused once stale — the banker has no expiry, no
// revocation, and no liveness re-check. This applies with equal force to
// sub-realm identities (cur.Sub): the sub token is ephemeral, but a
// banker minted from it is not, so any access check your realm performs
// before spending (a frozen flag, an allowlist, a paused switch) is
// bypassed entirely by a retained banker, which reaches the bank keeper
// without re-entering your code.
//
// BankerTypeOriginSend is the one exception, and is NOT persistable. Its
// authority is the send-envelope of a single message, so it is pinned to
// that message: it carries a live realm value, which the persistence walk
// refuses to store, so any attempt to keep one in realm state aborts the
// transaction before any coins move. For that type alone the paragraph
// above does not apply — the grant expires when the message ends. See
// [NewBanker]. Every other banker type behaves exactly as described.
//
// Banker panics on errors instead of returning errors.
// This also helps simplify the interface and prevent
// hidden bugs (e.g. ignoring errors)
//
// NOTE: this Gno interface is satisfied by a native go
// type, and those can't return non-primitive objects
// (without confusion).
type Banker interface {
	GetCoins(addr address) (dst chain.Coins)
	GetCoin(addr address, denom string) int64
	SendCoins(from, to address, amt chain.Coins)
	TotalCoin(denom string) int64
	IssueCoin(addr address, denom string, amount int64)
	RemoveCoin(addr address, denom string, amount int64)
}

// BankerType represents the "permission level" requested for a banker,
// retrievable through [NewBanker].
type BankerType uint8

// Available types of banker.
const (
	// Can only read state.
	BankerTypeReadonly BankerType = iota
	// Can only send from tx send: at most the coins this message's
	// send-envelope credited to the realm's own address. Valid only
	// inside the message that created it (it cannot be persisted) and
	// only while that realm is the envelope's recipient.
	BankerTypeOriginSend
	// Can send from all realm coins.
	BankerTypeRealmSend
	// Can issue and remove realm coins.
	BankerTypeRealmIssue

	maxBanker
)

func (b BankerType) String() string {
	switch b {
	case BankerTypeReadonly:
		return "BankerTypeReadonly"
	case BankerTypeOriginSend:
		return "BankerTypeOriginSend"
	case BankerTypeRealmSend:
		return "BankerTypeRealmSend"
	case BankerTypeRealmIssue:
		return "BankerTypeRealmIssue"
	default:
		panic("invalid BankerType: " + strconv.Itoa(int(b)))
	}
}

//----------------------------------------
// adapter for native banker

// NewReadonlyBanker returns a banker that can only read coin state.
// Readonly bankers have no authority to mutate anything (no SendCoins,
// IssueCoin, or RemoveCoin) so they do not require a realm capability
// — anyone, including Render functions and non-crossing helpers, can
// construct one. Use this in queries; use NewBanker for any operation
// that moves coins.
func NewReadonlyBanker() Banker {
	return banker{bt: BankerTypeReadonly}
}

// NewBanker returns a new Banker bound to rlm — anyone holding rlm can
// construct a banker for it. The realm value is the authorization
// capability: a function that takes `cur realm` (crossing function, or
// a non-first-arg-threaded function) can call NewBanker; a function
// without a realm handle in scope cannot.
//
// bt must be one of BankerTypeOriginSend, BankerTypeRealmSend, or
// BankerTypeRealmIssue. Use NewReadonlyBanker for BankerTypeReadonly —
// the readonly path requires no realm capability.
func NewBanker(bt BankerType, rlm realm) Banker {
	if bt == BankerTypeReadonly {
		panic("use NewReadonlyBanker for BankerTypeReadonly")
	}
	if bt >= maxBanker {
		panic("invalid banker type")
	}
	// rlm authorizes spending/issuing from rlm.Address(); require it to
	// be the live cur. Without this, a callee could pass cur.Previous()
	// (IsCurrent()==false, addr = caller) and the pkgAddr==from gate in
	// SendCoins would then authorize draining the CALLER's balance. Live
	// sub-tokens still pass (IsCurrent is parent-anchored for them).
	if !rlm.IsCurrent() {
		panic("banker can only be instantiated for the current realm")
	}
	// Sub-realm tokens (synthesized "host#subpath" pkgpath) may only
	// construct RealmSend bankers: RealmIssue would mint denoms whose
	// pkgpath component carries the "#" separator, and OriginSend's
	// guard below is meaningless for a sub-address that never holds the
	// origin-send envelope. Gate on bt first (an int compare) so the
	// common Readonly/RealmSend paths skip the pkgpath split;
	// fail-closed for any future banker type.
	if bt != BankerTypeRealmSend {
		if _, _, isSub := chain.SplitPkgSubPath(rlm.PkgPath()); isSub {
			panic(bt.String() + " not supported for sub-realm tokens")
		}
	}
	if bt == BankerTypeOriginSend && !rlm.Previous().IsUserCall() {
		panic("banker with type BankerTypeOriginSend can only be instantiated by the origin package")
	}
	b := banker{
		bt:      bt,
		pkgAddr: rlm.Address(),
		pkgPath: rlm.PkgPath(),
	}
	if bt == BankerTypeOriginSend {
		// Pin the banker to the message that created it. An OriginSend
		// banker's authority IS the send-envelope of one message, so it
		// must not outlive that message; the construction-time check
		// above is worthless against a value stashed in realm state and
		// woken up later. Carrying a live (non-origin) realm handle makes
		// the value structurally unpersistable — the realm finalization
		// walk refuses realm values outright (refusePersistRealmHIV in
		// gnovm/pkg/gnolang/uverse.go), so any attempt to store this
		// banker in realm state aborts the whole transaction. It stays
		// freely passable to other realms WITHIN the message, which is
		// the delegation the type exists for.
		//
		// This is structural, not a token: there is no per-message value
		// for a caller to forge and no context field for a construction
		// site to forget to populate.
		b.rlm = rlm
	}
	return b
}

// These are native bindings to the banker's functions.
func bankerGetCoins(bt uint8, addr string) (denoms []string, amounts []int64)
func bankerGetCoin(bt uint8, addr string, denom string) int64
func bankerSendCoins(bt uint8, from, to string, denoms []string, amounts []int64)
func bankerTotalCoin(bt uint8, denom string) int64
func bankerIssueCoin(bt uint8, addr string, denom string, amount int64)
func bankerRemoveCoin(bt uint8, addr string, denom string, amount int64)

type banker struct {
	bt      BankerType
	pkgAddr address
	pkgPath string // realm pkgpath, for IssueCoin/RemoveCoin denom-prefix.
	// rlm is set only for BankerTypeOriginSend, where it is the
	// message-lifetime pin described in [NewBanker]: realm values are
	// refused by the persistence walk, so a banker holding one cannot be
	// written to realm state. Nil for every other banker type — those
	// are genuine cross-transaction capabilities (see
	// gno.land/pkg/integration/testdata/banker_persistence.txtar) and
	// keep persisting exactly as before.
	rlm realm
}

// IsCanonical reports whether b is the canonical bank-keeper-backed
// Banker produced by NewBanker. The concrete `banker` type is unexported,
// so only this package can construct values that satisfy this assertion;
// hand-rolled Banker implementations (no-op fakes, mocks, decorators)
// return false. Use this in cap-token wrappers (e.g. AllowanceSender) to
// reject untrusted Banker substitutions at construction time, before
// any state mutation depends on the banker actually moving coins.
//
// Forward-compat: if this package ever introduces additional canonical
// Banker implementations alongside `banker`, this function MUST be
// updated to recognize them. New bank-keeper-backed types must be added
// to the type-switch explicitly — the current single-type assertion is
// a deliberate allowlist, not an inadvertent restriction.
func IsCanonical(b Banker) bool {
	_, ok := b.(banker)
	return ok
}

func (b banker) GetCoins(addr address) (dst chain.Coins) {
	denoms, amounts := bankerGetCoins(uint8(b.bt), string(addr))
	dst = make(chain.Coins, len(denoms))
	for i := range dst {
		dst[i] = chain.Coin{denoms[i], amounts[i]}
	}
	return dst
}

// GetCoin returns addr's balance of a single denom.
//
// Prefer this to GetCoins when you want one denom: GetCoins returns every denom
// the address holds and costs accordingly, and an address can be sent denoms it
// never asked for (see IssueCoin), so its cost is not under your control.
// GetCoin returns addr's balance of a single denom, reading nothing else. Returns
// zero for a denom addr does not hold.
func (b banker) GetCoin(addr address, denom string) int64 {
	return bankerGetCoin(uint8(b.bt), string(addr), denom)
}

func (b banker) SendCoins(from, to address, amt chain.Coins) {
	if b.bt == BankerTypeReadonly {
		panic("BankerTypeReadonly cannot send coins")
	}
	if b.pkgAddr != from {
		msg := `can only send coins from realm that created banker "` + b.pkgAddr + `", not "` + from + `"`
		panic(msg)
	}
	if b.bt == BankerTypeOriginSend && b.rlm == nil {
		// Unreachable for anything NewBanker built in this message: it
		// always sets rlm for OriginSend. Reaching here means the value
		// crossed a message boundary through some route the persistence
		// refusal did not catch (e.g. a realm object written before this
		// rule existed). Fail closed rather than fall through to the
		// envelope check, which alone cannot tell a revived banker from
		// a fresh one.
		panic("BankerTypeOriginSend banker is not valid outside the message that created it")
	}
	denoms, amounts := expandNative(amt)
	bankerSendCoins(uint8(b.bt), string(from), string(to), denoms, amounts)
}

func (b banker) TotalCoin(denom string) int64 {
	return bankerTotalCoin(uint8(b.bt), denom)
}

func (b banker) IssueCoin(addr address, denom string, amount int64) {
	if b.bt != BankerTypeRealmIssue {
		panic(b.bt.String() + " cannot issue coins")
	}
	assertCoinDenom(denom, b.pkgPath)
	bankerIssueCoin(uint8(b.bt), string(addr), denom, amount)
}

func (b banker) RemoveCoin(addr address, denom string, amount int64) {
	if b.bt != BankerTypeRealmIssue {
		panic(b.bt.String() + " cannot remove coins")
	}
	assertCoinDenom(denom, b.pkgPath)
	bankerRemoveCoin(uint8(b.bt), string(addr), denom, amount)
}

// assertCoinDenom checks that denom carries the prefix bound to the realm
// the banker was constructed for.
func assertCoinDenom(denom string, pkgPath string) {
	prefix := "/" + pkgPath + ":"
	if !strings.HasPrefix(denom, prefix) {
		panic("invalid denom, can only issue/remove coins with the realm's prefix: " + prefix)
	}

	baseDenom := denom[len(prefix):]
	if !isValidBaseDenom(baseDenom) {
		panic("cannot issue coins with invalid denom base name, it should start by a lowercase letter and be followed by 2-15 lowercase letters or digits")
	}
}

// check start by a lowercase letter and be followed by 2-15 lowercase letters or digits
func isValidBaseDenom(denom string) bool {
	length := len(denom)
	if length < 3 || length > 16 {
		return false
	}
	for i, c := range denom {
		switch {
		case c >= 'a' && c <= 'z',
			i > 0 && (c >= '0' && c <= '9'): // continue
		default:
			return false
		}
	}
	return true
}

// expandNative expands coins for usage within natively bound functions.
func expandNative(cz chain.Coins) (denoms []string, amounts []int64) {
	denoms = make([]string, len(cz))
	amounts = make([]int64, len(cz))
	for i, coin := range cz {
		denoms[i] = coin.Denom
		amounts[i] = coin.Amount
	}

	return denoms, amounts
}
