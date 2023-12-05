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

	"github.com/ignite/cli/v28/ignite/pkg/cache"
	"github.com/ignite/cli/v28/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/templates/field"
	modulecreate "github.com/ignite/cli/v28/ignite/templates/module/create"
)

// CreateConfigs creates a new configs in the scaffolded module.
func (s Scaffolder) CreateConfigs(
	ctx context.Context,
	cacheStorage cache.Storage,
	tracer *placeholder.Tracer,
	moduleName string,
	configs ...string,
) (sm xgenny.SourceModification, err error) {
	appName := s.modpath.Package
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

	// Parse config with the associated type
	configsFields, err := field.ParseFields(configs, checkForbiddenTypeIndex)
	if err != nil {
		return sm, err
	}

	if err := checkConfigCreated(s.path, appName, moduleName, configsFields); err != nil {
		return sm, err
	}

	opts := modulecreate.ConfigsOptions{
		ModuleName: moduleName,
		Configs:    configsFields,
		AppName:    s.modpath.Package,
		AppPath:    s.path,
	}

	g, err := modulecreate.NewModuleConfigs(opts)
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

// checkConfigCreated checks if the config has been already created.
func checkConfigCreated(appPath, appName, moduleName string, configs field.Fields) (err error) {
	absPath, err := filepath.Abs(filepath.Join(appPath, "api", appName, moduleName, "module"))
	if err != nil {
		return err
	}
	fileSet := token.NewFileSet()
	all, err := parser.ParseDir(fileSet, absPath, func(os.FileInfo) bool { return true }, parser.ParseComments)
	if err != nil {
		return err
	}

	configsName := make(map[string]struct{})
	for _, config := range configs {
		configsName[config.Name.LowerCase] = struct{}{}
	}

	for _, pkg := range all {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(x ast.Node) bool {
				typeSpec, ok := x.(*ast.TypeSpec)
				if !ok {
					return true
				}

				if _, ok := typeSpec.Type.(*ast.StructType); !ok ||
					typeSpec.Name.Name != "Module" ||
					typeSpec.Type == nil {
					return true
				}

				// Check if the struct has fields.
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					return true
				}

				// Iterate through the fields of the struct.
				for _, configField := range structType.Fields.List {
					for _, fieldName := range configField.Names {
						if _, ok := configsName[strings.ToLower(fieldName.Name)]; !ok {
							continue
						}
						err = fmt.Errorf(
							"config field '%s' already exist for module %s",
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
