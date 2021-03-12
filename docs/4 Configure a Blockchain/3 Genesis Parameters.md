# Genesis overwrites for development

The `genesis.json` file for the new blockchain is automatically created from the `config.yml` file.

In development environments, it is useful to test different scenarios after the blockchain is created. The `genesis.json` file for the blockchain is overwritten by the top-level genesis parameter in `config.yml`.

To set and test different values, add the `genesis` parameter to `config.yml`.

## Change value of one parameter

To change the value of one parameter, add the key-value pair under the `genesis` parameter. For example, change the chain-id:

```yml
genesis:
  chain_id: "foobar"
```

## Change values in modules

You can change one or more parameters of different modules. For example, in the `staking` module you can add a key-value pair to `bond_demon` to change which token gets staked:

```yml
genesis:
  app_state:
    staking:
      params:
        bond_denom: "denom"
```

## Genesis file

For genesis file details and field definitions, see [Using Tendermint > Genesis](https://docs.tendermint.com/master/tendermint-core/using-tendermint.html#genesis).

## Summary

- The genesis block is the first block of a blockchain.

- The `genesis.json` file for the blockchain is overwritten by the top-level genesis parameter in `config.yml`.

- After the blockchain is created, add the `genesis` parameter and key-value pairs to set and test different values in your development environment.
