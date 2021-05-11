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

	"github.com/gobuffalo/genny"
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
	appPath string,
	opts *modulecreate.MsgServerOptions,
) error {
	// Check if convention used
	msgServerDefined, err := isMsgServerDefined(appPath, opts.ModuleName)
	if err != nil {
		return err
	}
	if !msgServerDefined {
		// Patch the module to support the convention
		g, err := modulecreate.AddMsgServerConventionToLegacyModule(opts)
		if err != nil {
			return err
		}
		run := genny.WetRunner(context.Background())
		run.With(g)
		return run.Run()
	}
	return nil
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
func checkComponentValidity(appPath, moduleName, compName string) error {
	ok, err := moduleExists(appPath, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("the module %s doesn't exist", moduleName)
	}

	// Ensure the name is valid, otherwise it would generate an incorrect code
	if err := checkForbiddenComponentName(compName); err != nil {
		return fmt.Errorf("%s can't be used as a component name: %s", compName, err.Error())
	}

	// Check component name is not already used
	return checkComponentCreated(appPath, moduleName, compName)
}

// checkForbiddenComponentName returns true if the name is forbidden as a component name
func checkForbiddenComponentName(name string) error {
	// Check with names already used from the scaffolded code
	switch name {
	case
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

// ErrComponentAlreadyCreated is the error returned when a specific component is already created
type ErrComponentAlreadyCreated struct {
	name          string
	checkedType   string
	componentType string
}

// NewErrComponentAlreadyCreated returns a new ErrComponentAlreadyCreated error
func NewErrComponentAlreadyCreated(name, checkedType, componentType string) *ErrComponentAlreadyCreated {
	return &ErrComponentAlreadyCreated{
		name,
		checkedType,
		componentType,
	}
}

func (e *ErrComponentAlreadyCreated) Error() string {
	return fmt.Sprintf("component %s with name %s is already created (type %s exists)",
		e.componentType, e.name, e.checkedType)
}

// checkComponentCreated checks if the component has been already created with Starport in the project
func checkComponentCreated(appPath, moduleName, compName string) (err error) {
	compNameTitle := strings.Title(compName)

	// associate the type to check with the component that scaffold this type
	typesToCheck := map[string]string{
		compNameTitle:                           componentType,
		"MsgCreate" + compNameTitle:             componentType,
		"MsgUpdate" + compNameTitle:             componentType,
		"MsgDelete" + compNameTitle:             componentType,
		"QueryAll" + compNameTitle + "Request":  componentType,
		"QueryAll" + compNameTitle + "Response": componentType,
		"QueryGet" + compNameTitle + "Request":  componentType,
		"QueryGet" + compNameTitle + "Response": componentType,
		"Msg" + compNameTitle:                   componentMessage,
		"Query" + compNameTitle + "Request":     componentQuery,
		"Query" + compNameTitle + "Response":    componentQuery,
		"MsgSend" + compNameTitle:               componentPacket,
		compNameTitle + "PacketData":            componentPacket,
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
					err = NewErrComponentAlreadyCreated(compName, typeSpec.Name.Name, compType)
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
