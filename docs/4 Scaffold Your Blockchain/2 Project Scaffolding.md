# Project Scaffolding

Scaffolding a project using Starport is done with the `starport app` command.

The default project version scaffolded in Starport is Launchpad version.

To scaffold a project , or a Stargate version. The version can be specified by passing the `sdk-version` flag, followed by either `stargate` or `launchpad`.

ie.

```
starport app github.com/user/app --sdk-version=stargate
```

## Address prefix

You can change the way addresses look in your blockchain. On the Cosmos SDK Main Hub, addresses have a `cosmos` prefix, like `cosmos12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf`.

To specify the address prefix on the command line, use the `--address-prefix` parameter. For example, to change the blockchain prefix to moonlight:

```
starport app github.com/foo/bar --address-prefix moonlight
```

To change the address prefix for subsequent blockchain builds:

1. Change the `AccountAddressPrefix` variable in `/app/prefix.go`. Do not change other variables in the file.
2. To recognize the new prefix, change the `VUE_APP_ADDRESS_PREFIX` variable in `/vue/.env`.

## Stargate app

Scaffolding a Stargate app currently uses version `^0.40` of the Cosmos SDK.

A typical directory structure for a Stargate app `foo` will contain the following structure:

```
├── app
│   ├── app.go
│   ├── encoding.go
│   ├── export.go
│   ├── genesis.go
│   ├── params
│   │   ├── encoding.go
│   │   └── proto.go
│   └── types.go
├── cmd
│   └── food
│       ├── cmd
│       │   ├── app.go
│       │   ├── genaccounts.go
│       │   └── root.go
│       └── main.go
├── config.yml
├── go.mod
├── go.sum
├── internal
│   └── tools
│       └── tools.go
├── proto
│   ├── cosmos
│   │   └── base
│   │       └── query
│   │           └── v1beta1
│   │               └── pagination.proto
│   └── foo
│       └── v1beta
│           ├── genesis.proto
│           └── querier.proto
├── readme.md
├── scripts
│   └── protocgen
├── third_party
│   └── proto
│       ├── confio
│       │   └── proofs.proto
│       ├── cosmos_proto
│       │   └── cosmos.proto
│       ├── gogoproto
│       │   └── gogo.proto
│       ├── google
│       │   ├── api
│       │   │   ├── annotations.proto
│       │   │   ├── http.proto
│       │   │   └── httpbody.proto
│       │   └── protobuf
│       │       ├── any.proto
│       │       └── descriptor.proto
│       └── tendermint
│           ├── abci
│           │   └── types.proto
│           ├── crypto
│           │   ├── keys.proto
│           │   └── proof.proto
│           ├── libs
│           │   └── bits
│           │       └── types.proto
│           ├── types
│           │   ├── evidence.proto
│           │   ├── params.proto
│           │   ├── types.proto
│           │   └── validator.proto
│           └── version
│               └── types.proto
└── x
    └── foo
        ├── client
        │   ├── cli
        │   │   ├── query.go
        │   │   └── tx.go
        │   └── rest
        │       └── rest.go
        ├── genesis.go
        ├── handler.go
        ├── keeper
        │   ├── grpc_query.go
        │   ├── keeper.go
        │   └── querier.go
        ├── module.go
        └── types
            ├── codec.go
            ├── errors.go
            ├── genesis.go
            ├── genesis.pb.go
            ├── keys.go
            ├── querier.go
            ├── querier.pb.go
            └── types.go
```

## Launchpad app

Scaffolding a Launchpad app is currently the default that is being used by Starport, and uses version `0.39.x` of the Cosmos SDK.

A typical directory structure for a Launchpad app `bar` will contain the following structure:

```
├── app
│   ├── app.go
│   ├── export.go
│   └── prefix.go
├── cmd
│   ├── barcli
│   │   └── main.go
│   └── bard
│       ├── genaccounts.go
│       └── main.go
├── config.yml
├── go.mod
├── go.sum
├── readme.md
├── vue
│   ├── README.md
│   ├── babel.config.js
│   ├── package-lock.json
│   ├── package.json
│   ├── public
│   │   ├── favicon.ico
│   │   └── index.html
│   ├── src
│   │   ├── App.vue
│   │   ├── main.js
│   │   ├── router
│   │   │   └── index.js
│   │   ├── store
│   │   │   └── index.js
│   │   └── views
│   │       └── Index.vue
│   └── vue.config.js
└── x
    └── bar
        ├── abci.go
        ├── client
        │   ├── cli
        │   │   ├── query.go
        │   │   └── tx.go
        │   └── rest
        │       └── rest.go
        ├── genesis.go
        ├── handler.go
        ├── keeper
        │   ├── keeper.go
        │   ├── params.go
        │   └── querier.go
        ├── module.go
        ├── spec
        │   └── README.md
        └── types
            ├── codec.go
            ├── errors.go
            ├── events.go
            ├── expected_keepers.go
            ├── genesis.go
            ├── key.go
            ├── msg.go
            ├── params.go
            ├── querier.go
            └── types.go
```
