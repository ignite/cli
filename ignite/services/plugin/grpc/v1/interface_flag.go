package v1

import (
	"strconv"
	"strings"

	"github.com/spf13/pflag"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
)

const (
	cobraFlagTypeBool        = "bool"
	cobraFlagTypeInt         = "int"
	cobraFlagTypeInt64       = "int64"
	cobraFlagTypeString      = "string"
	cobraFlagTypeStringSlice = "stringSlice"
	cobraFlagTypeUint        = "uint"
	cobraFlagTypeUint64      = "uint64"
)

var flagTypes = map[string]Flag_Type{
	cobraFlagTypeBool:        Flag_TYPE_FLAG_BOOL,
	cobraFlagTypeInt:         Flag_TYPE_FLAG_INT,
	cobraFlagTypeInt64:       Flag_TYPE_FLAG_INT64,
	cobraFlagTypeString:      Flag_TYPE_FLAG_STRING_UNSPECIFIED,
	cobraFlagTypeStringSlice: Flag_TYPE_FLAG_STRING_SLICE,
	cobraFlagTypeUint:        Flag_TYPE_FLAG_UINT,
	cobraFlagTypeUint64:      Flag_TYPE_FLAG_UINT64,
}

func newDefaultFlagValueError(typeName, value string) error {
	return errors.Errorf("invalid default value for plugin command %s flag: %s", typeName, value)
}

func (f *Flag) exportToFlagSet(fs *pflag.FlagSet) error {
	switch f.Type {
	case Flag_TYPE_FLAG_BOOL:
		v, err := strconv.ParseBool(f.DefaultValue)
		if err != nil {
			return newDefaultFlagValueError(cobraFlagTypeBool, f.DefaultValue)
		}

		fs.BoolP(f.Name, f.Shorthand, v, f.Usage)
		fs.Set(f.Name, f.Value)
	case Flag_TYPE_FLAG_INT:
		v, err := strconv.Atoi(f.DefaultValue)
		if err != nil {
			return newDefaultFlagValueError(cobraFlagTypeInt, f.DefaultValue)
		}

		fs.IntP(f.Name, f.Shorthand, v, f.Usage)
		fs.Set(f.Name, f.Value)
	case Flag_TYPE_FLAG_UINT:
		v, err := strconv.ParseUint(f.DefaultValue, 10, 64)
		if err != nil {
			return newDefaultFlagValueError(cobraFlagTypeUint, f.DefaultValue)
		}

		fs.UintP(f.Name, f.Shorthand, uint(v), f.Usage)
		fs.Set(f.Name, f.Value)
	case Flag_TYPE_FLAG_INT64:
		v, err := strconv.ParseInt(f.DefaultValue, 10, 64)
		if err != nil {
			return newDefaultFlagValueError(cobraFlagTypeInt64, f.DefaultValue)
		}

		fs.Int64P(f.Name, f.Shorthand, v, f.Usage)
		fs.Set(f.Name, f.Value)
	case Flag_TYPE_FLAG_UINT64:
		v, err := strconv.ParseUint(f.DefaultValue, 10, 64)
		if err != nil {
			return newDefaultFlagValueError(cobraFlagTypeInt64, f.DefaultValue)
		}

		fs.Uint64P(f.Name, f.Shorthand, v, f.Usage)
		fs.Set(f.Name, f.Value)
	case Flag_TYPE_FLAG_STRING_SLICE:
		s := strings.Trim(f.DefaultValue, "[]")
		fs.StringSliceP(f.Name, f.Shorthand, strings.Fields(s), f.Usage)
		fs.Set(f.Name, strings.Trim(f.Value, "[]"))
	case Flag_TYPE_FLAG_STRING_UNSPECIFIED:
		fs.StringP(f.Name, f.Shorthand, f.DefaultValue, f.Usage)
		fs.Set(f.Name, f.Value)
	}
	return nil
}

type flagger interface {
	Flags() *pflag.FlagSet
	PersistentFlags() *pflag.FlagSet
}

func extractCobraFlags(cmd flagger) []*Flag {
	var flags []*Flag

	if cmd.Flags() != nil {
		cmd.Flags().VisitAll(func(pf *pflag.Flag) {
			// Skip persistent flags
			if cmd.PersistentFlags().Lookup(pf.Name) != nil {
				return
			}

			flags = append(flags, &Flag{
				Name:         pf.Name,
				Shorthand:    pf.Shorthand,
				Usage:        pf.Usage,
				DefaultValue: pf.DefValue,
				Value:        pf.Value.String(),
				Type:         flagTypes[pf.Value.Type()],
			})
		})
	}

	if cmd.PersistentFlags() != nil {
		cmd.PersistentFlags().VisitAll(func(pf *pflag.Flag) {
			flags = append(flags, &Flag{
				Name:         pf.Name,
				Shorthand:    pf.Shorthand,
				Usage:        pf.Usage,
				DefaultValue: pf.DefValue,
				Value:        pf.Value.String(),
				Type:         flagTypes[pf.Value.Type()],
				Persistent:   true,
			})
		})
	}

	return flags
}
