// Package cometbls provides native bindings for verifying CometBLS Groth16
// proofs over BN254 from Gno contracts.
//
// The host-side verifier is a line-for-line Go port of the Rust
// cometbls-groth16-verifier crate (unionlabs/union: lib/cometbls-groth16-verifier/).
// Pre-negated verifying-key constants live in constants.go (generated from
// verifying_key.bin; see the union repo's cmd/gen).
//
// Gno contracts call VerifyZKP declared in cometbls.gno, which encodes the
// LightHeader using the bespoke fixed-width layout (see EncodedLightHeaderSize)
// and forwards to X_verifyZKP.
package cometbls

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math/big"

	bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"golang.org/x/crypto/sha3"
)

// Field and point sizes (match the Rust crate).
const (
	FqSize = 32
	G1Size = 2 * FqSize
	G2Size = 2 * G1Size

	// ExpectedProofSize is the byte length of a serialized proof:
	//   A(G1) | B(G2) | C(G1) | ProofCommitment(G1) | ProofCommitmentPoK(G1)
	ExpectedProofSize = G1Size + G2Size + G1Size + G1Size + G1Size

	// EncodedLightHeaderSize is the byte length of the bespoke LightHeader
	// encoding accepted by VerifyZKP:
	//   8 height (BE i64) | 8 seconds (BE i64) | 4 nanos (BE i32)
	//   | 32 validators_hash | 32 next_validators_hash | 32 app_hash
	EncodedLightHeaderSize = 8 + 8 + 4 + 32 + 32 + 32
)

// hmacOPad / hmacIPad are the pre-computed HMAC pads for the key "CometBLS"
// under Keccak-256 (block size 136). Hardcoded in the Rust crate; reproduced
// here so we never materialize the key outside of this package.
var (
	hmacOPad = []byte{
		0x1F, 0x33, 0x31, 0x39, 0x28, 0x1E, 0x10, 0x0F,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
		0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C, 0x5C,
	}
	hmacIPad = []byte{
		0x75, 0x59, 0x5B, 0x53, 0x42, 0x74, 0x7A, 0x65,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
		0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36, 0x36,
	}

	// primeRMinusOne = fr.Modulus() - 1, used by hash_to_field.
	primeRMinusOne = new(big.Int).Sub(fr.Modulus(), big.NewInt(1))
)

// Proof is the three canonical Groth16 points.
type Proof struct {
	A bn254.G1Affine
	B bn254.G2Affine
	C bn254.G1Affine
}

// ZKP bundles a Groth16 proof with the Pedersen commitment and its
// proof-of-knowledge, matching the layout expected by the CometBLS circuit.
type ZKP struct {
	Proof              Proof
	ProofCommitment    bn254.G1Affine
	ProofCommitmentPoK bn254.G1Affine
}

// LightHeader matches cometbls-light-client-types::LightHeader.
type LightHeader struct {
	Height             int64
	TimeSeconds        int64
	TimeNanos          int32
	ValidatorsHash     [32]byte
	NextValidatorsHash [32]byte
	AppHash            [32]byte
}

// Sentinel errors mirroring the Rust `Error` enum. Tests compare against these
// by identity using errors.Is.
var (
	ErrInvalidPublicInput = errors.New("invalid public input")
	ErrInvalidPoint       = errors.New("invalid point")
	ErrInvalidProof       = errors.New("invalid proof")
	ErrInvalidPok         = errors.New("invalid pok")
	ErrInvalidCommitment  = errors.New("invalid commitment")
	ErrInvalidRawProof    = errors.New("invalid raw proof")
	ErrInvalidHeight      = errors.New("invalid height")
	ErrInvalidTimestamp   = errors.New("invalid timestamp")
	ErrInvalidHeaderLen   = errors.New("invalid header encoding length")
	ErrInvalidChainIDLen  = errors.New("chain id must be at most 31 bytes")
)

