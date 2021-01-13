package module_create

import (
	"fmt"
	"github.com/tendermint/starport/starport/templates/module"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
)

// New ...
func NewCreateStargate(opts *CreateOptions) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(appModifyStargate(opts))

	if err := g.Box(templates[cosmosver.Stargate]); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("ownerName", opts.OwnerName)
	ctx.Set("title", strings.Title)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	return g, nil
}

// app.go modification on Stargate when creating a module
func appModifyStargate(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := module.PathAppGo
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		template := `%[1]v
		"%[3]v/x/%[2]v"
		%[2]vkeeper "%[3]v/x/%[2]v/keeper"
		%[2]vtypes "%[3]v/x/%[2]v/types"`
		replacement := fmt.Sprintf(template, module.PlaceholderSgAppModuleImport, opts.ModuleName, opts.ModulePath)
		content := strings.Replace(f.String(), module.PlaceholderSgAppModuleImport, replacement, 1)

		// ModuleBasic
		template = `%[1]v
		%[2]v.AppModuleBasic{},`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppModuleBasic, opts.ModuleName)
		content = strings.Replace(content, module.PlaceholderSgAppModuleBasic, replacement, 1)

		// Keeper declaration
		template = `%[1]v
		%[2]vKeeper %[2]vkeeper.Keeper`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppKeeperDeclaration, opts.ModuleName)
		content = strings.Replace(content, module.PlaceholderSgAppKeeperDeclaration, replacement, 1)

		// Store key
		template = `%[1]v
		%[2]vtypes.StoreKey,`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppStoreKey, opts.ModuleName)
		content = strings.Replace(content, module.PlaceholderSgAppStoreKey, replacement, 1)

		// Keeper definition
		template = `%[1]v
		app.%[2]vKeeper = *%[2]vkeeper.NewKeeper(
			appCodec,
			keys[%[2]vtypes.StoreKey],
			keys[%[2]vtypes.MemStoreKey],
		)`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppKeeperDefinition, opts.ModuleName)
		content = strings.Replace(content, module.PlaceholderSgAppKeeperDefinition, replacement, 1)

		// App Module
		template = `%[1]v
		%[2]v.NewAppModule(appCodec, app.%[2]vKeeper),`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppAppModule, opts.ModuleName)
		content = strings.Replace(content, module.PlaceholderSgAppAppModule, replacement, 1)

		// Init genesis
		template = `%[1]v
		%[2]vtypes.ModuleName,`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppInitGenesis, opts.ModuleName)
		content = strings.Replace(content, module.PlaceholderSgAppInitGenesis, replacement, 1)

		// Param subspace
		template = `%[1]v
		paramsKeeper.Subspace(%[2]vtypes.ModuleName)`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppParamSubspace, opts.ModuleName)
		content = strings.Replace(content, module.PlaceholderSgAppParamSubspace, replacement, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
