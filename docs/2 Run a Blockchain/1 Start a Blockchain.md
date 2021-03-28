# Start a Blockchain

One of the first things you might want to do is start your blockchain. You can start your blockchain with the `starport serve` command.

The command does the following:

* Installs dependencies
* Imports state, if possible
* Builds protocol buffer files
* Generates JS, TS and Vuex clients (optional)
* Builds blockchain's binary
* Creates accounts
* Initializes a blockchain node
* Starts processes
  * Tendermint RPC
  * Cosmos SDK API
  * Faucet (optional)
  * Welcome screen
  * Web scaffold (optional)
* Watches for file changes and restarts
* Exports state

Optional `--config` flag: specify a custom configuration file. By default value is `config.yml`. Useful for launching two blockchains on the same machine, from the same source code using two different configuration files.

Optional `--reset-once` flag: resets state once. By default Starport will try to import existing state and resume the blockchain (state persistence). If that fails or if you want to initialise a blockchain from an empty state, use `--reset-once`.

Optional `--force-reset` flag: reset state on every file change, state persistence is turned off and Starport will not try to import state.

Optional `--rebuild-proto-once` flag: forces Starport to perform code generation from proto files for both custom and third-party modules. Use in combination with `--reset-once`. By default Starport scaffolds files generated from Cosmos SDK standard proto files, instead of generating them on the fly. Useful if you've scaffolded a blockchain with a previous version of Starport or you've upgraded Cosmos SDK version and you want to make sure that code generation is performed for all modules.

Optional `--verbose` flag: enters verbose (detailed) mode with extensive logging.

Optional `--home` flag: specify a custom home directory.
