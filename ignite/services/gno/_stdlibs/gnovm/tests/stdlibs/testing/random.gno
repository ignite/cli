package testing

import "math"

// Internal state for the random number generator.
var x uint64 = 42

// UniformRand generates a uniformly distributed random number.
// It uses the linear congrential generator method to produce the random number.
// and the result is in the range from 0 to m-1. here, m is 32768.
// To produce random number in [0, m-1], repeat this method as many times as needed.
// [1] https://en.wikipedia.org/wiki/Linear_congruential_generator
func UniformRand() uint64 {
	var a uint64 = 950213
	var c uint64 = 12345
	var m uint64 = 32768
	x = x*a + c
	return (x >> 16) % m
}

// _srand function sets the seed for the random number generator.
// This function provides an initial starting point for the sequence of random numbers.
func _srand(seed int64) {
	x = uint64(seed)
}

// nrand function generates a number approximating a normal distribution[1].
// It uses the Central Limit Theorem[2] by summing multiple uniformly distributed random numbers
// to approximate a normal distribution.
//
// y = Sum(k=1, K) (x_k - K/2) / sqrt(K/12)
//
// Here, K is some integer ans x_k are uniformly distributed numbers,
// even for K as small as 10, the approximation is quite good.
// [1] https://en.wikipedia.org/wiki/Normal_distribution
// [2] https://en.wikipedia.org/wiki/Central_limit_theorem
func nrand() float64 {
	var i, K uint64 = 0, 10
	var m uint64 = 32768
	var y float64 = 0

	for i = 0; i < K; i++ {
		y += float64(UniformRand()) / float64(m)
	}
	y = (y - float64(K)/2) / math.Sqrt(float64(K)/12)
	return y
}

// randRange generates a random integer between min and max (inclusive).
// This function leverages the UniformRand function to generate a random number in a specified range.
// Note: max should be greater than min.
func randRange(min, max int) uint64 {
	_min := uint64(min)
	_max := uint64(max)
	if _min >= _max {
		return _min
	}

	rangeSize := _max - _min + 1
	// adjust UniformRand to fit into our range.
	return _min + (UniformRand() % rangeSize)
}

func GenerateRandomBool(bias float64) bool {
	// Modify to use fuzz's random function for generating boolean with bias
	if bias < 0 || bias > 1 {
		panic("bias should be in the range [0, 1]")
	}
	// Convert fuzz's normalized range random float [-1, 1] to [0, 1]
	res := (nrand() + 1) / 2
	return res > bias
}
