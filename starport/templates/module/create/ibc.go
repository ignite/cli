package modulecreate

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/xstrings"
	"github.com/tendermint/starport/starport/templates/module"
)

// NewIBC returns the generator to scaffold the implementation of the IBCModule interface inside a module
func NewIBC(opts *CreateOptions) (*genny.Generator, error) {
	g := genny.New()

	g.RunFn(moduleModify(opts))
	g.RunFn(genesisModify(opts))
	g.RunFn(errorsModify(opts))
	g.RunFn(genesisTypeModify(opts))
	g.RunFn(genesisProtoModify(opts))
	g.RunFn(keysModify(opts))
	g.RunFn(keeperModify(opts))
	g.RunFn(appModify(opts))

	if err := g.Box(ibcTemplate); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("ownerName", opts.OwnerName)
	ctx.Set("ibcOrdering", opts.IBCOrdering)
	ctx.Set("title", strings.Title)

	// Used for proto package name
	ctx.Set("formatOwnerName", xstrings.FormatUsername)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	return g, nil
}

func moduleModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/module.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		templateImport := `porttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/05-port/types"`
		content := strings.Replace(f.String(), module.PlaceholderIBCModuleImport, templateImport, 1)

		// Interface to implement
		templateInterface := `_ porttypes.IBCModule   = AppModule{}`
		content = strings.Replace(content, module.PlaceholderIBCModuleInterface, templateInterface, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/genesis.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Genesis init
		templateInit := `k.SetPort(ctx, genState.PortId)
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
		content := strings.Replace(f.String(), module.PlaceholderIBCGenesisInit, templateInit, 1)

		// Genesis export
		templateExport := `genesis.PortId = k.GetPort(ctx)`
		content = strings.Replace(content, module.PlaceholderIBCGenesisExport, templateExport, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func errorsModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/errors.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// IBC errors
		template := `ErrInvalidPacketTimeout = sdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
ErrInvalidVersion = sdkerrors.Register(ModuleName, 1501, "invalid version")`
		content := strings.Replace(f.String(), module.PlaceholderIBCErrors, template, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTypeModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/genesis.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		templateImport := `host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"`
		content := strings.Replace(f.String(), module.PlaceholderIBCGenesisTypeImport, templateImport, 1)

		// Default genesis
		templateDefault := `PortId: PortID,`
		content = strings.Replace(content, module.PlaceholderIBCGenesisTypeDefault, templateDefault, 1)

		// Validate genesis
		// PlaceholderIBCGenesisTypeValidate
		templateValidate := `if err := host.PortIdentifierValidator(gs.PortId); err != nil {
	return err
}`
		content = strings.Replace(content, module.PlaceholderIBCGenesisTypeValidate, templateValidate, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisProtoModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("proto/%s/genesis.proto", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Determine the new field number
		content := f.String()
		fieldNumber := strings.Count(content, module.PlaceholderGenesisProtoStateField) + 1

		template := `string port_id = %[1]v; %[2]v`
		replacement := fmt.Sprintf(template, fieldNumber, module.PlaceholderGenesisProtoStateField)
		content = strings.Replace(content, module.PlaceholderIBCGenesisProto, replacement, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func keysModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/keys.go", opts.ModuleName)
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
		content := strings.Replace(f.String(), module.PlaceholderIBCKeysName, replacementName, 1)

		// PlaceholderIBCKeysPort
		templatePort := `var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("%[1]v-port-")
)`
		replacementPort := fmt.Sprintf(templatePort, opts.ModuleName)
		content = strings.Replace(content, module.PlaceholderIBCKeysPort, replacementPort, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func keeperModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/keeper/keeper.go", opts.ModuleName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Keeper new attributes
		templateAttribute := `channelKeeper types.ChannelKeeper
portKeeper    types.PortKeeper
scopedKeeper  types.ScopedKeeper`
		content := strings.Replace(f.String(), module.PlaceholderIBCKeeperAttribute, templateAttribute, 1)

		// New parameter for the constructor
		templateParameter := `channelKeeper types.ChannelKeeper,
portKeeper types.PortKeeper,
scopedKeeper types.ScopedKeeper,`
		content = strings.Replace(content, module.PlaceholderIBCKeeperParameter, templateParameter, 1)

		// New return values for the constructor
		templateReturn := `channelKeeper: channelKeeper,
portKeeper:    portKeeper,
scopedKeeper:  scopedKeeper,`
		content = strings.Replace(content, module.PlaceholderIBCKeeperReturn, templateReturn, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func appModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := module.PathAppGo
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Add route to IBC router
		templateRouter := `%[1]v
ibcRouter.AddRoute(%[2]vtypes.ModuleName, %[2]vModule)`
		replacementRouter := fmt.Sprintf(
			templateRouter,
			module.PlaceholderIBCAppRouter,
			opts.ModuleName,
		)
		content := strings.Replace(f.String(), module.PlaceholderIBCAppRouter, replacementRouter, 1)

		// Scoped keeper declaration for the module
		templateScopedKeeperDeclaration := `Scoped%[1]vKeeper capabilitykeeper.ScopedKeeper`
		replacementScopedKeeperDeclaration := fmt.Sprintf(templateScopedKeeperDeclaration, strings.Title(opts.ModuleName))
		content = strings.Replace(content, module.PlaceholderIBCAppScopedKeeperDeclaration, replacementScopedKeeperDeclaration, 1)

		// Scoped keeper definition
		templateScopedKeeperDefinition := `scoped%[1]vKeeper := app.CapabilityKeeper.ScopeToModule(%[2]vtypes.ModuleName)
app.Scoped%[1]vKeeper = scoped%[1]vKeeper`
		replacementScopedKeeperDefinition := fmt.Sprintf(
			templateScopedKeeperDefinition,
			strings.Title(opts.ModuleName),
			opts.ModuleName,
		)
		content = strings.Replace(content, module.PlaceholderIBCAppScopedKeeperDefinition, replacementScopedKeeperDefinition, 1)

		// New argument passed to the module keeper
		templateKeeperArgument := `app.IBCKeeper.ChannelKeeper,
&app.IBCKeeper.PortKeeper,
scoped%[1]vKeeper,`
		replacementKeeperArgument := fmt.Sprintf(
			templateKeeperArgument,
			strings.Title(opts.ModuleName),
		)
		content = strings.Replace(content, module.PlaceholderIBCAppKeeperArgument, replacementKeeperArgument, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
