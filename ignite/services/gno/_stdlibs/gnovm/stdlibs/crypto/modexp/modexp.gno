// Package modexp exposes arbitrary-precision modular exponentiation as a
// native binding, matching the semantics of the EIP-198 (0x05) precompile.
//
// Use this when a gno contract needs a single modular reduction or a modular
// inverse over a large prime — both reducible to this primitive without
// pulling in a full big-integer stdlib.
package modexp

// MaxOperandLen is the largest accepted length, in bytes, of each of base, exp
// and modulus, matching the limit EIP-7823 places on the EVM precompile.
// 8192-bit operands cover RSA up to 8192 bits, BN254 field reduction, and
// modular inverse over an 8192-bit prime.
const MaxOperandLen = 1024

// ModExp returns (base^exp) mod modulus. Inputs are big-endian byte slices.
// The output has len(modulus) bytes, left-padded with zeros.
//
// Matches EIP-198 semantics:
//   - If modulus is zero or empty, returns len(modulus) zero bytes.
//   - If exp is zero or empty, returns 1 (or 0 if modulus is 1).
//   - Inputs are taken as non-negative big-endian integers with leading zeros
//     tolerated.
//
// Operands are length-limited, as in EIP-7823: if any of base, exp or modulus
// is longer than MaxOperandLen, ModExp returns an empty slice without computing
// anything. Callers handling operands they did not size themselves should check
// len(result) rather than assume len(modulus) bytes came back.
//
// Gas scales on bits(exp) × words(modulus)², so a call costs roughly the work it
// actually causes: a small exponent against a wide modulus is cheap, while
// wide-exponent-and-wide-modulus calls are expensive enough that only a handful
// fit in one block.
func ModExp(base, exp, modulus []byte) []byte {
	// Checked here, not only in the native, so an oversized operand is refused
	// before the native call: the dispatcher converts every operand from Gno to
	// Go before the native body runs, so rejecting on the Go side still pays for
	// an allocate-and-copy of each one.
	if len(base) > MaxOperandLen || len(exp) > MaxOperandLen || len(modulus) > MaxOperandLen {
		return nil
	}
	return modExp(base, exp, modulus)
}

func modExp(base, exp, modulus []byte) []byte // injected
