package scaffolder

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis"
	"github.com/ignite/cli/v29/ignite/templates/field/datatype"
)

const (
	componentType    = "type"
	componentMessage = "message"
	componentQuery   = "query"
	componentPacket  = "packet"
)

// checkComponentValidity performs various checks common to all components to verify if it can be scaffolded.
func checkComponentValidity(appPath, moduleName string, compName multiformatname.Name, noMessage bool) error {
	ok, err := moduleExists(appPath, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf("the module %s doesn't exist", moduleName)
	}

	// Ensure the name is valid, otherwise it would generate an incorrect code
	if err := checkForbiddenComponentName(compName); err != nil {
		return errors.Errorf("%s can't be used as a component name: %w", compName.LowerCamel, err)
	}

	// Check component name is not already used
	return checkComponentCreated(appPath, moduleName, compName, noMessage)
}

// checkComponentCreated checks if the component has been already created with Ignite in the project.
func checkComponentCreated(appPath, moduleName string, compName multiformatname.Name, noMessage bool) (err error) {
	// associate the type to check with the component that scaffold this type
	typesToCheck := map[string]string{
		compName.UpperCamel: componentType,
		fmt.Sprintf("queryall%srequest", compName.LowerCase):  componentType,
		fmt.Sprintf("queryall%sresponse", compName.LowerCase): componentType,
		fmt.Sprintf("queryget%srequest", compName.LowerCase):  componentType,
		fmt.Sprintf("queryget%sresponse", compName.LowerCase): componentType,
		fmt.Sprintf("query%srequest", compName.LowerCase):     componentQuery,
		fmt.Sprintf("query%sresponse", compName.LowerCase):    componentQuery,
		fmt.Sprintf("%spacketdata", compName.LowerCase):       componentPacket,
	}

	if !noMessage {
		typesToCheck[fmt.Sprintf("msgcreate%s", compName.LowerCase)] = componentType
		typesToCheck[fmt.Sprintf("msgupdate%s", compName.LowerCase)] = componentType
		typesToCheck[fmt.Sprintf("msgdelete%s", compName.LowerCase)] = componentType
		typesToCheck[fmt.Sprintf("msg%s", compName.LowerCase)] = componentMessage
		typesToCheck[fmt.Sprintf("msgsend%s", compName.LowerCase)] = componentPacket
	}

	absPath, err := filepath.Abs(filepath.Join(appPath, "x", moduleName, "types"))
	if err != nil {
		return err
	}
	fileSet := token.NewFileSet()
	all, err := parser.ParseDir(fileSet, absPath, func(os.FileInfo) bool { return true }, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range all {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(x ast.Node) bool {
				typeSpec, ok := x.(*ast.TypeSpec)
				if !ok {
					return true
				}

				if _, ok := typeSpec.Type.(*ast.StructType); !ok {
					return true
				}

				// Check if the parsed type is from a scaffolded component with the name
				if compType, ok := typesToCheck[strings.ToLower(typeSpec.Name.Name)]; ok {
					err = errors.Errorf("component %s with name %s is already created (type %s exists)",
						compType,
						compName.Original,
						typeSpec.Name.Name,
					)
					return false
				}

				return true
			})
			if err != nil {
				return err
			}
		}
	}
	return err
}

// checkCustomTypes returns error if one of the types is invalid.
func checkCustomTypes(ctx context.Context, appPath, appName, protoDir, module string, fields []string) error {
	path := filepath.Join(appPath, protoDir, appName, module)
	customFieldTypes := make([]string, 0)
	for _, field := range fields {
		ft, ok := fieldType(field)
		if !ok {
			continue
		}

		if _, ok := datatype.IsSupportedType(datatype.Name(ft)); !ok {
			customFieldTypes = append(customFieldTypes, ft)
		}
	}
	return protoanalysis.HasMessages(ctx, path, customFieldTypes...)
}

// checkForbiddenComponentName returns true if the name is forbidden as a component name.
func checkForbiddenComponentName(name multiformatname.Name) error {
	// Check with names already used from the scaffolded code
	switch name.LowerCase {
	case
		"logger",
		"keeper",
		"query",
		"genesis",
		"types",
		"tx",
		datatype.TypeCustom:
		return errors.Errorf("%s is used by Ignite scaffolder", name.LowerCamel)
	}

	if strings.HasSuffix(name.LowerCase, "test") {
		return errors.New(`name cannot end with "test"`)
	}

	return checkGoReservedWord(name.LowerCamel)
}

// checkGoReservedWord checks if the name can't be used because it is a go reserved keyword.
func checkGoReservedWord(name string) error {
	// Check keyword or literal
	if token.Lookup(name).IsKeyword() {
		return errors.Errorf("%s is a Go keyword", name)
	}

	// Check with builtin identifier
	switch name {
	case
		"panic",
		"recover",
		"append",
		"bool",
		"byte",
		"cap",
		"close",
		"complex",
		"complex64",
		"complex128",
		"uint16",
		"copy",
		"false",
		"float32",
		"float64",
		"imag",
		"int",
		"int8",
		"int16",
		"uint32",
		"int32",
		"int64",
		"iota",
		"len",
		"make",
		"new",
		"nil",
		"uint64",
		"print",
		"println",
		"real",
		"string",
		"true",
		"uint",
		"uint8",
		"uintptr":
		return errors.Errorf("%s is a Go built-in identifier", name)
	}
	return checkMaxLength(name)
}

// containsCustomTypes returns true if the list of fields contains at least one custom type.
func containsCustomTypes(fields []string) bool {
	for _, field := range fields {
		ft, ok := fieldType(field)
		if !ok {
			continue
		}

		if _, ok := datatype.IsSupportedType(datatype.Name(ft)); !ok {
			return true
		}
	}
	return false
}

// checks if a field is given.  Returns type if true.
func fieldType(field string) (fieldType string, isCustom bool) {
	fieldSplit := strings.Split(field, datatype.Separator)
	if len(fieldSplit) <= 1 {
		return "", false
	}

	return fieldSplit[1], true
}
