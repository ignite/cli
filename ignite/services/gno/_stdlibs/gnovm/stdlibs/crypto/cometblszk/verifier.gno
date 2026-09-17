package cometblszk

import "crypto/bn254"

// publicInputMSM computes the Groth16 public-input linear combination:
//
//	out = CONSTANT + ProofCommitment + pub0*PUB_0 + pub1*PUB_1
//
// Mirrors CometblsZKVerifier.sol's publicInputMSM assembly block, but without
// the scalar < R check — our scalars come from hashToField (guaranteed < R)
// and from a SHA-256 digest with the top byte zeroed (~2^248 < R).
//
// Returns the 64-byte x|y encoding of the resulting G1 point.
func publicInputMSM(proofCommitment []byte, pub0, pub1 []byte) ([]byte, bool) {
	// accum = CONSTANT + ProofCommitment
	accum := g1AddConcat(constantX, constantY, proofCommitment[0:32], proofCommitment[32:64])
	if accum == nil {
		return nil, false
	}

	// accum += pub0 * PUB_0
	term0 := g1MulConcat(pub0X, pub0Y, pub0)
	if term0 == nil {
		return nil, false
	}
	accum = g1Add(accum, term0)
	if accum == nil {
		return nil, false
	}

	// accum += pub1 * PUB_1
	term1 := g1MulConcat(pub1X, pub1Y, pub1)
	if term1 == nil {
		return nil, false
	}
	accum = g1Add(accum, term1)
	if accum == nil {
		return nil, false
	}
	return accum, true
}

// verifyProofCommitmentPOK runs the Pedersen proof-of-knowledge pairing check:
//
//	e(ProofCommitment, PedersenG) * e(ProofCommitmentPoK, PedersenGRootSigmaNeg) == 1
//
// Mirrors CometblsZKVerifier.sol's verifyProofCommitmentPOK.
func verifyProofCommitmentPOK(proofCommitment, proofCommitmentPoK []byte) bool {
	input := make([]byte, 0, 2*bn254.PairingChunkSize)
	input = append(input, proofCommitment...)
	input = append(input, pedersenGX1...)
	input = append(input, pedersenGX0...)
	input = append(input, pedersenGY1...)
	input = append(input, pedersenGY0...)
	input = append(input, proofCommitmentPoK...)
	input = append(input, pedersenGRootSigmaNegX1...)
	input = append(input, pedersenGRootSigmaNegX0...)
	input = append(input, pedersenGRootSigmaNegY1...)
	input = append(input, pedersenGRootSigmaNegY0...)
	return pairingEqualsOne(input)
}

// verifyProof runs the Groth16 pairing check:
//
//	e(A, B) * e(C, -delta) * e(alpha, -beta) * e(msm, -gamma) == 1
//
// a is 64 bytes (A.x|A.y), b is 128 bytes (B.x_imag|B.x_real|B.y_imag|B.y_real),
// c is 64 bytes, msm is the 64-byte output of publicInputMSM.
//
// Mirrors CometblsZKVerifier.sol's verifyProof.
func verifyProof(a, b, c, msm []byte) bool {
	input := make([]byte, 0, 4*bn254.PairingChunkSize)
	// e(A, B)
	input = append(input, a...)
	input = append(input, b...)
	// e(C, -delta)
	input = append(input, c...)
	input = append(input, deltaNegX1...)
	input = append(input, deltaNegX0...)
	input = append(input, deltaNegY1...)
	input = append(input, deltaNegY0...)
	// e(alpha, -beta)
	input = append(input, alphaX...)
	input = append(input, alphaY...)
	input = append(input, betaNegX1...)
	input = append(input, betaNegX0...)
	input = append(input, betaNegY1...)
	input = append(input, betaNegY0...)
	// e(msm, -gamma)
	input = append(input, msm...)
	input = append(input, gammaNegX1...)
	input = append(input, gammaNegX0...)
	input = append(input, gammaNegY1...)
	input = append(input, gammaNegY0...)
	return pairingEqualsOne(input)
}

// g1AddConcat builds a 128-byte ECADD input from four 32-byte coordinates
// and returns the 64-byte sum, or nil on precompile failure.
func g1AddConcat(x1, y1, x2, y2 []byte) []byte {
	buf := make([]byte, 0, bn254.G1AddInputSize)
	buf = append(buf, x1...)
	buf = append(buf, y1...)
	buf = append(buf, x2...)
	buf = append(buf, y2...)
	return bn254.G1Add(buf)
}

// g1Add sums two G1 points already packed as 64-byte (x|y) buffers.
func g1Add(p1, p2 []byte) []byte {
	buf := make([]byte, 0, bn254.G1AddInputSize)
	buf = append(buf, p1...)
	buf = append(buf, p2...)
	return bn254.G1Add(buf)
}

// g1MulConcat builds a 96-byte ECMUL input and returns the product.
func g1MulConcat(x, y, scalar []byte) []byte {
	buf := make([]byte, 0, bn254.G1MulInputSize)
	buf = append(buf, x...)
	buf = append(buf, y...)
	buf = append(buf, scalar...)
	return bn254.G1Mul(buf)
}

// pairingEqualsOne returns true iff the pairing precompile reports the product
// of pairings equals the identity in GT.
func pairingEqualsOne(input []byte) bool {
	out := bn254.PairingCheck(input)
	if len(out) != bn254.PairingOutputSize {
		return false
	}
	// EIP-197 returns 32-byte 0x...01 on success, 0x...00 otherwise.
	for i := 0; i < 31; i++ {
		if out[i] != 0 {
			return false
		}
	}
	return out[31] == 1
}
