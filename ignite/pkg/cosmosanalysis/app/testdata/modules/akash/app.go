package app

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"cosmossdk.io/client/v2/autocli"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/simapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v3/modules/core"
	porttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/ovrclk/akash/x/inflation"

	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ovrclk/akash/x/audit"
	"github.com/ovrclk/akash/x/cert"
	escrowkeeper "github.com/ovrclk/akash/x/escrow/keeper"

	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	abci "github.com/tendermint/tendermint/abci/types"
	tmos "github.com/tendermint/tendermint/libs/os"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/version"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	ica "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v3/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v3/modules/core/03-connection/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"

	dkeeper "github.com/ovrclk/akash/x/deployment/keeper"
	mkeeper "github.com/ovrclk/akash/x/market/keeper"
	pkeeper "github.com/ovrclk/akash/x/provider/keeper"

	icaauth "github.com/ovrclk/akash/x/icaauth"
	icaauthkeeper "github.com/ovrclk/akash/x/icaauth/keeper"
	icaauthtypes "github.com/ovrclk/akash/x/icaauth/types"

	// unnamed import of statik for swagger UI support
	_ "github.com/ovrclk/akash/client/docs/statik"
)

const (
	AppName = "akash"
)

var (
	DefaultHome                         = os.ExpandEnv("$HOME/.akash")
	_           servertypes.Application = (*AkashApp)(nil)

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{}
)

// AkashApp extends ABCI appplication
type AkashApp struct {
	*bam.BaseApp
	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry codectypes.InterfaceRegistry

	invCheckPeriod uint

	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memkeys map[string]*sdk.MemoryStoreKey

	keeper struct {
		acct                authkeeper.AccountKeeper
		authz               authzkeeper.Keeper
		bank                bankkeeper.Keeper
		cap                 *capabilitykeeper.Keeper
		staking             stakingkeeper.Keeper
		slashing            slashingkeeper.Keeper
		mint                mintkeeper.Keeper
		distr               distrkeeper.Keeper
		gov                 govkeeper.Keeper
		crisis              crisiskeeper.Keeper
		upgrade             upgradekeeper.Keeper
		params              paramskeeper.Keeper
		ibc                 *ibckeeper.Keeper
		evidence            evidencekeeper.Keeper
		transfer            ibctransferkeeper.Keeper
		icaHostKeeper       icahostkeeper.Keeper
		icaControllerKeeper icacontrollerkeeper.Keeper
		icaAuthKeeper       icaauthkeeper.Keeper

		// make scoped keepers public for test purposes
		scopedIBCKeeper           capabilitykeeper.ScopedKeeper
		scopedTransferKeeper      capabilitykeeper.ScopedKeeper
		scopedICAControllerKeeper capabilitykeeper.ScopedKeeper
		scopedICAHostKeeper       capabilitykeeper.ScopedKeeper
		scopedIcaAuthKeeper       capabilitykeeper.ScopedKeeper

		// akash keepers
		escrow     escrowkeeper.Keeper
		deployment dkeeper.IKeeper
		market     mkeeper.IKeeper
		provider   pkeeper.IKeeper
		audit      audit.Keeper
		cert       cert.Keeper
		inflation  inflation.Keeper
	}

	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager

	// module configurator
	configurator module.Configurator
}

// https://github.com/cosmos/sdk-tutorials/blob/c6754a1e313eb1ed973c5c91dcc606f2fd288811/app.go#L73

