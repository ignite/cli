// Package cometbls verifies CometBLS Groth16 proofs over BN254 via a native
// host binding. The heavy lifting (pairings, scalar multiplications, hash
// derivation) happens in Go; the gno side only encodes the LightHeader and
// forwards to verifyZKP.
package cometbls

import "errors"

// EncodedLightHeaderSize is the byte length of the bespoke fixed-width
// LightHeader encoding consumed by VerifyZKP:
//
//	8 height (BE i64) | 8 time_seconds (BE i64) | 4 time_nanos (BE i32)
//	| 32 validators_hash | 32 next_validators_hash | 32 app_hash
const EncodedLightHeaderSize = 8 + 8 + 4 + 32 + 32 + 32

// LightHeader mirrors cometbls-light-client-types::LightHeader. All byte fields
// must be exactly 32 bytes — VerifyZKP rejects anything else.
type LightHeader struct {
	Height             int64
	TimeSeconds        int64
	TimeNanos          int32
	ValidatorsHash     []byte
	NextValidatorsHash []byte
	AppHash            []byte
}

// EncodeLightHeader produces the canonical wire encoding accepted by VerifyZKP.
// Returns nil if any hash field is not exactly 32 bytes.
func EncodeLightHeader(h LightHeader) []byte {
	if len(h.ValidatorsHash) != 32 || len(h.NextValidatorsHash) != 32 || len(h.AppHash) != 32 {
		return nil
	}
	out := make([]byte, EncodedLightHeaderSize)
	putUint64BE(out[0:8], uint64(h.Height))
	putUint64BE(out[8:16], uint64(h.TimeSeconds))
	putUint32BE(out[16:20], uint32(h.TimeNanos))
	copy(out[20:52], h.ValidatorsHash)
	copy(out[52:84], h.NextValidatorsHash)
	copy(out[84:116], h.AppHash)
	return out
}

// VerifyZKP verifies a CometBLS proof against the given trusted validators
// hash and light header. chainID must be at most 31 bytes, trustedValidatorsHash
// must be 32 bytes. Returns nil on successful verification.
func VerifyZKP(chainID string, trustedValidatorsHash []byte, header LightHeader, zkp []byte) error {
	encoded := EncodeLightHeader(header)
	if encoded == nil {
		return errors.New("cometbls: invalid header encoding length")
	}
	msg := verifyZKP(chainID, trustedValidatorsHash, encoded, zkp)
	if msg == "" {
		return nil
	}
	return errors.New("cometbls: " + msg)
}

func putUint64BE(dst []byte, v uint64) {
	dst[0] = byte(v >> 56)
	dst[1] = byte(v >> 48)
	dst[2] = byte(v >> 40)
	dst[3] = byte(v >> 32)
	dst[4] = byte(v >> 24)
	dst[5] = byte(v >> 16)
	dst[6] = byte(v >> 8)
	dst[7] = byte(v)
}

func putUint32BE(dst []byte, v uint32) {
	dst[0] = byte(v >> 24)
	dst[1] = byte(v >> 16)
	dst[2] = byte(v >> 8)
	dst[3] = byte(v)
}

func verifyZKP(chainID string, trustedValidatorsHash []byte, headerEncoded []byte, zkp []byte) string // injected
