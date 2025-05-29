package modulecreate

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis/protoutil"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/pkg/xstrings"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
	"github.com/ignite/cli/v29/ignite/templates/module"
)

// NewIBC returns the generator to scaffold the implementation of the IBCModule interface inside a module.
func NewIBC(replacer placeholder.Replacer, opts *CreateOptions) (*genny.Generator, error) {
	subFs, err := fs.Sub(fsIBC, "files/ibc")
	if err != nil {
		return nil, errors.Errorf("fail to generate sub: %w", err)
	}

	g := genny.New()
	g.RunFn(genesisModify(opts))
	g.RunFn(genesisTypesModify(opts))
	g.RunFn(genesisProtoModify(opts))
	g.RunFn(keysModify(replacer, opts))

	if err := g.OnlyFS(subFs, nil, nil); err != nil {
		return g, errors.Errorf("generator fs: %w", err)
	}

	appModulePath := gomodulepath.ExtractAppPath(opts.ModulePath)

	ctx := plush.NewContext()
	ctx.Set("moduleName", opts.ModuleName)
	ctx.Set("modulePath", opts.ModulePath)
	ctx.Set("appName", opts.AppName)
	ctx.Set("protoVer", opts.ProtoVer)
	ctx.Set("ibcOrdering", opts.IBCOrdering)
	ctx.Set("dependencies", opts.Dependencies)
	ctx.Set("protoPkgName", module.ProtoPackageName(appModulePath, opts.ModuleName, opts.ProtoVer))

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))
	g.Transformer(genny.Replace("{{protoDir}}", opts.ProtoDir))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	g.Transformer(genny.Replace("{{moduleName}}", opts.ModuleName))
	g.Transformer(genny.Replace("{{protoVer}}", opts.ProtoVer))

	return g, nil
}

func genesisModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "keeper/genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Genesis init
		replacementModuleInit := `if err := k.Port.Set(ctx, genState.PortId); err != nil {
		return err
	}`
		content, err := xast.ModifyFunction(
			f.String(),
			"InitGenesis",
			xast.AppendFuncCode(replacementModuleInit),
		)
		if err != nil {
			return err
		}

		// Genesis export
		replacementModuleExport := `genesis.PortId, err = k.Port.Get(ctx)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, err
	}`
		content, err = xast.ModifyFunction(
			content,
			"ExportGenesis",
			xast.AppendFuncCode(replacementModuleExport),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func genesisTypesModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "types/genesis.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		content, err := xast.AppendImports(
			f.String(),
			xast.WithNamedImport("host", "github.com/cosmos/ibc-go/v10/modules/core/24-host"),
		)
		if err != nil {
			return err
		}

		// Default genesis
		content, err = xast.ModifyFunction(
			content,
			"DefaultGenesis",
			xast.AppendFuncStruct("GenesisState", "PortId", "PortID"),
		)
		if err != nil {
			return err
		}

		// Validate genesis
		// PlaceholderIBCGenesisTypeValidate
		replacementTypesValidate := `if err := host.PortIdentifierValidator(gs.PortId); err != nil {
	return err
}`
		content, err = xast.ModifyFunction(
			content,
			"Validate",
			xast.AppendFuncCode(replacementTypesValidate),
		)
		if err != nil {
			return err
		}

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// Modifies genesis.proto to add a new field.
//
// What it depends on:
//   - Existence of a message named 'GenesisState' in genesis.proto.
func genesisProtoModify(opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := opts.ProtoFile("genesis.proto")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		protoFile, err := protoutil.ParseProtoFile(f)
		if err != nil {
			return err
		}

		// Grab GenesisState and add next (always 2, I gather) available field.
		// TODO: typed.ProtoGenesisStateMessage exists but in subfolder, so we can't use it here, refactor?
		genesisState, err := protoutil.GetMessageByName(protoFile, "GenesisState")
		if err != nil {
			return errors.Errorf("couldn't find message 'GenesisState' in %s: %w", path, err)
		}
		field := protoutil.NewField("port_id", "string", protoutil.NextUniqueID(genesisState))
		protoutil.Append(genesisState, field)

		newFile := genny.NewFileS(path, protoutil.Print(protoFile))
		return r.File(newFile)
	}
}

func keysModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join("x", opts.ModuleName, "types/keys.go")
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
	PortKey = collections.NewPrefix("%[1]v-port-")
)`
		replacementPort := fmt.Sprintf(templatePort, opts.ModuleName)
		content = replacer.Replace(content, module.PlaceholderIBCKeysPort, replacementPort)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func appIBCModify(replacer placeholder.Replacer, opts *CreateOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := module.PathIBCConfigGo
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// Import
		content, err := xast.AppendImports(
			f.String(),
			xast.WithNamedImport(
				fmt.Sprintf("%[1]vmodule", opts.ModuleName),
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

		// create IBC module
		templateIBCModule := `%[2]vIBCModule := %[2]vmodule.NewIBCModule(app.appCodec, app.%[3]vKeeper)
		ibcRouter.AddRoute(%[2]vmoduletypes.ModuleName, %[2]vIBCModule)
%[1]v`
		replacementIBCModule := fmt.Sprintf(
			templateIBCModule,
			module.PlaceholderIBCNewModule,
			opts.ModuleName,
			xstrings.Title(opts.ModuleName),
		)
		content = replacer.Replace(content, module.PlaceholderIBCNewModule, replacementIBCModule)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
