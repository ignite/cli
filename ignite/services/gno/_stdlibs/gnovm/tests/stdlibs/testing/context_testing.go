package testing

import (
	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/stdlibs/chain/banker"
	"github.com/gnolang/gno/gnovm/tests/stdlibs/chain/runtime"
	"github.com/gnolang/gno/tm2/pkg/crypto"
)

func X_getContext(m *gno.Machine) (
	originCaller string,
	origSendDenoms []string, origSendAmounts []int64,
	origSpendDenoms []string, origSpendAmounts []int64,
	chainID string,
	height int64,
	timeUnix int64, timeNano int64,
) {
	ctx := m.Context.(*runtime.TestExecContext)

	originCaller = ctx.OriginCaller.String()

	for _, coin := range ctx.OriginSend {
		origSendDenoms = append(origSendDenoms, coin.Denom)
		origSendAmounts = append(origSendAmounts, coin.Amount)
	}

	for _, coin := range *ctx.OriginSendSpent {
		origSpendDenoms = append(origSpendDenoms, coin.Denom)
		origSpendAmounts = append(origSpendAmounts, coin.Amount)
	}

	chainID = ctx.ChainID
	height = ctx.Height
	timeUnix = ctx.Timestamp
	timeNano = ctx.TimestampNano
	return
}

func X_setContext(
	m *gno.Machine,
	originCaller string,
	currRealmAddr string, currRealmPkgPath string,
	origSendDenoms []string, origSendAmounts []int64,
	origSpendDenoms []string, origSpendAmounts []int64,
	chainID string,
	height int64,
	timeUnix int64, timeNano int64,
) {
	ctx := m.Context.(*runtime.TestExecContext)

	ctx.ChainID = chainID
	ctx.Height = height
	ctx.Timestamp = timeUnix
	ctx.TimestampNano = timeNano
	ctx.OriginCaller = crypto.Bech32Address(originCaller)

	if currRealmAddr != "" {
		// Associate the given Realm with the caller's frame.
		var frameIdx int
		// NOTE: the frames are different from when calling testing.SetRealm (has been refactored to this code)
		//
		// When calling this function from Gno, the 3 top frames are the following:
		// #7: [FRAME FUNC:setContext RECV:(undefined) (15 args) 11/3/0/6/4 LASTPKG:testing ...]
		// #6: [FRAME FUNC:SetContext RECV:(undefined) (1 args) 8/2/0/4/3 LASTPKG:testing ...]
		// #5: [FRAME FUNC:SetRealm RECV:(undefined) (1 args) 5/1/0/2/2 LASTPKG:gno.land/r/demo/groups ...]
		// We want to set the Realm of the frame where testing.SetRealm is being called, hence -3-1.
		for i := m.NumFrames() - 4; i >= 0; i-- {
			// Must be a frame from calling a function.
			if fr := m.Frames[i]; fr.Func != nil && fr.Func.PkgPath != "testing" {
				frameIdx = i
				break
			}
		}

		m.Frames[frameIdx].TestOverridden = true // in case frame gets popped
		ctx.RealmFrames[frameIdx] = runtime.RealmOverride{
			Addr:    crypto.Bech32Address(currRealmAddr),
			PkgPath: currRealmPkgPath,
		}
		// Also mutate the captured `cur` value for this frame so that
		// reads through the uverse `realm` handle reflect the override,
		// matching what runtime.{Current,Previous}Realm() returns after
		// the X_getRealm walk:
		//
		//   - addr/pkgPath: overwrite with override values (CurrentRealm
		//     parity).
		//   - prev: depends on the override shape.
		//       * UserRealm override (pkgPath==""): set prev to a true
		//         nil — there's no "previous" beyond an EOA caller,
		//         matching runtime.PreviousRealm()'s walk panic.
		//       * CodeRealm override (pkgPath!=""): set prev to a fresh
		//         realm carrying the pre-override addr/pkgPath. That's
		//         the realm X_getRealm surfaces as "previous" of the
		//         override frame.
		fr := &m.Frames[frameIdx]
		if pv, ok := fr.Cur.V.(gno.PointerValue); ok && pv.TV != nil {
			if sv, ok := pv.TV.V.(*gno.StructValue); ok && len(sv.Fields) >= 3 {
				sv.Fields[0].V = gno.StringValue(currRealmAddr)
				sv.Fields[1].V = gno.StringValue(currRealmPkgPath)
				if currRealmPkgPath == "" {
					// UserRealm override — no previous.
					sv.Fields[2] = gno.TypedValue{}
				} else {
					// CodeRealm override — prev is the frame's
					// underlying package realm (what X_getRealm
					// surfaces as PreviousRealm of an override
					// frame: m.Frames[0].LastPackage.PkgPath in
					// the filetest case, or the frame's func
					// PkgPath more generally). Use the frame's
					// function package as the stable identity —
					// it doesn't shift across successive overrides.
					pkgPath := ""
					if fr.Func != nil {
						pkgPath = fr.Func.PkgPath
					}
					addr := ""
					if pkgPath != "" {
						addr = string(gno.DerivePkgBech32Addr(pkgPath))
					}
					sv.Fields[2] = gno.BuildOverridePrevField(addr, pkgPath)
				}
			}
		}
	}

	ctx.OriginSend = banker.CompactCoins(origSendDenoms, origSendAmounts)
	coins := banker.CompactCoins(origSpendDenoms, origSpendAmounts)
	ctx.OriginSendSpent = &coins
	// Keep OriginSendRecipient consistent with the envelope being set.
	//
	// On chain the -send envelope always lands on the entry realm, so a
	// test that re-points the current realm to a *code* realm is
	// modelling a different entry realm and must move the recipient with
	// it. Otherwise BankerTypeOriginSend sends fail against a stale
	// recipient.
	//
	// A *user* realm override (testing.SetRealm(testing.NewUserRealm))
	// means "an EOA is calling", not "this EOA received the envelope".
	// A user address can never equal any realm's address, so moving the
	// recipient there would point it at something no realm can match and
	// break every payment test. Leave it on the package under test, which
	// gnovm/pkg/test.Context already seeded.
	//
	// testing.SetOriginSend passes an empty currRealmAddr (GetContext
	// does not round-trip CurrentRealm), so it leaves the recipient
	// alone either way.
	//
	// Move both spellings together. OriginSendRecipient (address) and
	// OriginSendRecipientPath (package path) name the same realm and must
	// never disagree: the banker gate reads the address, the payable check
	// reads the path. Setting one without the other leaves the harness
	// describing two different realms.
	if currRealmAddr != "" && currRealmPkgPath != "" {
		ctx.OriginSendRecipient = crypto.Bech32Address(currRealmAddr)
		ctx.OriginSendRecipientPath = currRealmPkgPath
	}

	m.Context = ctx
}

