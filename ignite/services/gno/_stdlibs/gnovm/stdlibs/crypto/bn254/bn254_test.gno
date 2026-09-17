// Test vectors from:
//   G1Add:       go-ethereum bn256Add.json     https://github.com/ethereum/go-ethereum/blob/master/core/vm/testdata/precompiles/bn256Add.json
//   PairingCheck: go-ethereum bn256Pairing.json https://github.com/ethereum/go-ethereum/blob/master/core/vm/testdata/precompiles/bn256Pairing.json

package bn254_test

import (
	"bytes"
	"crypto/bn254"
	"encoding/hex"
	"testing"
)

func mustDecodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

// TestG1Add uses verbatim vectors from bn256Add.json.
func TestG1Add(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// cdetrio4: empty input → identity (zero-padded per EIP-196)
		{
			name:     "cdetrio4_empty",
			input:    "",
			expected: "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
		// cdetrio11: G + G = 2G
		{
			name:     "cdetrio11_G_plus_G",
			input:    "0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002",
			expected: "030644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd315ed738c0e0a7c92e7845f96b2ae9c0a68a6a449e3538fc7ff3ebf7a5a18a2c4",
		},
		// cdetrio14: P + (-P) = identity
		{
			name:     "cdetrio14_P_plus_negP",
			input:    "17c139df0efee0f766bc0204762b774362e4ded88953a39ce849a8a7fa163fa901e0559bacb160664764a357af8a9fe70baa9258e0b959273ffc5718c6d4cc7c17c139df0efee0f766bc0204762b774362e4ded88953a39ce849a8a7fa163fa92e83f8d734803fc370eba25ed1f6b8768bd6d83887b87165fc2434fe11a830cb00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			expected: "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
		// chfast1: non-trivial point addition
		{
			name:     "chfast1",
			input:    "18b18acfb4c2c30276db5411368e7185b311dd124691610c5d3b74034e093dc9063c909c4720840cb5134cb9f59fa749755796819658d32efc0d288198f3726607c2b7f58a84bd6145f00c9c2bc0bb1a187f20ff2c92963a88019e7c6a014eed06614e20c147e940f2d70da3f74c9a17df361706a4485c742bd6788478fa17d7",
			expected: "2243525c5efd4b9c3d3c45ac0ca3fe4dd85e830a4ce6b65fa1eeaee202839703301d1d33be6da8e509df21cc35964723180eed7532537db9ae5e7d48f195c915",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var input []byte
			if tc.input != "" {
				input = mustDecodeHex(tc.input)
			}
			want := mustDecodeHex(tc.expected)
			got := bn254.G1Add(input)
			if !bytes.Equal(got, want) {
				t.Errorf("G1Add mismatch:\n got  %x\n want %x", got, want)
			}
		})
	}
}

// TestG1Mul checks scalar multiplication edge cases (no external JSON for G1Mul).
func TestG1Mul(t *testing.T) {
	// G1 generator: x=1, y=2
	g := make([]byte, 64)
	g[31] = 1
	g[63] = 2

	// 0 * G = identity
	input := make([]byte, bn254.G1MulInputSize)
	copy(input[:64], g)
	got := bn254.G1Mul(input)
	if got == nil {
		t.Fatal("G1Mul returned nil for 0*G")
	}
	if !bytes.Equal(got, make([]byte, bn254.G1OutputSize)) {
		t.Errorf("0*G != identity: %x", got)
	}

	// 1 * G = G
	input = make([]byte, bn254.G1MulInputSize)
	copy(input[:64], g)
	input[bn254.G1MulInputSize-1] = 1
	got = bn254.G1Mul(input)
	if !bytes.Equal(got, g) {
		t.Errorf("1*G != G: %x", got)
	}
}

