// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math_test

import (
	"math"
	"testing"
)

func TestMaxUint(t *testing.T) {
	// XXX: https://github.com/gnolang/gno/issues/1085
	var v uint = math.MaxUint
	if v+1 != 0 {
		t.Errorf("MaxUint should wrap around to zero: %d", v+1)
	}
	if v := uint8(math.MaxUint8); v+1 != 0 {
		t.Errorf("MaxUint8 should wrap around to zero: %d", v+1)
	}
	if v := uint16(math.MaxUint16); v+1 != 0 {
		t.Errorf("MaxUint16 should wrap around to zero: %d", v+1)
	}
	if v := uint32(math.MaxUint32); v+1 != 0 {
		t.Errorf("MaxUint32 should wrap around to zero: %d", v+1)
	}
	var v1 uint64 = math.MaxUint64
	if v1+1 != 0 {
		t.Errorf("MaxUint64 should wrap around to zero: %d", v1+1)
	}
}

func TestMaxInt(t *testing.T) {
	if v := int(math.MaxInt); v+1 != math.MinInt {
		t.Errorf("MaxInt should wrap around to MinInt: %d", v+1)
	}
	if v := int8(math.MaxInt8); v+1 != math.MinInt8 {
		t.Errorf("MaxInt8 should wrap around to MinInt8: %d", v+1)
	}
	if v := int16(math.MaxInt16); v+1 != math.MinInt16 {
		t.Errorf("MaxInt16 should wrap around to MinInt16: %d", v+1)
	}
	if v := int32(math.MaxInt32); v+1 != math.MinInt32 {
		t.Errorf("MaxInt32 should wrap around to MinInt32: %d", v+1)
	}
	if v := int64(math.MaxInt64); v+1 != math.MinInt64 {
		t.Errorf("MaxInt64 should wrap around to MinInt64: %d", v+1)
	}
}
