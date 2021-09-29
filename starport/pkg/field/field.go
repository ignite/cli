// Package field provides methods to parse a field provided in a command with the format name:type
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

// DataType returns the field Datatype
func (f Field) DataType() string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return datatype.DataType(f.Datatype)
}

// ProtoType returns the field proto Datatype
func (f Field) ProtoType(index int) string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return datatype.ProtoType(f.Datatype, f.Name.LowerCamel, index)
}

// ValueDefault returns the Datatype value default
func (f Field) ValueDefault() string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return datatype.ValueDefault
}

// ValueLoop returns the Datatype value for loop iteration
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

// ValueIndex returns the Datatype value for indexes
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

// ValueInvalidIndex returns the Datatype value for invalid indexes
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

// GenesisArgs returns the Datatype genesis args
func (f Field) GenesisArgs(value int) string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return datatype.GenesisArgs(f.Name, value)
}

// CLIArgs returns the Datatype CLI args
func (f Field) CLIArgs(prefix string, argIndex int) string {
	datatype, ok := SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return datatype.CLIArgs(f.Name, f.Datatype, prefix, argIndex)
}

// ToBytes returns the Datatype byte array cast
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

// ToString returns the Datatype byte array cast
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

// GoCLIImports returns the Datatype imports for CLI package
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
