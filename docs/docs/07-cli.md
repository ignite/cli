---
sidebar_position: 7
description: Ignite CLI docs.
---

# CLI Reference

Documentation for Ignite CLI.
## ignite

Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain

**Synopsis**

Ignite CLI is a tool for creating sovereign blockchains built with Cosmos SDK, the world’s
most popular modular blockchain framework. Ignite CLI offers everything you need to scaffold,
test, build, and launch your blockchain.

To get started, create a blockchain:

ignite scaffold chain github.com/username/mars

**Options**

```
  -h, --help   help for ignite
```

**SEE ALSO**

* [ignite account](#ignite-account)	 - Commands for managing accounts
* [ignite chain](#ignite-chain)	 - Build, initialize and start a blockchain node or perform other actions on the blockchain
* [ignite docs](#ignite-docs)	 - Show Ignite CLI docs
* [ignite generate](#ignite-generate)	 - Generate clients, API docs from source code
* [ignite relayer](#ignite-relayer)	 - Connect blockchains by using IBC protocol
* [ignite scaffold](#ignite-scaffold)	 - Scaffold a new blockchain, module, message, query, and more
* [ignite tools](#ignite-tools)	 - Tools for advanced users
* [ignite version](#ignite-version)	 - Print the current build information


## ignite account

Commands for managing accounts

**Synopsis**

Commands for managing accounts. An account is a pair of a private key and a public key.
Ignite CLI uses accounts to interact with the Ignite blockchain, use an IBC relayer, and more.

**Options**

```
  -h, --help   help for account
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain
* [ignite account create](#ignite-account-create)	 - Create a new account
* [ignite account delete](#ignite-account-delete)	 - Delete an account by name
* [ignite account export](#ignite-account-export)	 - Export an account as a private key
* [ignite account import](#ignite-account-import)	 - Import an account by using a mnemonic or a private key
* [ignite account list](#ignite-account-list)	 - Show a list of all accounts
* [ignite account show](#ignite-account-show)	 - Show detailed information about a particular account


## ignite account create

Create a new account

```
ignite account create [name] [flags]
```

**Options**

```
  -h, --help                     help for create
      --keyring-backend string   Keyring backend to store your account keys (default "test")
```

**SEE ALSO**

* [ignite account](#ignite-account)	 - Commands for managing accounts


## ignite account delete

Delete an account by name

```
ignite account delete [name] [flags]
```

**Options**

```
  -h, --help                     help for delete
      --keyring-backend string   Keyring backend to store your account keys (default "test")
```

**SEE ALSO**

* [ignite account](#ignite-account)	 - Commands for managing accounts


## ignite account export

Export an account as a private key

```
ignite account export [name] [flags]
```

**Options**

```
  -h, --help                     help for export
      --keyring-backend string   Keyring backend to store your account keys (default "test")
      --non-interactive          Do not enter into interactive mode
      --passphrase string        Account passphrase
      --path string              path to export private key. default: ./key_[name]
```

**SEE ALSO**

* [ignite account](#ignite-account)	 - Commands for managing accounts


## ignite account import

Import an account by using a mnemonic or a private key

```
ignite account import [name] [flags]
```

**Options**

```
  -h, --help                     help for import
      --keyring-backend string   Keyring backend to store your account keys (default "test")
      --non-interactive          Do not enter into interactive mode
      --passphrase string        Account passphrase
      --secret string            Your mnemonic or path to your private key (use interactive mode instead to securely pass your mnemonic)
```

**SEE ALSO**

* [ignite account](#ignite-account)	 - Commands for managing accounts


## ignite account list

Show a list of all accounts

```
ignite account list [flags]
```

**Options**

```
      --address-prefix string    Account address prefix (default "cosmos")
  -h, --help                     help for list
      --keyring-backend string   Keyring backend to store your account keys (default "test")
```

**SEE ALSO**

* [ignite account](#ignite-account)	 - Commands for managing accounts


## ignite account show

Show detailed information about a particular account

```
ignite account show [name] [flags]
```

**Options**

```
      --address-prefix string    Account address prefix (default "cosmos")
  -h, --help                     help for show
      --keyring-backend string   Keyring backend to store your account keys (default "test")
```

**SEE ALSO**

* [ignite account](#ignite-account)	 - Commands for managing accounts


## ignite chain

Build, initialize and start a blockchain node or perform other actions on the blockchain

**Synopsis**

Build, initialize and start a blockchain node or perform other actions on the blockchain.

**Options**

```
  -h, --help   help for chain
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain
* [ignite chain build](#ignite-chain-build)	 - Build a node binary
* [ignite chain faucet](#ignite-chain-faucet)	 - Send coins to an account
* [ignite chain init](#ignite-chain-init)	 - Initialize your chain
* [ignite chain serve](#ignite-chain-serve)	 - Start a blockchain node in development
* [ignite chain simulate](#ignite-chain-simulate)	 - Run simulation testing for the blockchain


## ignite chain build

Build a node binary

**Synopsis**

By default, build your node binaries
and add the binaries to your $(go env GOPATH)/bin path.

To build binaries for a release, use the --release flag. The app binaries
for one or more specified release targets are built in a release/ dir under the app's
source. Specify the release targets with GOOS:GOARCH build tags.
If the optional --release.targets is not specified, a binary is created for your current environment.

Sample usages:
	- ignite chain build
	- ignite chain build --release -t linux:amd64 -t darwin:amd64 -t darwin:arm64

```
ignite chain build [flags]
```

**Options**

```
      --clear-cache               Clear the build cache (advanced)
  -h, --help                      help for build
      --home string               Home directory used for blockchains
  -o, --output string             binary output path
  -p, --path string               path of the app (default ".")
      --proto-all-modules         Enables proto code generation for 3rd party modules used in your chain. Available only without the --release flag
      --release                   build for a release
      --release.prefix string     tarball prefix for each release target. Available only with --release flag
  -t, --release.targets strings   release targets. Available only with --release flag
  -v, --verbose                   Verbose output
```

**SEE ALSO**

* [ignite chain](#ignite-chain)	 - Build, initialize and start a blockchain node or perform other actions on the blockchain


## ignite chain faucet

Send coins to an account

```
ignite chain faucet [address] [coin<,...>] [flags]
```

**Options**

```
  -h, --help          help for faucet
      --home string   Home directory used for blockchains
  -p, --path string   path of the app (default ".")
  -v, --verbose       Verbose output
```

**SEE ALSO**

* [ignite chain](#ignite-chain)	 - Build, initialize and start a blockchain node or perform other actions on the blockchain


## ignite chain init

Initialize your chain

```
ignite chain init [flags]
```

**Options**

```
      --clear-cache   Clear the build cache (advanced)
  -h, --help          help for init
      --home string   Home directory used for blockchains
  -p, --path string   path of the app (default ".")
```

**SEE ALSO**

* [ignite chain](#ignite-chain)	 - Build, initialize and start a blockchain node or perform other actions on the blockchain


## ignite chain serve

Start a blockchain node in development

**Synopsis**

Start a blockchain node with automatic reloading

```
ignite chain serve [flags]
```

**Options**

```
      --clear-cache         Clear the build cache (advanced)
  -c, --config string       Ignite config file (default: ./config.yml)
  -f, --force-reset         Force reset of the app state on start and every source change
  -h, --help                help for serve
      --home string         Home directory used for blockchains
  -p, --path string         path of the app (default ".")
      --proto-all-modules   Enables proto code generation for 3rd party modules used in your chain
  -r, --reset-once          Reset of the app state on first start
  -v, --verbose             Verbose output
```

**SEE ALSO**

* [ignite chain](#ignite-chain)	 - Build, initialize and start a blockchain node or perform other actions on the blockchain


## ignite chain simulate

Run simulation testing for the blockchain

**Synopsis**

Run simulation testing for the blockchain. It sends many randomized-input messages of each module to a simulated node and checks if invariants break

```
ignite chain simulate [flags]
```

**Options**

```
      --blockSize int             operations per block (default 30)
      --exportParamsHeight int    height to which export the randomly generated params
      --exportParamsPath string   custom file path to save the exported params JSON
      --exportStatePath string    custom file path to save the exported app state JSON
      --exportStatsPath string    custom file path to save the exported simulation statistics JSON
      --genesis string            custom simulation genesis file; cannot be used with params file
      --genesisTime int           override genesis UNIX time instead of using a random UNIX time
  -h, --help                      help for simulate
      --initialBlockHeight int    initial block to start the simulation (default 1)
      --lean                      lean simulation log output
      --numBlocks int             number of new blocks to simulate from the initial block height (default 200)
      --params string             custom simulation params file which overrides any random params; cannot be used with genesis
      --period uint               run slow invariants only once every period assertions
      --printAllInvariants        print all invariants if a broken invariant is found
      --seed int                  simulation random seed (default 42)
      --simulateEveryOperation    run slow invariants every operation
  -v, --verbose                   verbose log output
```

**SEE ALSO**

* [ignite chain](#ignite-chain)	 - Build, initialize and start a blockchain node or perform other actions on the blockchain


## ignite docs

Show Ignite CLI docs

```
ignite docs [flags]
```

**Options**

```
  -h, --help   help for docs
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain


## ignite generate

Generate clients, API docs from source code

**Synopsis**

Generate clients, API docs from source code.

Such as compiling protocol buffer files into Go or implement particular functionality, for example, generating an OpenAPI spec.

Produced source code can be regenerated by running a command again and is not meant to be edited by hand.

**Options**

```
      --clear-cache   Clear the build cache (advanced)
  -h, --help          help for generate
  -p, --path string   path of the app (default ".")
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain
* [ignite generate dart](#ignite-generate-dart)	 - Generate a Dart client
* [ignite generate openapi](#ignite-generate-openapi)	 - Generate generates an OpenAPI spec for your chain from your config.yml
* [ignite generate proto-go](#ignite-generate-proto-go)	 - Generate proto based Go code needed for the app's source code
* [ignite generate vuex](#ignite-generate-vuex)	 - Generate Vuex store for you chain's frontend from your config.yml


## ignite generate dart

Generate a Dart client

```
ignite generate dart [flags]
```

**Options**

```
  -h, --help   help for dart
  -y, --yes    Answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
      --clear-cache   Clear the build cache (advanced)
  -p, --path string   path of the app (default ".")
```

**SEE ALSO**

* [ignite generate](#ignite-generate)	 - Generate clients, API docs from source code


## ignite generate openapi

Generate generates an OpenAPI spec for your chain from your config.yml

```
ignite generate openapi [flags]
```

**Options**

```
  -h, --help   help for openapi
  -y, --yes    Answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
      --clear-cache   Clear the build cache (advanced)
  -p, --path string   path of the app (default ".")
```

**SEE ALSO**

* [ignite generate](#ignite-generate)	 - Generate clients, API docs from source code


## ignite generate proto-go

Generate proto based Go code needed for the app's source code

```
ignite generate proto-go [flags]
```

**Options**

```
  -h, --help   help for proto-go
  -y, --yes    Answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
      --clear-cache   Clear the build cache (advanced)
  -p, --path string   path of the app (default ".")
```

**SEE ALSO**

* [ignite generate](#ignite-generate)	 - Generate clients, API docs from source code


## ignite generate vuex

Generate Vuex store for you chain's frontend from your config.yml

```
ignite generate vuex [flags]
```

**Options**

```
  -h, --help                help for vuex
      --proto-all-modules   Enables proto code generation for 3rd party modules used in your chain
  -y, --yes                 Answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
      --clear-cache   Clear the build cache (advanced)
  -p, --path string   path of the app (default ".")
```

**SEE ALSO**

* [ignite generate](#ignite-generate)	 - Generate clients, API docs from source code


## ignite relayer

Connect blockchains by using IBC protocol

**Options**

```
  -h, --help   help for relayer
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain
* [ignite relayer configure](#ignite-relayer-configure)	 - Configure source and target chains for relaying
* [ignite relayer connect](#ignite-relayer-connect)	 - Link chains associated with paths and start relaying tx packets in between


## ignite relayer configure

Configure source and target chains for relaying

```
ignite relayer configure [flags]
```

**Options**

```
  -a, --advanced                  Advanced configuration options for custom IBC modules
  -h, --help                      help for configure
      --keyring-backend string    Keyring backend to store your account keys (default "test")
      --ordered                   Set the channel as ordered
  -r, --reset                     Reset the relayer config
      --source-account string     Source Account
      --source-client-id string   use a custom client id for source
      --source-faucet string      Faucet address of the source chain
      --source-gaslimit int       Gas limit used for transactions on source chain
      --source-gasprice string    Gas price used for transactions on source chain
      --source-port string        IBC port ID on the source chain
      --source-prefix string      Address prefix of the source chain
      --source-rpc string         RPC address of the source chain
      --source-version string     Module version on the source chain
      --target-account string     Target Account
      --target-client-id string   use a custom client id for target
      --target-faucet string      Faucet address of the target chain
      --target-gaslimit int       Gas limit used for transactions on target chain
      --target-gasprice string    Gas price used for transactions on target chain
      --target-port string        IBC port ID on the target chain
      --target-prefix string      Address prefix of the target chain
      --target-rpc string         RPC address of the target chain
      --target-version string     Module version on the target chain
```

**SEE ALSO**

* [ignite relayer](#ignite-relayer)	 - Connect blockchains by using IBC protocol


## ignite relayer connect

Link chains associated with paths and start relaying tx packets in between

```
ignite relayer connect [<path>,...] [flags]
```

**Options**

```
  -h, --help                     help for connect
      --keyring-backend string   Keyring backend to store your account keys (default "test")
```

**SEE ALSO**

* [ignite relayer](#ignite-relayer)	 - Connect blockchains by using IBC protocol


## ignite scaffold

Scaffold a new blockchain, module, message, query, and more

**Synopsis**

Scaffold commands create and modify the source code files to add functionality.

CRUD stands for "create, read, update, delete".

**Options**

```
  -h, --help   help for scaffold
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain
* [ignite scaffold band](#ignite-scaffold-band)	 - Scaffold an IBC BandChain query oracle to request real-time data
* [ignite scaffold chain](#ignite-scaffold-chain)	 - Fully-featured Cosmos SDK blockchain
* [ignite scaffold flutter](#ignite-scaffold-flutter)	 - A Flutter app for your chain
* [ignite scaffold list](#ignite-scaffold-list)	 - CRUD for data stored as an array
* [ignite scaffold map](#ignite-scaffold-map)	 - CRUD for data stored as key-value pairs
* [ignite scaffold message](#ignite-scaffold-message)	 - Message to perform state transition on the blockchain
* [ignite scaffold module](#ignite-scaffold-module)	 - Scaffold a Cosmos SDK module
* [ignite scaffold packet](#ignite-scaffold-packet)	 - Message for sending an IBC packet
* [ignite scaffold query](#ignite-scaffold-query)	 - Query to get data from the blockchain
* [ignite scaffold single](#ignite-scaffold-single)	 - CRUD for data stored in a single location
* [ignite scaffold type](#ignite-scaffold-type)	 - Scaffold only a type definition
* [ignite scaffold vue](#ignite-scaffold-vue)	 - Vue 3 web app template


## ignite scaffold band

Scaffold an IBC BandChain query oracle to request real-time data

**Synopsis**

Scaffold an IBC BandChain query oracle to request real-time data from BandChain scripts in a specific IBC-enabled Cosmos SDK module

```
ignite scaffold band [queryName] --module [moduleName] [flags]
```

**Options**

```
      --clear-cache     Clear the build cache (advanced)
  -h, --help            help for band
      --module string   IBC Module to add the packet into
  -p, --path string     path of the app (default ".")
      --signer string   Label for the message signer (default: creator)
  -y, --yes             Answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## ignite scaffold chain

Fully-featured Cosmos SDK blockchain

**Synopsis**

Create a new application-specific Cosmos SDK blockchain.

For example, the following command will create a blockchain called "hello" in the "hello/" directory:

  ignite scaffold chain hello

A project name can be a simple name or a URL. The name will be used as the Go module path for the project. Examples of project names:

  ignite scaffold chain foo
  ignite scaffold chain foo/bar
  ignite scaffold chain example.org/foo
  ignite scaffold chain github.com/username/foo
		
A new directory with source code files will be created in the current directory. To use a different path use the "--path" flag.

Most of the logic of your blockchain is written in custom modules. Each module effectively encapsulates an independent piece of functionality. Following the Cosmos SDK convention, custom modules are stored inside the "x/" directory. By default, Ignite creates a module with a name that matches the name of the project. To create a blockchain without a default module use the "--no-module" flag. Additional modules can be added after a project is created with "ignite scaffold module" command.

Account addresses on Cosmos SDK-based blockchains have string prefixes. For example, the Cosmos Hub blockchain uses the default "cosmos" prefix, so that addresses look like this: "cosmos12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf". To use a custom address prefix use the "--address-prefix" flag. For example:

  ignite scaffold chain foo --address-prefix bar

By default when compiling a blockchain's source code Ignite creates a cache to speed up the build process. To clear the cache when building a blockchain use the "--clear-cache" flag. It is very unlikely you will ever need to use this flag.

The blockchain is using the Cosmos SDK modular blockchain framework. Learn more about Cosmos SDK on https://docs.cosmos.network

```
ignite scaffold chain [name] [flags]
```

**Options**

```
      --address-prefix string   Account address prefix (default "cosmos")
      --clear-cache             Clear the build cache (advanced)
  -h, --help                    help for chain
      --no-module               Create a project without a default module
  -p, --path string             Create a project in a specific path (default ".")
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## ignite scaffold flutter

A Flutter app for your chain

```
ignite scaffold flutter [flags]
```

**Options**

```
  -h, --help          help for flutter
  -p, --path string   path to scaffold content of the Flutter app (default "./flutter")
  -y, --yes           Answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## ignite scaffold list

CRUD for data stored as an array

```
ignite scaffold list NAME [field]... [flags]
```

**Options**

```
      --clear-cache     Clear the build cache (advanced)
  -h, --help            help for list
      --module string   Module to add into. Default is app's main module
      --no-message      Disable CRUD interaction messages scaffolding
      --no-simulation   Disable CRUD simulation scaffolding
  -p, --path string     path of the app (default ".")
      --signer string   Label for the message signer (default: creator)
  -y, --yes             Answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## ignite scaffold map

CRUD for data stored as key-value pairs

```
ignite scaffold map NAME [field]... [flags]
```

**Options**

```
      --clear-cache     Clear the build cache (advanced)
  -h, --help            help for map
      --index strings   fields that index the value (default [index])
      --module string   Module to add into. Default is app's main module
      --no-message      Disable CRUD interaction messages scaffolding
      --no-simulation   Disable CRUD simulation scaffolding
  -p, --path string     path of the app (default ".")
      --signer string   Label for the message signer (default: creator)
  -y, --yes             Answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## ignite scaffold message

Message to perform state transition on the blockchain

```
ignite scaffold message [name] [field1] [field2] ... [flags]
```

**Options**

```
      --clear-cache        Clear the build cache (advanced)
  -d, --desc string        Description of the command
  -h, --help               help for message
      --module string      Module to add the message into. Default: app's main module
      --no-simulation      Disable CRUD simulation scaffolding
  -p, --path string        path of the app (default ".")
  -r, --response strings   Response fields
      --signer string      Label for the message signer (default: creator)
  -y, --yes                Answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## ignite scaffold module

Scaffold a Cosmos SDK module

**Synopsis**

Scaffold a new Cosmos SDK module in the `x` directory

```
ignite scaffold module [name] [flags]
```

**Options**

```
      --clear-cache            Clear the build cache (advanced)
      --dep strings            module dependencies (e.g. --dep account,bank)
  -h, --help                   help for module
      --ibc                    scaffold an IBC module
      --ordering string        channel ordering of the IBC module [none|ordered|unordered] (default "none")
      --params strings         scaffold module params
  -p, --path string            path of the app (default ".")
      --require-registration   if true command will fail if module can't be registered
  -y, --yes                    Answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## ignite scaffold packet

Message for sending an IBC packet

**Synopsis**

Scaffold an IBC packet in a specific IBC-enabled Cosmos SDK module

```
ignite scaffold packet [packetName] [field1] [field2] ... --module [moduleName] [flags]
```

**Options**

```
      --ack strings     Custom acknowledgment type (field1,field2,...)
      --clear-cache     Clear the build cache (advanced)
  -h, --help            help for packet
      --module string   IBC Module to add the packet into
      --no-message      Disable send message scaffolding
  -p, --path string     path of the app (default ".")
      --signer string   Label for the message signer (default: creator)
  -y, --yes             Answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## ignite scaffold query

Query to get data from the blockchain

```
ignite scaffold query [name] [request_field1] [request_field2] ... [flags]
```

**Options**

```
      --clear-cache        Clear the build cache (advanced)
  -d, --desc string        Description of the command
  -h, --help               help for query
      --module string      Module to add the query into. Default: app's main module
      --paginated          Define if the request can be paginated
  -p, --path string        path of the app (default ".")
  -r, --response strings   Response fields
  -y, --yes                Answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## ignite scaffold single

CRUD for data stored in a single location

```
ignite scaffold single NAME [field]... [flags]
```

**Options**

```
      --clear-cache     Clear the build cache (advanced)
  -h, --help            help for single
      --module string   Module to add into. Default is app's main module
      --no-message      Disable CRUD interaction messages scaffolding
      --no-simulation   Disable CRUD simulation scaffolding
  -p, --path string     path of the app (default ".")
      --signer string   Label for the message signer (default: creator)
  -y, --yes             Answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## ignite scaffold type

Scaffold only a type definition

```
ignite scaffold type NAME [field]... [flags]
```

**Options**

```
      --clear-cache     Clear the build cache (advanced)
  -h, --help            help for type
      --module string   Module to add into. Default is app's main module
      --no-message      Disable CRUD interaction messages scaffolding
      --no-simulation   Disable CRUD simulation scaffolding
  -p, --path string     path of the app (default ".")
      --signer string   Label for the message signer (default: creator)
  -y, --yes             Answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## ignite scaffold vue

Vue 3 web app template

```
ignite scaffold vue [flags]
```

**Options**

```
  -h, --help          help for vue
  -p, --path string   path to scaffold content of the Vue.js app (default "./vue")
  -y, --yes           Answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## ignite tools

Tools for advanced users

**Options**

```
  -h, --help   help for tools
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain
* [ignite tools completions](#ignite-tools-completions)	 - Generate completions script
* [ignite tools ibc-relayer](#ignite-tools-ibc-relayer)	 - Typescript implementation of an IBC relayer
* [ignite tools ibc-setup](#ignite-tools-ibc-setup)	 - Collection of commands to quickly setup a relayer
* [ignite tools protoc](#ignite-tools-protoc)	 - Execute the protoc command


## ignite tools completions

Generate completions script

**Synopsis**

 The completions command outputs a completion script you can use in your shell. The output script requires 
				that [bash-completion](https://github.com/scop/bash-completion)	is installed and enabled in your 
				system. Since most Unix-like operating systems come with bash-completion by default, bash-completion 
				is probably already installed and operational.

Bash:

  $ source <(ignite  tools completions bash)

  To load completions for every new session, run:

  ** Linux **
  $ ignite  tools completions bash > /etc/bash_completion.d/ignite

  ** macOS **
  $ ignite  tools completions bash > /usr/local/etc/bash_completion.d/ignite

Zsh:

  If shell completions is not already enabled in your environment, you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  To load completions for each session, execute once:
  
  $ ignite  tools completions zsh > "${fpath[1]}/_ignite"

  You will need to start a new shell for this setup to take effect.

fish:

  $ ignite  tools completions fish | source

  To load completions for each session, execute once:
  
  $ ignite  tools completions fish > ~/.config/fish/completionss/ignite.fish

PowerShell:

  PS> ignite  tools completions powershell | Out-String | Invoke-Expression

  To load completions for every new session, run:
  
  PS> ignite  tools completions powershell > ignite.ps1
  
  and source this file from your PowerShell profile.


```
ignite tools completions
```

**Options**

```
  -h, --help   help for completions
```

**SEE ALSO**

* [ignite tools](#ignite-tools)	 - Tools for advanced users


## ignite tools ibc-relayer

Typescript implementation of an IBC relayer

```
ignite tools ibc-relayer [--] [...] [flags]
```

**Examples**

```
ignite tools ibc-relayer -- -h
```

**Options**

```
  -h, --help   help for ibc-relayer
```

**SEE ALSO**

* [ignite tools](#ignite-tools)	 - Tools for advanced users


## ignite tools ibc-setup

Collection of commands to quickly setup a relayer

```
ignite tools ibc-setup [--] [...] [flags]
```

**Examples**

```
ignite tools ibc-setup -- -h
ignite tools ibc-setup -- init --src relayer_test_1 --dest relayer_test_2
```

**Options**

```
  -h, --help   help for ibc-setup
```

**SEE ALSO**

* [ignite tools](#ignite-tools)	 - Tools for advanced users


## ignite tools protoc

Execute the protoc command

**Synopsis**

The protoc command. You don't need to setup the global protoc include folder with -I, it's automatically handled

```
ignite tools protoc [--] [...] [flags]
```

**Examples**

```
ignite tools protoc -- --version
```

**Options**

```
  -h, --help   help for protoc
```

**SEE ALSO**

* [ignite tools](#ignite-tools)	 - Tools for advanced users


## ignite version

Print the current build information

```
ignite version [flags]
```

**Options**

```
  -h, --help   help for version
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain

