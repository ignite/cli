package math

// Float32bits returns the IEEE 754 binary representation of f, with the sign
// bit of f and the result in the same bit position.
// Float32bits(Float32frombits(x)) == x.
func Float32bits(f float32) uint32 // injected

// Float32frombits returns the floating-point number corresponding to the IEEE
// 754 binary representation b, with the sign bit of b and the result in the
// same bit position. Float32frombits(Float32bits(x)) == x.
func Float32frombits(b uint32) float32 // injected

// Float64bits returns the IEEE 754 binary representation of f, with the sign
// bit of f and the result in the same bit position.
// Float64bits(Float64frombits(x)) == x.
func Float64bits(f float64) uint64 // injected

// Float64frombits returns the floating-point number corresponding to the IEEE
// 754 binary representation b, with the sign bit of b and the result in the
// same bit position. Float64frombits(Float64bits(x)) == x.
func Float64frombits(b uint64) float64 // injected
