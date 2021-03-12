# Project Scaffold Reference

The `starport app` command scaffolds a project. By default, the Cosmos SDK version is Stargate. <!-- what is a project? compared to a "blockchain" or "app" -->

## Address prefix

You can change the way addresses look in your blockchain.

On the Cosmos SDK Hub, addresses have a `cosmos` prefix, like `cosmos12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf`.

To specify a custom address prefix on the command line, use the `--address-prefix` flag. For example, to change the blockchain prefix to moonlight:

```
starport app github.com/foo/bar --address-prefix moonlight
```

To change the address prefix for subsequent blockchain builds:

1. Change the `AccountAddressPrefix` variable in the `/app/prefix.go` file. Be sure to preserve other variables in the file.
2. To recognize the new prefix, change the `VUE_APP_ADDRESS_PREFIX` variable in `/vue/.env`.

## Stargate app

```
starport app github.com/foo/bar
```

Scaffolding a Stargate app currently uses version `^0.42` of the Cosmos SDK.
<!-- is there a way to use a release variable here? -->
A typical directory structure for a Stargate app `foo` contains the following structure: <!-- how can I verify this? and did this nifty structure get auto-generated? -->

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

To scaffold an app using the earlier `launchpad` release, use the `--sdk-version` flag:

```
starport app github.com/foo/bar --sdk-version launchpad
```
<!-- do we want to use a different value for [github.com/org/repo] in the launchpad command? -->
Scaffolding a Starport app on launchpad uses version `^0.39` of the Cosmos SDK.

A typical directory structure for a Stargate app `foo` contains the following structure:


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
