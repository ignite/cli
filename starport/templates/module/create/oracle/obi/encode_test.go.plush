package obi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type ID uint8

type Inner struct {
	A ID    `obi:"a"`
	B uint8 `obi:"b"`
}

type SimpleData struct {
	X uint8 `obi:"x"`
	Y uint8 `obi:"y"`
}

type ExampleData struct {
	Symbol string  `obi:"symbol"`
	Px     uint64  `obi:"px"`
	In     Inner   `obi:"in"`
	Arr    []int16 `obi:"arr"`
}

type InvalidStruct struct {
	IsBool bool
}

func TestEncodeBytes(t *testing.T) {
	require.Equal(t, MustEncode(ExampleData{
		Symbol: "BTC",
		Px:     9000,
		In: Inner{
			A: 1,
			B: 2,
		},
		Arr: []int16{10, 11},
	}), []byte{0x0, 0x0, 0x0, 0x3, 0x42, 0x54, 0x43, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x23, 0x28, 0x1, 0x2, 0x0, 0x0, 0x0, 0x2, 0x0, 0xa, 0x0, 0xb})
}
func TestEncodeBytesMulti(t *testing.T) {
	require.Equal(t, MustEncode(SimpleData{X: 1, Y: 2}, ExampleData{
		Symbol: "BTC",
		Px:     9000,
		In: Inner{
			A: 1,
			B: 2,
		},
		Arr: []int16{10, 11},
	}), []byte{0x1, 0x2, 0x0, 0x0, 0x0, 0x3, 0x42, 0x54, 0x43, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x23, 0x28, 0x1, 0x2, 0x0, 0x0, 0x0, 0x2, 0x0, 0xa, 0x0, 0xb})
}

func TestEncodeStructFail(t *testing.T) {
	invalid := InvalidStruct{
		IsBool: true,
	}
	require.PanicsWithError(t, "obi: unsupported value type: bool", func() { MustEncode(invalid) })
}

// Uint8
func TestEncodeBytesUint8(t *testing.T) {
	num := uint8(123)
	require.Equal(t, []byte{num}, MustEncode(num))
}

func TestEncodeBytesAliasUint8(t *testing.T) {
	type ID uint8
	num := uint8(123)
	id := ID(num)
	require.Equal(t, []byte{num}, MustEncode(id))
}

// Uint16
func TestEncodeBytesUint16(t *testing.T) {
	num := uint16(123)
	require.Equal(t, []byte{0x00, 0x7b}, MustEncode(num))
}

func TestEncodeBytesAliasUint16(t *testing.T) {
	type ID uint16
	num := uint16(123)
	id := ID(num)
	require.Equal(t, []byte{0x0, 0x7b}, MustEncode(id))
}

// Uint32
func TestEncodeBytesUint32(t *testing.T) {
	num := uint32(123)
	require.Equal(t, []byte{0x0, 0x0, 0x0, 0x7b}, MustEncode(num))
}

func TestEncodeBytesAliasUint32(t *testing.T) {
	type ID uint32
	num := uint32(123)
	id := ID(num)
	require.Equal(t, []byte{0x0, 0x0, 0x0, 0x7b}, MustEncode(id))
}

// Uint64
func TestEncodeBytesUint64(t *testing.T) {
	num := uint64(123)
	require.Equal(t, []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7b}, MustEncode(num))
}

func TestEncodeBytesAliasUint64(t *testing.T) {
	type ID uint64
	num := uint64(123)
	id := ID(num)
	require.Equal(t, []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7b}, MustEncode(id))
}

// Int8
func TestEncodeBytesInt8(t *testing.T) {
	num := int8(-123)
	require.Equal(t, []byte{0x85}, MustEncode(num))
}

func TestEncodeBytesAliasInt8(t *testing.T) {
	type ID int8
	num := int8(-123)
	id := ID(num)
	require.Equal(t, []byte{0x85}, MustEncode(id))
}

// Int16
func TestEncodeBytesInt16(t *testing.T) {
	num := int16(-123)
	require.Equal(t, []byte{0xff, 0x85}, MustEncode(num))
}

func TestEncodeBytesAliasInt16(t *testing.T) {
	type ID int16
	num := int16(-123)
	id := ID(num)
	require.Equal(t, []byte{0xff, 0x85}, MustEncode(id))
}

// Int32
func TestEncodeBytesInt32(t *testing.T) {
	num := int32(-123)
	require.Equal(t, []byte{0xff, 0xff, 0xff, 0x85}, MustEncode(num))
}

func TestEncodeBytesAliasInt32(t *testing.T) {
	type ID int32
	num := int32(-123)
	id := ID(num)
	require.Equal(t, []byte{0xff, 0xff, 0xff, 0x85}, MustEncode(id))
}

// Int64
func TestEncodeBytesInt64(t *testing.T) {
	num := int64(-123)
	require.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x85}, MustEncode(num))
}

func TestEncodeBytesAliasInt64(t *testing.T) {
	type ID int64
	num := int32(-123)
	id := ID(num)
	require.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x85}, MustEncode(id))
}

func TestEncodeString(t *testing.T) {
	testString := "hello alice and bob"
	expectedEncodeBytes := []byte{0x00, 0x00, 0x00, 0x13, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x61, 0x6c, 0x69, 0x63, 0x65, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x62, 0x6f, 0x62}
	require.Equal(t, []byte(expectedEncodeBytes), MustEncode(testString))
}

func TestEncodeSlice(t *testing.T) {
	testSlice := []int32{1, 2, 3, 4, 5, 6}
	expectedEncodeBytes := []byte{0x0, 0x0, 0x0, 0x6, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x6}
	require.Equal(t, expectedEncodeBytes, MustEncode(testSlice))
}

func TestEncodeSliceFail(t *testing.T) {
	testSlice := []bool{true, false, true, true}
	require.PanicsWithError(t, "obi: unsupported value type: bool", func() { MustEncode(testSlice) })
}

func TestEncodeByteArray(t *testing.T) {
	testByteArray := []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6}
	expectedEncodeByte := []byte{0x0, 0x0, 0x0, 0x6, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6}
	require.Equal(t, expectedEncodeByte, MustEncode(testByteArray))
}

func TestEncodeNotSupported(t *testing.T) {
	notSupportBool := true
	byteArray, err := Encode(notSupportBool)
	require.EqualError(t, err, "obi: unsupported value type: bool")
	require.Nil(t, byteArray)
}

func TestEncodeNotSupport(t *testing.T) {
	notSupportBool := true
	require.PanicsWithError(t, "obi: unsupported value type: bool", func() { MustEncode(notSupportBool) })
}
