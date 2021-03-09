# Configure Starport

The `config.yml` file generated in your blockchain folder describes the development environment for your blockchain.

<!-- TOC depthFrom:2 depthTo:2 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Genesis file](#genesis-file)
- [config.yml Attributes and Examples](#configyml-attributes-and-examples)
- [accounts](#accounts)
- [`build`](#build)
- [Example build](#example-build)
- [`build.proto`](#buildproto)
- [`faucet`](#faucet)
- [`validator`](#validator)
- [`init.home`](#inithome)
- [`init.config`](#initconfig)
- [`init.app`](#initapp)
- [`init.keyring-backend`](#initkeyring-backend)
- [`genesis`](#genesis)
- [Learn more](#learn-more)

<!-- /TOC -->

 The generated `config.yml` look like:

```yml
accounts:
  - name: user1
    coins: ["1000token", "100000000stake"]
  - name: user2
    coins: ["500token"]
validator:
  name: user1
  staked: "100000000stake"
```

## Genesis file

With Starport, a genesis file for the new blockchain is automatically created from your `config.yml` file. See [Using Tendermint > Genesis](https://docs.tendermint.com/master/tendermint-core/using-tendermint.html#genesis) in the Tendermint Core documentation.

There is an optional genesis parameter in the `config.yml`, which you can use to define parameters in your genesis file directly, such as the chain-id as follows:

```yml
genesis:
  chain_id: "foobar"
```

You can also manipulate parameters of different modules. If you wanted for example change the `staking` module, which contains staking parameters such as which token gets staked, you would add the following to your `config.yml``

```yml
genesis:
  app_state:
    staking:
      params:
        bond_denom: "denom"
```

After manipulating the `config.yml` to your likings, you can start the blockchain app with `starport serve`. This will create (or override) the folder located at your user homefolder `~/.myappd` (the name of your app with a `d` for `daemon` attached) and initiate your blockchain with the genesis file, located under `~/.myappd/config`. The second folder you can find in the `~/.myappd` folder is `data` - this is where the blockchain will write the consecutive blocks and transactions. The other folder created is the `~/.myappcli` folder, which contains a configuration file for your current command line interface, such as `chain-id`, output parameters such as `json` or `indent` mode. If you want to make sure all of your data from the blockchain setup is deleted, make sure to remove the `~/.myappd` and `~/.myappcli` folder.


 Let us examine the parts of the configuration of this file.

```yml
accounts:
  - name: user1
    coins: ["1000token", "100000000stake"]
  - name: user2
    coins: ["500token"]
validator:
  name: user1
  staked: "100000000stake"
```

The `accounts` define the initial distribution of Tokens on the blockchain. Here is the place where you can define original holders of the tokens on your blockchain. These accounts will get translated into the genesis block and after launching your blockchain the users mentioned have access to their respective tokens. The `name` parameter in `accounts` will create a random new keypair in your blockchain app keychain, which you can access on the command line. You can also reference these names under `validator` to define the starting validators with a bounded stake that you can configure. The bounding stake has to be equal to or less the stake given in the `accounts` parameter. The `coins` specify the amount of coins and their denomination on the blockchain. Here you can list a variety of coin denominations and their respective amounts to be used on your blockchain.

There is an optional genesis parameter in the `config.yml`, which you can use to define parameters in your genesis file directly, such as the chain-id as follows:

```yml
genesis:
  chain_id: "foobar"
```

(useful for development, for example, building a module that has params and you can set these params and can test parameters, app easy to build and test it for certain parameters ) You can also manipulate parameters of different modules. If you wanted for example change the `staking` module, which contains staking parameters such as which token gets staked, you would add the following to your `config.yml``

```yml
genesis:
  app_state:
    staking:
      params:
        bond_denom: "denom"
```

After manipulating the `config.yml` to your likings, you can start the blockchain app with `starport serve`. This will create (or override) the folder located at your user homefolder `~/.myappd` (the name of your app with a `d` for `daemon` attached) and initiate your blockchain with the genesis file, located under `~/.myappd/config`. The second folder you can find in the `~/.myappd` folder is `data` - this is where the blockchain will write the consecutive blocks and transactions. The other folder created is the `~/.myappcli` folder, which contains a configuration file for your current command line interface, such as `chain-id`, output parameters such as `json` or `indent` mode.

If you want to make sure all of your data from the blockchain setup is deleted, make sure to remove the `~/.myappd` and `~/.myappcli` folder.

## config.yml Attributes and Examples

Configure The `config.yml` file

## accounts

A list of user accounts created during genesis of the blockchain that define the initial distribution of tokens on the blockchain. Only the named users .

Key     | Required | Type            | Description
------- | -------- | --------------- | --------------------------------------------------------------------------------------------------------------------
name    | Y        | String          | Local name of a key pair. Only specified account names have access to their tokens after the blockchain is launched.
coins   | Y        | List of Strings | Initial coins with denominations. For example, "1000token"
address | N        | String          | Account address in Bech32 address format

### Example accounts

```yaml
accounts:
  - name: alice
    coins: ["1000token", "100000000stake"]
  - name: bob
    coins: ["500token"]
    address: cosmos1adn9gxjmrc3hrsdx5zpc9sj2ra7kgqkmphf8yw
```

## `build`

Key    | Required | Type   | Description
------ | -------- | ------ | -------------------------------------------------------------------------------
binary | N        | String | Name of the node binary that Starport builds and uses, typically ends with `d-`

## Example build

```yaml
build:
  binary: "mychaind"
```

## `build.proto`

Key               | Required | Type            | Description
----------------- | -------- | --------------- | -------------------------------------------------------------------------------------------
path              | N        | String          | Path to protocol buffer files (default: `"proto"`)
third_party_paths | N        | List of Strings | Path to thid-party protocol buffer files (default: `["third_party/proto", "proto_vendor"]`)

## `faucet`

The faucet service sends tokens to addresses. The default address for the web user interface is <http://localhost:4500>.

Key       | Required | Type            | Description
--------- | -------- | --------------- | ------------------------------------------------
name      | Y        | String          | Name of a key pair. `name` must be in `accounts`
coins     | Y        | List of Strings | Coins with denominations sent per request
coins_max | N        | List of Strings | Maximum amount of tokens sent for each address
port      | N        | Integer         | Port number. Default: `4500`

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

Key    | Required | Type   | Description
------ | -------- | ------ | ---------------------------------------------------------------------------------
name   | Y        | String | Name of a key pair. `name` must be in `accounts`
staked | Y        | String | Amount of coins to bond. Must be less or equal to the amount of coins account has

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
