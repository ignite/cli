package sha256

const Size = 32

// Sum returns the SHA-256 checksum of the data.
func Sum256(data []byte) [Size]byte { return sum256(data) }

func sum256(data []byte) [32]byte // injected