// NewApp creates and returns a new Akash App.
func NewApp(
	logger log.Logger, db dbm.DB, tio io.Writer, loadLatest bool, invCheckPeriod uint, skipUpgradeHeights map[int64]bool,
	homePath string, appOpts servertypes.AppOptions, options ...func(*bam.BaseApp),
) *AkashApp {
	// find out the genesis time, to be used later in inflation calculation
	// genesisTime := getGenesisTime(appOpts, homePath)

	// TODO: Remove cdc in favor of appCodec once all modules are migrated.
	encodingConfig := MakeEncodingConfig()
	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bapp := bam.NewBaseApp(AppName, logger, db, encodingConfig.TxConfig.TxDecoder(), options...)
	bapp.SetCommitMultiStoreTracer(tio)
	bapp.SetVersion(version.Version)
	bapp.SetInterfaceRegistry(interfaceRegistry)

	keys := kvStoreKeys()
	tkeys := transientStoreKeys()
	memkeys := memStoreKeys()

	app := &AkashApp{
		BaseApp:           bapp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memkeys:           memkeys,
	}
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())

	app.keeper.params = initParamsKeeper(appCodec, cdc, app.keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	bapp.SetParamStore(app.keeper.params.Subspace(bam.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.keeper.cap = capabilitykeeper.NewKeeper(appCodec, app.keys[capabilitytypes.StoreKey], app.memkeys[capabilitytypes.MemStoreKey])

	scopedIBCKeeper := app.keeper.cap.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.keeper.cap.ScopeToModule(ibctransfertypes.ModuleName)
	scopedIcaAuthKeeper := app.keeper.cap.ScopeToModule(icaauthtypes.ModuleName)
	scopedICAControllerKeeper := app.keeper.cap.ScopeToModule(icacontrollertypes.SubModuleName)
	scopedICAHostKeeper := app.keeper.cap.ScopeToModule(icahosttypes.SubModuleName)

	// seal the capability keeper so all persistent capabilities are loaded in-memory and prevent
	// any further modules from creating scoped sub-keepers.
	app.keeper.cap.Seal()

	app.keeper.acct = authkeeper.NewAccountKeeper(
		appCodec,
		app.keys[authtypes.StoreKey],
		app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		MacPerms(),
	)

	// add authz keeper
	app.keeper.authz = authzkeeper.NewKeeper(app.keys[authzkeeper.StoreKey], appCodec, app.MsgServiceRouter())

	app.keeper.bank = bankkeeper.NewBaseKeeper(
		appCodec,
		app.keys[banktypes.StoreKey],
		app.keeper.acct,
		app.GetSubspace(banktypes.ModuleName),
		app.BlockedAddrs(),
	)

	skeeper := stakingkeeper.NewKeeper(
		appCodec,
		app.keys[stakingtypes.StoreKey],
		app.keeper.acct,
		app.keeper.bank,
		app.GetSubspace(stakingtypes.ModuleName),
	)

	app.keeper.mint = mintkeeper.NewKeeper(
		appCodec,
		app.keys[minttypes.StoreKey],
		app.GetSubspace(minttypes.ModuleName),
		&skeeper,
		app.keeper.acct,
		app.keeper.bank,
		authtypes.FeeCollectorName,
	)

	app.keeper.distr = distrkeeper.NewKeeper(
		appCodec,
		app.keys[distrtypes.StoreKey],
		app.GetSubspace(distrtypes.ModuleName),
		app.keeper.acct,
		app.keeper.bank,
		&skeeper,
		authtypes.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)

	app.keeper.slashing = slashingkeeper.NewKeeper(
		appCodec,
		app.keys[slashingtypes.StoreKey],
		&skeeper,
		app.GetSubspace(slashingtypes.ModuleName),
	)

	app.keeper.crisis = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName),
		invCheckPeriod,
		app.keeper.bank,
		authtypes.FeeCollectorName,
	)

	app.keeper.upgrade = upgradekeeper.NewKeeper(skipUpgradeHeights, app.keys[upgradetypes.StoreKey], appCodec, homePath, app.BaseApp)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.keeper.staking = *skeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.keeper.distr.Hooks(),
			app.keeper.slashing.Hooks(),
		),
	)

	// register IBC Keeper
	app.keeper.ibc = ibckeeper.NewKeeper(
		appCodec, app.keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName),
		app.keeper.staking, app.keeper.upgrade, scopedIBCKeeper,
	)

	// register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.keeper.params)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.keeper.distr)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.keeper.upgrade)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.keeper.ibc.ClientKeeper))

	app.keeper.gov = govkeeper.NewKeeper(
		appCodec,
		app.keys[govtypes.StoreKey],
		app.GetSubspace(govtypes.ModuleName),
		app.keeper.acct,
		app.keeper.bank,
		&skeeper,
		govRouter,
	)

	// register Transfer Keepers
	app.keeper.transfer = ibctransferkeeper.NewKeeper(
		appCodec, app.keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.keeper.ibc.ChannelKeeper, app.keeper.ibc.ChannelKeeper, &app.keeper.ibc.PortKeeper,
		app.keeper.acct, app.keeper.bank, scopedTransferKeeper,
	)

	transferModule := transfer.NewAppModule(app.keeper.transfer)
	transferIBCModule := transfer.NewIBCModule(app.keeper.transfer)

	app.keeper.icaControllerKeeper = icacontrollerkeeper.NewKeeper(
		appCodec, keys[icacontrollertypes.StoreKey], app.GetSubspace(icacontrollertypes.SubModuleName),
		app.keeper.ibc.ChannelKeeper, // may be replaced with middleware such as ics29 fee
		app.keeper.ibc.ChannelKeeper, &app.keeper.ibc.PortKeeper,
		scopedICAControllerKeeper, app.MsgServiceRouter(),
	)

	app.keeper.icaHostKeeper = icahostkeeper.NewKeeper(
		appCodec, keys[icahosttypes.StoreKey], app.GetSubspace(icahosttypes.SubModuleName),
		app.keeper.ibc.ChannelKeeper, &app.keeper.ibc.PortKeeper,
		app.keeper.acct, scopedICAHostKeeper, app.MsgServiceRouter(),
	)

	icaModule := ica.NewAppModule(&app.keeper.icaControllerKeeper, &app.keeper.icaHostKeeper)

	app.keeper.icaAuthKeeper = icaauthkeeper.NewKeeper(appCodec, keys[icaauthtypes.StoreKey], app.keeper.icaControllerKeeper, scopedIcaAuthKeeper)
	icaAuthModule := icaauth.NewAppModule(appCodec, app.keeper.icaAuthKeeper)
	icaAuthIBCModule := icaauth.NewIBCModule(app.keeper.icaAuthKeeper)

	icaControllerIBCModule := icacontroller.NewIBCModule(app.keeper.icaControllerKeeper, icaAuthIBCModule)
	icaHostIBCModule := icahost.NewIBCModule(app.keeper.icaHostKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(icacontrollertypes.SubModuleName, icaControllerIBCModule).
		AddRoute(icahosttypes.SubModuleName, icaHostIBCModule).
		AddRoute(ibctransfertypes.ModuleName, transferIBCModule).
		AddRoute(icaauthtypes.ModuleName, icaControllerIBCModule)

	app.keeper.ibc.SetRouter(ibcRouter)

	// create evidence keeper with evidence router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, app.keys[evidencetypes.StoreKey], &app.keeper.staking, app.keeper.slashing,
	)

	// if evidence needs to be handled for the app, set routes in router here and seal
	app.keeper.evidence = *evidenceKeeper

	app.setAkashKeepers()

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	app.mm = module.NewManager(
		append([]module.AppModule{
			genutil.NewAppModule(app.keeper.acct, app.keeper.staking, app.BaseApp.DeliverTx, encodingConfig.TxConfig),
			auth.NewAppModule(appCodec, app.keeper.acct, nil),
			authzmodule.NewAppModule(appCodec, app.keeper.authz, app.keeper.acct, app.keeper.bank, app.interfaceRegistry),
			vesting.NewAppModule(app.keeper.acct, app.keeper.bank),
			bank.NewAppModule(appCodec, app.keeper.bank, app.keeper.acct),
			capability.NewAppModule(appCodec, *app.keeper.cap),
			crisis.NewAppModule(&app.keeper.crisis, skipGenesisInvariants),
			gov.NewAppModule(appCodec, app.keeper.gov, app.keeper.acct, app.keeper.bank),
			mint.NewAppModule(appCodec, app.keeper.mint, app.keeper.acct, nil),
			slashing.NewAppModule(appCodec, app.keeper.slashing, app.keeper.acct, app.keeper.bank, app.keeper.staking),
			distr.NewAppModule(appCodec, app.keeper.distr, app.keeper.acct, app.keeper.bank, app.keeper.staking),
			staking.NewAppModule(appCodec, app.keeper.staking, app.keeper.acct, app.keeper.bank),
			upgrade.NewAppModule(app.keeper.upgrade),
			evidence.NewAppModule(app.keeper.evidence),
			ibc.NewAppModule(app.keeper.ibc),
			params.NewAppModule(app.keeper.params),
			transferModule,
			icaModule,
			icaAuthModule,
		}, app.akashAppModules()...)...,
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
	// NOTE: As of v0.45.0 of cosmos SDK, all modules need to be here.
	app.mm.SetOrderBeginBlockers(
		app.akashBeginBlockModules()...,
	)
	app.mm.SetOrderEndBlockers(
		app.akashEndBlockModules()...,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	//       properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		app.akashInitGenesisOrder()...,
	)

	app.mm.RegisterInvariants(&app.keeper.crisis)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.mm.RegisterServices(app.configurator)

	// add test gRPC service for testing gRPC queries in isolation
	testdata.RegisterQueryServer(app.GRPCQueryRouter(), testdata.QueryImpl{})

	app.sm = module.NewSimulationManager(
		append([]module.AppModuleSimulation{
			auth.NewAppModule(appCodec, app.keeper.acct, authsims.RandomGenesisAccounts),
			authzmodule.NewAppModule(appCodec, app.keeper.authz, app.keeper.acct, app.keeper.bank, app.interfaceRegistry),
			bank.NewAppModule(appCodec, app.keeper.bank, app.keeper.acct),
			capability.NewAppModule(appCodec, *app.keeper.cap),
			gov.NewAppModule(appCodec, app.keeper.gov, app.keeper.acct, app.keeper.bank),
			mint.NewAppModule(appCodec, app.keeper.mint, app.keeper.acct, nil),
			staking.NewAppModule(appCodec, app.keeper.staking, app.keeper.acct, app.keeper.bank),
			distr.NewAppModule(appCodec, app.keeper.distr, app.keeper.acct, app.keeper.bank, app.keeper.staking),
			slashing.NewAppModule(appCodec, app.keeper.slashing, app.keeper.acct, app.keeper.bank, app.keeper.staking),
			params.NewAppModule(app.keeper.params),
			evidence.NewAppModule(app.keeper.evidence),
			ibc.NewAppModule(app.keeper.ibc),
			transferModule,
			NewICAHostSimModule(icaModule, appCodec),
		},
			app.akashSimModules()...,
		)...,
	)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memkeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	handler, err := ante.NewAnteHandler(ante.HandlerOptions{
		AccountKeeper:   app.keeper.acct,
		BankKeeper:      app.keeper.bank,
		SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
	})
	if err != nil {
		panic(err)
	}
	app.SetAnteHandler(handler)

	app.SetEndBlocker(app.EndBlocker)

	// register the upgrade handler
	app.registerUpgradeHandlers(icaModule)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit("app initialization:" + err.Error())
		}
	}

	app.keeper.scopedIBCKeeper = scopedIBCKeeper
	app.keeper.scopedTransferKeeper = scopedTransferKeeper
	app.keeper.scopedICAControllerKeeper = scopedICAControllerKeeper
	app.keeper.scopedICAHostKeeper = scopedICAHostKeeper
	app.keeper.scopedIcaAuthKeeper = scopedIcaAuthKeeper

	return app
}

