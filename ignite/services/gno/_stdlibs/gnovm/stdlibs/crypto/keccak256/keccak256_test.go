// Test vectors sourced from:
//   https://eips.ethereum.org/EIPS/eip-1
//   https://emn178.github.io/online-tools/keccak_256.html

package keccak256

import (
	"encoding/hex"
	"testing"
)

func TestSum256KnownAnswers(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string // hex-encoded expected digest
	}{
		// The empty-string digest is the sanity check that pins down Legacy
		// Keccak-256 vs NIST SHA3-256 (the two disagree on their padding
		// byte).
		{"empty", "", "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"},
		{"abc", "abc", "4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45"},
		{"hello", "hello", "1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := X_sum256([]byte(tc.in))
			if hex.EncodeToString(got[:]) != tc.want {
				t.Fatalf("X_sum256(%q) = %x, want %s", tc.in, got, tc.want)
			}
		})
	}
}

// Sanity check: the native binding must reject NIST SHA3-256 output if
// someone ever swapped the implementation. SHA3-256 of "" is
// a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a, which
// must NOT equal our digest.
func TestSum256NotSHA3(t *testing.T) {
	got := X_sum256(nil)
	sha3Empty := "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"
	if hex.EncodeToString(got[:]) == sha3Empty {
		t.Fatalf("X_sum256 returns NIST SHA3-256 instead of Legacy Keccak-256")
	}
}
