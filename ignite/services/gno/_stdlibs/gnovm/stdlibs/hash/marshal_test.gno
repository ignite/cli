// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test that the hashes in the standard library implement
// BinaryMarshaler, BinaryUnmarshaler,
// and lock in the current representations.

package hash_test

import (
	"bytes"
	"encoding"
	"encoding/hex"
	"hash"
	"hash/adler32"
	"testing"
)

func fromHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

var marshalTests = []struct {
	name   string
	new    func() hash.Hash
	golden []byte
}{
	{"adler32", func() hash.Hash { return adler32.New() }, fromHex("61646c01460a789d")},
}

func TestMarshalHash(t *testing.T) {
	for _, tt := range marshalTests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 256)
			for i := range buf {
				buf[i] = byte(i)
			}

			h := tt.new()
			h.Write(buf[:256])
			sum := h.Sum(nil)

			h2 := tt.new()
			h3 := tt.new()
			const split = 249
			for i := 0; i < split; i++ {
				h2.Write(buf[i : i+1])
			}
			h2m, ok := h2.(encoding.BinaryMarshaler)
			if !ok {
				t.Fatalf("Hash does not implement MarshalBinary")
			}
			enc, err := h2m.MarshalBinary()
			if err != nil {
				t.Fatalf("MarshalBinary: %v", err)
			}
			if !bytes.Equal(enc, tt.golden) {
				t.Errorf("MarshalBinary = %x, want %x", enc, tt.golden)
			}
			h3u, ok := h3.(encoding.BinaryUnmarshaler)
			if !ok {
				t.Fatalf("Hash does not implement UnmarshalBinary")
			}
			if err := h3u.UnmarshalBinary(enc); err != nil {
				t.Fatalf("UnmarshalBinary: %v", err)
			}
			h2.Write(buf[split:])
			h3.Write(buf[split:])
			sum2 := h2.Sum(nil)
			sum3 := h3.Sum(nil)
			if !bytes.Equal(sum2, sum) {
				t.Fatalf("Sum after MarshalBinary = %x, want %x", sum2, sum)
			}
			if !bytes.Equal(sum3, sum) {
				t.Fatalf("Sum after UnmarshalBinary = %x, want %x", sum3, sum)
			}
		})
	}
}
