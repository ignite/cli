package module

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
)

// New ...
func NewImportStargate(opts *ImportOptions) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(importAppModifyStargate(opts))
	//g.RunFn(exportModify(opts))
	//g.RunFn(cmdMainModify(opts))
	if err := g.Box(packr.New("wasm", "./wasm")); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("AppName", opts.AppName)
	ctx.Set("title", strings.Title)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", opts.AppName))
	return g, nil
}

// app.go modification on Stargate when importing wasm
func importAppModifyStargate(opts *ImportOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := PathAppGo
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateImport := `%[1]v
		"github.com/CosmWasm/wasmd/x/wasm"
		wasmclient "github.com/CosmWasm/wasmd/x/wasm/client"`
		replacementImport := fmt.Sprintf(templateImport, placeholderSgAppModuleImport)
		content := strings.Replace(f.String(), placeholderSgAppModuleImport, replacementImport, 1)

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
		content = strings.Replace(f.String(), placeholderSgWasmAppEnabledProposals, templateEnabledProposals, 1)

		templateGovProposalHandler := `%[1]v
		wasmclient.ProposalHandlers,`
		replacementProposalHandler := fmt.Sprintf(templateGovProposalHandler, placeholderSgAppGovProposalHandler)
		content = strings.Replace(f.String(), placeholderSgAppGovProposalHandler, replacementProposalHandler, 1)

		templateModuleBasic := `%[1]v
		wasm.AppModuleBasic{},`
		replacementModuleBasic := fmt.Sprintf(templateModuleBasic, placeholderSgAppModuleBasic)
		content = strings.Replace(f.String(), placeholderSgAppModuleBasic, replacementModuleBasic, 1)

		templateKeeperDeclaration := `%[1]v
		wasmKeeper wasm.Keeper`
		replacementKeeperDeclaration := fmt.Sprintf(templateKeeperDeclaration, placeholderSgAppKeeperDeclaration)
		content = strings.Replace(f.String(), placeholderSgAppKeeperDeclaration, replacementKeeperDeclaration, 1)

		templateEnabledProposalsArgument := `%[1]v
		enabledProposals []wasm.ProposalType,`
		replacementEnabledProposalsArgument := fmt.Sprintf(templateEnabledProposalsArgument, placeholderSgAppNewArgument)
		content = strings.Replace(f.String(), placeholderSgAppNewArgument, replacementEnabledProposalsArgument, 1)

		templateStoreKey := `%[1]v
		wasm.StoreKey,`
		replacementStoreKey := fmt.Sprintf(templateStoreKey, placeholderSgAppStoreKey)
		content = strings.Replace(f.String(), placeholderSgAppStoreKey, replacementStoreKey, 1)

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
		replacementKeeperDefinition := fmt.Sprintf(templateKeeperDefinition, placeholderSgAppKeeperDefinition)
		content = strings.Replace(f.String(), placeholderSgAppKeeperDefinition, replacementKeeperDefinition, 1)

		templateAppModule := `%[1]v
		wasm.NewAppModule(&app.wasmKeeper, app.stakingKeeper),`
		replacementAppModule := fmt.Sprintf(templateAppModule, placeholderSgAppAppModule)
		content = strings.Replace(f.String(), placeholderSgAppAppModule, replacementAppModule, 1)

		templateInitGenesis := `%[1]v
		wasm.ModuleName,`
		replacementInitGenesis := fmt.Sprintf(templateInitGenesis, placeholderSgAppInitGenesis)
		content = strings.Replace(f.String(), placeholderSgAppInitGenesis, replacementInitGenesis, 1)

		templateParamSubspace := `%[1]v
		paramsKeeper.Subspace(wasm.ModuleName)`
		replacementParamSubspace := fmt.Sprintf(templateParamSubspace, placeholderSgAppParamSubspace)
		content = strings.Replace(f.String(), placeholderSgAppParamSubspace, replacementParamSubspace, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// app.go modification on Stargate when importing wasm
func importRootModifyStargate(opts *ImportOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := "cmd/" + opts.BinaryNamePrefix + "d/cmd/root.go"
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateImport := `%[1]v
		"github.com/CosmWasm/wasmd/x/wasm"`
		replacementImport := fmt.Sprintf(templateImport, placeholderSgRootImport)
		content := strings.Replace(f.String(), placeholderSgRootImport, replacementImport, 1)

		templateCommand := `%[1]v
		"AddGenesisWasmMsgCmd(app.DefaultNodeHome),`
		replacementCommand := fmt.Sprintf(templateCommand, placeholderSgRootCommands)
		content = strings.Replace(f.String(), placeholderSgRootCommands, replacementCommand, 1)

		templateInitFlags := `%[1]v
		"wasm.AddModuleInitFlags(startCmd)`
		replacementInitFlags := fmt.Sprintf(templateInitFlags, placeholderSgRootInitFlags)
		content = strings.Replace(f.String(), placeholderSgRootInitFlags, replacementInitFlags, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}