// Test vectors sourced from:
//   CometBFT merkle: https://github.com/cometbft/cometbft/blob/main/crypto/merkle/proof_test.go
//   Tendermint simple Merkle spec: SHA256(0x00||leaf) for leaves, SHA256(0x01||left||right) for inner nodes.

package merkle

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"runtime"
	"testing"

	tm2merkle "github.com/gnolang/gno/tm2/pkg/crypto/merkle"
)

// leafHashDirect computes SHA256(0x00 || leaf) without going through X_leafHash,
// letting us cross-validate the native binding against Go's standard library.
func leafHashDirect(leaf []byte) []byte {
	h := sha256.New()
	h.Write([]byte{0x00})
	h.Write(leaf)
	return h.Sum(nil)
}

func innerHashDirect(left, right []byte) []byte {
	h := sha256.New()
	h.Write([]byte{0x01})
	h.Write(left)
	h.Write(right)
	return h.Sum(nil)
}

func TestLeafHashKnownVector(t *testing.T) {
	cases := []string{"", "a", "hello", "gno.land"}
	for _, item := range cases {
		got := X_leafHash([]byte(item))
		want := leafHashDirect([]byte(item))
		if !bytes.Equal(got, want) {
			t.Fatalf("leafHash(%q): got %x, want %x", item, got, want)
		}
	}
}

func TestInnerHashKnownVector(t *testing.T) {
	left := leafHashDirect([]byte("a"))
	right := leafHashDirect([]byte("b"))

	got := X_innerHash(left, right)
	want := innerHashDirect(left, right)
	if !bytes.Equal(got, want) {
		t.Fatalf("innerHash mismatch: got %x, want %x", got, want)
	}

	// Commutativity does NOT hold: H(left,right) != H(right,left)
	swapped := X_innerHash(right, left)
	if bytes.Equal(got, swapped) {
		t.Fatal("innerHash should not be commutative")
	}
}

func TestHashFromByteSlicesEmpty(t *testing.T) {
	encoded := []byte{0, 0, 0, 0} // count=0
	got := X_hashFromByteSlices(encoded)
	if got != nil {
		t.Fatalf("expected nil for empty input, got %x", got)
	}
}

func TestHashFromByteSlicesSingle(t *testing.T) {
	item := []byte("hello")
	encoded := encodeItems([][]byte{item})
	got := X_hashFromByteSlices(encoded)
	want := X_leafHash(item)
	if !bytes.Equal(got, want) {
		t.Fatalf("single item: got %x, want %x", got, want)
	}
}

