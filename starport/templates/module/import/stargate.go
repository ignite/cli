package moduleimport

import (
	"context"
	"embed"
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/pkg/xgenny"
	"github.com/tendermint/starport/starport/templates/module"
	"github.com/tendermint/starport/starport/templates/testutil"
)

var (
	//go:embed stargate/* stargate/**/*
	fsStargate embed.FS

	stargateTemplate = xgenny.NewEmbedWalker(fsStargate, "stargate/")
)

// NewStargate returns the generator to scaffold code to import wasm module inside a Stargate app
func NewStargate(ctx context.Context, opts *ImportOptions) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(appModifyStargate(ctx))
	g.RunFn(rootModifyStargate(ctx, opts))
	if err := g.Box(stargateTemplate); err != nil {
		return g, err
	}
	pctx := plush.NewContext()
	pctx.Set("AppName", opts.AppName)
	pctx.Set("title", strings.Title)

	testutil.WASMRegister(ctx, pctx, g)

	g.Transformer(plushgen.Transformer(pctx))
	g.Transformer(genny.Replace("{{binaryNamePrefix}}", opts.BinaryNamePrefix))
	return g, nil
}

// app.go modification on Stargate when importing wasm
func appModifyStargate(ctx context.Context) genny.RunFn {
	return func(r *genny.Runner) error {
		path := module.PathAppGo
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateImport := `%[1]v
		"strings"
		"github.com/CosmWasm/wasmd/x/wasm"
		wasmclient "github.com/CosmWasm/wasmd/x/wasm/client"`
		replacementImport := fmt.Sprintf(templateImport, module.PlaceholderSgAppModuleImport)
		content := placeholder.Replace(ctx, f.String(), module.PlaceholderSgAppModuleImport, replacementImport)

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
		content = placeholder.Replace(ctx, content, module.PlaceholderSgWasmAppEnabledProposals, templateEnabledProposals)

		templateGovProposalHandlers := `%[1]v
		govProposalHandlers = wasmclient.ProposalHandlers`
		replacementProposalHandlers := fmt.Sprintf(templateGovProposalHandlers, module.PlaceholderSgAppGovProposalHandlers)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgAppGovProposalHandlers, replacementProposalHandlers)

		templateModuleBasic := `%[1]v
		wasm.AppModuleBasic{},`
		replacementModuleBasic := fmt.Sprintf(templateModuleBasic, module.PlaceholderSgAppModuleBasic)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgAppModuleBasic, replacementModuleBasic)

		templateKeeperDeclaration := `%[1]v
		wasmKeeper       wasm.Keeper
		scopedWasmKeeper capabilitykeeper.ScopedKeeper
		`
		replacementKeeperDeclaration := fmt.Sprintf(templateKeeperDeclaration, module.PlaceholderSgAppKeeperDeclaration)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgAppKeeperDeclaration, replacementKeeperDeclaration)

		templateDeclaration := `%[1]v
		scopedWasmKeeper := app.CapabilityKeeper.ScopeToModule(wasm.ModuleName)
		`
		replacementDeclaration := fmt.Sprintf(templateDeclaration, module.PlaceholderSgAppScopedKeeper)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgAppScopedKeeper, replacementDeclaration)

		templateDeclaration = `%[1]v
		app.scopedWasmKeeper = scopedWasmKeeper
		`
		replacementDeclaration = fmt.Sprintf(templateDeclaration, module.PlaceholderSgAppBeforeInitReturn)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgAppBeforeInitReturn, replacementDeclaration)

		templateEnabledProposalsArgument := `%[1]v
		enabledProposals []wasm.ProposalType, wasmOpts []wasm.Option,`
		replacementEnabledProposalsArgument := fmt.Sprintf(templateEnabledProposalsArgument, module.PlaceholderSgAppNewArgument)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgAppNewArgument, replacementEnabledProposalsArgument)

		templateStoreKey := `%[1]v
		wasm.StoreKey,`
		replacementStoreKey := fmt.Sprintf(templateStoreKey, module.PlaceholderSgAppStoreKey)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgAppStoreKey, replacementStoreKey)

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
				wasmOpts...,
		)
	
		// The gov proposal types can be individually enabled
		if len(enabledProposals) != 0 {
			govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.wasmKeeper, enabledProposals))
		}`
		replacementKeeperDefinition := fmt.Sprintf(templateKeeperDefinition, module.PlaceholderSgAppKeeperDefinition)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgAppKeeperDefinition, replacementKeeperDefinition)

		templateAppModule := `%[1]v
		wasm.NewAppModule(appCodec, &app.wasmKeeper, app.StakingKeeper),`
		replacementAppModule := fmt.Sprintf(templateAppModule, module.PlaceholderSgAppAppModule)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgAppAppModule, replacementAppModule)

		templateInitGenesis := `%[1]v
		wasm.ModuleName,`
		replacementInitGenesis := fmt.Sprintf(templateInitGenesis, module.PlaceholderSgAppInitGenesis)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgAppInitGenesis, replacementInitGenesis)

		templateParamSubspace := `%[1]v
		paramsKeeper.Subspace(wasm.ModuleName)`
		replacementParamSubspace := fmt.Sprintf(templateParamSubspace, module.PlaceholderSgAppParamSubspace)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgAppParamSubspace, replacementParamSubspace)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// app.go modification on Stargate when importing wasm
func rootModifyStargate(ctx context.Context, opts *ImportOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := "cmd/" + opts.BinaryNamePrefix + "d/cmd/root.go"
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		templateImport := `%[1]v
		"github.com/CosmWasm/wasmd/x/wasm"
		"github.com/prometheus/client_golang/prometheus"
		wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"`
		replacementImport := fmt.Sprintf(templateImport, module.PlaceholderSgRootImport)
		content := placeholder.Replace(ctx, f.String(), module.PlaceholderSgRootImport, replacementImport)

		templateCommand := `%[1]v
		AddGenesisWasmMsgCmd(app.DefaultNodeHome),`
		replacementCommand := fmt.Sprintf(templateCommand, module.PlaceholderSgRootCommands)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgRootCommands, replacementCommand)

		templateInitFlags := `%[1]v
		wasm.AddModuleInitFlags(startCmd)`
		replacementInitFlags := fmt.Sprintf(templateInitFlags, module.PlaceholderSgRootInitFlags)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgRootInitFlags, replacementInitFlags)

		template := `%[1]v
		var wasmOpts []wasm.Option
		if cast.ToBool(appOpts.Get("telemetry.enabled")) {
			   wasmOpts = append(wasmOpts, wasmkeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
		}`
		replacement := fmt.Sprintf(template, module.PlaceholderSgRootAppBeforeInit)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgRootAppBeforeInit, replacement)

		template = `%[1]v
		app.GetEnabledProposals(),
		wasmOpts,`
		replacement = fmt.Sprintf(template, module.PlaceholderSgRootAppArgument)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgRootAppArgument, replacement)

		template = `%[1]v
		app.GetEnabledProposals(),
		nil,`
		replacement = fmt.Sprintf(template, module.PlaceholderSgRootExportArgument)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgRootExportArgument, replacement)

		replacement = fmt.Sprintf(template, module.PlaceholderSgRootNoHeightExportArgument)
		content = placeholder.Replace(ctx, content, module.PlaceholderSgRootNoHeightExportArgument, replacement)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