func (app *AkashApp) registerUpgradeHandlers(icaModule ica.AppModule) {
	app.keeper.upgrade.SetUpgradeHandler("akash_v0.15.0_cosmos_v0.44.x", func(ctx sdk.Context,
		plan upgradetypes.Plan, _ module.VersionMap,
	) (module.VersionMap, error) {
		// set max expected block time parameter. Replace the default with your expected value
		app.keeper.ibc.ConnectionKeeper.SetParams(ctx, ibcconnectiontypes.DefaultParams())

		// 1st-time running in-store migrations, using 1 as fromVersion to
		// avoid running InitGenesis.
		fromVM := map[string]uint64{
			"auth":         1,
			"bank":         1,
			"capability":   1,
			"crisis":       1,
			"distribution": 1,
			"evidence":     1,
			"gov":          1,
			"mint":         1,
			"params":       1,
			"slashing":     1,
			"staking":      1,
			"upgrade":      1,
			"vesting":      1,
			"ibc":          1,
			"genutil":      1,
			"transfer":     1,

			// akash modules
			"audit":      1,
			"cert":       1,
			"deployment": 1,
			"escrow":     1,
			"market":     1,
			"provider":   1,
		}

		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	})

	upgradeInfo, err := app.keeper.upgrade.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == "akash_v0.15.0_cosmos_v0.44.x" && !app.keeper.upgrade.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{"authz", "inflation"},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}

	// ica upgrade
	const upgradeName = "01-ica-upgrade"
	app.keeper.upgrade.SetUpgradeHandler(
		upgradeName,
		func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			fromVM[icatypes.ModuleName] = icaModule.ConsensusVersion()

			// create ICS27 Controller submodule params
			// enable the controller chain
			controllerParams := icacontrollertypes.Params{ControllerEnabled: true}

			// create ICS27 Host submodule params
			hostParams := icahosttypes.Params{
				// enable the host chain
				HostEnabled: true,
				// allowing the all messages
				AllowMessages: []string{"*"},
			}

			ctx.Logger().Info("start to init interchainaccount module...")
			// initialize ICS27 module
			icaModule.InitModule(ctx, controllerParams, hostParams)

			ctx.Logger().Info("start to run module migrations...")

			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)

	if upgradeInfo.Name == upgradeName && !app.keeper.upgrade.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{icacontrollertypes.StoreKey, icahosttypes.StoreKey},
		}

		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

