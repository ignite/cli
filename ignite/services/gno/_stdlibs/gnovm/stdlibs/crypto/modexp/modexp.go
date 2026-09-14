package modexp

import "math/big"

// maxOperandLen bounds every operand of X_modExp at 1024 bytes (8192 bits),
// matching the limit EIP-7823 sets on the equivalent EVM precompile.
//
// The gas charge, not this cap, is what prices the exponentiation: the row for
// this native in gnovm/stdlibs/native_gas.go scales on bits(exp) ×
// words(modulus)², the quantity big.Int.Exp's runtime is proportional to.
//
// The cap covers what that charge does not:
//
//   - The envelope the slope was fit over. "This fit was only validated to 1024
//     bytes" is per-native knowledge and belongs here, not in the gas table.
//   - base, which no term of the metric reads. Its cost — the reduction of base
//     modulo modulus — is bounded by this cap rather than priced. That is only
//     sound because the bound is small; do not raise it without adding a slope
//     on len(base).
//
// The length check is duplicated in ModExp (modexp.gno) and that copy is the
// load-bearing one: the dispatcher converts every operand to Go before this
// function is entered, so rejecting here still pays for an allocate-and-copy of
// each operand. Rejecting on the Gno side skips the native call entirely. This
// copy cannot be fully bounded by any length check, here or there — Gno2GoValue
// copies a byte slice's *capacity* (gnovm/pkg/gnolang/gonative.go), so a
// zero-length slice with a large capacity still copies it all. That is
// pre-existing dispatcher behaviour shared by every []byte native, not specific
// to modexp.
//
// 1024 bytes covers every real use of modular exponentiation — RSA up to 8192
// bits, BN254 field reduction (crypto/cometblszk's hashToField uses 32-byte
// operands), and modular inverse over an 8192-bit prime via Fermat, where the
// exponent is the same width as the modulus. Nothing legitimate needs more: for
// gcd(a,n)=1, a^e ≡ a^(e mod φ(n)), so an exponent wider than the modulus is
// reducible by construction. EIP-7823's survey of Ethereum mainnet from April
// 2018 to January 2025 found no successful call exceeding 513 bytes in any
// field. The gas ceiling also binds well before this cap does: at the shipped
// slope a symmetric 1024-byte call costs ~526M gas, about a sixth of a
// 3B-gas block.
const maxOperandLen = 1024

// X_modExp mirrors EIP-198's MODEXP precompile.
//
// Operands longer than maxOperandLen yield nil, matching how the sibling
// EIP-precompile natives (crypto/bn254, crypto/merkle) report inputs outside
// their contract.
func X_modExp(base, exp, modulus []byte) []byte {
	if len(base) > maxOperandLen || len(exp) > maxOperandLen || len(modulus) > maxOperandLen {
		return nil
	}
	out := make([]byte, len(modulus))
	m := new(big.Int).SetBytes(modulus)
	if m.Sign() == 0 {
		return out
	}
	b := new(big.Int).SetBytes(base)
	e := new(big.Int).SetBytes(exp)
	r := new(big.Int).Exp(b, e, m)
	r.FillBytes(out)
	return out
}