// TestOversizedInputMatchesTruncated pins the length handling that G1Add and
// G1Mul now do before the native call. Those checks exist so an oversized
// input is discarded before the dispatcher copies it into Go memory — both
// natives are priced flat, so that copy would otherwise be free — and this
// test is what says the early exits kept the EIP-196 semantics they replaced.
func TestOversizedInputMatchesTruncated(t *testing.T) {
	// G + G, then a megabyte of trailing garbage. EIP-196 ignores everything
	// past byte 128, so the answer must not change.
	gg := mustDecodeHex("0000000000000000000000000000000000000000000000000000000000000001" +
		"0000000000000000000000000000000000000000000000000000000000000002" +
		"0000000000000000000000000000000000000000000000000000000000000001" +
		"0000000000000000000000000000000000000000000000000000000000000002")
	oversized := append(append([]byte{}, gg...), make([]byte, 1<<20)...)
	if !bytes.Equal(bn254.G1Add(oversized), bn254.G1Add(gg)) {
		t.Errorf("G1Add: oversized input changed the result; excess bytes must be ignored")
	}

	// G1Mul takes no padding at all: any length but 96 is rejected.
	mul := make([]byte, bn254.G1MulInputSize)
	mul[31] = 1
	mul[63] = 2
	mul[bn254.G1MulInputSize-1] = 1
	if got := bn254.G1Mul(append(append([]byte{}, mul...), make([]byte, 1<<20)...)); got != nil {
		t.Errorf("G1Mul: oversized input returned %x, want nil", got)
	}
	if got := bn254.G1Mul(mul[:95]); got != nil {
		t.Errorf("G1Mul: short input returned %x, want nil", got)
	}
}

// TestPairingCheck uses verbatim vectors from bn256Pairing.json.
func TestPairingCheck(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		wantLast byte // last byte of 32-byte output: 0x01 = true, 0x00 = false
	}{
		// empty_data: product of zero pairings = 1
		{name: "empty_data", input: "", wantLast: 1},
		// jeff1: valid 2-pair equality check
		{
			name:     "jeff1",
			input:    "1c76476f4def4bb94541d57ebba1193381ffa7aa76ada664dd31c16024c43f593034dd2920f673e204fee2811c678745fc819b55d3e9d294e45c9b03a76aef41209dd15ebff5d46c4bd888e51a93cf99a7329636c63514396b4a452003a35bf704bf11ca01483bfa8b34b43561848d28905960114c8ac04049af4b6315a416782bb8324af6cfc93537a2ad1a445cfd0ca2a71acd7ac41fadbf933c2a51be344d120a2a4cf30c1bf9845f20c6fe39e07ea2cce61f0c9bb048165fe5e4de877550111e129f1cf1097710d41c4ac70fcdfa5ba2023c6ff1cbeac322de49d1b6df7c2032c61a830e3c17286de9462bf242fca2883585b93870a73853face6a6bf411198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa",
			wantLast: 1,
		},
		// jeff6: pairing inequality → 0
		{
			name:     "jeff6",
			input:    "1c76476f4def4bb94541d57ebba1193381ffa7aa76ada664dd31c16024c43f593034dd2920f673e204fee2811c678745fc819b55d3e9d294e45c9b03a76aef41209dd15ebff5d46c4bd888e51a93cf99a7329636c63514396b4a452003a35bf704bf11ca01483bfa8b34b43561848d28905960114c8ac04049af4b6315a416782bb8324af6cfc93537a2ad1a445cfd0ca2a71acd7ac41fadbf933c2a51be344d120a2a4cf30c1bf9845f20c6fe39e07ea2cce61f0c9bb048165fe5e4de877550111e129f1cf1097710d41c4ac70fcdfa5ba2023c6ff1cbeac322de49d1b6df7c103188585e2364128fe25c70558f1560f4f9350baf3959e603cc91486e110936198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa",
			wantLast: 0,
		},
		// one_point: single (G1, G2) pair → 0
		{
			name:     "one_point",
			input:    "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c21800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa",
			wantLast: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var input []byte
			if tc.input != "" {
				input = mustDecodeHex(tc.input)
			}
			got := bn254.PairingCheck(input)
			if len(got) != bn254.PairingOutputSize {
				t.Fatalf("expected %d-byte output, got %d", bn254.PairingOutputSize, len(got))
			}
			if got[bn254.PairingOutputSize-1] != tc.wantLast {
				t.Errorf("last byte: got 0x%02x, want 0x%02x", got[bn254.PairingOutputSize-1], tc.wantLast)
			}
		})
	}
}
