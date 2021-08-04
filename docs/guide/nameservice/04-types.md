---
order: 4
---

# Types

Now that you've defined messages that trigger state transitions, it's time to implement types and methods that operate on the state.

A keeper is an abstraction that let's your blockchain app interact with the state. Functions like create, update and delete are defined as keeper methods. In the Nameservice blockchain a `whois` type needs to be defined, along with create and delete methods.

Starport has several commands that scaffold the code for CRUD functionality for a list-like data structure, a map (key-value pairs) and a single element in the state. In this example, `whois` will be stored in a map-like data structure.

## Whois

Use `starport scaffold map` command to scaffold the `whois` type and the code for creating, reading, updating, and deleting (CRUD) information about names. The first argument is the name of the type being created (`whois`), the rest is list of fields. By default, generic CRUD messages are scaffolded, but since you've already created messages specifically for this blockchain, skip messages with a `--no-message` flag.

```
starport scaffold map whois name value price --no-message
```

`starport scaffold map` created and mofidied several files:

* `proto/nameservice/whois.proto`: the `Whois` type defined as a proto message.
* `proto/nameservice/query.proto`: queries to get data from the blockchain defined as proto messages and registered in the `Query` service.
* `proto/nameservice/genesis.proto`: a type for exporting the state of the blockchain (for example, during software upgrades)
* `x/nameservice/keeper/grpc_query_whois.go`: keeper methods to query the blockchain.
* `x/nameservice/keeper/grpc_query_whois_test.go`: tests for query keeper methods.
* `x/nameservice/keeper/whois.go`: keper methods to get, set and remove whois information from the store.
* `x/nameservice/keeper/whois_test.go`: tests for keeper methods.
* `x/nameservice/client/cli/query_whois.go`: CLI commands for querying the blockchain.
* `x/nameservice/client/cli/query.go`: registering CLI commands.
* `x/nameservice/client/cli/query_whois_test.go`: tests for CLI commands.
* `x/nameservice/types/keys.go`: a string as a prefix in the key used to store whois information in the state.
* `x/nameservice/genesis.go`: logic for exporting and exporting the state.
* `x/nameservice/types/genesis.go`: logic for validating the genesis file.
* `x/nameservice/module.go`: registering gRPC gateway routes.

`SetWhois` in the `keeper` package uses a key-value store using a prefix for the whois type (`Whois-value-`) encodes the `Whois` type (that is generated from a protocol buffer definition) and uses `store.Set` method to write a Whois into the store.

`GetWhois` selects a store using the whois prefix, and uses `store.Get` to fetch a `whois` with a particular index.