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
func NewCreateLaunchpad(opts *CreateOptions) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(appModifyLaunchpad(opts))

	if err := g.Box(templates[cosmosver.Launchpad]); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("title", strings.Title)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	return g, nil
}

// app.go modification on Launchpad when creating a module
func appModifyLaunchpad(opts *CreateOptions) genny.RunFn {
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
		replacement := fmt.Sprintf(template, module.Placeholder, opts.ModuleName, opts.ModulePath)
		content := strings.Replace(f.String(), module.Placeholder, replacement, 1)

		// ModuleBasic
		template = `%[1]v
		%[2]v.AppModuleBasic{},`
		replacement = fmt.Sprintf(template, module.Placeholder2, opts.ModuleName)
		content = strings.Replace(content, module.Placeholder2, replacement, 1)

		// Keeper declaration
		template = `%[1]v
		%[2]vKeeper %[2]vkeeper.Keeper`
		replacement = fmt.Sprintf(template, module.Placeholder3, opts.ModuleName)
		content = strings.Replace(content, module.Placeholder3, replacement, 1)

		// Store key
		template = `%[1]v
		%[2]vtypes.StoreKey,`
		replacement = fmt.Sprintf(template, module.Placeholder5, opts.ModuleName)
		content = strings.Replace(content, module.Placeholder5, replacement, 1)

		// Param subspace
		template = `%[1]v
		app.subspaces[%[2]vtypes.ModuleName] = app.paramsKeeper.Subspace(%[2]vtypes.DefaultParamspace)`
		replacement = fmt.Sprintf(template, module.Placeholder5_1, opts.ModuleName)
		content = strings.Replace(content, module.Placeholder5_1, replacement, 1)

		// Keeper definition
		template = `%[1]v
		app.%[2]vKeeper = %[2]vkeeper.NewKeeper(
			app.cdc,
			keys[%[2]vtypes.StoreKey],
			app.subspaces[%[2]vtypes.ModuleName],
		)`
		replacement = fmt.Sprintf(template, module.Placeholder5_2, opts.ModuleName)
		content = strings.Replace(content, module.Placeholder5_2, replacement, 1)

		// Module manager
		template = `%[1]v
		%[2]v.NewAppModule(app.%[2]vKeeper),`
		replacement = fmt.Sprintf(template, module.Placeholder6, opts.ModuleName)
		content = strings.Replace(content, module.Placeholder6, replacement, 1)

		// Genesis
		template = `%[1]v
		%[2]vtypes.ModuleName,`
		replacement = fmt.Sprintf(template, module.Placeholder7, opts.ModuleName)
		content = strings.Replace(content, module.Placeholder7, replacement, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}


