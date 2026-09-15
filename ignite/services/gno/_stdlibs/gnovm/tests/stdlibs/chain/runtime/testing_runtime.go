package runtime

import (
	"fmt"
	"strings"

	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/stdlibs"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/overflow"
	tm2std "github.com/gnolang/gno/tm2/pkg/std"
)

// TestExecContext is the testing extension of the exec context.
type TestExecContext struct {
	stdlibs.ExecContext

	// These are used to set up the result of CurrentRealm() and PreviousRealm().
	RealmFrames map[int]RealmOverride
}

var _ stdlibs.ExecContexter = &TestExecContext{}

type RealmOverride struct {
	Addr    crypto.Bech32Address
	PkgPath string
}

func AssertOriginCall(m *gno.Machine) {
	if !isOriginCall(m) {
		m.Panic(typedString("invalid non-origin call"))
		return
	}
}

func typedString(s gno.StringValue) gno.TypedValue {
	tv := gno.TypedValue{T: gno.StringType}
	tv.SetString(s)
	return tv
}

func isOriginCall(m *gno.Machine) bool {
	tname := m.Frames[0].Func.Name
	// Count only actual function call frames (excludes closures and
	// control-flow basic frames like for/range/switch).
	callFrames := m.NumCallFrames()
	switch tname {
	case "main": // test is a _filetest
		// Non-closure frames expected:
		// 0. main
		// 1. $RealmFuncName
		// 2. runtime.AssertOriginCall
		return callFrames == 3
	case "RunTest", "runTest_cur": // _test, with or without (cur realm, t *testing.T)
		// Non-closure frames expected:
		// 0. testing.RunTest / runTest_cur
		// 1. tRunner / tRunner_cur
		// 2. $TestFuncName / $TestFuncName_cur
		// 3. $RealmFuncName
		// 4. runtime.AssertOriginCall
		return callFrames == 5
	}
	// support init() in _filetest
	// XXX do we need to distinguish from 'runtest'/_test?
	// XXX pretty hacky even if not.
	if strings.HasPrefix(string(tname), "init.") {
		return callFrames == 3
	}
	panic("unable to determine if test is a _test or a _filetest")
}

// anyLiveOverride reports whether any on-stack frame carries a live
// testing.SetRealm override (map entry AND the frame's TestOverridden
// flag, matching getOverride's validity rule). Stale map entries left
// by popped frames don't count — they must not disable the identity
// chain for the rest of the machine.
func anyLiveOverride(m *gno.Machine, ctx *TestExecContext) bool {
	if len(ctx.RealmFrames) == 0 {
		return false
	}
	for i := m.NumFrames() - 1; i >= 0; i-- {
		if !m.Frames[i].TestOverridden {
			continue
		}
		if _, ok := ctx.RealmFrames[i]; ok {
			return true
		}
	}
	return false
}

func getOverride(m *gno.Machine, i int) (RealmOverride, bool) {
	fr := &m.Frames[i]
	ctx := m.Context.(*TestExecContext)
	override, overridden := ctx.RealmFrames[i]
	if overridden && !fr.TestOverridden {
		return RealmOverride{}, false // override was replaced
	}
	return override, overridden
}

func X_getRealm(m *gno.Machine, height int) (addr string, pkgPath string) {
	// NOTE: keep in sync with stdlibs/std.getRealm

	ctx := m.Context.(*TestExecContext)

	// Identity-chain walk (see execctx.GetRealm / gno.PresentedRealmAt):
	// serves presented identities, including sub-realm tokens. Applied
	// only when no live testing.SetRealm override is on the stack:
	// overrides are frame-index-keyed and interleave with the legacy
	// walk below, and UserRealm overrides on non-crossing frames are
	// invisible to the chain — the gate is load-bearing.
	if !anyLiveOverride(m, ctx) {
		if a, p, ok := gno.PresentedRealmAt(m, height); ok {
			return a, p
		}
	}

	var (
		lfr     = m.LastFrame() // last call frame
		crosses int             // track realm crosses
	)

	for i := m.NumFrames() - 1; i >= 0; i-- {
		fr := &m.Frames[i]

		// Skip over (non-realm) non-crosses.
		override, overridden := getOverride(m, i)
		if overridden {
			if override.PkgPath == "" && crosses < height {
				m.Panic(typedString("frame not found: cannot seek beyond origin caller override"))
			}
		}
		if !overridden {
			if !fr.IsCall() {
				continue
			}
			if !fr.WithCross {
				lfr = fr
				continue
			}
		}

		// Sanity check XXX move check elsewhere
		if !overridden {
			if !fr.DidCrossing {
				panic(fmt.Sprintf(
					"cross(fn) but fn didn't call crossing(): %s.%s",
					fr.Func.PkgPath,
					fr.Func.String()))
			}
		}

		crosses++
		if crosses > height {
			if overridden {
				caller, pkgPath := override.Addr, override.PkgPath
				return string(caller), pkgPath
			} else {
				currlm := lfr.LastRealm
				caller, rlmPath := gno.DerivePkgBech32Addr(currlm.Path), currlm.Path
				return string(caller), rlmPath
			}
		}
		lfr = fr
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
			m.Panic(typedString("frame not found"))
			return "", ""
		}
	case gno.StageRun:
		switch height {
		case crosses:
			fr := m.Frames[0]
			path := fr.LastPackage.PkgPath
			if path == "" {
				// Not sure what would cause this.
				panic("should not happen")
			} else {
				// e.g. TestFoo(t *testing.Test) in *_test.gno
				// or main() in *_filetest.gno
				return string(gno.DerivePkgBech32Addr(path)), path
			}
		case crosses + 1:
			return string(ctx.OriginCaller), ""
		default:
			m.Panic(typedString("frame not found"))
			return "", ""
		}
	default:
		panic("exec kind unspecified")
	}
}

