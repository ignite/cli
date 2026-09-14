// Test vectors cross-validated against:
//   CometBFT merkle: https://github.com/cometbft/cometbft/blob/main/crypto/merkle/proof_test.go
//   Tendermint simple Merkle spec: SHA256(0x00||leaf) for leaves, SHA256(0x01||left||right) for inner nodes.

package merkle_test

import (
	"bytes"
	"crypto/merkle"
	"crypto/sha256"
	"testing"
)

// leafHashDirect computes SHA256(0x00 || leaf) via crypto/sha256, letting us
// cross-validate LeafHash without depending on the native binding under test.
func leafHashDirect(leaf []byte) []byte {
	b := make([]byte, 1+len(leaf))
	b[0] = 0x00
	copy(b[1:], leaf)
	s := sha256.Sum256(b)
	return s[:]
}

func innerHashDirect(left, right []byte) []byte {
	b := make([]byte, 1+len(left)+len(right))
	b[0] = 0x01
	copy(b[1:], left)
	copy(b[1+len(left):], right)
	s := sha256.Sum256(b)
	return s[:]
}

// encodeItems packs items into the [4-byte BE count][4-byte BE len][data]...
// wire format expected by HashFromByteSlices.
func encodeItems(items [][]byte) []byte {
	n := len(items)
	out := []byte{byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)}
	for _, item := range items {
		l := len(item)
		out = append(out, byte(l>>24), byte(l>>16), byte(l>>8), byte(l))
		out = append(out, item...)
	}
	return out
}

// flatAunts concatenates sibling hashes into the flat byte slice expected by
// VerifySimpleProof.
func flatAunts(aunts ...[]byte) []byte {
	var out []byte
	for _, a := range aunts {
		out = append(out, a...)
	}
	return out
}

// TestLeafHashMatchesSHA256 cross-validates LeafHash against a direct SHA256
// computation for a range of inputs including the empty byte slice.
func TestLeafHashMatchesSHA256(t *testing.T) {
	cases := []string{"", "a", "hello", "gno.land"}
	for _, item := range cases {
		got := merkle.LeafHash([]byte(item))
		want := leafHashDirect([]byte(item))
		if !bytes.Equal(got, want) {
			t.Errorf("LeafHash(%q): got %x, want %x", item, got, want)
		}
	}
}

// TestInnerHashMatchesSHA256 cross-validates InnerHash and asserts that it is
// not commutative: H(left, right) must differ from H(right, left).
func TestInnerHashMatchesSHA256(t *testing.T) {
	left := leafHashDirect([]byte("a"))
	right := leafHashDirect([]byte("b"))

	got := merkle.InnerHash(left, right)
	want := innerHashDirect(left, right)
	if !bytes.Equal(got, want) {
		t.Errorf("InnerHash mismatch: got %x, want %x", got, want)
	}

	swapped := merkle.InnerHash(right, left)
	if bytes.Equal(got, swapped) {
		t.Error("InnerHash must not be commutative")
	}
}

func TestHashFromByteSlicesEmpty(t *testing.T) {
	encoded := []byte{0, 0, 0, 0} // count = 0
	got := merkle.HashFromByteSlices(encoded)
	if len(got) != 0 {
		t.Errorf("expected empty result for empty slice list, got %x", got)
	}
}

func TestHashFromByteSlicesSingle(t *testing.T) {
	item := []byte("hello")
	encoded := encodeItems([][]byte{item})
	got := merkle.HashFromByteSlices(encoded)
	want := merkle.LeafHash(item)
	if !bytes.Equal(got, want) {
		t.Errorf("single item: got %x, want %x", got, want)
	}
}

// TestHashFromByteSlicesThreeItems pins the root for ["a","b","c"] by building
// the tree manually: inner(inner(leafA, leafB), leafC).
func TestHashFromByteSlicesThreeItems(t *testing.T) {
	items := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	got := merkle.HashFromByteSlices(encodeItems(items))

	leafA := leafHashDirect([]byte("a"))
	leafB := leafHashDirect([]byte("b"))
	leafC := leafHashDirect([]byte("c"))
	want := innerHashDirect(innerHashDirect(leafA, leafB), leafC)

	if !bytes.Equal(got, want) {
		t.Errorf("3-item root mismatch:\n got  %x\n want %x", got, want)
	}
}

// TestVerifySimpleProofValid builds a balanced 4-item tree manually and checks
// that every leaf's membership proof is accepted.
//
// Tree shape:
//
//	    root
//	   /    \
//	 i01    i23
//	 / \    / \
//	la  lb lc  ld
func TestVerifySimpleProofValid(t *testing.T) {
	la := merkle.LeafHash([]byte("a"))
	lb := merkle.LeafHash([]byte("b"))
	lc := merkle.LeafHash([]byte("c"))
	ld := merkle.LeafHash([]byte("d"))
	i01 := merkle.InnerHash(la, lb)
	i23 := merkle.InnerHash(lc, ld)
	root := merkle.InnerHash(i01, i23)

	cases := []struct {
		leaf  []byte
		index int
		aunts []byte
	}{
		{[]byte("a"), 0, flatAunts(lb, i23)},
		{[]byte("b"), 1, flatAunts(la, i23)},
		{[]byte("c"), 2, flatAunts(ld, i01)},
		{[]byte("d"), 3, flatAunts(lc, i01)},
	}

	for _, tc := range cases {
		if !merkle.VerifySimpleProof(root, tc.leaf, tc.index, 4, tc.aunts) {
			t.Errorf("valid proof for index %d rejected", tc.index)
		}
	}
}

// TestVerifySimpleProofTamperedLeaf checks that substituting a different leaf
// causes proof verification to fail.
func TestVerifySimpleProofTamperedLeaf(t *testing.T) {
	la := merkle.LeafHash([]byte("a"))
	lb := merkle.LeafHash([]byte("b"))
	lc := merkle.LeafHash([]byte("c"))
	ld := merkle.LeafHash([]byte("d"))
	i01 := merkle.InnerHash(la, lb)
	i23 := merkle.InnerHash(lc, ld)
	root := merkle.InnerHash(i01, i23)

	aunts := flatAunts(lb, i23)
	if merkle.VerifySimpleProof(root, []byte("tampered"), 0, 4, aunts) {
		t.Error("tampered leaf should not verify")
	}
}

// TestVerifySimpleProofWrongRoot checks that a valid proof against the wrong
// root is rejected.
func TestVerifySimpleProofWrongRoot(t *testing.T) {
	la := merkle.LeafHash([]byte("a"))
	lb := merkle.LeafHash([]byte("b"))
	lc := merkle.LeafHash([]byte("c"))
	ld := merkle.LeafHash([]byte("d"))
	i01 := merkle.InnerHash(la, lb)
	i23 := merkle.InnerHash(lc, ld)
	merkle.InnerHash(i01, i23) // root (unused — we deliberately pass the wrong one)

	wrongRoot := make([]byte, 32)
	aunts := flatAunts(lb, i23)
	if merkle.VerifySimpleProof(wrongRoot, []byte("a"), 0, 4, aunts) {
		t.Error("wrong root should not verify")
	}
}
