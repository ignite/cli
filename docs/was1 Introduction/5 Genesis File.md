# Genesis Block

The first block in a blockchain, block 0, is traditionally called the _genesis_ or _genesis block_.

For all blockchains, the genesis block is the starting point of history. Every block in a blockchain contains a hash of all transactions that it embeds and a pointer to the previous block. Because the genesis block is the starting point of a blockchain, the genesis block does not have a pointer to a previous block. With Starport, the `genesis.json` file for the new blockchain is automatically created from your `config.yml` file. For information about the genesis file and field definitions, see [Using Tendermint > Genesis](https://docs.tendermint.com/master/tendermint-core/using-tendermint.html#genesis) in the Tendermint Core documentation.

## Summary

- The genesis block is the first block of a blockchain.

- The genesis block contains initial stakeholders and starting validators.
