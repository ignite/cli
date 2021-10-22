---
order: 1
description: Starport CLI docs. 
parent:
  order: 8
  title: CLI Reference
---

# CLI Reference

Documentation for Starport CLI.

## starport

Starport offers everything you need to scaffold, test, build, and launch your blockchain

**Synopsis**

Starport is a tool for creating sovereign blockchains built with Cosmos SDK, the world’s
most popular modular blockchain framework. Starport offers everything you need to scaffold,
test, build, and launch your blockchain.

To get started, create a blockchain:

starport scaffold chain github.com/cosmonaut/mars

**Options**

```
  -h, --help   help for starport
```

**SEE ALSO**

* [starport account](#starport-account)	 - Commands for managing accounts
* [starport chain](#starport-chain)	 - Build, initialize and start a blockchain node or perform other actions on the blockchain
* [starport docs](#starport-docs)	 - Show Starport docs
* [starport generate](#starport-generate)	 - Generate clients, API docs from source code
* [starport network](#starport-network)	 - Launch a blockchain network in production
* [starport relayer](#starport-relayer)	 - Connect blockchains by using IBC protocol
* [starport scaffold](#starport-scaffold)	 - Scaffold a new blockchain, module, message, query, and more
* [starport tools](#starport-tools)	 - Tools for advanced users
* [starport version](#starport-version)	 - Print the current build information


## starport account

Commands for managing accounts

**Synopsis**

Commands for managing accounts. An account is a pair of a private key and a public key.
Starport uses accounts to interact with the Starport Network blockchain, use an IBC relayer, and more.

**Options**

```
  -h, --help   help for account
```

**SEE ALSO**

* [starport](#starport)	 - Starport offers everything you need to scaffold, test, build, and launch your blockchain
* [starport account create](#starport-account-create)	 - Create a new account
* [starport account delete](#starport-account-delete)	 - Delete an account by name
* [starport account export](#starport-account-export)	 - Export an account as a private key
* [starport account import](#starport-account-import)	 - Import an account by using a mnemonic or a private key
* [starport account list](#starport-account-list)	 - Show a list of all accounts
* [starport account show](#starport-account-show)	 - Show detailed information about a particular account


## starport account create

Create a new account

```
starport account create [name] [flags]
```

**Options**

```
  -h, --help                     help for create
      --keyring-backend string   Keyring backend to store your account keys (default "test")
```

**SEE ALSO**

* [starport account](#starport-account)	 - Commands for managing accounts


## starport account delete

Delete an account by name

```
starport account delete [name] [flags]
```

**Options**

```
  -h, --help                     help for delete
      --keyring-backend string   Keyring backend to store your account keys (default "test")
```

**SEE ALSO**

* [starport account](#starport-account)	 - Commands for managing accounts


## starport account export

Export an account as a private key

```
starport account export [name] [flags]
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

* [starport account](#starport-account)	 - Commands for managing accounts


## starport account import

Import an account by using a mnemonic or a private key

```
starport account import [name] [flags]
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

* [starport account](#starport-account)	 - Commands for managing accounts


## starport account list

Show a list of all accounts

```
starport account list [flags]
```

**Options**

```
      --address-prefix string    Account address prefix (default "cosmos")
  -h, --help                     help for list
      --keyring-backend string   Keyring backend to store your account keys (default "test")
```

**SEE ALSO**

* [starport account](#starport-account)	 - Commands for managing accounts


## starport account show

Show detailed information about a particular account

```
starport account show [name] [flags]
```

**Options**

```
      --address-prefix string    Account address prefix (default "cosmos")
  -h, --help                     help for show
      --keyring-backend string   Keyring backend to store your account keys (default "test")
```

**SEE ALSO**

* [starport account](#starport-account)	 - Commands for managing accounts


## starport chain

Build, initialize and start a blockchain node or perform other actions on the blockchain

**Synopsis**

Build, initialize and start a blockchain node or perform other actions on the blockchain.

**Options**

```
  -h, --help          help for chain
  -p, --path string   path of the app (default ".")
```

**SEE ALSO**

* [starport](#starport)	 - Starport offers everything you need to scaffold, test, build, and launch your blockchain
* [starport chain build](#starport-chain-build)	 - Build a node binary
* [starport chain faucet](#starport-chain-faucet)	 - Send coins to an account
* [starport chain init](#starport-chain-init)	 - Initialize your chain
* [starport chain serve](#starport-chain-serve)	 - Start a blockchain node in development


## starport chain build

Build a node binary

**Synopsis**

By default, build your node binaries
and add the binaries to your $(go env GOPATH)/bin path.

To build binaries for a release, use the --release flag. The app binaries
for one or more specified release targets are built in a release/ dir under the app's
source. Specify the release targets with GOOS:GOARCH build tags.
If the optional --release.targets is not specified, a binary is created for your current environment.

Sample usages:
	- starport chain build
	- starport chain build --release -t linux:amd64 -t darwin:amd64 -t darwin:arm64

```
starport chain build [flags]
```

**Options**

```
  -h, --help                      help for build
      --home string               Home directory used for blockchains
  -o, --output string             binary output path
      --proto-all-modules         Enables proto code generation for 3rd party modules used in your chain. Available only without the --release flag
      --release                   build for a release
      --release.prefix string     tarball prefix for each release target. Available only with --release flag
  -t, --release.targets strings   release targets. Available only with --release flag
  -v, --verbose                   Verbose output
```

**Options inherited from parent commands**

```
  -p, --path string   path of the app (default ".")
```

**SEE ALSO**

* [starport chain](#starport-chain)	 - Build, initialize and start a blockchain node or perform other actions on the blockchain


## starport chain faucet

Send coins to an account

```
starport chain faucet [address] [coin<,...>] [flags]
```

**Options**

```
  -h, --help          help for faucet
      --home string   Home directory used for blockchains
  -v, --verbose       Verbose output
```

**Options inherited from parent commands**

```
  -p, --path string   path of the app (default ".")
```

**SEE ALSO**

* [starport chain](#starport-chain)	 - Build, initialize and start a blockchain node or perform other actions on the blockchain


## starport chain init

Initialize your chain

```
starport chain init [flags]
```

**Options**

```
  -h, --help          help for init
      --home string   Home directory used for blockchains
```

**Options inherited from parent commands**

```
  -p, --path string   path of the app (default ".")
```

**SEE ALSO**

* [starport chain](#starport-chain)	 - Build, initialize and start a blockchain node or perform other actions on the blockchain


## starport chain serve

Start a blockchain node in development

**Synopsis**

Start a blockchain node with automatic reloading

```
starport chain serve [flags]
```

**Options**

```
  -c, --config string       Starport config file (default: ./config.yml)
  -f, --force-reset         Force reset of the app state on start and every source change
  -h, --help                help for serve
      --home string         Home directory used for blockchains
      --proto-all-modules   Enables proto code generation for 3rd party modules used in your chain
  -r, --reset-once          Reset of the app state on first start
  -v, --verbose             Verbose output
```

**Options inherited from parent commands**

```
  -p, --path string   path of the app (default ".")
```

**SEE ALSO**

* [starport chain](#starport-chain)	 - Build, initialize and start a blockchain node or perform other actions on the blockchain


## starport docs

Show Starport docs

```
starport docs [flags]
```

**Options**

```
  -h, --help   help for docs
```

**SEE ALSO**

* [starport](#starport)	 - Starport offers everything you need to scaffold, test, build, and launch your blockchain


## starport generate

Generate clients, API docs from source code

**Synopsis**

Generate clients, API docs from source code.

Such as compiling protocol buffer files into Go or implement particular functionality, for example, generating an OpenAPI spec.

Produced source code can be regenerated by running a command again and is not meant to be edited by hand.

**Options**

```
  -h, --help          help for generate
  -p, --path string   path of the app (default ".")
```

**SEE ALSO**

* [starport](#starport)	 - Starport offers everything you need to scaffold, test, build, and launch your blockchain
* [starport generate openapi](#starport-generate-openapi)	 - Generate generates an OpenAPI spec for your chain from your config.yml
* [starport generate proto-go](#starport-generate-proto-go)	 - Generate proto based Go code needed for the app's source code
* [starport generate vuex](#starport-generate-vuex)	 - Generate Vuex store for you chain's frontend from your config.yml


## starport generate openapi

Generate generates an OpenAPI spec for your chain from your config.yml

```
starport generate openapi [flags]
```

**Options**

```
  -h, --help   help for openapi
```

**Options inherited from parent commands**

```
  -p, --path string   path of the app (default ".")
```

**SEE ALSO**

* [starport generate](#starport-generate)	 - Generate clients, API docs from source code


## starport generate proto-go

Generate proto based Go code needed for the app's source code

```
starport generate proto-go [flags]
```

**Options**

```
  -h, --help   help for proto-go
```

**Options inherited from parent commands**

```
  -p, --path string   path of the app (default ".")
```

**SEE ALSO**

* [starport generate](#starport-generate)	 - Generate clients, API docs from source code


## starport generate vuex

Generate Vuex store for you chain's frontend from your config.yml

```
starport generate vuex [flags]
```

**Options**

```
  -h, --help                help for vuex
      --proto-all-modules   Enables proto code generation for 3rd party modules used in your chain
```

**Options inherited from parent commands**

```
  -p, --path string   path of the app (default ".")
```

**SEE ALSO**

* [starport generate](#starport-generate)	 - Generate clients, API docs from source code


## starport network

Launch a blockchain network in production

**Options**

```
  -h, --help                        help for network
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport](#starport)	 - Starport offers everything you need to scaffold, test, build, and launch your blockchain
* [starport network account](#starport-network-account)	 - Show the underlying SPN account
* [starport network chain](#starport-network-chain)	 - Build networks
* [starport network proposal](#starport-network-proposal)	 - Proposals related to starting network


## starport network account

Show the underlying SPN account

**Synopsis**

Show the underlying SPN account.
To pick another account see "starport network account use -h"
If no account is picked, default "spn" account is used.


```
starport network account [flags]
```

**Options**

```
  -h, --help   help for account
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network](#starport-network)	 - Launch a blockchain network in production
* [starport network account create](#starport-network-account-create)	 - Create an account
* [starport network account export](#starport-network-account-export)	 - Export an account
* [starport network account import](#starport-network-account-import)	 - Import an account
* [starport network account use](#starport-network-account-use)	 - Pick an account to be used with Starport Network


## starport network account create

Create an account

```
starport network account create [name] [flags]
```

**Options**

```
  -h, --help   help for create
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network account](#starport-network-account)	 - Show the underlying SPN account


## starport network account export

Export an account

```
starport network account export [flags]
```

**Options**

```
  -a, --account string   path to save private key (default "[account in use]")
  -h, --help             help for export
  -p, --path string      path to save private key (default "[account].key")
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network account](#starport-network-account)	 - Show the underlying SPN account


## starport network account import

Import an account

```
starport network account import [name] [password] [path-to-private-key] [flags]
```

**Options**

```
  -h, --help   help for import
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network account](#starport-network-account)	 - Show the underlying SPN account


## starport network account use

Pick an account to be used with Starport Network

**Synopsis**

Pick one of the accounts in OS keyring to put into use or provide one with --name flag.
Picked account will be used while interacting with Starport Network.

```
starport network account use [flags]
```

**Options**

```
  -h, --help          help for use
  -n, --name string   Account name to put into use
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network account](#starport-network-account)	 - Show the underlying SPN account


## starport network chain

Build networks

**Options**

```
  -h, --help   help for chain
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network](#starport-network)	 - Launch a blockchain network in production
* [starport network chain create](#starport-network-chain-create)	 - Create a new network
* [starport network chain join](#starport-network-chain-join)	 - Propose to join to a network as a validator
* [starport network chain list](#starport-network-chain-list)	 - List all chains with proposals summary
* [starport network chain show](#starport-network-chain-show)	 - Show details of a chain
* [starport network chain start](#starport-network-chain-start)	 - Start network


## starport network chain create

Create a new network

```
starport network chain create [chain] [source] [flags]
```

**Options**

```
      --branch string    Git branch to use
      --genesis string   URL to a custom Genesis
  -h, --help             help for create
      --home string      Home directory used for blockchains
      --tag string       Git tag to use
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network chain](#starport-network-chain)	 - Build networks


## starport network chain join

Propose to join to a network as a validator

```
starport network chain join [chain-id] [flags]
```

**Options**

```
      --gentx string             Path to a gentx file (optional)
  -h, --help                     help for join
      --home string              Home directory used for blockchains
      --keyring-backend string   Keyring backend used for the blockchain account (default "os")
      --peer string              Configure peer in node-id@host:port format (optional)
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network chain](#starport-network-chain)	 - Build networks


## starport network chain list

List all chains with proposals summary

```
starport network chain list [flags]
```

**Options**

```
  -h, --help            help for list
      --search string   List chains with the specified prefix in chain id
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network chain](#starport-network-chain)	 - Build networks


## starport network chain show

Show details of a chain

```
starport network chain show [chain-id] [flags]
```

**Options**

```
      --genesis   Show exclusively the genesis of the chain
  -h, --help      help for show
      --peers     Show exclusively the peers of the chain
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network chain](#starport-network-chain)	 - Build networks


## starport network chain start

Start network

```
starport network chain start [chain-id] [-- <flags>...] [flags]
```

**Options**

```
  -h, --help          help for start
      --home string   Home directory used for blockchains
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network chain](#starport-network-chain)	 - Build networks


## starport network proposal

Proposals related to starting network

**Options**

```
  -h, --help   help for proposal
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network](#starport-network)	 - Launch a blockchain network in production
* [starport network proposal approve](#starport-network-proposal-approve)	 - Approve proposals
* [starport network proposal list](#starport-network-proposal-list)	 - List all pending proposals
* [starport network proposal reject](#starport-network-proposal-reject)	 - Reject proposals
* [starport network proposal show](#starport-network-proposal-show)	 - Show details of a proposal
* [starport network proposal verify](#starport-network-proposal-verify)	 - Simulate and verify proposals validity


## starport network proposal approve

Approve proposals

```
starport network proposal approve [chain-id] [number<,...>] [flags]
```

**Options**

```
  -h, --help              help for approve
      --home string       Home directory used for blockchains
      --no-verification   approve the proposals without verifying them
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network proposal](#starport-network-proposal)	 - Proposals related to starting network


## starport network proposal list

List all pending proposals

```
starport network proposal list [chain-id] [flags]
```

**Options**

```
  -h, --help            help for list
      --status string   Filter proposals by status (pending|approved|rejected)
      --type string     Filter proposals by type (add-account|add-validator)
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network proposal](#starport-network-proposal)	 - Proposals related to starting network


## starport network proposal reject

Reject proposals

```
starport network proposal reject [chain-id] [number<,...>] [flags]
```

**Options**

```
  -h, --help   help for reject
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network proposal](#starport-network-proposal)	 - Proposals related to starting network


## starport network proposal show

Show details of a proposal

```
starport network proposal show [chain-id] [number] [flags]
```

**Options**

```
  -h, --help   help for show
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network proposal](#starport-network-proposal)	 - Proposals related to starting network


## starport network proposal verify

Simulate and verify proposals validity

```
starport network proposal verify [chain-id] [number<,...>] [flags]
```

**Options**

```
      --debug         show output of verification command in the console
  -h, --help          help for verify
      --home string   Home directory used for blockchains
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-api-address string      SPN api address (default "https://rest.alpha.starport.network")
      --spn-faucet-address string   SPN Faucet address (default "https://faucet.alpha.starport.network")
      --spn-node-address string     SPN node address (default "https://rpc.alpha.starport.network:443")
```

**SEE ALSO**

* [starport network proposal](#starport-network-proposal)	 - Proposals related to starting network


## starport relayer

Connect blockchains by using IBC protocol

**Options**

```
  -h, --help   help for relayer
```

**SEE ALSO**

* [starport](#starport)	 - Starport offers everything you need to scaffold, test, build, and launch your blockchain
* [starport relayer configure](#starport-relayer-configure)	 - Configure source and target chains for relaying
* [starport relayer connect](#starport-relayer-connect)	 - Link chains associated with paths and start relaying tx packets in between


## starport relayer configure

Configure source and target chains for relaying

```
starport relayer configure [flags]
```

**Options**

```
  -a, --advanced                 Advanced configuration options for custom IBC modules
  -h, --help                     help for configure
      --keyring-backend string   Keyring backend to store your account keys (default "test")
      --ordered                  Set the channel as ordered
      --source-account string    Source Account
      --source-faucet string     Faucet address of the source chain
      --source-gaslimit int      Gas limit used for transactions on source chain
      --source-gasprice string   Gas price used for transactions on source chain
      --source-port string       IBC port ID on the source chain
      --source-prefix string     Address prefix of the source chain
      --source-rpc string        RPC address of the source chain
      --source-version string    Module version on the source chain
      --target-account string    Target Account
      --target-faucet string     Faucet address of the target chain
      --target-gaslimit int      Gas limit used for transactions on target chain
      --target-gasprice string   Gas price used for transactions on target chain
      --target-port string       IBC port ID on the target chain
      --target-prefix string     Address prefix of the target chain
      --target-rpc string        RPC address of the target chain
      --target-version string    Module version on the target chain
```

**SEE ALSO**

* [starport relayer](#starport-relayer)	 - Connect blockchains by using IBC protocol


## starport relayer connect

Link chains associated with paths and start relaying tx packets in between

```
starport relayer connect [<path>,...] [flags]
```

**Options**

```
  -h, --help                     help for connect
      --keyring-backend string   Keyring backend to store your account keys (default "test")
```

**SEE ALSO**

* [starport relayer](#starport-relayer)	 - Connect blockchains by using IBC protocol


## starport scaffold

Scaffold a new blockchain, module, message, query, and more

**Synopsis**

Scaffold commands create and modify the source code files to add functionality.

CRUD stands for "create, read, update, delete".

**Options**

```
  -h, --help   help for scaffold
```

**SEE ALSO**

* [starport](#starport)	 - Starport offers everything you need to scaffold, test, build, and launch your blockchain
* [starport scaffold band](#starport-scaffold-band)	 - Scaffold an IBC BandChain query oracle to request real-time data
* [starport scaffold chain](#starport-scaffold-chain)	 - Fully-featured Cosmos SDK blockchain
* [starport scaffold list](#starport-scaffold-list)	 - CRUD for data stored as an array
* [starport scaffold map](#starport-scaffold-map)	 - CRUD for data stored as key-value pairs
* [starport scaffold message](#starport-scaffold-message)	 - Message to perform state transition on the blockchain
* [starport scaffold module](#starport-scaffold-module)	 - Scaffold a Cosmos SDK module
* [starport scaffold packet](#starport-scaffold-packet)	 - Message for sending an IBC packet
* [starport scaffold query](#starport-scaffold-query)	 - Query to get data from the blockchain
* [starport scaffold single](#starport-scaffold-single)	 - CRUD for data stored in a single location
* [starport scaffold type](#starport-scaffold-type)	 - Scaffold only a type definition
* [starport scaffold vue](#starport-scaffold-vue)	 - Vue 3 web app template


## starport scaffold band

Scaffold an IBC BandChain query oracle to request real-time data

**Synopsis**

Scaffold an IBC BandChain query oracle to request real-time data from BandChain scripts in a specific IBC-enabled Cosmos SDK module

```
starport scaffold band [queryName] --module [moduleName] [flags]
```

**Options**

```
  -h, --help            help for band
      --module string   IBC Module to add the packet into
  -p, --path string     path of the app (default ".")
      --signer string   Label for the message signer (default: creator)
```

**SEE ALSO**

* [starport scaffold](#starport-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## starport scaffold chain

Fully-featured Cosmos SDK blockchain

**Synopsis**

Scaffold a new Cosmos SDK blockchain with a default directory structure

```
starport scaffold chain [github.com/org/repo] [flags]
```

**Options**

```
      --address-prefix string   Address prefix (default "cosmos")
  -h, --help                    help for chain
      --no-module               Prevent scaffolding a default module in the app
  -p, --path string             path to scaffold the chain (default ".")
```

**SEE ALSO**

* [starport scaffold](#starport-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## starport scaffold list

CRUD for data stored as an array

```
starport scaffold list NAME [field]... [flags]
```

**Options**

```
  -h, --help            help for list
      --module string   Module to add into. Default is app's main module
      --no-message      Disable CRUD interaction messages scaffolding
  -p, --path string     path of the app (default ".")
      --signer string   Label for the message signer (default: creator)
```

**SEE ALSO**

* [starport scaffold](#starport-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## starport scaffold map

CRUD for data stored as key-value pairs

```
starport scaffold map NAME [field]... [flags]
```

**Options**

```
  -h, --help            help for map
      --index strings   fields that index the value (default [index])
      --module string   Module to add into. Default is app's main module
      --no-message      Disable CRUD interaction messages scaffolding
  -p, --path string     path of the app (default ".")
      --signer string   Label for the message signer (default: creator)
```

**SEE ALSO**

* [starport scaffold](#starport-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## starport scaffold message

Message to perform state transition on the blockchain

```
starport scaffold message [name] [field1] [field2] ... [flags]
```

**Options**

```
  -d, --desc string        Description of the command
  -h, --help               help for message
      --module string      Module to add the message into. Default: app's main module
  -p, --path string        path of the app (default ".")
  -r, --response strings   Response fields
      --signer string      Label for the message signer (default: creator)
```

**SEE ALSO**

* [starport scaffold](#starport-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## starport scaffold module

Scaffold a Cosmos SDK module

**Synopsis**

Scaffold a new Cosmos SDK module in the `x` directory

```
starport scaffold module [name] [flags]
```

**Options**

```
      --dep strings            module dependencies (e.g. --dep account,bank)
  -h, --help                   help for module
      --ibc                    scaffold an IBC module
      --ordering string        channel ordering of the IBC module [none|ordered|unordered] (default "none")
  -p, --path string            path of the app (default ".")
      --require-registration   if true command will fail if module can't be registered
```

**SEE ALSO**

* [starport scaffold](#starport-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## starport scaffold packet

Message for sending an IBC packet

**Synopsis**

Scaffold an IBC packet in a specific IBC-enabled Cosmos SDK module

```
starport scaffold packet [packetName] [field1] [field2] ... --module [moduleName] [flags]
```

**Options**

```
      --ack strings     Custom acknowledgment type (field1,field2,...)
  -h, --help            help for packet
      --module string   IBC Module to add the packet into
      --no-message      Disable send message scaffolding
  -p, --path string     path of the app (default ".")
      --signer string   Label for the message signer (default: creator)
```

**SEE ALSO**

* [starport scaffold](#starport-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## starport scaffold query

Query to get data from the blockchain

```
starport scaffold query [name] [request_field1] [request_field2] ... [flags]
```

**Options**

```
  -d, --desc string        Description of the command
  -h, --help               help for query
      --module string      Module to add the query into. Default: app's main module
      --paginated          Define if the request can be paginated
  -p, --path string        path of the app (default ".")
  -r, --response strings   Response fields
```

**SEE ALSO**

* [starport scaffold](#starport-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## starport scaffold single

CRUD for data stored in a single location

```
starport scaffold single NAME [field]... [flags]
```

**Options**

```
  -h, --help            help for single
      --module string   Module to add into. Default is app's main module
      --no-message      Disable CRUD interaction messages scaffolding
  -p, --path string     path of the app (default ".")
      --signer string   Label for the message signer (default: creator)
```

**SEE ALSO**

* [starport scaffold](#starport-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## starport scaffold type

Scaffold only a type definition

```
starport scaffold type NAME [field]... [flags]
```

**Options**

```
  -h, --help            help for type
      --module string   Module to add into. Default is app's main module
      --no-message      Disable CRUD interaction messages scaffolding
  -p, --path string     path of the app (default ".")
      --signer string   Label for the message signer (default: creator)
```

**SEE ALSO**

* [starport scaffold](#starport-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## starport scaffold vue

Vue 3 web app template

```
starport scaffold vue [flags]
```

**Options**

```
  -h, --help          help for vue
  -p, --path string   path to scaffold content of the Vue.js app (default "./vue")
```

**SEE ALSO**

* [starport scaffold](#starport-scaffold)	 - Scaffold a new blockchain, module, message, query, and more


## starport tools

Tools for advanced users

**Options**

```
  -h, --help   help for tools
```

**SEE ALSO**

* [starport](#starport)	 - Starport offers everything you need to scaffold, test, build, and launch your blockchain
* [starport tools completions](#starport-tools-completions)	 - Generate completions script
* [starport tools ibc-relayer](#starport-tools-ibc-relayer)	 - Typescript implementation of an IBC relayer
* [starport tools ibc-setup](#starport-tools-ibc-setup)	 - Collection of commands to quickly setup a relayer
* [starport tools protoc](#starport-tools-protoc)	 - Execute the protoc command


## starport tools completions

Generate completions script

**Synopsis**

 The completions command outputs a completion script you can use in your shell. The output script requires 
				that [bash-completion](https://github.com/scop/bash-completion)	is installed and enabled in your 
				system. Since most Unix-like operating systems come with bash-completion by default, bash-completion 
				is probably already installed and operational.

Bash:

  $ source <(starport  tools completions bash)

  To load completions for every new session, run:

  ** Linux **
  $ starport  tools completions bash > /etc/bash_completion.d/starport

  ** macOS **
  $ starport  tools completions bash > /usr/local/etc/bash_completion.d/starport

Zsh:

  If shell completions is not already enabled in your environment, you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  To load completions for each session, execute once:
  
  $ starport  tools completions zsh > "${fpath[1]}/_starport"

  You will need to start a new shell for this setup to take effect.

fish:

  $ starport  tools completions fish | source

  To load completions for each session, execute once:
  
  $ starport  tools completions fish > ~/.config/fish/completionss/starport.fish

PowerShell:

  PS> starport  tools completions powershell | Out-String | Invoke-Expression

  To load completions for every new session, run:
  
  PS> starport  tools completions powershell > starport.ps1
  
  and source this file from your PowerShell profile.


```
starport tools completions
```

**Options**

```
  -h, --help   help for completions
```

**SEE ALSO**

* [starport tools](#starport-tools)	 - Tools for advanced users


## starport tools ibc-relayer

Typescript implementation of an IBC relayer

```
starport tools ibc-relayer [--] [...] [flags]
```

**Examples**

```
starport tools ibc-relayer -- -h
```

**Options**

```
  -h, --help   help for ibc-relayer
```

**SEE ALSO**

* [starport tools](#starport-tools)	 - Tools for advanced users


## starport tools ibc-setup

Collection of commands to quickly setup a relayer

```
starport tools ibc-setup [--] [...] [flags]
```

**Examples**

```
starport tools ibc-setup -- -h
starport tools ibc-setup -- init --src relayer_test_1 --dest relayer_test_2
```

**Options**

```
  -h, --help   help for ibc-setup
```

**SEE ALSO**

* [starport tools](#starport-tools)	 - Tools for advanced users


## starport tools protoc

Execute the protoc command

**Synopsis**

The protoc command. You don't need to setup the global protoc include folder with -I, it's automatically handled

```
starport tools protoc [--] [...] [flags]
```

**Examples**

```
starport tools protoc -- --version
```

**Options**

```
  -h, --help   help for protoc
```

**SEE ALSO**

* [starport tools](#starport-tools)	 - Tools for advanced users


## starport version

Print the current build information

```
starport version [flags]
```

**Options**

```
  -h, --help   help for version
```

**SEE ALSO**

* [starport](#starport)	 - Starport offers everything you need to scaffold, test, build, and launch your blockchain

