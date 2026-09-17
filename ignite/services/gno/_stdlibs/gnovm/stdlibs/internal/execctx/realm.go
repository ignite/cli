package execctx

import (
	"fmt"

	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
)

func GetRealm(m *gno.Machine, height int) (addr, pkgPath string) {
	// NOTE: keep in sync with test/stdlibs/std.getRealm (which keeps a
	// full legacy walk to interleave testing.SetRealm overrides).

	// Identity-chain walk (presented identities): serves every height
	// below the origin terminal — see gno.PresentedRealmAt. For
	// ordinary crosses the presented identity coincides with the
	// crossed-from context; sub-realm tokens are where the two diverge,
	// and the chain is what keeps unsafe.{Current,Previous}Realm in
	// agreement with cur/cur.Previous().
	if a, p, ok := gno.PresentedRealmAt(m, height); ok {
		return a, p
	}

	// Boundary fallback: the requested height is at/past the origin
	// terminal, or no crossing frame carries a Cur. Count crossings for
	// the stage-dependent boundary answers below. Heights below the
	// crossing count are always served by the identity chain above
	// (every WithCross frame gets its Cur at precall), so this loop no
	// longer samples per-frame realms; if that invariant ever broke,
	// the switch below fails loudly ("frame not found") rather than
	// serving context-chain answers.
	var (
		ctx     = GetContext(m)
		crosses int // track realm crosses
	)

	for i := m.NumFrames() - 1; i >= 0; i-- {
		fr := &m.Frames[i]

		// Skip over (non-realm) non-crosses.
		if !fr.IsCall() {
			continue
		}
		if !fr.WithCross {
			continue
		}

		// Sanity check
		if !fr.DidCrossing {
			panic(fmt.Sprintf(
				"call to cross(fn) did not call crossing : %s",
				fr.Func.String()))
		}

		crosses++
	}

	switch m.Stage {
	case gno.StageAdd:
		switch height {
		case crosses:
			fr := m.Frames[0]
			path := fr.LastPackage.PkgPath
			return string(gno.DerivePkgBech32Addr(path)), path
		case crosses + 1:
			return string(ctx.OriginCaller), ""
		default:
			m.PanicString("frame not found")
			return "", ""
		}
	case gno.StageRun:
		switch height {
		case crosses:
			fr := m.Frames[0]
			path := fr.LastPackage.PkgPath
			if path == "" {
				// e.g. MsgCall, cross-call a public realm function
				return string(ctx.OriginCaller), ""
			} else {
				// e.g. MsgRun, non-cross-call main()
				return string(gno.DerivePkgBech32Addr(path)), path
			}
		case crosses + 1:
			return string(ctx.OriginCaller), ""
		default:
			m.PanicString("frame not found")
			return "", ""
		}
	default:
		panic("exec kind unspecified")
	}
}

// CurrentRealm retrieves the current realm's address and pkgPath.
// It's not a native binding; but is used as a helper function here and
// elsewhere to clarify usage.
func CurrentRealm(m *gno.Machine) (address, pkgPath string) {
	return GetRealm(m, 0)
}
