package unsafe

import (
	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	testruntime "github.com/gnolang/gno/gnovm/tests/stdlibs/chain/runtime"
)

// X_getRealm delegates to the chain/runtime test stdlib's override-aware
// implementation so testing.SetRealm / RealmFrames hooks apply to
// unsafe.PreviousRealm() and unsafe.CurrentRealm() in test code.
//
// originCaller and originSend are NOT overridden here — the real
// implementations in chain/runtime/unsafe read from execctx.GetContext,
// which testing helpers already populate via SetOriginCaller /
// OriginSend setup.
func X_getRealm(m *gno.Machine, height int) (addr string, pkgPath string) {
	return testruntime.X_getRealm(m, height)
}
