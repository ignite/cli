package module

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
)

// these needs to be created in the compiler time, otherwise packr2 won't be
// able to find boxes.
var templates = map[cosmosver.MajorVersion]*packr.Box{
	cosmosver.Launchpad: packr.New("module/templates/launchpad", "./launchpad"),
	cosmosver.Stargate:  packr.New("module/templates/stargate", "./stargate"),
}

// New ...
func NewCreateStargate(opts *CreateOptions) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(createAppModifyStargate(opts))

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
func createAppModifyStargate(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := PathAppGo
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		template := `%[1]v
		"%[3]v/x/%[2]v"
		%[2]vkeeper "%[3]v/x/%[2]v/keeper"
		%[2]vtypes "%[3]v/x/%[2]v/types"`
		replacement := fmt.Sprintf(template, placeholderSgAppModuleImport, opts.ModuleName, opts.ModulePath)
		content := strings.Replace(f.String(), placeholderSgAppModuleImport, replacement, 1)

		// ModuleBasic
		template = `%[1]v
		%[2]v.AppModuleBasic{},`
		replacement = fmt.Sprintf(template, placeholderSgAppModuleBasic, opts.ModuleName)
		content = strings.Replace(content, placeholderSgAppModuleBasic, replacement, 1)

		// Keeper declaration
		template = `%[1]v
		%[2]vKeeper %[2]vkeeper.Keeper`
		replacement = fmt.Sprintf(template, placeholderSgAppKeeperDeclaration, opts.ModuleName)
		content = strings.Replace(content, placeholderSgAppKeeperDeclaration, replacement, 1)

		// Store key
		template = `%[1]v
		%[2]vtypes.StoreKey,`
		replacement = fmt.Sprintf(template, placeholderSgAppStoreKey, opts.ModuleName)
		content = strings.Replace(content, placeholderSgAppStoreKey, replacement, 1)

		// Keeper definition
		template = `%[1]v
		app.%[2]vKeeper = *%[2]vkeeper.NewKeeper(
			appCodec,
			keys[%[2]vtypes.StoreKey],
			keys[%[2]vtypes.MemStoreKey],
		)`
		replacement = fmt.Sprintf(template, placeholderSgAppKeeperDefinition, opts.ModuleName)
		content = strings.Replace(content, placeholderSgAppKeeperDefinition, replacement, 1)

		// App Module
		template = `%[1]v
		%[2]v.NewAppModule(appCodec, app.%[2]vKeeper),`
		replacement = fmt.Sprintf(template, placeholderSgAppAppModule, opts.ModuleName)
		content = strings.Replace(content, placeholderSgAppAppModule, replacement, 1)

		// Init genesis
		template = `%[1]v
		%[2]vtypes.ModuleName,`
		replacement = fmt.Sprintf(template, placeholderSgAppInitGenesis, opts.ModuleName)
		content = strings.Replace(content, placeholderSgAppInitGenesis, replacement, 1)

		// Param subspace
		template = `%[1]v
		paramsKeeper.Subspace(%[2]vtypes.ModuleName)`
		replacement = fmt.Sprintf(template, placeholderSgAppParamSubspace, opts.ModuleName)
		content = strings.Replace(content, placeholderSgAppParamSubspace, replacement, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
