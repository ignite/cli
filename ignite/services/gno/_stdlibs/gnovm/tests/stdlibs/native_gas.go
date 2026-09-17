package stdlibs

import (
	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
)

// Test-stdlib gas entries. These natives only run in test mode (via
// teststdlibs.NativeResolver) and never on the live chain. Registered
// here purely so chargeNativeGas in gnolang doesn't panic via the
// uncalibrated-stdlib forcing function when TestFiles or other test
// harnesses exercise them.
//
// All entries charge ZERO gas. Test runs aren't gas-budget-bound, the
// numbers wouldn't matter even if they were, and benching this set
// properly through gnovm/cmd/calibrate/ is a separate exercise. If a
// future test setup ever cares about gas attribution to a test-only
// native, calibrate that specific native and update its row.

var testNativeFns = [][2]string{
	// chain/runtime test overrides intentionally OMITTED here —
	// gnostdlibs already registered AssertOriginCall and getRealm at
	// init, and the production gas values cover the test variant fine.
	// Duplicates would panic via RegisterNativeGas.

	// fmt — typedvalue inspection helpers used by fmt.Println/Sprintf.
	{"fmt", "typeString"},
	{"fmt", "valueOfInternal"},
	{"fmt", "getAddr"},
	{"fmt", "getPtrElem"},
	{"fmt", "mapKeyValues"},
	{"fmt", "arrayIndex"},
	{"fmt", "fieldByIndex"},
	{"fmt", "asByteSlice"},

	// os — write/sleep test helpers.
	{"os", "write"},
	{"os", "sleep"},

	// runtime — GC / MemStats test helpers.
	{"runtime", "GC"},
	{"runtime", "MemStats"},

	// testing — context / assertions / regex.
	{"testing", "getContext"},
	{"testing", "isRealm"},
	{"testing", "makeRealm"},
	{"testing", "matchString"},
	{"testing", "newRealm"},
	{"testing", "originRealm"},
	{"testing", "recoverWithStacktrace"},
	{"testing", "setContext"},
	{"testing", "setSysParamBool"},
	{"testing", "setSysParamInt64"},
	{"testing", "setSysParamStrings"},
	{"testing", "setSysParamUint64"},
	{"testing", "testIssueCoins"},
	{"testing", "unixNano"},

	// unicode — read-only character class checks.
	{"unicode", "IsGraphic"},
	{"unicode", "IsPrint"},
	{"unicode", "IsUpper"},
	{"unicode", "SimpleFold"},
}

func init() {
	for _, fn := range testNativeFns {
		gno.RegisterNativeGas(fn[0], gno.Name(fn[1]), &gno.NativeGasInfo{
			SlopeIdx: -1,
		})
	}
}
