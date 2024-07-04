---
sidebar_position: 3
description: Primary configuration file to describe the development environment for your blockchain.
title: config.yml reference
---

# config.yml reference

The `config.yml` file generated in your blockchain folder uses key-value pairs to describe the development environment for your blockchain.

Only a default set of parameters is provided. If more nuanced configuration is required, you can add these parameters to the `config.yml` file.

## accounts

A list of user accounts created during genesis of the blockchain.

| Key      | Required | Type            | Description                                                                                                                     |
| -------- | -------- | --------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| name     | Y        | String          | Local name of a key pair. An account name must be listed to gain access to the account tokens after the blockchain is launched. |
| coins    | Y        | List of Strings | Initial coins with denominations. For example, "1000token"                                                                      |
| address  | N        | String          | Account address in Bech32 address format.                                                                                        |
| mnemonic | N        | String          | Mnemonic used to generate an account. This field is ignored if `address` is specified.                                           |

Note that you can only use `address` OR `mnemonic` for an account. You can't use both, because an address is derived from a mnemonic.

If an account is a validator account (`alice` by default), it cannot have an `address` field.

**accounts example**

```yaml
accounts:
  - name: alice
    coins: ["1000token", "100000000stake"]
  - name: bob
    coins: ["500token"]
    address: cosmos1adn9gxjmrc3hrsdx5zpc9sj2ra7kgqkmphf8yw
```

## build

| Key      | Required | Type             | Description                                                                                                  |
|----------|----------|------------------|--------------------------------------------------------------------------------------------------------------|
| main     | N        | String           | When an app contains more than one main Go package, required to define the path of the chain's main package. |
| binary   | N        | String           | Name of the node binary that is built, typically ends with `d`.                                              |
| ldflags  | N        | List of Strings  | ldflags to set version information for go applications.                                                      |

**build example**

```yaml
build:
  binary: "mychaind"
  ldflags: [ "-X main.Version=development", "-X main.Date=01/05/2022T19:54" ]
```

### build.proto

| Key               | Required | Type            | Description                                                                                |
| ----------------- | -------- | --------------- | ------------------------------------------------------------------------------------------ |
| path              | N        | String          | Path to protocol buffer files. Default: `"proto"`.                                         |
| third_party_paths | N        | List of Strings | Path to third-party protocol buffer files. Default: `["third_party/proto", "proto_vendor"]`. |

## client

Configures and enables client code generation. To prevent Ignite CLI from regenerating the client, remove the `client` property.

### client.vuex

```yaml
client:
  vuex:
    path: "vue/src/store"
```

Generates Vuex stores for the blockchain in `path` on `serve` and `build` commands.

### client.typescript

```yaml
client:
  typescript:
    path: "vue/src/generated"
```

Generates TypeScript clients for the blockchain in `path` on `serve` and `build` commands.

### client.openapi

```yaml
client:
  openapi:
    path: "docs/static/openapi.yml"
```

Generates OpenAPI YAML file in `path`. By default this file is embedded in the node's binary.

## faucet

The faucet service sends tokens to addresses. The default address for the web user interface is <http://localhost:4500>.

| Key               | Required | Type            | Description                                                 |
| ----------------- | -------- | --------------- | ----------------------------------------------------------- |
| name              | Y        | String          | Name of a key pair. The `name` key pair must be in `accounts`.            |
| coins             | Y        | List of Strings | One or more coins with denominations sent per request.       |
| coins_max         | N        | List of Strings | One or more maximum amounts of tokens sent for each address. |
| host              | N        | String          | Host and port number. Default: `:4500`. Cannot be higher than 65536 |
| rate_limit_window | N        | String          | Time after which the token limit is reset (in seconds).      |

**faucet example**

```yaml
faucet:
  name: faucet
  coins: ["100token", "5foo"]
  coins_max: ["2000token", "1000foo"]
  port: 4500
```

## validator

A blockchain requires one or more validators.

| Key    | Required | Type   | Description                                                                                     |
| ------ | -------- | ------ | ----------------------------------------------------------------------------------------------- |
| name   | Y        | String | The account that is used to initialize the validator. The `name` key pair must be in `accounts`. |
| staked | Y        | String | Amount of coins to bond. Must be less than or equal to the amount of coins in the account.       |

**validator example**

```yaml
accounts:
  - name: alice
    coins: ["1000token", "100000000stake"]
validator:
  name: user1
  staked: "100000000stake"
```

## init.home

The path to the data directory that stores blockchain data and blockchain configuration.

**init example**

```yaml
init:
  home: "~/.myblockchain"
```

## init.config

Overwrites properties in `config/config.toml` in the data directory.

## init.app

Overwrites properties in `config/app.toml` in the data directory.

## init.client

Overwrites properties in `config/client.toml` in the data directory.

**init.client example**

```yaml
init:
  client:
    keyring-backend: "os"
```

## host

Configuration of host names and ports for processes started by Ignite CLI. Port numbers can't exceed 65536.

**host example**

```yaml
host:
  rpc: ":26659"
  p2p: ":26658"
  prof: ":6061"
  grpc: ":9091"
  api: ":1318"
```

## genesis

Use to overwrite values in `genesis.json` in the data directory to test different values in development environments. See [Genesis Overwrites for Development](../kb/04-genesis.md).
