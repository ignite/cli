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
func NewImport(opts *ImportOptions) (*genny.Generator, error) {
	g := genny.New()
	g.RunFn(appModify(opts))
	g.RunFn(exportModify(opts))
	g.RunFn(cmdMainModify(opts))
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

func appModify(opts *ImportOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := "app/app.go"
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
	"path/filepath"
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/spf13/viper"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"`
		replacement := fmt.Sprintf(template, placeholder)
		content := strings.Replace(f.String(), placeholder, replacement, 1)

		template2 := `%[1]v
		distr.AppModuleBasic{},
		wasm.AppModuleBasic{},`
		replacement2 := fmt.Sprintf(template2, placeholder2)
		content = strings.Replace(content, placeholder2, replacement2, 1)

		template2_1 := `%[1]v
		distr.ModuleName: nil,`
		replacement2_1 := fmt.Sprintf(template2_1, placeholder2_1)
		content = strings.Replace(content, placeholder2_1, replacement2_1, 1)

		template3 := `%[1]v
		distrKeeper    distr.Keeper
		wasmKeeper    wasm.Keeper`
		replacement3 := fmt.Sprintf(template3, placeholder3)
		content = strings.Replace(content, placeholder3, replacement3, 1)

		template5 := `%[1]v
		distr.StoreKey,
		wasm.StoreKey,`
		replacement5 := fmt.Sprintf(template5, placeholder5)
		content = strings.Replace(content, placeholder5, replacement5, 1)

		template5_1 := `%[1]v
		app.subspaces[distr.ModuleName] = app.paramsKeeper.Subspace(distr.DefaultParamspace)`
		replacement5_1 := fmt.Sprintf(template5_1, placeholder5_1)
		content = strings.Replace(content, placeholder5_1, replacement5_1, 1)

		template5_2 := `%[1]v
		app.distrKeeper = distr.NewKeeper(
			app.cdc, keys[distr.StoreKey], app.subspaces[distr.ModuleName], &stakingKeeper,
			app.supplyKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs(),
		)`
		replacement5_2 := fmt.Sprintf(template5_2, placeholder5_2)
		content = strings.Replace(content, placeholder5_2, replacement5_2, 1)

		template5_3 := `%[1]v
		app.distrKeeper.Hooks(),`
		replacement5_3 := fmt.Sprintf(template5_3, placeholder5_3)
		content = strings.Replace(content, placeholder5_3, replacement5_3, 1)

		template4 := placeholder4 + "\n" +
			"type WasmWrapper struct { Wasm wasm.Config `mapstructure:\"wasm\"`}" + `
		var wasmRouter = bApp.Router()
		homeDir := viper.GetString(cli.HomeFlag)
		wasmDir := filepath.Join(homeDir, "wasm")

		wasmWrap := WasmWrapper{Wasm: wasm.DefaultWasmConfig()}
		err := viper.Unmarshal(&wasmWrap)
		if err != nil {
			panic("error while reading wasm config: " + err.Error())
		}
		wasmConfig := wasmWrap.Wasm
		supportedFeatures := "staking"
		app.subspaces[wasm.ModuleName] = app.paramsKeeper.Subspace(wasm.DefaultParamspace)
		app.wasmKeeper = wasm.NewKeeper(app.cdc, keys[wasm.StoreKey], app.subspaces[wasm.ModuleName], app.accountKeeper, app.bankKeeper, app.stakingKeeper, app.distrKeeper, wasmRouter, wasmDir, wasmConfig, supportedFeatures, nil, nil)`
		content = strings.Replace(content, placeholder4, template4, 1)

		template6 := `%[1]v
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
		wasm.NewAppModule(app.wasmKeeper),`
		replacement6 := fmt.Sprintf(template6, placeholder6)
		content = strings.Replace(content, placeholder6, replacement6, 1)

		template6_1 := `%[1]v
		distr.ModuleName,`
		replacement6_1 := fmt.Sprintf(template6_1, placeholder6_1)
		content = strings.Replace(content, placeholder6_1, replacement6_1, 1)

		template6_2 := `%[1]v
		distr.ModuleName,`
		replacement6_2 := fmt.Sprintf(template6_2, placeholder6_2)
		content = strings.Replace(content, placeholder6_2, replacement6_2, 1)

		template7 := `%[1]v
		wasm.ModuleName,`
		replacement7 := fmt.Sprintf(template7, placeholder7)
		content = strings.Replace(content, placeholder7, replacement7, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

// Append Distr modules in export.go
func exportModify(opts *ImportOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := "app/export.go"
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}

		template := `%[1]v
		/* Handle fee distribution state. */

		// withdraw all validator commission
		app.stakingKeeper.IterateValidators(ctx, func(_ int64, val staking.ValidatorI) (stop bool) {
			_, err := app.distrKeeper.WithdrawValidatorCommission(ctx, val.GetOperator())
			if err != nil {
				log.Fatal(err)
			}
			return false
		})

		// withdraw all delegator rewards
		dels := app.stakingKeeper.GetAllDelegations(ctx)
		for _, delegation := range dels {
			_, err := app.distrKeeper.WithdrawDelegationRewards(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
			if err != nil {
				log.Fatal(err)
			}
		}

		// clear validator slash events
		app.distrKeeper.DeleteAllValidatorSlashEvents(ctx)

		// clear validator historical rewards
		app.distrKeeper.DeleteAllValidatorHistoricalRewards(ctx)`
		replacement := fmt.Sprintf(template, placeholder)
		content := strings.Replace(f.String(), placeholder, replacement, 1)

		template2 := `%[1]v
		// donate any unwithdrawn outstanding reward fraction tokens to the community pool
		scraps := app.distrKeeper.GetValidatorOutstandingRewards(ctx, val.GetOperator())
		feePool := app.distrKeeper.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
		app.distrKeeper.SetFeePool(ctx, feePool)

		app.distrKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())`
		replacement2 := fmt.Sprintf(template2, placeholder2)
		content = strings.Replace(content, placeholder2, replacement2, 1)

		template3 := `%[1]v
		// reinitialize all delegations
		for _, del := range dels {
			app.distrKeeper.Hooks().BeforeDelegationCreated(ctx, del.DelegatorAddress, del.ValidatorAddress)
			app.distrKeeper.Hooks().AfterDelegationModified(ctx, del.DelegatorAddress, del.ValidatorAddress)
		}`
		replacement3 := fmt.Sprintf(template3, placeholder3)
		content = strings.Replace(content, placeholder3, replacement3, 1)

		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func cmdMainModify(opts *ImportOptions) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("cmd/%[1]vcli/main.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		template := `%[1]v
	wasmrest "github.com/CosmWasm/wasmd/x/wasm/client/rest"`
		replacement := fmt.Sprintf(template, placeholder)
		content := strings.Replace(f.String(), placeholder, replacement, 1)

		template2 := `%[1]v
	wasmrest.RegisterRoutes(rs.CliCtx, rs.Mux)`
		replacement2 := fmt.Sprintf(template2, placeholder2)
		content = strings.Replace(content, placeholder2, replacement2, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
