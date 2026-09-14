package merkle

// LeafHash returns SHA256(0x00 || leaf).
func LeafHash(leaf []byte) []byte { return leafHash(leaf) }

func leafHash(leaf []byte) []byte // injected

// InnerHash returns SHA256(0x01 || left || right).
func InnerHash(left, right []byte) []byte { return innerHash(left, right) }

func innerHash(left, right []byte) []byte // injected

// HashFromByteSlices computes the Tendermint simple Merkle root.
// items is encoded as: [4-byte BE count][4-byte BE len][data]...
func HashFromByteSlices(encoded []byte) []byte { return hashFromByteSlices(encoded) }

func hashFromByteSlices(encoded []byte) []byte // injected

// VerifySimpleProof verifies a Tendermint SimpleProof.
// aunts is a flat concatenation of consecutive 32-byte sibling hashes.
func VerifySimpleProof(rootHash, leaf []byte, index, total int, aunts []byte) bool {
	return verifySimpleProof(rootHash, leaf, index, total, aunts)
}

func verifySimpleProof(rootHash, leaf []byte, index, total int, aunts []byte) bool // injected
