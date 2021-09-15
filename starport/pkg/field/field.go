// Package field provides methods to parse a field provided in a command with the format name:type
package field

import (
	"fmt"
	"math/rand"

	"github.com/tendermint/starport/starport/pkg/multiformatname"
)

const (
	TypeCustom      = "custom"
	TypeString      = "string"
	TypeStringSlice = "[]string"
	TypeBool        = "bool"
	TypeInt         = "int"
	TypeIntSlice    = "[]int"
	TypeUint        = "uint"
	TypeUintSlice   = "[]uint"

	TypeNameInt  = "int32"
	TypeNameUint = "uint64"

	TypeSeparator = ":"
)

var (
	StaticDataTypes = map[string]string{
		TypeString:      TypeString,
		TypeStringSlice: TypeString,
		TypeBool:        TypeBool,
		TypeInt:         TypeNameInt,
		TypeIntSlice:    TypeNameInt,
		TypeUint:        TypeNameUint,
		TypeUintSlice:   TypeNameUint,
	}
)

type (
	// Field represents a field inside a structure for a component
	// it can be a field contained in a type or inside the response of a query, etc...
	Field struct {
		Name         multiformatname.Name
		Datatype     string
		DatatypeName string
	}

	// Fields represents a Field slice
	Fields []Field
)

// GetDatatype return the Datatype based in the DatatypeName
func (f Field) GetDatatype() string {
	switch f.DatatypeName {
	case TypeString, TypeBool, TypeInt, TypeUint:
		return f.Datatype
	case TypeStringSlice, TypeIntSlice, TypeUintSlice:
		return fmt.Sprintf("[]%s", f.Datatype)
	case TypeCustom:
		return fmt.Sprintf("*%s", f.Datatype)
	default:
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
}

// GetProtoDatatype return the proto Datatype based in the DatatypeName
func (f Field) GetProtoDatatype() string {
	switch f.DatatypeName {
	case TypeString, TypeBool, TypeInt, TypeUint, TypeCustom:
		return f.Datatype
	case TypeStringSlice, TypeIntSlice, TypeUintSlice:
		return fmt.Sprintf("repeated %s", f.Datatype)
	default:
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
}

// NeedCastImport return true if the field slice
// needs import the cast library
func (f Fields) NeedCastImport() bool {
	for _, field := range f {
		if field.DatatypeName != TypeString &&
			field.DatatypeName != TypeCustom {
			return true
		}
	}
	return false
}

// IsComplex return true if the field slice
// needs import the json library
func (f Fields) IsComplex() bool {
	for _, field := range f {
		if field.DatatypeName == TypeCustom {
			return true
		}
	}
	return false
}

// String return all inline fields args for command usage
func (f Fields) String() string {
	args := ""
	for _, field := range f {
		args += fmt.Sprintf(" [%s]", field.Name.Kebab)
	}
	return args
}

// Custom return a list of custom fields
func (f Fields) Custom() []string {
	fields := make([]string, 0)
	for _, field := range f {
		if field.DatatypeName == TypeCustom {
			dataType, err := multiformatname.NewName(field.Datatype)
			if err != nil {
				panic(err)
			}
			fields = append(fields, dataType.Snake)
		}
	}
	return fields
}

// GenesisField create a genesis field
func (f Field) GenesisField(value int) string {
	switch f.DatatypeName {
	case TypeString:
		return fmt.Sprintf("%s: \"%s\",\n", f.Name.UpperCamel, f.Name.LowerCamel)
	case TypeStringSlice:
		return fmt.Sprintf("%s: []string{\"%s\"},\n", f.Name.UpperCamel, f.Name.LowerCamel)
	case TypeInt, TypeUint:
		return fmt.Sprintf("%s: %d,\n", f.Name.UpperCamel, rand.Intn(value))
	case TypeIntSlice:
		return fmt.Sprintf("%s: []int32{%d},\n", f.Name.UpperCamel, rand.Intn(value))
	case TypeUintSlice:
		return fmt.Sprintf("%s: []uint64{%d},\n", f.Name.UpperCamel, rand.Intn(value))
	case TypeBool:
		return fmt.Sprintf("%s: %t,\n", f.Name.UpperCamel, rand.Intn(value)%2 == 0)
	default:
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
}
