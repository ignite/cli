package ed25519

import (
	"crypto/ed25519"
)

func X_verify(publicKey []byte, message []byte, signature []byte) bool {
	return ed25519.Verify(publicKey, message, signature)
}
