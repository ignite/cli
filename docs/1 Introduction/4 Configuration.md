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
The `name` parameter in `accounts` will create a random new keypair in your blockchain app keychain, which you can access on the command line. You can also reference these names under `validator` to define the starting validators with a bounded stake that you can configure. The bounding stake has to be equal to or less the stake given in the `accounts` parameter.
The `coins` specify the amount of coins and their denomination on the blockchain. Here you can list a variety of coin denominations and their respective amounts to be used on your blockchain.

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

After manipulating the `config.yml` to your likings, you can start the blockchain application with `starport serve`. This will create (or override) the folder located at your user homefolder `~/.myappd` (the name of your application with a `d` for `daemon` attached) and initiate your blockchain with the genesis file, located under `~/.myappd/config`. The second folder you can find in the `~/.myappd` folder is `data` - this is where the blockchain will write the consecutive blocks and transactions.
The other folder created is the `~/.myappcli` folder, which contains a configuration file for your current command line interface, such as `chain-id`, output parameters such as `json` or `indent` mode.

If you want to make sure all of your data from the blockchain setup is deleted, make sure to remove the `~/.myappd` and `~/.myappcli` folder.

## Address denomination

You can change the way addresses look in your blockchain application. Namely what they have attached in the beginning. On the Cosmos SDK Main Hub addresses are displayed with a `cosmos` in front of their address, e.g.

`cosmos12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf`

You can change the first prefix by changing the `AccountAddressPrefix` variable in your created application `/app/prefix.go`. It is recommended that you leave the other variables as they are, because those are standards used in the other Cosmos SDK chains and can therefore be recognized. These have security implications such as not sending to addresses that might not be able to spend it.

To have your frontend working properly with the new denomination, you need to change the `VUE_APP_ADDRESS_PREFIX` variable in `/vue/.env`.

To have all of it done automatically, when creating your app with the command `starport app github.com/foo/bar`, just append the `--address-prefix prefix` parameter.

## Summary

- The `config.yml` defines your genesis accounts and validators.
- It lets you bootstrap your blockchain with different tokens and specify the amount of each account in the first block.
- Changing the prefix for addresses can be done in the `/app/prefix.go` file.
