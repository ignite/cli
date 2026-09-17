// Package stdlibs provides supplemental stdlibs for the testing environment.
package stdlibs

import (
	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/stdlibs"
	"github.com/gnolang/gno/gnovm/tests/stdlibs/chain/runtime"
)

//go:generate go run github.com/gnolang/gno/misc/genstd -skip-init-order

func NativeResolver(pkgPath string, name gno.Name) func(*gno.Machine) {
	for _, nf := range nativeFuncs {
		if nf.gnoPkg == pkgPath && name == nf.gnoFunc {
			return nf.f
		}
	}
	return stdlibs.NativeResolver(pkgPath, name)
}

type TestExecContext = runtime.TestExecContext

func HasNativePkg(pkgPath string) bool {
	for _, nf := range nativeFuncs {
		if nf.gnoPkg == pkgPath {
			return true
		}
	}
	return stdlibs.HasNativePkg(pkgPath)
}
