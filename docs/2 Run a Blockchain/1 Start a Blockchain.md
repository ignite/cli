# Start a Blockchain

One of your first actions is to start your blockchain. To start your blockchain, use the `starport serve` command.

The `starport serve` command performs the following tasks:

* Installs dependencies
* Imports state, if possible
* Builds protocol buffer files
* Optionally generates JavaScript (JS), TypeScript (TS), and Vuex clients
* Builds blockchain's binary
* Creates accounts
* Initializes a blockchain node
* Starts the following processes:
  * Tendermint RPC
  * Cosmos SDK API
  * Faucet (optional)
  * Welcome screen
  * Web scaffold (optional)
* Watches for file changes and restarts
* Exports state

Flags for the `starport serve` command determine how your blockchain starts:

- **--config** optional, default is `config.yml`

    Custom configuration file. Using unique configuration files is required to launch two blockchains on the same machine from the same source code. 

- **--reset-once** optional

    Reset the state only once. Use this flag to resume a failed reset or to initialize a blockchain from an empty state. The default state persistence imports the existing state and resumes the blockchain. 

Optional `--force-reset` flag: reset state on every file change, state persistence is turned off and Starport will not try to import state.

Optional `--rebuild-proto-once` flag: forces Starport to perform code generation from proto files for both custom and third-party modules. Use in combination with `--reset-once`. By default Starport scaffolds files generated from Cosmos SDK standard proto files, instead of generating them on the fly. Useful if you've scaffolded a blockchain with a previous version of Starport or you've upgraded Cosmos SDK version and you want to make sure that code generation is performed for all modules.

Optional `--verbose` flag: enters verbose (detailed) mode with extensive logging.

Optional `--home` flag: specify a custom home directory.
