package module_import

import (
	"fmt"
	"github.com/tendermint/starport/starport/templates/module"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
)


// New ...
func NewImportStargate(opts *ImportOptions) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(appModifyStargate(opts))
	g.RunFn(rootModifyStargate(opts))
	//g.RunFn(cmdMainModify(opts))
	if err := g.Box(packr.New("module/import/templates/stargate", "./stargate")); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("AppName", opts.AppName)
	ctx.Set("title", strings.Title)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{binaryNamePrefix}}", opts.BinaryNamePrefix))
	return g, nil
}

// app.go modification on Stargate when importing wasm
func appModifyStargate(opts *ImportOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := module.PathAppGo
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateImport := `%[1]v
		"github.com/CosmWasm/wasmd/x/wasm"
		wasmclient "github.com/CosmWasm/wasmd/x/wasm/client"`
		replacementImport := fmt.Sprintf(templateImport, module.PlaceholderSgAppModuleImport)
		content := strings.Replace(f.String(), module.PlaceholderSgAppModuleImport, replacementImport, 1)

		templateEnabledProposals := `var (
			// If EnabledSpecificProposals is "", and this is "true", then enable all x/wasm proposals.
			// If EnabledSpecificProposals is "", and this is not "true", then disable all x/wasm proposals.
			ProposalsEnabled = "false"
			// If set to non-empty string it must be comma-separated list of values that are all a subset
			// of "EnableAllProposals" (takes precedence over ProposalsEnabled)
			// https://github.com/CosmWasm/wasmd/blob/02a54d33ff2c064f3539ae12d75d027d9c665f05/x/wasm/internal/types/proposal.go#L28-L34
			EnableSpecificProposals = ""
		)
		
		// GetEnabledProposals parses the ProposalsEnabled / EnableSpecificProposals values to
		// produce a list of enabled proposals to pass into wasmd app.
		func GetEnabledProposals() []wasm.ProposalType {
			if EnableSpecificProposals == "" {
				if ProposalsEnabled == "true" {
					return wasm.EnableAllProposals
				}
				return wasm.DisableAllProposals
			}
			chunks := strings.Split(EnableSpecificProposals, ",")
			proposals, err := wasm.ConvertToProposals(chunks)
			if err != nil {
				panic(err)
			}
			return proposals
		}`
		content = strings.Replace(f.String(), module.PlaceholderSgWasmAppEnabledProposals, templateEnabledProposals, 1)

		templateGovProposalHandler := `%[1]v
		wasmclient.ProposalHandlers,`
		replacementProposalHandler := fmt.Sprintf(templateGovProposalHandler, module.PlaceholderSgAppGovProposalHandler)
		content = strings.Replace(f.String(), module.PlaceholderSgAppGovProposalHandler, replacementProposalHandler, 1)

		templateModuleBasic := `%[1]v
		wasm.AppModuleBasic{},`
		replacementModuleBasic := fmt.Sprintf(templateModuleBasic, module.PlaceholderSgAppModuleBasic)
		content = strings.Replace(f.String(), module.PlaceholderSgAppModuleBasic, replacementModuleBasic, 1)

		templateKeeperDeclaration := `%[1]v
		wasmKeeper wasm.Keeper`
		replacementKeeperDeclaration := fmt.Sprintf(templateKeeperDeclaration, module.PlaceholderSgAppKeeperDeclaration)
		content = strings.Replace(f.String(), module.PlaceholderSgAppKeeperDeclaration, replacementKeeperDeclaration, 1)

		templateEnabledProposalsArgument := `%[1]v
		enabledProposals []wasm.ProposalType,`
		replacementEnabledProposalsArgument := fmt.Sprintf(templateEnabledProposalsArgument, module.PlaceholderSgAppNewArgument)
		content = strings.Replace(f.String(), module.PlaceholderSgAppNewArgument, replacementEnabledProposalsArgument, 1)

		templateStoreKey := `%[1]v
		wasm.StoreKey,`
		replacementStoreKey := fmt.Sprintf(templateStoreKey, module.PlaceholderSgAppStoreKey)
		content = strings.Replace(f.String(), module.PlaceholderSgAppStoreKey, replacementStoreKey, 1)

		templateKeeperDefinition := `%[1]v
		// The last arguments can contain custom message handlers, and custom query handlers,
		// if we want to allow any custom callbacks
		supportedFeatures := "staking"
		app.wasmKeeper = wasm.NewKeeper(
			appCodec,
			keys[wasm.StoreKey],
			app.getSubspace(wasm.ModuleName),
			app.accountKeeper,
			app.bankKeeper,
			app.stakingKeeper,
			app.distrKeeper,
			wasmRouter,
			wasmDir,
			wasmConfig,
			supportedFeatures,
			nil,
			nil,
		)
	
		// The gov proposal types can be individually enabled
		if len(enabledProposals) != 0 {
			govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.wasmKeeper, enabledProposals))
		}`
		replacementKeeperDefinition := fmt.Sprintf(templateKeeperDefinition, module.PlaceholderSgAppKeeperDefinition)
		content = strings.Replace(f.String(), module.PlaceholderSgAppKeeperDefinition, replacementKeeperDefinition, 1)

		templateAppModule := `%[1]v
		wasm.NewAppModule(&app.wasmKeeper, app.stakingKeeper),`
		replacementAppModule := fmt.Sprintf(templateAppModule, module.PlaceholderSgAppAppModule)
		content = strings.Replace(f.String(), module.PlaceholderSgAppAppModule, replacementAppModule, 1)

		templateInitGenesis := `%[1]v
		wasm.ModuleName,`
		replacementInitGenesis := fmt.Sprintf(templateInitGenesis, module.PlaceholderSgAppInitGenesis)
		content = strings.Replace(f.String(), module.PlaceholderSgAppInitGenesis, replacementInitGenesis, 1)

		templateParamSubspace := `%[1]v
		paramsKeeper.Subspace(wasm.ModuleName)`
		replacementParamSubspace := fmt.Sprintf(templateParamSubspace, module.PlaceholderSgAppParamSubspace)
		content = strings.Replace(f.String(), module.PlaceholderSgAppParamSubspace, replacementParamSubspace, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// app.go modification on Stargate when importing wasm
func rootModifyStargate(opts *ImportOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := "cmd/" + opts.BinaryNamePrefix + "d/cmd/root.go"
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateImport := `%[1]v
		"github.com/CosmWasm/wasmd/x/wasm"`
		replacementImport := fmt.Sprintf(templateImport, module.PlaceholderSgRootImport)
		content := strings.Replace(f.String(), module.PlaceholderSgRootImport, replacementImport, 1)

		templateCommand := `%[1]v
		"AddGenesisWasmMsgCmd(app.DefaultNodeHome),`
		replacementCommand := fmt.Sprintf(templateCommand, module.PlaceholderSgRootCommands)
		content = strings.Replace(f.String(), module.PlaceholderSgRootCommands, replacementCommand, 1)

		templateInitFlags := `%[1]v
		"wasm.AddModuleInitFlags(startCmd)`
		replacementInitFlags := fmt.Sprintf(templateInitFlags, module.PlaceholderSgRootInitFlags)
		content = strings.Replace(f.String(), module.PlaceholderSgRootInitFlags, replacementInitFlags, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}