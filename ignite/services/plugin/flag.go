package plugin

import (
	"strconv"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

var (
	// ErrFlagNotFound error key flag not found.
	ErrFlagNotFound = errors.New("flag not found")
	// ErrInvalidFlagType error invalid flag type.
	ErrInvalidFlagType = errors.New("invalid flag type")
	// ErrFlagAssertion error flag type assertion failed.
	ErrFlagAssertion = errors.New("flag type assertion failed")
)

// Flags represents a slice of Flag pointers.
type Flags []*Flag

// getValue returns the value of the flag with the specified key and type.
// It uses the provided conversion function to convert the string value to the desired type.
func (f Flags) getValue(key string, flagType FlagType, convFunc func(v string) (interface{}, error)) (interface{}, error) {
	for _, flag := range f {
		if flag.Name == key {
			if flag.Type != flagType {
				return nil, errors.Wrapf(ErrInvalidFlagType, "invalid flag type %v for key %s", flag.Type, key)
			}
			return convFunc(flagValue(flag))
		}
	}
	return nil, errors.Wrap(ErrFlagNotFound, key)
}

// GetString retrieves the string value of the flag with the specified key.
func (f Flags) GetString(key string) (string, error) {
	v, err := f.getValue(key, FlagTypeString, func(v string) (interface{}, error) {
		return strings.TrimSpace(v), nil
	})
	if err != nil {
		return "", err
	}
	result, ok := v.(string)
	if !ok {
		return "", errors.Wrapf(ErrFlagAssertion, "invalid assertion type %T for key %s", v, key)
	}
	return result, nil
}

// GetStringSlice retrieves the string slice value of the flag with the specified key.
func (f Flags) GetStringSlice(key string) ([]string, error) {
	v, err := f.getValue(key, FlagTypeStringSlice, func(v string) (interface{}, error) {
		v = strings.Trim(v, "[]")
		s := strings.Split(v, ",")
		if len(s) == 0 || (len(s) == 1 && s[0] == "") {
			return []string{}, nil
		}
		return s, nil
	})
	if err != nil {
		return []string{}, err
	}
	result, ok := v.([]string)
	if !ok {
		return []string{}, errors.Wrapf(ErrFlagAssertion, "invalid string slice assertion type %T for key %s", v, key)
	}
	return result, nil
}

// GetBool retrieves the boolean value of the flag with the specified key.
func (f Flags) GetBool(key string) (bool, error) {
	v, err := f.getValue(key, FlagTypeBool, func(v string) (interface{}, error) {
		return strconv.ParseBool(v)
	})
	if err != nil {
		return false, err
	}
	result, ok := v.(bool)
	if !ok {
		return false, errors.Wrapf(ErrFlagAssertion, "invalid bool assertion type %T for key %s", v, key)
	}
	return result, nil
}

// GetInt retrieves the integer value of the flag with the specified key.
func (f Flags) GetInt(key string) (int, error) {
	v, err := f.getValue(key, FlagTypeInt, func(v string) (interface{}, error) {
		return strconv.Atoi(v)
	})
	if err != nil {
		return 0, err
	}
	result, ok := v.(int)
	if !ok {
		return 0, errors.Wrapf(ErrFlagAssertion, "invalid int assertion type %T for key %s", v, key)
	}
	return result, nil
}

// GetInt64 retrieves the int64 value of the flag with the specified key.
func (f Flags) GetInt64(key string) (int64, error) {
	v, err := f.getValue(key, FlagTypeInt64, func(v string) (interface{}, error) {
		return strconv.ParseInt(v, 10, 64)
	})
	if err != nil {
		return int64(0), err
	}
	result, ok := v.(int64)
	if !ok {
		return int64(0), errors.Wrapf(ErrFlagAssertion, "invalid int64 assertion type %T for key %s", v, key)
	}
	return result, nil
}

// GetUint retrieves the uint value of the flag with the specified key.
func (f Flags) GetUint(key string) (uint, error) {
	v, err := f.getValue(key, FlagTypeUint, func(v string) (interface{}, error) {
		return strconv.ParseUint(v, 10, 64)
	})
	if err != nil {
		return uint(0), err
	}
	result, ok := v.(uint64)
	if !ok {
		return uint(0), errors.Wrapf(ErrFlagAssertion, "invalid uint assertion type %T for key %s", v, key)
	}
	return uint(result), nil
}

// GetUint64 retrieves the uint64 value of the flag with the specified key.
func (f Flags) GetUint64(key string) (uint64, error) {
	v, err := f.getValue(key, FlagTypeUint64, func(v string) (interface{}, error) {
		return strconv.ParseUint(v, 10, 64)
	})
	if err != nil {
		return uint64(0), err
	}
	result, ok := v.(uint64)
	if !ok {
		return uint64(0), errors.Wrapf(ErrFlagAssertion, "invalid uint64 assertion type %T for key %s", v, key)
	}
	return result, nil
}

// flagValue returns the value of the flag if set, otherwise returns the default value.
func flagValue(flag *Flag) string {
	if flag.Value != "" {
		return flag.Value
	}
	return flag.DefaultValue
}
