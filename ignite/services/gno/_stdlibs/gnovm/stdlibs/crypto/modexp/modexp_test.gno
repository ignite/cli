package modexp_test

import (
	"bytes"
	"crypto/modexp"
	"testing"
)

// ones returns n bytes of 0xFF.
func ones(n int) []byte { return bytes.Repeat([]byte{0xFF}, n) }

// oddMod returns an n-byte modulus that is large and odd, so accepted calls take
// the Montgomery path rather than a shortcut.
func oddMod(n int) []byte {
	m := ones(n)
	m[n-1] = 0xFD
	return m
}

func TestModExpKnownAnswer(t *testing.T) {
	// 2^10 mod 1000 = 24.
	got := modexp.ModExp([]byte{2}, []byte{10}, []byte{0x03, 0xe8})
	if !bytes.Equal(got, []byte{0x00, 0x18}) {
		t.Fatalf("2^10 mod 1000: got %x, want 0018", got)
	}
}

// TestModExpOperandLenBoundary checks each operand independently at
// MaxOperandLen and one byte over it. Only one operand sits at the limit per
// case — a symmetric all-at-1024 call is a real ~140ms modular exponentiation,
// and it proves nothing about a length boundary that a 1-byte exponent does not.
//
// This exercises the guard in ModExp itself, which is the one that matters: it
// refuses oversized operands before the native is dispatched, so the VM never
// converts them to Go. The native repeats the check as defense in depth.
func TestModExpOperandLenBoundary(t *testing.T) {
	lim := modexp.MaxOperandLen

	// len(nil) == 0, so length covers acceptance and rejection alike.
	cases := []struct {
		name           string
		base, exp, mod []byte
		wantLen        int
	}{
		{"base at limit", ones(lim), ones(1), oddMod(32), 32},
		{"exp at limit", ones(32), ones(lim), oddMod(32), 32},
		{"modulus at limit", ones(32), ones(1), oddMod(lim), lim},
		{"base over limit", ones(lim + 1), ones(1), oddMod(32), 0},
		{"exp over limit", ones(32), ones(lim + 1), oddMod(32), 0},
		{"modulus over limit", ones(32), ones(1), oddMod(lim + 1), 0},
	}
	for _, tc := range cases {
		got := modexp.ModExp(tc.base, tc.exp, tc.mod)
		if len(got) != tc.wantLen {
			t.Errorf("%s: got %d bytes, want %d", tc.name, len(got), tc.wantLen)
		}
	}
}