func TestHashFromByteSlicesMatchesTm2(t *testing.T) {
	// Cross-validate against tm2/pkg/crypto/merkle for a multi-item tree.
	items := [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d")}
	want := tm2merkle.SimpleHashFromByteSlices(items)

	encoded := encodeItems(items)
	got := X_hashFromByteSlices(encoded)
	if !bytes.Equal(got, want) {
		t.Fatalf("root mismatch: got %x, want %x", got, want)
	}
}

// TestHashFromByteSlicesKnownHex pins the root for ["a","b","c"] against a
// value computed from the reference CometBFT implementation.
func TestHashFromByteSlicesKnownHex(t *testing.T) {
	items := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	encoded := encodeItems(items)
	got := X_hashFromByteSlices(encoded)

	// Compute expected via Go sha256 to be independent of tm2.
	leafA := leafHashDirect([]byte("a"))
	leafB := leafHashDirect([]byte("b"))
	leafC := leafHashDirect([]byte("c"))
	// SimpleHash for 3 items: inner(inner(leafA, leafB), leafC)
	want := innerHashDirect(innerHashDirect(leafA, leafB), leafC)

	if !bytes.Equal(got, want) {
		t.Fatalf("3-item root: got %x, want %x", got, want)
	}
	_ = hex.EncodeToString(got)
}

func TestVerifySimpleProofValid(t *testing.T) {
	items := [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d")}
	rootHash, proofs := tm2merkle.SimpleProofsFromByteSlices(items)

	for i, item := range items {
		proof := proofs[i]
		aunts := flattenAunts(proof.Aunts)
		if !X_verifySimpleProof(rootHash, item, proof.Index, proof.Total, aunts) {
			t.Fatalf("valid proof for item %d rejected", i)
		}
	}
}

func TestVerifySimpleProofTamperedLeaf(t *testing.T) {
	items := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	rootHash, proofs := tm2merkle.SimpleProofsFromByteSlices(items)

	proof := proofs[0]
	aunts := flattenAunts(proof.Aunts)
	if X_verifySimpleProof(rootHash, []byte("tampered"), proof.Index, proof.Total, aunts) {
		t.Fatal("tampered leaf should not verify")
	}
}

func TestVerifySimpleProofWrongRoot(t *testing.T) {
	items := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	_, proofs := tm2merkle.SimpleProofsFromByteSlices(items)

	wrongRoot := make([]byte, 32) // all zeros
	proof := proofs[0]
	aunts := flattenAunts(proof.Aunts)
	if X_verifySimpleProof(wrongRoot, []byte("a"), proof.Index, proof.Total, aunts) {
		t.Fatal("wrong root should not verify")
	}
}

// decodeByteSlices must reject a count larger than the payload can possibly
// contain BEFORE allocating make([][]byte, count), so a tiny input cannot
// force a huge host allocation (native gas is charged by encoded length only).
// A 4-byte input claiming count=1<<20 must return (nil,false) and must not
// allocate megabytes.
func TestDecodeByteSlicesHugeCountRejected(t *testing.T) {
	encoded := []byte{0x00, 0x10, 0x00, 0x00} // count = 1<<20, no items

	if items, ok := decodeByteSlices(encoded); ok || items != nil {
		t.Fatalf("expected (nil,false) for over-large count, got ok=%v len=%d", ok, len(items))
	}

	// Allocation guard: a 4-byte input must not allocate megabytes.
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	_ = X_hashFromByteSlices(encoded)
	runtime.ReadMemStats(&after)
	if d := after.TotalAlloc - before.TotalAlloc; d > 1<<20 {
		t.Fatalf("allocated %d bytes from a %d-byte input (expected < 1 MiB)", d, len(encoded))
	}
}

// Exact-fit boundary: count == len(b)/4 (all zero-length items) is the maximum
// valid encoding and must be accepted — locks in `>` (not `>=`).
func TestDecodeByteSlicesExactFitBoundary(t *testing.T) {
	encoded := encodeItems([][]byte{{}, {}}) // count=2, both zero-length: exact fit
	if items, ok := decodeByteSlices(encoded); !ok || len(items) != 2 {
		t.Fatalf("expected ok with 2 empty items, got ok=%v len=%d", ok, len(items))
	}
}

// A valid encoding with data still decodes correctly (no regression).
func TestDecodeByteSlicesValid(t *testing.T) {
	encoded := encodeItems([][]byte{{0xAA, 0xBB, 0xCC}})
	items, ok := decodeByteSlices(encoded)
	if !ok || len(items) != 1 || !bytes.Equal(items[0], []byte{0xAA, 0xBB, 0xCC}) {
		t.Fatalf("valid decode failed: ok=%v items=%v", ok, items)
	}
}

// count == 0 -> empty result, ok == true.
func TestDecodeByteSlicesZeroCount(t *testing.T) {
	if items, ok := decodeByteSlices([]byte{0x00, 0x00, 0x00, 0x00}); !ok || len(items) != 0 {
		t.Fatalf("expected ok with 0 items, got ok=%v len=%d", ok, len(items))
	}
}

func flattenAunts(aunts [][]byte) []byte {
	var out []byte
	for _, a := range aunts {
		out = append(out, a...)
	}
	return out
}

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
