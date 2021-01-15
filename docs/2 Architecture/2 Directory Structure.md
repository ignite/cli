# Directory Structure

A typical Cosmos SDK Module directory structure is as follows

```
x/{module}
├── client
│   ├── cli
│   │   ├── query.go
│   │   └── tx.go
│   └── rest
│       └── rest.go
├── genesis.go
├── handler.go
├── keeper
│   ├── grpc_query.go
│   ├── keeper.go
│   └── query.go
├── module.go
└── types
    ├── codec.go
    ├── errors.go
    ├── genesis.go
    ├── genesis.pb.go
    ├── keys.go
    ├── query.go
    ├── query.pb.go
    └── types.go
```

`client/`: The module's CLI and REST client functionality implementation and testing.

`handler.go`: The module's message handlers.

`keeper/`: The module's keeper implementation along with any auxiliary implementations such as the querier and invariants.

`types/`: The module's type definitions such as messages, KVStore keys, parameter types, Protocol Buffer definitions, and expected_keepers.go contracts.

`module.go`: The module's implementation of the AppModule and AppModuleBasic interfaces.
