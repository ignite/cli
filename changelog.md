# Changelog

## [`v0.20.3`](https://github.com/ignite-hq/cli/releases/tag/v0.20.3)

### Fixes

- Use latest version of CLI in templates to fix Linux ARM support _(It's now possible to develop chains in Linux ARM machines and since the chain depends on the CLI in its go.mod, it needs to use the latest version that support ARM targets)_

## [`v0.20.2`](https://github.com/ignite-hq/cli/releases/tag/v0.20.2)

### Fixes

- Use `unsafe-reset-all` cmd under `tendermint` cmd for chains that use `=> v0.45.3` version of Cosmos SDK

## [`v0.20.1`](https://github.com/ignite-hq/cli/releases/tag/v0.20.1)

### Features

- Release the CLI with Linux ARM and native M1 binaries

## [`v0.20.0`](https://github.com/ignite-hq/cli/releases/tag/v0.20.0)

Our new name is **Ignite CLI**!

**IMPORTANT!** This upgrade renames `starport` command to `ignite`. From now on, use `ignite` command to access the CLI.

### Features

- Upgraded Cosmos SDK version to `v0.45.2`
- Added support for in memory backend in `pkg/cosmosclient` package
- Improved our tutorials and documentation

## [`v0.19.5`](https://github.com/ignite-hq/cli/pull/2158/commits)

### Features

- Enable client code and Vuex code generation for query only modules as well.
- Upgraded the Vue template to `v0.3.5`.

### Fixes:
- Fixed snake case in code generation.
- Fixed plugin installations for Go =>v1.18.

### Changes:
- Dropped transpilation of TS to JS. Code generation now only produces TS files.

## `v0.19.4`

### Features

- Upgraded Vue template to `v0.3.0`.

## `v0.19.3`

### Features

- Upgraded Flutter template to `v2.0.3`

## [`v0.19.2`](https://github.com/ignite-hq/cli/milestone/14)

### Fixes

- Fixed race condition during faucet transfer
- Fixed account sequence mismatch issue on faucet and relayer
- Fixed templates for IBC code scaffolding

### Features

- Upgraded blockchain templates to use IBC v2.0.2

### Breaking Changes

- Deprecated the Starport Modules [tendermint/spm](https://github.com/tendermint/spm) repo and moved the contents to the Ignite CLI repo [`ignite/pkg/`](https://github.com/ignite-hq/cli/tree/develop/ignite/pkg/) in [PR 1971](https://github.com/ignite-hq/cli/pull/1971/files) 
 
    Updates are required if your chain uses these packages: 

    - `spm/ibckeeper` is now `pkg/cosmosibckeeper`
    - `spm/cosmoscmd` is now `pkg/cosmoscmd` 
    - `spm/openapiconsole` is now `pkg/openapiconsole`
    - `testutil/sample` is now `cosmostestutil/sample`

- Updated the faucet HTTP API schema. See API changes in [fix: improve faucet reliability #1974](https://github.com/ignite-hq/cli/pull/1974/files#diff-0e157f4f60d6fbd95e695764df176c8978d85f1df61475fbfa30edef62fe35cd)

## `v0.19.1`

### Fixes

- Enabled the `scaffold flutter` command

## `v0.19.0`

### Features

- `starport scaffold` commands support `ints`, `uints`, `strings`, `coin`, `coins` as field types (#1579)
- Added simulation testing with `simapp` to the default template (#1731)
- Added `starport generate dart` to generate a Dart client from protocol buffer files
- Added `starport scaffold flutter` to scaffold a Flutter mobile app template
- Parameters can be specified with a new `--params` flag when scaffolding modules (#1716)
- Simulations can be run with `starport chain simulate`
- Set `cointype` for accounts in  `config.yml` (#1663)

### Fixes

- Allow using a `creator` field when scaffolding a model with a `--no-message` flag (#1730)
- Improved error handling when generating code (#1907)
- Ensure account has funds after faucet transfer when using `cosmosclient` (#1846)
- Move from `io/ioutil` to `io` and `os` package (refactoring) (#1746)

## `v0.18.0`

### Breaking Changes

- Starport v0.18 comes with Cosmos SDK v0.44 that introduced changes that are not compatible with chains that were scaffolded with Starport versions lower than v0.18. After upgrading from Starport v0.17.3 to Starport v0.18, you must update the default blockchain template to use blockchains that were scaffolded with earlier versions. See [Migration](./docs/migration/index.md).

### Features:

- Scaffold commands allow using previously scaffolded types as fields
- Added `--signer` flag to `message`, `list`, `map`, and `single` scaffolding to allow customizing the name of the signer of the message
- Added `--index` flag to `scaffold map` to provide a custom list of indices
- Added `scaffold type` to scaffold a protocol buffer definition of a type
- Automatically check for new Starport versions
- Added `starport tools completions` to generate CLI completions
- Added `starport account` commands to manage accounts (key pairs)
- `starport version` now prints detailed information about OS, Go version, and more
- Modules are scaffolded with genesis validation tests
- Types are scaffolded with tests for `ValidateBasic` methods
- `cosmosclient` has been refactored and can be used as a library for interacting with Cosmos SDK chains
- `starport relayer` uses `starport account`
- Added `--path` flag for all `scaffold`, `generate` and `chain` commands
- Added `--output` flag to the `build` command
- Configure port of gRPC web in `config.yml` with the `host.grpc-web` property
- Added `build.main` field to `config.yml` for apps to specify the path of the chain's main package. This property is required to be set only when an app contains multiple main packages.

### Fixes

- Scaffolding a message no longer prevents scaffolding a map, list, or single that has the same type name when using the `--no-message` flag
- Generate Go code from proto files only from default directories or directories specified in `config.yml`
- Fixed faucet token transfer calculation
- Removed `creator` field for types scaffolded with the `--no-message` flag
- Encode the count value in the store with `BigEndian`

## `v0.17.3`

### Fixes

- oracle: add a specific BandChain pkg version to avoid Cosmos SDK version conflicts

## `v0.17.2`

### Features

- `client.toml` is initialized and used by node's CLI, can be configured through `config.yml` with the `init.client` property
- Support serving Cosmos SDK `v0.43.x` based chains

## `v0.17.1`

### Fixes

- Set visibility to `public` on Gitpod's port 7575 to enable peer discovery for SPN
- Fixed GitHub action that releases blockchain node's binary
- Fixed an error in chain scaffolding due to "unknown revision"
- Fixed an error in `starport chain serve` by limiting the scope where proto files are searched for

## `v0.17`

### Features

- Added GitHub action that automatically builds and releases a binary
- The `--release` flag for the `build` command adds the ability to release binaries in a tarball with a checksum file.
- Added the flag `--no-module` to the command `starport app` to prevent scaffolding a default module when creating a new app
- Added `--dep` flag to specify module dependency when scaffolding a module
- Added support for multiple naming conventions for component names and field names
- Print created and modified files when scaffolding a new component
- Added `starport generate` namespace with commands to generate Go, Vuex and OpenAPI
- Added `starport chain init` command to initialize a chain without starting a node
- Scaffold a type that contains a single instance in the store
- Introduced `starport tools` command for advanced users. Existing `starport relayer lowlevel *` commands are also moved under `tools`
- Added `faucet.rate_limit_window` property to `config.yml`
- Simplified the `cmd` package in the template
- Added `starport scaffold band` oracle query scaffolding
- Updated TypeScript relayer to 0.2.0
- Added customizable gas limits for the relayer

### Fixes

- Use snake case for generated files
- Prevent using incorrect module name
- Fixed permissions issue when using Starport in Docker
- Ignore hidden directories when building a chain
- Fix error when scaffolding an IBC module in non-Starport chains

## `v0.16.2`

### Fix

- Prevent indirect Buf dependency

## `v0.16.1`

### Features

- Ensure that CLI operates fine even if the installation directory (bin) of Go programs is not configured properly

## `v0.16.0`

### Features

- The new `join` flag adds the ability to pass a `--genesis` file and `--peers` address list with `starport network chain join`
- The new `show` flag adds the ability to show `--genesis` and `--peers` list with `starport network chain show`
- `protoc` is now bundled with Ignite CLI. You don't need to install it anymore.
- Starport is now published automatically on the Docker Hub
- `starport relayer` `configure` and `connect` commands now use the [confio/ts-relayer](https://github.com/confio/ts-relayer) under the hood. Also, checkout the new `starport relayer lowlevel` command
- An OpenAPI spec for your chain is now automatically generated with `serve` and `build` commands: a console is available at `localhost:1317` and spec at `localhost:1317/static/openapi.yml` by default for the newly scaffolded chains
- Keplr extension is supported on web apps created with Starport
- Added tests to the scaffold
- Improved reliability of scaffolding by detecting placeholders
- Added ability to scaffold modules in chains not created with Starport
- Added the ability to scaffold Cosmos SDK queries
- IBC relayer support is available on web apps created with Starport
- New types without CRUD operations can be added with the `--no-message` flag in the `type` command
- New packet without messages can be added with the `--no-message` flag in the `packet` command
- Added `docs` command to read Starport documentation on the CLI
- Published documentation on https://docs.starport.network
- Added `mnemonic` property to account in the `accounts` list to generate a key from a mnemonic

### Fixes

- `starport network chain join` hanging issue when creating an account
- Error when scaffolding a chain with an underscore in the repo name (thanks @bensooraj!)

### Changes

- `starport serve` no longer starts the web app in the `vue` directory (use `npm` to start it manually)
- Default scaffold no longer includes legacy REST API endpoints (thanks @bensooraj!)
- Removed support for Cosmos SDK v0.39 Launchpad

## `v0.15.0`

### Features

- IBC module scaffolding
- IBC packet scaffolding with acknowledgements
- JavaScript and Vuex client code generation for Cosmos SDK and custom modules
- Standalone relayer with `configure` and `connect` commands
- Advanced relayer options for configuring ports and versions
- Scaffold now follows `MsgServer` convention
- Message scaffolding
- Added `starport type ... --indexed` to scaffold indexed types
- Custom config file support with `starport serve -c custom.yml`
- Detailed terminal output for created accounts: name, address, mnemonic
- Added spinners to indicate progress for long-running commands
- Updated to Cosmos SDK v0.42.1

### Changes

- Replaced `packr` with Go 1.16 `embed`
- Renamed `servers` top-level property to `host`

## `v0.14.0`

### Features

- Chain state persistence between `starport serve` launches
- Integrated Stargate app's `scripts/protocgen` into Starport as a native feature. Running `starport build/serve` will automatically take care of building proto files without a need of script in the app's source code.
- Integrated third-party proto-files used by Cosmos SDK modules into Ignite CLI
- Added ability to customize binary name with `build.binary` in `config.yml`
- Added ability to change path to home directory with `
.home` in `config.yml`
- Added ability to add accounts by `address` with in `config.yml`
- Added faucet functionality available on port 4500 and configurable with `faucet` in `config.yml`
- Added `starport faucet [address] [coins]` command
- Updated scaffold to Cosmos SDK v0.41.0
- Distroless multiplatform docker containers for starport that can be used for `starport serve`
- UI containers for chains scaffolded with Starport
- Use SOS-lite and Docker instead of systemD
- Arch PKGBUILD in `scripts`

### Fixes:

- Support for CosmWasm on Stargate
- Bug with dashes in Github username breaking proto package name
- Bug with custom address prefix
- use docker buildx as a single command with multiple platforms to make multi-manifest work properly

## `v0.13.0`

### Features

- Added `starport network` commands for launching blockchains
- Added proxy (Chisel) to support launching blockchains from Gitpod
- Upgraded the template (Stargate) to Cosmos SDK v0.40.0-rc3
- Added a gRPC-Web proxy that is available under http://localhost:12345/grpc
- Added chain id configurability by recognizing `chain_id` from `genesis` section of `config.yml`.
- Added `config/app.toml` and `config/config.toml` configurability for appd under new `init.app` and `init.config` sections of `config.yml`
- Point to Stargate as default SDK version for scaffolding
- Covered CRUD operations for Stargate scaffolding
- Added docs on gopath to build from source directions
- Arch Linux Based Raspberry Pi development environment
- Calculate the necessary gas for sending transactions to SPN

### Fixes

- Routing REST API endpoints of querier on Stargate
- Evaluate `--address-prefix` option when scaffolding for Stargate
- Use a deterministic method to generate scaffolded type IDs
- Modify scaffolded type's creator type from address to string
- Copy built starport arm64 binary from tendermintdevelopment/starport:arm64 for device images
- Added git to amd64 docker image
- Comment out Gaia's seeds in the systemd unit template for downstream chains

## `v0.12.0`

### Features

- Added Github CLI to gitpod environment for greater ease of use
- Added `starport build` command to build and install app binaries
- Improved the first-time experience for readers of the Starport readme and parts of the Starport Handbook
- Added `starport module create` command to scaffold custom modules
- Raspberry Pi now installs, builds, and serves the Vue UI
- Improved documentation for Raspberry Pi Device Images
- Added IBC and some other modules
- Added an option to configure server addresses under `servers` section in `config.yml`

### Fixes

- `--address-prefix` will always be translated to lowercase while scaffolding with `app` command
- HTTP API: accept strings in JSON and cast them to int and bool
- Update @tendermint/vue to `v0.1.7`
- Removed "Starport Pi"
- Removed Makefile from Downstream Pi
- Fixed Downstream Pi image Github Action
- Prevent duplicated fields with `type` command
- Fixed handling of protobufs profiler: prof_laddr -> pprof_laddr
- Fix an error, when a Stargate `serve` cmd doesn't start if a user doesn't have a relayer installed

## `v0.11.1`

### Features

- Published on Snapcraft

## `v0.11.0`

### Features

- Added experimental [Stargate](https://stargate.cosmos.network/) scaffolding option with `--sdk-version stargate` flag on `starport app` command
- Pi Image Generation for chains generated with Starport
- Github action with capture of binary artifacts for chains generated with Starport
- Gitpod: added guidelines and changed working directory into `docs`
- Updated web scaffold with an improved sign in, balance list and a simple wallet
- Added CRUD actions for scaffolded types: delete, update, and get

## `v0.0.10`

### Features

- Add ARM64 releases
- OS Image Generation for Raspberry Pi 3 and 4
- Added `version` command
- Added support for _validator_ configuration in _config.yml_.
- Starport can be launched on Gitpod
- Added `make clean`

### Fixes

- Compile with go1.15
- Running `starport add type...` multiple times no longer breaks the app
- Running `appcli tx app create-x` now checks for all required args
- Removed unused `--denom` flag from the `app` command. It previously has moved as a prop to the `config.yml` under `accounts` section
- Disabled proxy server in the Vue app (this was causing to some compatibilitiy issues) and enabled CORS for `appcli rest-server` instead
- `type` command supports dashes in app names

## `v0.0.10-rc.3`

### Features

- Configure `genesis.json` through `genesis` field in `config.yml`
- Initialize git repository on `app` scaffolding
- Check Go and GOPATH when running `serve`

### Changes

- verbose is --verbose, not -v, in the cli
- Renamed `frontend` directory to `vue`
- Added first E2E tests (for `app` and `add wasm` subcommands)

### Fixes

- No longer crashes when git is initialized but doesn't have commits
- Failure to start the frontend doesn't prevent Starport from running
- Changes to `config.yml` trigger reinitialization of the app
- Running `starport add wasm` multiple times no longer breaks the app

## `v0.0.10-rc.X`

### Features

- Initialize with accounts defined `config.yml`
- `starport serve --verbose` shows detailed output from every process
- Custom address prefixes with `--address-prefix` flag
- Cosmos SDK Launchpad support
- Rebuild and reinitialize on file change

## `v0.0.9`

Initial release.
