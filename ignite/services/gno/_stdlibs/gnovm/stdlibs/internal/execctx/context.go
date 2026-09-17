package execctx

import (
	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/sdk"
	"github.com/gnolang/gno/tm2/pkg/std"
)

type BankerInterface interface {
	GetCoins(addr crypto.Bech32Address) (dst std.Coins)
	// GetCoin reads one denom. GetCoins costs O(denoms held); this does not.
	GetCoin(addr crypto.Bech32Address, denom string) int64
	SendCoins(from, to crypto.Bech32Address, amt std.Coins)
	TotalCoin(denom string) int64
	IssueCoin(addr crypto.Bech32Address, denom string, amount int64)
	RemoveCoin(addr crypto.Bech32Address, denom string, amount int64)
}

type ParamsInterface interface {
	SetString(key, val string)
	SetBool(key string, val bool)
	SetInt64(key string, val int64)
	SetUint64(key string, val uint64)
	SetBytes(key string, val []byte)
	SetStrings(key string, val []string)
	UpdateStrings(key string, val []string, add bool)
	// GetXxx writes the stored value (if any) into *ptr and returns
	// whether the key existed. A return of false leaves *ptr at its
	// zero value, distinguishing "never set" from "set to zero" —
	// which the in-memory backed types alone cannot.
	GetString(key string, ptr *string) bool
	GetBool(key string, ptr *bool) bool
	GetInt64(key string, ptr *int64) bool
	GetUint64(key string, ptr *uint64) bool
	GetBytes(key string, ptr *[]byte) bool
	GetStrings(key string, ptr *[]string) bool
}

type ExecContext struct {
	ChainID         string
	ChainDomain     string
	Height          int64
	Timestamp       int64 // seconds
	TimestampNano   int64 // nanoseconds, only used for testing.
	OriginCaller    crypto.Bech32Address
	OriginSend      std.Coins
	OriginSendSpent *std.Coins // mutable
	// OriginSendRecipient is the single address that OriginSend was
	// actually credited to for this message — the entry realm's address
	// (for MsgRun, the caller's, since /e/<addr>/run derives to it).
	// BankerTypeOriginSend spending is gated on it: a banker is a plain
	// persistable value, so its construction-time authority check
	// (rlm.Previous().IsUserCall(), see chain/banker/banker.gno) can be
	// separated arbitrarily in time from its use. Without this field the
	// use-time limit would re-arm against the ambient envelope of any
	// later message, one the banker's realm never received.
	//
	// The zero value means "no envelope was delivered in this message",
	// which is fail-closed: no BankerTypeOriginSend send can succeed.
	// That is the correct value for envelope-less contexts (queries,
	// internal realm callouts).
	OriginSendRecipient crypto.Bech32Address
	// OriginSendRecipientPath is the same realm as OriginSendRecipient,
	// written as a package path instead of an address. It exists so the
	// payable check can tell "the realm that got paid looked at the
	// envelope" apart from "some other realm looked at it" using the
	// realm the VM is already tracking (Machine.Realm), with no stack
	// walk and no change to the gas table.
	//
	// Empty in contexts that cannot carry a send, which is fail-closed:
	// no realm has an empty path, so nothing matches and nothing is
	// marked as observed.
	OriginSendRecipientPath string
	// OriginSendObserved records whether the realm that was paid ever made
	// the send-envelope observable — set via MarkOriginSendObservedBy, or
	// MarkOriginSendObserved on the banker path. MsgCall uses
	// it to reject a non-empty envelope that the callee never looked at,
	// which is what would otherwise silently strand coins in a realm that
	// has no notion of being paid.
	//
	// "The code read the envelope" is the operational definition of
	// payable. It cannot be decided statically: Gno has interface dispatch
	// and function values, so whether a call path reaches a read is
	// undecidable in general, and any conservative approximation marks
	// most of the ecosystem payable.
	//
	// This is a safety net against stranded funds, NOT an access-control
	// boundary: a realm opts in simply by reading the envelope and
	// discarding the result. Nil in contexts that cannot carry a send.
	//
	// Must stay a pointer. The Mark* methods below have value receivers and
	// GetContext returns a copy, so they can only write through this
	// indirection. Change it to a plain bool and they become silent no-ops,
	// which would fail every MsgCall that carries coins.
	OriginSendObserved *bool // mutable
	Banker             BankerInterface
	Params             ParamsInterface
	EventLogger        *sdk.EventLogger
	SessionAccount     std.DelegatedAccount // nil for master-key txs
}

// MarkOriginSendObservedBy records that the realm at realmPath made the
// message's send-envelope observable — but only if that realm is the one
// the envelope was credited to.
//
// The check is the point. Any realm can read the ambient envelope at any
// call depth, so without it a realm deeper in the chain could satisfy the
// payable check on behalf of an entry realm that never looked. That is
// exactly the stranding this is meant to catch: the coins sit in the entry
// realm, which has no idea it was paid.
//
// Callers pass Machine.Realm.Path — the realm the VM is currently
// executing in, which it already tracks, so this costs one string compare
// and no stack walk.
//
// Nil-safe and empty-safe: contexts that cannot carry a send leave the
// field unset, and no realm has an empty path, so nothing is marked.
//
// Use MarkOriginSendObserved instead when the caller has already proved
// the envelope belongs to the realm some other way.
func (e ExecContext) MarkOriginSendObservedBy(realmPath string) {
	if e.OriginSendObserved == nil {
		return
	}
	if realmPath == "" || realmPath != e.OriginSendRecipientPath {
		return
	}
	*e.OriginSendObserved = true
}

// MarkOriginSendObserved marks the envelope observed without checking who
// is running.
//
// Only call this once you have already proved the envelope belongs to the
// caller some other way. The one such place is the BankerTypeOriginSend
// send path, which first checks that the spending address is the address
// the envelope was credited to. Once that holds, the realm that owns the
// banker is by definition the realm that was paid, even if it handed the
// banker to another realm to spend within this message — building a
// banker over your own envelope is itself proof you noticed the payment.
//
// Nil-safe.
func (e ExecContext) MarkOriginSendObserved() {
	if e.OriginSendObserved != nil {
		*e.OriginSendObserved = true
	}
}

// GetContext returns the execution context.
// This is used to allow extending the exec context using interfaces,
// for instance when testing.
func (e ExecContext) GetExecContext() ExecContext {
	return e
}

var _ ExecContexter = ExecContext{}

// ExecContexter is a type capable of returning the parent [ExecContext]. When
// using these standard libraries, m.Context should always implement this
// interface. This can be obtained by embedding [ExecContext].
type ExecContexter interface {
	GetExecContext() ExecContext
}

// NOTE: In order to make this work by simply embedding ExecContext in another
// context (like TestExecContext), the method needs to be something other than
// the field name.

// GetContext returns the context from the Gno machine.
func GetContext(m *gno.Machine) ExecContext {
	return m.Context.(ExecContexter).GetExecContext()
}

// Wire the per-tx OriginCaller into gnolang so it can build the per-tx
// origin realm value (the EOA-shaped value at the chain root used by
// captured `cur.Previous()`). This must match what runtime.PreviousRealm()
// surfaces in the same context — (OriginCaller, "") at the EOA boundary.
func init() {
	gno.OriginCallerExtractor = func(ctx any) string {
		if ec, ok := ctx.(ExecContexter); ok {
			return string(ec.GetExecContext().OriginCaller)
		}
		return ""
	}
}
