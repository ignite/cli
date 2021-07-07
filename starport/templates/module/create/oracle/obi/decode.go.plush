package obi

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

func decodeImpl(data []byte, v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil, errors.New("obi: decode into non-ptr type")
	}
	ev := rv.Elem()
	switch ev.Kind() {
	case reflect.Uint8:
		val, rem, err := DecodeUnsigned8(data)
		ev.SetUint(uint64(val))
		return rem, err
	case reflect.Uint16:
		val, rem, err := DecodeUnsigned16(data)
		ev.SetUint(uint64(val))
		return rem, err
	case reflect.Uint32:
		val, rem, err := DecodeUnsigned32(data)
		ev.SetUint(uint64(val))
		return rem, err
	case reflect.Uint64:
		val, rem, err := DecodeUnsigned64(data)
		ev.SetUint(uint64(val))
		return rem, err
	case reflect.Int8:
		val, rem, err := DecodeSigned8(data)
		ev.SetInt(int64(val))
		return rem, err
	case reflect.Int16:
		val, rem, err := DecodeSigned16(data)
		ev.SetInt(int64(val))
		return rem, err
	case reflect.Int32:
		val, rem, err := DecodeSigned32(data)
		ev.SetInt(int64(val))
		return rem, err
	case reflect.Int64:
		val, rem, err := DecodeSigned64(data)
		ev.SetInt(int64(val))
		return rem, err
	case reflect.String:
		val, rem, err := DecodeString(data)
		ev.SetString(val)
		return rem, err
	case reflect.Slice:
		if ev.Type().Elem().Kind() == reflect.Uint8 {
			val, rem, err := DecodeBytes(data)
			ev.SetBytes(val)
			return rem, err
		}
		length, rem, err := DecodeUnsigned32(data)
		if err != nil {
			return nil, err
		}
		slice := reflect.MakeSlice(ev.Type(), int(length), int(length))
		for idx := 0; idx < int(length); idx++ {
			var err error
			rem, err = decodeImpl(rem, slice.Index(idx).Addr().Interface())
			if err != nil {
				return nil, err
			}
		}
		ev.Set(slice)
		return rem, nil
	case reflect.Struct:
		rem := data
		for idx := 0; idx < ev.NumField(); idx++ {
			var err error
			rem, err = decodeImpl(rem, ev.Field(idx).Addr().Interface())
			if err != nil {
				return nil, err
			}
		}
		return rem, nil
	default:
		return nil, fmt.Errorf("obi: unsupported value type: %s", ev.Kind())
	}
}

// Decode uses obi encoding scheme to decode the given input(s).
func Decode(data []byte, v ...interface{}) error {
	var err error
	rem := data
	for _, each := range v {
		rem, err = decodeImpl(rem, each)
		if err != nil {
			return err
		}
	}
	if len(rem) != 0 {
		return errors.New("obi: not all data was consumed while decoding")
	}
	return nil
}

// MustDecode uses obi encoding scheme to decode the given input. Panics on error.
func MustDecode(data []byte, v ...interface{}) {
	err := Decode(data, v...)
	if err != nil {
		panic(err)
	}
}

// DecodeUnsigned16 decodes the input bytes into `uint8` and returns the remaining bytes.
func DecodeUnsigned8(data []byte) (uint8, []byte, error) {
	if len(data) < 1 {
		return 0, nil, errors.New("obi: out of range")
	}
	return data[0], data[1:], nil
}

// DecodeUnsigned16 decodes the input bytes into `uint16` and returns the remaining bytes.
func DecodeUnsigned16(data []byte) (uint16, []byte, error) {
	if len(data) < 2 {
		return 0, nil, errors.New("obi: out of range")
	}
	return binary.BigEndian.Uint16(data[:2]), data[2:], nil
}

// DecodeUnsigned32 decodes the input bytes into `uint32` and returns the remaining bytes.
func DecodeUnsigned32(data []byte) (uint32, []byte, error) {
	if len(data) < 4 {
		return 0, nil, errors.New("obi: out of range")
	}
	return binary.BigEndian.Uint32(data[:4]), data[4:], nil
}

// DecodeUnsigned64 decodes the input bytes into `uint64` and returns the remaining bytes.
func DecodeUnsigned64(data []byte) (uint64, []byte, error) {
	if len(data) < 8 {
		return 0, nil, errors.New("obi: out of range")
	}
	return binary.BigEndian.Uint64(data[:8]), data[8:], nil
}

// DecodeSigned8 decodes the input bytes into `uint64` and returns the remaining bytes.
func DecodeSigned8(data []byte) (int8, []byte, error) {
	unsigned, rem, err := DecodeUnsigned8(data)
	return int8(unsigned), rem, err
}

// DecodeSigned16 decodes the input bytes into `uint64` and returns the remaining bytes.
func DecodeSigned16(data []byte) (int16, []byte, error) {
	unsigned, rem, err := DecodeUnsigned16(data)
	return int16(unsigned), rem, err
}

// DecodeSigned32 decodes the input bytes into `uint64` and returns the remaining bytes.
func DecodeSigned32(data []byte) (int32, []byte, error) {
	unsigned, rem, err := DecodeUnsigned32(data)
	return int32(unsigned), rem, err
}

// DecodeSigned64 decodes the input bytes into `uint64` and returns the remaining bytes.
func DecodeSigned64(data []byte) (int64, []byte, error) {
	unsigned, rem, err := DecodeUnsigned64(data)
	return int64(unsigned), rem, err
}

// DecodeBytes decodes the input bytes and returns bytes result and the remaining bytes.
func DecodeBytes(data []byte) ([]byte, []byte, error) {
	length, rem, err := DecodeUnsigned32(data)
	if err != nil {
		return nil, nil, err
	}
	if uint32(len(rem)) < length {
		return nil, nil, errors.New("obi: out of range")
	}
	return rem[:length], rem[length:], nil
}

// DecodeString decodes the input bytes and returns string result and the remaining bytes.
func DecodeString(data []byte) (string, []byte, error) {
	length, rem, err := DecodeUnsigned32(data)
	if err != nil {
		return "", nil, err
	}
	if uint32(len(rem)) < length {
		return "", nil, errors.New("obi: out of range")
	}
	return string(rem[:length]), rem[length:], nil
}
