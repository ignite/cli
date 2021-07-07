package obi

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func getSchemaImpl(s *strings.Builder, t reflect.Type) error {
	switch t.Kind() {
	case reflect.Uint8:
		s.WriteString("u8")
		return nil
	case reflect.Uint16:
		s.WriteString("u16")
		return nil
	case reflect.Uint32:
		s.WriteString("u32")
		return nil
	case reflect.Uint64:
		s.WriteString("u64")
		return nil
	case reflect.Int8:
		s.WriteString("i8")
		return nil
	case reflect.Int16:
		s.WriteString("i16")
		return nil
	case reflect.Int32:
		s.WriteString("i32")
		return nil
	case reflect.Int64:
		s.WriteString("i64")
		return nil
	case reflect.String:
		s.WriteString("string")
		return nil
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			s.WriteString("bytes")
			return nil
		}
		s.WriteString("[")
		err := getSchemaImpl(s, t.Elem())
		if err != nil {
			return err
		}
		s.WriteString("]")
		return nil
	case reflect.Struct:
		if t.NumField() == 0 {
			return errors.New("obi: empty struct is not supported")
		}
		s.WriteString("{")
		for idx := 0; idx < t.NumField(); idx++ {
			field := t.Field(idx)
			name, ok := field.Tag.Lookup("obi")
			if !ok {
				return fmt.Errorf("obi: no obi tag found for field %s of %s", field.Name, t.Name())
			}
			if idx != 0 {
				s.WriteString(",")
			}
			s.WriteString(name)
			s.WriteString(":")
			err := getSchemaImpl(s, field.Type)
			if err != nil {
				return err
			}
		}
		s.WriteString("}")
		return nil
	default:
		return fmt.Errorf("obi: unsupported value type: %s", t.Kind())
	}
}

// GetSchema returns the compact OBI individual schema of the given value.
func GetSchema(v interface{}) (string, error) {
	s := &strings.Builder{}
	err := getSchemaImpl(s, reflect.TypeOf(v))
	if err != nil {
		return "", err
	}
	return s.String(), nil
}

// MustGetSchema returns the compact OBI individual schema of the given value. Panics on error.
func MustGetSchema(v interface{}) string {
	schema, err := GetSchema(v)
	if err != nil {
		panic(err)
	}
	return schema
}