// ParseZKP decodes a raw proof byte slice in big-endian layout.
//
// Layout (bytes):
//
//	0     .. G1                    = A
//	G1    .. G1+G2                 = B
//	G1+G2 .. 2*G1+G2               = C
//	2*G1+G2 .. 3*G1+G2             = ProofCommitment
//	3*G1+G2 .. 4*G1+G2             = ProofCommitmentPoK
func ParseZKP(raw []byte) (*ZKP, error) {
	if len(raw) != ExpectedProofSize {
		return nil, ErrInvalidRawProof
	}
	var zkp ZKP
	if _, err := zkp.Proof.A.SetBytes(raw[0:G1Size]); err != nil {
		return nil, ErrInvalidPoint
	}
	if _, err := zkp.Proof.B.SetBytes(raw[G1Size : G1Size+G2Size]); err != nil {
		return nil, ErrInvalidPoint
	}
	if _, err := zkp.Proof.C.SetBytes(raw[G1Size+G2Size : 2*G1Size+G2Size]); err != nil {
		return nil, ErrInvalidPoint
	}
	if _, err := zkp.ProofCommitment.SetBytes(raw[2*G1Size+G2Size : 3*G1Size+G2Size]); err != nil {
		return nil, ErrInvalidPoint
	}
	if _, err := zkp.ProofCommitmentPoK.SetBytes(raw[3*G1Size+G2Size : 4*G1Size+G2Size]); err != nil {
		return nil, ErrInvalidPoint
	}
	return &zkp, nil
}

// EncodeLightHeader produces the bespoke fixed-width encoding accepted by
// VerifyZKP. A matching encoder is provided gno-side in cometbls.gno.
func EncodeLightHeader(h LightHeader) []byte {
	out := make([]byte, EncodedLightHeaderSize)
	binary.BigEndian.PutUint64(out[0:8], uint64(h.Height))
	binary.BigEndian.PutUint64(out[8:16], uint64(h.TimeSeconds))
	binary.BigEndian.PutUint32(out[16:20], uint32(h.TimeNanos))
	copy(out[20:52], h.ValidatorsHash[:])
	copy(out[52:84], h.NextValidatorsHash[:])
	copy(out[84:116], h.AppHash[:])
	return out
}

// DecodeLightHeader parses the bespoke encoding, validating bounds.
func DecodeLightHeader(buf []byte) (LightHeader, error) {
	var h LightHeader
	if len(buf) != EncodedLightHeaderSize {
		return h, ErrInvalidHeaderLen
	}
	h.Height = int64(binary.BigEndian.Uint64(buf[0:8]))
	h.TimeSeconds = int64(binary.BigEndian.Uint64(buf[8:16]))
	h.TimeNanos = int32(binary.BigEndian.Uint32(buf[16:20]))
	copy(h.ValidatorsHash[:], buf[20:52])
	copy(h.NextValidatorsHash[:], buf[52:84])
	copy(h.AppHash[:], buf[84:116])
	if h.Height < 0 {
		return h, ErrInvalidHeight
	}
	if h.TimeSeconds < 0 || h.TimeNanos < 0 {
		return h, ErrInvalidTimestamp
	}
	return h, nil
}

// hmacKeccak computes HMAC(Keccak256, key="CometBLS", message) using the
// pre-computed opad/ipad constants.
func hmacKeccak(message []byte) [32]byte {
	inner := sha3.NewLegacyKeccak256()
	inner.Write(hmacIPad)
	inner.Write(message)
	innerDigest := inner.Sum(nil)

	outer := sha3.NewLegacyKeccak256()
	outer.Write(hmacOPad)
	outer.Write(innerDigest)
	var out [32]byte
	copy(out[:], outer.Sum(nil))
	return out
}

// hashToField = (HMAC_Keccak256(msg) mod (r - 1)) + 1, where r is the BN254
// scalar field modulus.
func hashToField(message []byte) fr.Element {
	h := hmacKeccak(message)
	n := new(big.Int).SetBytes(h[:])
	n.Mod(n, primeRMinusOne)
	n.Add(n, big.NewInt(1))
	var e fr.Element
	e.SetBigInt(n)
	return e
}

// hashCommitment matches Gnark's commitment hashing: serialize the G1 as
// 64 big-endian bytes (x || y) and apply hashToField.
func hashCommitment(p *bn254.G1Affine) fr.Element {
	var buf [64]byte
	xBytes := p.X.Bytes()
	yBytes := p.Y.Bytes()
	copy(buf[0:32], xBytes[:])
	copy(buf[32:64], yBytes[:])
	return hashToField(buf[:])
}

