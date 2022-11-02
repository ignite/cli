---
sidebar_position: 3
description: Primary configuration file to describe the development environment for your blockchain.
title: config.yml reference
---

# Configuration file reference

The `config.yml` file generated in your blockchain folder uses key-value pairs to describe the development environment for your blockchain.

Only a default set of parameters is provided. If more nuanced configuration is required, you can add these parameters to the `config.yml` file.

## Accounts

A list of user accounts created during genesis of the blockchain.

```yml
accounts:
  - name: alice
    coins: ['20000token', '200000000stake']
  - name: bob
    coins: ['10000token', '100000000stake']
```

Ignite uses information from `accounts` when initializing the chain with `ignite chain init` and `ignite chain start`. In the example above Ignite will add two accounts to the `genesis.json` file of the chain.

`name` is a local name of a key pair associated with an account. Once the chain is initialized and started, you will be able to use `name` when signing transactions. With the configuration above, you'd be able to sign transactions both with Alice's and Bob's accounts like so `exampled tx bank send ... --from alice`.

`coins` is a list of token balances for the account. If a token denomination is in this list, it will exist in the genesis balance and will be a valid token. When initialized with the config file above, a chain will only have two accounts at genesis (Alice and Bob) and two native tokens (with denominations `token` and `stake`).

By default, every time a chain is re-initialized, Ignite will create a new key pair for each account. So even though the account name can remain the same (`bob`), every chain reinitialize it will have a different mnemonic and address.

If you want an account to have a specific address, provide the `address` field with a valid bech32 address. The prefix (by default, `cosmos`) should match the one expected by your chain. When an account is provided with an `address` a key pair will not be generated, because it's impossible to derive a key from an address.

```yml
accounts:
  - name: bob
    address: cosmos1s39200s6v4c96ml2xzuh389yxpd0guk2mzn3mz
    coins: ['20000token', '200000000stake']
```

If you want an account to be initialized from a specific mnemonic, provide the `mnemonic` field with a valid mnemonic. A private key, a public key and an address will be derived from a mnemonic.

```yml
accounts:
  - name: bob
    mnemonic: cargo ramp supreme review change various throw air figure humble soft steel slam pole betray inhale already dentist enough away office apple sample glue
    coins: ['20000token', '200000000stake']
```

You cannot have both `address` and `mnemonic` defined for a single account.

Some accounts are used as validator accounts (see `validators` section). Validator accounts cannot have an `address` field, because Ignite needs to be able to derive a private key (either from a random mnemonic or from a specific one provided in the `mnemonic` field). Validator accounts should have enough tokens of the staking denomination for self-delegation.

By default, the `alice` account is used as a validator account, its key is derived from a mnemonic generated randomly at genesis, the staking denomination is `stake`, and this account has enough `stake` for self-delegation.