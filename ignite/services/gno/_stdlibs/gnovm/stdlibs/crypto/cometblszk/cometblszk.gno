// Package cometblszk is a pure-gno port of Union's CometBLS Groth16 verifier
// (evm/contracts/lib/CometblsZKVerifier.sol + clients/CometblsClient.sol).
//
// Unlike crypto/cometbls — which delegates the full verification to a single
// native call — this package keeps the verification logic in gno and leans on
// thin native primitives that mirror the EVM precompiles (crypto/bn254,
// crypto/keccak256, crypto/modexp). The design trades a bit of per-step
// overhead for auditable gno code paths and a minimal native surface.
//
// The public API is:
//
//	VerifyZKP(chainID, trustedValidatorsHash, header, zkp) error
//
// The light header is encoded with the same bespoke fixed-width layout as
// crypto/cometbls (see EncodedLightHeaderSize). The proof is 384 bytes (the
// same layout as the Rust cometbls-groth16-verifier).
package cometblszk

import (
	"crypto/sha256"
	"errors"
)

// EncodedLightHeaderSize is the fixed byte length of the bespoke LightHeader
// encoding, shared with crypto/cometbls:
//
//	8 height (BE i64) | 8 time_seconds (BE i64) | 4 time_nanos (BE i32)
//	| 32 validators_hash | 32 next_validators_hash | 32 app_hash
const EncodedLightHeaderSize = 8 + 8 + 4 + 32 + 32 + 32

// ExpectedProofSize is the byte length of a serialized proof:
//
//	A(64) | B(128) | C(64) | ProofCommitment(64) | ProofCommitmentPoK(64) = 384.
//
// G2 points are packed with the imaginary coefficient first
// (x_imag|x_real|y_imag|y_real), matching EIP-197.
const ExpectedProofSize = 64 + 128 + 64 + 64 + 64

// LightHeader mirrors cometbls-light-client-types::LightHeader. ValidatorsHash,
// NextValidatorsHash, AppHash must each be 32 bytes.
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

// VerifyZKP verifies a CometBLS Groth16 proof, doing the bulk of the work in
// gno and delegating only to bn254/keccak256/modexp/sha256 native primitives.
func VerifyZKP(chainID string, trustedValidatorsHash []byte, header LightHeader, zkp []byte) error {
	if len(chainID) > 31 {
		return errors.New("cometblszk: chain id exceeds 31 bytes")
	}
	if len(trustedValidatorsHash) != 32 {
		return errors.New("cometblszk: trusted validators hash must be 32 bytes")
	}
	if header.Height < 0 {
		return errors.New("cometblszk: invalid height")
	}
	if header.TimeSeconds < 0 || header.TimeNanos < 0 {
		return errors.New("cometblszk: invalid timestamp")
	}
	if len(header.ValidatorsHash) != 32 || len(header.NextValidatorsHash) != 32 || len(header.AppHash) != 32 {
		return errors.New("cometblszk: invalid header hash length")
	}
	if len(zkp) != ExpectedProofSize {
		return errors.New("cometblszk: invalid proof length")
	}

	// Slice the proof into its logical pieces. All slices alias zkp — this
	// package never mutates them.
	a := zkp[0:64]      // A point in G1
	b := zkp[64:192]    // B point in G2 (imag-first)
	c := zkp[192:256]   // C point in G1
	pc := zkp[256:320]  // ProofCommitment in G1
	pok := zkp[320:384] // ProofCommitmentPoK in G1

	// Compute the two scalar public inputs.
	commitmentHash := hashToField(pc)
	inputsHashBytes := inputsHash(chainID, trustedValidatorsHash, header)

	// publicInputMSM: CONSTANT + pc + inputsHash*PUB_0 + commitmentHash*PUB_1
	msm, ok := publicInputMSM(pc, inputsHashBytes, commitmentHash)
	if !ok {
		return errors.New("cometblszk: MSM failed")
	}

	// Pedersen POK pairing check.
	if !verifyProofCommitmentPOK(pc, pok) {
		return errors.New("cometblszk: invalid pok")
	}

	// Groth16 pairing check.
	if !verifyProof(a, b, c, msm) {
		return errors.New("cometblszk: invalid proof")
	}
	return nil
}

// inputsHash mirrors CometblsClient.sol's sha256-based inputsHash, with the
// most significant byte zeroed so the value fits in F_r.
func inputsHash(chainID string, trustedValidatorsHash []byte, h LightHeader) []byte {
	var buf [32]byte
	payload := make([]byte, 0, 32*8)

	// chainID: left-padded with zeros to 32 bytes.
	copy(buf[:], make([]byte, 32))
	copy(buf[32-len(chainID):], []byte(chainID))
	payload = append(payload, buf[:]...)

	payload = append(payload, u256FromUint64(uint64(h.Height))...)
	payload = append(payload, u256FromUint64(uint64(h.TimeSeconds))...)
	payload = append(payload, u256FromUint64(uint64(h.TimeNanos))...)
	payload = append(payload, h.ValidatorsHash...)
	payload = append(payload, h.NextValidatorsHash...)
	payload = append(payload, h.AppHash...)
	payload = append(payload, trustedValidatorsHash...)

	digest := sha256.Sum256(payload)
	digest[0] = 0
	return digest[:]
}

func u256FromUint64(v uint64) []byte {
	out := make([]byte, 32)
	putUint64BE(out[24:32], v)
	return out
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
