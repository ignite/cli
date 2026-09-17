package testing

import (
	"math"
	"time"
)

func updateSeed() {
	seed := time.Now().UnixNano()
	_srand(seed)
}

func Test_UniformRand(t *T) {
	valueMap := make(map[uint64]int)
	maxIter := 1000

	for i := 0; i < maxIter; i++ {
		result := UniformRand()

		if result < 0 || result >= 32768 {
			t.Errorf("rand() = %v, want in range [0, 32767]", result)
		}

		if _, ok := valueMap[result]; ok {
			valueMap[result]++
		} else {
			valueMap[result] = 1
		}
	}

	lvm := len(valueMap)
	if lvm > maxIter || lvm == 0 {
		t.Errorf("len(valueMap) = %v, want <= %v", lvm, maxIter)
	}
}

func Test_nrand(t *T) {
	sum := 0.0
	for i := 0; i < 1000; i++ {
		result := nrand()
		sum += result
	}
	avg := sum / float64(1000)

	// The value returned by nrand() should be close to 0
	// on average for large sample sizes.
	// The expectation of the standard deviation should be
	// close to 1 for large sample sizes.
	if math.Abs(avg) > 0.1 { // can be adjusted based on sample size
		t.Errorf("nrand() = %v, want in range [-0.1, 0.1]", avg)
	}
}

func Test_GenerateRandomBool(t *T) {
	updateSeed()

	for _, bias := range []float64{0, 0.5, 1} {
		trueCount, falseCount := 0, 0
		sampleSize := 1000

		for i := 0; i < sampleSize; i++ {
			result := GenerateRandomBool(bias)
			if result {
				trueCount++
			} else {
				falseCount++
			}
		}

		if trueCount == 0 || falseCount == 0 {
			t.Errorf("Bias = %v, trueCount = %v, falseCount = %v, want both > 0", bias, trueCount, falseCount)
		}

		if bias < 0 || bias > 1 {
			t.Errorf("Bias = %v, want in range [0, 1]", bias)
		}
	}
}

func TestRandRange(t *T) {
	nums := make(map[uint64]int)
	for i := 0; i < 1000; i++ {
		res := randRange(0, 10)
		if res < 0 || res > 10 {
			t.Errorf("gerandRangenerateRange() = %v, want in range [0, 9]", res)
		}

		if _, ok := nums[res]; ok {
			nums[res]++
		} else {
			nums[res] = 1
		}
	}

	if len(nums) != 11 {
		t.Errorf("len(nums) = %v, want in range [0, 10]", len(nums))
	}
}
