// Package field provides methods to parse a field provided in a command with the format name:type
package field

import (
	"fmt"

	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"github.com/tendermint/starport/starport/templates/field/datatype"
)

// Field represents a field inside a structure for a component
// it can be a field contained in a type or inside the response of a query, etc...
type Field struct {
	Name         multiformatname.Name
	DatatypeName datatype.Name
	Datatype     string
}

// DataType returns the field Datatype
func (f Field) DataType() string {
	dt, ok := datatype.SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return dt.DataType(f.Datatype)
}

// ProtoType returns the field proto Datatype
func (f Field) ProtoType(index int) string {
	dt, ok := datatype.SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return dt.ProtoType(f.Datatype, f.Name.LowerCamel, index)
}

// ValueDefault returns the Datatype value default
func (f Field) ValueDefault() string {
	dt, ok := datatype.SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return dt.ValueDefault
}

// ValueLoop returns the Datatype value for loop iteration
func (f Field) ValueLoop() string {
	dt, ok := datatype.SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	if dt.NonIndex {
		panic(fmt.Sprintf("non index type %s", f.DatatypeName))
	}
	return dt.ValueLoop
}

// ValueIndex returns the Datatype value for indexes
func (f Field) ValueIndex() string {
	dt, ok := datatype.SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	if dt.NonIndex {
		panic(fmt.Sprintf("non index type %s", f.DatatypeName))
	}
	return dt.ValueIndex
}

// ValueInvalidIndex returns the Datatype value for invalid indexes
func (f Field) ValueInvalidIndex() string {
	dt, ok := datatype.SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	if dt.NonIndex {
		panic(fmt.Sprintf("non index type %s", f.DatatypeName))
	}
	return dt.ValueInvalidIndex
}

// GenesisArgs returns the Datatype genesis args
func (f Field) GenesisArgs(value int) string {
	dt, ok := datatype.SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return dt.GenesisArgs(f.Name, value)
}

// CLIArgs returns the Datatype CLI args
func (f Field) CLIArgs(prefix string, argIndex int) string {
	dt, ok := datatype.SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return dt.CLIArgs(f.Name, f.Datatype, prefix, argIndex)
}

// ToBytes returns the Datatype byte array cast
func (f Field) ToBytes(name string) string {
	dt, ok := datatype.SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	if dt.NonIndex {
		panic(fmt.Sprintf("non index type %s", f.DatatypeName))
	}
	return dt.ToBytes(name)
}

// ToString returns the Datatype byte array cast
func (f Field) ToString(name string) string {
	dt, ok := datatype.SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	if dt.NonIndex {
		panic(fmt.Sprintf("non index type %s", f.DatatypeName))
	}
	return dt.ToString(name)
}

// GoCLIImports returns the Datatype imports for CLI package
func (f Field) GoCLIImports() []datatype.GoImport {
	dt, ok := datatype.SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return dt.GoCLIImports
}

// ProtoImports return the Datatype imports for proto files
func (f Field) ProtoImports() []string {
	dt, ok := datatype.SupportedTypes[f.DatatypeName]
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return dt.ProtoImports
}
