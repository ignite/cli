package testing

import (
	"strings"

	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/stdlibs"
)

// X_setSysParam* are un-gated counterparts to the production
// X_setSysParam* in gnovm/stdlibs/sys/params/params.go. They write
// directly into the test machine's Params, skipping
// assertSysParamsRealm. Path-based gating (this file ships only
// under tests/stdlibs and is loaded only when Testing=true) keeps
// these unreachable from realm code at chain runtime.
//
// They also skip the production keeper's WillSetParam validation
// hook by design — tests can seed deliberately invalid data to
// exercise error branches like GetValsetEffective's "valset
// corrupted" panic.

func X_setSysParamStrings(m *gno.Machine, module, submodule, name string, val []string) {
	stdlibs.GetContext(m).Params.SetStrings(prmkey(module, submodule, name), val)
}

func X_setSysParamBool(m *gno.Machine, module, submodule, name string, val bool) {
	stdlibs.GetContext(m).Params.SetBool(prmkey(module, submodule, name), val)
}

func X_setSysParamUint64(m *gno.Machine, module, submodule, name string, val uint64) {
	stdlibs.GetContext(m).Params.SetUint64(prmkey(module, submodule, name), val)
}

func X_setSysParamInt64(m *gno.Machine, module, submodule, name string, val int64) {
	stdlibs.GetContext(m).Params.SetInt64(prmkey(module, submodule, name), val)
}

// prmkey duplicates the unexported helper in
// gnovm/stdlibs/sys/params/params.go — keep in sync. Could be
// lifted to a shared internal package if a third site ever needs it.
func prmkey(module, submodule, name string) string {
	if strings.Contains(name, ":") {
		panic("invalid param name: " + name)
	}
	if submodule == "" {
		panic("submodule cannot be empty")
	}
	return module + ":" + submodule + ":" + name
}
