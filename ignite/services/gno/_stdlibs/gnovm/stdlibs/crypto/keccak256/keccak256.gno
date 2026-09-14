// Package keccak256 exposes the Legacy Keccak-256 hash (EVM's keccak256) as a
// native binding. The padding differs from NIST SHA3-256, so this is offered
// as a separate package rather than under crypto/sha3.
package keccak256

const Size = 32

// Sum256 returns the Keccak-256 hash of data.
func Sum256(data []byte) [Size]byte { return sum256(data) }

func sum256(data []byte) [32]byte // injected
