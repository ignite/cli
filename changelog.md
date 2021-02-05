# Changelog

## `v0.14.0`

### Features:

- Chain state persistence between `starport serve` launches
- Integrated Stargate app's `scripts/protocgen` into Starport as a native feature. Running `starport build/serve` will automatically take care of building proto files without a need of script in the app's source code.
- Integrated third-party proto-files used by Cosmos SDK modules into Starport CLI
- Added ability to customize binary name with `build.binary` in `config.yml`
- Added ability to change path to home directory with `init.home` in `config.yml`
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

### Features:

- Added `starport network` commands for launching blockchains
- Added proxy (Chisel) to support launching blockchains from Gitpod
- Upgraded the template (Stargate) to Cosmos SDK v0.40.0-rc3
- Added a gRPC-Web proxy, which is available under http://localhost:12345/grpc.
- Added chain id configurability by recognizing `chain_id` from `genesis` section of `config.yml`.
- Added `config/app.toml` and `config/config.toml` configurability for appd under new `init.app` and `init.config` sections of `config.yml`.
- Point to Stargate as default SDK version for scaffolding.
- Covered CRUD operations for Stargate scaffolding.
- Added docs on gopath to build from source directions
- Arch Linux Based Raspberry Pi development environment
- Calculate the necessary gas for sending transactions to SPN

### Fixes:

- Routing REST API endpoints of querier on Stargate.
- Evaluate `--address-prefix` option when scaffolding for Stargate.
- Use a deterministic method to generate scaffolded type IDs
- Modify scaffolded type's creator type from address to string
- Copy built starport arm64 binary from tendermintdevelopment/starport:arm64 for device images
- Added git to amd64 docker image
- Comment out Gaia's seeds in the systemd unit template for downstream chains

## `v0.12.0`

### Features:

- Added Github CLI to gitpod environment for greater ease of use
- Added `starport build` command to build and install app binaries.
- Improved the first-time experience for readers of the Starport readme and parts of the Starport Handbook.
- Added `starport module create` command to scaffold custom modules
- Raspberry Pi now installs, builds, and serves the Vue UI
- Improved documentation for Raspberry Pi Device Images
- Added IBC and some other modules.
- Added an option to configure server addresses under `servers` section in `config.yml`.

### Fixes:

- `--address-prefix` will always be translated to lowercase while scaffolding with `app` command.
- HTTP API: accept strings in JSON and cast them to int and bool
- Update @tendermint/vue to `v0.1.7`
- Removed "Starport Pi"
- Removed Makefile from Downstream Pi
- Fixed Downstream Pi image Github Action
- Prevent duplicated fields with `type` command
- Fixed handling of protobufs profiler: prof_laddr -> pprof_laddr
- Fix an error, when a Stargate `serve` cmd doesn't start if a user doesn't have a relayer installed.

## `v0.11.1`

### Features:

- Published on Snapcraft.

## `v0.11.0`

### Features:

- Added experimental [Stargate](https://stargate.cosmos.network/) scaffolding option with `--sdk-version stargate` flag on `starport app` command.
- Pi Image Generation for chains generated with Starport
- Github action with capture of binary artifacts for chains generted with starport
- Gitpod: added guidelines and changed working directory into `docs`.
- Updated web scaffold with an improved sign in, balance list and a simple wallet.
- Added CRUD actions for scaffolded types: delete, update and get.

## `v0.0.10`

### Features:

- Add ARM64 releases.
- OS Image Generation for Raspberry Pi 3 and 4
- Added `version` command
- Added support for _validator_ configuration in _config.yml_.
- Starport can be launched on Gitpod
- Added `make clean`

### Fixes:

- Compile with go1.15
- Running `starport add type...` multiple times no longer breaks the app
- Running `appcli tx app create-x` now checks for all required args. -#173.
- Removed unused `--denom` flag from the `app` command. It previously has moved as a prop to the `config.yml` under `accounts` section.
- Disabled proxy server in the Vue app (this was causing to some compatibilitiy issues) and enabled CORS for `appcli rest-server` instead.
- `type` command supports dashes in app names.

## `v0.0.10-rc.3`

### Features:

- Configure `genesis.json` through `genesis` field in `config.yml`
- Initialize git repository on `app` scaffolding
- Check Go and GOPATH when running `serve`

### Changes:

- verbose is --verbose, not -v, in the cli
- Renamed `frontend` directory to `vue`
- Added first E2E tests (for `app` and `add wasm` subcommands)

### Fixes:

- No longer crashes, when git is initialized, but doesn't have commits
- Failure to start the frontend doesn't prevent Starport from running
- Changes to `config.yml` trigger reinitialization of the app
- Running `starport add wasm` multiple times no longer breaks the app

## `v0.0.10-rc.X`

### Features:

- Initialize with accounts defined `config.yml`
- `starport serve --verbose` shows detailed output from every process
- Custom address prefixes with `--address-prefix` flag
- Cosmos SDK Launchpad support
- Rebuild and reinitialize on file change

## `v0.0.9`

Initial release.
