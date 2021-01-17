[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/tendermint/starport)

# Starport

Starport is the easiest way to build blockchains. It is a developer-friendly interface to the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk), the world's most widely-used blockchain application framework. Starport generates boilerplate code for you, so you can focus on writing business logic.

![Banner](./docs/banner.jpeg)

Overview: https://www.youtube.com/watch?v=rmbPjCGDXek

## Install

### NPM

```
npm i -g @tendermint/starport
```

### macOS with Homebrew

```
brew install tendermint/tap/starport
```

### Build from source

```
git clone https://github.com/tendermint/starport && cd starport && make
```

Requirements: Go 1.14 and Node.js (optional but highly recommended, used for the welcome screen and web app scaffolding).

## Get started

### Create an application

This command creates an empty template for a Cosmos SDK application. By default it also includes a module with the same name as the package. To create a new application called `blog`, run:

```
starport app github.com/your-github-id/blog
```

| Flag               | Default  | Description                |
| ------------------ | -------- | -------------------------- |
| `--address-prefix` | `cosmos` | Prefix, used for addresses |

### Run an application

```
starport serve
```

To start the server, go into you application's directory and run `starport serve`. This commands installs dependencies, builds and initializes the app and runs both Tendermint RPC server (by default on `localhost:26657`) as well as LCD (by default on `localhost:1317`) with hot reloading enabled.

`starport serve` uses `config.yml` to initialize your application, make sure you have it in your project directory (see [Configure](#configure)).

Note: depending on your OS and firewall settings, you may have to accept a prompt asking if your application's binary (`blogd` in this case) can accept external connections.

| Flag        | Default | Description                          |
| ----------- | ------- | ------------------------------------ |
| `--verbose` | `false` | Enable verbose output from processes |
| `--path`    |         | Path to the project                  |

### Create data types

```
starport type [typeName] [field1] [field2:bool] ...
```

This command generates messages, handlers, keepers, CLI and REST clients and type definition for `typeName` type. A type can have any number of `field` arguments. By default fields are strings, but `bool` and `int` are supported.

For example,

```
starport type post title body
```

This command generates a type `Post` with two fields: `title` and `body`.

To add a post run `blogcli tx blog create-post "My title" "This is a blog" --from=user1`.


### Configure

Initialization parameters of your app are stored in `config.yml`.

The simple configuration file includes a list of accounts and their initial coins:

```
version: 1
accounts:
  - name: user1
    coins: ["1000token", "100000000stake"]
  - name: user2
    coins: ["500token"]
validator:
  name: user1
  staked: "100000000stake"
```

#### `accounts`

A list of user accounts created during genesis of your application.

| Key   | Required | Type            | Description                                       |
| ----- | -------- | --------------- | ------------------------------------------------- |
| name  | Y        | String          | Local name of the key pair                        |
| coins | Y        | List of Strings | Initial coins with denominations (e.g. "100coin") |

#### `validator`

A property that describes your local validator. `name` should be one of the names, specified in the `accounts` array. The account should have enough tokens for staking purposes.

| Key    | Required | Type   | Description                                                                         |
| ------ | -------- | ------ | ----------------------------------------------------------------------------------- |
| name   | Y        | String | Name of one the accounts                                                            |
| staked | Y        | String | Amount of coins staked by your validator, should be >= 10^6 (e.g. "100000000stake") |

### Add smart contract support

```
starport add wasm
```

Adds smart contracts with [CosmWasm](https://docs.cosmwasm.com). Follow a short [smart contract tutorial](https://www.notion.so/Smart-contracts-with-CosmWasm-c6fbcd584b78437a843e738b922dc108) to get started.

## More tutorials

- [Blog (video) tutorial](https://www.youtube.com/watch?v=rmbPjCGDXek): get started with your first blockchain
- [Poll tutorial](https://tutorials.cosmos.network/starport-polling-app/): build a voting application with a web-based UI
- [Smart contract tutorial](https://www.notion.so/Smart-contracts-with-CosmWasm-c6fbcd584b78437a843e738b922dc108): add smart contracts to your app with CosmWasm: build, upload, instantiate and run a smart contract
- [Blog (from scratch) tutorial](https://tutorials.cosmos.network/starport-blog/01-index.html): learn how Starport works by building a blog without scaffolding

## Questions & comments

For questions and support please join the #starport channel in the [Cosmos Community Discord](https://discord.com/invite/W8trcGV). The issue list of this repo is exclusively for bug reports and feature requests.

## Contributing

`develop` contains the latest version. Please branch of from `develop` and create a PR.

## Stay in touch

Starport is a product built by [Tendermint](https://tendermint.com). Follow us to get the latest updates!

- [Twitter](https://twitter.com/tendermint_team)
- [Blog](https://medium.com/tendermint)
- [Jobs](https://tendermint.com/careers)
