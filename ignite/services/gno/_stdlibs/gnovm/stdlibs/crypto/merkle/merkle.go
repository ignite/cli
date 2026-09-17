package merkle

import (
	"github.com/gnolang/gno/tm2/pkg/crypto/merkle"
	"github.com/gnolang/gno/tm2/pkg/crypto/tmhash"
)

// X_leafHash returns SHA256(0x00 || leaf), the Tendermint leaf hash.
func X_leafHash(leaf []byte) []byte {
	return tmhash.Sum(append([]byte{0x00}, leaf...))
}

// X_innerHash returns SHA256(0x01 || left || right), the Tendermint inner node hash.
func X_innerHash(left, right []byte) []byte {
	return tmhash.Sum(append(append([]byte{0x01}, left...), right...))
}

// X_hashFromByteSlices computes the Tendermint simple Merkle root.
// items is encoded as: [4-byte BE count][4-byte BE len][data]...
func X_hashFromByteSlices(encoded []byte) []byte {
	items, ok := decodeByteSlices(encoded)
	if !ok {
		return nil
	}
	return merkle.SimpleHashFromByteSlices(items)
}

// X_verifySimpleProof verifies a Tendermint SimpleProof.
// aunts is encoded as consecutive 32-byte hashes.
func X_verifySimpleProof(rootHash, leaf []byte, index, total int, aunts []byte) bool {
	if len(aunts)%32 != 0 {
		return false
	}
	// Copy each 32-byte chunk into its own allocation to avoid capacity-aliasing
	// bugs: tm2's innerHash uses append(left, right...) which mutates left's
	// backing array if len(left) < cap(left), corrupting later aunts.
	auntSlices := make([][]byte, len(aunts)/32)
	for i := range auntSlices {
		tmp := make([]byte, 32)
		copy(tmp, aunts[i*32:i*32+32])
		auntSlices[i] = tmp
	}
	sp := &merkle.SimpleProof{
		Total:    total,
		Index:    index,
		LeafHash: X_leafHash(leaf),
		Aunts:    auntSlices,
	}
	return sp.Verify(rootHash, leaf) == nil
}

// decodeByteSlices decodes [4-byte BE count][4-byte BE len][data]... encoding.
func decodeByteSlices(b []byte) ([][]byte, bool) {
	if len(b) < 4 {
		return nil, false
	}
	count := int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])
	b = b[4:]
	// Each item consumes at least 4 bytes (its length header), so a valid
	// encoding of `count` items requires len(b) >= 4*count. Reject before the
	// allocation so a tiny input cannot force a huge make([][]byte, count).
	// The division form avoids any multiplication overflow.
	if count > len(b)/4 {
		return nil, false
	}
	items := make([][]byte, count)
	for i := range count {
		if len(b) < 4 {
			return nil, false
		}
		n := int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])
		b = b[4:]
		if len(b) < n {
			return nil, false
		}
		items[i] = b[:n]
		b = b[n:]
	}
	return items, true
}
