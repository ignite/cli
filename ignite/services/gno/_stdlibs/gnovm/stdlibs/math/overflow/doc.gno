// Package overflow offers overflow-checked integer arithmetic operations
// for all signed and unsigned integer types. Each of the operations returns a
// result,bool combination. The modulo % operation is not present, as it is
// always safe.
//
// This implementation is based on its Go equivalent in tm2/pkg/overflow, using
// code generation in place of generics. All functions are named matching the
// following regex:
//
//	(Add|Sub|Mul|Div)(u)?(8|16|32|64)?(p)?
//
// The "u" identifies a function which works on an unsigned value. The number
// identifies the bit size, and can be omitted for the int and uint types.
// The trailing "p" identifies the panicking versions of each function.
package overflow
