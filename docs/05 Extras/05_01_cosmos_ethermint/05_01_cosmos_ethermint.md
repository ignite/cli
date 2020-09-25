# Add Ethermint manually to a Cosmos SDK chain

Ethermint extends the Cosmos SDK with the Ethereum Virtual Machine and allows for uploading and interacting with smart contracts. It is designed to mimic the Ethereum blockchain and allow for direct interoperability of the two blockchains as well as using smart contracts written in Solidity or any other language that is supported by Ethereum.

When starting with the Cosmos SDK ethermint is not part of the default modules that are installed on a basic blockchain application. In this tutorial we will be looking what steps are necessary to get started with the Ethereum Virtual Machine on Cosmos.

Let's use Starport to bootstrap our basic blockchain

`starport app github.com/username/ethapp && cd ethapp`

Currently Ethermint is being developed by ChainSafe, until the repository `github.com/cosmos/ethermint` is up to date, we need to modify the `go.mod` file, at the end of it, place the following `replace` repository command:

```go
replace github.com/cosmos/ethermint => github.com/ChainSafe/ethermint v0.2.0-rc4
```
 
The Ethereum Virtual Machine, in short `evm` takes a few additions to our `/app/app.go` file. First, we need to add them as imports

Add to import
```go
"github.com/cosmos/ethermint/app/ante"
"github.com/cosmos/ethermint/x/evm"
```

Next, we need to add it as module to our `ModuleBasics`

```go

	ModuleBasics    = module.NewBasicManager(
		genutil.AppModuleBasic{},
        auth.AppModuleBasic{},
        ...

		// this line is for Ethermint
		evm.AppModuleBasic{}, // <-----
        // this line is used by starport scaffolding # 2
    )
```

When browsing a bit further in the code, you can see the struct `NewApp`, in here we want to add our `EvmKeeper`

```go
type NewApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey

	subspaces map[string]params.Subspace

    accountKeeper auth.AccountKeeper
    ...
    // this line is for Ethermint
	EvmKeeper evm.Keeper // <-----
    // this line is used by starport scaffolding # 3
    ...

}

```

Everything easy so far, now we have to make a change in our `NewInitApp` function, right at the beginning there is a function for transaction decoding, we want to change this to accept Ethereum transactions, a few lines later we will add the evm Key-Value store in `NewKVStoreKeys`.

```go
func NewInitApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	invCheckPeriod uint,
	baseAppOptions ...func(*bam.BaseApp),
) *NewApp {
	cdc := MakeCodec()

	// the TxDecoder has to be changed to evm txdecoder
	// NOTE we use custom Ethermint transaction decoder that supports the sdk.Tx interface instead of sdk.StdTx
	bApp := bam.NewBaseApp(appName, logger, db, evm.TxDecoder(cdc), baseAppOptions...) // <------
	bApp.SetCommitMultiStoreTracer(traceStore)
    bApp.SetAppVersion(version.Version)
    
    keys := sdk.NewKVStoreKeys(
    bam.MainStoreKey,
    auth.StoreKey,
    staking.StoreKey,
    supply.StoreKey,
    params.StoreKey,
    ethapptypes.StoreKey,
    // this line is used for the Ethermint
    evm.StoreKey, // <------
    // this line is used by starport scaffolding # 5
    )

    ...
```

and again a few lines later, we initialize the params keeper and subspaces, we add the `evm` subspace as follows:

```go
app.subspaces[evm.ModuleName] = app.paramsKeeper.Subspace(evm.DefaultParamspace)
```

and later initialize our evm keeper

```go
app.EvmKeeper = evm.NewKeeper(
		app.cdc, keys[evm.StoreKey], app.subspaces[evm.ModuleName], app.accountKeeper,
)
```

Next we add the `evm` to the module manager, also add the evm module to `SetOrderEndBlockers` 

```go
	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		ethapp.NewAppModule(app.ethappKeeper, app.bankKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		// this line is used by Ethermint
		evm.NewAppModule(app.EvmKeeper, app.accountKeeper), // <-----
		// this line is used by starport scaffolding # 6
    )
    

	// Next line is changed by Ethermint
	app.mm.SetOrderEndBlockers(staking.ModuleName, evm.ModuleName)
```

In genesis we should also add the `evm`:

```go
	app.mm.SetOrderInitGenesis(
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		ethapptypes.ModuleName,
		supply.ModuleName,
		genutil.ModuleName,
		// this line is used for Ethermint
		evm.ModuleName, // <-----
		// this line is used by starport scaffolding # 7
	)
```

Additionally, we have to set the `ante` handler, to compile ethereum transactions:

replace the `SetAnteHandler` with:


```go
	app.SetAnteHandler(
		ante.NewAnteHandler(
			app.accountKeeper,
			app.EvmKeeper,
			app.supplyKeeper,
		),
	)
```

After this, you can see the evm available via the command line.

Run the application with 

```bash
starport serve
```

Now you can see all transaction types with

```bash
ethappcli tx --help
```

In the response you should see the `evm` listed as

```bash
 evm         EVM transaction subcommands
```

Find more information on how to use the module in on the help files or documentation:

```bash
ethappcli tx evm --help
```

https://docs.ethermint.zone/

Let's start creating and uploading contracts.

Next:

- [Create a Fungible Token](04%20usecases/02_erc20/02_erc20.md)
- [Create a Non-Fungible Token](04%20usecases/04_nft/04_nft.md)