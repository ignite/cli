---
order: 19
---

# Complete App

When you used the `starport type` command, your application has already been incorporated in the `./app.go` file. 

> _*NOTE*_: Your application needs to import the code you just wrote. Here the import path is set to this repository (`github.com/cosmos/sdk-tutorials/nameservice/x/nameservice`). If you are following along in your own repo you will need to change the import path to reflect that (`github.com/{ .Username }/{ .Project.Repo }/x/nameservice`).

Inside the `./app/app.go` file, it made the following changes:

- Instantiate required `Keepers` from each desired module.
- Generate `storeKeys` required by each `Keeper`.
- Register `Handler`s from each module. The `AddRoute()` method from `baseapp`'s `router` is used to this end.
- Register `Querier`s from each module. The `AddRoute()` method from `baseapp`'s `queryRouter` is used to this end.
- Mount `KVStore`s to the provided keys in the `baseApp` multistore.
- Set the `initChainer` for defining the initial application state.

As a result, the file should look like this - 

<<< @/nameservice/nameservice/app/app.go

> _*NOTE*_: The TransientStore mentioned above is an in-memory implementation of the KVStore for state that is not persisted.

> _*NOTE*_: Pay attention to how the modules are initiated: the order matters! Here the sequence goes Auth --> Bank --> Feecollection --> Staking --> Distribution --> Slashing, then the hooks were set for the staking module. This is because some of these modules depend on others existing before they can be used.

You'll notice a few functions at the end of the file. The `initChainer` defines how accounts in `genesis.json` are mapped into the application state on initial chain start. The `ExportAppStateAndValidators` function helps bootstrap the initial state for the application. `BeginBlocker` and `EndBlocker` are optional methods module developers can implement in their module. They will be triggered at the beginning and at the end of each block respectively, when the `BeginBlock` and `EndBlock` ABCI messages are received from the underlying consensus engine.

### Now, it's time to update your entrypoints.
