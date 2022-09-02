---
sidebar_position: 5
---

# Types

Now that you've defined messages that trigger state transitions, it's time to implement types and methods that operate on the state.

A keeper is an abstraction that lets your blockchain app interact with the state. Functions like create, update, and delete are defined as keeper methods. In the scavenge blockchain, you need to define the `scavenge` and `commit` types along with create and update methods.

Several Ignite CLI commands are available to scaffold the code for CRUD functionality for a list-like data structure, a map (key-value pairs), and a single element in the state. In this example, both `scavenge` and `commit` are stored in a map-like data structure.

## Scavenge

Use the `ignite scaffold map` command to scaffold the `scavenge` type and the code for creating, reading, updating, and deleting (CRUD) scavenges.

The first argument is the name of the type to create (`scavenge`), the rest is a list of fields. By default, generic CRUD messages are scaffolded. However, since you already created messages specifically for this scavenge blockchain, use the `--no-message` flag to skip message creation.

```bash
ignite scaffold map scavenge solutionHash solution description reward scavenger --no-message
```

The `ignite scaffold map` command creates and modifies several files:

```
modify proto/scavenge/genesis.proto
modify proto/scavenge/query.proto
create proto/scavenge/scavenge.proto
modify x/scavenge/client/cli/query.go
create x/scavenge/client/cli/query_scavenge.go
create x/scavenge/client/cli/query_scavenge_test.go
modify x/scavenge/genesis.go
modify x/scavenge/genesis_test.go
create x/scavenge/keeper/grpc_query_scavenge.go
create x/scavenge/keeper/grpc_query_scavenge_test.go
create x/scavenge/keeper/scavenge.go
create x/scavenge/keeper/scavenge_test.go
modify x/scavenge/module.go
modify x/scavenge/types/genesis.go
modify x/scavenge/types/genesis_test.go
create x/scavenge/types/key_scavenge.go

🎉 scavenge added.
```

The `scaffold map` command does all of these code updates for you:

* `proto/scavenge/scavenge.proto`

  * Defines the `Scavenge` type as a proto message

* `proto/scavenge/query.proto`

  * Defines queries to get data from the blockchain as proto messages and registers the queries in the `Query` service

* `proto/scavenge/genesis.proto`

  * Creates type for exporting the state of the blockchain (for example, during software upgrades)

* `x/scavenge/keeper/grpc_query_scavenge.go`

  * Defines keeper methods to query the blockchain

* `x/scavenge/keeper/grpc_query_scavenge_test.go`

  * Creates tests for query keeper methods

* `x/scavenge/keeper/scavenge.go`

  * Defines keeper methods to get, set, and remove scavenges from the store

* `x/scavenge/keeper/scavenge_test.go`

  * Creates tests for the keeper methods

* `x/scavenge/client/cli/query_scavenge.go`

  * Creates CLI commands for querying the blockchain

* `x/scavenge/client/cli/query.go`

  * Registers the CLI commands

* `x/scavenge/client/cli/query_scavenge_test.go`

  * Createstests for the CLI commands

* `x/scavenge/types/keys.go`

  * Creates a string as a prefix in the key used to store scavenges in the state

* `x/scavenge/genesis.go`

  * Creates logic for exporting and exporting the state

* `x/scavenge/types/genesis.go`

  * Createslogic for validating the genesis file

* `x/scavenge/module.go`

  * Registers the gRPC gateway routes

Review the `x/scavenge/keeper/scavenge.go` file to see the `SetScavenge` updates that were made in the `keeper` package, like the `store.Set` method that writes a Scavenge into the store:

```go
// SetScavenge set a specific scavenge in the store from its index
func (k Keeper) SetScavenge(ctx sdk.Context, scavenge types.Scavenge) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ScavengeKeyPrefix))
	b := k.cdc.MustMarshal(&scavenge)
	store.Set(types.ScavengeKey(
		scavenge.Index,
	), b)
}
```

Review the update for `GetScavenge` that selects a store using the scavenge prefix and uses `store.Get` to fetch a scavenge with a particular index.

## Commit

Use `ignite scaffold map` to create the same logic for a `commit` type.

```bash
ignite scaffold map commit solutionHash solutionScavengerHash --no-message
```

## Save changes

Now is a good time to store your project in a git commit:

```bash
git add .
git commit -m "add scavenge types"
```
