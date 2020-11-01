# Introduction

Starport creates a blockchain for you in Golang. Requirements for this is to have Golang installed. You can get all the information here https://golang.org/doc/install.

Starport installation instructions can be found here: https://github.com/tendermint/starport#install

To the machine you are executing on there are not many requirements. It runs on Linux or Mac Operating Systems and can be run from a Raspberry Pi.

## Your blockchain application

To create a blockchain application we use the command `app`

```bash
starport app github.com/username/myapp
```

| Flag               | Default  | Description                |
| ------------------ | -------- | -------------------------- |
| `--address-prefix` | `cosmos` | Prefix, used for addresses |

This will create the folder `myapp` and is a usable blockchain blueprint. If you want to dive directly into looking at the details of your blockchain you can run it with entering your `myapp` folder and use the command `serve` to initialise your blockchain and start it.

## Serve

`starport serve`

To start the server, go into you application's directory and run `starport serve`. This commands installs dependencies, builds and initializes the app and runs both Tendermint RPC server (by default on `localhost:26657`) as well as LCD (by default on `localhost:1317`) with hot reloading enabled.

`starport serve` uses `config.yml` to initialize your application, make sure you have it in your project directory (see [Configure](#configure)).

Note: depending on your OS and firewall settings, you may have to accept a prompt asking if your application's binary (`blogd` in this case) can accept external connections.

| Flag        | Default | Description                          |
| ----------- | ------- | ------------------------------------ |
| `--verbose` | `false` | Enable verbose output from processes |
| `--path`    |         | Path to the project                  |

The first step of your own blockchain is already done. Using the default settings, a blockchain that has networking, consensus protocol with an own token is hereby established. From here on, you can implement logic that makes your own blockchain unique. 

## The Key-Value Store (KV)
### How to use types

In the SDK, data is stored in the multistore. Key-Value pairs are saved in the KVStores. Multiple stores can be created and managed at the same time. We will use the store to save our data to the blockchain.
Starport assists us in setting up the Key-Value Store with the command `type`. 
In order to use `type` we should give our type a fitting `typeName` with the intended fields that we want to use. If we wanted to store user with username and age, we would use the command

```
starport type [typeName] [field1] [field2:bool] ...
```

More specific

`starport type user username age:int` 

This command generates messages, handlers, keepers, CLI and REST clients and type definition for `typeName` type. A type can have any number of `field` arguments. By default fields are strings, but `bool` and `int` are supported.

Now a Key-Value Store for the user with fields username and age is created. We can create a new user with the command

`myappcli tx myapp create-user "my-first-username" 35`

Which creates the user with username `my-first-username` and age of `35`. 

Another example,

```
starport type post title body
```

This command generates a type `Post` with two fields: `title` and `body`.

To add a post run `blogcli tx blog create-post "My title" "This is a blog" --from=user1`.

These are the basic commands for getting started with starport. From creating a first blockchain to adding your own data types and accessing the User Interface. In the next two chapters, we will be looking closer at the initial setup for starport and how to configure it. Afterwards, we will be looking into more complex usecases, where each of the commands and more will be explained in detail.

#### Accounts on your blockchain

An account on the blockchain is a keypair of private and public keys.
When you start your blockchain with starport, you can define the name of the keys and the amount of coins they start with. The keys are created for you and displayed on startup. You can use these keys when interacting with your blockchain.
A list of user accounts is created during genesis of your application. You can define them as follows in your `config.yml` file. See an example in chapter [configuration](../03_configuration/03_configuration.md).

| Key   | Required | Type            | Description                                       |
| ----- | -------- | --------------- | ------------------------------------------------- |
| name  | Y        | String          | Local name of the key pair                        |
| coins | Y        | List of Strings | Initial coins with denominations (e.g. "100coin") |

#### The initial validator

Blocks on a tendermint blockchain are created and validated by the so called `validators`. You can define the set of validators your blockchain starts with in your `config.yml`.
The validator property describes your set of validators. Use a `name` that you have specified in the `accounts` array. The account should have enough tokens for staking purposes. See an example in chapter [configuration](../03_configuration/03_configuration.md).

| Key    | Required | Type   | Description                                                                         |
| ------ | -------- | ------ | ----------------------------------------------------------------------------------- |
| name   | Y        | String | Name of one the accounts                                                            |
| staked | Y        | String | Amount of coins staked by your validator, should be >= 10^6 (e.g. "100000000stake") |

## Summary

- With the command `starport app` a new blockchain can be initialised.
- A combination `starport app` and `starport serve` already let's you manage your blockchain out of the box.
- The default blockchain includes networking and a consensus protocol with your own token.
- Data is managed with the Key-Value Store and data types can be added with `starport type`.
- Accounts are created during genesis of the application. These can be configured in the `config.yml`.
