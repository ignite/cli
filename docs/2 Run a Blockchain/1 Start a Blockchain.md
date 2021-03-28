# Start a Blockchain

One of your first actions is to start your blockchain. To start your blockchain, use the `starport serve` command.

The `starport serve` command performs the following tasks:

* Installs dependencies
* Imports state, if possible
* Builds protocol buffer files
* Optionally generates JavaScript (JS), TypeScript (TS), and Vuex clients
* Builds a compiled blockchain binary
* Creates accounts
* Initializes a blockchain node
* Starts the following processes:
  * Tendermint RPC
  * Cosmos SDK API
  * Faucet, optional
  * Welcome screen
  * Web scaffold, optional
* Watches for file changes and restarts
* Exports state
The `starport serve` command starts a fully operational blockchain. You can use flags to configure how the blockchain runs. All flags are optional.
Flags for the `starport serve` command determine how your blockchain starts:

- **--config** optional, default is `config.yml`

    Custom configuration file. Using unique configuration files is required to launch two blockchains on the same machine from the same source code. 

- **--reset-once** optional

    Reset the state only once. Use this flag to resume a failed reset or to initialize a blockchain from an empty state. The default state persistence imports the existing state and resumes the blockchain. 

- **--force-reset** optional

    Reset state on every file change. Do not import state and turn off state persistence.

- **--rebuild-proto-once** use with `--reset-once`

    Force code generation from proto files for custom and third-party modules. By default, Starport statically scaffolds files generated from Cosmos SDK standard proto files, instead of generating them dynamically. Use this flag to perform code generation on all modules if a blockchain was scaffolded on an earlier Starport version or after a Cosmos SDK upgrade.

Optional `--verbose` flag: enters verbose (detailed) mode with extensive logging.

Optional `--home` flag: specify a custom home directory.
