# Configuration

When creating a new app with starport, you will see a `config.yml` file in your blockchain application folder. This file defines the genesis file, the network you will be using and the first validators of your blockchain.

Let us examine the parts of the configuration of this file.

```yml
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

The first line - `version` is used to indicate changes, which is especially important when working on public code or in a team.

The `accounts` define the initial distribution of Tokens on the blockchain. Here is the place where you can define original holders of the tokens on your blockchain. These accounts will get translated into the genesis block and after launching your blockchain the users mentioned have access to their respective tokens.
The `name` parameter in `accounts` will create a random new keypair in your blockchain app keychain, which you can access on the command line. You can also reference these names under `validator` to define the starting validators with a bounded stake that you can configure. The bounding stake has to be equal to or less the stake given in the `accounts` paramenter.

There is an optional genesis parameter in the `config.yml`, which you can use to define parameters in your genesis file directly, such as the chain-id as follows:

```yml
genesis:
  chain_id: "foobar"
```

The `coins` specify amount of coins and their denomination on the blockchain. Here you can list a variety of coin denominations and their respective amounts to be used on your blockchain.

After manipulating the `config.yml` to your likings, you can start the blockchain application with `starport serve`. This will create (or override) the folder located at your user homefolder `~/.myappd` (the name of your application with a `d` for `daemon` attached) and initiate your blockchain with the genesis file, located under `~/.myappd/config`. The second folder you can find in the `~/.myappd` folder is `data` - this is where the blockchain will write the consecutive blocks and transactions.

If you want to get sure all data from a blockchain setup get deleted, make sure to remove the `~/.myappd` folder.

## Address denomination

You can change the way addresses look in your blockchain application. Namely what they have attached in the beginning. On the Cosmos SDK Main Hub addresses are displayed with a `cosmos` in front of their address, e.g.

`cosmos12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf`

You can change the first prefix by changing the `AccountAddressPrefix` variable in your created application `/app/prefix.go`. It is recommended to leave the other variables as they are, because those are standards used in the other Cosmos SDK chains and can therefore be recognised. These have security implications such as not sending to addresses that might not be able to spend it.

## Summary

- The `config.yml` defines your genesis accounts and validators.
- It lets you bootrstrap your blockchain with different tokens and specify the amount of each in the first block.
- Changing the prefix for addresses can be done in the `/app/prefix.go` file.

[◀️ Previous - Genesis File](../../02_using_starport/02_genesis_file/02_genesis_file.md) [▶️ Next - What are Modules](../../03_modules/01_what_are_modules/01_what_are_modules.md)  