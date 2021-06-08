---
order: 1
description: Use the Starport serve command to start your blockchain.
parent:
  order: 3
  title: Run a Blockchain
---

# Start a Blockchain

Blockchains are decentralized applications.

- In production, blockchains often run the same software on many validator nodes that are run by different people and entities. To launch a blockchain in production, the validator entities coordinate the launch process to start their nodes simultaneously.
- During development, a blockchain can be started locally on a single validator node. This convenient process lets you can restart a chain quickly and iterate faster. Starting a chain on a single node in development is similar to starting a traditional web application on a local server. 

## Start a Blockchain Node in Development

Switch to a directory that contains a blockchain scaffolded with Starport and run the following command:

```
starport serve
```

This command initializes a chain, builds the code, starts a single validator node and starts watching for file changes. Whenever a file is changed, Starport will reinitialise the chain, rebuild and start the chain again. Starport tries to preserve the state of the chain, if the changes to source code are compatible with the previous state. This is very useful for development purposes. `starport serve` command is a development tool, it should not be used in production environment. See below how to run a chain in production.

## The Magic of `starport serve`

The `starport serve` command performs the following tasks:

- Installs dependencies
- Imports state, if possible
- Builds protocol buffer files
- Optionally generates JavaScript (JS), TypeScript (TS), and Vuex clients
- Builds a compiled blockchain binary
- Creates accounts
- Initializes a blockchain node
- Starts the following processes:
  - Tendermint RPC
  - Cosmos SDK API
  - Faucet, optional
- Watches for file changes and restarts
- Exports state

The `starport serve` command starts a fully operational blockchain. You can use flags to configure how the blockchain runs. All flags are optional.

## Define How Your Blockchain Starts

Flags for the `starport serve` command determine how your blockchain starts:

`--config`, default is `config.yml`

Custom configuration file. Using unique configuration files is required to launch two blockchains on the same machine from the same source code.

`--reset-once`

Reset the state only once. Use this flag to resume a failed reset or to initialize a blockchain from an empty state. The default state persistence imports the existing state and resumes the blockchain.

`--force-reset`

Reset state on every file change. Do not import state and turn off state persistence.

`--rebuild-proto-once` use with `--reset-once`

Force code generation from proto files for custom and third-party modules. By default, Starport statically scaffolds files generated from Cosmos SDK standard proto files, instead of generating them dynamically. Use this flag to perform code generation on all modules if a blockchain was scaffolded on an earlier Starport version or after a Cosmos SDK upgrade.

`--verbose`

Enters verbose detailed mode with extensive logging.

`--home`

Specify a custom home directory.

## Start a Blockchain Node in Production

Both `starport serve` and `starport build` compile the source code of the chain in a binary file and install it in `~/go/bin`. By default, the binary name matches the name of the repository with a `d` suffix in the end. For example, if a chain was scaffolded using `starport app github.com/alice/chain`, then the binary will be named `chaind`.

The binary name can be customised in `config.yml`:

```yml
build:
  binary: "newchaind"
```

Learn more about how to use the binary to [run a chain in production](https://docs.cosmos.network/v0.42/run-node/run-node.html).
