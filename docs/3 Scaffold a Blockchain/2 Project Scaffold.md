# Project Scaffold Reference

The `starport app` command scaffolds a new Cosmos SDK blockchain project.

starport app github.com/hello/planet

This command will create a directory called `planet`, which contains all the files for your project. The `github.com` URL in the argument is a string that will be used for Go module's path. The repository name (`planet`, in this case) will be used as the project's name. A git repository will be initialized locally.

The project directory structure:

* `app`: files that wire the blockchain together
* `cmd`: blockchain node's binary
* `proto`: protocol buffer files for custom modules
* `x`: directory with custom modules
* `vue`: scaffolded web application (optional)
* `config.yml`: configuration file

Most of the logic of your application-specific blockchain is written in custom modules. Each module effectively encapsulates an independent piece of functionality. Custom modules are stored inside the `x` directory. By default, `starport app` scaffolds a module with a name that matches the name of the project. In our example, it will be `x/planet`.

Every Cosmos SDK module has protocol buffer files defining data structures, messages, queries, RPCs, etc. `proto` contains a directory with proto files per each custom module in `x`.

Global changes to your blockchain are defined in files inside the `app` directory. This includes importing third-party modules, defining relationships between modules, and configuring blockchain-wide settings.

`config.yml` is a file that contains configuration options that Starport uses to build, initialize and start your blockchain node in development.

## Address prefix

Account addresses on Cosmos SDK-based blockchains have string prefixes. For example, Cosmos Hub blockchain uses the default `cosmos` prefix, so that addresses look like this: `cosmos12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf`. 

When creating a new blockchain, pass a prefix as a value to the `--address-prefix` flag like so:

starport app github.com/hello/planet --address-prefix moonlight

Using this prefix, account addresses on your blockchain look like this: `moonlight12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf`.

To change the prefix manually after the blockchain has been scaffolded, modify the `AccountAddressPrefix` in `app/prefix.go`.



```


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
