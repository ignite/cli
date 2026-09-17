package ed25519

// Verify returns true if the signature is valid for the message and public key.
func Verify(publicKey []byte, message []byte, signature []byte) bool {
	return verify(publicKey, message, signature)
}

func verify(publicKey []byte, message []byte, signature []byte) bool // injected
