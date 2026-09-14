package math

import "math"

func Float32bits(f float32) uint32     { return math.Float32bits(f) }
func Float32frombits(b uint32) float32 { return math.Float32frombits(b) }
func Float64bits(f float64) uint64     { return math.Float64bits(f) }
func Float64frombits(b uint64) float64 { return math.Float64frombits(b) }
