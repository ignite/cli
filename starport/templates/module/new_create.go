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

// New ...
func NewCreateStargate(opts *CreateOptions) (*genny.Generator, error) {
	g := genny.New()
	if err := g.Box(templates[cosmosver.Stargate]); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("title", strings.Title)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	return g, nil
}

func appModifyLaunchpad(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := "app/app.go"
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		template := `%[1]v
		"%[3]v/x/%[2]v"
		%[2]vkeeper "%[3]v/x/%[2]v/keeper"
		%[2]vtypes "%[3]v/x/%[2]v/types"`
		replacement := fmt.Sprintf(template, placeholder, opts.ModuleName, opts.ModulePath)
		content := strings.Replace(f.String(), placeholder, replacement, 1)

		// ModuleBasic
		template2 := `%[1]v
		%[2]v.AppModuleBasic{},`
		replacement2 := fmt.Sprintf(template2, placeholder2, opts.ModuleName)
		content = strings.Replace(content, placeholder2, replacement2, 1)

		// Keeper declaration
		template3 := `%[1]v
		%[2]vKeeper %[2]vkeeper.Keeper`
		replacement3 := fmt.Sprintf(template3, placeholder3, opts.ModuleName)
		content = strings.Replace(content, placeholder3, replacement3, 1)

		// Store key
		template5 := `%[1]v
		%[2]vtypes.StoreKey,`
		replacement5 := fmt.Sprintf(template5, placeholder5, opts.ModuleName)
		content = strings.Replace(content, placeholder5, replacement5, 1)

		// Param subspace
		template5_1 := `%[1]v
		app.subspaces[%[2]vtypes.ModuleName] = app.paramsKeeper.Subspace(%[2]vtypes.DefaultParamspace)`
		replacement5_1 := fmt.Sprintf(template5_1, placeholder5_1, opts.ModuleName)
		content = strings.Replace(content, placeholder5_1, replacement5_1, 1)

		// Keeper definition
		template5_2 := `%[1]v
		app.%[2]vKeeper = %[2]vkeeper.NewKeeper(
			app.cdc,
			keys[%[2]vtypes.StoreKey],
			app.subspaces[%[2]vtypes.ModuleName],
		)`
		replacement5_2 := fmt.Sprintf(template5_2, placeholder5_2, opts.ModuleName)
		content = strings.Replace(content, placeholder5_2, replacement5_2, 1)

		// Module manager
		template6 := `%[1]v
		%[2]v.NewAppModule(app.%[2]vKeeper),`
		replacement6 := fmt.Sprintf(template6, placeholder6, opts.ModuleName)
		content = strings.Replace(content, placeholder6, replacement6, 1)

		// Genesis
		template7 := `%[1]v
		%[2]vtypes.ModuleName,`
		replacement7 := fmt.Sprintf(template7, placeholder7, opts.ModuleName)
		content = strings.Replace(content, placeholder7, replacement7, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
