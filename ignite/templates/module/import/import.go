package moduleimport

import (
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/templates/field/plushhelpers"
	"github.com/ignite/cli/ignite/templates/module"
)

// NewGenerator returns the generator to scaffold code to import wasm module inside an app.
func NewGenerator(replacer placeholder.Replacer, opts *ImportOptions) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(appModify(replacer, opts))
	g.RunFn(cmdModify(replacer, opts))

	ctx := plush.NewContext()
	ctx.Set("AppName", opts.AppName)
	plushhelpers.ExtendPlushContext(ctx)

	return g, nil
}

// app.go modification when importing wasm.
func appModify(replacer placeholder.Replacer, opts *ImportOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, module.PathAppGo)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateImport := `%[1]v
		"github.com/tendermint/spm-extras/wasmcmd"
		"github.com/CosmWasm/wasmd/x/wasm"
		wasmclient "github.com/CosmWasm/wasmd/x/wasm/client"`
		replacementImport := fmt.Sprintf(templateImport, module.PlaceholderSgAppModuleImport)
		content := replacer.Replace(f.String(), module.PlaceholderSgAppModuleImport, replacementImport)

		templateEnabledProposals := `var (
			// If EnabledSpecificProposals is "", and this is "true", then enable all x/wasm proposals.
			// If EnabledSpecificProposals is "", and this is not "true", then disable all x/wasm proposals.
			ProposalsEnabled = "false"
			// If set to non-empty string it must be comma-separated list of values that are all a subset
			// of "EnableAllProposals" (takes precedence over ProposalsEnabled)
			// https://github.com/CosmWasm/wasmd/blob/02a54d33ff2c064f3539ae12d75d027d9c665f05/x/wasm/internal/types/proposal.go#L28-L34
			EnableSpecificProposals = ""
		)
		`
		content = replacer.Replace(content, module.PlaceholderSgWasmAppEnabledProposals, templateEnabledProposals)

		templateGovProposalHandlers := `%[1]v
		govProposalHandlers = wasmclient.ProposalHandlers`
		replacementProposalHandlers := fmt.Sprintf(templateGovProposalHandlers, module.PlaceholderSgAppGovProposalHandlers)
		content = replacer.Replace(content, module.PlaceholderSgAppGovProposalHandlers, replacementProposalHandlers)

		templateModuleBasic := `%[1]v
		wasm.AppModuleBasic{},`
		replacementModuleBasic := fmt.Sprintf(templateModuleBasic, module.PlaceholderSgAppModuleBasic)
		content = replacer.Replace(content, module.PlaceholderSgAppModuleBasic, replacementModuleBasic)

		templateKeeperDeclaration := `%[1]v
		wasmKeeper       wasm.Keeper
		scopedWasmKeeper capabilitykeeper.ScopedKeeper
		`
		replacementKeeperDeclaration := fmt.Sprintf(templateKeeperDeclaration, module.PlaceholderSgAppKeeperDeclaration)
		content = replacer.Replace(content, module.PlaceholderSgAppKeeperDeclaration, replacementKeeperDeclaration)

		templateDeclaration := `%[1]v
		scopedWasmKeeper := app.CapabilityKeeper.ScopeToModule(wasm.ModuleName)
		`
		replacementDeclaration := fmt.Sprintf(templateDeclaration, module.PlaceholderSgAppScopedKeeper)
		content = replacer.Replace(content, module.PlaceholderSgAppScopedKeeper, replacementDeclaration)

		templateDeclaration = `%[1]v
		app.scopedWasmKeeper = scopedWasmKeeper
		`
		replacementDeclaration = fmt.Sprintf(templateDeclaration, module.PlaceholderSgAppBeforeInitReturn)
		content = replacer.Replace(content, module.PlaceholderSgAppBeforeInitReturn, replacementDeclaration)

		templateStoreKey := `%[1]v
		wasm.StoreKey,`
		replacementStoreKey := fmt.Sprintf(templateStoreKey, module.PlaceholderSgAppStoreKey)
		content = replacer.Replace(content, module.PlaceholderSgAppStoreKey, replacementStoreKey)

		templateKeeperDefinition := `%[1]v
		wasmDir := filepath.Join(homePath, "wasm")
	
		wasmConfig, err := wasm.ReadWasmConfig(appOpts)
		if err != nil {
			panic("error while reading wasm config: " + err.Error())
		}

		// The last arguments can contain custom message handlers, and custom query handlers,
		// if we want to allow any custom callbacks
		supportedFeatures := "staking"
		app.wasmKeeper = wasm.NewKeeper(
				appCodec,
				keys[wasm.StoreKey],
				app.GetSubspace(wasm.ModuleName),
				app.AccountKeeper,
				app.BankKeeper,
				app.StakingKeeper,
				app.DistrKeeper,
				app.IBCKeeper.ChannelKeeper,
				&app.IBCKeeper.PortKeeper,
				scopedWasmKeeper,
				app.TransferKeeper,
				app.Router(),
				app.GRPCQueryRouter(),
				wasmDir,
				wasmConfig,
				supportedFeatures,
		)
	
		// The gov proposal types can be individually enabled
		enabledProposals := wasmcmd.GetEnabledProposals(ProposalsEnabled, EnableSpecificProposals)
		if len(enabledProposals) != 0 {
			govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.wasmKeeper, enabledProposals))
		}`
		replacementKeeperDefinition := fmt.Sprintf(templateKeeperDefinition, module.PlaceholderSgAppKeeperDefinition)
		content = replacer.Replace(content, module.PlaceholderSgAppKeeperDefinition, replacementKeeperDefinition)

		templateAppModule := `%[1]v
		wasm.NewAppModule(appCodec, &app.wasmKeeper, app.StakingKeeper),`
		replacementAppModule := fmt.Sprintf(templateAppModule, module.PlaceholderSgAppAppModule)
		content = replacer.Replace(content, module.PlaceholderSgAppAppModule, replacementAppModule)

		templateInitGenesis := `%[1]v
		wasm.ModuleName,`
		replacementInitGenesis := fmt.Sprintf(templateInitGenesis, module.PlaceholderSgAppInitGenesis)
		content = replacer.Replace(content, module.PlaceholderSgAppInitGenesis, replacementInitGenesis)

		templateParamSubspace := `%[1]v
		paramsKeeper.Subspace(wasm.ModuleName)`
		replacementParamSubspace := fmt.Sprintf(templateParamSubspace, module.PlaceholderSgAppParamSubspace)
		content = replacer.Replace(content, module.PlaceholderSgAppParamSubspace, replacementParamSubspace)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// main.go modification when importing wasm.
func cmdModify(replacer placeholder.Replacer, opts *ImportOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := filepath.Join(opts.AppPath, "cmd", opts.BinaryNamePrefix+"d/cmd/root.go")
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		// add wasm import
		templateImport := `%[1]v
		"github.com/tendermint/spm-extras/wasmcmd"`
		replacementImport := fmt.Sprintf(templateImport, module.PlaceholderSgRootModuleImport)
		content := replacer.Replace(f.String(), module.PlaceholderSgRootModuleImport, replacementImport)

		// add wasm command
		templateCommands := `wasmcmd.GenesisWasmMsgCmd(app.DefaultNodeHome),
		%[1]v`
		replacementCommands := fmt.Sprintf(templateCommands, module.PlaceholderSgRootCommands)
		content = replacer.Replace(content, module.PlaceholderSgRootCommands, replacementCommands)

		// add wasm start args
		templateArgs := `wasmcmd.AddModuleInitFlags(startCmd)
		%[1]v`
		replacementArgs := fmt.Sprintf(templateArgs, module.PlaceholderSgRootArgument)
		content = replacer.Replace(content, module.PlaceholderSgRootArgument, replacementArgs)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
