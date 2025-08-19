---
sidebar_position: 1
description: Primary configuration file to describe the development environment for your blockchain.
title: Configuration File Documentation
---

# Configuration File Reference

After scaffolding a blockchain with Ignite CLI, you will find a configuration file at the root of your newly created directory.

The `config.yml` file generated in your blockchain folder uses key-value pairs
to describe the development environment for your blockchain.

Only a default set of parameters is provided. If more nuanced configuration is
required, you can add these parameters to the `config.yml` file.

## Genesis

The genesis file is the initial block of your blockchain. It is required to launch a chain because it contains important
information such as token balances and modules' state.
By default, genesis is stored at `$DATA_DIR/config/genesis.json`.

Since the genesis file is frequently reinitialized during development, you can persistently set options using the
`genesis` property in your `config.yml`:

```yml
genesis:
  app_state:
    staking:
      params:
        bond_denom: "denom"
```

To discover which properties the genesis file supports, initialize a chain and inspect the generated genesis file in the
data directory.

### Overriding Genesis Parameters (e.g., chain_id, balances, etc.)

You may need to customize specific parameters in the genesis file, such as `chain_id`, token balances, module
parameters, or custom state.

To override genesis values with Ignite CLI, persistently set overrides in the `genesis`property of your `config.yml`.
Any YAML structure under `genesis` will be merged into the generated `genesis.json` during initialization.

Eg: Changing `chain_id` and `staking` parameters

```yml
genesis:
  chain_id: "my-custom-chain"
  app_state:
    staking:
      params:
        bond_denom: "mytoken"
    bank:
      balances:
        - address: "cosmos1..."
          coins:
            - denom: "mytoken"
              amount: "1000000"
```

- `chain_id`: Sets the chain ID for your blockchain.
- `app_state`: Allows you to modify module states (e.g., staking, bank, etc.).

> ⚠️ If you set `chain_id` in the `genesis`, it will persist across `ignite chain init` or `ignite chain serve` runs.

The `genesis` property supports deep merging and can override any field present in the generated genesis file.
For more complex setups, you can use the `include` field in `config.yml` to split overrides into multiple files.

## Validation

Ignite uses the `validation` field to determine the kind of validation
of your blockchain. There are currently two supported kinds of validation:

- `sovereign` which is the standard kind of validation where your blockchain
  has its own validator set. This is the default value when this field is not
  in the config file.
- `consumer` indicates your blockchain is a consumer chain, in the sense of
  Replicated Security. That means it doesn't have a validator set, but
  inherits the one of a provider chain.

While the `sovereign` chain is the default validation when you run the `ignite scaffold
chain`, to scaffold a consumer chain, you have to run `ignite scaffold chain
--consumer`.

This field is, at this time of writing, only used by Ignite at the genesis
generation step, because the genesis of a sovereign chain and a consumer chain
are different.

## Accounts

A list of user accounts created during genesis of the blockchain.

```yml
accounts:
  - name: alice
    coins: [ '20000token', '200000000stake' ]
  - name: bob
    coins: [ '10000token', '100000000stake' ]
```

Ignite uses information from `accounts` when initializing the chain with `ignite
chain init` and `ignite chain start`. In the example above Ignite will add two
accounts to the `genesis.json` file of the chain.

`name` is a local name of a key pair associated with an account. Once the chain
is initialized and started, you will be able to use `name` when signing
transactions. With the configuration above, you'd be able to sign transactions
both with Alice's and Bob's accounts like so `exampled tx bank send ... --from
alice`.

`coins` is a list of token balances for the account. If a token denomination is
in this list, it will exist in the genesis balance and will be a valid token.
When initialized with the config file above, a chain will only have two accounts
at genesis (Alice and Bob) and two native tokens (with denominations `token` and
`stake`).

By default, every time a chain is re-initialized, Ignite will create a new key
pair for each account. So even though the account name can remain the same
(`bob`), every chain reinitialize it will have a different mnemonic and address.

If you want an account to have a specific address, provide the `address` field
with a valid bech32 address. The prefix (by default, `cosmos`) should match the
one expected by your chain. When an account is provided with an `address` a key
pair will not be generated, because it's impossible to derive a key from an
address. An account with a given address will be added to the genesis file (with
an associated token balance), but because there is no key pair, you will not be
able to broadcast transactions from that address. This is useful when you have
generated a key pair outside of Ignite (for example, using your chain's CLI or
in an extension wallet) and want to have a token balance associated with the
address of this key pair.

```yml
accounts:
  - name: bob
    coins: [ '20000token', '200000000stake' ]
    address: cosmos1s39200s6v4c96ml2xzuh389yxpd0guk2mzn3mz
```

If you want an account to be initialized from a specific mnemonic, provide the
`mnemonic` field with a valid mnemonic. A private key, a public key and an
address will be derived from a mnemonic.

```yml
accounts:
  - name: bob
    coins: [ '20000token', '200000000stake' ]
    mnemonic: cargo ramp supreme review change various throw air figure humble soft steel slam pole betray inhale already dentist enough away office apple sample glue
```

You cannot have both `address` and `mnemonic` defined for a single account.

Some accounts are used as validator accounts (see `validators` section).
Validator accounts cannot have an `address` field, because Ignite needs to be
able to derive a private key (either from a random mnemonic or from a specific
one provided in the `mnemonic` field). Validator accounts should have enough
tokens of the staking denomination for self-delegation.

