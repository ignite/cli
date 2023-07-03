package modulecreate

import (
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"
	"github.com/iancoleman/strcase"

	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/pkg/xstrings"
	"github.com/ignite/cli/ignite/templates/field/plushhelpers"
	"github.com/ignite/cli/ignite/templates/module"
)

// NewGenerator returns the generator to scaffold a module inside an app.
func NewGenerator(opts *CreateOptions) (*genny.Generator, error) {
	var (
		g = genny.New()

		msgServerTemplate = xgenny.NewEmbedWalker(
			fsMsgServer,
			"files/msgserver/",
			opts.AppPath,
		)
		genesisTestTemplate = xgenny.NewEmbedWalker(
			fsGenesisTest,
			"files/genesistest/",
			opts.AppPath,
		)
		baseTemplate = xgenny.NewEmbedWalker(
			fsBase,
			"files/base/",
			opts.AppPath,
		)
	)

	if err := g.Box(msgServerTemplate); err != nil {
		return g, err
	}
	if err := g.Box(genesisTestTemplate); err != nil {
		return g, err
	}
	if err := g.Box(baseTemplate); err != nil {
		return g, err
	}

	appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)

	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("dependencies", opts.Dependencies)
	ctx.Set("params", opts.Params)
	ctx.Set("isIBC", opts.IsIBC)
	ctx.Set("apiPath", fmt.Sprintf("/%s/%s", appModulePath, opts.ModuleName))
	ctx.Set("protoPkgName", module.ProtoPackageName(appModulePath, opts.ModuleName))
	ctx.Set("toVariableName", strcase.ToLowerCamel)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))

	gSimapp, err := AddSimulation(opts.AppPath, opts.ModulePath, opts.ModuleName, opts.Params...)
	if err != nil {
		return g, err
	}
	g.Merge(gSimapp)

	return g, nil
}

// NewAppModify returns generator with modifications required to register a module in the app.
func NewAppModify(replacer placeholder.Replacer, opts *CreateOptions) *genny.Generator {
	g := genny.New()
	g.RunFn(appModify(replacer, opts))
	g.RunFn(appConfigModify(replacer, opts))
	if opts.IsIBC {
		g.RunFn(appIBCModify(replacer, opts))
	}
	return g
}

func appConfigModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		configPath := filepath.Join(opts.AppPath, module.PathAppConfigGo)
		fConfig, err := r.Disk.Find(configPath)
		if err != nil {
			return err
		}

		// Import
		template := `%[2]vmoduletypes "%[3]v/x/%[2]v/types"
%[2]vmodulev1 "%[3]v/api/%[4]v/%[2]v/module"
%[1]v`
		replacement := fmt.Sprintf(template, module.PlaceholderSgAppModuleImport, opts.ModuleName, opts.ModulePath, opts.AppName)
		content := replacer.Replace(fConfig.String(), module.PlaceholderSgAppModuleImport, replacement)

		// Init genesis
		template = `%[2]vmoduletypes.ModuleName,
%[1]v`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppInitGenesis, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderSgAppInitGenesis, replacement)
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppBeginBlockers, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderSgAppBeginBlockers, replacement)
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppEndBlockers, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderSgAppEndBlockers, replacement)

		template = `{
				Name:   %[2]vmoduletypes.ModuleName,
				Config: appconfig.WrapAny(&%[2]vmodulev1.Module{}),
			},
%[1]v`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppModuleConfig, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderSgAppModuleConfig, replacement)

		// Module dependencies
		for _, dep := range opts.Dependencies {
			// If bank is a dependency, add account permissions to the module
			if dep.Name == "Bank" {
				template = `{Account: %[2]vmoduletypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner, authtypes.Staking}},
%[1]v`

				replacement = fmt.Sprintf(
					template,
					module.PlaceholderSgAppMaccPerms,
					opts.ModuleName,
				)
				content = replacer.Replace(content, module.PlaceholderSgAppMaccPerms, replacement)
			}
		}

		newFile := genny.NewFileS(configPath, content)

		return r.File(newFile)
	}
}

// app.go modification when creating a module.
func appModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		appPath := filepath.Join(opts.AppPath, module.PathAppGo)
		f, err := r.Disk.Find(appPath)
		if err != nil {
			return err
		}

		// Import
		template := `%[2]vmodule "%[3]v/x/%[2]v"
		%[2]vmodulekeeper "%[3]v/x/%[2]v/keeper"
%[1]v`
		replacement := fmt.Sprintf(template, module.PlaceholderSgAppModuleImport, opts.ModuleName, opts.ModulePath)
		content := replacer.Replace(f.String(), module.PlaceholderSgAppModuleImport, replacement)

		// ModuleBasic
		template = `%[2]vmodule.AppModuleBasic{},
%[1]v`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppModuleBasic, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderSgAppModuleBasic, replacement)

		// Keeper declaration
		template = `%[2]vKeeper %[3]vmodulekeeper.Keeper
%[1]v`
		replacement = fmt.Sprintf(
			template,
			module.PlaceholderSgAppKeeperDeclaration,
			xstrings.Title(opts.ModuleName),
			opts.ModuleName,
		)
		content = replacer.Replace(content, module.PlaceholderSgAppKeeperDeclaration, replacement)

		// Keeper definition
		template = `&app.%[2]vKeeper,
		%[1]v`
		replacement = fmt.Sprintf(
			template,
			module.PlaceholderSgAppKeeperDefinition,
			xstrings.Title(opts.ModuleName),
		)
		content = replacer.Replace(content, module.PlaceholderSgAppKeeperDefinition, replacement)

		newFile := genny.NewFileS(appPath, content)
		return r.File(newFile)
	}
}
