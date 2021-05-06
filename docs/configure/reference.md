---
order: 2
description: Primary configuration file to describe the development environment for your blockchain.
---

# config.yml Reference

The `config.yml` file generated in your blockchain folder uses key-value pairs to describe the development environment for your blockchain.

Only a default set of parameters is provided. If more nuanced configuration is required, you can add these parameters to the `config.yml` file.


## `accounts`

A list of user accounts created during genesis of the blockc

| Key     | Required | Type            | Description                                                                                                                |
| ------- | -------- | --------------- | -------------------------------------------------------------------------------------------------------------------------- |
| name    | Y        | String          | Local name of a key pair. An account names must be listed to have access to their tokens after the blockchain is launched. |
| coins   | Y        | List of Strings | Initial coins with denominations. For example, "1000token"                                                                 |
| address | N        | String          | Account address in Bech32 address format                                                                                   |

**accounts example**

```yaml
accounts:
  - name: alice
    coins: ["1000token", "100000000stake"]
  - name: bob
    coins: ["500token"]
    address: cosmos1adn9gxjmrc3hrsdx5zpc9sj2ra7kgqkmphf8yw
```

## `build`

| Key    | Required | Type   | Description                                                    |
| ------ | -------- | ------ | -------------------------------------------------------------- |
| binary | N        | String | Name of the node binary that is built, typically ends with `d` |

**build example**

```yaml
build:
  binary: "mychaind"
```

## `build.proto`

| Key               | Required | Type            | Description                                                                                |
| ----------------- | -------- | --------------- | ------------------------------------------------------------------------------------------ |
| path              | N        | String          | Path to protocol buffer files. Default: `"proto"`                                          |
| third_party_paths | N        | List of Strings | Path to thid-party protocol buffer files. Default: `["third_party/proto", "proto_vendor"]` |

## `client`

Configures client code generation.

```
client:
  vuex:
    path: "vue/src/store"
```

Generates TypeScript/Vuex client for the blockchain in `path` on `serve` and `build` commands. Remove `client` property to prevent Starport from regenerating the client.

## `faucet`

The faucet service sends tokens to addresses. The default address for the web user interface is <http://localhost:4500>.

| Key       | Required | Type            | Description                                                 |
| --------- | -------- | --------------- | ----------------------------------------------------------- |
| name      | Y        | String          | Name of a key pair. `name` must be in `accounts`            |
| coins     | Y        | List of Strings | One or more coins with denominations sent per request       |
| coins_max | N        | List of Strings | One or more maximum amounts of tokens sent for each address |
| host      | N        | String          | Host and port number. Default: `:4500`                      |

**faucet example**

```yaml
faucet:
  name: faucet
  coins: ["100token", "5foo"]
  coins_max: ["2000token", "1000foo"]
  port: 4500
```

## `validator`

A blockchain requires one or more validators.

| Key    | Required | Type   | Description                                                                                     |
| ------ | -------- | ------ | ----------------------------------------------------------------------------------------------- |
| name   | Y        | String | The account that is used to initialize the validator. The `name` key pair must be in `accounts` |
| staked | Y        | String | Amount of coins to bond. Must be less than or equal to the amount of coins in the account       |

**validator example**

```yaml
accounts:
  - name: alice
    coins: ["1000token", "100000000stake"]
validator:
  name: user1
  staked: "100000000stake"
```

## `init.home`

The path to the data directory that stores blockchain data and blockchain configuration.

**init example**

```yaml
init:
  home: "~/.myblockchain"
```

## `init.config`

Overwrites properties in `config/config.toml` in the data directory.

## `init.app`

Overwrites properties in `config/app.toml` in the data directory.

## `init.keyring-backend`

The [keyring backend](https://docs.cosmos.network/master/run-node/keyring.html) to store the private key. Default value is `test`.

**init.keyring-backend example**

```yaml
init:
  keyring-backend: "os"
```

## `host`

Configuration of host names and ports for processes started by Starport:

**host example**

```
host:
  rpc: ":26659"
  p2p: ":26658"
  prof: ":6061"
  grpc: ":9091"
  api: ":1318"
  frontend: ":8081"
  dev-ui: ":12346"
```

## `genesis`

Use to overwrite values in `genesis.json` in the data directory to test different values in development environments. See [Genesis Overwrites for Development](https://docs.starport.network/configure/genesis.html).
