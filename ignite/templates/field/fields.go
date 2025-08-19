package field

import (
	"fmt"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/templates/field/datatype"
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
	}
	args = strings.TrimSpace(args)
	return strings.Trim(args, ",")
}

// CLIUsage returns all inline fields args for CLI command usage.
func (f Fields) CLIUsage() string {
	args := ""
	for _, field := range f {
		args += fmt.Sprintf(" [%s]", field.CLIUsage())
	}
	return strings.TrimSpace(args)
}

// Custom returns a list of custom fields.
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
