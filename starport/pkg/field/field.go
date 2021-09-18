// Package field provides methods to parse a field provided in a command with the format Name:type
package field

import (
	"fmt"

	"github.com/tendermint/starport/starport/pkg/multiformatname"
)

// Field represents a field inside a structure for a component
// it can be a field contained in a type or inside the response of a query, etc...
type Field struct {
	Name         multiformatname.Name
	DatatypeName DataTypeName
	Datatype     string
}

// DataDeclaration return the Datatype data declaration
func (f Field) DataDeclaration() string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return datatype.DataDeclaration(f.Datatype)
}

// ProtoDeclaration return the Datatype proto declaration
func (f Field) ProtoDeclaration(index int) string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return datatype.ProtoDeclaration(f.Datatype, f.Name.LowerCamel, index)
}

// ValueDefault return the Datatype value default
func (f Field) ValueDefault() string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return datatype.ValueDefault
}

// ValueLoop return the Datatype value for loop iteration
func (f Field) ValueLoop() string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	if datatype.NonIndex {
		panic(fmt.Sprintf("non index type %s", f.DatatypeName))
	}
	return datatype.ValueLoop
}

// ValueIndex return the Datatype value for indexes
func (f Field) ValueIndex() string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	if datatype.NonIndex {
		panic(fmt.Sprintf("non index type %s", f.DatatypeName))
	}
	return datatype.ValueIndex
}

// ValueInvalidIndex return the Datatype value for invalid indexes
func (f Field) ValueInvalidIndex() string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	if datatype.NonIndex {
		panic(fmt.Sprintf("non index type %s", f.DatatypeName))
	}
	return datatype.ValueInvalidIndex
}

// GenesisArgs return the Datatype genesis args
func (f Field) GenesisArgs(value int) string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return datatype.GenesisArgs(f.Name, value)
}

// CLIArgs return the Datatype CLI args
func (f Field) CLIArgs(prefix string, argIndex int) string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return datatype.CLIArgs(f.Name, f.Datatype, prefix, argIndex)
}

// ToBytes return the Datatype byte array cast
func (f Field) ToBytes(name string) string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	if datatype.NonIndex {
		panic(fmt.Sprintf("non index type %s", f.DatatypeName))
	}
	return datatype.ToBytes(name)
}

// ToString return the Datatype byte array cast
func (f Field) ToString(name string) string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	if datatype.NonIndex {
		panic(fmt.Sprintf("non index type %s", f.DatatypeName))
	}
	return datatype.ToString(name)
}

// GoCLIImports return the Datatype imports for CLI package
func (f Field) GoCLIImports() []GoImport {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return datatype.GoCLIImports
}

// ProtoImports return the Datatype imports for proto files
func (f Field) ProtoImports() []string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return datatype.ProtoImports
}