func X_testIssueCoins(m *gno.Machine, addr string, denom []string, amt []int64) {
	ctx := m.Context.(*runtime.TestExecContext)
	banker := ctx.Banker
	for i := range denom {
		banker.IssueCoin(crypto.Bech32Address(addr), denom[i], amt[i])
	}
}

// X_makeRealm builds a uverse realm value with the given (addr, pkgPath,
// prev) tuple. Tests use this to construct cur values explicitly when
// the SetRealm/SetCodeRealm composition semantics don't match the
// scenario being tested — e.g. to express "alice EOA crossed into
// r/foo, cur.Previous() == alice realm" without relying on chained
// SetRealm calls (which overwrite prev with the test pkg's own addr).
func X_makeRealm(m *gno.Machine, addr, pkgPath string, prev gno.TypedValue) gno.TypedValue {
	return gno.MakeRealmValue(m.Alloc, addr, pkgPath, prev)
}

// X_originRealm returns the EOA-origin realm value (addr=OriginCaller,
// pkgPath="", prev=truly-nil). Useful as the seed prev when assembling
// explicit cur chains in tests.
func X_originRealm(m *gno.Machine) gno.TypedValue {
	return gno.OriginRealmTV()
}

func X_newRealm(m *gno.Machine, addr, pkgPath string) gno.TypedValue {
	rlmType := m.Store.GetType("chain/runtime.Realm")
	return gno.TypedValue{
		// testing imports chain/runtime, so this type is always available.
		T: rlmType,
		V: m.Alloc.NewStructWithFields(
			rlmType,
			// addr address
			gno.TypedValue{T: m.Store.GetType(".uverse.address"), V: gno.StringValue(addr)},
			// pkgPath string
			gno.TypedValue{T: gno.StringType, V: gno.StringValue(pkgPath)},
		),
	}
}

func X_isRealm(m *gno.Machine, pkgPath string) bool {
	return gno.IsRealmPath(pkgPath)
}
