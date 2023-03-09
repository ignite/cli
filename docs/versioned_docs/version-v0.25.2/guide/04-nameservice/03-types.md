---
sidebar_position: 3
description: Implement types and methods that operate on the state.
---

# Implement Types

Now that you've defined messages that trigger state transitions, it's time to implement types and methods that operate on the state.

> The Cosmos SDK relies on keepers. A keeper is an abstraction that lets your blockchain app interact with the state. Functions like create, read, update, and delete (CRUD) are defined as keeper methods.

For the nameservice blockchain, define a `whois` type and the create and delete methods.

Because Ignite CLI does the heavy lifting for you, choose from several [ignite scaffold](../../../references/cli#ignite-scaffold) commands to create CRUD functionality code for data stored in different ways:

- Array, a list-like data structure
- Map (key-value pairs)
- In a single location  

## Add the whois Type

Use the `ignite scaffold map` command to scaffold the `whois` type and create the code that implements CRUD functionality to create, read, update, and delete information about names.

In this example, the `whois` type is stored in a map-like data structure:

```bash
ignite scaffold map whois name value price owner --no-message
```

where:

- whois is the type
- name is the name the user sets
- value is the name that name resolves to
- price is the bid
- `--no-message` flag skips message creation

    By default, generic CRUD messages are scaffolded. However, you've already created messages specifically for this blockchain, so you can skip message creation with the `--no-message` flag.

The `ignite scaffold map whois name value price --no-message` command created and modified several files:

* `proto/nameservice/whois.proto`

    Defines the `Whois` type as a proto message.

* `proto/nameservice/query.proto`

    * Queries to get data from the blockchain.
    * Define queries as proto messages.
    * Register the messages in the `Query` service.

* `proto/nameservice/genesis.proto`

    A type for exporting the state of the blockchain, for example, during software upgrades.

* `x/nameservice/keeper/grpc_query_whois.go`

    Keeper methods to query the blockchain.

* `x/nameservice/keeper/grpc_query_whois_test.go`

    Tests for query keeper methods.

* `x/nameservice/keeper/whois.go`

    Keeper methods to get, set, and remove whois information from the store.

* `x/nameservice/keeper/whois_test.go`

    Tests for keeper methods.

* `x/nameservice/client/cli/query_whois.go`

    CLI commands for querying the blockchain.

* `x/nameservice/client/cli/query.go`

    Registers the CLI commands.

* `x/nameservice/client/cli/query_whois_test.go`

    Tests for CLI commands.

* `x/nameservice/types/keys.go`

    String prefix in the key to store whois information in the state.

* `x/nameservice/genesis.go`

    Logic for exporting the state.

* `x/nameservice/types/genesis.go`

    Logic for validating the genesis file.

* `x/nameservice/module.go`

    Registers gRPC gateway routes.

## Keeper Package

In the `x/nameservice/keeper/whois.go` file, take at a look at the keeper package.

- `SetWhois` uses a key-value store with a prefix for the `Whois` type and uses a `store.Set` method to write a `Whois` into the store.

<!-- where is this? teach me please
`Whois-value-` encodes the `Whois` type that is generated from a protocol buffer definition-->

- `GetWhois` selects a store using the `Whois` prefix and uses a `store.Get` method to fetch a `Whois` with a particular index.

The keeper package also includes `RemoveWhois` and `GetAllWhois`.
