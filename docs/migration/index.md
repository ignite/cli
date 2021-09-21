---
order: 1
parent:
  title: Migration
  order: 3
description: Required changes when upgrading from Starport v0.17.3 to Starport v0.18.
---

The default template in Cosmos SDK versions lower than v0.44 are not compatible with Starport v0.18.

Changes are required when upgrade from Starport v0.17.3 to Starport v0.18.0 and later. 

To update the default template so your Starport installation is compatible with Cosmos SDK v0.44, make these changes to the blockchain template after you upgrade to Starport v0.18.

# v0.18

## Blockchain

### app/app.go

```go
import (
  //...
  "github.com/cosmos/cosmos-sdk/x/feegrant"
  feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
  feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"

  "github.com/cosmos/ibc-go/modules/apps/transfer"
  ibctransferkeeper "github.com/cosmos/ibc-go/modules/apps/transfer/keeper"
  ibctransfertypes "github.com/cosmos/ibc-go/modules/apps/transfer/types"
  ibc "github.com/cosmos/ibc-go/modules/core"
  ibcclient "github.com/cosmos/ibc-go/modules/core/02-client"
  ibcporttypes "github.com/cosmos/ibc-go/modules/core/05-port/types"
  ibchost "github.com/cosmos/ibc-go/modules/core/24-host"
  ibckeeper "github.com/cosmos/ibc-go/modules/core/keeper"
  
  // transfer "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer"
  // ibctransferkeeper "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/keeper"
  // ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
  // ibc "github.com/cosmos/cosmos-sdk/x/ibc/core"
  // ibcclient "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client"
  // porttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/05-port/types"
  // ibchost "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
  // ibckeeper "github.com/cosmos/cosmos-sdk/x/ibc/core/keeper"
)

var (
  //...
  ModuleBasics = module.NewBasicManager(
    //...
    slashing.AppModuleBasic{},
    // Add feegrantmodule.AppModuleBasic{},
    feegrantmodule.AppModuleBasic{}, // <--
    ibc.AppModuleBasic{},
    //...
  )
  //...
)

type App struct {
  //...
  // Replace codec.Marshaler with codec.Codec
  appCodec          codec.Codec // <--
  // Add FeeGrantKeeper
  FeeGrantKeeper   feegrantkeeper.Keeper // <--
}

``go
func New(...) {
  //bApp.SetAppVersion(version.Version)
  bApp.SetVersion(version.Version) // <--

  keys := sdk.NewKVStoreKeys(
    //...
    upgradetypes.StoreKey,
    // Add feegrant.StoreKey
    feegrant.StoreKey, // <--
    evidencetypes.StoreKey,
    //...
  )

  app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, keys[feegrant.StoreKey], app.AccountKeeper)  // <--
  // Add app.BaseApp as the last argument to upgradekeeper.NewKeeper
  app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, app.BaseApp)
  
  app.IBCKeeper = ibckeeper.NewKeeper(
    // Add app.UpgradeKeeper
    appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), app.StakingKeeper, app.UpgradeKeeper, scopedIBCKeeper,
  )

  govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
    //...
    // Replace NewClientUpdateProposalHandler with NewClientUpdateProposalHandler
    AddRoute(ibchost.RouterKey, ibcclient.NewClientUpdateProposalHandler(app.IBCKeeper.ClientKeeper))

  // Replace porttypes with ibcporttypes
  ibcRouter := ibcporttypes.NewRouter()

  app.mm.SetOrderBeginBlockers(
    upgradetypes.ModuleName,
    // Add capabilitytypes.ModuleName,
    capabilitytypes.ModuleName,
    minttypes.ModuleName,
    //...
    // Add feegrant.ModuleName,
    feegrant.ModuleName,
  )

  app.mm.RegisterServices(module.NewConfigurator(
    // Add app.appCodec
    app.appCodec,
    //...
  )
  
  // Replace:
  // app.SetAnteHandler(
  // 	ante.NewAnteHandler(
  // 		app.AccountKeeper, app.BankKeeper, ante.DefaultSigVerificationGasConsumer,
  // 		encodingConfig.TxConfig.SignModeHandler(),
  // 	),
  // )

  // With the following:
  anteHandler, err := ante.NewAnteHandler(
    ante.HandlerOptions{
      AccountKeeper:   app.AccountKeeper,
      BankKeeper:      app.BankKeeper,
      SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
      FeegrantKeeper:  app.FeeGrantKeeper,
      SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
    },
  )
  if err != nil {
    panic(err)
  }
  app.SetAnteHandler(anteHandler)

  // Remove the following:
  // ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})
  // app.CapabilityKeeper.InitializeAndSeal(ctx)
}

func (app *App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
  var genesisState GenesisState
  // Replace tmjson with json
  if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
    panic(err)
  }
  // Add the following:
  app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
  return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// Replace Marshaler with Codec
func (app *App) AppCodec() codec.Codec {
  return app.appCodec
}

// Replace BinaryMarshaler with BinaryCodec
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
  //...
}
```

