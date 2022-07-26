package modulecreate

import (
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"

	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/pkg/xstrings"
	"github.com/ignite/cli/ignite/templates/field/plushhelpers"
	"github.com/ignite/cli/ignite/templates/module"
	"github.com/ignite/cli/ignite/templates/typed"
)

// NewIBC returns the generator to scaffold the implementation of the IBCModule interface inside a module
func NewIBC(replacer placeholder.Replacer, opts *CreateOptions) (*genny.Generator, error) {
	var (
		g        = genny.New()
		template = xgenny.NewEmbedWalker(fsIBC, "ibc/", opts.AppPath)
	)

	g.RunFn(genesisModify(replacer, opts))
	g.RunFn(genesisTypesModify(replacer, opts))
	g.RunFn(genesisProtoModify(replacer, opts))
	g.RunFn(keysModify(replacer, opts))

	if err := g.Box(template); err != nil {
		return g, err
	}

	appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)

	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("ibcOrdering", opts.IBCOrdering)
	ctx.Set("dependencies", opts.Dependencies)
	ctx.Set("protoPkgName", module.ProtoPackageName(appModulePath, opts.ModuleName))

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	return g, nil
}

func genesisModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Genesis init
		templateInit := `%s
k.SetPort(ctx, genState.PortId)
// Only try to bind to port if it is not already bound, since we may already own
// port capability from capability InitGenesis
if !k.IsBound(ctx, genState.PortId) {
	// module binds to the port on InitChain
	// and claims the returned capability
	err := k.BindPort(ctx, genState.PortId)
	if err != nil {
		panic("could not claim port capability: " + err.Error())
	}
}`
		replacementInit := fmt.Sprintf(templateInit, typed.PlaceholderGenesisModuleInit)
		content := replacer.Replace(f.String(), typed.PlaceholderGenesisModuleInit, replacementInit)

		// Genesis export
		templateExport := `genesis.PortId = k.GetPort(ctx)
%s`
		replacementExport := fmt.Sprintf(templateExport, typed.PlaceholderGenesisModuleExport)
		content = replacer.Replace(content, typed.PlaceholderGenesisModuleExport, replacementExport)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTypesModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		templateImport := `host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
%s`
		replacementImport := fmt.Sprintf(templateImport, typed.PlaceholderGenesisTypesImport)
		content := replacer.Replace(f.String(), typed.PlaceholderGenesisTypesImport, replacementImport)

		// Default genesis
		templateDefault := `PortId: PortID,
%s`
		replacementDefault := fmt.Sprintf(templateDefault, typed.PlaceholderGenesisTypesDefault)
		content = replacer.Replace(content, typed.PlaceholderGenesisTypesDefault, replacementDefault)

		// Validate genesis
		// PlaceholderIBCGenesisTypeValidate
		templateValidate := `if err := host.PortIdentifierValidator(gs.PortId); err != nil {
	return err
}
%s`
		replacementValidate := fmt.Sprintf(templateValidate, typed.PlaceholderGenesisTypesValidate)
		content = replacer.Replace(content, typed.PlaceholderGenesisTypesValidate, replacementValidate)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisProtoModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "proto", opts.ModuleName, "genesis.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Determine the new field number
		content := f.String()

		template := `string port_id = 2;
  %s`
		replacement := fmt.Sprintf(template, typed.PlaceholderGenesisProtoState)
		content = replacer.Replace(content, typed.PlaceholderGenesisProtoState, replacement)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func keysModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "x", opts.ModuleName, "types/keys.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Append version and the port ID in keys
		templateName := `// Version defines the current version the IBC module supports
Version = "%[1]v-1"

// PortID is the default port id that module binds to
PortID = "%[1]v"`
		replacementName := fmt.Sprintf(templateName, opts.ModuleName)
		content := replacer.Replace(f.String(), module.PlaceholderIBCKeysName, replacementName)

		// PlaceholderIBCKeysPort
		templatePort := `var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("%[1]v-port-")
)`
		replacementPort := fmt.Sprintf(templatePort, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderIBCKeysPort, replacementPort)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func appIBCModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, module.PathAppGo)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Add route to IBC router
		templateRouter := `ibcRouter.AddRoute(%[2]vmoduletypes.ModuleName, %[2]vModule)
%[1]v`
		replacementRouter := fmt.Sprintf(
			templateRouter,
			module.PlaceholderIBCAppRouter,
			opts.ModuleName,
		)
		content := replacer.Replace(f.String(), module.PlaceholderIBCAppRouter, replacementRouter)

		// Scoped keeper declaration for the module
		templateScopedKeeperDeclaration := `Scoped%[1]vKeeper capabilitykeeper.ScopedKeeper`
		replacementScopedKeeperDeclaration := fmt.Sprintf(templateScopedKeeperDeclaration, xstrings.Title(opts.ModuleName))
		content = replacer.Replace(content, module.PlaceholderIBCAppScopedKeeperDeclaration, replacementScopedKeeperDeclaration)

		// Scoped keeper definition
		templateScopedKeeperDefinition := `scoped%[1]vKeeper := app.CapabilityKeeper.ScopeToModule(%[2]vmoduletypes.ModuleName)
app.Scoped%[1]vKeeper = scoped%[1]vKeeper`
		replacementScopedKeeperDefinition := fmt.Sprintf(
			templateScopedKeeperDefinition,
			xstrings.Title(opts.ModuleName),
			opts.ModuleName,
		)
		content = replacer.Replace(content, module.PlaceholderIBCAppScopedKeeperDefinition, replacementScopedKeeperDefinition)

		// New argument passed to the module keeper
		templateKeeperArgument := `app.IBCKeeper.ChannelKeeper,
&app.IBCKeeper.PortKeeper,
scoped%[1]vKeeper,`
		replacementKeeperArgument := fmt.Sprintf(
			templateKeeperArgument,
			xstrings.Title(opts.ModuleName),
		)
		content = replacer.Replace(content, module.PlaceholderIBCAppKeeperArgument, replacementKeeperArgument)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
