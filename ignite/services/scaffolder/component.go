package scaffolder

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/pkg/protoanalysis"
	"github.com/ignite/cli/ignite/templates/field/datatype"
)

const (
	componentType    = "type"
	componentMessage = "message"
	componentQuery   = "query"
	componentPacket  = "packet"

	protoFolder = "proto"
)

// checkComponentValidity performs various checks common to all components to verify if it can be scaffolded.
func checkComponentValidity(appPath, moduleName string, compName multiformatname.Name, noMessage bool) error {
	ok, err := moduleExists(appPath, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("the module %s doesn't exist", moduleName)
	}

	// Ensure the name is valid, otherwise it would generate an incorrect code
	if err := checkForbiddenComponentName(compName); err != nil {
		return fmt.Errorf("%s can't be used as a component name: %w", compName.LowerCamel, err)
	}

	// Check component name is not already used
	return checkComponentCreated(appPath, moduleName, compName, noMessage)
}

// checkComponentCreated checks if the component has been already created with Ignite in the project.
func checkComponentCreated(appPath, moduleName string, compName multiformatname.Name, noMessage bool) (err error) {
	// associate the type to check with the component that scaffold this type
	typesToCheck := map[string]string{
		compName.UpperCamel:                          componentType,
		"queryall" + compName.LowerCase + "request":  componentType,
		"queryall" + compName.LowerCase + "response": componentType,
		"queryget" + compName.LowerCase + "request":  componentType,
		"queryget" + compName.LowerCase + "response": componentType,
		"query" + compName.LowerCase + "request":     componentQuery,
		"query" + compName.LowerCase + "response":    componentQuery,
		compName.LowerCase + "packetdata":            componentPacket,
	}

	if !noMessage {
		typesToCheck["msgcreate"+compName.LowerCase] = componentType
		typesToCheck["msgupdate"+compName.LowerCase] = componentType
		typesToCheck["msgdelete"+compName.LowerCase] = componentType
		typesToCheck["msg"+compName.LowerCase] = componentMessage
		typesToCheck["msgsend"+compName.LowerCase] = componentPacket
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
					err = fmt.Errorf("component %s with name %s is already created (type %s exists)",
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
func checkCustomTypes(ctx context.Context, path, appName, module string, fields []string) error {
	protoPath := filepath.Join(path, protoFolder, appName, module)
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
	return protoanalysis.HasMessages(ctx, protoPath, customFieldTypes...)
}

// checkForbiddenComponentName returns true if the name is forbidden as a component name.
func checkForbiddenComponentName(name multiformatname.Name) error {
	// Check with names already used from the scaffolded code
	switch name.LowerCase {
	case
		"oracle",
		"logger",
		"keeper",
		"query",
		"genesis",
		"types",
		"tx",
		datatype.TypeCustom:
		return fmt.Errorf("%s is used by Ignite scaffolder", name.LowerCamel)
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
		return fmt.Errorf("%s is a Go keyword", name)
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
		return fmt.Errorf("%s is a Go built-in identifier", name)
	}
	return nil
}

// checkForbiddenOracleFieldName returns true if the name is forbidden as an oracle field name.
//
// Deprecated: This function is no longer maintained.
func checkForbiddenOracleFieldName(name string) error {
	mfName, err := multiformatname.NewName(name, multiformatname.NoNumber)
	if err != nil {
		return err
	}

	// Check with names already used from the scaffolded code
	switch mfName.UpperCase {
	case
		"CLIENTID",
		"ORACLESCRIPTID",
		"SOURCECHANNEL",
		"CALLDATA",
		"ASKCOUNT",
		"MINCOUNT",
		"FEELIMIT",
		"PREPAREGAS",
		"EXECUTEGAS":
		return fmt.Errorf("%s is used by Ignite scaffolder", name)
	}
	return nil
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
