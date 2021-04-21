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

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	modulecreate "github.com/tendermint/starport/starport/templates/module/create"
	"github.com/tendermint/starport/starport/templates/typed"
	"github.com/tendermint/starport/starport/templates/typed/indexed"
)

const (
	TypeString = "string"
	TypeBool   = "bool"
	TypeInt32  = "int32"
)

type AddTypeOption struct {
	Legacy    bool
	Indexed   bool
	NoMessage bool
}

// AddType adds a new type stype to scaffolded app by using optional type fields.
func (s *Scaffolder) AddType(addTypeOptions AddTypeOption, moduleName string, typeName string, fields ...string) error {
	version, err := s.version()
	if err != nil {
		return err
	}
	majorVersion := version.Major()
	path, err := gomodulepath.ParseAt(s.path)
	if err != nil {
		return err
	}

	// If no module is provided, we add the type to the app's module
	if moduleName == "" {
		moduleName = path.Package
	}
	ok, err := moduleExists(s.path, moduleName)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("the module %s doesn't exist", moduleName)
	}

	// Ensure the type name is valid, otherwise it would generate an incorrect code
	if isForbiddenComponentName(typeName) {
		return fmt.Errorf("%s can't be used as a type name", typeName)
	}

	// Check component name is not already used
	ok, err = isComponentCreated(s.path, moduleName, typeName)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("%s component is already added", typeName)
	}

	// Parse provided field
	tFields, err := parseFields(fields, isForbiddenTypeField)
	if err != nil {
		return err
	}

	var (
		g    *genny.Generator
		opts = &typed.Options{
			AppName:    path.Package,
			ModulePath: path.RawPath,
			ModuleName: moduleName,
			OwnerName:  owner(path.RawPath),
			TypeName:   typeName,
			Fields:     tFields,
			Legacy:     addTypeOptions.Legacy,
			NoMessage:  addTypeOptions.NoMessage,
		}
	)
	// generate depending on the version
	if majorVersion == cosmosver.Launchpad {
		if addTypeOptions.Indexed {
			return errors.New("indexed types not supported on Launchpad")
		}

		g, err = typed.NewLaunchpad(opts)
	} else {
		// Check and support MsgServer convention
		if err := supportMsgServer(
			s.path,
			&modulecreate.MsgServerOptions{
				ModuleName: opts.ModuleName,
				ModulePath: opts.ModulePath,
				AppName:    opts.AppName,
				OwnerName:  opts.OwnerName,
			},
		); err != nil {
			return err
		}

		// Check if indexed type
		if addTypeOptions.Indexed {
			g, err = indexed.NewStargate(opts)
		} else {
			// Scaffolding a type with ID
			g, err = typed.NewStargate(opts)
		}
	}
	if err != nil {
		return err
	}
	run := genny.WetRunner(context.Background())
	run.With(g)
	if err := run.Run(); err != nil {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := s.protoc(pwd, path.RawPath, majorVersion); err != nil {
		return err
	}
	return fmtProject(pwd)
}

// parseFields parses the provided fields, analyses the types and checks there is no duplicated field
func parseFields(fields []string, isForbiddenField func(string) bool) ([]typed.Field, error) {
	// Used to check duplicated field
	existingFields := make(map[string]bool)

	var tFields []typed.Field
	for _, f := range fields {
		fs := strings.Split(f, ":")
		name := fs[0]

		// Ensure the field name is not a Go reserved name, it would generate an incorrect code
		if isForbiddenField(name) {
			return tFields, fmt.Errorf("%s can't be used as a field name", name)
		}

		// Ensure the field is not duplicated
		if _, exists := existingFields[name]; exists {
			return tFields, fmt.Errorf("the field %s is duplicated", name)
		}
		existingFields[name] = true

		datatypeName, datatype := TypeString, TypeString
		acceptedTypes := map[string]string{
			"string": TypeString,
			"bool":   TypeBool,
			"int":    TypeInt32,
		}
		isTypeSpecified := len(fs) == 2
		if isTypeSpecified {
			if t, ok := acceptedTypes[fs[1]]; ok {
				datatype = t
				datatypeName = fs[1]
			} else {
				return tFields, fmt.Errorf("the field type %s doesn't exist", fs[1])
			}
		}
		tFields = append(tFields, typed.Field{
			Name:         name,
			Datatype:     datatype,
			DatatypeName: datatypeName,
		})
	}

	return tFields, nil
}

func isTypeCreated(appPath, moduleName, typeName string) (isCreated bool, err error) {
	abspath, err := filepath.Abs(filepath.Join(appPath, "x", moduleName, "types"))
	if err != nil {
		return false, err
	}
	fset := token.NewFileSet()
	all, err := parser.ParseDir(fset, abspath, func(os.FileInfo) bool { return true }, parser.ParseComments)
	if err != nil {
		return false, err
	}
	// To check if the file is created, we check if the message MsgCreate[TypeName] or Msg[TypeName] is defined
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
				if ("MsgCreate"+strings.Title(typeName) != typeSpec.Name.Name) && ("Msg"+strings.Title(typeName) != typeSpec.Name.Name) {
					return true
				}
				isCreated = true
				return false
			})
		}
	}
	return
}

// isForbiddenTypeField returns true if the name is forbidden as a field name
func isForbiddenTypeField(name string) bool {
	switch name {
	case
		"id",
		"index",
		"appendedValue",
		"creator":
		return true
	}

	return isGoReservedWord(name)
}
