---
order: 5
---

# Types

Now that you've defined messages that trigger state transitions, it's time to implement types and methods that operate on the state.

A keeper is an abstraction that let's your blockchain app interact with the state. Functions like create, update and delete are defined as keeper methods. In the Scavenge blockchain a `scavenge` and `commit` types need to be defined, along with create and update methods.

Starport has several commands that scaffold the code for CRUD functionality for a list-like data structure, a map (key-value pairs) and a single element in the state. In this example, both `scavenge` and `commit` will be stored in a map-like data structure.

## Scavenge

Use `starport scaffold map` command to scaffold the `scavenge` type and the code for creating, reading, updating, and deleting (CRUD) scavenges. The first argument is the name of the type being created (`scavenge`), the rest is list of fields. By default, generic CRUD messages are scaffolded, but since you've already created messages specifically for this blockchain, skip messages with a `--no-message` flag.

```
starport scaffold map scavenge solutionHash solution description reward scavenger --no-message
```

`starport scaffold map` created and mofidied several files:

* `proto/scavenge/scavenge.proto`: the `Scavenge` type defined as a proto message.
* `proto/scavenge/query.proto`: queries to get data from the blockchain defined as proto messages and registered in the `Query` service.
* `proto/scavenge/genesis.proto`: a type for exporting the state of the blockchain (for example, during software upgrades)
* `x/scavenge/keeper/grpc_query_scavenge.go`: keeper methods to query the blockchain.
* `x/scavenge/keeper/grpc_query_scavenge_test.go`: tests for query keeper methods.
* `x/scavenge/keeper/scavenge.go`: keper methods to get, set and remove scavenges from the store.
* `x/scavenge/keeper/scavenge_test.go`: tests for keeper methods.
* `x/scavenge/client/cli/query_scavenge.go`: CLI commands for querying the blockchain.
* `x/scavenge/client/cli/query.go`: registering CLI commands.
* `x/scavenge/client/cli/query_scavenge_test.go`: tests for CLI commands.
* `x/scavenge/types/keys.go`: a string as a prefix in the key used to store scavenges in the state.
* `x/scavenge/genesis.go`: logic for exporting and exporting the state.
* `x/scavenge/types/genesis.go`: logic for validating the genesis file.
* `x/scavenge/module.go`: registering gRPC gateway routes.

`SetScavenge` in the `keeper` package uses a key-value store using a prefix for the scavenge type (`Scavenge-value-`) encodes the `Scavenge` type (that is generated from a protocol buffer definition) and uses `store.Set` method to write a Scavenge into the store.

`GetScavenge` selects a store using the scavenge prefix, and uses `store.Get` to fetch a scavenge with a particular index.

## Commit

Use `starport scaffold map` to create the same logic for a `commit` type.

```
starport scaffold map commit solutionHash solutionScavengerHash --no-message
```