func getGenesisTime(appOpts servertypes.AppOptions, homePath string) time.Time { // nolint: unused,deadcode
	if v := appOpts.Get("GenesisTime"); v != nil {
		// in tests, GenesisTime is supplied using appOpts
		genTime, ok := v.(time.Time)
		if !ok {
			panic("expected GenesisTime to be a Time value")
		}
		return genTime
	}

	genDoc, err := tmtypes.GenesisDocFromFile(filepath.Join(homePath, "config/genesis.json"))
	if err != nil {
		panic(err)
	}

	return genDoc.GenesisTime
}

// MakeCodecs constructs the *std.Codec and *codec.LegacyAmino instances used by
// simapp. It is useful for tests and clients who do not want to construct the
// full simapp
func MakeCodecs() (codec.Codec, *codec.LegacyAmino) {
	config := MakeEncodingConfig()
	return config.Marshaler, config.Amino
}

// Name returns the name of the App
func (app *AkashApp) Name() string { return app.BaseApp.Name() }

// InitChainer application update at chain initialization
func (app *AkashApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	app.keeper.upgrade.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// BeginBlocker is a function in which application updates every begin block
func (app *AkashApp) BeginBlocker(
	ctx sdk.Context, req abci.RequestBeginBlock,
) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker is a function in which application updates every end block
func (app *AkashApp) EndBlocker(
	ctx sdk.Context, req abci.RequestEndBlock,
) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// LegacyAmino returns AkashApp's amino codec.
func (app *AkashApp) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns AkashApp's app codec.
func (app *AkashApp) AppCodec() codec.Codec {
	return app.appCodec
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *AkashApp) ModuleAccountAddrs() map[string]bool {
	return MacAddrs()
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *AkashApp) BlockedAddrs() map[string]bool {
	perms := MacPerms()
	blockedAddrs := make(map[string]bool)
	for acc := range perms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}

// InterfaceRegistry returns AkashApp's InterfaceRegistry
func (app *AkashApp) InterfaceRegistry() codectypes.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
func (app *AkashApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
func (app *AkashApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
func (app *AkashApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.keeper.params.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *AkashApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *AkashApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	// Register legacy tx routes
	authrest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	// Register new tx routes from grpc-gateway
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics().RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics().RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *AkashApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *AkashApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(ctx client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// LoadHeight method of AkashApp loads baseapp application version with given height
func (app *AkashApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

func (AkashApp) TxConfig() client.TxConfig       { return nil }
func (AkashApp) AutoCliOpts() autocli.AppOptions { return autocli.AppOptions{} }

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key,
	tkey sdk.StoreKey,
) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)

	return akashSubspaces(paramsKeeper)
}
