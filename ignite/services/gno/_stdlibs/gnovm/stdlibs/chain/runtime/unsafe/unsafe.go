// Package unsafe quarantines the stack-walking and tx-origin
// primitives. See unsafe.gno for the security rationale.
package unsafe

import (
	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/stdlibs/internal/execctx"
)

func X_getRealm(m *gno.Machine, height int) (address, pkgPath string) {
	return execctx.GetRealm(m, height)
}

func X_originCaller(m *gno.Machine) string {
	return string(execctx.GetContext(m).OriginCaller)
}

// X_originSend is the native binding for the OriginSend envelope, a tx-origin
// primitive (see Class-2 in docs/resources/gno-security.md), so it sits with
// the other tx-origin natives rather than with the banker.
//
// The expansion below repeats banker.ExpandCoins rather than calling it, so
// this package need not import the banker it is quarantined from. The two
// write one wire shape and have to stay in step.
func X_originSend(m *gno.Machine) (denoms []string, amounts []int64) {
	ctx := execctx.GetContext(m)
	// Reading the envelope is what makes a realm payable; MsgCall rejects
	// a non-empty envelope no code ever looked at. Only the realm that was
	// actually paid may mark — a callee reading the ambient envelope must
	// not vouch for an entry realm that never looked.
	// MarkOriginSendObservedBy does that comparison; m.Realm is the realm
	// currently executing, which the VM already tracks, so this costs one
	// string compare and no stack walk.
	//
	// Marking for an empty envelope is harmless and kept: the flag records
	// observation, not amount, and the keeper only consults it when send is
	// non-empty.
	var realmPath string
	if m.Realm != nil {
		realmPath = m.Realm.Path
	}
	ctx.MarkOriginSendObservedBy(realmPath)
	os := ctx.OriginSend
	denoms = make([]string, len(os))
	amounts = make([]int64, len(os))
	for i, coin := range os {
		denoms[i] = coin.Denom
		amounts[i] = coin.Amount
	}
	return denoms, amounts
}