By default, the `alice` account is used as a validator account, its key is
derived from a mnemonic generated randomly at genesis, the staking denomination
is `stake`, and this account has enough `stake` for self-delegation.

If your chain is using its own
[cointype](https://github.com/satoshilabs/slips/blob/master/slip-0044.md), you
can use the `cointype` field to provide the integer value

```yml
accounts:
  - name: bob
    coins: [ '20000token', '200000000stake' ]
    cointype: 7777777
```

## Validators

Commands like `ignite chain init` and `ignite chain serve` initialize and launch
a validator node for development purposes.

```yml
validators:
  - name: alice
    bonded: '100000000stake'
```

`name` refers to key name in the `accounts` list.

`bonded` is the self-delegation amount of a validator. The `bonded` amount
should not be lower than `1000000` nor higher than the account's
balance in the `account` list.

Validators store their node configuration files in the data directory. By
default, Ignite uses the name of the project as the name of the data directory,
for example, `$HOME/.example/`. To use a different path for the data directory
you can customize the `home` property.

Configuration in the data directory is reset frequently by Ignite. To persist
some changes to configuration files you can use `app`, `config` and `client`
properties that correspond to `$HOME/.example/config/app.toml`,
`$HOME/.example/config/config.toml` and `$HOME/.example/config/client.toml`.

```yml
validators:
  - name: alice
    bonded: '100000000stake'
    home: "~/.mychain"
    app:
      pruning: "nothing"
    config:
      moniker: "mychain"
    client:
      output: "json"
```

To see which properties are available for `config.toml`, `app.toml` and
`client.toml`, initialize a chain with `ignite chain init` and open the file you
want to know more about.

Currently, Ignite starts only one validator node, so the first item in the
`validators` list is used (the rest is ignored). Support for multiple validators
is in progress.

## Build

The `build` property lets you customize how Ignite builds your chain's binary.

By default, Ignite builds the `main` package from `cmd/PROJECT_NAME/main.go`. If
you more than one `main` package in your project, or you have renamed the
directory, use the `main` property to provide the path to the `main` Go package:

```yml
build:
  main: cmd/hello/cmd
```

Ignite compiles your project into a binary and uses the project's name with a
`d` suffix as name for the binary. To customize the binary name use the `binary`
property:

```yml
build:
  binary: "helloworldd"
```

To customize the linker flags used in the build process:

```yml
build:
  ldflags: [ "-X main.Version=development", "-X main.Date=01/05/2022T19:54" ]
```

By default, custom protocol buffer (proto) files are located in the `proto`
directory. If your project keeps proto files in a different directory, you
should tell Ignite about this:

```yml
build:
  proto:
    path: "myproto"
```

## Faucet

The faucet service sends tokens to addresses.

```yml
faucet:
  name: bob
  coins: [ "5token", "100000stake" ]
```

`name` refers to a key name in the `accounts` list. This is a required property.

`coins` is the amount of tokens that will be sent to a user by the faucet. This
is a required property.

`coins_max` is a maximum amount of tokens that can be sent to a single address.
To reset the token limit use the `rate_limit_window` property (in seconds).

The default the faucet works on port `4500`. To use a different port number use
the `port` property.

```yml
faucet:
  name: faucet
  coins: [ "100token", "5foo" ]
  coins_max: [ "2000token", "1000foo" ]
  port: 4500
  rate_limit_window: 3600
```

## Genesis

Genesis file is the initial block in the blockchain. It is required to launch a
blockchain, because it contains important information like token balances, and
modules' state. Genesis is stored in `$DATA_DIR/config/genesis.json`.

Since the genesis file is reinitialized frequently during development, you can
set persistent options in the `genesis` property:

```yml
genesis:
  app_state:
    staking:
      params:
        bond_denom: "denom"
```

To know which properties a genesis file supports, initialize a chain and look up
the genesis file in the data directory.

## Client code generation

Ignite can generate client-side code for interacting with your chain with the
`ignite generate` set of commands. Use the following properties to customize the
paths where the client-side code is generated.

```yml
client:
  openapi:
    path: "docs/static/openapi.json"
  typescript:
    path: "ts-client"
  composables:
    path: "vue/src/composables"
  hooks:
    path: "react/src/hooks"
```

## Include

In your main `config.yml`, use the `include` field to reference other local or remote YAML files.
It allows you to split your chain configuration across multiple files, making it easier to manage and reuse configuration parts.  

```yml
version: 1
include:
  - "./accounts.yml"
  - "./validators.yml"
```

Include remote files via URL or server path are also valid:

```yml
version: 1
include:
  - "localhost:3045/accounts.yml"
  - "https://ignite.com/config/validators.yml"
```

#### Common Use Cases:

Split your config into a base setup and an external `accounts.yml` for better separation of concerns:

- `config.yml`
```yml
version: 1
include:
  - "./accounts.yml"
client:
  typescript:
    path: ts-client
```

- `accounts.yml`
```yml
accounts:
  - name: alice
    coins:
      - 20000token
      - 200000000stake
  - name: bob
    coins:
      - 20000token
      - 200000000stake
faucet:
  name: alice
  coins:
    - 5token
    - 100000stake
```
