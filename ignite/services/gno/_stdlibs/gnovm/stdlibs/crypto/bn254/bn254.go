package bn254

import (
	"math/big"

	bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
)

var fpModulus = fp.Modulus()

// X_g1Add mirrors EIP-196's ECADD precompile.
// Per EIP-196, short inputs are right-padded with zeros to 128 bytes; excess bytes are ignored.
//
// The excess is normally already gone: bn254.gno's G1Add truncates to 128
// bytes before calling in, because this native is priced flat and the
// dispatcher's Gno→Go copy is not. Keep that check there — the one below only
// covers Go callers, by which point the copy has happened.
func X_g1Add(input []byte) []byte {
	padded := make([]byte, 128)
	if len(input) < 128 {
		copy(padded, input)
	} else {
		copy(padded, input[:128])
	}
	p1, ok := parseG1(padded[0:64])
	if !ok {
		return nil
	}
	p2, ok := parseG1(padded[64:128])
	if !ok {
		return nil
	}
	var sum bn254.G1Affine
	sum.Add(&p1, &p2)
	return marshalG1(&sum)
}

// X_g1Mul mirrors EIP-196's ECMUL precompile.
//
// bn254.gno's G1Mul applies the same length check before calling in; see
// X_g1Add for why that matters.
func X_g1Mul(input []byte) []byte {
	if len(input) != 96 {
		return nil
	}
	p, ok := parseG1(input[0:64])
	if !ok {
		return nil
	}
	scalar := new(big.Int).SetBytes(input[64:96])
	var out bn254.G1Affine
	out.ScalarMultiplication(&p, scalar)
	return marshalG1(&out)
}

// X_pairingCheck mirrors EIP-197's ECPAIRING precompile. An empty input
// reports "1" (product of zero pairings), matching the spec.
func X_pairingCheck(input []byte) []byte {
	if len(input)%192 != 0 {
		return nil
	}
	n := len(input) / 192
	if n == 0 {
		out := make([]byte, 32)
		out[31] = 1
		return out
	}
	g1s := make([]bn254.G1Affine, n)
	g2s := make([]bn254.G2Affine, n)
	for i := range n {
		off := i * 192
		g1, ok := parseG1(input[off : off+64])
		if !ok {
			return nil
		}
		g2, ok := parseG2(input[off+64 : off+192])
		if !ok {
			return nil
		}
		g1s[i] = g1
		g2s[i] = g2
	}
	ok, err := bn254.PairingCheck(g1s, g2s)
	if err != nil {
		return nil
	}
	out := make([]byte, 32)
	if ok {
		out[31] = 1
	}
	return out
}

// parseG1 interprets 64 bytes (x|y, BE) as an affine G1 point. The encoding
// (0, 0) is accepted as the point at infinity. All other points must satisfy
// the curve equation; coordinates must be reduced modulo p.
func parseG1(buf []byte) (bn254.G1Affine, bool) {
	var p bn254.G1Affine
	x := new(big.Int).SetBytes(buf[0:32])
	y := new(big.Int).SetBytes(buf[32:64])
	if x.Cmp(fpModulus) >= 0 || y.Cmp(fpModulus) >= 0 {
		return p, false
	}
	if x.Sign() == 0 && y.Sign() == 0 {
		// Identity. gnark-crypto represents infinity with zero coordinates;
		// IsOnCurve returns true for it.
		p.X.SetZero()
		p.Y.SetZero()
		return p, true
	}
	p.X.SetBigInt(x)
	p.Y.SetBigInt(y)
	if !p.IsOnCurve() {
		return p, false
	}
	return p, true
}

// parseG2 interprets 128 bytes as an affine G2 point using EIP-197's
// imaginary-first Fp2 layout: x_imag | x_real | y_imag | y_real.
func parseG2(buf []byte) (bn254.G2Affine, bool) {
	var p bn254.G2Affine
	xImag := new(big.Int).SetBytes(buf[0:32])
	xReal := new(big.Int).SetBytes(buf[32:64])
	yImag := new(big.Int).SetBytes(buf[64:96])
	yReal := new(big.Int).SetBytes(buf[96:128])
	for _, v := range []*big.Int{xImag, xReal, yImag, yReal} {
		if v.Cmp(fpModulus) >= 0 {
			return p, false
		}
	}
	if xImag.Sign() == 0 && xReal.Sign() == 0 && yImag.Sign() == 0 && yReal.Sign() == 0 {
		p.X.SetZero()
		p.Y.SetZero()
		return p, true
	}
	p.X.A1.SetBigInt(xImag)
	p.X.A0.SetBigInt(xReal)
	p.Y.A1.SetBigInt(yImag)
	p.Y.A0.SetBigInt(yReal)
	if !p.IsOnCurve() {
		return p, false
	}
	// Reject points not in the correct subgroup — the EIP-197 precompile
	// rejects these too.
	if !p.IsInSubGroup() {
		return p, false
	}
	return p, true
}

// marshalG1 emits the 64-byte big-endian x|y encoding.
func marshalG1(p *bn254.G1Affine) []byte {
	out := make([]byte, 64)
	if p.IsInfinity() {
		return out
	}
	xBytes := p.X.Bytes()
	yBytes := p.Y.Bytes()
	copy(out[0:32], xBytes[:])
	copy(out[32:64], yBytes[:])
	return out
}