// TestBanker is a banker that can be used as a mock banker in test contexts.
type TestBanker struct {
	CoinTable map[crypto.Bech32Address]tm2std.Coins
}

var _ stdlibs.BankerInterface = &TestBanker{}

// GetCoins implements the Banker interface.
func (tb *TestBanker) GetCoins(addr crypto.Bech32Address) (dst tm2std.Coins) {
	return tb.CoinTable[addr]
}

// GetCoin implements the Banker interface.
func (tb *TestBanker) GetCoin(addr crypto.Bech32Address, denom string) int64 {
	// AmountOf validates the denom before reading anything, so a malformed one
	// panics here as it does on chain, where SDKBanker.GetCoin checks explicitly
	// because the keeper reaches a store key without going through Coins.
	return tb.CoinTable[addr].AmountOf(denom)
}

// SendCoins implements the Banker interface.
func (tb *TestBanker) SendCoins(from, to crypto.Bech32Address, amt tm2std.Coins) {
	fcoins, fexists := tb.CoinTable[from]
	if !fexists {
		panic(fmt.Sprintf(
			"source address %s does not exist",
			from.String()))
	}
	if !fcoins.IsAllGTE(amt) {
		panic(fmt.Sprintf(
			"source address %s has %s; cannot send %s",
			from.String(), fcoins, amt))
	}
	// First, subtract from 'from'.
	frest := fcoins.Sub(amt)
	tb.CoinTable[from] = frest
	// Second, add to 'to'.
	// NOTE: even works when from==to, due to 2-step isolation.
	tcoins := tb.CoinTable[to]
	tsum := tcoins.Add(amt)
	tb.CoinTable[to] = tsum
}

// TotalCoin implements the Banker interface.
func (tb *TestBanker) TotalCoin(denom string) int64 {
	// Summed from the table rather than kept as a counter: the test banker has no
	// mint/burn asymmetry, so the sum *is* the supply. Total, like the chain — a
	// denom nobody holds has zero supply, and a malformed one panics as SDKBanker's
	// does, so `gno test` and production agree.
	//
	// Checked up front rather than left to AmountOf below, which is only reached
	// once per held denom: an empty table would skip the loop entirely and report
	// a malformed denom as zero, where the chain panics.
	if err := tm2std.ValidateDenom(denom); err != nil {
		panic(err)
	}
	// Overflow-checked, so this cannot report a wrapped negative where the chain
	// would refuse the mint outright — the supply cap is exactly what the chain's
	// counter enforces, and `gno test` must not disagree with it.
	var total int64
	for _, coins := range tb.CoinTable {
		sum, ok := overflow.Add(total, coins.AmountOf(denom))
		if !ok {
			panic("total supply of " + denom + " overflows int64")
		}
		total = sum
	}
	return total
}

// IssueCoin implements the Banker interface.
//
// Deliberately does not enforce the chain's per-denom supply cap: there is no counter
// here, so an aggregate past MaxInt64 is reachable under `gno test` where MintCoins
// would refuse the second mint. TotalCoin above panics rather than reporting a wrapped
// total, so the state cannot be read as a number, and the divergence is in the
// permissive direction. See the "Known limitation" note in
// tm2/adr/pr6034_coin_supply.md before changing this — an overflow check here alone
// refuses at the right point for the wrong reason and still does not model burn.
func (tb *TestBanker) IssueCoin(addr crypto.Bech32Address, denom string, amt int64) {
	coins := tb.CoinTable[addr]
	sum := coins.Add(tm2std.Coins{{Denom: denom, Amount: amt}})
	tb.CoinTable[addr] = sum
}

// RemoveCoin implements the Banker interface.
func (tb *TestBanker) RemoveCoin(addr crypto.Bech32Address, denom string, amt int64) {
	coins := tb.CoinTable[addr]
	rest := coins.Sub(tm2std.Coins{{Denom: denom, Amount: amt}})
	tb.CoinTable[addr] = rest
}

func X_testIssueCoins(m *gno.Machine, addr string, denom []string, amt []int64) {
	ctx := m.Context.(*TestExecContext)
	banker := ctx.Banker
	for i := range denom {
		banker.IssueCoin(crypto.Bech32Address(addr), denom[i], amt[i])
	}
}
