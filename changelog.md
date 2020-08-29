# Changelog

## `develop`

### Features:

* Add ARM64 releases.
* OS Image Generation for Raspberry Pi 3 and 4
* Added `version` command
* Added support for _validator_ configuration in _config.yml_.

### Fixes:

* Compile with go1.15
* Running `starport add type...` multiple times no longer breaks the app
* Running `appcli tx app create-x` now checks for all required args. -#173.
* Removed unused `--denom` flag from the `app` command. It previously has moved as a prop to the `config.yml` under `accounts` section.


## `v0.0.10-rc.3`

### Features:

* Configure `genesis.json` through `genesis` field in `config.yml`
* Initialize git repository on `app` scaffolding
* Check Go and GOPATH when running `serve`

### Changes:

* Renamed `frontend` directory to `vue`
* Added first E2E tests (for `app` and `add wasm` subcommands)

### Fixes:

* No longer crashes, when git is initialized, but doesn't have commits
* Failure to start the frontend doesn't prevent Starport from running
* Changes to `config.yml` trigger reinitialization of the app
* Running `starport add wasm` multiple times no longer breaks the app

## `v0.0.10-rc.X`

### Features:

* Initialize with accounts defined `config.yml`
* `starport serve --verbose` shows detailed output from every process
* Custom address prefixes with `--address-prefix` flag
* Cosmos SDK Launchpad support
* Rebuild and reinitialize on file change

## `v0.0.9`

Initial release.
