package modulecreate

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	appanalysis "github.com/tendermint/starport/starport/pkg/cosmosanalysis/app"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xstrings"
	"github.com/tendermint/starport/starport/templates/module"
)

// NewStargate returns the generator to scaffold a module inside a Stargate app
func NewStargate(opts *CreateOptions) (*genny.Generator, error) {
	g := genny.New()
	if err := g.Box(msgServerTemplate); err != nil {
		return g, err
	}
	if err := g.Box(stargateTemplate); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("ownerName", opts.OwnerName)
	ctx.Set("title", strings.Title)
	ctx.Set("dependencies", opts.Dependencies)

	// Used for proto package name
	ctx.Set("formatOwnerName", xstrings.FormatUsername)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	return g, nil
}

// NewStargateAppModify returns generator with modifications required to register a module in the app.
func NewStargateAppModify(replacer placeholder.Replacer, opts *CreateOptions) *genny.Generator {
	g := genny.New()
	g.RunFn(appModifyStargate(replacer, opts))
	if opts.IsIBC {
		g.RunFn(appIBCModify(replacer, opts))
	}
	return g
}

// app.go modification on Stargate when creating a module
func appModifyStargate(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := module.PathAppGo
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		template := `%[1]v
		%[2]vmodule "%[3]v/x/%[2]v"
		%[2]vmodulekeeper "%[3]v/x/%[2]v/keeper"
		%[2]vmoduletypes "%[3]v/x/%[2]v/types"`
		replacement := fmt.Sprintf(template, module.PlaceholderSgAppModuleImport, opts.ModuleName, opts.ModulePath)
		content := replacer.Replace(f.String(), module.PlaceholderSgAppModuleImport, replacement)

		// ModuleBasic
		template = `%[1]v
		%[2]vmodule.AppModuleBasic{},`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppModuleBasic, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderSgAppModuleBasic, replacement)

		// Keeper declaration
		var scopedKeeperDeclaration string
		if opts.IsIBC {
			// Scoped keeper declaration for IBC module
			// We set this placeholder so it is modified by the IBC module scaffolder
			scopedKeeperDeclaration = module.PlaceholderIBCAppScopedKeeperDeclaration
		}
		template = `%[1]v
		%[3]v
		%[4]vKeeper %[2]vmodulekeeper.Keeper`
		replacement = fmt.Sprintf(
			template,
			module.PlaceholderSgAppKeeperDeclaration,
			opts.ModuleName,
			scopedKeeperDeclaration,
			strings.Title(opts.ModuleName),
		)
		content = replacer.Replace(content, module.PlaceholderSgAppKeeperDeclaration, replacement)

		// Store key
		template = `%[1]v
		%[2]vmoduletypes.StoreKey,`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppStoreKey, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderSgAppStoreKey, replacement)

		// Module dependencies
		var depArgs string
		for _, dep := range opts.Dependencies {
			// if module has dependencies, we must verify the keeper of the dependency is defined in app.go
			if err := appanalysis.CheckKeeper(path, dep.KeeperName); err != nil {
				replacer.AppendMiscError(fmt.Sprintf(
					"the module cannot have %s as a dependency: %s",
					dep.Name,
					err.Error(),
				))
			}
			depArgs = fmt.Sprintf("%sapp.%s,\n", depArgs, dep.KeeperName)

			// If bank is a dependency, add account permissions to the module
			if dep.Name == "bank" {
				template = `%[1]v
				%[2]vmoduletypes.ModuleName: {authtypes.Minter, authtypes.Burner, authtypes.Staking},`

				replacement = fmt.Sprintf(
					template,
					module.PlaceholderSgAppMaccPerms,
					opts.ModuleName,
				)
				content = replacer.Replace(content, module.PlaceholderSgAppMaccPerms, replacement)
			}
		}

		// Keeper definition
		var scopedKeeperDefinition string
		var ibcKeeperArgument string
		if opts.IsIBC {
			// Scoped keeper definition for IBC module
			// We set this placeholder so it is modified by the IBC module scaffolder
			scopedKeeperDefinition = module.PlaceholderIBCAppScopedKeeperDefinition
			ibcKeeperArgument = module.PlaceholderIBCAppKeeperArgument
		}
		template = `%[3]v
		app.%[5]vKeeper = *%[2]vmodulekeeper.NewKeeper(
			appCodec,
			keys[%[2]vmoduletypes.StoreKey],
			keys[%[2]vmoduletypes.MemStoreKey],
			%[4]v
			%[6]v)
		%[2]vModule := %[2]vmodule.NewAppModule(appCodec, app.%[5]vKeeper)

		%[1]v`
		replacement = fmt.Sprintf(
			template,
			module.PlaceholderSgAppKeeperDefinition,
			opts.ModuleName,
			scopedKeeperDefinition,
			ibcKeeperArgument,
			strings.Title(opts.ModuleName),
			depArgs,
		)
		content = replacer.Replace(content, module.PlaceholderSgAppKeeperDefinition, replacement)

		// App Module
		template = `%[1]v
		%[2]vModule,`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppAppModule, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderSgAppAppModule, replacement)

		// Init genesis
		template = `%[1]v
		%[2]vmoduletypes.ModuleName,`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppInitGenesis, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderSgAppInitGenesis, replacement)

		// Param subspace
		template = `%[1]v
		paramsKeeper.Subspace(%[2]vmoduletypes.ModuleName)`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppParamSubspace, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderSgAppParamSubspace, replacement)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
