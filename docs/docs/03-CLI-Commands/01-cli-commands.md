---
description: Ignite CLI docs.
---

# CLI commands

Documentation for Ignite CLI.
## ignite

Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain

**Synopsis**

Ignite CLI is a tool for creating sovereign blockchains built with Cosmos SDK, the world's
most popular modular blockchain framework. Ignite CLI offers everything you need to scaffold,
test, build, and launch your blockchain.

To get started, create a blockchain:

$ ignite scaffold chain example

Announcements:

⋆ Check out how to integrate the EVM or POA in our latest tutorials: https://tutorials.ignite.com 📖
⋆ Satisfied with Ignite? Or totally fed-up with it? Tell us: https://bit.ly/3WZS2uS


**Options**

```
  -h, --help   help for ignite
```

**SEE ALSO**

* [ignite account](#ignite-account)	 - Create, delete, and show Ignite accounts
* [ignite app](#ignite-app)	 - Create and manage Ignite Apps
* [ignite appregistry](#ignite-appregistry)	 - Browse the Ignite App Registry App
* [ignite chain](#ignite-chain)	 - Build, init and start a blockchain node
* [ignite completion](#ignite-completion)	 - Generates shell completion script.
* [ignite docs](#ignite-docs)	 - Show Ignite CLI docs
* [ignite generate](#ignite-generate)	 - Generate clients, API docs from source code
* [ignite relayer](#ignite-relayer)	 - Connect blockchains with an IBC relayer
* [ignite scaffold](#ignite-scaffold)	 - Create a new blockchain, module, message, query, and more
* [ignite testnet](#ignite-testnet)	 - Simulate and manage test networks
* [ignite version](#ignite-version)	 - Print the current build information


## ignite account

Create, delete, and show Ignite accounts

**Synopsis**

Commands for managing Ignite accounts. An Ignite account is a private/public
keypair stored in a keyring. Currently Ignite accounts are used when interacting
with Ignite Apps (namely ignite relayer, ignite network and ignite connect).

Note: Ignite account commands are not for managing your chain's keys and accounts. Use
you chain's binary to manage accounts from "config.yml". For example, if your
blockchain is called "mychain", use "mychaind keys" to manage keys for the
chain.


**Options**

```
  -h, --help                     help for account
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
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
      --coin-type uint32   coin type to use for the account (default 118)
  -h, --help               help for create
```

**Options inherited from parent commands**

```
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**SEE ALSO**

* [ignite account](#ignite-account)	 - Create, delete, and show Ignite accounts


## ignite account delete

Delete an account by name

```
ignite account delete [name] [flags]
```

**Options**

```
  -h, --help   help for delete
```

**Options inherited from parent commands**

```
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**SEE ALSO**

* [ignite account](#ignite-account)	 - Create, delete, and show Ignite accounts


## ignite account export

Export an account as a private key

```
ignite account export [name] [flags]
```

**Options**

```
  -h, --help                help for export
      --non-interactive     do not enter into interactive mode
      --passphrase string   passphrase to encrypt the exported key
      --path string         path to export private key. default: ./key_[name]
```

**Options inherited from parent commands**

```
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**SEE ALSO**

* [ignite account](#ignite-account)	 - Create, delete, and show Ignite accounts


## ignite account import

Import an account by using a mnemonic or a private key

```
ignite account import [name] [flags]
```

**Options**

```
      --coin-type uint32    coin type to use for the account (default 118)
  -h, --help                help for import
      --non-interactive     do not enter into interactive mode
      --passphrase string   passphrase to decrypt the imported key (ignored when secret is a mnemonic)
      --secret string       Your mnemonic or path to your private key (use interactive mode instead to securely pass your mnemonic)
```

**Options inherited from parent commands**

```
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**SEE ALSO**

* [ignite account](#ignite-account)	 - Create, delete, and show Ignite accounts


## ignite account list

Show a list of all accounts

```
ignite account list [flags]
```

**Options**

```
      --address-prefix string   account address prefix (default "cosmos")
  -h, --help                    help for list
```

**Options inherited from parent commands**

```
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**SEE ALSO**

* [ignite account](#ignite-account)	 - Create, delete, and show Ignite accounts


## ignite account show

Show detailed information about a particular account

```
ignite account show [name] [flags]
```

**Options**

```
      --address-prefix string   account address prefix (default "cosmos")
  -h, --help                    help for show
```

**Options inherited from parent commands**

```
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**SEE ALSO**

* [ignite account](#ignite-account)	 - Create, delete, and show Ignite accounts


## ignite app

Create and manage Ignite Apps

**Options**

```
  -h, --help   help for app
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain
* [ignite app describe](#ignite-app-describe)	 - Print information about installed apps
* [ignite app install](#ignite-app-install)	 - Install app
* [ignite app list](#ignite-app-list)	 - List installed apps
* [ignite app scaffold](#ignite-app-scaffold)	 - Scaffold a new Ignite App
* [ignite app uninstall](#ignite-app-uninstall)	 - Uninstall app
* [ignite app update](#ignite-app-update)	 - Update app


## ignite app describe

Print information about installed apps

**Synopsis**

Print information about an installed Ignite App commands and hooks.

```
ignite app describe [path] [flags]
```

**Examples**

```
ignite app describe github.com/org/my-app/
```

**Options**

```
  -h, --help   help for describe
```

**SEE ALSO**

* [ignite app](#ignite-app)	 - Create and manage Ignite Apps


## ignite app install

Install app

**Synopsis**

Installs an Ignite App.

Respects key value pairs declared after the app path to be added to the generated configuration definition.

```
ignite app install [path] [key=value]... [flags]
```

**Examples**

```
ignite app install github.com/org/my-app/ foo=bar baz=qux
```

**Options**

```
  -g, --global   use global plugins configuration ($HOME/.ignite/apps/igniteapps.yml)
  -h, --help     help for install
```

**SEE ALSO**

* [ignite app](#ignite-app)	 - Create and manage Ignite Apps


## ignite app list

List installed apps

**Synopsis**

Prints status and information of all installed Ignite Apps.

```
ignite app list [flags]
```

**Options**

```
  -h, --help   help for list
```

**SEE ALSO**

* [ignite app](#ignite-app)	 - Create and manage Ignite Apps


## ignite app scaffold

Scaffold a new Ignite App

**Synopsis**

Scaffolds a new Ignite App in the current directory.

A git repository will be created with the given module name, unless the current directory is already a git repository.

```
ignite app scaffold [name] [flags]
```

**Examples**

```
ignite app scaffold github.com/org/my-app/
```

**Options**

```
  -h, --help   help for scaffold
```

**SEE ALSO**

* [ignite app](#ignite-app)	 - Create and manage Ignite Apps


## ignite app uninstall

Uninstall app

**Synopsis**

Uninstalls an Ignite App specified by path.

```
ignite app uninstall [path] [flags]
```

**Examples**

```
ignite app uninstall github.com/org/my-app/
```

**Options**

```
  -g, --global   use global plugins configuration ($HOME/.ignite/apps/igniteapps.yml)
  -h, --help     help for uninstall
```

**SEE ALSO**

* [ignite app](#ignite-app)	 - Create and manage Ignite Apps


## ignite app update

Update app

**Synopsis**

Updates an Ignite App specified by path.

If no path is specified all declared apps are updated.

```
ignite app update [path] [flags]
```

**Examples**

```
ignite app update github.com/org/my-app/
```

**Options**

```
  -h, --help   help for update
```

**SEE ALSO**

* [ignite app](#ignite-app)	 - Create and manage Ignite Apps


## ignite appregistry

Browse the Ignite App Registry App

```
ignite appregistry [flags]
```

**Options**

```
  -h, --help   help for appregistry
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain


## ignite chain

Build, init and start a blockchain node

**Synopsis**

Commands in this namespace let you to build, initialize, and start your
blockchain node locally for development purposes.

To run these commands you should be inside the project's directory so that
Ignite can find the source code. To ensure that you are, run "ls", you should
see the following files in the output: "go.mod", "x", "proto", "app", etc.

By default the "build" command will identify the "main" package of the project,
install dependencies if necessary, set build flags, compile the project into a
binary and install the binary. The "build" command is useful if you just want
the compiled binary, for example, to initialize and start the chain manually. It
can also be used to release your chain's binaries automatically as part of
continuous integration workflow.

The "init" command will build the chain's binary and use it to initialize a
local validator node. By default the validator node will be initialized in your
$HOME directory in a hidden directory that matches the name of your project.
This directory is called a data directory and contains a chain's genesis file
and a validator key. This command is useful if you want to quickly build and
initialize the data directory and use the chain's binary to manually start the
blockchain. The "init" command is meant only for development purposes, not
production.

The "serve" command builds, initializes, and starts your blockchain locally with
a single validator node for development purposes. "serve" also watches the
source code directory for file changes and intelligently
re-builds/initializes/starts the chain, essentially providing "code-reloading".
The "serve" command is meant only for development purposes, not production.

To distinguish between production and development consider the following.

In production, blockchains often run the same software on many validator nodes
that are run by different people and entities. To launch a blockchain in
production, the validator entities coordinate the launch process to start their
nodes simultaneously.

During development, a blockchain can be started locally on a single validator
node. This convenient process lets you restart a chain quickly and iterate
faster. Starting a chain on a single node in development is similar to starting
a traditional web application on a local server.

The "faucet" command lets you send tokens to an address from the "faucet"
account defined in "config.yml". Alternatively, you can use the chain's binary
to send token from any other account that exists on chain.

The "simulate" command helps you start a simulation testing process for your
chain.


**Options**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -h, --help            help for chain
  -y, --yes             answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain
* [ignite chain build](#ignite-chain-build)	 - Build a node binary
* [ignite chain debug](#ignite-chain-debug)	 - Launch a debugger for a blockchain app
* [ignite chain faucet](#ignite-chain-faucet)	 - Send coins to an account
* [ignite chain init](#ignite-chain-init)	 - Initialize your chain
* [ignite chain lint](#ignite-chain-lint)	 - Lint codebase using golangci-lint
* [ignite chain modules](#ignite-chain-modules)	 - Manage modules
* [ignite chain serve](#ignite-chain-serve)	 - Start a blockchain node in development
* [ignite chain simulate](#ignite-chain-simulate)	 - Run simulation testing for the blockchain


## ignite chain build

Build a node binary

**Synopsis**


The build command compiles the source code of the project into a binary and
installs the binary in the $(go env GOPATH)/bin directory.

You can customize the output directory for the binary using a flag:

	ignite chain build --output dist

To compile the binary Ignite first compiles protocol buffer (proto) files into
Go source code. Proto files contain required type and services definitions. If
you're using another program to compile proto files, you can use a flag to tell
Ignite to skip the proto compilation step:

	ignite chain build --skip-proto

Afterwards, Ignite install dependencies specified in the go.mod file. By default
Ignite doesn't check that dependencies of the main module stored in the module
cache have not been modified since they were downloaded. To enforce dependency
checking (essentially, running "go mod verify") use a flag:

	ignite chain build --check-dependencies

Next, Ignite identifies the "main" package of the project. By default the "main"
package is located in "cmd/{app}d" directory, where "{app}" is the name of the
scaffolded project and "d" stands for daemon. If your project contains more
than one "main" package, specify the path to the one that Ignite should compile
in config.yml:

	build:
	  main: custom/path/to/main

By default the binary name will match the top-level module name (specified in
go.mod) with a suffix "d". This can be customized in config.yml:

	build:
	  binary: mychaind

You can also specify custom linker flags:

	build:
	  ldflags:
	    - "-X main.Version=development"
	    - "-X main.Date=01/05/2022T19:54"

To build binaries for a release, use the --release flag. The binaries for one or
more specified release targets are built in a "release/" directory in the
project's source directory. Specify the release targets with GOOS:GOARCH build
tags. If the optional --release.targets is not specified, a binary is created
for your current environment.

	ignite chain build --release -t linux:amd64 -t darwin:amd64 -t darwin:arm64


```
ignite chain build [flags]
```

**Options**

```
      --build.tags strings        parameters to build the chain binary
      --check-dependencies        verify that cached dependencies have not been modified since they were downloaded
      --clear-cache               clear the build cache (advanced)
      --debug                     build a debug binary
  -h, --help                      help for build
  -o, --output string             binary output path
  -p, --path string               path of the app (default ".")
      --release                   build for a release
      --release.prefix string     tarball prefix for each release target. Available only with --release flag
  -t, --release.targets strings   release targets. Available only with --release flag
      --skip-proto                skip file generation from proto
  -v, --verbose                   verbose output
```

**Options inherited from parent commands**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -y, --yes             answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite chain](#ignite-chain)	 - Build, init and start a blockchain node


## ignite chain debug

Launch a debugger for a blockchain app

**Synopsis**

The debug command starts a debug server and launches a debugger.

Ignite uses the Delve debugger by default. Delve enables you to interact with
your program by controlling the execution of the process, evaluating variables,
and providing information of thread / goroutine state, CPU register state and
more.

A debug server can optionally be started in cases where default terminal client
is not desirable. When the server starts it first runs the blockchain app,
attaches to it and finally waits for a client connection. It accepts both
JSON-RPC or DAP client connections.

To start a debug server use the following flag:

	ignite chain debug --server

To start a debug server with a custom address use the following flags:

	ignite chain debug --server --server-address 127.0.0.1:30500

The debug server stops automatically when the client connection is closed.


```
ignite chain debug [flags]
```

**Options**

```
  -h, --help                    help for debug
  -p, --path string             path of the app (default ".")
      --server                  start a debug server
      --server-address string   debug server address (default "127.0.0.1:30500")
```

**Options inherited from parent commands**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -y, --yes             answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite chain](#ignite-chain)	 - Build, init and start a blockchain node


## ignite chain faucet

Send coins to an account

```
ignite chain faucet [address] [coin<,...>] [flags]
```

**Options**

```
  -h, --help          help for faucet
      --home string   directory where the blockchain node is initialized
  -p, --path string   path of the app (default ".")
  -v, --verbose       verbose output
```

**Options inherited from parent commands**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -y, --yes             answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite chain](#ignite-chain)	 - Build, init and start a blockchain node


## ignite chain init

Initialize your chain

**Synopsis**

The init command compiles and installs the binary (like "ignite chain build")
and uses that binary to initialize the blockchain's data directory for one
validator. To learn how the build process works, refer to "ignite chain build
--help".

By default, the data directory will be initialized in $HOME/.mychain, where
"mychain" is the name of the project. To set a custom data directory use the
--home flag or set the value in config.yml:

	validators:
	  - name: alice
	    bonded: '100000000stake'
	    home: "~/.customdir"

The data directory contains three files in the "config" directory: app.toml,
config.toml, client.toml. These files let you customize the behavior of your
blockchain node and the client executable. When a chain is re-initialized the
data directory can be reset. To make some values in these files persistent, set
them in config.yml:

	validators:
	  - name: alice
	    bonded: '100000000stake'
	    app:
	      minimum-gas-prices: "0.025stake"
	    config:
	      consensus:
	        timeout_commit: "5s"
	        timeout_propose: "5s"
	    client:
	      output: "json"

The configuration above changes the minimum gas price of the validator (by
default the gas price is set to 0 to allow "free" transactions), sets the block
time to 5s, and changes the output format to JSON. To see what kind of values
this configuration accepts see the generated TOML files in the data directory.

As part of the initialization process Ignite creates on-chain accounts with
token balances. By default, config.yml has two accounts in the top-level
"accounts" property. You can add more accounts and change their token balances.
Refer to config.yml guide to see which values you can set.

One of these accounts is a validator account and the amount of self-delegated
tokens can be set in the top-level "validator" property.

One of the most important components of an initialized chain is the genesis
file, the 0th block of the chain. The genesis file is stored in the data
directory "config" subdirectory and contains the initial state of the chain,
including consensus and module parameters. You can customize the values of the
genesis in config.yml:

	genesis:
	  app_state:
	    staking:
	      params:
	        bond_denom: "foo"

The example above changes the staking token to "foo". If you change the staking
denom, make sure the validator account has the right tokens.

The init command is meant to be used ONLY FOR DEVELOPMENT PURPOSES. Under the
hood it runs commands like "appd init", "appd add-genesis-account", "appd
gentx", and "appd collect-gentx". For production, you may want to run these
commands manually to ensure a production-level node initialization.


```
ignite chain init [flags]
```

**Options**

```
      --build.tags strings   parameters to build the chain binary
      --check-dependencies   verify that cached dependencies have not been modified since they were downloaded
      --clear-cache          clear the build cache (advanced)
      --debug                build a debug binary
  -h, --help                 help for init
      --home string          directory where the blockchain node is initialized
  -p, --path string          path of the app (default ".")
      --skip-proto           skip file generation from proto
  -v, --verbose              verbose output
```

**Options inherited from parent commands**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -y, --yes             answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite chain](#ignite-chain)	 - Build, init and start a blockchain node


## ignite chain lint

Lint codebase using golangci-lint

**Synopsis**

The lint command runs the golangci-lint tool to lint the codebase.

```
ignite chain lint [flags]
```

**Options**

```
  -h, --help   help for lint
```

**Options inherited from parent commands**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -y, --yes             answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite chain](#ignite-chain)	 - Build, init and start a blockchain node


## ignite chain modules

Manage modules

**Synopsis**

The modules command allows you to manage modules in the codebase.

**Options**

```
  -h, --help   help for modules
```

**Options inherited from parent commands**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -y, --yes             answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite chain](#ignite-chain)	 - Build, init and start a blockchain node
* [ignite chain modules list](#ignite-chain-modules-list)	 - List all Cosmos SDK modules in the app


## ignite chain modules list

List all Cosmos SDK modules in the app

**Synopsis**

The list command lists all modules in the app.

```
ignite chain modules list [flags]
```

**Options**

```
  -h, --help   help for list
```

**Options inherited from parent commands**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -y, --yes             answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite chain modules](#ignite-chain-modules)	 - Manage modules


## ignite chain serve

Start a blockchain node in development

**Synopsis**

The serve command compiles and installs the binary (like "ignite chain build"),
uses that binary to initialize the blockchain's data directory for one validator
(like "ignite chain init"), and starts the node locally for development purposes
with automatic code reloading.

Automatic code reloading means Ignite starts watching the project directory.
Whenever a file change is detected, Ignite automatically rebuilds, reinitializes
and restarts the node.

Whenever possible Ignite will try to keep the current state of the chain by
exporting and importing the genesis file.

To force Ignite to start from a clean slate even if a genesis file exists, use
the following flag:

	ignite chain serve --reset-once

To force Ignite to reset the state every time the source code is modified, use
the following flag:

	ignite chain serve --force-reset

With Ignite it's possible to start more than one blockchain from the same source
code using different config files. This is handy if you're building
inter-blockchain functionality and, for example, want to try sending packets
from one blockchain to another. To start a node using a specific config file:

	ignite chain serve --config mars.yml

The serve command is meant to be used ONLY FOR DEVELOPMENT PURPOSES. Under the
hood, it runs "appd start", where "appd" is the name of your chain's binary. For
production, you may want to run "appd start" manually.


```
ignite chain serve [flags]
```

**Options**

```
      --build.tags strings   parameters to build the chain binary
      --check-dependencies   verify that cached dependencies have not been modified since they were downloaded
      --clear-cache          clear the build cache (advanced)
  -f, --force-reset          force reset of the app state on start and every source change
      --generate-clients     generate code for the configured clients on reset or source code change
  -h, --help                 help for serve
      --home string          directory where the blockchain node is initialized
  -o, --output-file string   output file logging the chain output (no UI, no stdin, listens for SIGTERM, implies --yes) (default: stdout)
  -p, --path string          path of the app (default ".")
      --quit-on-fail         quit program if the app fails to start
  -r, --reset-once           reset the app state once on init
      --skip-build           skip initial build of the app (uses local binary)
      --skip-proto           skip file generation from proto
  -v, --verbose              verbose output
```

**Options inherited from parent commands**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -y, --yes             answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite chain](#ignite-chain)	 - Build, init and start a blockchain node


## ignite chain simulate

Run simulation testing for the blockchain

**Synopsis**

Run simulation testing for the blockchain. It sends many randomized-input messages of each module to a simulated node.

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
      --seed int                  simulation random seed (default 42)
      --simName string            name of the simulation to run (default "TestFullAppSimulation")
```

**Options inherited from parent commands**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -y, --yes             answers interactive yes/no questions with yes
```

**SEE ALSO**

* [ignite chain](#ignite-chain)	 - Build, init and start a blockchain node


## ignite completion

Generates shell completion script.

```
ignite completion [bash|zsh|fish|powershell] [flags]
```

**Options**

```
  -h, --help   help for completion
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain


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

Such as compiling protocol buffer files into Go or implement particular
functionality, for example, generating an OpenAPI spec.

Produced source code can be regenerated by running a command again and is not
meant to be edited by hand.


**Options**

```
      --clear-cache           clear the build cache (advanced)
      --enable-proto-vendor   enable proto package vendor for missing Buf dependencies
  -h, --help                  help for generate
  -p, --path string           path of the app (default ".")
  -v, --verbose               verbose output
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain
* [ignite generate composables](#ignite-generate-composables)	 - TypeScript frontend client and Vue 3 composables
* [ignite generate openapi](#ignite-generate-openapi)	 - OpenAPI spec for your chain
* [ignite generate proto-go](#ignite-generate-proto-go)	 - Compile protocol buffer files to Go source code required by Cosmos SDK
* [ignite generate ts-client](#ignite-generate-ts-client)	 - TypeScript frontend client


## ignite generate composables

TypeScript frontend client and Vue 3 composables

```
ignite generate composables [flags]
```

**Options**

```
  -h, --help            help for composables
  -o, --output string   Vue 3 composables output path
  -y, --yes             answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
      --clear-cache           clear the build cache (advanced)
      --enable-proto-vendor   enable proto package vendor for missing Buf dependencies
  -p, --path string           path of the app (default ".")
  -v, --verbose               verbose output
```

**SEE ALSO**

* [ignite generate](#ignite-generate)	 - Generate clients, API docs from source code


## ignite generate openapi

OpenAPI spec for your chain

```
ignite generate openapi [flags]
```

**Options**

```
  -h, --help   help for openapi
  -y, --yes    answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
      --clear-cache           clear the build cache (advanced)
      --enable-proto-vendor   enable proto package vendor for missing Buf dependencies
  -p, --path string           path of the app (default ".")
  -v, --verbose               verbose output
```

**SEE ALSO**

* [ignite generate](#ignite-generate)	 - Generate clients, API docs from source code


## ignite generate proto-go

Compile protocol buffer files to Go source code required by Cosmos SDK

```
ignite generate proto-go [flags]
```

**Options**

```
  -h, --help   help for proto-go
  -y, --yes    answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
      --clear-cache           clear the build cache (advanced)
      --enable-proto-vendor   enable proto package vendor for missing Buf dependencies
  -p, --path string           path of the app (default ".")
  -v, --verbose               verbose output
```

**SEE ALSO**

* [ignite generate](#ignite-generate)	 - Generate clients, API docs from source code


## ignite generate ts-client

TypeScript frontend client

**Synopsis**

Generate a framework agnostic TypeScript client for your blockchain project.

By default the TypeScript client is generated in the "ts-client/" directory. You
can customize the output directory in config.yml:

	client:
	  typescript:
	    path: new-path

Output can also be customized by using a flag:

	ignite generate ts-client --output new-path

TypeScript client code can be automatically regenerated on reset or source code
changes when the blockchain is started with a flag:

	ignite chain serve --generate-clients


```
ignite generate ts-client [flags]
```

**Options**

```
      --disable-cache   disable build cache
  -h, --help            help for ts-client
  -o, --output string   TypeScript client output path
  -y, --yes             answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
      --clear-cache           clear the build cache (advanced)
      --enable-proto-vendor   enable proto package vendor for missing Buf dependencies
  -p, --path string           path of the app (default ".")
  -v, --verbose               verbose output
```

**SEE ALSO**

* [ignite generate](#ignite-generate)	 - Generate clients, API docs from source code


## ignite relayer

Connect blockchains with an IBC relayer

```
ignite relayer [flags]
```

**Options**

```
  -h, --help   help for relayer
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain


## ignite scaffold

Create a new blockchain, module, message, query, and more

**Synopsis**

Scaffolding is a quick way to generate code for major pieces of your
application.

For details on each scaffolding target (chain, module, message, etc.) run the
corresponding command with a "--help" flag, for example, "ignite scaffold chain
--help".

The Ignite team strongly recommends committing the code to a version control
system before running scaffolding commands. This will make it easier to see the
changes to the source code as well as undo the command if you've decided to roll
back the changes.

This blockchain you create with the chain scaffolding command uses the modular
Cosmos SDK framework and imports many standard modules for functionality like
proof of stake, token transfer, inter-blockchain connectivity, governance, and
more. Custom functionality is implemented in modules located by convention in
the "x/" directory. By default, your blockchain comes with an empty custom
module. Use the module scaffolding command to create an additional module.

An empty custom module doesn't do much, it's basically a container for logic
that is responsible for processing transactions and changing the application
state. Cosmos SDK blockchains work by processing user-submitted signed
transactions, which contain one or more messages. A message contains data that
describes a state transition. A module can be responsible for handling any
number of messages.

A message scaffolding command will generate the code for handling a new type of
Cosmos SDK message. Message fields describe the state transition that the
message is intended to produce if processed without errors.

Scaffolding messages is useful to create individual "actions" that your module
can perform. Sometimes, however, you want your blockchain to have the
functionality to create, read, update and delete (CRUD) instances of a
particular type. Depending on how you want to store the data there are three
commands that scaffold CRUD functionality for a type: list, map, and single.
These commands create four messages (one for each CRUD action), and the logic to
add, delete, and fetch the data from the store. If you want to scaffold only the
logic, for example, you've decided to scaffold messages separately, you can do
that as well with the "--no-message" flag.

Reading data from a blockchain happens with a help of queries. Similar to how
you can scaffold messages to write data, you can scaffold queries to read the
data back from your blockchain application.

You can also scaffold a type, which just produces a new protocol buffer file
with a proto message description. Note that proto messages produce (and
correspond with) Go types whereas Cosmos SDK messages correspond to proto "rpc"
in the "Msg" service.

If you're building an application with custom IBC logic, you might need to
scaffold IBC packets. An IBC packet represents the data sent from one blockchain
to another. You can only scaffold IBC packets in IBC-enabled modules scaffolded
with an "--ibc" flag. Note that the default module is not IBC-enabled.


**Options**

```
  -h, --help      help for scaffold
  -v, --verbose   verbose output
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain
* [ignite scaffold chain](#ignite-scaffold-chain)	 - New Cosmos SDK blockchain
* [ignite scaffold chain-registry](#ignite-scaffold-chain-registry)	 - Configs for the chain registry
* [ignite scaffold configs](#ignite-scaffold-configs)	 - Configs for a custom Cosmos SDK module
* [ignite scaffold list](#ignite-scaffold-list)	 - CRUD for data stored as an array
* [ignite scaffold map](#ignite-scaffold-map)	 - CRUD for data stored as key-value pairs
* [ignite scaffold message](#ignite-scaffold-message)	 - Message to perform state transition on the blockchain
* [ignite scaffold module](#ignite-scaffold-module)	 - Custom Cosmos SDK module
* [ignite scaffold packet](#ignite-scaffold-packet)	 - Message for sending an IBC packet
* [ignite scaffold params](#ignite-scaffold-params)	 - Parameters for a custom Cosmos SDK module
* [ignite scaffold query](#ignite-scaffold-query)	 - Query for fetching data from a blockchain
* [ignite scaffold single](#ignite-scaffold-single)	 - CRUD for data stored in a single location
* [ignite scaffold type](#ignite-scaffold-type)	 - Type definition
* [ignite scaffold type-list](#ignite-scaffold-type-list)	 - List scaffold types
* [ignite scaffold vue](#ignite-scaffold-vue)	 - Vue 3 web app template


## ignite scaffold chain

New Cosmos SDK blockchain

**Synopsis**

Create a new application-specific Cosmos SDK blockchain.

For example, the following command will create a blockchain called "hello" in
the "hello/" directory:

	ignite scaffold chain hello

A project name can be a simple name or a URL. The name will be used as the Go
module path for the project. Examples of project names:

	ignite scaffold chain foo
	ignite scaffold chain foo/bar
	ignite scaffold chain example.org/foo
	ignite scaffold chain github.com/username/foo
		
A new directory with source code files will be created in the current directory.
To use a different path use the "--path" flag.

Most of the logic of your blockchain is written in custom modules. Each module
effectively encapsulates an independent piece of functionality. Following the
Cosmos SDK convention, custom modules are stored inside the "x/" directory. By
default, Ignite creates a module with a name that matches the name of the
project. To create a blockchain without a default module use the "--no-module"
flag. Additional modules can be added after a project is created with "ignite
scaffold module" command.

Account addresses on Cosmos SDK-based blockchains have string prefixes. For
example, the Cosmos Hub blockchain uses the default "cosmos" prefix, so that
addresses look like this: "cosmos12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf". To
use a custom address prefix use the "--address-prefix" flag. For example:

	ignite scaffold chain foo --address-prefix bar

By default when compiling a blockchain's source code Ignite creates a cache to
speed up the build process. To clear the cache when building a blockchain use
the "--clear-cache" flag. It is very unlikely you will ever need to use this
flag.

The blockchain is using the Cosmos SDK modular blockchain framework. Learn more
about Cosmos SDK on https://docs.cosmos.network


```
ignite scaffold chain [name] [flags]
```

**Options**

```
      --address-prefix string    account address prefix (default "cosmos")
      --clear-cache              clear the build cache (advanced)
      --coin-type uint32         coin type to use for the account (default 118)
      --default-denom string     default staking denom (default "stake")
  -h, --help                     help for chain
      --minimal                  create a minimal blockchain (with the minimum required Cosmos SDK modules)
      --module-configs strings   add module configs
      --no-module                create a project without a default module
      --params strings           add default module parameters
  -p, --path string              create a project in a specific path
      --proto-dir string         chain proto directory (default "proto")
      --skip-git                 skip Git repository initialization
      --skip-proto               skip proto generation
```

**Options inherited from parent commands**

```
  -v, --verbose   verbose output
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Create a new blockchain, module, message, query, and more


## ignite scaffold chain-registry

Configs for the chain registry

**Synopsis**

Scaffold the chain registry chain.json and assets.json files.

The chain registry is a GitHub repo, hosted at https://github.com/cosmos/chain-registry, that
contains the chain.json and assets.json files of most of chains in the Cosmos ecosystem.
It is good practices, when creating a new chain, and about to launch a testnet or mainnet, to
publish the chain's metadata in the chain registry.

Read more about the chain.json at https://github.com/cosmos/chain-registry?tab=readme-ov-file#chainjson
Read more about the assets.json at https://github.com/cosmos/chain-registry?tab=readme-ov-file#assetlists

```
ignite scaffold chain-registry [flags]
```

**Options**

```
      --clear-cache   clear the build cache (advanced)
  -h, --help          help for chain-registry
  -p, --path string   path of the app (default ".")
  -y, --yes           answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
  -v, --verbose   verbose output
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Create a new blockchain, module, message, query, and more


## ignite scaffold configs

Configs for a custom Cosmos SDK module

**Synopsis**

Scaffold a new config for a Cosmos SDK module.

A Cosmos SDK module can have configurations. An example of a config is "address prefix" of the
"auth" module. A config can be scaffolded into a module using the "--module-configs" into
the scaffold module command or using the "scaffold configs" command. By default 
configs are of type "string", but you can specify a type for each config. For example:

	ignite scaffold configs foo baz:uint bar:bool

Refer to Cosmos SDK documentation to learn more about modules, dependencies and
configs.


```
ignite scaffold configs [configs]... [flags]
```

**Options**

```
      --clear-cache     clear the build cache (advanced)
  -h, --help            help for configs
      --module string   module to add the query into (default: app's main module)
  -p, --path string     path of the app (default ".")
  -y, --yes             answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
  -v, --verbose   verbose output
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Create a new blockchain, module, message, query, and more


## ignite scaffold list

CRUD for data stored as an array

**Synopsis**

The "list" scaffolding command is used to generate files that implement the
logic for storing and interacting with data stored as a list in the blockchain
state.

The command accepts a NAME argument that will be used as the name of a new type
of data. It also accepts a list of FIELDs that describe the type.

The interaction with the data follows the create, read, updated, and delete
(CRUD) pattern. For each type three Cosmos SDK messages are defined for writing
data to the blockchain: MsgCreate{Name}, MsgUpdate{Name}, MsgDelete{Name}. For
reading data two queries are defined: {Name} and {Name}All. The type, messages,
and queries are defined in the "proto/" directory as protocol buffer messages.
Messages and queries are mounted in the "Msg" and "Query" services respectively.

When messages are handled, the appropriate keeper methods are called. By
convention, the methods are defined in
"x/{moduleName}/keeper/msg_server_{name}.go". Helpful methods for getting,
setting, removing, and appending are defined in the same "keeper" package in
"{name}.go".

The "list" command essentially allows you to define a new type of data and
provides the logic to create, read, update, and delete instances of the type.
For example, let's review a command that generates the code to handle a list of
posts and each post has "title" and "body" fields:

	ignite scaffold list post title body

This provides you with a "Post" type, MsgCreatePost, MsgUpdatePost,
MsgDeletePost and two queries: Post and PostAll. The compiled CLI, let's say the
binary is "blogd" and the module is "blog", has commands to query the chain (see
"blogd q blog") and broadcast transactions with the messages above (see "blogd
tx blog").

The code generated with the list command is meant to be edited and tailored to
your application needs. Consider the code to be a "skeleton" for the actual
business logic you will implement next.

By default, all fields are assumed to be strings. If you want a field of a
different type, you can specify it after a colon ":". The following types are
supported: string, bool, int, uint, coin, array.string, array.int, array.uint,
array.coin. An example of using field types:

	ignite scaffold list pool amount:coin tags:array.string height:int

For detailed type information use ignite scaffold type --help

"Index" indicates whether the type can be used as an index in
"ignite scaffold map".

Ignite also supports custom types:

	ignite scaffold list product-details name desc
	ignite scaffold list product price:coin details:ProductDetails

In the example above the "ProductDetails" type was defined first, and then used
as a custom type for the "details" field. Ignite doesn't support arrays of
custom types yet.

Your chain will accept custom types in JSON-notation:

	exampled tx example create-product 100coin '{"name": "x", "desc": "y"}' --from alice

By default the code will be scaffolded in the module that matches your project's
name. If you have several modules in your project, you might want to specify a
different module:

	ignite scaffold list post title body --module blog

By default, each message comes with a "creator" field that represents the
address of the transaction signer. You can customize the name of this field with
a flag:

	ignite scaffold list post title body --signer author

It's possible to scaffold just the getter/setter logic without the CRUD
messages. This is useful when you want the methods to handle a type, but would
like to scaffold messages manually. Use a flag to skip message scaffolding:

	ignite scaffold list post title body --no-message

The "creator" field is not generated if a list is scaffolded with the
"--no-message" flag.


```
ignite scaffold list NAME [field]... [flags]
```

**Options**

```
      --clear-cache     clear the build cache (advanced)
  -h, --help            help for list
      --module string   specify which module to generate code in
      --no-message      skip generating message handling logic
      --no-simulation   skip simulation logic
  -p, --path string     path of the app (default ".")
      --signer string   label for the message signer (default: creator)
  -y, --yes             answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
  -v, --verbose   verbose output
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Create a new blockchain, module, message, query, and more


## ignite scaffold map

CRUD for data stored as key-value pairs

**Synopsis**

The "map" scaffolding command is used to generate files that implement the logic
for storing and interacting with data stored as key-value pairs (or a
dictionary) in the blockchain state.

The "map" command is very similar to "ignite scaffold list" with the main
difference in how values are indexed. With "list" values are indexed by an
incrementing integer, whereas "map" values are indexed by a user-provided value
(or multiple values).

Let's use the same blog post example:

	ignite scaffold map post title body:string

This command scaffolds a "Post" type and CRUD functionality to create, read,
updated, and delete posts. However, when creating a new post with your chain's
binary (or by submitting a transaction through the chain's API) you will be
required to provide an "index":

	blogd tx blog create-post [index] [title] [body]
	blogd tx blog create-post hello "My first post" "This is the body"

This command will create a post and store it in the blockchain's state under the
"hello" index. You will be able to fetch back the value of the post by querying
for the "hello" key.

	blogd q blog show-post hello

By default, the index is called "index", to customize the index, use the "--index" flag.

Since the behavior of "list" and "map" scaffolding is very similar, you can use
the "--no-message", "--module", "--signer" flags as well as the colon syntax for
custom types.

For detailed type information use ignite scaffold type --help


```
ignite scaffold map NAME [field]... [flags]
```

**Options**

```
      --clear-cache     clear the build cache (advanced)
  -h, --help            help for map
      --index string    field that index the value (default "index")
      --module string   specify which module to generate code in
      --no-message      skip generating message handling logic
      --no-simulation   skip simulation logic
  -p, --path string     path of the app (default ".")
      --signer string   label for the message signer (default: creator)
  -y, --yes             answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
  -v, --verbose   verbose output
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Create a new blockchain, module, message, query, and more


## ignite scaffold message

Message to perform state transition on the blockchain

**Synopsis**

Message scaffolding is useful for quickly adding functionality to your
blockchain to handle specific Cosmos SDK messages.

Messages are objects whose end goal is to trigger state transitions on the
blockchain. A message is a container for fields of data that affect how the
blockchain's state will change. You can think of messages as "actions" that a
user can perform.

For example, the bank module has a "Send" message for token transfers between
accounts. The send message has three fields: from address (sender), to address
(recipient), and a token amount. When this message is successfully processed,
the token amount will be deducted from the sender's account and added to the
recipient's account.

Ignite's message scaffolding lets you create new types of messages and add them
to your chain. For example:

	ignite scaffold message add-pool amount:coins denom active:bool --module dex

The command above will create a new message MsgAddPool with three fields: amount
(in tokens), denom (a string), and active (a boolean). The message will be added
to the "dex" module.

For detailed type information use ignite scaffold type --help

By default, the message is defined as a proto message in the
"proto/{app}/{module}/tx.proto" and registered in the "Msg" service. A CLI command to
create and broadcast a transaction with MsgAddPool is created in the module's
"cli" package. Additionally, Ignite scaffolds a message constructor and the code
to satisfy the sdk.Msg interface and register the message in the module.

Most importantly in the "keeper" package Ignite scaffolds an "AddPool" function.
Inside this function, you can implement message handling logic.

When successfully processed a message can return data. Use the —response flag to
specify response fields and their types. For example

	ignite scaffold message create-post title body --response id:int,title

The command above will scaffold MsgCreatePost which returns both an ID (an
integer) and a title (a string).

Message scaffolding follows the rules as "ignite scaffold list/map/single" and
supports fields with standard and custom types. See "ignite scaffold list —help"
for details.


```
ignite scaffold message [name] [field1:type1] [field2:type2] ... [flags]
```

**Options**

```
      --clear-cache        clear the build cache (advanced)
  -d, --desc string        description of the command
  -h, --help               help for message
      --module string      module to add the message into. Default: app's main module
      --no-simulation      disable CRUD simulation scaffolding
  -p, --path string        path of the app (default ".")
  -r, --response strings   response fields
      --signer string      label for the message signer (default: creator)
  -y, --yes                answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
  -v, --verbose   verbose output
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Create a new blockchain, module, message, query, and more


## ignite scaffold module

Custom Cosmos SDK module

**Synopsis**

Scaffold a new Cosmos SDK module.

Cosmos SDK is a modular framework and each independent piece of functionality is
implemented in a separate module. By default your blockchain imports a set of
standard Cosmos SDK modules. To implement custom functionality of your
blockchain, scaffold a module and implement the logic of your application.

This command does the following:

* Creates a directory with module's protocol buffer files in "proto/"
* Creates a directory with module's boilerplate Go code in "x/"
* Imports the newly created module by modifying "app/app.go"

This command will proceed with module scaffolding even if "app/app.go" doesn't
have the required default placeholders. If the placeholders are missing, you
will need to modify "app/app.go" manually to import the module. If you want the
command to fail if it can't import the module, use the "--require-registration"
flag.

To scaffold an IBC-enabled module use the "--ibc" flag. An IBC-enabled module is
like a regular module with the addition of IBC-specific logic and placeholders
to scaffold IBC packets with "ignite scaffold packet".

A module can depend on one or more other modules and import their keeper
methods. To scaffold a module with a dependency use the "--dep" flag

For example, your new custom module "foo" might have functionality that requires
sending tokens between accounts. The method for sending tokens is a defined in
the "bank"'s module keeper. You can scaffold a "foo" module with the dependency
on "bank" with the following command:

	ignite scaffold module foo --dep bank

You can then define which methods you want to import from the "bank" keeper in
"expected_keepers.go".

You can also scaffold a module with a list of dependencies that can include both
standard and custom modules (provided they exist):

	ignite scaffold module bar --dep foo,mint,account,FeeGrant

Note: the "--dep" flag doesn't install third-party modules into your
application, it just generates extra code that specifies which existing modules
your new custom module depends on.

A Cosmos SDK module can have parameters (or "params"). Params are values that
can be set at the genesis of the blockchain and can be modified while the
blockchain is running. An example of a param is "Inflation rate change" of the
"mint" module. A module can be scaffolded with params using the "--params" flag
that accepts a list of param names. By default params are of type "string", but
you can specify a type for each param. For example:

	ignite scaffold module foo --params baz:uint,bar:bool

Refer to Cosmos SDK documentation to learn more about modules, dependencies and
params.


```
ignite scaffold module [name] [flags]
```

**Options**

```
      --clear-cache              clear the build cache (advanced)
      --dep strings              add a dependency on another module
  -h, --help                     help for module
      --ibc                      add IBC functionality
      --module-configs strings   add module configs
      --ordering string          channel ordering of the IBC module [none|ordered|unordered] (default "none")
      --params strings           add module parameters
  -p, --path string              path of the app (default ".")
      --require-registration     fail if module can't be registered
  -y, --yes                      answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
  -v, --verbose   verbose output
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Create a new blockchain, module, message, query, and more


## ignite scaffold packet

Message for sending an IBC packet

**Synopsis**

Scaffold an IBC packet in a specific IBC-enabled Cosmos SDK module

```
ignite scaffold packet [packetName] [field1] [field2] ... --module [moduleName] [flags]
```

**Options**

```
      --ack strings     custom acknowledgment type (field1,field2,...)
      --clear-cache     clear the build cache (advanced)
  -h, --help            help for packet
      --module string   IBC Module to add the packet into
      --no-message      disable send message scaffolding
  -p, --path string     path of the app (default ".")
      --signer string   label for the message signer (default: creator)
  -y, --yes             answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
  -v, --verbose   verbose output
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Create a new blockchain, module, message, query, and more


## ignite scaffold params

Parameters for a custom Cosmos SDK module

**Synopsis**

Scaffold a new parameter for a Cosmos SDK module.

A Cosmos SDK module can have parameters (or "params"). Params are values that
can be set at the genesis of the blockchain and can be modified while the
blockchain is running. An example of a param is "Inflation rate change" of the
"mint" module. A params can be scaffolded into a module using the "--params" into
the scaffold module command or using the "scaffold params" command. By default 
params are of type "string", but you can specify a type for each param. For example:

	ignite scaffold params foo baz:uint bar:bool

Refer to Cosmos SDK documentation to learn more about modules, dependencies and
params.


```
ignite scaffold params [param]... [flags]
```

**Options**

```
      --clear-cache     clear the build cache (advanced)
  -h, --help            help for params
      --module string   module to add the query into. Default: app's main module
  -p, --path string     path of the app (default ".")
  -y, --yes             answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
  -v, --verbose   verbose output
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Create a new blockchain, module, message, query, and more


## ignite scaffold query

Query for fetching data from a blockchain

**Synopsis**

Query for fetching data from a blockchain.
		
For detailed type information use ignite scaffold type --help.

```
ignite scaffold query [name] [field1:type1] [field2:type2] ... [flags]
```

**Options**

```
      --clear-cache        clear the build cache (advanced)
  -d, --desc string        description of the CLI to broadcast a tx with the message
  -h, --help               help for query
      --module string      module to add the query into. Default: app's main module
      --paginated          define if the request can be paginated
  -p, --path string        path of the app (default ".")
  -r, --response strings   response fields
  -y, --yes                answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
  -v, --verbose   verbose output
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Create a new blockchain, module, message, query, and more


## ignite scaffold single

CRUD for data stored in a single location

**Synopsis**

CRUD for data stored in a single location.
		
For detailed type information use ignite scaffold type --help.

```
ignite scaffold single NAME [field:type]... [flags]
```

**Examples**

```
  ignite scaffold single todo-single title:string done:bool
```

**Options**

```
      --clear-cache     clear the build cache (advanced)
  -h, --help            help for single
      --module string   specify which module to generate code in
      --no-message      skip generating message handling logic
      --no-simulation   skip simulation logic
  -p, --path string     path of the app (default ".")
      --signer string   label for the message signer (default: creator)
  -y, --yes             answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
  -v, --verbose   verbose output
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Create a new blockchain, module, message, query, and more


## ignite scaffold type

Type definition

**Synopsis**

Type information

Types 		Usage 																										
address 	use '<FIELD_NAME>:address' to scaffold string types (eg: cosmos1abcdefghijklmnopqrstuvwxyz0123456). 														
array.coin 	use '<FIELD_NAME>:array.coin' to scaffold sdk.Coins types (eg: 20stake). Disclaimer: Only one `coins` or `dec.coins` field can accept multiple CLI values per command due to AutoCLI limitations. 		
array.dec.coin 	use '<FIELD_NAME>:array.dec.coin' to scaffold sdk.DecCoins types (eg: 20000002stake). Disclaimer: Only one `coins` or `dec.coins` field can accept multiple CLI values per command due to AutoCLI limitations. 	
array.int 	use '<FIELD_NAME>:array.int' to scaffold []int64 types (eg: 5,4,3,2,1). 																	
array.string 	use '<FIELD_NAME>:array.string' to scaffold []string types (eg: abc,xyz). 																	
array.uint 	use '<FIELD_NAME>:array.uint' to scaffold []uint64 types (eg: 13,26,31,40). 																	
bool 		use '<FIELD_NAME>:bool' to scaffold bool types (eg: true). 																			
bytes 		use '<FIELD_NAME>:bytes' to scaffold []byte types (eg: 3,2,3,5). 																		
coin 		use '<FIELD_NAME>:coin' to scaffold sdk.Coin types (eg: 10token). 																		
coins 		use '<FIELD_NAME>:array.coin' to scaffold sdk.Coins types (eg: 20stake). Disclaimer: Only one `coins` or `dec.coins` field can accept multiple CLI values per command due to AutoCLI limitations. 		
custom 		use the custom type to scaffold already created chain types. 																			
dec.coin 	use '<FIELD_NAME>:dec.coin' to scaffold sdk.DecCoin types (eg: 100001token). 																	
dec.coins 	use '<FIELD_NAME>:array.dec.coin' to scaffold sdk.DecCoins types (eg: 20000002stake). Disclaimer: Only one `coins` or `dec.coins` field can accept multiple CLI values per command due to AutoCLI limitations. 	
int 		use '<FIELD_NAME>:int' to scaffold int64 types (eg: 111). 																			
int64 		use '<FIELD_NAME>:int' to scaffold int64 types (eg: 111). 																			
ints 		use '<FIELD_NAME>:array.int' to scaffold []int64 types (eg: 5,4,3,2,1). 																	
string 		use '<FIELD_NAME>:string' to scaffold string types (eg: xyz). 																			
strings 	use '<FIELD_NAME>:array.string' to scaffold []string types (eg: abc,xyz). 																	
uint 		use '<FIELD_NAME>:uint' to scaffold uint64 types (eg: 111). 																			
uint64 		use '<FIELD_NAME>:uint' to scaffold uint64 types (eg: 111). 																			
uints 		use '<FIELD_NAME>:array.uint' to scaffold []uint64 types (eg: 13,26,31,40). 																	

Field Usage:
    - fieldName
    - fieldName:fieldType

If no :fieldType, default (string) is used



```
ignite scaffold type NAME [field:type] ... [flags]
```

**Examples**

```
  ignite scaffold type todo-item priority:int desc:string tags:array.string done:bool
```

**Options**

```
      --clear-cache     clear the build cache (advanced)
  -h, --help            help for type
      --module string   specify which module to generate code in
      --no-message      skip generating message handling logic
      --no-simulation   skip simulation logic
  -p, --path string     path of the app (default ".")
      --signer string   label for the message signer (default: creator)
  -y, --yes             answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
  -v, --verbose   verbose output
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Create a new blockchain, module, message, query, and more


## ignite scaffold type-list

List scaffold types

**Synopsis**

List all available scaffold types

```
ignite scaffold type-list [flags]
```

**Options**

```
  -h, --help   help for type-list
```

**Options inherited from parent commands**

```
  -v, --verbose   verbose output
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Create a new blockchain, module, message, query, and more


## ignite scaffold vue

Vue 3 web app template

```
ignite scaffold vue [flags]
```

**Options**

```
  -h, --help   help for vue
  -y, --yes    answers interactive yes/no questions with yes
```

**Options inherited from parent commands**

```
  -v, --verbose   verbose output
```

**SEE ALSO**

* [ignite scaffold](#ignite-scaffold)	 - Create a new blockchain, module, message, query, and more


## ignite testnet

Simulate and manage test networks

**Synopsis**

Comprehensive toolset for managing and simulating blockchain test networks. It allows users to either run a test network in place using mainnet data or set up a multi-node environment for more complex testing scenarios. Additionally, it includes a subcommand for simulating the chain, which is useful for fuzz testing and other testing-related tasks.

**Options**

```
  -h, --help   help for testnet
```

**SEE ALSO**

* [ignite](#ignite)	 - Ignite CLI offers everything you need to scaffold, test, build, and launch your blockchain
* [ignite testnet in-place](#ignite-testnet-in-place)	 - Create and start a testnet from current local net state
* [ignite testnet multi-node](#ignite-testnet-multi-node)	 - Initialize and provide multi-node on/off functionality
* [ignite testnet simulate](#ignite-testnet-simulate)	 - Run simulation testing for the blockchain


## ignite testnet in-place

Create and start a testnet from current local net state

**Synopsis**

Testnet in-place command is used to create and start a testnet from current local net state(including mainnet).
After using this command in the repo containing the config.yml file, the network will start.
We can create a testnet from the local network state and mint additional coins for the desired accounts from the config.yml file.

```
ignite testnet in-place [flags]
```

**Options**

```
      --address-prefix string   account address prefix (default "cosmos")
      --check-dependencies      verify that cached dependencies have not been modified since they were downloaded
      --clear-cache             clear the build cache (advanced)
      --coin-type uint32        coin type to use for the account (default 118)
  -h, --help                    help for in-place
      --home string             directory where the blockchain node is initialized
  -p, --path string             path of the app (default ".")
      --skip-proto              skip file generation from proto
  -v, --verbose                 verbose output
```

**SEE ALSO**

* [ignite testnet](#ignite-testnet)	 - Simulate and manage test networks


## ignite testnet multi-node

Initialize and provide multi-node on/off functionality

**Synopsis**

Initialize the test network with the number of nodes and bonded from the config.yml file::
			...
                  validators:
                        - name: alice
                        bonded: 100000000stake
                        - name: validator1
                        bonded: 100000000stake
                        - name: validator2
                        bonded: 200000000stake
                        - name: validator3
                        bonded: 300000000stake


			The "multi-node" command allows developers to easily set up, initialize, and manage multiple nodes for a 
			testnet environment. This command provides full flexibility in enabling or disabling each node as desired, 
			making it a powerful tool for simulating a multi-node blockchain network during development.

			Usage:
					ignite testnet multi-node [flags]

		

```
ignite testnet multi-node [flags]
```

**Options**

```
      --check-dependencies       verify that cached dependencies have not been modified since they were downloaded
      --clear-cache              clear the build cache (advanced)
  -h, --help                     help for multi-node
      --home string              directory where the blockchain node is initialized
      --node-dir-prefix string   prefix of dir node (default "validator")
  -p, --path string              path of the app (default ".")
  -r, --reset-once               reset the app state once on init
      --skip-proto               skip file generation from proto
  -v, --verbose                  verbose output
```

**SEE ALSO**

* [ignite testnet](#ignite-testnet)	 - Simulate and manage test networks


## ignite testnet simulate

Run simulation testing for the blockchain

**Synopsis**

Run simulation testing for the blockchain. It sends many randomized-input messages of each module to a simulated node.

```
ignite testnet simulate [flags]
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
      --seed int                  simulation random seed (default 42)
      --simName string            name of the simulation to run (default "TestFullAppSimulation")
```

**SEE ALSO**

* [ignite testnet](#ignite-testnet)	 - Simulate and manage test networks


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

# Scaffold Type

Ignites provides a set of scaffold types that can be used to generate code for your application.
These types are used in the `ignite scaffold` command.

## Available Scaffold Types

| Type | Usage |
| --- | --- |
| address | use '<FIELD_NAME>:address' to scaffold string types (eg: cosmos1abcdefghijklmnopqrstuvwxyz0123456). |
| array.coin | use '<FIELD_NAME>:array.coin' to scaffold sdk.Coins types (eg: 20stake). Disclaimer: Only one `coins` or `dec.coins` field can accept multiple CLI values per command due to AutoCLI limitations. |
| array.dec.coin | use '<FIELD_NAME>:array.dec.coin' to scaffold sdk.DecCoins types (eg: 20000002stake). Disclaimer: Only one `coins` or `dec.coins` field can accept multiple CLI values per command due to AutoCLI limitations. |
| array.int | use '<FIELD_NAME>:array.int' to scaffold []int64 types (eg: 5,4,3,2,1). |
| array.string | use '<FIELD_NAME>:array.string' to scaffold []string types (eg: abc,xyz). |
| array.uint | use '<FIELD_NAME>:array.uint' to scaffold []uint64 types (eg: 13,26,31,40). |
| bool | use '<FIELD_NAME>:bool' to scaffold bool types (eg: true). |
| bytes | use '<FIELD_NAME>:bytes' to scaffold []byte types (eg: 3,2,3,5). |
| coin | use '<FIELD_NAME>:coin' to scaffold sdk.Coin types (eg: 10token). |
| coins | use '<FIELD_NAME>:array.coin' to scaffold sdk.Coins types (eg: 20stake). Disclaimer: Only one `coins` or `dec.coins` field can accept multiple CLI values per command due to AutoCLI limitations. |
| custom | use the custom type to scaffold already created chain types. |
| dec.coin | use '<FIELD_NAME>:dec.coin' to scaffold sdk.DecCoin types (eg: 100001token). |
| dec.coins | use '<FIELD_NAME>:array.dec.coin' to scaffold sdk.DecCoins types (eg: 20000002stake). Disclaimer: Only one `coins` or `dec.coins` field can accept multiple CLI values per command due to AutoCLI limitations. |
| int | use '<FIELD_NAME>:int' to scaffold int64 types (eg: 111). |
| int64 | use '<FIELD_NAME>:int' to scaffold int64 types (eg: 111). |
| ints | use '<FIELD_NAME>:array.int' to scaffold []int64 types (eg: 5,4,3,2,1). |
| string | use '<FIELD_NAME>:string' to scaffold string types (eg: xyz). |
| strings | use '<FIELD_NAME>:array.string' to scaffold []string types (eg: abc,xyz). |
| uint | use '<FIELD_NAME>:uint' to scaffold uint64 types (eg: 111). |
| uint64 | use '<FIELD_NAME>:uint' to scaffold uint64 types (eg: 111). |
| uints | use '<FIELD_NAME>:array.uint' to scaffold []uint64 types (eg: 13,26,31,40). |


Field Usage:

    - fieldName
    - fieldName:fieldType


If no :fieldType, default (string) is used
