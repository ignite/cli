# What are modules

In the Cosmos SDK modules are the basis for the logic of your blockchain. Each module serves specific purposes and functions. The Cosmos SDK offers a variety of native modules to make a blockchain work. These modules handle authentication for users, token transfers, governance functions, staking of tokens, supply of tokens and many more.

If you want to change the default functionality of a module or just change certain hardcoded parameter, you can fork a module and change it, therefore owning your own logic for your blockchain. While forking and editing a module should be done carefully, this approach marks the Cosmos SDK as especially powerful, as you can experiment with different parameters as the standard implementation suggests.

Modules do not need to be created by a specific company or individual. They can be created by anyone and offered for general use to the public. Although there do exist standards that projects look into before integrating a module to their blockchain. It is recommended that a module has understandable specifications, handles one thing good and is well tested - optimally battle-tested on a live blockchain.
When growing more complex, sometimes it makes more sense to have two modules instead of one module trying to "solve-it-all", this consideration can make it more attractive for other projects to use a module in their blockchain project.

## Standard modules

When creating a blockchain with starport or manually with the Cosmos SDK, there is a set of modules that you should be using in order to have a set of rules that define a blockchain in the first place.
These modules are 

- [What are modules](#what-are-modules)
  - [Standard modules](#standard-modules)
  - [Auth](#auth)
  - [Bank](#bank)
  - [Staking](#staking)
  - [Distribution](#distribution)
  - [Params](#params)
  - [Using modules](#using-modules)
  - [Summary](#summary)

## Auth

The `auth` module is responsible for accounts on the blockchain and basic transaction types that accounts can use. 

It introduces Gas and Fees as concepts in order to prevent the blockchain to bloat by not-identifyable accounts, as on public blockchains you do not have more information about accounts as the public key or balance of an account or the previous transaction history. 

The interface of an Account is defined as

```go
type BaseAccount struct {
  Address       AccAddress
  Coins         Coins
  PubKey        PubKey
  AccountNumber uint64
  Sequence      uint64
}
```

The `auth` module exposes the account keeper where accounts can be stored or changed. Furthermore it exposes the types for Standard transactions, fees, signatures or replay-prevention. It also allows for vesting accounts, as that certain coins can be made accessible over a period of time to an entity. The vesting logic is mostly used for unbonding of staking but can also be used for other purposes.

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/auth/spec/README.md)_

## Bank

The `bank` module has its name because it acts as the module that allows for token transfers and checks for the validity of each transfer. Furthermore, it is responsible for checking the whole supply of the chain, the sum of all account balances.

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/bank/spec/README.md)_

## Staking

The `staking` module allows for an advanced Proof of Stake system, where validators can be created and tokens delegated to validators. 

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/staking/spec/02_state_transitions.md#slashing)_

## Distribution

The `distribution` module is responsible to distribute the inflation of a Token. When new Tokens get created, they get distributed to the validators and their respective delegators, with a potential commission the validator takes. Each validator can choose a commission of those Token when creating a validator, this commission is editable.

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/distribution/spec/README.md)_

## Params

The `params` module allows for a global parameter store in your blockchain application. It is designed to hold the chain parameters so that they can be changed during runtime by governance. It allows to upgrade the blockchain parameters via the `government` module and take effect on an agreed upon time when the majority of the shareholders decide to make a change.


_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/params/spec/README.md)_

Those modules are typically installed on default when using starport. There are a range of modules that are also part of the Cosmos SDK, additionally some other public modules have already reached a major level of usage and acceptance. We will look at more advanced modules in the next Chapter.

## Using modules

With starport you can add a module with the command `starport module create modulename`. When adding a module manually to a blockchain application, it requires to edit the `app/app.go` and the `myappcli/main.go` with the according entries. Starport manages the code edits and additions for you conveniently.

## Summary

- Importing modules in a Cosmos SDK built blockchain exposes new functionalities for the blockchain.
- Any combination of modules is allowed.
- The modules define what can be done on the blockchain.
- Modules are editable, but the success of your blockchain will be dependend on choosing the correct modules for your blockchain, for functionality and security sake.
- `starport module import <modulename>` lets you import modules into your blockchain application.
