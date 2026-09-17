package stdlibs

import (
	"fmt"

	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
)

// Per-native gas charging. Calibrated values come from
// gnovm/cmd/calibrate/native_bench_test.go (and its machine-harness
// companion); the python fitter at gnovm/cmd/calibrate/gen_native_table.py
// emits the literal table below.
//
// Registration is at init time into a package global in gnolang
// (RegisterNativeGas). Mirrors how OpCPU* constants live in gnolang —
// no per-Machine resolver plumbing.
//
// Calibration is end-to-end through doOpCallNativeBody (Gno↔Go reflect +
// X_ work + return push); no separate dispatch overhead constant.
//
// Two-slope (Slope2/PostSlope2) is supported by the runtime for natives
// whose cost depends on two independent dimensions. Only
// crypto/merkle.innerHash uses it today: it hashes two unbounded byte
// slices, so each one carries its own per-byte slope. Every other row
// uses a single slope — for the store-backed natives the per-byte work
// happens inside the metered KVStore (gctx), so the empirical per-byte
// CPU cost in the dispatcher is negligible.

// Re-export the gnolang SizeKind constants for readable table literals.
const (
	SizeFlat            = gno.SizeFlat
	SizeLenBytes        = gno.SizeLenBytes
	SizeLenString       = gno.SizeLenString
	SizeLenSlice        = gno.SizeLenSlice
	SizeNumCallFrames   = gno.SizeNumCallFrames
	SizeReturnLen       = gno.SizeReturnLen
	SizeSliceTotalBytes = gno.SizeSliceTotalBytes
	SizeModExpWork      = gno.SizeModExpWork
)

// nativeGasEntry is the on-disk shape of a row, copied into a
// gno.NativeGasInfo at init time.
type nativeGasEntry struct {
	Pkg, Fn   string
	Base      int64
	Slope     int64
	SlopeIdx  int8
	SlopeKind gno.NativeGasSize

	// Optional second pre-call slope for natives whose cost depends on
	// two independent dimensions (e.g. count and total inner bytes).
	Slope2     int64
	Slope2Idx  int8
	Slope2Kind gno.NativeGasSize

	// Post-call charge (zero = none). PostSlopeIdx is the stack offset
	// from top (1 = last-pushed return); PostSlopeKind must be a
	// SizeReturn* kind (or SizeSliceTotalBytes for slice-of-string).
	PostBase      int64
	PostSlope     int64
	PostSlopeIdx  int8
	PostSlopeKind gno.NativeGasSize

	// Optional second post-call slope, summed independently.
	PostSlope2     int64
	PostSlope2Idx  int8
	PostSlope2Kind gno.NativeGasSize
}