// PublicInputs computes the two scalar public inputs consumed by the circuit:
// the SHA-256-derived inputs hash (top byte zeroed to fit in F_r) and the
// commitment hash.
func PublicInputs(chainID string, trustedValidatorsHash [32]byte, header LightHeader, zkp *ZKP) ([2]fr.Element, error) {
	var out [2]fr.Element
	if len(chainID) > 31 {
		return out, ErrInvalidChainIDLen
	}
	if header.Height < 0 {
		return out, ErrInvalidHeight
	}
	if header.TimeSeconds < 0 || header.TimeNanos < 0 {
		return out, ErrInvalidTimestamp
	}

	h := sha256.New()

	// chain_id: left-padded with zeros to 32 bytes.
	pad := make([]byte, 32-len(chainID))
	h.Write(pad)
	h.Write([]byte(chainID))

	var u256 [32]byte
	putU64AsU256(&u256, uint64(header.Height))
	h.Write(u256[:])
	putU64AsU256(&u256, uint64(header.TimeSeconds))
	h.Write(u256[:])
	// Rust casts nanos (i32) → i64 → u64 → U256(BE). We mirror that: nanos is
	// non-negative so uint64 cast is identity on the low 64 bits.
	putU64AsU256(&u256, uint64(header.TimeNanos))
	h.Write(u256[:])

	h.Write(header.ValidatorsHash[:])
	h.Write(header.NextValidatorsHash[:])
	h.Write(header.AppHash[:])
	h.Write(trustedValidatorsHash[:])

	var digest [32]byte
	h.Sum(digest[:0])
	// Zero the most significant byte so the value fits in F_r.
	digest[0] = 0

	var inputsHashElem fr.Element
	inputsHashElem.SetBytes(digest[:])
	out[0] = inputsHashElem
	out[1] = hashCommitment(&zkp.ProofCommitment)
	return out, nil
}

// VerifyZKP verifies a CometBLS Groth16 proof. The header is accepted in the
// bespoke fixed-width encoding produced by EncodeLightHeader (see
// EncodedLightHeaderSize for the layout).
//
// Returns nil on success or one of the sentinel errors on failure.
func VerifyZKP(chainID string, trustedValidatorsHash []byte, headerEncoded []byte, zkpBytes []byte) error {
	if len(trustedValidatorsHash) != 32 {
		return ErrInvalidCommitment
	}
	var tvh [32]byte
	copy(tvh[:], trustedValidatorsHash)

	header, err := DecodeLightHeader(headerEncoded)
	if err != nil {
		return err
	}

	zkp, err := ParseZKP(zkpBytes)
	if err != nil {
		return err
	}

	inputs, err := PublicInputs(chainID, tvh, header, zkp)
	if err != nil {
		return err
	}

	// public_inputs_msm = GammaAbcG1[0] + ProofCommitment
	//                     + sum(inputs[i] * GammaAbcG1[i+1])
	var msm bn254.G1Affine
	msm.Add(&GammaAbcG1[0], &zkp.ProofCommitment)

	var scalar big.Int
	var term bn254.G1Affine
	for i := range NbPublicInputs {
		inputs[i].BigInt(&scalar)
		term.ScalarMultiplication(&GammaAbcG1[i+1], &scalar)
		msm.Add(&msm, &term)
	}

	// Pedersen Proof-of-Knowledge check:
	//   e(pc, PedersenG) * e(pok, PedersenGRootSigmaNeg) == 1
	okPok, err := bn254.PairingCheck(
		[]bn254.G1Affine{zkp.ProofCommitment, zkp.ProofCommitmentPoK},
		[]bn254.G2Affine{PedersenG, PedersenGRootSigmaNeg},
	)
	if err != nil {
		return ErrInvalidPok
	}
	if !okPok {
		return ErrInvalidPok
	}

	// Groth16 pairing check:
	//   e(A, B) * e(msm, GammaNegG2) * e(C, DeltaNegG2) * e(AlphaG1, BetaNegG2) == 1
	okProof, err := bn254.PairingCheck(
		[]bn254.G1Affine{zkp.Proof.A, msm, zkp.Proof.C, AlphaG1},
		[]bn254.G2Affine{zkp.Proof.B, GammaNegG2, DeltaNegG2, BetaNegG2},
	)
	if err != nil {
		return ErrInvalidProof
	}
	if !okProof {
		return ErrInvalidProof
	}
	return nil
}

// putU64AsU256 writes v as a 32-byte big-endian integer (high 24 bytes zero).
func putU64AsU256(dst *[32]byte, v uint64) {
	for i := range dst {
		dst[i] = 0
	}
	binary.BigEndian.PutUint64(dst[24:32], v)
}
