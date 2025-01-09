# Changelog

## Unreleased

### Features

- [#4457](https://github.com/ignite/cli/pull/4457) Add `skip-build` flag to `chain serve` command to avoid (re)building the chain

## [`v28.6.1`](https://github.com/ignite/cli/releases/tag/v28.6.1)

### Changes

- [#4449](https://github.com/ignite/cli/pull/4449) Bump scaffolded chain to Cosmos SDK `v0.50.11`. Previous version have a high security vulnerability.

## [`v28.6.0`](https://github.com/ignite/cli/releases/tag/v28.6.0)

### Features

- [#4377](https://github.com/ignite/cli/pull/4377) Add multi node (validator) testnet
- [#4362](https://github.com/ignite/cli/pull/4362) Scaffold `Makefile`

### Changes

- [#4376](https://github.com/ignite/cli/pull/4376) Set different chain-id for in place testnet

### Bug Fixes

- [#4421](https://github.com/ignite/cli/pull/4422) Fix typo in simulation template

## [`v28.5.3`](https://github.com/ignite/cli/releases/tag/v28.5.3)

### Changes

- [#4372](https://github.com/ignite/cli/pull/4372) Bump Cosmos SDK to `v0.50.10`
- [#4357](https://github.com/ignite/cli/pull/4357) Bump chain dependencies (store, ics, log, etc)
- [#4328](https://github.com/ignite/cli/pull/4328) Send ignite bug report to sentry. Opt out the same way as for usage analytics

## [`v28.5.2`](https://github.com/ignite/cli/releases/tag/v28.5.2)

### Features

- [#4297](https://github.com/ignite/cli/pull/4297) Add in-place testnet creation command for apps

### Changes

- [#4292](https://github.com/ignite/cli/pull/4292) Bump Cosmos SDK to `v0.50.9`
- [#4341](https://github.com/ignite/cli/pull/4341) Bump `ibc-go` to `8.5.0`
- [#4345](https://github.com/ignite/cli/pull/4345) Added survey link

### Fixes

- [#4319](https://github.com/ignite/cli/pull/4319) Remove fee abstraction module from open api code generation
- [#4309](https://github.com/ignite/cli/pull/4309) Fix chain id for chain simulations
- [#4322](https://github.com/ignite/cli/pull/4322) Create a message for authenticate buf for generate ts-client
- [#4323](https://github.com/ignite/cli/pull/4323) Add missing `--config` handling in the `chain` commands
- [#4350](https://github.com/ignite/cli/pull/4350) Skip upgrade prefix for sim tests

## [`v28.5.1`](https://github.com/ignite/cli/releases/tag/v28.5.1)

### Features

- [#4276](https://github.com/ignite/cli/pull/4276) Add `cosmosclient.CreateTxWithOptions` method to facilite more custom tx creation

### Changes

- [#4262](https://github.com/ignite/cli/pull/4262) Bring back relayer command
- [#4269](https://github.com/ignite/cli/pull/4269) Add custom flag parser for extensions
- [#4270](https://github.com/ignite/cli/pull/4270) Add flags to the extension hooks commands
- [#4286](https://github.com/ignite/cli/pull/4286) Add missing verbose mode flags

## [`v28.5.0`](https://github.com/ignite/cli/releases/tag/v28.5.0)

### Features

- [#4183](https://github.com/ignite/cli/pull/4183) Set `chain-id` in the client.toml
- [#4090](https://github.com/ignite/cli/pull/4090) Remove `protoc` pkg and also nodetime helpers `ts-proto` and `sta`
- [#4076](https://github.com/ignite/cli/pull/4076) Remove the ignite `relayer` and `tools` commands with all ts-relayer logic
- [#4133](https://github.com/ignite/cli/pull/4133) Improve buf rate limit

### Changes

- [#4095](https://github.com/ignite/cli/pull/4095) Migrate to matomo analytics
- [#4149](https://github.com/ignite/cli/pull/4149) Bump cometbft to `v0.38.7`
- [#4168](https://github.com/ignite/cli/pull/4168) Bump IBC to `v8.3.1`
  If you are upgrading manually from `v8.2.0` to `v8.3.1`, add the following to your `ibc.go` file:

  ```diff
  app.ICAHostKeeper = ...
  + app.ICAHostKeeper.WithQueryRouter(app.GRPCQueryRouter())`
  app.ICAControllerKeeper = ...
  ```

- [#4178](https://github.com/ignite/cli/pull/4178) Bump cosmos-sdk to `v0.50.7`
- [#4194](https://github.com/ignite/cli/pull/4194) Bump client/v2 to `v2.0.0-beta.2`
  If you are uprading manually, check out the recommended changes in `root.go` from the above PR.
- [#4210](https://github.com/ignite/cli/pull/4210) Improve default home wiring
- [#4077](https://github.com/ignite/cli/pull/4077) Merge the swagger files manually instead use nodetime `swagger-combine`
- [#4249](https://github.com/ignite/cli/pull/4249) Prevent creating a chain with number in the name
- [#4253](https://github.com/ignite/cli/pull/4253) Bump cosmos-sdk to `v0.50.8`

### Fixes

- [#4184](https://github.com/ignite/cli/pull/4184) Set custom `InitChainer` because of manually registered modules
- [#4198](https://github.com/ignite/cli/pull/4198) Set correct prefix overwriting in `buf.gen.pulsar.yaml`
- [#4199](https://github.com/ignite/cli/pull/4199) Set and seal SDK global config in `app/config.go`
- [#4212](https://github.com/ignite/cli/pull/4212) Set default values for extension flag to dont crash ignite
- [#4216](https://github.com/ignite/cli/pull/4216) Avoid create duplicated scopedKeppers
- [#4242](https://github.com/ignite/cli/pull/4242) Use buf build binary from the gobin path
- [#4250](https://github.com/ignite/cli/pull/4250) Set gas adjustment before calculating

## [`v28.4.0`](https://github.com/ignite/cli/releases/tag/v28.4.0)

### Features

- [#4108](https://github.com/ignite/cli/pull/4108) Add `xast` package (cherry-picked from [#3770](https://github.com/ignite/cli/pull/3770))
- [#4110](https://github.com/ignite/cli/pull/4110) Scaffold a consumer chain with `interchain-security` v5.0.0.
- [#4117](https://github.com/ignite/cli/pull/4117), [#4125](https://github.com/ignite/cli/pull/4125) Support relative path when installing local plugins

### Changes

- [#3959](https://github.com/ignite/cli/pull/3959) Remove app name prefix from the `.gitignore` file
- [#4103](https://github.com/ignite/cli/pull/4103) Bump cosmos-sdk to `v0.50.6`

### Fixes

- [#3969](https://github.com/ignite/cli/pull/3969) Get first config validator using a getter to avoid index errors
- [#4033](https://github.com/ignite/cli/pull/4033) Fix cobra completion using `fishshell`
- [#4062](https://github.com/ignite/cli/pull/4062) Avoid nil `scopedKeeper` in `TransmitXXX` functions
- [#4086](https://github.com/ignite/cli/pull/4086) Retry to get the IBC balance if it fails the first time
- [#4096](https://github.com/ignite/cli/pull/4096) Add new reserved names module and remove duplicated genesis order
- [#4112](https://github.com/ignite/cli/pull/4112) Remove duplicate SetCmdClientContextHandler
- [#4219](https://github.com/ignite/cli/pull/4219) Remove deprecated `sdk.MustSortJSON`

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
