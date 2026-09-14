// Package unsafe quarantines the stack-walking and tx-origin
// primitives that are footgun-prone for caller- and payment-
// authentication. Forcing call sites to read unsafe.PreviousRealm(),
// unsafe.CurrentRealm(), unsafe.OriginCaller(), unsafe.OriginSend()
// makes the bug class greppable.
//
// Prefer cur.Previous() / cur.PkgPath() / cur.Address() on a threaded
// `cur realm` parameter for caller-identity. Reserve the symbols here
// for cases where you intentionally want tx-level identity (event
// emission, fee attribution) and pair them with
// runtime.AssertOriginCall() so misuse panics.
package unsafe

import (
	"chain"
	"chain/runtime"
)

// PreviousRealm returns the realm-context of the caller frame.
//
// SECURITY: stack-walks. In a non-crossing helper it identifies
// whichever realm was outermost-crossing in the chain, NOT the
// immediate caller. The idiomatic substitute is cur.Previous() on
// a threaded `cur realm` parameter, which is statically scoped to
// a crossing function and cannot lie.
func PreviousRealm() runtime.Realm {
	addr, path := getRealm(1)
	return runtime.NewRealm(address(addr), path)
}

// CurrentRealm returns the current realm-context.
//
// SECURITY: stack-walks. Inside a /p/ helper called via borrow rule #2 borrow
// it returns the BORROWER's realm — using equality against an admin
// address here ships a delegation bug. The idiomatic substitute, when
// you have cur, is cur.PkgPath() / cur.Address(); when you don't,
// accept `_ int, rlm realm` from the caller.
func CurrentRealm() runtime.Realm {
	addr, path := getRealm(0)
	return runtime.NewRealm(address(addr), path)
}

// OriginCaller returns the EOA that signed the current transaction.
//
// SECURITY: gno's tx.origin. A malicious realm called by the EOA can
// act AS the EOA against any other realm relying on OriginCaller as
// auth — identical to Solidity's well-known tx.origin phishing class.
// Prefer cur.Previous() for caller-auth; reserve OriginCaller for
// intentionally EOA-scoped uses (event emission, fee attribution) and
// pair with runtime.AssertOriginCall() so misuse panics.
func OriginCaller() address {
	return address(originCaller())
}

// OriginSend returns the coin envelope attached at the chain root.
//
// SECURITY: OriginCaller's payment twin. The envelope is shared across
// every realm in the call chain, so a malicious intermediate realm can
// consume it after a target's payment-verification check passes
// (TOCTOU/phishing). When used for payment verification, pair with
// runtime.AssertOriginCall() AND runtime.PreviousRealm().IsUserCall()
// (NOT IsUser(), which accepts maketx-run ephemeral realms that can
// pre-consume the envelope).
func OriginSend() chain.Coins {
	den, amt := originSend()
	coins := make(chain.Coins, len(den))
	for i := range coins {
		coins[i] = chain.Coin{Denom: den[i], Amount: amt[i]}
	}
	return coins
}

// Native bindings.
func getRealm(height int) (address string, pkgPath string)
func originCaller() string
func originSend() (denoms []string, amounts []int64)
