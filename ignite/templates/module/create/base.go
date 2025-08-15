package modulecreate

import (
	"fmt"
	"io/fs"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"
	"github.com/iancoleman/strcase"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/pkg/xstrings"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
	"github.com/ignite/cli/v29/ignite/templates/module"
)

// NewGenerator returns the generator to scaffold a module inside an app.
func NewGenerator(opts *CreateOptions) (*genny.Generator, error) {
	subBase, err := fs.Sub(fsBase, "files/base")
	if err != nil {
		return nil, errors.Errorf("fail to generate sub: %w", err)
	}

	g := genny.New()
	if err := g.OnlyFS(subBase, nil, nil); err != nil {
		return g, err
	}

	appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)

	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("protoVer", opts.ProtoVer)
	ctx.Set("dependencies", opts.Dependencies)
	ctx.Set("params", opts.Params)
	ctx.Set("configs", opts.Configs)
	ctx.Set("isIBC", opts.IsIBC)
	ctx.Set("apiPath", fmt.Sprintf("/%s/%s/%s", appModulePath, opts.ModuleName, opts.ProtoVer))
	ctx.Set("protoPkgName", module.ProtoPackageName(appModulePath, opts.ModuleName, opts.ProtoVer))
	ctx.Set("protoModulePkgName", module.ProtoModulePackageName(appModulePath, opts.ModuleName, opts.ProtoVer))
	ctx.Set("toVariableName", strcase.ToLowerCamel)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{protoDir}}", opts.ProtoDir))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{protoVer}}", opts.ProtoVer))

	return g, nil
}

// NewAppModify returns generator with modifications required to register a module in the app.
func NewAppModify(replacer placeholder.Replacer, opts *CreateOptions) *genny.Generator {
	g := genny.New()
	g.RunFn(appModify(opts))
	g.RunFn(appConfigModify(replacer, opts))
	if opts.IsIBC {
		g.RunFn(appIBCModify(replacer, opts))
	}
	return g
}

func appConfigModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		configPath := module.PathAppConfigGo
		fConfig, err := r.Disk.Find(configPath)
		if err != nil {
			return err
		}

		// Import
		content, err := xast.AppendImports(
			fConfig.String(),
			xast.WithNamedImport(
				"_",
				fmt.Sprintf("%[1]v/x/%[2]v/module", opts.ModulePath, opts.ModuleName),
			),
			xast.WithNamedImport(
				fmt.Sprintf("%[1]vmoduletypes", opts.ModuleName),
				fmt.Sprintf("%[1]v/x/%[2]v/types", opts.ModulePath, opts.ModuleName),
			),
		)
		if err != nil {
			return err
		}

		// Init genesis
		template := `%[2]vmoduletypes.ModuleName,
%[1]v`
		replacement := fmt.Sprintf(template, module.PlaceholderSgAppInitGenesis, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderSgAppInitGenesis, replacement)
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppBeginBlockers, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderSgAppBeginBlockers, replacement)
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppEndBlockers, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderSgAppEndBlockers, replacement)

		template = `{
				Name:   %[2]vmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&%[2]vmoduletypes.Module{}),
			},
%[1]v`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppModuleConfig, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderSgAppModuleConfig, replacement)

		// Module dependencies
		for _, dep := range opts.Dependencies {
			// If bank is a dependency, add account permissions to the module
			if dep.Name == "Bank" {
				replacement = fmt.Sprintf(
					"{Account: %[1]vmoduletypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner, authtypes.Staking}}",
					opts.ModuleName,
				)

				// Keeper definition
				content, err = xast.ModifyGlobalArrayVar(content, "moduleAccPerms", xast.AppendGlobalArrayValue(replacement))
				if err != nil {
					return err
				}

			}
		}

		newFile := genny.NewFileS(configPath, content)

		return r.File(newFile)
	}
}

// app.go modification when creating a module.
func appModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		appPath := module.PathAppGo
		f, err := r.Disk.Find(appPath)
		if err != nil {
			return err
		}

		// Import
		content, err := xast.AppendImports(
			f.String(),
			xast.WithNamedImport(
				fmt.Sprintf("%[1]vmodulekeeper", opts.ModuleName),
				fmt.Sprintf("%[1]v/x/%[2]v/keeper", opts.ModulePath, opts.ModuleName),
			),
		)
		if err != nil {
			return err
		}

		// Keeper declaration
		content, err = xast.ModifyStruct(
			content,
			"App",
			xast.AppendStructValue(
				fmt.Sprintf("%[1]vKeeper", xstrings.Title(opts.ModuleName)),
				fmt.Sprintf("%[1]vmodulekeeper.Keeper", opts.ModuleName),
			),
		)
		if err != nil {
			return err
		}

		// Keeper definition
		content, err = xast.ModifyFunction(
			content,
			"New",
			xast.AppendInsideFuncCall(
				"Inject",
				fmt.Sprintf("\n&app.%[1]vKeeper", xstrings.Title(opts.ModuleName)),
				-1,
			),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(appPath, content)
		return r.File(newFile)
	}
}
