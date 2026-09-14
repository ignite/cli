// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package subtle

// XORBytes computes x[i] ^ y[i] for all i < n = min(len(x), len(y)),
// returning buf, the resulting bytes slice.
//
// XXX: In the original implementation, XORBytes also checked for memory overlap
// between the destination and the input slices. Since buf is always freshly
// allocated here, overlap with x or y is impossible and the check is omitted.
func XORBytes(x, y []byte) []byte {
	n := min(len(x), len(y))
	buf := make([]byte, n)
	xorBytes(buf, x, y, n)
	return buf
}

// XORBytesUnsafe sets dst[i] = x[i] ^ y[i] for all i < n = min(len(x), len(y)), returning n, the number of bytes written to dst.
// If dst does not have length at least n, XORBytes panics without writing anything to dst.
//
// WARNING: dst must not overlap with x or y unless it is exactly the same slice.
// Otherwise, the behavior may be surprising and incorrect because overlap checks were removed.
func XORBytesUnsafe(dst, x, y []byte) int {
	n := min(len(x), len(y))
	if n == 0 {
		return 0
	}
	if n > len(dst) {
		panic("subtle.XORBytes: dst too short")
	}

	xorBytes(dst, x, y, n)
	return n
}

// xorBytes XORs bytes from x and y into dst.
// All slices must have length >= n.
func xorBytes(dst, x, y []byte, n int) {
	for i := 0; i < n; i++ {
		dst[i] = x[i] ^ y[i]
	}
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
