# Changelog

### Features:

* Pi Image Generation for chains generated with Starport
* Github action with capture of binary artifacts for for chains generted with starport


## `v0.0.10`

### Features:

* Added `version` command
* Added support for _validator_ configuration in _config.yml_.
* Starport can be launched on Gitpod

### Fixes:

* Running `starport add type...` multiple times no longer breaks the app
* Running `appcli tx app create-x` now checks for all required args. -#173.
* Removed unused `--denom` flag from the `app` command. It previously has moved as a prop to the `config.yml` under `accounts` section.
* Disabled proxy server in the Vue app (this was causing to some compatibilitiy issues) and enabled CORS for `appcli rest-server` instead.
* `type` command supports dashes in app names.


## `v0.0.10-rc.3`

### Features:

* Configure `genesis.json` through `genesis` field in `config.yml`
* Initialize git repository on `app` scaffolding
* Check Go and GOPATH when running `serve`

### Changes:

* verbose is --verbose, not -v, in the cli
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
