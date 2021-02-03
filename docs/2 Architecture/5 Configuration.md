# Configuration

With Starport your blockchain can be configured with `config.yml`.

## `accounts`

A list of user accounts created during genesis of your application.

| Key     | Required | Type            | Description                                       |
| ------- | -------- | --------------- | ------------------------------------------------- |
| name    | Y        | String          | Local name of a key pair                          |
| coins   | Y        | List of Strings | Initial coins with denominations (e.g. "100coin") |
| address | N        | String          | Address of the account in bech32                  |

### Example

```yaml
accounts:
  - name: alice
    coins: ["1000token", "100000000stake"]
  - name: bob
    coins: ["500token"]
    address: cosmos1adn9gxjmrc3hrsdx5zpc9sj2ra7kgqkmphf8yw
```

## `build`

| Key    | Required | Type   | Description                                                     |
| ------ | -------- | ------ | --------------------------------------------------------------- |
| binary | N        | String | Name of the node binary that will be built and used by Starport |

## Example

```yaml
build:
  binary: "mychaind"
```

## `build.proto`

| Key               | Required | Type            | Description                                                                                 |
| ----------------- | -------- | --------------- | ------------------------------------------------------------------------------------------- |
| path              | N        | String          | Path to protocol buffer files (default: `"proto"`)                                          |
| third_party_paths | N        | List of Strings | Path to thid-party protocol buffer files (default: `["third_party/proto", "proto_vendor"]`) |


## `faucet`

A faucet service that sends tokens to addresses. Web UI is available by default on http://localhost:4500.

| Key       | Required | Type            | Description                                      |
| --------- | -------- | --------------- | ------------------------------------------------ |
| name      | Y        | String          | Name of a key pair. `name` must be in `accounts` |
| coins     | Y        | List of Strings | Coins with denominations sent per request        |
| coins_max | N        | List of Strings | Maximum amount of tokens sent per address        |
| port      | N        | Integer         | Port number (default: `4500`)                    |

### Example

```yaml
faucet:
  name: faucet
  coins: ["100token", "5foo"]
  coins_max: ["2000token", "1000foo"]
  port: 4500
```

## `validator`

A blockchain has to have at least one validator-node. `validator` specifies the account that will be used to initialize the validator and parameters of the validator.

| Key    | Required | Type   | Description                                                                       |
| ------ | -------- | ------ | --------------------------------------------------------------------------------- |
| name   | Y        | String | Name of a key pair. `name` must be in `accounts`                                  |
| staked | Y        | String | Amount of coins to bond. Must be less or equal to the amount of coins account has |

### Example

```yaml
accounts:
  - name: alice
    coins: ["1000token", "100000000stake"]
validator:
  name: user1
  staked: "100000000stake"
```

## `init.home`

A blockchain stores data and configuration in a data directory. This property specifies a path to the data directory.

### Example

```yaml
init:
  home: "~/.myblockchain"
```

## `init.config`

Overwrites properties in `config/config.toml` in the data directory.

## `init.app`

Overwrites properties in `config/app.toml` in the data directory.

## `init.keyring-backend`

Specifies a [keyring backend](https://docs.cosmos.network/master/run-node/keyring.html).

### Example

```yaml
init:
  keyring-backend: "os"
```

## `genesis`

Overwrites properties in `config.genesis.json` in the data directory.

### Example

```yaml
genesis:
  chain-id: "foobar"
```

## Learn more

- [Starport](https://github.com/tendermint/starport)
- [Cosmos SDK documentation](https://docs.cosmos.network)
- [Cosmos Tutorials](https://tutorials.cosmos.network)
- [Channel on Discord](https://discord.gg/W8trcGV)