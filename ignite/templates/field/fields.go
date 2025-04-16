package field

import (
	"fmt"

	"github.com/ignite/cli/v28/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v28/ignite/templates/field/datatype"
)

// Fields represents a Field slice.
type Fields []Field

// GoCLIImports returns all go CLI imports.
func (f Fields) GoCLIImports() []datatype.GoImport {
	allImports := make([]datatype.GoImport, 0)
	exist := make(map[string]struct{})
	for _, fields := range f {
		for _, goImport := range fields.GoCLIImports() {
			if _, ok := exist[goImport.Name]; ok {
				continue
			}
			exist[goImport.Name] = struct{}{}
			allImports = append(allImports, goImport)
		}
	}
	return allImports
}

// ProtoImports returns all proto imports.
func (f Fields) ProtoImports() []string {
	allImports := make([]string, 0)
	exist := make(map[string]struct{})
	for _, fields := range f {
		for _, protoImport := range fields.ProtoImports() {
			if _, ok := exist[protoImport]; ok {
				continue
			}
			exist[protoImport] = struct{}{}
			allImports = append(allImports, protoImport)
		}
	}
	return allImports
}

<<<<<<< HEAD
// String returns all inline fields args for command usage.
func (f Fields) String() string {
	args := ""
	for _, field := range f {
		args += fmt.Sprintf(" [%s]", field.Name.Kebab)
=======
// ProtoFieldNameAutoCLI returns  all inline fields args for name used in proto.
// It should be used in AutoCLI to generate the field name.
func (f Fields) ProtoFieldNameAutoCLI() string {
	args := ""
	for i, field := range f {
		// only the last field can be a variadic field
		if i == len(f)-1 && field.IsSlice() {
			args += fmt.Sprintf(`{ProtoField: "%s", Varargs: true}, `, field.ProtoFieldName())
			continue
		}

		args += fmt.Sprintf(`{ProtoField: "%s"}, `, field.ProtoFieldName())
>>>>>>> 7c532030 (feat(templates): set varargs in autocli config when needed (#4624))
	}
	return args
}

// Custom return a list of custom fields.
func (f Fields) Custom() []string {
	fields := make([]string, 0)
	for _, field := range f {
		if field.DatatypeName == datatype.TypeCustom {
			dataType, err := multiformatname.NewName(field.Datatype)
			if err != nil {
				panic(err)
			}
			fields = append(fields, dataType.Snake)
		}
	}
	return fields
}
