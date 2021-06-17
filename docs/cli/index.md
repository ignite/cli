---
order: 1
description: Starport CLI docs. 
parent:
  order: 8
  title: CLI 
---

# CLI
Documentation for Starport CLI.

## starport

A developer tool for building Cosmos SDK blockchains

**Options**

```
  -h, --help   help for starport
```

**SEE ALSO**

* [starport app](#starport-app)	 - Scaffold a new blockchain
* [starport build](#starport-build)	 - Build a node binary
* [starport docs](#starport-docs)	 - Show Starport docs
* [starport faucet](#starport-faucet)	 - Send coins to an account
* [starport message](#starport-message)	 - Scaffold a Cosmos SDK message
* [starport module](#starport-module)	 - Manage Cosmos SDK modules for your blockchain
* [starport network](#starport-network)	 - Launch a blockchain network in a decentralized way
* [starport packet](#starport-packet)	 - Scaffold an IBC packet
* [starport query](#starport-query)	 - Scaffold a Cosmos SDK query
* [starport relayer](#starport-relayer)	 - Connects blockchains via IBC protocol
* [starport serve](#starport-serve)	 - Start a blockchain node in development
* [starport type](#starport-type)	 - Scaffold a type with CRUD operations
* [starport version](#starport-version)	 - Print the current build information


## starport app

Scaffold a new blockchain

**Synopsis**

Scaffold a new Cosmos SDK blockchain with a default directory structure

```
starport app [github.com/org/repo] [flags]
```

**Options**

```
      --address-prefix string   Address prefix (default "cosmos")
  -h, --help                    help for app
      --no-default-module       Prevent scaffolding a default module in the app
```

**SEE ALSO**

* [starport](#starport)	 - A developer tool for building Cosmos SDK blockchains


## starport build

Build a node binary

**Synopsis**

By default, build your node binaries
and add the binaries to your $(go env GOPATH)/bin path.

To build binaries for a release, use the --release flag. The app binaries
for one or more specified release targets are built in a release/ dir under the app's
source. Specify the release targets with GOOS:GOARCH build tags.
If the optional --release.targets is not specified, a binary is created for your current environment.

Sample usages:
	- starport build
	- starport build --release -t linux:amd64 -t darwin:amd64 -t darwin:arm64

```
starport build [flags]
```

**Options**

```
  -h, --help                      help for build
      --home string               Home directory used for blockchains
  -p, --path string               path of the app (default ".")
      --rebuild-proto-once        Enables proto code generation for 3rd party modules. Available only without the --release flag
      --release                   build for a release
      --release.prefix string     tarball prefix for each release target. Available only with --release flag
  -t, --release.targets strings   release targets. Available only with --release flag
  -v, --verbose                   Verbose output
```

**SEE ALSO**

* [starport](#starport)	 - A developer tool for building Cosmos SDK blockchains


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

* [starport](#starport)	 - A developer tool for building Cosmos SDK blockchains


## starport faucet

Send coins to an account

```
starport faucet [address] [coin<,...>] [flags]
```

**Options**

```
  -h, --help          help for faucet
      --home string   Home directory used for blockchains
  -p, --path string   path of the app
  -v, --verbose       Verbose output
```

**SEE ALSO**

* [starport](#starport)	 - A developer tool for building Cosmos SDK blockchains


## starport message

Scaffold a Cosmos SDK message

```
starport message [name] [field1] [field2] ... [flags]
```

**Options**

```
  -d, --desc string        Description of the command
  -h, --help               help for message
      --module string      Module to add the message into. Default: app's main module
  -p, --path string        path of the app
  -r, --response strings   Response fields
```

**SEE ALSO**

* [starport](#starport)	 - A developer tool for building Cosmos SDK blockchains


## starport module

Manage Cosmos SDK modules for your blockchain

**Options**

```
  -h, --help   help for module
```

**SEE ALSO**

* [starport](#starport)	 - A developer tool for building Cosmos SDK blockchains
* [starport module create](#starport-module-create)	 - Scaffold a Cosmos SDK module
* [starport module import](#starport-module-import)	 - Import a new module to app.


## starport module create

Scaffold a Cosmos SDK module

**Synopsis**

Scaffold a new Cosmos SDK module in the `x` directory

```
starport module create [name] [flags]
```

**Options**

```
  -h, --help                   help for create
      --ibc                    scaffold an IBC module
      --ordering string        channel ordering of the IBC module [none|ordered|unordered] (default "none")
      --require-registration   if true command will fail if module can't be registered
```

**SEE ALSO**

* [starport module](#starport-module)	 - Manage Cosmos SDK modules for your blockchain


## starport module import

Import a new module to app.

**Synopsis**

Add support for WebAssembly smart contracts to your blockchain.

```
starport module import [feature] [flags]
```

**Options**

```
  -h, --help   help for import
```

**SEE ALSO**

* [starport module](#starport-module)	 - Manage Cosmos SDK modules for your blockchain


## starport network

Launch a blockchain network in a decentralized way

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

* [starport](#starport)	 - A developer tool for building Cosmos SDK blockchains
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

* [starport network](#starport-network)	 - Launch a blockchain network in a decentralized way
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

* [starport network](#starport-network)	 - Launch a blockchain network in a decentralized way
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

* [starport network](#starport-network)	 - Launch a blockchain network in a decentralized way
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


## starport packet

Scaffold an IBC packet

**Synopsis**

Scaffold an IBC packet in a specific IBC-enabled Cosmos SDK module

```
starport packet [packetName] [field1] [field2] ... --module [moduleName] [flags]
```

**Options**

```
      --ack strings     Custom acknowledgment type (field1,field2,...)
  -h, --help            help for packet
      --module string   IBC Module to add the packet into
      --no-message      Disable send message scaffolding
```

**SEE ALSO**

* [starport](#starport)	 - A developer tool for building Cosmos SDK blockchains


## starport query

Scaffold a Cosmos SDK query

```
starport query [name] [request_field1] [request_field2] ... [flags]
```

**Options**

```
  -d, --desc string        Description of the command
  -h, --help               help for query
      --module string      Module to add the query into. Default: app's main module
      --paginated          Define if the request can be paginated
  -p, --path string        path of the app
  -r, --response strings   Response fields
```

**SEE ALSO**

* [starport](#starport)	 - A developer tool for building Cosmos SDK blockchains


## starport relayer

Connects blockchains via IBC protocol

**Options**

```
  -h, --help   help for relayer
```

**SEE ALSO**

* [starport](#starport)	 - A developer tool for building Cosmos SDK blockchains
* [starport relayer configure](#starport-relayer-configure)	 - Configure source and target chains for relaying
* [starport relayer connect](#starport-relayer-connect)	 - Link chains associated with paths and start relaying tx packets in between
* [starport relayer lowlevel](#starport-relayer-lowlevel)	 - Low-level relayer commands from @confio/relayer
* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer


## starport relayer configure

Configure source and target chains for relaying

```
starport relayer configure [flags]
```

**Options**

```
  -a, --advanced                 Advanced configuration options for custom IBC modules
  -h, --help                     help for configure
      --ordered                  Set the channel as ordered
      --source-faucet string     Faucet address of the source chain
      --source-gasprice string   Gas price used for transactions on source chain
      --source-port string       IBC port ID on the source chain
      --source-prefix string     Address prefix of the source chain
      --source-rpc string        RPC address of the source chain
      --source-version string    Module version on the source chain
      --target-faucet string     Faucet address of the target chain
      --target-gasprice string   Gas price used for transactions on target chain
      --target-port string       IBC port ID on the target chain
      --target-prefix string     Address prefix of the target chain
      --target-rpc string        RPC address of the target chain
      --target-version string    Module version on the target chain
```

**SEE ALSO**

* [starport relayer](#starport-relayer)	 - Connects blockchains via IBC protocol


## starport relayer connect

Link chains associated with paths and start relaying tx packets in between

```
starport relayer connect [<path>,...] [flags]
```

**Options**

```
  -h, --help   help for connect
```

**SEE ALSO**

* [starport relayer](#starport-relayer)	 - Connects blockchains via IBC protocol


## starport relayer lowlevel

Low-level relayer commands from @confio/relayer

**Options**

```
  -h, --help   help for lowlevel
```

**SEE ALSO**

* [starport relayer](#starport-relayer)	 - Connects blockchains via IBC protocol
* [starport relayer lowlevel ibc-relayer](#starport-relayer-lowlevel-ibc-relayer)	 - Typescript implementation of an IBC relayer
* [starport relayer lowlevel ibc-setup](#starport-relayer-lowlevel-ibc-setup)	 - Collection of commands to quickly setup a relayer


## starport relayer lowlevel ibc-relayer

Typescript implementation of an IBC relayer

```
starport relayer lowlevel ibc-relayer [--] [...] [flags]
```

**Examples**

```
starport relayer lowlevel ibc-relayer -- -h
```

**Options**

```
  -h, --help   help for ibc-relayer
```

**SEE ALSO**

* [starport relayer lowlevel](#starport-relayer-lowlevel)	 - Low-level relayer commands from @confio/relayer


## starport relayer lowlevel ibc-setup

Collection of commands to quickly setup a relayer

```
starport relayer lowlevel ibc-setup [--] [...] [flags]
```

**Examples**

```
starport relayer lowlevel ibc-setup -- -h
starport relayer lowlevel ibc-setup -- init --src relayer_test_1 --dest relayer_test_2
```

**Options**

```
  -h, --help   help for ibc-setup
```

**SEE ALSO**

* [starport relayer lowlevel](#starport-relayer-lowlevel)	 - Low-level relayer commands from @confio/relayer


## starport relayer rly

Low-level commands from github.com/cosmos/relayer

**Synopsis**

The relayer has commands for:
  1. Configuration of the Chains and Paths that the relayer with transfer packets over
  2. Management of keys and light clients on the local machine that will be used to sign and verify txs
  3. Query and transaction functionality for IBC
  4. A responsive relaying application that listens on a path
  5. Commands to assist with development, testnets, and versioning.

NOTE: Most of the commands have aliases that make typing them much quicker (i.e. 'rly tx', 'rly q', etc...)

**Options**

```
  -d, --debug         debug output
  -h, --help          help for rly
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer](#starport-relayer)	 - Connects blockchains via IBC protocol
* [starport relayer rly ](#starport-relayer-rly-)	 - 
* [starport relayer rly ](#starport-relayer-rly-)	 - 
* [starport relayer rly ](#starport-relayer-rly-)	 - 
* [starport relayer rly ](#starport-relayer-rly-)	 - 
* [starport relayer rly api](#starport-relayer-rly-api)	 - Start the relayer API
* [starport relayer rly chains](#starport-relayer-rly-chains)	 - manage chain configurations
* [starport relayer rly config](#starport-relayer-rly-config)	 - manage configuration file
* [starport relayer rly development](#starport-relayer-rly-development)	 - commands for developers either deploying or hacking on the relayer
* [starport relayer rly keys](#starport-relayer-rly-keys)	 - manage keys held by the relayer for each chain
* [starport relayer rly light](#starport-relayer-rly-light)	 - manage light clients held by the relayer for each chain
* [starport relayer rly paths](#starport-relayer-rly-paths)	 - manage path configurations
* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands
* [starport relayer rly start](#starport-relayer-rly-start)	 - Start the listening relayer on a given path
* [starport relayer rly testnets](#starport-relayer-rly-testnets)	 - commands for joining and running relayer testnets
* [starport relayer rly transact](#starport-relayer-rly-transact)	 - IBC transaction commands
* [starport relayer rly version](#starport-relayer-rly-version)	 - Print the relayer version info


## starport relayer rly 



```
starport relayer rly  [flags]
```

**Options**

```
  -h, --help   help for this command
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer


## starport relayer rly 



```
starport relayer rly  [flags]
```

**Options**

```
  -h, --help   help for this command
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer


## starport relayer rly 



```
starport relayer rly  [flags]
```

**Options**

```
  -h, --help   help for this command
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer


## starport relayer rly 



```
starport relayer rly  [flags]
```

**Options**

```
  -h, --help   help for this command
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer


## starport relayer rly api

Start the relayer API

```
starport relayer rly api [flags]
```

**Options**

```
  -h, --help   help for api
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer


## starport relayer rly chains

manage chain configurations

**Options**

```
  -h, --help   help for chains
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer
* [starport relayer rly chains add](#starport-relayer-rly-chains-add)	 - Add a new chain to the configuration file by passing a file (-f) or url (-u), or user input
* [starport relayer rly chains add-dir](#starport-relayer-rly-chains-add-dir)	 - Add new chains to the configuration file from a directory 
		full of chain configuration, useful for adding testnet configurations
* [starport relayer rly chains address](#starport-relayer-rly-chains-address)	 - Returns a chain's configured key's address
* [starport relayer rly chains delete](#starport-relayer-rly-chains-delete)	 - Returns chain configuration data
* [starport relayer rly chains edit](#starport-relayer-rly-chains-edit)	 - Returns chain configuration data
* [starport relayer rly chains list](#starport-relayer-rly-chains-list)	 - Returns chain configuration data
* [starport relayer rly chains show](#starport-relayer-rly-chains-show)	 - Returns a chain's configuration data


## starport relayer rly chains add

Add a new chain to the configuration file by passing a file (-f) or url (-u), or user input

```
starport relayer rly chains add [flags]
```

**Examples**

```
$ rly chains add
$ rly ch a
$ rly chains add --file chains/ibc0.json
$ rly chains add --url http://relayer.com/ibc0.json
```

**Options**

```
  -f, --file string   fetch json data from specified file
  -h, --help          help for add
  -u, --url string    url to fetch data from
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly chains](#starport-relayer-rly-chains)	 - manage chain configurations


## starport relayer rly chains add-dir

Add new chains to the configuration file from a directory 
		full of chain configuration, useful for adding testnet configurations

```
starport relayer rly chains add-dir [dir] [flags]
```

**Examples**

```
$ rly chains add-dir testnet/chains/
$ rly ch ad testnet/chains/
```

**Options**

```
  -h, --help   help for add-dir
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly chains](#starport-relayer-rly-chains)	 - manage chain configurations


## starport relayer rly chains address

Returns a chain's configured key's address

```
starport relayer rly chains address [chain-id] [flags]
```

**Examples**

```
$ rly chains address ibc-0
$ rly ch addr ibc-0
```

**Options**

```
  -h, --help   help for address
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly chains](#starport-relayer-rly-chains)	 - manage chain configurations


## starport relayer rly chains delete

Returns chain configuration data

```
starport relayer rly chains delete [chain-id] [flags]
```

**Examples**

```
$ rly chains delete ibc-0
$ rly ch d ibc-0
```

**Options**

```
  -h, --help   help for delete
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly chains](#starport-relayer-rly-chains)	 - manage chain configurations


## starport relayer rly chains edit

Returns chain configuration data

```
starport relayer rly chains edit [chain-id] [key] [value] [flags]
```

**Examples**

```
$ rly chains edit ibc-0 trusting-period 32h
$ rly ch e ibc-0 trusting-period 32h
```

**Options**

```
  -h, --help   help for edit
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly chains](#starport-relayer-rly-chains)	 - manage chain configurations


## starport relayer rly chains list

Returns chain configuration data

```
starport relayer rly chains list [flags]
```

**Examples**

```
$ rly chains list
$ rly ch l
```

**Options**

```
  -h, --help   help for list
  -j, --json   returns the response in json format
  -y, --yaml   output using yaml
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly chains](#starport-relayer-rly-chains)	 - manage chain configurations


## starport relayer rly chains show

Returns a chain's configuration data

```
starport relayer rly chains show [chain-id] [flags]
```

**Examples**

```
$ rly chains show ibc-0 --json
$ rly chains show ibc-0 --yaml
$ rly ch s ibc-0 --json
$ rly ch s ibc-0 --yaml
```

**Options**

```
  -h, --help   help for show
  -j, --json   returns the response in json format
  -y, --yaml   output using yaml
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly chains](#starport-relayer-rly-chains)	 - manage chain configurations


## starport relayer rly config

manage configuration file

**Options**

```
  -h, --help   help for config
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer
* [starport relayer rly config add-chains](#starport-relayer-rly-config-add-chains)	 - Add new chains to the configuration file from a
		 directory full of chain configurations, useful for adding testnet configurations
* [starport relayer rly config add-paths](#starport-relayer-rly-config-add-paths)	 - Add new paths to the configuration file from a directory full of path configurations, useful for adding testnet configurations. 
		Chain configuration files must be added before calling this command.
* [starport relayer rly config init](#starport-relayer-rly-config-init)	 - Creates a default home directory at path defined by --home
* [starport relayer rly config show](#starport-relayer-rly-config-show)	 - Prints current configuration


## starport relayer rly config add-chains

Add new chains to the configuration file from a
		 directory full of chain configurations, useful for adding testnet configurations

```
starport relayer rly config add-chains [/path/to/chains/] [flags]
```

**Examples**

```
$ rly config add-chains configs/chains
```

**Options**

```
  -h, --help   help for add-chains
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly config](#starport-relayer-rly-config)	 - manage configuration file


## starport relayer rly config add-paths

Add new paths to the configuration file from a directory full of path configurations, useful for adding testnet configurations. 
		Chain configuration files must be added before calling this command.

```
starport relayer rly config add-paths [/path/to/paths/] [flags]
```

**Examples**

```
$ rly config add-paths configs/paths
```

**Options**

```
  -h, --help   help for add-paths
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly config](#starport-relayer-rly-config)	 - manage configuration file


## starport relayer rly config init

Creates a default home directory at path defined by --home

```
starport relayer rly config init [flags]
```

**Examples**

```
$ rly config init --home /home/runner/.relayer
$ rly cfg i
```

**Options**

```
  -h, --help   help for init
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly config](#starport-relayer-rly-config)	 - manage configuration file


## starport relayer rly config show

Prints current configuration

```
starport relayer rly config show [flags]
```

**Examples**

```
$ rly config show --home /home/runner/.relayer
$ rly cfg list
```

**Options**

```
  -h, --help   help for show
  -j, --json   returns the response in json format
  -y, --yaml   output using yaml
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly config](#starport-relayer-rly-config)	 - manage configuration file


## starport relayer rly development

commands for developers either deploying or hacking on the relayer

**Options**

```
  -h, --help   help for development
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer
* [starport relayer rly development faucet](#starport-relayer-rly-development-faucet)	 - faucet returns a sample faucet service file
* [starport relayer rly development gaia](#starport-relayer-rly-development-gaia)	 - gaia returns a sample gaiad service file
* [starport relayer rly development genesis](#starport-relayer-rly-development-genesis)	 - fetch the genesis file for a configured chain
* [starport relayer rly development listen](#starport-relayer-rly-development-listen)	 - listen to all transaction and block events from a given chain and output them to stdout
* [starport relayer rly development relayer](#starport-relayer-rly-development-relayer)	 - relayer returns a service file for the relayer to relay over an individual path


## starport relayer rly development faucet

faucet returns a sample faucet service file

```
starport relayer rly development faucet [user] [home] [chain-id] [key-name] [amount] [flags]
```

**Examples**

```
$ rly dev faucet faucetuser /home/faucetuser ibc-0 testkey 1000000stake
$ rly development faucet root /home/root ibc-1 testkey2 1000000stake
```

**Options**

```
  -h, --help   help for faucet
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly development](#starport-relayer-rly-development)	 - commands for developers either deploying or hacking on the relayer


## starport relayer rly development gaia

gaia returns a sample gaiad service file

```
starport relayer rly development gaia [user] [home] [flags]
```

**Examples**

```
$ rly dev gaia [user-name] [path-to-home]
$ rly development gaia user /home/user
```

**Options**

```
  -h, --help   help for gaia
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly development](#starport-relayer-rly-development)	 - commands for developers either deploying or hacking on the relayer


## starport relayer rly development genesis

fetch the genesis file for a configured chain

```
starport relayer rly development genesis [chain-id] [flags]
```

**Examples**

```
$ rly dev genesis ibc-0
$ rly dev gen ibc-0
$ rly development genesis ibc-2
```

**Options**

```
  -h, --help   help for genesis
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly development](#starport-relayer-rly-development)	 - commands for developers either deploying or hacking on the relayer


## starport relayer rly development listen

listen to all transaction and block events from a given chain and output them to stdout

```
starport relayer rly development listen [chain-id] [flags]
```

**Examples**

```
$ rly dev listen ibc-0 --data --no-tx
$ rly dev l ibc-1 --no-block
$ rly development listen ibc-2 --no-tx
```

**Options**

```
      --data       output full event data
  -h, --help       help for listen
  -b, --no-block   don't output block events
  -t, --no-tx      don't output transaction events
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly development](#starport-relayer-rly-development)	 - commands for developers either deploying or hacking on the relayer


## starport relayer rly development relayer

relayer returns a service file for the relayer to relay over an individual path

```
starport relayer rly development relayer [path-name] [flags]
```

**Examples**

```
$ rly dev rly demo-path
$ rly development relayer path-test
```

**Options**

```
  -h, --help   help for relayer
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly development](#starport-relayer-rly-development)	 - commands for developers either deploying or hacking on the relayer


## starport relayer rly keys

manage keys held by the relayer for each chain

**Options**

```
  -h, --help   help for keys
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer
* [starport relayer rly keys add](#starport-relayer-rly-keys-add)	 - adds a key to the keychain associated with a particular chain
* [starport relayer rly keys delete](#starport-relayer-rly-keys-delete)	 - deletes a key from the keychain associated with a particular chain
* [starport relayer rly keys export](#starport-relayer-rly-keys-export)	 - exports a privkey from the keychain associated with a particular chain
* [starport relayer rly keys list](#starport-relayer-rly-keys-list)	 - lists keys from the keychain associated with a particular chain
* [starport relayer rly keys restore](#starport-relayer-rly-keys-restore)	 - restores a mnemonic to the keychain associated with a particular chain
* [starport relayer rly keys show](#starport-relayer-rly-keys-show)	 - shows a key from the keychain associated with a particular chain


## starport relayer rly keys add

adds a key to the keychain associated with a particular chain

```
starport relayer rly keys add [chain-id] [[name]] [flags]
```

**Examples**

```
$ rly keys add ibc-0
$ rly keys add ibc-1 key2
$ rly k a ibc-2 testkey
```

**Options**

```
      --coin-type uint32   coin type number for HD derivation (default 118)
  -h, --help               help for add
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly keys](#starport-relayer-rly-keys)	 - manage keys held by the relayer for each chain


## starport relayer rly keys delete

deletes a key from the keychain associated with a particular chain

```
starport relayer rly keys delete [chain-id] [[name]] [flags]
```

**Examples**

```
$ rly keys delete ibc-0 -y
$ rly keys delete ibc-1 key2 -y
$ rly k d ibc-2 testkey
```

**Options**

```
  -h, --help   help for delete
  -y, --skip   output using yaml
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly keys](#starport-relayer-rly-keys)	 - manage keys held by the relayer for each chain


## starport relayer rly keys export

exports a privkey from the keychain associated with a particular chain

```
starport relayer rly keys export [chain-id] [name] [flags]
```

**Examples**

```
$ rly keys export ibc-0 testkey
$ rly k e ibc-2 testkey
```

**Options**

```
  -h, --help   help for export
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly keys](#starport-relayer-rly-keys)	 - manage keys held by the relayer for each chain


## starport relayer rly keys list

lists keys from the keychain associated with a particular chain

```
starport relayer rly keys list [chain-id] [flags]
```

**Examples**

```
$ rly keys list ibc-0
$ rly k l ibc-1
```

**Options**

```
  -h, --help   help for list
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly keys](#starport-relayer-rly-keys)	 - manage keys held by the relayer for each chain


## starport relayer rly keys restore

restores a mnemonic to the keychain associated with a particular chain

```
starport relayer rly keys restore [chain-id] [name] [mnemonic] [flags]
```

**Examples**

```
$ rly keys restore ibc-0 testkey "[mnemonic-words]"
$ rly k r ibc-1 faucet-key "[mnemonic-words]"
```

**Options**

```
      --coin-type uint32   coin type number for HD derivation (default 118)
  -h, --help               help for restore
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly keys](#starport-relayer-rly-keys)	 - manage keys held by the relayer for each chain


## starport relayer rly keys show

shows a key from the keychain associated with a particular chain

```
starport relayer rly keys show [chain-id] [[name]] [flags]
```

**Examples**

```
$ rly keys show ibc-0
$ rly keys show ibc-1 key2
$ rly k s ibc-2 testkey
```

**Options**

```
  -h, --help   help for show
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly keys](#starport-relayer-rly-keys)	 - manage keys held by the relayer for each chain


## starport relayer rly light

manage light clients held by the relayer for each chain

**Options**

```
  -h, --help   help for light
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer
* [starport relayer rly light delete](#starport-relayer-rly-light-delete)	 - wipe the light client database, forcing re-initialzation on the next run
* [starport relayer rly light header](#starport-relayer-rly-light-header)	 - Get a header from the light client database
* [starport relayer rly light init](#starport-relayer-rly-light-init)	 - Initiate the light client
* [starport relayer rly light update](#starport-relayer-rly-light-update)	 - Update the light client to latest header from configured node


## starport relayer rly light delete

wipe the light client database, forcing re-initialzation on the next run

```
starport relayer rly light delete [chain-id] [flags]
```

**Examples**

```
$ rly light delete ibc-0
$ rly l d ibc-2
```

**Options**

```
  -h, --help   help for delete
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly light](#starport-relayer-rly-light)	 - manage light clients held by the relayer for each chain


## starport relayer rly light header

Get a header from the light client database

**Synopsis**

Get a header from the light client database. 0 returns lasttrusted header and all others return the header at that height if stored

```
starport relayer rly light header [chain-id] [[height]] [flags]
```

**Examples**

```
$ rly light header ibc-0
$ rly light header ibc-1 1400
$ rly l hdr ibc-2
```

**Options**

```
  -h, --help   help for header
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly light](#starport-relayer-rly-light)	 - manage light clients held by the relayer for each chain


## starport relayer rly light init

Initiate the light client

**Synopsis**

Initiate the light client by:
	1. passing it a root of trust as a --hash/-x and --height
	2. Use --force/-f to initialize from the configured node

```
starport relayer rly light init [chain-id] [flags]
```

**Examples**

```
$ rly light init ibc-0 --force
$ rly light init ibc-1 --height 1406 --hash <hash>
$ rly l i ibc-2 --force
```

**Options**

```
  -f, --force           option to force non-standard behavior such as initialization of light client fromconfigured chain or generation of new path
  -x, --hash bytesHex   Trusted header's hash
      --height int      Trusted header's height (default -1)
  -h, --help            help for init
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly light](#starport-relayer-rly-light)	 - manage light clients held by the relayer for each chain


## starport relayer rly light update

Update the light client to latest header from configured node

```
starport relayer rly light update [chain-id] [flags]
```

**Examples**

```
$ rly light update ibc-0
$ rly l u ibc-1
```

**Options**

```
  -h, --help   help for update
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly light](#starport-relayer-rly-light)	 - manage light clients held by the relayer for each chain


## starport relayer rly paths

manage path configurations

**Synopsis**


A path represents the "full path" or "link" for communication between two chains. This includes the client, 
connection, and channel ids from both the source and destination chains as well as the strategy to use when relaying

**Options**

```
  -h, --help   help for paths
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer
* [starport relayer rly paths add](#starport-relayer-rly-paths-add)	 - add a path to the list of paths
* [starport relayer rly paths delete](#starport-relayer-rly-paths-delete)	 - delete a path with a given index
* [starport relayer rly paths generate](#starport-relayer-rly-paths-generate)	 - generate a new path between src and dst, reusing any identifiers that exist
* [starport relayer rly paths list](#starport-relayer-rly-paths-list)	 - print out configured paths
* [starport relayer rly paths show](#starport-relayer-rly-paths-show)	 - show a path given its name


## starport relayer rly paths add

add a path to the list of paths

```
starport relayer rly paths add [src-chain-id] [dst-chain-id] [path-name] [flags]
```

**Examples**

```
$ rly paths add ibc-0 ibc-1 demo-path
$ rly paths add ibc-0 ibc-1 demo-path --file paths/demo.json
$ rly pth a ibc-0 ibc-1 demo-path
```

**Options**

```
  -f, --file string   fetch json data from specified file
  -h, --help          help for add
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly paths](#starport-relayer-rly-paths)	 - manage path configurations


## starport relayer rly paths delete

delete a path with a given index

```
starport relayer rly paths delete [index] [flags]
```

**Examples**

```
$ rly paths delete demo-path
$ rly pth d path-name
```

**Options**

```
  -h, --help   help for delete
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly paths](#starport-relayer-rly-paths)	 - manage path configurations


## starport relayer rly paths generate

generate a new path between src and dst, reusing any identifiers that exist

```
starport relayer rly paths generate [src-chain-id] [dst-chain-id] [name] [flags]
```

**Examples**

```
$ rly paths generate ibc-0 ibc-1 demo-path
$ rly pth gen ibc-0 ibc-1 demo-path --unordered false --version ics20-2
```

**Options**

```
  -h, --help              help for generate
  -p, --port string       port to use when generating path (default "transfer")
  -s, --strategy string   specify strategy of path to generate (default "naive")
  -o, --unordered         create an unordered channel (default true)
  -v, --version string    version of channel to create (default "ics20-1")
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly paths](#starport-relayer-rly-paths)	 - manage path configurations


## starport relayer rly paths list

print out configured paths

```
starport relayer rly paths list [flags]
```

**Examples**

```
$ rly paths list --yaml
$ rly paths list --json
$ rly pth l
```

**Options**

```
  -h, --help   help for list
  -j, --json   returns the response in json format
  -y, --yaml   output using yaml
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly paths](#starport-relayer-rly-paths)	 - manage path configurations


## starport relayer rly paths show

show a path given its name

```
starport relayer rly paths show [path-name] [flags]
```

**Examples**

```
$ rly paths show demo-path --yaml
$ rly paths show demo-path --json
$ rly pth s path-name
```

**Options**

```
  -h, --help   help for show
  -j, --json   returns the response in json format
  -y, --yaml   output using yaml
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly paths](#starport-relayer-rly-paths)	 - manage path configurations


## starport relayer rly query

IBC query commands

**Synopsis**

Commands to query IBC primitives and other useful data on configured chains.

**Options**

```
  -h, --help   help for query
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer
* [starport relayer rly query ](#starport-relayer-rly-query-)	 - 
* [starport relayer rly query ](#starport-relayer-rly-query-)	 - 
* [starport relayer rly query account](#starport-relayer-rly-query-account)	 - query the relayer's account on a given network by chain ID
* [starport relayer rly query balance](#starport-relayer-rly-query-balance)	 - query the relayer's account balance on a given network by chain-ID
* [starport relayer rly query channel](#starport-relayer-rly-query-channel)	 - query a channel by channel and port ID on a network by chain ID
* [starport relayer rly query channels](#starport-relayer-rly-query-channels)	 - query for all channels on a network by chain ID
* [starport relayer rly query client](#starport-relayer-rly-query-client)	 - query the state of a light client on a network by chain ID
* [starport relayer rly query client-connections](#starport-relayer-rly-query-client-connections)	 - query for all connections for a given client on a network by chain ID
* [starport relayer rly query clients](#starport-relayer-rly-query-clients)	 - query for all light client states on a network by chain ID
* [starport relayer rly query connection](#starport-relayer-rly-query-connection)	 - query the connection state for a given connection id on a network by chain ID
* [starport relayer rly query connection-channels](#starport-relayer-rly-query-connection-channels)	 - query all channels associated with a given connection on a network by chain ID
* [starport relayer rly query connections](#starport-relayer-rly-query-connections)	 - query for all connections on a network by chain ID
* [starport relayer rly query header](#starport-relayer-rly-query-header)	 - query the header of a network by chain ID at a given height or the latest height
* [starport relayer rly query ibc-denoms](#starport-relayer-rly-query-ibc-denoms)	 - query denomination traces for a given network by chain ID
* [starport relayer rly query node-state](#starport-relayer-rly-query-node-state)	 - query the consensus state of a network by chain ID
* [starport relayer rly query packet-commit](#starport-relayer-rly-query-packet-commit)	 - query for the packet commitment given a sequence and channel ID on a network by chain ID
* [starport relayer rly query tx](#starport-relayer-rly-query-tx)	 - query for a transaction on a given network by transaction hash and chain ID
* [starport relayer rly query txs](#starport-relayer-rly-query-txs)	 - query for transactions on a given network by chain ID and a set of transaction events
* [starport relayer rly query unrelayed-acknowledgements](#starport-relayer-rly-query-unrelayed-acknowledgements)	 - query for unrelayed acknowledgement sequence numbers that remain to be relayed on a given path
* [starport relayer rly query unrelayed-packets](#starport-relayer-rly-query-unrelayed-packets)	 - query for the packet sequence numbers that remain to be relayed on a given path
* [starport relayer rly query valset](#starport-relayer-rly-query-valset)	 - query the validator set at particular height for a network by chain ID


## starport relayer rly 



```
starport relayer rly  [flags]
```

**Options**

```
  -h, --help   help for this command
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer


## starport relayer rly 



```
starport relayer rly  [flags]
```

**Options**

```
  -h, --help   help for this command
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer


## starport relayer rly query account

query the relayer's account on a given network by chain ID

```
starport relayer rly query account [chain-id] [flags]
```

**Examples**

```
$ rly query account ibc-0
$ rly q acc ibc-1
```

**Options**

```
  -h, --help   help for account
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query balance

query the relayer's account balance on a given network by chain-ID

```
starport relayer rly query balance [chain-id] [[key-name]] [flags]
```

**Examples**

```
$ rly query balance ibc-0
$ rly query balance ibc-0 testkey
```

**Options**

```
  -h, --help         help for balance
  -i, --ibc-denoms   Display IBC denominations for sending tokens back to other chains
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query channel

query a channel by channel and port ID on a network by chain ID

```
starport relayer rly query channel [chain-id] [channel-id] [port-id] [flags]
```

**Examples**

```
$ rly query channel ibc-0 ibczerochannel transfer
$ rly query channel ibc-2 ibctwochannel transfer --height 1205
```

**Options**

```
      --height int   Height of headers to fetch
  -h, --help         help for channel
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query channels

query for all channels on a network by chain ID

```
starport relayer rly query channels [chain-id] [flags]
```

**Examples**

```
$ rly query channels ibc-0
$ rly query channels ibc-2 --offset 2 --limit 30
```

**Options**

```
  -h, --help          help for channels
  -l, --limit uint    pagination limit for query (default 1000)
  -o, --offset uint   pagination offset for query
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query client

query the state of a light client on a network by chain ID

```
starport relayer rly query client [chain-id] [client-id] [flags]
```

**Examples**

```
$ rly query client ibc-0 ibczeroclient
$ rly query client ibc-0 ibczeroclient --height 1205
```

**Options**

```
      --height int   Height of headers to fetch
  -h, --help         help for client
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query client-connections

query for all connections for a given client on a network by chain ID

```
starport relayer rly query client-connections [chain-id] [client-id] [flags]
```

**Examples**

```
$ rly query client-connections ibc-0 ibczeroclient
$ rly query client-connections ibc-0 ibczeroclient --height 1205
```

**Options**

```
      --height int   Height of headers to fetch
  -h, --help         help for client-connections
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query clients

query for all light client states on a network by chain ID

```
starport relayer rly query clients [chain-id] [flags]
```

**Examples**

```
$ rly query clients ibc-0
$ rly query clients ibc-2 --offset 2 --limit 30
```

**Options**

```
  -h, --help          help for clients
  -l, --limit uint    pagination limit for query (default 1000)
  -o, --offset uint   pagination offset for query
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query connection

query the connection state for a given connection id on a network by chain ID

```
starport relayer rly query connection [chain-id] [connection-id] [flags]
```

**Examples**

```
$ rly query connection ibc-0 ibconnection0
$ rly q conn ibc-1 ibconeconn
```

**Options**

```
  -h, --help   help for connection
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query connection-channels

query all channels associated with a given connection on a network by chain ID

```
starport relayer rly query connection-channels [chain-id] [connection-id] [flags]
```

**Examples**

```
$ rly query connection-channels ibc-0 ibcconnection1
$ rly query connection-channels ibc-2 ibcconnection2 --offset 2 --limit 30
```

**Options**

```
  -h, --help          help for connection-channels
  -l, --limit uint    pagination limit for query (default 1000)
  -o, --offset uint   pagination offset for query
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query connections

query for all connections on a network by chain ID

```
starport relayer rly query connections [chain-id] [flags]
```

**Examples**

```
$ rly query connections ibc-0
$ rly query connections ibc-2 --offset 2 --limit 30
$ rly q conns ibc-1
```

**Options**

```
  -h, --help          help for connections
  -l, --limit uint    pagination limit for query (default 1000)
  -o, --offset uint   pagination offset for query
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query header

query the header of a network by chain ID at a given height or the latest height

```
starport relayer rly query header [chain-id] [[height]] [flags]
```

**Examples**

```
$ rly query header ibc-0
$ rly query header ibc-0 1400
```

**Options**

```
  -h, --help   help for header
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query ibc-denoms

query denomination traces for a given network by chain ID

```
starport relayer rly query ibc-denoms [chain-id] [flags]
```

**Examples**

```
$ rly query ibc-denoms ibc-0
$ rly q ibc-denoms ibc-0
```

**Options**

```
  -h, --help   help for ibc-denoms
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query node-state

query the consensus state of a network by chain ID

```
starport relayer rly query node-state [chain-id] [flags]
```

**Examples**

```
$ rly query node-state ibc-0
$ rly q node-state ibc-1
```

**Options**

```
  -h, --help   help for node-state
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query packet-commit

query for the packet commitment given a sequence and channel ID on a network by chain ID

```
starport relayer rly query packet-commit [chain-id] [channel-id] [port-id] [seq] [flags]
```

**Examples**

```
$ rly query packet-commit ibc-0 ibczerochannel transfer 32
$ rly q packet-commit ibc-1 ibconechannel transfer 31
```

**Options**

```
  -h, --help   help for packet-commit
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query tx

query for a transaction on a given network by transaction hash and chain ID

```
starport relayer rly query tx [chain-id] [tx-hash] [flags]
```

**Examples**

```
$ rly query tx ibc-0 [tx-hash]
$ rly q tx ibc-0 A5DF8D272F1C451CFF92BA6C41942C4D29B5CF180279439ED6AB038282F956BE
```

**Options**

```
  -h, --help   help for tx
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query txs

query for transactions on a given network by chain ID and a set of transaction events

**Synopsis**

Search for a paginated list of transactions that match the given set of
events. Each event takes the form of '{eventType}.{eventAttribute}={value}' with multiple events
separated by '&'.

Please refer to each module's documentation for the full set of events to query for. Each module
documents its respective events under 'cosmos-sdk/x/{module}/spec/xx_events.md'.

```
starport relayer rly query txs [chain-id] [events] [flags]
```

**Examples**

```
$ rly query txs ibc-0 "message.action=transfer" --offset 1 --limit 10
$ rly q txs ibc-0 "message.action=transfer"
```

**Options**

```
  -h, --help          help for txs
  -l, --limit uint    pagination limit for query (default 1000)
  -o, --offset uint   pagination offset for query
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query unrelayed-acknowledgements

query for unrelayed acknowledgement sequence numbers that remain to be relayed on a given path

```
starport relayer rly query unrelayed-acknowledgements [path] [flags]
```

**Examples**

```
$ rly q unrelayed-acknowledgements demo-path
$ rly query unrelayed-acknowledgements demo-path
$ rly query unrelayed-acks demo-path
```

**Options**

```
  -h, --help   help for unrelayed-acknowledgements
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query unrelayed-packets

query for the packet sequence numbers that remain to be relayed on a given path

```
starport relayer rly query unrelayed-packets [path] [flags]
```

**Examples**

```
$ rly q unrelayed-packets demo-path
$ rly query unrelayed-packets demo-path
$ rly query unrelayed-pkts demo-path
```

**Options**

```
  -h, --help   help for unrelayed-packets
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly query valset

query the validator set at particular height for a network by chain ID

```
starport relayer rly query valset [chain-id] [flags]
```

**Examples**

```
$ rly query valset ibc-0
$ rly q valset ibc-1
```

**Options**

```
  -h, --help   help for valset
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly query](#starport-relayer-rly-query)	 - IBC query commands


## starport relayer rly start

Start the listening relayer on a given path

```
starport relayer rly start [path-name] [flags]
```

**Examples**

```
$ rly start demo-path --max-msgs 3
$ rly start demo-path2 --max-tx-size 10
```

**Options**

```
  -h, --help                      help for start
  -l, --max-msgs string           maximum number of messages in a relay transaction (default "5")
  -s, --max-tx-size string        strategy of path to generate of the messages in a relay transaction (default "2")
      --time-threshold duration   time before to expiry time to update client (default 6h0m0s)
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer


## starport relayer rly testnets

commands for joining and running relayer testnets

**Options**

```
  -h, --help   help for testnets
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer
* [starport relayer rly testnets faucet](#starport-relayer-rly-testnets-faucet)	 - listens on a port for requests for tokens
* [starport relayer rly testnets request](#starport-relayer-rly-testnets-request)	 - request tokens from a relayer faucet


## starport relayer rly testnets faucet

listens on a port for requests for tokens

```
starport relayer rly testnets faucet [chain-id] [key-name] [amount] [flags]
```

**Examples**

```
$ rly testnets faucet ibc-0 testkey 100000stake --listen http://0.0.0.0:8081
$ rly tst faucet ibc-0 testkey 100000stake
```

**Options**

```
  -h, --help            help for faucet
  -l, --listen string   sets the faucet listener addresss (default "0.0.0.0:8000")
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly testnets](#starport-relayer-rly-testnets)	 - commands for joining and running relayer testnets


## starport relayer rly testnets request

request tokens from a relayer faucet

```
starport relayer rly testnets request [chain-id] [[key-name]] [flags]
```

**Examples**

```
$ rly testnets request ibc-0 --url http://0.0.0.0:8000
$ rly testnets request ibc-0 testkey --url http://0.0.0.0:8000
$ rly tst req ibc-0
```

**Options**

```
  -h, --help         help for request
  -u, --url string   url to fetch data from
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly testnets](#starport-relayer-rly-testnets)	 - commands for joining and running relayer testnets


## starport relayer rly transact

IBC transaction commands

**Synopsis**

Commands to create IBC transactions on pre-configured chains.
Most of these commands take a [path] argument. Make sure:
  1. Chains are properly configured to relay over by using the 'rly chains list' command
  2. Path is properly configured to relay over by using the 'rly paths list' command

**Options**

```
  -h, --help   help for transact
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer
* [starport relayer rly transact ](#starport-relayer-rly-transact-)	 - 
* [starport relayer rly transact ](#starport-relayer-rly-transact-)	 - 
* [starport relayer rly transact ](#starport-relayer-rly-transact-)	 - 
* [starport relayer rly transact channel-close](#starport-relayer-rly-transact-channel-close)	 - close a channel between two configured chains with a configured path
* [starport relayer rly transact clients](#starport-relayer-rly-transact-clients)	 - create a clients between two configured chains with a configured path
* [starport relayer rly transact connection](#starport-relayer-rly-transact-connection)	 - create a connection between two configured chains with a configured path
* [starport relayer rly transact link](#starport-relayer-rly-transact-link)	 - create clients, connection, and channel between two configured chains with a configured path
* [starport relayer rly transact link-then-start](#starport-relayer-rly-transact-link-then-start)	 - a shorthand command to execute 'link' followed by 'start'
* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands
* [starport relayer rly transact relay-acknowledgements](#starport-relayer-rly-transact-relay-acknowledgements)	 - relay any remaining non-relayed acknowledgements on a given path, in both directions
* [starport relayer rly transact relay-packets](#starport-relayer-rly-transact-relay-packets)	 - relay any remaining non-relayed packets on a given path, in both directions
* [starport relayer rly transact send](#starport-relayer-rly-transact-send)	 - send funds to a different address on the same chain
* [starport relayer rly transact transfer](#starport-relayer-rly-transact-transfer)	 - initiate a transfer from one network to another
* [starport relayer rly transact update-clients](#starport-relayer-rly-transact-update-clients)	 - update IBC clients between two configured chains with a configured path
* [starport relayer rly transact upgrade-chain](#starport-relayer-rly-transact-upgrade-chain)	 - upgrade an IBC-enabled network with a given upgrade plan
* [starport relayer rly transact upgrade-clients](#starport-relayer-rly-transact-upgrade-clients)	 - upgrades IBC clients between two configured chains with a configured path and chain-id


## starport relayer rly 



```
starport relayer rly  [flags]
```

**Options**

```
  -h, --help   help for this command
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer


## starport relayer rly 



```
starport relayer rly  [flags]
```

**Options**

```
  -h, --help   help for this command
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer


## starport relayer rly 



```
starport relayer rly  [flags]
```

**Options**

```
  -h, --help   help for this command
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer


## starport relayer rly transact channel-close

close a channel between two configured chains with a configured path

```
starport relayer rly transact channel-close [path-name] [flags]
```

**Examples**

```
$ rly transact channel-close demo-path
$ rly tx channel-close demo-path --timeout 5s
$ rly tx channel-close demo-path
$ rly tx channel-close demo-path -o 3s
```

**Options**

```
  -h, --help             help for channel-close
  -o, --timeout string   timeout between relayer runs (default "10s")
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact](#starport-relayer-rly-transact)	 - IBC transaction commands


## starport relayer rly transact clients

create a clients between two configured chains with a configured path

**Synopsis**

Creates a working ibc client for chain configured on each end of the path by querying headers from each chain and then sending the corresponding create-client messages

```
starport relayer rly transact clients [path-name] [flags]
```

**Examples**

```
$ rly transact clients demo-path
```

**Options**

```
  -h, --help                        help for clients
      --override                    option to not reuse existing client
  -e, --update-after-expiry         allow governance to update the client if expiry occurs (default true)
  -m, --update-after-misbehaviour   allow governance to update the client if misbehaviour freezing occurs (default true)
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact](#starport-relayer-rly-transact)	 - IBC transaction commands


## starport relayer rly transact connection

create a connection between two configured chains with a configured path

**Synopsis**

Create or repair a connection between two IBC-connected networks
along a specific path.

```
starport relayer rly transact connection [path-name] [flags]
```

**Examples**

```
$ rly transact connection demo-path
$ rly tx conn demo-path --timeout 5s
```

**Options**

```
  -h, --help                        help for connection
  -r, --max-retries uint            maximum retries after failed message send (default 3)
      --override                    option to not reuse existing client
  -o, --timeout string              timeout between relayer runs (default "10s")
  -e, --update-after-expiry         allow governance to update the client if expiry occurs (default true)
  -m, --update-after-misbehaviour   allow governance to update the client if misbehaviour freezing occurs (default true)
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact](#starport-relayer-rly-transact)	 - IBC transaction commands


## starport relayer rly transact link

create clients, connection, and channel between two configured chains with a configured path

**Synopsis**

Create an IBC client between two IBC-enabled networks, in addition
to creating a connection and a channel between the two networks on a configured path.

```
starport relayer rly transact link [path-name] [flags]
```

**Examples**

```
$ rly transact link demo-path
$ rly tx link demo-path
$ rly tx connect demo-path
```

**Options**

```
  -h, --help                        help for link
  -r, --max-retries uint            maximum retries after failed message send (default 3)
      --override                    option to not reuse existing client
  -o, --timeout string              timeout between relayer runs (default "10s")
  -e, --update-after-expiry         allow governance to update the client if expiry occurs (default true)
  -m, --update-after-misbehaviour   allow governance to update the client if misbehaviour freezing occurs (default true)
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact](#starport-relayer-rly-transact)	 - IBC transaction commands


## starport relayer rly transact link-then-start

a shorthand command to execute 'link' followed by 'start'

**Synopsis**

Create IBC clients, connection, and channel between two configured IBC
networks with a configured path and then start the relayer on that path.

```
starport relayer rly transact link-then-start [path-name] [flags]
```

**Examples**

```
$ rly transact link-then-start demo-path
$ rly tx link-then-start demo-path --timeout 5s
```

**Options**

```
  -h, --help                 help for link-then-start
  -l, --max-msgs string      maximum number of messages in a relay transaction (default "5")
  -r, --max-retries uint     maximum retries after failed message send (default 3)
  -s, --max-tx-size string   strategy of path to generate of the messages in a relay transaction (default "2")
      --override             option to not reuse existing client
  -o, --timeout string       timeout between relayer runs (default "10s")
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact](#starport-relayer-rly-transact)	 - IBC transaction commands


## starport relayer rly transact raw

raw IBC transaction commands

**Options**

```
  -h, --help   help for raw
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact](#starport-relayer-rly-transact)	 - IBC transaction commands
* [starport relayer rly transact raw chan-ack](#starport-relayer-rly-transact-raw-chan-ack)	 - chan-ack
* [starport relayer rly transact raw chan-close-confirm](#starport-relayer-rly-transact-raw-chan-close-confirm)	 - chan-close-confirm
* [starport relayer rly transact raw chan-close-init](#starport-relayer-rly-transact-raw-chan-close-init)	 - chan-close-init
* [starport relayer rly transact raw chan-confirm](#starport-relayer-rly-transact-raw-chan-confirm)	 - chan-confirm
* [starport relayer rly transact raw chan-init](#starport-relayer-rly-transact-raw-chan-init)	 - chan-init
* [starport relayer rly transact raw chan-try](#starport-relayer-rly-transact-raw-chan-try)	 - chan-try
* [starport relayer rly transact raw channel-step](#starport-relayer-rly-transact-raw-channel-step)	 - create the next step in creating a channel between chains with the passed identifiers
* [starport relayer rly transact raw client](#starport-relayer-rly-transact-raw-client)	 - create a client for dst-chain on src-chain
* [starport relayer rly transact raw close-channel-step](#starport-relayer-rly-transact-raw-close-channel-step)	 - create the next step in closing a channel between chains with the passed identifiers
* [starport relayer rly transact raw conn-ack](#starport-relayer-rly-transact-raw-conn-ack)	 - conn-ack
* [starport relayer rly transact raw conn-confirm](#starport-relayer-rly-transact-raw-conn-confirm)	 - conn-confirm
* [starport relayer rly transact raw conn-init](#starport-relayer-rly-transact-raw-conn-init)	 - conn-init
* [starport relayer rly transact raw conn-try](#starport-relayer-rly-transact-raw-conn-try)	 - conn-try
* [starport relayer rly transact raw connection-step](#starport-relayer-rly-transact-raw-connection-step)	 - create a connection between chains, passing in identifiers
* [starport relayer rly transact raw transfer](#starport-relayer-rly-transact-raw-transfer)	 - initiate a transfer from one network to another
* [starport relayer rly transact raw update-client](#starport-relayer-rly-transact-raw-update-client)	 - update client for dst-chain on src-chain


## starport relayer rly transact raw chan-ack

chan-ack

```
starport relayer rly transact raw chan-ack [src-chain-id] [dst-chain-id] [src-client-id] 
		[src-chan-id] [dst-chan-id] [src-port-id] [dst-port-id] [flags]
```

**Examples**

```
$ rly transact raw chan-ack ibc-0 ibc-1 ibczeroclient ibcchan1 ibcchan2 transfer transfer
$ rly tx raw chan-ack ibc-0 ibc-1 ibczeroclient ibcchan1 ibcchan2 transfer transfer
```

**Options**

```
  -h, --help   help for chan-ack
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact raw chan-close-confirm

chan-close-confirm

```
starport relayer rly transact raw chan-close-confirm [src-chain-id] [dst-chain-id] 
		[src-client-id] [src-chan-id] [dst-chan-id] [src-port-id] [dst-port-id] [flags]
```

**Examples**

```
$ rly tx raw chan-close-confirm ibc-0 ibc-1 ibczeroclient ibcchan1 ibcchan2 transfer transfer
```

**Options**

```
  -h, --help   help for chan-close-confirm
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact raw chan-close-init

chan-close-init

```
starport relayer rly transact raw chan-close-init [chain-id] [channel-id] [port-id] [flags]
```

**Examples**

```
$ rly transact raw chan-close-init ibc-0 ibcchan1 transfer
$ rly tx raw chan-close-init ibc-0 ibcchan1 transfer
```

**Options**

```
  -h, --help   help for chan-close-init
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact raw chan-confirm

chan-confirm

```
starport relayer rly transact raw chan-confirm [src-chain-id] [dst-chain-id] [src-client-id] 
		[src-chan-id] [dst-chan-id] [src-port-id] [dst-port-id] [flags]
```

**Examples**

```
$ rly transact raw chan-confirm ibc-0 ibc-1 ibczeroclient ibcchan1 ibcchan2 transfer transfer
$ rly tx raw chan-confirm ibc-0 ibc-1 ibczeroclient ibcchan1 ibcchan2 transfer transfer
```

**Options**

```
  -h, --help   help for chan-confirm
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact raw chan-init

chan-init

```
starport relayer rly transact raw chan-init [src-chain-id] [dst-chain-id] [src-client-id] 
		[dst-client-id] [src-conn-id] [dst-conn-id] [src-chan-id] [dst-chan-id] [src-port-id] [dst-port-id] [ordering] [flags]
```

**Examples**

```
$ rly tx raw chan-init ibc-0 ibc-1 ibczeroclient ibconeclient 
ibcconn1 ibcconn2 ibcchan1 ibcchan2 transfer transfer ordered
```

**Options**

```
  -h, --help   help for chan-init
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact raw chan-try

chan-try

```
starport relayer rly transact raw chan-try [src-chain-id] [dst-chain-id] 
		[src-client-id] [src-conn-id] [src-chan-id] [dst-chan-id] [src-port-id] [dst-port-id] [flags]
```

**Examples**

```
$ rly transact raw chan-try ibc-0 ibc-1 ibczeroclient ibcconn0 ibcchan1 ibcchan2 transfer transfer
$ rly tx raw chan-try ibc-0 ibc-1 ibczeroclient ibcconn0 ibcchan1 ibcchan2 transfer transfer
```

**Options**

```
  -h, --help   help for chan-try
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact raw channel-step

create the next step in creating a channel between chains with the passed identifiers

```
starport relayer rly transact raw channel-step [src-chain-id] [dst-chain-id] [src-client-id] [dst-client-id] 
		[src-connection-id] [dst-connection-id] [src-channel-id] [dst-channel-id] [src-port-id] [dst-port-id] [ordering] [flags]
```

**Examples**

```
$ rly transact raw chan-step ibc-0 ibc-1 ibczeroclient ibconeclient ibcconn1 
ibcconn2 ibcchan1 ibcchan2 transfer transfer ordered
$ rly tx raw channel-step ibc-0 ibc-1 ibczeroclient ibconeclient ibcconn1
 ibcconn2 ibcchan1 ibcchan2 transfer transfer ordered
```

**Options**

```
  -h, --help   help for channel-step
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact raw client

create a client for dst-chain on src-chain

```
starport relayer rly transact raw client [src-chain-id] [dst-chain-id] [client-id] [flags]
```

**Examples**

```
$ rly transact raw client ibc-0 ibc-1 ibczeroclient
$ rly tx raw clnt ibc-1 ibc-0 ibconeclient
```

**Options**

```
  -h, --help                        help for client
  -e, --update-after-expiry         allow governance to update the client if expiry occurs (default true)
  -m, --update-after-misbehaviour   allow governance to update the client if misbehaviour freezing occurs (default true)
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact raw close-channel-step

create the next step in closing a channel between chains with the passed identifiers

```
starport relayer rly transact raw close-channel-step [src-chain-id] [dst-chain-id] [src-client-id] [dst-client-id] 
		[src-connection-id] [dst-connection-id] [src-channel-id] [dst-channel-id] [src-port-id] [dst-port-id] [flags]
```

**Examples**

```
$ rly tx raw close-channel-step ibc-0 ibc-1 ibczeroclient ibconeclient ibcconn1 
ibcconn2 ibcchan1 ibcchan2 transfer transfer
```

**Options**

```
  -h, --help   help for close-channel-step
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact raw conn-ack

conn-ack

```
starport relayer rly transact raw conn-ack [src-chain-id] [dst-chain-id] [dst-client-id] [src-client-id] [src-conn-id] [dst-conn-id] [flags]
```

**Examples**

```
$ rly transact raw conn-ack ibc-0 ibc-1 ibconeclient ibczeroclient ibcconn1 ibcconn2
$ rly tx raw conn-ack ibc-0 ibc-1 ibconeclient ibczeroclient ibcconn1 ibcconn2
```

**Options**

```
  -h, --help   help for conn-ack
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact raw conn-confirm

conn-confirm

```
starport relayer rly transact raw conn-confirm [src-chain-id] [dst-chain-id] [src-client-id] [dst-client-id] [src-conn-id] [dst-conn-id] [flags]
```

**Examples**

```
$ rly transact raw conn-confirm ibc-0 ibc-1 ibczeroclient ibconeclient ibcconn1 ibcconn2
$ rly tx raw conn-confirm ibc-0 ibc-1 ibczeroclient ibconeclient ibcconn1 ibcconn2
```

**Options**

```
  -h, --help   help for conn-confirm
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact raw conn-init

conn-init

```
starport relayer rly transact raw conn-init [src-chain-id] [dst-chain-id] [src-client-id] [dst-client-id] [src-conn-id] [dst-conn-id] [flags]
```

**Examples**

```
$ rly transact raw conn-init ibc-0 ibc-1 ibczeroclient ibconeclient ibcconn1 ibcconn2
$ rly tx raw conn-init ibc-0 ibc-1 ibczeroclient ibconeclient ibcconn1 ibcconn2
```

**Options**

```
  -h, --help   help for conn-init
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact raw conn-try

conn-try

```
starport relayer rly transact raw conn-try [src-chain-id] [dst-chain-id] [src-client-id] [dst-client-id] [src-conn-id] [dst-conn-id] [flags]
```

**Examples**

```
$ rly transact raw conn-try ibc-0 ibc-1 ibczeroclient ibconeclient ibcconn1 ibcconn2
$ rly tx raw conn-try ibc-0 ibc-1 ibczeroclient ibconeclient ibcconn1 ibcconn2
```

**Options**

```
  -h, --help   help for conn-try
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact raw connection-step

create a connection between chains, passing in identifiers

**Synopsis**

This command creates the next handshake message given a specifc set of identifiers. 
		If the command fails, you can safely run it again to repair an unfinished connection

```
starport relayer rly transact raw connection-step [src-chain-id] [dst-chain-id] [src-client-id] [dst-client-id] 
		[src-connection-id] [dst-connection-id] [flags]
```

**Examples**

```
$ rly transact raw connection-step ibc-0 ibc-1 ibczeroclient ibconeclient ibcconn1 ibcconn2
$ rly tx raw conn-step ibc-0 ibc-1 ibczeroclient ibconeclient ibcconn1 ibcconn2
```

**Options**

```
  -h, --help   help for connection-step
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact raw transfer

initiate a transfer from one network to another

**Synopsis**

Initiate a token transfer via IBC between two networks. The created packet
must be relayed to the destination chain.

```
starport relayer rly transact raw transfer [src-chain-id] [dst-chain-id] [amount] [dst-addr] [flags]
```

**Examples**

```
$ rly tx transfer ibc-0 ibc-1 100000stake cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk --path demo-path
$ rly tx transfer ibc-0 ibc-1 100000stake cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk --path demo -y 2 -c 10
$ rly tx transfer ibc-0 ibc-1 100000stake raw:non-bech32-address --path demo
$ rly tx raw send ibc-0 ibc-1 100000stake cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk --path demo -c 5
```

**Options**

```
  -h, --help                           help for transfer
  -p, --path string                    specify the path to relay over
  -y, --timeout-height-offset uint     set timeout height offset for 
  -c, --timeout-time-offset duration   specify the path to relay over
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact raw update-client

update client for dst-chain on src-chain

```
starport relayer rly transact raw update-client [src-chain-id] [dst-chain-id] [client-id] [flags]
```

**Examples**

```
$ rly transact raw update-client ibc-0 ibc-1 ibczeroclient
$ rly tx raw uc ibc-0 ibc-1 ibconeclient
```

**Options**

```
  -h, --help   help for update-client
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact raw](#starport-relayer-rly-transact-raw)	 - raw IBC transaction commands


## starport relayer rly transact relay-acknowledgements

relay any remaining non-relayed acknowledgements on a given path, in both directions

```
starport relayer rly transact relay-acknowledgements [path-name] [flags]
```

**Examples**

```
$ rly transact relay-acknowledgements demo-path
$ rly tx relay-acks demo-path -l 3 -s 6
```

**Options**

```
  -h, --help                 help for relay-acknowledgements
  -l, --max-msgs string      maximum number of messages in a relay transaction (default "5")
  -s, --max-tx-size string   strategy of path to generate of the messages in a relay transaction (default "2")
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact](#starport-relayer-rly-transact)	 - IBC transaction commands


## starport relayer rly transact relay-packets

relay any remaining non-relayed packets on a given path, in both directions

```
starport relayer rly transact relay-packets [path-name] [flags]
```

**Examples**

```
$ rly transact relay-packets demo-path
$ rly tx relay-pkts demo-path
```

**Options**

```
  -h, --help                 help for relay-packets
  -l, --max-msgs string      maximum number of messages in a relay transaction (default "5")
  -s, --max-tx-size string   strategy of path to generate of the messages in a relay transaction (default "2")
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact](#starport-relayer-rly-transact)	 - IBC transaction commands


## starport relayer rly transact send

send funds to a different address on the same chain

```
starport relayer rly transact send [chain-id] [from-key] [to-address] [amount] [flags]
```

**Examples**

```
$ rly tx send testkey cosmos10yft4nc8tacpngwlpyq3u4t88y7qzc9xv0q4y8 10000uatom
```

**Options**

```
  -h, --help   help for send
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact](#starport-relayer-rly-transact)	 - IBC transaction commands


## starport relayer rly transact transfer

initiate a transfer from one network to another

**Synopsis**

Initiate a token transfer via IBC between two networks. The created packet
must be relayed to the destination chain.

```
starport relayer rly transact transfer [src-chain-id] [dst-chain-id] [amount] [dst-addr] [flags]
```

**Examples**

```
$ rly tx transfer ibc-0 ibc-1 100000stake cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk --path demo-path
$ rly tx transfer ibc-0 ibc-1 100000stake cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk --path demo -y 2 -c 10
$ rly tx transfer ibc-0 ibc-1 100000stake raw:non-bech32-address --path demo
$ rly tx raw send ibc-0 ibc-1 100000stake cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk --path demo -c 5
```

**Options**

```
  -h, --help                           help for transfer
  -p, --path string                    specify the path to relay over
  -y, --timeout-height-offset uint     set timeout height offset for 
  -c, --timeout-time-offset duration   specify the path to relay over
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact](#starport-relayer-rly-transact)	 - IBC transaction commands


## starport relayer rly transact update-clients

update IBC clients between two configured chains with a configured path

**Synopsis**

Updates IBC client for chain configured on each end of the supplied path.
Clients are updated by querying headers from each chain and then sending the
corresponding update-client messages.

```
starport relayer rly transact update-clients [path-name] [flags]
```

**Examples**

```
$ rly transact update-clients demo-path
```

**Options**

```
  -h, --help   help for update-clients
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact](#starport-relayer-rly-transact)	 - IBC transaction commands


## starport relayer rly transact upgrade-chain

upgrade an IBC-enabled network with a given upgrade plan

**Synopsis**

Upgrade an IBC-enabled network by providing the chain-id of the
network being upgraded, the new unbonding period, the proposal deposit and the JSN file of the
upgrade plan without the upgrade client state.

```
starport relayer rly transact upgrade-chain [path-name] [chain-id] [new-unbonding-period] [deposit] [path/to/upgradePlan.json] [flags]
```

**Options**

```
  -h, --help   help for upgrade-chain
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact](#starport-relayer-rly-transact)	 - IBC transaction commands


## starport relayer rly transact upgrade-clients

upgrades IBC clients between two configured chains with a configured path and chain-id

```
starport relayer rly transact upgrade-clients [path-name] [chain-id] [flags]
```

**Options**

```
      --height int   Height of headers to fetch
  -h, --help         help for upgrade-clients
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly transact](#starport-relayer-rly-transact)	 - IBC transaction commands


## starport relayer rly version

Print the relayer version info

```
starport relayer rly version [flags]
```

**Examples**

```
$ rly version --json
$ rly v
```

**Options**

```
  -h, --help   help for version
  -j, --json   returns the response in json format
```

**Options inherited from parent commands**

```
  -d, --debug         debug output
      --home string   set home directory (default "/home/runner/.relayer")
```

**SEE ALSO**

* [starport relayer rly](#starport-relayer-rly)	 - Low-level commands from github.com/cosmos/relayer


## starport serve

Start a blockchain node in development

**Synopsis**

Start a blockchain node with automatic reloading

```
starport serve [flags]
```

**Options**

```
  -c, --config string        Starport config file (default: ./config.yml)
  -f, --force-reset          Force reset of the app state on start and every source change
  -h, --help                 help for serve
      --home string          Home directory used for blockchains
  -p, --path string          Path of the app
      --rebuild-proto-once   Enables proto code generation for 3rd party modules
  -r, --reset-once           Reset of the app state on first start
  -v, --verbose              Verbose output
```

**SEE ALSO**

* [starport](#starport)	 - A developer tool for building Cosmos SDK blockchains


## starport type

Scaffold a type with CRUD operations

**Synopsis**

Scaffold a type with create, read, update and delete operations

```
starport type [typeName] [field1] [field2] ... [flags]
```

**Options**

```
  -h, --help            help for type
      --indexed         Scaffold an indexed type
      --module string   Module to add the type into. Default: app's main module
      --no-message      Disable CRUD interaction messages scaffolding
  -p, --path string     path of the app
```

**SEE ALSO**

* [starport](#starport)	 - A developer tool for building Cosmos SDK blockchains


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

* [starport](#starport)	 - A developer tool for building Cosmos SDK blockchains

