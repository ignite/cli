// Package field provides methods to parse a field provided in a command with the format name:type
package field

import (
	"fmt"

	"github.com/emicklei/proto"

	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/templates/field/datatype"
)

// Field represents a field inside a structure for a component
// it can be a field contained in a type or inside the response of a query, etc...
type Field struct {
	Name         multiformatname.Name
	DatatypeName datatype.Name
	Datatype     string
}

// DataType returns the field Datatype.
func (f Field) DataType() string {
	dt, ok := datatype.IsSupportedType(f.DatatypeName)
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return dt.DataType(f.Datatype)
}

// ProtoFieldName returns the field name used in proto.
func (f Field) ProtoFieldName() string {
	return f.Name.LowerCamel
}

// ProtoType returns the field proto Datatype.
func (f Field) ProtoType(index int) string {
	dt, ok := datatype.IsSupportedType(f.DatatypeName)
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return dt.ProtoType(f.Datatype, f.ProtoFieldName(), index)
}

// DefaultTestValue returns the Datatype value default.
func (f Field) DefaultTestValue() string {
	dt, ok := datatype.IsSupportedType(f.DatatypeName)
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return dt.DefaultTestValue
}

// ValueLoop returns the Datatype value for loop iteration.
func (f Field) ValueLoop() string {
	dt, ok := datatype.IsSupportedType(f.DatatypeName)
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	if dt.NonIndex {
		panic(fmt.Sprintf("non index type %s", f.DatatypeName))
	}
	return dt.ValueLoop
}

// ValueIndex returns the Datatype value for indexes.
func (f Field) ValueIndex() string {
	dt, ok := datatype.IsSupportedType(f.DatatypeName)
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	if dt.NonIndex {
		panic(fmt.Sprintf("non index type %s", f.DatatypeName))
	}
	return dt.ValueIndex
}

// ValueInvalidIndex returns the Datatype value for invalid indexes.
func (f Field) ValueInvalidIndex() string {
	dt, ok := datatype.IsSupportedType(f.DatatypeName)
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	if dt.NonIndex {
		panic(fmt.Sprintf("non index type %s", f.DatatypeName))
	}
	return dt.ValueInvalidIndex
}

// GenesisArgs returns the Datatype genesis args.
func (f Field) GenesisArgs(value int) string {
	dt, ok := datatype.IsSupportedType(f.DatatypeName)
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return dt.GenesisArgs(f.Name, value)
}

// CLIArgs returns the Datatype CLI args.
func (f Field) CLIArgs(prefix string, argIndex int) string {
	dt, ok := datatype.IsSupportedType(f.DatatypeName)
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return dt.CLIArgs(f.Name, f.Datatype, prefix, argIndex)
}

// ToBytes returns the Datatype byte array cast.
func (f Field) ToBytes(name string) string {
	dt, ok := datatype.IsSupportedType(f.DatatypeName)
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	if dt.NonIndex {
		panic(fmt.Sprintf("non index type %s", f.DatatypeName))
	}
	return dt.ToBytes(name)
}

// ToString returns the Datatype byte array cast.
func (f Field) ToString(name string) string {
	dt, ok := datatype.IsSupportedType(f.DatatypeName)
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	if dt.NonIndex {
		panic(fmt.Sprintf("non index type %s", f.DatatypeName))
	}
	return dt.ToString(name)
}

// ToProtoField returns the Datatype as a *proto.Field node.
func (f Field) ToProtoField(index int) *proto.NormalField {
	// TODO: Do we can if it's an index type?
	dt, ok := datatype.IsSupportedType(f.DatatypeName)
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return dt.ToProtoField(f.Datatype, f.Name.LowerCamel, index)
}

// GoCLIImports returns the Datatype imports for CLI package.
func (f Field) GoCLIImports() []datatype.GoImport {
	dt, ok := datatype.IsSupportedType(f.DatatypeName)
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return dt.GoCLIImports
}

// ProtoImports returns the Datatype imports for proto files.
func (f Field) ProtoImports() []string {
	dt, ok := datatype.IsSupportedType(f.DatatypeName)
	if !ok {
		panic(fmt.Sprintf("unknown type %s", f.DatatypeName))
	}
	return dt.ProtoImports
}
