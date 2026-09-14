package binary_test

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestLittleEndianUint16(t *testing.T) {
	b := []byte{0xe8, 0x03, 0xd0, 0x07}
	x1 := binary.LittleEndian.Uint16(b[0:])
	x2 := binary.LittleEndian.Uint16(b[2:])

	if x1 != 0x03e8 || x2 != 0x07d0 {
		t.Fatal("littleEndian.Uint16 failed")
	}
}

func TestLittleEndianPutUint16(t *testing.T) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b[0:], 0x03e8)
	binary.LittleEndian.PutUint16(b[2:], 0x07d0)

	if !bytes.Equal(b, []byte{0xe8, 0x03, 0xd0, 0x07}) {
		t.Fatal("littleEndian.PutUint16 failed")
	}
}

func TestLittleEndianAppendUint16(t *testing.T) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b[0:], 0x03e8)
	binary.LittleEndian.PutUint16(b[2:], 0x07d0)

	v := uint16(0x0bb8)
	b = binary.LittleEndian.AppendUint16(b, v)

	if !bytes.Equal(b, []byte{0xe8, 0x03, 0xd0, 0x07, 0xb8, 0x0b}) {
		t.Fatal("littleEndian.AppendUint16 failed")
	}
}

func TestLittleEndianUint32(t *testing.T) {
	b := []byte{0xe8, 0x03, 0xd0, 0x07}
	x := binary.LittleEndian.Uint32(b[0:])
	if x != 0x07d003e8 {
		t.Fatal("littleEndian.Uint32 failed")
	}
}

func TestLittleEndianPutUint32(t *testing.T) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b[0:], 0x07d003e8)

	if !bytes.Equal(b, []byte{0xe8, 0x03, 0xd0, 0x07}) {
		t.Fatal("littleEndian.PutUint32 failed")
	}
}

func TestLittleEndianAppendUint32(t *testing.T) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b[0:], 0x07d003e8)

	v := uint32(0x0bb807d0)
	b = binary.LittleEndian.AppendUint32(b, v)

	if !bytes.Equal(b, []byte{0xe8, 0x03, 0xd0, 0x07, 0xd0, 0x07, 0xb8, 0x0b}) {
		t.Fatal("littleEndian.AppendUint32 failed")
	}
}

func TestLittleEndianUint64(t *testing.T) {
	b := []byte{0xe8, 0x03, 0xd0, 0x07, 0xe9, 0xf5, 0x14, 0xfe}
	x := binary.LittleEndian.Uint64(b[0:])
	if x != 0xfe14f5e907d003e8 {
		t.Fatal("littleEndian.Uint64 failed")
	}
}

func TestLittleEndianPutUint64(t *testing.T) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b[0:], 0xfe14f5e907d003e8)

	if !bytes.Equal(b, []byte{0xe8, 0x03, 0xd0, 0x07, 0xe9, 0xf5, 0x14, 0xfe}) {
		t.Fatal("littleEndian.PutUint64 failed")
	}
}

func TestLittleEndianAppendUint64(t *testing.T) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b[0:], 0xfe14f5e907d003e8)

	v := uint64(0x0bb807d0d007e9f5)
	b = binary.LittleEndian.AppendUint64(b, v)

	if !bytes.Equal(b, []byte{0xe8, 0x03, 0xd0, 0x07, 0xe9, 0xf5, 0x14, 0xfe, 0xf5, 0xe9, 0x07, 0xd0, 0xd0, 0x07, 0xb8, 0x0b}) {
		t.Fatal("littleEndian.AppendUint64 failed")
	}
}

func TestBigEndianUint16(t *testing.T) {
	b := []byte{0xe8, 0x03, 0xd0, 0x07}
	x1 := binary.BigEndian.Uint16(b[0:])
	x2 := binary.BigEndian.Uint16(b[2:])

	if x1 != 0xe803 || x2 != 0xd007 {
		t.Fatal("bigEndian.Uint16 failed")
	}
}

func TestBigEndianPutUint16(t *testing.T) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint16(b[0:], 0x03e8)
	binary.BigEndian.PutUint16(b[2:], 0x07d0)

	if !bytes.Equal(b, []byte{0x03, 0xe8, 0x07, 0xd0}) {
		t.Fatal("bigEndian.PutUint16 failed")
	}
}

func TestBigEndianAppendUint16(t *testing.T) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint16(b[0:], 0x03e8)
	binary.BigEndian.PutUint16(b[2:], 0x07d0)

	v := uint16(0x0bb8)
	b = binary.BigEndian.AppendUint16(b, v)

	if !bytes.Equal(b, []byte{0x03, 0xe8, 0x07, 0xd0, 0x0b, 0xb8}) {
		t.Fatal("BigEndian.AppendUint16 failed")
	}
}

func TestBigEndianUint32(t *testing.T) {
	b := []byte{0xe8, 0x03, 0xd0, 0x07}
	x := binary.BigEndian.Uint32(b[0:])
	if x != 0xe803d007 {
		t.Fatal("bigEndian.Uint32 failed")
	}
}

func TestBigEndianPutUint32(t *testing.T) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b[0:], 0x07d003e8)

	if !bytes.Equal(b, []byte{0x07, 0xd0, 0x03, 0xe8}) {
		t.Fatal("bigEndian.PutUint32 failed")
	}
}

func TestBigEndianAppendUint32(t *testing.T) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b[0:], 0x07d003e8)

	v := uint32(0x0bb807d0)
	b = binary.BigEndian.AppendUint32(b, v)

	if !bytes.Equal(b, []byte{0x07, 0xd0, 0x03, 0xe8, 0x0b, 0xb8, 0x07, 0xd0}) {
		t.Fatal("bigEndian.AppendUint32 failed")
	}
}

func TestBigEndianUint64(t *testing.T) {
	b := []byte{0xe8, 0x03, 0xd0, 0x07, 0xe9, 0xf5, 0x14, 0xfe}
	x := binary.BigEndian.Uint64(b[0:])
	if x != 0xe803d007e9f514fe {
		t.Fatal("bigEndian.Uint64 failed")
	}
}

func TestBigEndianPutUint64(t *testing.T) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b[0:], 0xfe14f5e907d003e8)

	if !bytes.Equal(b, []byte{0xfe, 0x14, 0xf5, 0xe9, 0x07, 0xd0, 0x03, 0xe8}) {
		t.Fatal("bigEndian.PutUint64 failed")
	}
}

func TestbigEndianAppendUint64(t *testing.T) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b[0:], 0xfe14f5e907d003e8)

	v := uint64(0x0bb807d0d007e9f5)
	b = binary.BigEndian.AppendUint64(b, v)

	if !bytes.Equal(b, []byte{0xfe, 0x14, 0xf5, 0xe9, 0x07, 0xd0, 0x03, 0xe8, 0x0b, 0xb8, 0x07, 0xd0, 0xd0, 0x07, 0xe9, 0xf5}) {
		t.Fatal("bigEndian.AppendUint64 failed")
	}
}
