package obi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Uint8
func TestDecodeUint8(t *testing.T) {
	var actual uint8
	byteArray := []byte{0x7b}
	MustDecode(byteArray, &actual)
	require.Equal(t, uint8(123), actual)
}

func TestDecodeAliasUint8(t *testing.T) {
	type ID uint8
	var actual ID
	byteArray := []byte{0x7b}
	MustDecode(byteArray, &actual)
	require.Equal(t, ID(123), actual)
}

func TestDecodeUint8Fail(t *testing.T) {
	var actual uint8
	byteArray := []byte{}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

// Uint16
func TestDecodeUint16(t *testing.T) {
	var actual uint16
	byteArray := []byte{0x00, 0x7b}
	MustDecode(byteArray, &actual)
	require.Equal(t, uint16(123), actual)
}

func TestDecodeAliasUint16(t *testing.T) {
	type ID uint16
	var actual ID
	byteArray := []byte{0x00, 0x7b}
	MustDecode(byteArray, &actual)
	require.Equal(t, ID(123), actual)
}

func TestDecodeUint16Fail(t *testing.T) {
	var actual uint16
	byteArray := []byte{}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

// Uint32
func TestDecodeUint32(t *testing.T) {
	var actual uint32
	byteArray := []byte{0x0, 0x0, 0x0, 0x7b}
	MustDecode(byteArray, &actual)
	require.Equal(t, uint32(123), actual)
}

func TestDecodeAliasUint32(t *testing.T) {
	type ID uint32
	var actual ID
	byteArray := []byte{0x0, 0x0, 0x0, 0x7b}
	MustDecode(byteArray, &actual)
	require.Equal(t, ID(123), actual)
}

func TestDecodeUint32Fail(t *testing.T) {
	var actual uint32
	byteArray := []byte{}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

// Uint64
func TestDecodeUint64(t *testing.T) {
	var actual uint64
	byteArray := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7b}
	MustDecode(byteArray, &actual)
	require.Equal(t, uint64(123), actual)
}

func TestDecodeAliasUint64(t *testing.T) {
	type ID uint64
	var actual ID
	byteArray := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7b}
	MustDecode(byteArray, &actual)
	require.Equal(t, ID(123), actual)
}

func TestDecodeUint64Fail(t *testing.T) {
	var actual uint64
	byteArray := []byte{}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

// Int8
func TestDecodeInt8(t *testing.T) {
	var actual int8
	byteArray := []byte{0x85}
	MustDecode(byteArray, &actual)
	require.Equal(t, int8(-123), actual)
}

func TestDecodeAliasInt8(t *testing.T) {
	type ID int8
	var actual ID
	byteArray := []byte{0x85}
	MustDecode(byteArray, &actual)
	require.Equal(t, ID(-123), actual)
}

func TestDecodeInt8Fail(t *testing.T) {
	var actual int8
	byteArray := []byte{}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

// Int16
func TestDecodeInt16(t *testing.T) {
	var actual int16
	byteArray := []byte{0xff, 0x85}
	MustDecode(byteArray, &actual)
	require.Equal(t, int16(-123), actual)
}

func TestDecodeAliasInt16(t *testing.T) {
	type ID int16
	var actual ID
	byteArray := []byte{0xff, 0x85}
	MustDecode(byteArray, &actual)
	require.Equal(t, ID(-123), actual)
}

func TestDecodeInt16Fail(t *testing.T) {
	var actual int16
	byteArray := []byte{}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

// Int32
func TestDecodeInt32(t *testing.T) {
	var actual int32
	byteArray := []byte{0xff, 0xff, 0xff, 0x85}
	MustDecode(byteArray, &actual)
	require.Equal(t, int32(-123), actual)
}

func TestDecodeAliasInt32(t *testing.T) {
	type ID int32
	var actual ID
	byteArray := []byte{0xff, 0xff, 0xff, 0x85}
	MustDecode(byteArray, &actual)
	require.Equal(t, ID(-123), actual)
}

func TestDecodeInt32Fail(t *testing.T) {
	var actual int32
	byteArray := []byte{}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

// Int64
func TestDecodeInt64(t *testing.T) {
	var actual int64
	byteArray := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x85}
	MustDecode(byteArray, &actual)
	require.Equal(t, int64(-123), actual)
}

func TestDecodeAliasInt64(t *testing.T) {
	type ID int64
	var actual ID
	byteArray := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x85}
	MustDecode(byteArray, &actual)
	require.Equal(t, ID(-123), actual)
}

func TestDecodeInt64Fail(t *testing.T) {
	var actual int64
	byteArray := []byte{}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

func TestDecodeString(t *testing.T) {
	var actual string
	byteArray := []byte{0x00, 0x00, 0x00, 0x13, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x61, 0x6c, 0x69, 0x63, 0x65, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x62, 0x6f, 0x62}
	MustDecode(byteArray, &actual)
	require.Equal(t, "hello alice and bob", actual)
}

func TestDecodeEmptyStringFail(t *testing.T) {
	var actual string
	byteArray := []byte{}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

func TestDecodeStringOutOfRangeFail(t *testing.T) {
	var actual string
	byteArray := []byte{0x00, 0x00, 0x00, 0x13}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

func TestDecodeSlice(t *testing.T) {
	var actual []int32
	byteArray := []byte{0x0, 0x0, 0x0, 0x6, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x6}
	MustDecode(byteArray, &actual)
	require.Equal(t, []int32{1, 2, 3, 4, 5, 6}, actual)
}

func TestDecodeSliceFail(t *testing.T) {
	var actual []int32
	byteArray := []byte{0x6, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

func TestDecodeSliceOutOfRangeFail(t *testing.T) {
	var actual []int32
	byteArray := []byte{0x6, 0x0, 0x0}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

func TestDecodeStruct(t *testing.T) {
	var actual ExampleData
	byteArray := []byte{0x0, 0x0, 0x0, 0x3, 0x42, 0x54, 0x43, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x23, 0x28, 0x1, 0x2, 0x0, 0x0, 0x0, 0x2, 0x0, 0xa, 0x0, 0xb}
	MustDecode(byteArray, &actual)
	require.Equal(t,
		ExampleData{
			Symbol: "BTC",
			Px:     9000,
			In: Inner{
				A: 1,
				B: 2,
			},
			Arr: []int16{10, 11},
		},
		actual,
	)
}

func TestDecodeMultipleValues(t *testing.T) {
	var a, b int8
	byteArray := []byte{0x20, 0x21}
	MustDecode(byteArray, &a, &b)
	require.Equal(t, int8(32), a)
	require.Equal(t, int8(33), b)
}

func TestDecodeStructFail(t *testing.T) {
	var actual ExampleData
	byteArray := []byte{0x0, 0x0, 0x0, 0x6, 0x1, 0x0, 0x0, 0x0}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

func TestDecodeStructOutOfRangeFail(t *testing.T) {
	var actual []ExampleData
	byteArray := []byte{0x0, 0x0, 0x0, 0x6}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

func TestDecodeByteArray(t *testing.T) {
	var actual []byte
	byteArray := []byte{0x0, 0x0, 0x0, 0x6, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6}
	MustDecode(byteArray, &actual)
	require.Equal(t, []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6}, actual)
}

func TestDecodeByteFail(t *testing.T) {
	var actual []byte
	byteArray := []byte{0x0, 0x6}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

func TestDecodeByteArrayOutOfRangeFail(t *testing.T) {
	var actual []byte
	byteArray := []byte{0x0, 0x0, 0x0, 0x6}
	require.PanicsWithError(t, "obi: out of range", func() { MustDecode(byteArray, &actual) })
}

func TestDecodeIntoNonPointer(t *testing.T) {
	var actual []byte
	byteArray := []byte{0x0, 0x0, 0x0, 0x6, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6}
	require.PanicsWithError(t, "obi: decode into non-ptr type", func() { MustDecode(byteArray, actual) })
}

func TestUnsupportedType(t *testing.T) {
	var actual bool
	byteArray := []byte{0x6, 0x0, 0x0, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6}
	require.PanicsWithError(t, "obi: unsupported value type: bool", func() { MustDecode(byteArray, &actual) })
}

func TestNotAllDataConsumed(t *testing.T) {
	var actual []byte
	byteArray := []byte{0x0, 0x0, 0x0, 0x6, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0x10}
	require.PanicsWithError(t, "obi: not all data was consumed while decoding", func() { MustDecode(byteArray, &actual) })
}
