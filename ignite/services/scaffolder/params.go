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

	"github.com/gobuffalo/genny/v2"

	"github.com/ignite/cli/ignite/pkg/cache"
	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/field"
	modulecreate "github.com/ignite/cli/ignite/templates/module/create"
)

// CreateParams creates a new params in the scaffolded module.
func (s Scaffolder) CreateParams(
	ctx context.Context,
	cacheStorage cache.Storage,
	tracer *placeholder.Tracer,
	moduleName string,
	params ...string,
) (sm xgenny.SourceModification, err error) {
	// If no module is provided, we add the type to the app's module
	if moduleName == "" {
		moduleName = s.modpath.Package
	}
	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return sm, err
	}
	moduleName = mfName.LowerCase

	// Check if the module already exist
	ok, err := moduleExists(s.path, moduleName)
	if err != nil {
		return sm, err
	}
	if !ok {
		return sm, fmt.Errorf("the module %v not exist", moduleName)
	}

	// Parse params with the associated type
	paramsFields, err := field.ParseFields(params, checkForbiddenTypeIndex)
	if err != nil {
		return sm, err
	}

	if err := checkParamCreated(s.path, moduleName, paramsFields); err != nil {
		return sm, err
	}

	opts := modulecreate.ParamsOptions{
		ModuleName: moduleName,
		Params:     paramsFields,
		AppName:    s.modpath.Package,
		AppPath:    s.path,
	}

	g, err := modulecreate.NewModuleParam(tracer, opts)
	if err != nil {
		return sm, err
	}
	gens := []*genny.Generator{g}

	sm, err = xgenny.RunWithValidation(tracer, gens...)
	if err != nil {
		return sm, err
	}

	return sm, finish(ctx, cacheStorage, opts.AppPath, s.modpath.RawPath)
}

// checkParamCreated checks if the parameter has been already created.
func checkParamCreated(appPath, moduleName string, params field.Fields) (err error) {
	absPath, err := filepath.Abs(filepath.Join(appPath, "x", moduleName, "types"))
	if err != nil {
		return err
	}
	fileSet := token.NewFileSet()
	all, err := parser.ParseDir(fileSet, absPath, func(os.FileInfo) bool { return true }, parser.ParseComments)
	if err != nil {
		return err
	}

	paramsName := make(map[string]struct{})
	for _, param := range params {
		paramsName[param.Name.LowerCase] = struct{}{}
	}

	for _, pkg := range all {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(x ast.Node) bool {
				typeSpec, ok := x.(*ast.TypeSpec)
				if !ok {
					return true
				}

				if _, ok := typeSpec.Type.(*ast.StructType); !ok ||
					typeSpec.Name.Name != "Params" ||
					typeSpec.Type == nil {
					return true
				}

				// Check if the struct has fields.
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					return true
				}

				// Iterate through the fields of the struct.
				for _, paramField := range structType.Fields.List {
					for _, fieldName := range paramField.Names {
						if _, ok := paramsName[strings.ToLower(fieldName.Name)]; !ok {
							continue
						}
						err = fmt.Errorf(
							"param field '%s' already exist for module %s",
							fieldName.Name,
							moduleName,
						)
						return false
					}
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
