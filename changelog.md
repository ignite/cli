# Changelog

## Unreleased

### Features

- [#4108](https://github.com/ignite/cli/pull/4108) Add `xast` package (cherry-picked from [#3770](https://github.com/ignite/cli/pull/3770))

### Changes

- [#3959](https://github.com/ignite/cli/pull/3959) Remove app name prefix from the `.gitignore` file

### Fixes

- [#4033](https://github.com/ignite/cli/pull/4033) Fix cobra completion using `fishshell`
- [#4062](https://github.com/ignite/cli/pull/4062) Avoid nil `scopedKeeper` in `TransmitXXX` functions
- [#3969](https://github.com/ignite/cli/pull/3969) Get first config validator using a getter to avoid index errors

## [`v28.3.0`](https://github.com/ignite/cli/releases/tag/v28.3.0)

### Features

- [#4019](https://github.com/ignite/cli/pull/4019) Add `skip-proto` flag to `s chain` command
- [#3985](https://github.com/ignite/cli/pull/3985) Make some `cmd` pkg functions public
- [#3956](https://github.com/ignite/cli/pull/3956) Prepare for wasm app
- [#3660](https://github.com/ignite/cli/pull/3660) Add ability to scaffold ICS consumer chain

### Changes

- [#4035](https://github.com/ignite/cli/pull/4035) Bump `cometbft` to `v0.38.6` and `ibc-go/v8` to `v8.1.1`
- [#4031](https://github.com/ignite/cli/pull/4031) Bump `cli-plugin-network` to `v0.2.2` due to dependencies issue.
- [#4013](https://github.com/ignite/cli/pull/4013) Bump `cosmos-sdk` to `v0.50.5`
- [#4010](https://github.com/ignite/cli/pull/4010) Use `AppName` instead `ModuleName` for scaffold a new App
- [#3972](https://github.com/ignite/cli/pull/3972) Skip Ignite app loading for some base commands that don't allow apps
- [#3983](https://github.com/ignite/cli/pull/3983) Bump `cosmos-sdk` to `v0.50.4` and `ibc-go` to `v8.1.0`

### Fixes

- [#4021](https://github.com/ignite/cli/pull/4021) Set correct custom signer in `s list --signer <signer>`
- [#3995](https://github.com/ignite/cli/pull/3995) Fix interface check for ibc modules
- [#3953](https://github.com/ignite/cli/pull/3953) Fix apps `Stdout` is redirected to `Stderr`
- [#3863](https://github.com/ignite/cli/pull/3963) Fix breaking issue for app client API when reading app chain info

## [`v28.2.0`](https://github.com/ignite/cli/releases/tag/v28.2.0)

### Features

- [#3924](https://github.com/ignite/cli/pull/3924) Scaffold NFT module by default
- [#3839](https://github.com/ignite/cli/pull/3839) New structure for app scaffolding
- [#3835](https://github.com/ignite/cli/pull/3835) Add `--minimal` flag to `scaffold chain` to scaffold a chain with the least amount of sdk modules
- [#3820](https://github.com/ignite/cli/pull/3820) Add integration tests for IBC chains

### Changes

- [#3899](https://github.com/ignite/cli/pull/3899) Introduce `plugin.Execute` function
- [#3903](https://github.com/ignite/cli/pull/3903) Don't specify a default build tag and deprecate notion of app version

### Fixes

- [#3905](https://github.com/ignite/cli/pull/3905) Fix `ignite completion`
- [#3931](https://github.com/ignite/cli/pull/3931) Fix `app update` command and duplicated apps

## [`v28.1.1`](https://github.com/ignite/cli/releases/tag/v28.1.1)

### Fixes

- [#3878](https://github.com/ignite/cli/pull/3878) Support local forks of Cosmos SDK in scaffolded chain.
- [#3869](https://github.com/ignite/cli/pull/3869) Fix .git in parent dir
- [#3867](https://github.com/ignite/cli/pull/3867) Fix genesis export for ibc modules.
- [#3850](https://github.com/ignite/cli/pull/3871) Fix app.go file detection in apps scaffolded before v28.0.0

### Changes

- [#3885](https://github.com/ignite/cli/pull/3885) Scaffold chain with Cosmos SDK `v0.50.3`
- [#3877](https://github.com/ignite/cli/pull/3877) Change Ignite App extension to "ign"
- [#3897](https://github.com/ignite/cli/pull/3897) Introduce alternative folder in templates

## [`v28.1.0`](https://github.com/ignite/cli/releases/tag/v28.1.0)

### Features

- [#3786](https://github.com/ignite/cli/pull/3786) Add artifacts for publishing Ignite to FlatHub and Snapcraft
- [#3830](https://github.com/ignite/cli/pull/3830) Remove gRPC info from Ignite Apps errors
- [#3861](https://github.com/ignite/cli/pull/3861) Send to the analytics if the user is using a GitPod

### Changes

- [#3822](https://github.com/ignite/cli/pull/3822) Improve default scaffolded AutoCLI config
- [#3838](https://github.com/ignite/cli/pull/3838) Scaffold chain with Cosmos SDK `v0.50.2`, and bump confix and x/upgrade to latest
- [#3829](https://github.com/ignite/cli/pull/3829) Support version prefix for cached values
- [#3723](https://github.com/ignite/cli/pull/3723) Create a wrapper for errors

### Fixes

- [#3827](https://github.com/ignite/cli/pull/3827) Change ignite apps to be able to run in any directory
- [#3831](https://github.com/ignite/cli/pull/3831) Correct ignite app gRPC server stop memory issue
- [#3825](https://github.com/ignite/cli/pull/3825) Fix a minor Keplr type-checking bug in TS client
- [#3836](https://github.com/ignite/cli/pull/3836), [#3858](https://github.com/ignite/cli/pull/3858) Add missing IBC commands for scaffolded chain
- [#3833](https://github.com/ignite/cli/pull/3833) Improve Cosmos SDK detection to support SDK forks
- [#3849](https://github.com/ignite/cli/pull/3849) Add missing `tx.go` file by default and enable cli if autocli does not exist
- [#3851](https://github.com/ignite/cli/pull/3851) Add missing ibc interfaces to chain client
- [#3860](https://github.com/ignite/cli/pull/3860) Fix analytics event name

## [`v28.0.0`](https://github.com/ignite/cli/releases/tag/v28.0.0)

### Features

- [#3659](https://github.com/ignite/cli/pull/3659) cosmos-sdk `v0.50.x` upgrade
- [#3694](https://github.com/ignite/cli/pull/3694) Query and Tx AutoCLI support
- [#3536](https://github.com/ignite/cli/pull/3536) Change app.go to v2 and add AppWiring feature
- [#3544](https://github.com/ignite/cli/pull/3544) Add bidirectional communication to app (plugin) system
- [#3756](https://github.com/ignite/cli/pull/3756) Add faucet compatibility for latest sdk chains
- [#3476](https://github.com/ignite/cli/pull/3476) Use `buf.build` binary to code generate from proto files
- [#3724](https://github.com/ignite/cli/pull/3724) Add or vendor proto packages from Go dependencies
- [#3561](https://github.com/ignite/cli/pull/3561) Add GetChainInfo method to plugin system API
- [#3626](https://github.com/ignite/cli/pull/3626) Add logging levels to relayer
- [#3614](https://github.com/ignite/cli/pull/3614) feat: use DefaultBaseappOptions for app.New method
- [#3715](https://github.com/ignite/cli/pull/3715) Add test suite for the cli tests

### Changes

- [#3793](https://github.com/ignite/cli/pull/3793) Refactor Ignite to follow semantic versioning (prepares v28.0.0). If you are using packages, do not forget to import the `/v28` version of the packages.
- [#3529](https://github.com/ignite/cli/pull/3529) Refactor plugin system to use gRPC
- [#3751](https://github.com/ignite/cli/pull/3751) Rename label to skip changelog check
- [#3745](https://github.com/ignite/cli/pull/3745) Set tx fee amount as option
- [#3748](https://github.com/ignite/cli/pull/3748) Change default rpc endpoint to a working one
- [#3621](https://github.com/ignite/cli/pull/3621) Change `pkg/availableport` to allow custom parameters in `Find` function and handle duplicated ports
- [#3810](https://github.com/ignite/cli/pull/3810) Bump network app version to `v0.2.1`
- [#3581](https://github.com/ignite/cli/pull/3581) Bump cometbft and cometbft-db in the template
- [#3522](https://github.com/ignite/cli/pull/3522) Remove indentation from `chain serve` output
- [#3346](https://github.com/ignite/cli/issues/3346) Improve scaffold query --help
- [#3601](https://github.com/ignite/cli/pull/3601) Update ts-relayer version to `0.10.0`
- [#3658](https://github.com/ignite/cli/pull/3658) Rename Marshaler to Codec in EncodingConfig
- [#3653](https://github.com/ignite/cli/pull/3653) Add "app" extension to plugin binaries
- [#3656](https://github.com/ignite/cli/pull/3656) Disable Go toolchain download
- [#3662](https://github.com/ignite/cli/pull/3662) Refactor CLI "plugin" command to "app"
- [#3669](https://github.com/ignite/cli/pull/3669) Rename `plugins` config file to `igniteapps`
- [#3683](https://github.com/ignite/cli/pull/3683) Resolve `--dep auth` as `--dep account` in `scaffold module`
- [#3795](https://github.com/ignite/cli/pull/3795) Bump cometbft to `v0.38.2`
- [#3599](https://github.com/ignite/cli/pull/3599) Add analytics as an option
- [#3670](https://github.com/ignite/cli/pull/3670) Remove binaries

### Fixes

- [#3386](https://github.com/ignite/cli/issues/3386) Prevent scaffolding of default module called "ibc"
- [#3592](https://github.com/ignite/cli/pull/3592) Fix `pkg/protoanalysis` to support HTTP rule parameter arguments
- [#3598](https://github.com/ignite/cli/pull/3598) Fix consensus param keeper constructor key in `app.go`
- [#3610](https://github.com/ignite/cli/pull/3610) Fix overflow issue of cosmos faucet in `pkg/cosmosfaucet/transfer.go` and `pkg/cosmosfaucet/cosmosfaucet.go`
- [#3618](https://github.com/ignite/cli/pull/3618) Fix TS client generation import path issue
- [#3631](https://github.com/ignite/cli/pull/3631) Fix unnecessary vue import in hooks/composables template
- [#3661](https://github.com/ignite/cli/pull/3661) Change `pkg/cosmosanalysis` to find Cosmos SDK runtime app registered modules
- [#3716](https://github.com/ignite/cli/pull/3716) Fix invalid plugin hook check
- [#3725](https://github.com/ignite/cli/pull/3725) Fix flaky TS client generation issues on linux
- [#3726](https://github.com/ignite/cli/pull/3726) Update TS client dependencies. Bump vue/react template versions
- [#3728](https://github.com/ignite/cli/pull/3728) Fix wrong parser for proto package names
- [#3729](https://github.com/ignite/cli/pull/3729) Fix broken generator due to caching /tmp include folders
- [#3767](https://github.com/ignite/cli/pull/3767) Fix `v0.50` ibc genesis issue
- [#3808](https://github.com/ignite/cli/pull/3808) Correct TS code generation to generate paginated fields
