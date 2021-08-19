package scaffolder

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/multiformatname"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
)

const (
	componentType    = "type"
	componentMessage = "message"
	componentQuery   = "query"
	componentPacket  = "packet"
)

// supportMsgServer checks if the module supports the MsgServer convention
// if not, the module codebase is modified to support it
// https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-031-msg-service.md
func supportMsgServer(
	replacer placeholder.Replacer,
	appPath string,
	opts *modulecreate.MsgServerOptions,
) (*genny.Generator, error) {
	// Check if convention used
	msgServerDefined, err := isMsgServerDefined(appPath, opts.ModuleName)
	if err != nil {
		return nil, err
	}
	if !msgServerDefined {
		// Patch the module to support the convention
		return modulecreate.AddMsgServerConventionToLegacyModule(replacer, opts)
	}
	return nil, nil
}

// isMsgServerDefined checks if the module uses the MsgServer convention for transactions
// this is checked by verifying the existence of the tx.proto file
func isMsgServerDefined(appPath, moduleName string) (bool, error) {
	txProto, err := filepath.Abs(filepath.Join(appPath, "proto", moduleName, "tx.proto"))
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(txProto); os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// checkComponentValidity performs various checks common to all components to verify if it can be scaffolded
func checkComponentValidity(appPath, moduleName string, compName multiformatname.Name, noMessage bool) error {
	ok, err := moduleExists(appPath, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("the module %s doesn't exist", moduleName)
	}

	// Ensure the name is valid, otherwise it would generate an incorrect code
	if err := checkForbiddenComponentName(compName.LowerCamel); err != nil {
		return fmt.Errorf("%s can't be used as a component name: %s", compName, err.Error())
	}

	// Check component name is not already used
	return checkComponentCreated(appPath, moduleName, compName, noMessage)
}

// checkForbiddenComponentName returns true if the name is forbidden as a component name
func checkForbiddenComponentName(name string) error {
	// Check with names already used from the scaffolded code
	switch name {
	case
		"oracle",
		"logger",
		"keeper",
		"query",
		"genesis",
		"types",
		"tx":
		return fmt.Errorf("%s is used by Starport scaffolder", name)
	}

	return checkGoReservedWord(name)
}

// checkGoReservedWord checks if the name can't be used because it is a go reserved keyword
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

// checkComponentCreated checks if the component has been already created with Starport in the project
func checkComponentCreated(appPath, moduleName string, compName multiformatname.Name, noMessage bool) (err error) {

	// associate the type to check with the component that scaffold this type
	typesToCheck := map[string]string{
		compName.UpperCamel:                           componentType,
		"QueryAll" + compName.UpperCamel + "Request":  componentType,
		"QueryAll" + compName.UpperCamel + "Response": componentType,
		"QueryGet" + compName.UpperCamel + "Request":  componentType,
		"QueryGet" + compName.UpperCamel + "Response": componentType,
		"Query" + compName.UpperCamel + "Request":     componentQuery,
		"Query" + compName.UpperCamel + "Response":    componentQuery,
		compName.UpperCamel + "PacketData":            componentPacket,
	}

	if !noMessage {
		typesToCheck["MsgCreate"+compName.UpperCamel] = componentType
		typesToCheck["MsgUpdate"+compName.UpperCamel] = componentType
		typesToCheck["MsgDelete"+compName.UpperCamel] = componentType
		typesToCheck["Msg"+compName.UpperCamel] = componentMessage
		typesToCheck["MsgSend"+compName.UpperCamel] = componentPacket
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
				if compType, ok := typesToCheck[typeSpec.Name.Name]; ok {
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
				return
			}
		}
	}
	return err
}
