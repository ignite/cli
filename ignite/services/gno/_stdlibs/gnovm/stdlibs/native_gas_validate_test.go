package stdlibs

import (
	"strings"
	"testing"
)

// TestValidateRowRejectsBadModExpWorkPair pins the boot-time guard. The
// crypto/modexp row carried SlopeIdx: 2 while it charged on len(modulus); with
// SizeModExpWork that same value makes the pair read run past the end of a
// 3-parameter call block, yielding zero work and collapsing the charge to Base
// for any operand size. That must not be reachable silently.
func TestValidateRowRejectsBadModExpWorkPair(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		row    nativeGasEntry
		params int
		want   string // substring of the expected panic; "" means must not panic
	}{
		{
			name:   "modexp's own shape is valid",
			row:    nativeGasEntry{Pkg: "crypto/modexp", Fn: "modExp", SlopeIdx: 1, SlopeKind: SizeModExpWork},
			params: 3,
		},
		{
			name:   "the superseded SlopeIdx would read past the block",
			row:    nativeGasEntry{Pkg: "crypto/modexp", Fn: "modExp", SlopeIdx: 2, SlopeKind: SizeModExpWork},
			params: 3,
			want:   "collapse to Base",
		},
		{
			name:   "negative SlopeIdx has no pair to read",
			row:    nativeGasEntry{Pkg: "x", Fn: "y", SlopeIdx: -1, SlopeKind: SizeModExpWork},
			params: 3,
			want:   "SizeModExpWork needs params",
		},
		{
			name:   "rejected as a second pre-call slope",
			row:    nativeGasEntry{Pkg: "x", Fn: "y", SlopeIdx: 0, SlopeKind: SizeLenBytes, Slope2Idx: 1, Slope2Kind: SizeModExpWork},
			params: 3,
			want:   "only valid as SlopeKind",
		},
		{
			name:   "rejected as a post-call slope",
			row:    nativeGasEntry{Pkg: "x", Fn: "y", SlopeIdx: -1, SlopeKind: SizeFlat, PostSlopeIdx: 1, PostSlopeKind: SizeModExpWork},
			params: 3,
			want:   "only valid as SlopeKind",
		},
		{
			name:   "unrelated kinds are untouched",
			row:    nativeGasEntry{Pkg: "x", Fn: "y", SlopeIdx: 2, SlopeKind: SizeLenBytes},
			params: 3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var got string
			func() {
				defer func() {
					if r := recover(); r != nil {
						got, _ = r.(string)
					}
				}()
				validateRow(tc.row, tc.params)
			}()
			switch {
			case tc.want == "" && got != "":
				t.Fatalf("unexpected panic: %s", got)
			case tc.want != "" && got == "":
				t.Fatalf("expected a panic containing %q, got none", tc.want)
			case tc.want != "" && !strings.Contains(got, tc.want):
				t.Fatalf("panic %q does not contain %q", got, tc.want)
			}
		})
	}
}

// TestEveryRowValidates runs the shipped table through the guard, so a future
// row that trips it fails here rather than at node boot.
func TestEveryRowValidates(t *testing.T) {
	t.Parallel()
	nParams := make(map[string]int, len(nativeFuncs))
	for _, nf := range nativeFuncs {
		nParams[nf.gnoPkg+"\x00"+string(nf.gnoFunc)] = len(nf.params)
	}
	for _, e := range calibratedNativeGas {
		n, ok := nParams[e.Pkg+"\x00"+e.Fn]
		if !ok {
			t.Errorf("%s.%s has a gas row but no generated binding", e.Pkg, e.Fn)
			continue
		}
		validateRow(e, n)
	}
}