### app/genesis.go

```go
// Replace codec.JSONMarshaler with codec.JSONCodec
func NewDefaultGenesisState(cdc codec.JSONCodec) GenesisState {
  //...
}
```

### testutil/keeper/mars.go

Add the following code:

```go
package keeper

import (
  "testing"

  "github.com/cosmonaut/mars/x/mars/keeper"
  "github.com/cosmonaut/mars/x/mars/types"
  "github.com/cosmos/cosmos-sdk/codec"
  codectypes "github.com/cosmos/cosmos-sdk/codec/types"
  "github.com/cosmos/cosmos-sdk/store"
  storetypes "github.com/cosmos/cosmos-sdk/store/types"
  sdk "github.com/cosmos/cosmos-sdk/types"
  "github.com/stretchr/testify/require"
  "github.com/tendermint/tendermint/libs/log"
  tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
  tmdb "github.com/tendermint/tm-db"
)

func MarsKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
  storeKey := sdk.NewKVStoreKey(types.StoreKey)
  memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

  db := tmdb.NewMemDB()
  stateStore := store.NewCommitMultiStore(db)
  stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
  stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
  require.NoError(t, stateStore.LoadLatestVersion())

  registry := codectypes.NewInterfaceRegistry()
  k := keeper.NewKeeper(
    codec.NewProtoCodec(registry),
    storeKey,
    memStoreKey,
  )

  ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
  return k, ctx
}
```

### testutil/network/network.go

```go
func DefaultConfig() network.Config {
  //...
  return network.Config{
    //...
    // Add sdk.DefaultPowerReduction
    AccountTokens:   sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction),
    StakingTokens:   sdk.TokensFromConsensusPower(500, sdk.DefaultPowerReduction),
    BondedTokens:    sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction),
    //...
  }
}
```

### testutil/sample/sample.go

Add the following code:

```go
package sample

import (
  "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
  sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccAddress returns a sample account address
func AccAddress() string {
  pk := ed25519.GenPrivKey().PubKey()
  addr := pk.Address()
  return sdk.AccAddress(addr).String()
}
```

## Module

### x/mars/keeper/keeper.go

```go
type (
  Keeper struct {
    // Replace Marshaler with BinaryCodec
    cdc      codec.BinaryCodec
    //...
  }
)

func NewKeeper(
  // Replace Marshaler with BinaryCodec
  cdc codec.BinaryCodec,
  //...
) *Keeper {
  // ...
}
```

### x/mars/keeper/msg_server_test.go

```go
package keeper_test

import (
  //...
  // Add the following:
  keepertest "github.com/cosmonaut/mars/testutil/keeper"
  "github.com/cosmonaut/mars/x/mars/keeper"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
  // Replace
  // keeper, ctx := setupKeeper(t)
  // return NewMsgServerImpl(*keeper), sdk.WrapSDKContext(ctx)

  // With the following:
  k, ctx := keepertest.MarsKeeper(t)
  return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
```

### x/mars/module.go

```go
type AppModuleBasic struct {
  // Replace Marshaler with BinaryCodec
  cdc codec.BinaryCodec
}

// Replace Marshaler with BinaryCodec
func NewAppModuleBasic(cdc codec.BinaryCodec) AppModuleBasic {
  return AppModuleBasic{cdc: cdc}
}

// Replace JSONMarshaler with JSONCodec
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
  return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// Replace JSONMarshaler with JSONCodec
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
  //...
}

// Replace JSONMarshaler with JSONCodec
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
  //...
}

// Replace JSONMarshaler with JSONCodec
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
  //...
}

// Replace JSONMarshaler with JSONCodec
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
  //...
}

// Add the following
func (AppModule) ConsensusVersion() uint64 { return 2 }
```