// Calibrated on Apple M2 ARM64 (NOT the reference Xeon 8168 — re-run
// gnovm/cmd/calibrate before any consensus-relevant deployment). 1 gas
// = 1 ns. Slope is ns per 1024 units of N. R² > 0.93 for all linear fits.
//
// Values come from gen_native_table.py over native_bench_output.txt.
// Two caveats on that provenance, both worth knowing before a re-run:
//
//   - The checked-in snapshot is stale, and not because of this table.
//     native_gas_table.go.txt holds 46 rows against production's 66, and
//     native_bench_output.txt contains no bench lines for any of the IBC
//     crypto natives (bn254, cometbls, keccak256, merkle, modexp). Those
//     rows came from a separate run (see below) and cannot be reproduced
//     from the committed input at all, so "re-running the fitter
//     reproduces this table verbatim" holds only for the 46 rows the
//     snapshot covers.
//   - The crypto/modexp row is additionally not a fitter output even given
//     its bench data: its Slope is an upper bound over the bench grid
//     rather than the least-squares coefficient, because that native's
//     cost is a product and a central fit underprices half the grid. See
//     the row's own comment.
//
// The 2D bench-grid extension
// (slice natives benched at multiple per-element byte sizes) confirmed
// the per-byte CPU slope is below noise for every native shipping
// today, so the table stays single-slope; the schema fields support
// future natives that genuinely scale on both dimensions.
//
// 72 entries — exhaustive coverage of gnovm/stdlibs/generated.go.
// The trailing 10 IBC-crypto entries (crypto/bn254, crypto/cometbls,
// crypto/keccak256, crypto/merkle, crypto/modexp) are draft fits measured
// on Intel Xeon Silver 4114; the chain/markdown rows and the rest are on
// Apple M2. The whole table must be regenerated on the reference Xeon 8168
// before any consensus-relevant deployment; the IBC rows are flagged
// "draft" in their trailing comment to make that obvious.
//
// The six chain/params Get* rows are a further exception, and a different one:
// they are not fitted at all. native_bench_output.txt holds no samples for
// them, so re-running the fitter drops those rows rather than reproducing
// them. Each Base is copied from the matching Set*, and GetBytes and
// GetStrings additionally borrow a PostSlope from the sys/params getter of the
// same shape. The benchmarks that would replace them exist
// (BenchmarkNative_Params_Get* in
// gnovm/cmd/calibrate/native_machine_bench_test.go).
//
// The borrowed Base is very unlikely to undercharge: a setter writes rather
// than reads and also runs recordParamsDelta for storage-deposit accounting,
// so it does strictly more work than the getter lending its price. That
// argument covers the flat part only. The two borrowed PostSlopes come from
// measured sys/params getters rather than from a setter, and the four scalar
// getters carry no length term at all.
var calibratedNativeGas = []nativeGasEntry{
	{Pkg: "crypto/sha256", Fn: "sum256", Base: 226, Slope: 8906, SlopeIdx: 0, SlopeKind: SizeLenBytes},                                                         // fit base=226.3ns slope=8.6969ns/N (=8906/1024) R²=1.000
	{Pkg: "crypto/ed25519", Fn: "verify", Base: 56534, Slope: 8975, SlopeIdx: 1, SlopeKind: SizeLenBytes},                                                      // fit base=56534.0ns slope=8.7645ns/N (=8975/1024) R²=0.991
	{Pkg: "chain", Fn: "packageAddress", Base: 552, Slope: 15201, SlopeIdx: 0, SlopeKind: SizeLenString},                                                       // fit base=552.1ns slope=14.8448ns/N (=15201/1024) R²=0.998; realm.Sub mirrors this — keep gno.OpCPUSubRealmBase/Slope in sync (TestSubRealmGasMirrorsPackageAddress)
	{Pkg: "chain", Fn: "deriveStorageDepositAddr", Base: 541, Slope: 471, SlopeIdx: 0, SlopeKind: SizeLenString},                                               // fit base=540.9ns slope=0.4602ns/N (=471/1024) R²=0.994
	{Pkg: "chain", Fn: "pubKeyAddress", Base: 2631, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                         // flat, median 2631.0ns
	{Pkg: "time", Fn: "loadFromEmbeddedTZData", Base: 16068, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                // flat, median 16068.0ns
	{Pkg: "math", Fn: "Float32bits", Base: 32, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                              // flat, median 32.5ns
	{Pkg: "math", Fn: "Float32frombits", Base: 32, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                          // flat, median 32.4ns
	{Pkg: "math", Fn: "Float64bits", Base: 29, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                              // flat, median 28.7ns
	{Pkg: "math", Fn: "Float64frombits", Base: 29, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                          // flat, median 28.8ns
	{Pkg: "chain/banker", Fn: "bankerSendCoins", Base: 322, Slope: 35318, SlopeIdx: 3, SlopeKind: SizeLenSlice},                                                // fit base=321.9ns slope=34.4898ns/N (=35318/1024) R²=0.999
	{Pkg: "chain/banker", Fn: "bankerGetCoins", Base: 349, SlopeIdx: -1, SlopeKind: SizeFlat, PostSlope: 36206, PostSlopeIdx: 2, PostSlopeKind: SizeReturnLen}, // post-call: base=349.1ns + 35.3578ns/N (=36206/1024) R²=0.998
	{Pkg: "chain/banker", Fn: "bankerGetCoin", Base: 129, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                   // flat, median 129.2ns
	{Pkg: "chain/banker", Fn: "bankerTotalCoin", Base: 87, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                  // flat, median 87.0ns
	{Pkg: "chain/banker", Fn: "bankerIssueCoin", Base: 141, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                 // flat, median 140.6ns
	{Pkg: "chain/banker", Fn: "bankerRemoveCoin", Base: 196, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                // flat, median 195.9ns
	{Pkg: "chain/params", Fn: "SetBytes", Base: 1912, Slope: 13213, SlopeIdx: 1, SlopeKind: SizeLenBytes},                                                      // fit base=1912.0ns slope=12.9035ns/N (=13213/1024) R²=1.000
	{Pkg: "chain/params", Fn: "SetString", Base: 1772, Slope: 135, SlopeIdx: 1, SlopeKind: SizeLenString},                                                      // fit base=1772.3ns slope=0.1323ns/N (=135/1024) R²=0.933
	{Pkg: "chain/params", Fn: "SetBool", Base: 1643, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                        // flat, median 1643.0ns
	{Pkg: "chain/params", Fn: "SetInt64", Base: 1201, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                       // flat, median 1201.0ns
	{Pkg: "chain/params", Fn: "SetUint64", Base: 1219, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                      // flat, median 1219.0ns
	{Pkg: "chain/params", Fn: "GetBytes", Base: 1912, SlopeIdx: -1, SlopeKind: SizeFlat, PostSlope: 10584, PostSlopeIdx: 2, PostSlopeKind: SizeReturnLen},      // mirrors chain/params keying plus sys/params bytes return cost
	{Pkg: "chain/params", Fn: "GetString", Base: 1772, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                      // mirrors chain/params SetString keying cost
	{Pkg: "chain/params", Fn: "GetBool", Base: 1643, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                        // mirrors chain/params SetBool keying cost
	{Pkg: "chain/params", Fn: "GetInt64", Base: 1201, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                       // mirrors chain/params SetInt64 keying cost
	{Pkg: "chain/params", Fn: "GetUint64", Base: 1219, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                      // mirrors chain/params SetUint64 keying cost
	{Pkg: "sys/params", Fn: "setSysParamBytes", Base: 323, Slope: 9703, SlopeIdx: 3, SlopeKind: SizeLenBytes},                                                  // fit base=323.3ns slope=9.4757ns/N (=9703/1024) R²=0.995
	{Pkg: "sys/params", Fn: "getSysParamBytes", Base: 416, SlopeIdx: -1, SlopeKind: SizeFlat, PostSlope: 10584, PostSlopeIdx: 2, PostSlopeKind: SizeReturnLen}, // post-call: base=415.7ns + 10.3357ns/N (=10584/1024) R²=1.000
	{Pkg: "sys/params", Fn: "setSysParamString", Base: 269, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                 // flat, median 269.1ns
	{Pkg: "sys/params", Fn: "setSysParamBool", Base: 217, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                   // flat, median 217.4ns
	{Pkg: "sys/params", Fn: "setSysParamInt64", Base: 228, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                  // flat, median 227.8ns
	{Pkg: "sys/params", Fn: "setSysParamUint64", Base: 299, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                 // flat, median 298.9ns
	{Pkg: "sys/params", Fn: "getSysParamBool", Base: 236, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                   // flat, median 236.5ns
	{Pkg: "sys/params", Fn: "getSysParamInt64", Base: 323, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                  // flat, median 322.9ns
	{Pkg: "sys/params", Fn: "getSysParamUint64", Base: 309, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                 // flat, median 308.7ns
	{Pkg: "sys/params", Fn: "getSysParamString", Base: 363, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                 // flat, median 362.8ns
	{Pkg: "chain/runtime", Fn: "ChainID", Base: 45, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                         // flat, median 44.8ns
	{Pkg: "chain/runtime", Fn: "ChainDomain", Base: 45, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                     // flat, median 44.5ns
	{Pkg: "chain/runtime", Fn: "ChainHeight", Base: 30, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                     // flat, median 30.2ns
	{Pkg: "chain/runtime", Fn: "getSessionInfo", Base: 148, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                 // flat, median 148.4ns
	{Pkg: "chain/runtime", Fn: "AssertOriginCall", Base: 5, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                 // flat, median 5.0ns
	// chain/runtime/unsafe natives — same implementations as their chain/runtime / chain/banker counterparts.
	{Pkg: "chain/runtime/unsafe", Fn: "originCaller", Base: 45, SlopeIdx: -1, SlopeKind: SizeFlat},                                                               // mirrors chain/runtime.originCaller
	{Pkg: "chain/runtime/unsafe", Fn: "getRealm", Base: 1003, Slope: 1319, SlopeIdx: -1, SlopeKind: SizeNumCallFrames},                                           // mirrors chain/runtime.getRealm
	{Pkg: "chain/runtime/unsafe", Fn: "originSend", Base: 280, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                // mirrors chain/banker.originSend
	{Pkg: "time", Fn: "now", Base: 47, SlopeIdx: -1, SlopeKind: SizeFlat},                                                                                        // flat, median 46.9ns
	{Pkg: "chain", Fn: "emit", Base: 362, Slope: 40218, SlopeIdx: 1, SlopeKind: SizeLenSlice},                                                                    // fit base=361.9ns slope=39.2750ns/N (=40218/1024) R²=0.955
	{Pkg: "chain/params", Fn: "SetStrings", Base: 1601, Slope: 39842, SlopeIdx: 1, SlopeKind: SizeLenSlice},                                                      // fit base=1601.1ns slope=38.9082ns/N (=39842/1024) R²=0.993
	{Pkg: "chain/params", Fn: "UpdateParamStrings", Base: 1298, Slope: 24077, SlopeIdx: 1, SlopeKind: SizeLenSlice},                                              // fit base=1298.0ns slope=23.5122ns/N (=24077/1024) R²=1.000
	{Pkg: "chain/params", Fn: "GetStrings", Base: 1601, SlopeIdx: -1, SlopeKind: SizeFlat, PostSlope: 23215, PostSlopeIdx: 2, PostSlopeKind: SizeReturnLen},      // mirrors chain/params keying plus sys/params strings return cost
	{Pkg: "sys/params", Fn: "setSysParamStrings", Base: 341, Slope: 27034, SlopeIdx: 3, SlopeKind: SizeLenSlice},                                                 // fit base=341.0ns slope=26.4006ns/N (=27034/1024) R²=0.997
	{Pkg: "sys/params", Fn: "updateSysParamStrings", Base: 413, Slope: 26861, SlopeIdx: 3, SlopeKind: SizeLenSlice},                                              // fit base=413.4ns slope=26.2318ns/N (=26861/1024) R²=0.998
	{Pkg: "sys/params", Fn: "getSysParamStrings", Base: 349, SlopeIdx: -1, SlopeKind: SizeFlat, PostSlope: 23215, PostSlopeIdx: 2, PostSlopeKind: SizeReturnLen}, // post-call: base=348.9ns + 22.6713ns/N (=23215/1024) R²=0.999
	// chain/markdown — calibrated from gnovm/cmd/calibrate (M2 ARM64 baseline,
	// same as the other rows above). Same Xeon 8168 re-calibration caveat
	// applies; values are stable across re-runs (R² ≥ 0.994 for all eight).
	{Pkg: "chain/markdown", Fn: "StripBidiAndZeroWidth", Base: 71, Slope: 1882, SlopeIdx: 0, SlopeKind: SizeLenString},    // fit base=71.1ns slope=1.8378ns/N (=1882/1024) R²=0.998
	{Pkg: "chain/markdown", Fn: "NormalizeBreaks", Base: 72, Slope: 573, SlopeIdx: 0, SlopeKind: SizeLenString},           // fit base=71.8ns slope=0.5600ns/N (=573/1024) R²=0.994
	{Pkg: "chain/markdown", Fn: "EscapeInline", Base: 76, Slope: 2041, SlopeIdx: 0, SlopeKind: SizeLenString},             // fit base=75.9ns slope=1.9932ns/N (=2041/1024) R²=0.995
	{Pkg: "chain/markdown", Fn: "EscapeTitle", Base: 76, Slope: 2060, SlopeIdx: 0, SlopeKind: SizeLenString},              // fit base=76.1ns slope=2.0120ns/N (=2060/1024) R²=0.995
	{Pkg: "chain/markdown", Fn: "PercentEncodeURL", Base: 77, Slope: 4281, SlopeIdx: 0, SlopeKind: SizeLenString},         // fit base=77.0ns slope=4.1804ns/N (=4281/1024) R²=0.998
	{Pkg: "chain/markdown", Fn: "MatchCharsetN", Base: 172, Slope: 685, SlopeIdx: 0, SlopeKind: SizeLenString},            // fit base=172.4ns slope=0.6685ns/N (=685/1024) R²=1.000
	{Pkg: "chain/markdown", Fn: "CodeFence", Base: 99, Slope: 1032, SlopeIdx: 0, SlopeKind: SizeLenString},                // fit base=98.9ns slope=1.0076ns/N (=1032/1024) R²=0.998
	{Pkg: "chain/markdown", Fn: "EscapeBlockHazards", Base: 136, Slope: 27361, SlopeIdx: 0, SlopeKind: SizeLenString},     // fit base=136ns slope=26.72ns/N (=27361/1024) R²=1.000 — shape `[a]: u\n(x\n` (post bracket-walker worst case, scanLRDTail paren-title budget exhausts)
	{Pkg: "chain/markdown", Fn: "EscapeBlockHazardsRich", Base: 136, Slope: 27597, SlopeIdx: 0, SlopeKind: SizeLenString}, // fit base=136ns slope=26.95ns/N (=27597/1024) R²=1.000 — same worst case as EscapeBlockHazards (bracket walker dominates); ~1% above strict variant despite skipping two per-line checks (within measurement noise)
	{Pkg: "chain/markdown", Fn: "MaxForeignBlocksPerConvert", Base: 32, SlopeIdx: -1, SlopeKind: SizeFlat},                // flat: returns a compile-time constant, no input

	// --- IBC crypto stdlibs (draft, Xeon Silver 4114) ---
	{Pkg: "crypto/keccak256", Fn: "sum256", Base: 4323, Slope: 23654, SlopeIdx: 0, SlopeKind: SizeLenBytes},       // draft fit base=4323ns slope=23.10ns/N (=23654/1024) on 0..16384 bytes
	{Pkg: "crypto/bn254", Fn: "g1Add", Base: 14883, SlopeIdx: -1, SlopeKind: SizeFlat},                            // draft, median 14883ns (input fixed 128B)
	{Pkg: "crypto/bn254", Fn: "g1Mul", Base: 44465, SlopeIdx: -1, SlopeKind: SizeFlat},                            // draft, median 44465ns (input fixed 96B)
	{Pkg: "crypto/bn254", Fn: "pairingCheck", Base: 457574, Slope: 1786890, SlopeIdx: 0, SlopeKind: SizeLenBytes}, // draft fit base=457574ns slope=1745.4ns/N (=1786890/1024) on 1..4 pairs
	{Pkg: "crypto/cometbls", Fn: "verifyZKP", Base: 2632556, SlopeIdx: -1, SlopeKind: SizeFlat},                   // draft, median 2.63ms (full Groth16 verify; proof always 384B, header 116B)
	{Pkg: "crypto/merkle", Fn: "leafHash", Base: 3528, Slope: 32911, SlopeIdx: 0, SlopeKind: SizeLenBytes},        // draft fit base=3528ns slope=32.14ns/N (=32911/1024) on 0..4096 bytes
	// innerHash hashes 0x01||left||right, so its cost is O(len(left)+len(right))
	// — leafHash's shape over two operands instead of one, and nothing bounds
	// either. Both operands therefore carry leafHash's per-byte rate, which
	// makes the total charge track total hashed bytes at the same rate leafHash
	// pays. It was flat until 1MiB+1MiB was measured running ~300x over its
	// price; see adr/native_input_bounds.md.
	//
	// Draft fit base=7513ns (still the 32+32B median, so it double-counts ~2µs
	// of per-byte cost at that size — conservative until the pending reference
	// recalibration turns it into a proper intercept) + 32.14ns/N (=32911/1024)
	// on each of len(left) and len(right).
	{Pkg: "crypto/merkle", Fn: "innerHash", Base: 7513, Slope: 32911, SlopeIdx: 0, SlopeKind: SizeLenBytes, Slope2: 32911, Slope2Idx: 1, Slope2Kind: SizeLenBytes},
	{Pkg: "crypto/merkle", Fn: "hashFromByteSlices", Base: 4839, Slope: 188621, SlopeIdx: 0, SlopeKind: SizeLenBytes}, // draft fit base=4839ns slope=184.2ns/N (=188621/1024) on encoded 1..512 items
	{Pkg: "crypto/merkle", Fn: "verifySimpleProof", Base: 4567, Slope: 53533, SlopeIdx: 4, SlopeKind: SizeLenBytes},   // draft fit base=4567ns slope=52.3ns/N (=53533/1024) on aunts 96..320 bytes
	// modExp is charged on two independent components, because it has two:
	//
	//   - Slope, on SizeModExpWork at SlopeIdx=1: the exponentiation itself. The
	//     kind counts the modular multiplications big.Int.Exp will perform,
	//     reading the branch structure of nat.go's expNN off the operand lengths;
	//     derivation lives with it in gnovm/pkg/gnolang/native_gas.go. Slope
	//     converts that count to nanoseconds, so it is the only hardware-dependent
	//     part of the model. SlopeIdx names the exponent and the modulus is read
	//     from the next parameter.
	//   - Slope2, on SizeLenBytes at the modulus: converting the operands across
	//     the dispatcher, allocating the result and filling it. This is linear in
	//     len(modulus) and runs even when the exponent is empty and no
	//     exponentiation happens at all, so folding it into Base makes the charge
	//     for a zero-length exponent independent of the modulus — a hole, since
	//     the call still allocates and fills len(modulus) bytes. It is a small
	//     term now (~2.9 ns/byte): byte slices became Data-backed in #97, so the
	//     dispatcher no longer converts one TypedValue per byte. Before that it
	//     was ~134 ns/byte and this slope was 49x larger.
	//
	// Both slopes are UPPER bounds over the ModExpGrid benches rather than the
	// least-squares coefficients: a central fit sits below cost on about half the
	// grid, which is fine for describing a cost and wrong for charging one. That is
	// what gen_native_table.py's fit_modexp now emits, so this row IS a fitter
	// output and regenerating reproduces it, rather than being a hand-edit the
	// fitter would silently undo.
	//
	// Over the recorded grid, projected onto reference hardware, every point lands
	// between 1.24x and 2.57x of measured cost, and the charge is monotonic in
	// len(exp) at every modulus — the previous row charged an 8-byte exponent less
	// than a 9-byte one while it measured more expensive. Re-fit with
	// gnovm/cmd/calibrate/ibc_native_bench_test.go, which records the run these came
	// from and the hardware it was taken on; pass --hw-factor unless you are on the
	// reference Xeon. TestModExpRowCoversRecordedGrid checks every point.
	//
	// X_modExp and ModExp cap each operand at 1024 bytes, matching EIP-7823. The cap
	// bounds the envelope this fit was checked over and the unpriced base reduction;
	// this row does the pricing.
	{
		Pkg: "crypto/modexp", Fn: "modExp", Base: 1400,
		Slope: 2200, SlopeIdx: 1, SlopeKind: SizeModExpWork,
		Slope2: 3500, Slope2Idx: 2, Slope2Kind: SizeLenBytes,
	}, // draft, upper bound over 26 ModExpGrid points (x2.3 to reference, +19% margin); thinnest at expLen=3/modLen=256
}

// validateRow cross-checks a row against the native's actual signature before
// registration. Every other misconfiguration in this table already fails loudly
// — an unregistered native panics on first call, and a bad SlopeIdx panics on
// the block index — but SizeModExpWork reads a *pair* of parameters and cannot
// panic on a missing second one without giving up the bounds check. Left
// unvalidated it would degrade silently, and toward free: the crypto/modexp row
// carried SlopeIdx: 2 until recently, and with that value the pair read runs off
// the end of the call block, returns zero work, and collapses the charge to
// Base for any operand size. Catch it at boot instead.
func validateRow(e nativeGasEntry, params int) {
	if e.SlopeKind == SizeModExpWork {
		if e.SlopeIdx < 0 || int(e.SlopeIdx)+1 >= params {
			panic(fmt.Sprintf("%s.%s: SizeModExpWork needs params at SlopeIdx and SlopeIdx+1, "+
				"but SlopeIdx=%d and the native takes %d params — the charge would silently "+
				"collapse to Base", e.Pkg, e.Fn, e.SlopeIdx, params))
		}
	}
	// The metric needs a parameter pair, so it can never be read off the
	// return stack. nativeSizeOf panics if it ever is; refuse it here too so
	// the mistake surfaces at boot rather than on the first call.
	if e.Slope2Kind == SizeModExpWork || e.PostSlopeKind == SizeModExpWork || e.PostSlope2Kind == SizeModExpWork {
		panic(fmt.Sprintf("%s.%s: SizeModExpWork is only valid as SlopeKind (pre-call)", e.Pkg, e.Fn))
	}
}

func init() {
	// Index the generated bindings so rows can be checked against the real
	// signatures rather than against a hand-maintained duplicate of them.
	nParams := make(map[string]int, len(nativeFuncs))
	for _, nf := range nativeFuncs {
		nParams[nf.gnoPkg+"\x00"+string(nf.gnoFunc)] = len(nf.params)
	}
	for _, e := range calibratedNativeGas {
		if n, ok := nParams[e.Pkg+"\x00"+e.Fn]; ok {
			validateRow(e, n)
		}
		gno.RegisterNativeGas(e.Pkg, gno.Name(e.Fn), &gno.NativeGasInfo{
			Base:           e.Base,
			Slope:          e.Slope,
			SlopeIdx:       e.SlopeIdx,
			SlopeKind:      e.SlopeKind,
			Slope2:         e.Slope2,
			Slope2Idx:      e.Slope2Idx,
			Slope2Kind:     e.Slope2Kind,
			PostBase:       e.PostBase,
			PostSlope:      e.PostSlope,
			PostSlopeIdx:   e.PostSlopeIdx,
			PostSlopeKind:  e.PostSlopeKind,
			PostSlope2:     e.PostSlope2,
			PostSlope2Idx:  e.PostSlope2Idx,
			PostSlope2Kind: e.PostSlope2Kind,
		})
	}
}
