package obi

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

// Encode uses obi encoding scheme to encode the given input into bytes.
func encodeImpl(v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Uint8:
		return EncodeUnsigned8(uint8(rv.Uint())), nil
	case reflect.Uint16:
		return EncodeUnsigned16(uint16(rv.Uint())), nil
	case reflect.Uint32:
		return EncodeUnsigned32(uint32(rv.Uint())), nil
	case reflect.Uint64:
		return EncodeUnsigned64(uint64(rv.Uint())), nil
	case reflect.Int8:
		return EncodeSigned8(int8(rv.Int())), nil
	case reflect.Int16:
		return EncodeSigned16(int16(rv.Int())), nil
	case reflect.Int32:
		return EncodeSigned32(int32(rv.Int())), nil
	case reflect.Int64:
		return EncodeSigned64(int64(rv.Int())), nil
	case reflect.String:
		return EncodeString(rv.String()), nil
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return EncodeBytes(rv.Bytes()), nil
		}

		res := EncodeUnsigned32(uint32(rv.Len()))
		for idx := 0; idx < rv.Len(); idx++ {
			each, err := Encode(rv.Index(idx).Interface())
			if err != nil {
				return nil, err
			}
			res = append(res, each...)
		}
		return res, nil
	case reflect.Struct:
		res := []byte{}
		for idx := 0; idx < rv.NumField(); idx++ {
			each, err := Encode(rv.Field(idx).Interface())
			if err != nil {
				return nil, err
			}
			res = append(res, each...)
		}
		return res, nil
	default:
		return nil, fmt.Errorf("obi: unsupported value type: %s", rv.Kind())
	}
}

// Encode uses obi encoding scheme to encode the given input(s) into bytes.
func Encode(v ...interface{}) ([]byte, error) {
	res := []byte{}
	for _, each := range v {
		encoded, err := encodeImpl(each)
		if err != nil {
			return nil, err
		}
		res = append(res, encoded...)
	}
	return res, nil
}

// MustEncode uses obi encoding scheme to encode the given input into bytes. Panics on error.
func MustEncode(v ...interface{}) []byte {
	res, err := Encode(v...)
	if err != nil {
		panic(err)
	}
	return res
}

// EncodeUnsigned8 takes an `uint8` variable and encodes it into a byte array
func EncodeUnsigned8(v uint8) []byte {
	return []byte{v}
}

// EncodeUnsigned16 takes an `uint16` variable and encodes it into a byte array
func EncodeUnsigned16(v uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, v)
	return bytes
}

// EncodeUnsigned32 takes an `uint32` variable and encodes it into a byte array
func EncodeUnsigned32(v uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, v)
	return bytes
}

// EncodeUnsigned64 takes an `uint64` variable and encodes it into a byte array
func EncodeUnsigned64(v uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, v)
	return bytes
}

// EncodeSigned8 takes an `int8` variable and encodes it into a byte array
func EncodeSigned8(v int8) []byte {
	return EncodeUnsigned8(uint8(v))
}

// EncodeSigned16 takes an `int16` variable and encodes it into a byte array
func EncodeSigned16(v int16) []byte {
	return EncodeUnsigned16(uint16(v))
}

// EncodeSigned32 takes an `int32` variable and encodes it into a byte array
func EncodeSigned32(v int32) []byte {
	return EncodeUnsigned32(uint32(v))
}

// EncodeSigned64 takes an `int64` variable and encodes it into a byte array
func EncodeSigned64(v int64) []byte {
	return EncodeUnsigned64(uint64(v))
}

// EncodeBytes takes a `[]byte` variable and encodes it into a byte array
func EncodeBytes(v []byte) []byte {
	return append(EncodeUnsigned32(uint32(len(v))), v...)
}

// EncodeString takes a `string` variable and encodes it into a byte array
func EncodeString(v string) []byte {
	return append(EncodeUnsigned32(uint32(len(v))), []byte(v)...)
}
