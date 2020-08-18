# Changelog

## `develop`

### Changes:

* Renamed `frontend` directory to `vue`
* Added first E2E tests (for `app` and `add wasm` subcommands)

### Fixes:

* No longer crashes, when git is initialized, but doesn't have commits
* Failure to start the frontend doesn't prevent Starport from running
* Changes to `config.yml` trigger reinitialization of the app
* Running `starport add wasm` multiple times no longer breaks the app

## `v0.0.9`
