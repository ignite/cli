---
order: 16
---

# AppModule Interface

The Cosmos SDK provides a standard interface for modules. This [`AppModule`](https://github.com/cosmos/cosmos-sdk/blob/master/types/module.go) interface requires modules to provide a set of methods used by the `ModuleBasicsManager` to incorporate them into your application.

We should already have a `module.go` file in `./nameservice`, and we don't need to change anything, but it should look like this.

<<< @/nameservice/nameservice/x/nameservice/module.go

To see more examples of AppModule implementation, check out some of the other modules in the SDK such as [x/staking](https://github.com/cosmos/cosmos-sdk/blob/master/x/staking/genesis.go)

### Next, we need to implement the genesis-specific methods called above.
