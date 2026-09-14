package sha256

import "crypto/sha256"

func X_sum256(data []byte) [32]byte {
	return sha256.Sum256(data)
}
