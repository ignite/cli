# Advanced Modules

As described in the previous chapter, you need a few basic modules of the Cosmos SDK to deploy a basic blockchain. The modules we will be looking into in this chapter will add optional functions to a blockchain application. These modules enable from live upgrading a software update for a blockchain - to inter-blockchain communication, smart contracts or more.

The modules chosen for this chapter are:

- [Governance](#governance)
- [IBC](#ibc)
- [Slashing](#slashing)
- [Upgrade](#upgrade)
- [Wasm](#wasm)

## Governance

The `governance` module allows on-chain voting mechanisms. The documentation lists the following four features:

- **Proposal submission:** Users can submit proposals with a deposit. Once the
minimum deposit is reached, proposal enters voting period
- **Vote:** Participants can vote on proposals that reached MinDeposit
- **Inheritance and penalties:** Delegators inherit their validator's vote if
they don't vote themselves.
- **Claiming deposit:** Users that deposited on proposals can recover their
deposits if the proposal was accepted OR if the proposal never entered voting period.

These governance proposals are used in a variety of ways. It can be a general voting mechanism for the blockchain, but historically the government proposals have been used for

- Changing parameters of the blockchain application like inflation, slashing or software upgrades.
- Community spending of a pool.
- Token holder Polls

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/gov/spec/README.md)_

## IBC

The `IBC` (Inter-blockchain communication protocol) allows to send tokens from one chain to another. This usually happens with locking the token on one chain and unlocking it on another chain. It is not only for the general blockchain unit but also supports tokens like UTXO (Bitcoin), NFT Standards (ERC721) or other tokens (like ERC20/ERC777).

The `IBC` module allows to create zones on blockchains that can then interact with each other. There might be more data to share between zones than tokens. These zones can furthermore act to secure tokens, create exchanges or certain interactions of tokens or transactions. This behaviour allows to be used as a Sharding Hub, where certain interactions can be outsurced. 

_Read the [specification](https://github.com/cosmos/ics)_

## Slashing

`Slashing` is an optional module but used on many live blockchains. The module allows to disincentivize non-conform blockchain actions from validators. Validators play a crucial role in keeping the blockchain online, alive, active and secure. In case a validator misbehaves and e.g. create blocks that lead to a fork of the blockchain, the validator can be automatically penalised by a set amount of tokens that the validator is responsible for. On many blockchains, these penalty is currently set at 5% share burn of the validators self-delegated and delegated tokens. Additionally the validator may be added to the `tumbstone`, in result not being able to create another block ever again. Another value the `slashing` module can be used for is penalizing a validator for not signing enough blocks. On many chains, when the validator misses 95% of the last 10,000 blocks it will get slashed for 0.1% of the tokens responsible.

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/slashing/spec/README.md)_

## Upgrade

`Upgrade` is an optional module that enables on-chain upgrades of your blockchain application. It allows for breaking changes to be included into the blockchain by an agreed-upon blockchain height or time. 
Before having the upgrade module available, blockchains would need to halt, create a new genesis block and restart the chain all validators together in order to reach a new consensus and get the blockchain started again. This module allows for seemless upgrades for a blockchain.

_Read the [specification](https://github.com/cosmos/cosmos-sdk/blob/master/x/upgrade/spec/README.md)_

## Wasm

The `wasm` module enables smart contracts written in rust on your blockchain applications. Rust is a programming language that focuses on performance and safety, especially safe concurrency. Smart contracts can get uploaded on the blockchain by anyone and executed by anyone. The module will support validation of contracts and the output of their respective contract functions. It allows to use contract standards from the Ethereum Virtual Machine or deploying new contracts that you develop on your own.

_Read the [specification](https://github.com/CosmWasm/cosmwasm)_

## Your own module

This is only a starting collection for available modules on the Cosmos SDK. Many blockchain projects have created their own module that has a specialised behaviour. If you want to develop a module on your own, there is a definition that describes the [standards for modules](https://github.com/cosmos/cosmos-sdk/blob/master/docs/building-modules/README.md) available.

## Summary

- The presented advanced modules are optional but can be added plug-and-play to your blockchain application.
- The modules allow extended behaviours such as polling, releasing shares to certain entities, changing the blockchain, smart contracts or slashing validators for misbehavior.
- Modules are developed by the teams behind Cosmos but also decentralized, modules can be created by any entity.

[◀️ Previous - Basic Modules](../../03%20Modules/02_basic_modules/02_basic_modules.md) | [▶️ Next - Smart Modules](../../03%20Modules/04_smart_modules/04_smart_modules.md)  