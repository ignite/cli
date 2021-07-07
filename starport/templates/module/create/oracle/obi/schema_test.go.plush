package obi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type EmptySchema struct{}

type NoOBITagStruct struct {
	NoOBITag string `"noOBItag"` //missing obi
}

type NotSupportedStruct struct {
	IsValid bool   `obi:"isValid"`
	Test    string `obi:"test"`
}

type AllData struct {
	NumUint8  uint8  `obi:"numUint8"`
	NumUint16 uint16 `obi:"numUint16"`
	NumUint32 uint32 `obi:"numUint32"`
	NumUint64 uint64 `obi:"numUint64"`
	NumInt8   int8   `obi:"numInt8"`
	NumInt16  int16  `obi:"numInt16"`
	NumInt32  int32  `obi:"numInt32"`
	NumInt64  int64  `obi:"numInt64"`
}

type Uint8ID uint8
type Uint16ID uint16
type Uint32ID uint32
type Uint64ID uint64
type Int8ID int8
type Int16ID int16
type Int32ID int32
type Int64ID int64

type AllDataAlias struct {
	NumUint8  Uint8ID  `obi:"numUint8"`
	NumUint16 Uint16ID `obi:"numUint16"`
	NumUint32 Uint32ID `obi:"numUint32"`
	NumUint64 Uint64ID `obi:"numUint64"`
	NumInt8   Int8ID   `obi:"numInt8"`
	NumInt16  Int16ID  `obi:"numInt16"`
	NumInt32  Int32ID  `obi:"numInt32"`
	NumInt64  Int64ID  `obi:"numInt64"`
}

type ByteArrayStruct struct {
	ByteArray []byte `obi:"byteArray"`
}

func TestSchema(t *testing.T) {
	require.Equal(t, "{symbol:string,px:u64,in:{a:u8,b:u8},arr:[i16]}", MustGetSchema(ExampleData{}))
}

func TestEmptySchemaFail(t *testing.T) {
	require.PanicsWithError(t, "obi: empty struct is not supported", func() { MustGetSchema(EmptySchema{}) })
}

func TestMissingOBISchemaFail(t *testing.T) {
	require.PanicsWithError(t, "obi: no obi tag found for field NoOBITag of NoOBITagStruct", func() { MustGetSchema(NoOBITagStruct{}) })
}

func TestUnsupportedTypeFail(t *testing.T) {
	require.PanicsWithError(t, "obi: unsupported value type: bool", func() { MustGetSchema(NotSupportedStruct{}) })
}

func TestSchemaSupportedNumberTypeSuccess(t *testing.T) {
	require.Equal(t, "{numUint8:u8,numUint16:u16,numUint32:u32,numUint64:u64,numInt8:i8,numInt16:i16,numInt32:i32,numInt64:i64}", MustGetSchema(AllData{}))
}

func TestSchemaSupportedNumberAliasTypeSuccess(t *testing.T) {
	require.Equal(t, "{numUint8:u8,numUint16:u16,numUint32:u32,numUint64:u64,numInt8:i8,numInt16:i16,numInt32:i32,numInt64:i64}", MustGetSchema(AllDataAlias{}))
}

func TestSchemaInvalidSliceFail(t *testing.T) {
	invalidSlice := []bool{false, false, true, true}
	require.PanicsWithError(t, "obi: unsupported value type: bool", func() { MustGetSchema(invalidSlice) })
}

func TestSchemaByteArraySuccess(t *testing.T) {
	byteArray := []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6}
	require.Equal(t, "bytes", MustGetSchema(byteArray))
}

func TestSchemaByteArrayInStructSuccess(t *testing.T) {
	require.Equal(t, "{byteArray:bytes}", MustGetSchema(ByteArrayStruct{}))
}
