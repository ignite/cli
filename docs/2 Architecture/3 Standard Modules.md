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
  - [Capability](#capability)
  - [Staking](#staking)
  - [Mint](#mint)
  - [Distribution](#distribution)
  - [Params](#params)
  - [Governance](#governance)
  - [Crisis](#crisis)
  - [Slashing](#slashing)
  - [IBC](#ibc)
  - [Upgrade](#upgrade)
  - [Evidence](#evidence)
  - [Using modules](#using-modules)
  - [Summary](#summary)

## Auth

The `auth` module is responsible for accounts on the blockchain and basic transaction types that accounts can use. 

It introduces Gas and Fees as concepts in order to prevent the blockchain to bloat by not-identifyable accounts, as on public blockchains you do not have more information about accounts as the public key or balance of an account or the previous transaction history. 

The interface of an Account is defined as

```proto
// BaseAccount defines a base account type. It contains all the necessary fields
// for basic account functionality. Any custom account type should extend this
// type for additional functionality (e.g. vesting).
message BaseAccount {
  option (gogoproto.goproto_getters)  = false;
  option (gogoproto.goproto_stringer) = false;
  option (gogoproto.equal)            = false;

  option (cosmos_proto.implements_interface) = "AccountI";

  string              address = 1;
  google.protobuf.Any pub_key = 2
      [(gogoproto.jsontag) = "public_key,omitempty", (gogoproto.moretags) = "yaml:\"public_key\""];
  uint64 account_number = 3 [(gogoproto.moretags) = "yaml:\"account_number\""];
  uint64 sequence       = 4;
}
```

The `auth` module exposes the account keeper where accounts can be stored or changed. Furthermore it exposes the types for Standard transactions, fees, signatures or replay-prevention. It also allows for vesting accounts, as that certain coins can be made accessible over a period of time to an entity. The vesting logic is mostly used for unbonding of staking but can also be used for other purposes.

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/auth/spec/README.md)_


## Bank

The `bank` module has its name because it acts as the module that allows for token transfers and checks for the validity of each transfer. Furthermore, it is responsible for checking the whole supply of the chain, the sum of all account balances.

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/bank/spec/README.md)_


## Capability

Full implementation of the IBC specification requires the ability to create and authenticate object-capability keys at runtime (i.e., during transaction execution), as described in ICS 5. In the IBC specification, capability keys are created for each newly initialised port & channel, and are used to authenticate future usage of the port or channel. Since channels and potentially ports can be initialised during transaction execution, the state machine must be able to create object-capability keys at this time.

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/capability/spec/README.md)_


## Staking

The `staking` module allows for an advanced Proof of Stake system, where validators can be created and tokens delegated to validators. 

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/staking/spec/02_state_transitions.md#slashing)_


## Mint

The minting mechanism is designed to allow for a flexible inflation rate determined by market demand targeting a particular bonded-stake ratio effect a balance between market liquidity and staked supply. 

It can be broken down in the following way:

If the inflation rate is below the goal %-bonded the inflation rate will increase until a maximum value is reached
If the goal % bonded (67% in Cosmos-Hub) is maintained, then the inflation rate will stay constant
If the inflation rate is above the goal %-bonded the inflation rate will decrease until a minimum value is reached

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/mint/spec/README.md)_


## Distribution

The `distribution` module is responsible to distribute the inflation of a Token. When new Tokens get created, they get distributed to the validators and their respective delegators, with a potential commission the validator takes. Each validator can choose a commission of those Token when creating a validator, this commission is editable.

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/distribution/spec/README.md)_


## Params

The `params` module allows for a global parameter store in your blockchain application. It is designed to hold the chain parameters so that they can be changed during runtime by governance. It allows to upgrade the blockchain parameters via the `government` module and take effect on an agreed upon time when the majority of the shareholders decide to make a change.


_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/params/spec/README.md)_

Those modules are typically installed on default when using starport. There are a range of modules that are also part of the Cosmos SDK, additionally some other public modules have already reached a major level of usage and acceptance. We will look at more advanced modules in the next Chapter.


## Governance

The module enables Cosmos-SDK based blockchains to support an on-chain governance system. In this system, holders of the native staking token of the chain can vote on proposals on a 1 token 1 vote basis (from there it can be parameterized). Next is a list of features the module currently supports:

- Proposal submission: Users can submit proposals with a deposit. Once the minimum deposit is reached, proposal enters voting period
- Vote: Participants can vote on proposals that reached MinDeposit
- Inheritance and penalties: Delegators inherit their validator's vote if they don't vote themselves.
- Claiming deposit: Users that deposited on proposals can recover their deposits if the proposal was accepted OR if the proposal never entered voting period.

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/gov/spec/README.md)_


## Crisis

The crisis module halts the blockchain under the circumstance that a blockchain invariant is broken. Invariants can be registered with the application during the application initialization process.

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/crisis/spec/README.md)_


## Slashing

The slashing module enables Cosmos SDK-based blockchains to disincentivize any attributable action by a protocol-recognized actor with value at stake by penalizing them ("slashing").

Penalties may include, but are not limited to:

Burning some amount of their stake
Removing their ability to vote on future blocks for a period of time.
This module will be used by the Cosmos Hub, the first hub in the Cosmos ecosystem.

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/slashing/spec/README.md)_


## IBC

IBC allows to relay packets between chains and could be used with any compatible modules between two chains.
The `IBC` module, as Inter-blockchain Communication, enables for example sending native tokens between blockchains. 
It is divided by a subset of specifications.

_Read the [specifications](https://github.com/cosmos/cosmos-sdk/blob/master/x/ibc/spec/README.md)_


## Upgrade

`x/upgrade` is an implementation of a Cosmos SDK module that facilitates smoothly upgrading a live Cosmos chain to a new (breaking) software version. It accomplishes this by providing a BeginBlocker hook that prevents the blockchain state machine from proceeding once a pre-defined upgrade block time or height has been reached.

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/upgrade/spec/README.md)_


## Evidence

`x/evidence` is an implementation of a Cosmos SDK module, per ADR 009, that allows for the submission and handling of arbitrary evidence of misbehavior such as equivocation and counterfactual signing.

The evidence module differs from standard evidence handling which typically expects the underlying consensus engine, e.g. Tendermint, to automatically submit evidence when it is discovered by allowing clients and foreign chains to submit more complex evidence directly.

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/evidence/spec/README.md)_


## Using modules

With starport you can add a module with the command `starport module create modulename`. When adding a module manually to a blockchain application, it requires to edit the `app/app.go` and the `myappcli/main.go` with the according entries. Starport manages the code edits and additions for you conveniently.

## Summary

- Importing modules in a Cosmos SDK built blockchain exposes new functionalities for the blockchain.
- Any combination of modules is allowed.
- The modules define what can be done on the blockchain.
- Modules are editable, but the success of your blockchain will be dependend on choosing the correct modules for your blockchain, for functionality and security sake.
- `starport module import <modulename>` lets you import modules into your blockchain application.
