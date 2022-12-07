---
order: 2
description: Use the Ignite CLI serve command to start your blockchain.
---

# Start a blockchain

Blockchains are decentralized applications.

- In production, blockchains often run the same software on many validator nodes that are run by different people and entities. To launch a blockchain in production, the validator entities coordinate the launch process to start their nodes simultaneously.
- During development, a blockchain can be started locally on a single validator node. This convenient process lets you restart a chain quickly and iterate faster. Starting a chain on a single node in development is similar to starting a traditional web application on a local server.

## Start a blockchain node in development

Switch to the directory that contains a blockchain that was scaffolded with Ignite CLI. To start the blockchain node, run the following command:

```bash
ignite chain serve
```

This command initializes a chain, builds the code, starts a single validator node, and starts watching for file changes.

Whenever a file is changed, the chain is automatically reinitialized, rebuilt, and started again. The chain's state is preserved if the changes to the source code are compatible with the previous state. This state preservation is beneficial for development purposes.

Because the `ignite chain serve` command is a development tool, it should not be used in a production environment. Read on to learn the process of running a blockchain in production.

## The Magic of `ignite chain serve`

The `ignite chain serve` command starts a fully operational blockchain.

The `ignite chain serve` command performs the following tasks:

- Installs dependencies
- Imports state, if possible
- Builds protocol buffer files
- Optionally generates TypeScript clients and Vuex stores 
- Builds a compiled blockchain binary
- Creates accounts
- Initializes a blockchain node
- Starts the following processes:
  - Tendermint RPC
  - Cosmos SDK API
  - Faucet, optional
- Watches for file changes and restarts
- Exports state

You can use flags to configure how the blockchain runs.

## Define how your blockchain starts

Flags for the `ignite chain serve` command determine how your blockchain starts. All flags are optional.

`--config`

Custom configuration file. Using unique configuration files is required to launch two blockchains on the same machine from the same source code. When omitted, the default is `config.yml`.

`--reset-once`

Reset the state only once. Use this flag to resume a failed reset or to initialize a blockchain from an empty state. The default state persistence imports the existing state and resumes the blockchain.

`--force-reset`

Reset state on every file change. Do not import state and turn off state persistence.

`--verbose`

Enter verbose detailed mode with extensive logging.

`--home`

Specify a custom home directory. 

## Start a blockchain node in production

The `ignite chain serve` and `ignite chain build` commands compile the source code of the chain in a binary file and install the binary in `~/go/bin`. By default, the binary name is the name of the repository appended with `d`. For example, if you scaffold a chain using `ignite scaffold chain mars`, then the binary is named `marsd`.

You can customize the binary name in `config.yml`:

```yaml
build:
  binary: "newchaind"
```

Or also add custom `ldflags` into your app binary:

```yaml
build:
  ldflags: [ "-X main.Env=prod", "-X main.Version=1.0.1" ]
```

Learn more about how to use the binary to [run a chain in production](https://docs.cosmos.network/main/run-node/run-node.html).
