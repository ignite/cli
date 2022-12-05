---
sidebar_position: 4
description: Test different scenarios after the blockchain is created.
---

# Genesis overwrites for development

The `genesis.json` file for all new blockchains is automatically created from the `config.yml` file to define the initial state upon genesis of the blockchain.

In development environments, it is useful to test different scenarios after the blockchain is created. The `genesis.json` file for the blockchain is overwritten by the top-level `genesis` parameter in `config.yml`.

To set and test different values, add the `genesis` parameter to `config.yml`.

## Change the value of a single parameter

To change the value of one parameter, add the key-value pair under the `genesis` parameter. For example, change the value of `chain-id`:

```yaml
genesis:
  chain_id: "foobar"
```

## Change values in modules

You can change one or more parameters of different modules. For example, in the `staking` module you can add a key-value pair to `bond_denom` to change which token gets staked:

```yaml
genesis:
  app_state:
    staking:
      params:
        bond_denom: "denom"
```

## Genesis file

For genesis file details and field definitions, see Cosmos Hub documentation for the [Genesis File](https://hub.cosmos.network/main/resources/genesis.html).

## Genesis block summary

- The genesis block is the first block of a blockchain.

- The `genesis.json` file for the blockchain is overwritten by the top-level genesis parameter in `config.yml`.

- After the blockchain is created, add the `genesis` parameter and key-value pairs to set and test different values in your development environment.
