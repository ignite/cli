# An introduction to Starport

Let us dive into Starport, what we can achieve with it, and which other technologies play well with it.

## What is Starport?

Starport is a tool that makes it easier to create blockchains.

Starport uses the Tendermint Consensus engine and the Cosmos SDK to create a blockchain application in the Go programming language. This blockchain has a Proof-of-Stake system with validators (https://en.longhash.com/news/how-cosmos-governance-works-and-how-you-can-become-a-validator) that can be defined in the genesis block.

With just a few command lines, you can create a blockchain, launch it, serve it on the cloud and have a GUI ready to start testing your application.

Bootstrapping blockchains was initially the job of the `scaffold` program, which helped to create a blockchain application. Starport takes it to the next level and also creates a user interface with Vue.js that provides a good starting point for developing a browser-based client-side application for your blockchain.

The scaffolded application still includes a command line interface that lets you manage keys, create validators, send tokens.

## Projects using Tendermint / Cosmos SDK

There are many projects already showcasing that the Tendermint BFT Consensus and the Cosmos SDK enables a variety of usecases and with a little effort can become part of the most robust blockchains that exist.

The technology stack is being used by projects such as:

- [Cosmos](https://github.com/cosmos/gaia) (Main IBC Hub and "Rolemodel" of the Cosmos SDK)
- [Binance Chain](https://github.com/binance-chain) (DEX and utility token)
- [Crypto.com Chain](https://github.com/crypto-com/chain-main) (Payments, DeFi, and utility token)
- [IRIS](https://github.com/irisnet) (IBC Hub and developer oriented)
- [Kava](https://github.com/Kava-Labs/kava) (DeFi and Stable Coins)
- [Aragon](https://docs.chain.aragon.org/) (DAO catalyst)
- [CosmWasm](https://cosmwasm.com/) (smart contracts using WASM)
- [Ethermint](https://ethermint.zone/) (Ethereum virtual machine)

[See the full list here](https://cosmonauts.world/).

The basis of the consensus protocol, network management and blockchain initialization has been proven to be very reliable. From there, you can build anything with this technology stack.

## Modules

What empowers building blockchains and your own application on top of it are the Cosmos modules. Modules are developed by companies and communities around the Cosmos SDK. Some of the modules are crucial to build a basic blockchain while others enable new features for a blockchain.

Each of the running blockchains most of the time have at least one module developed on their own to enable the features it specialises on. These standards solve the clearity of a project and the understanding of all the different parts that create a specific blockchain application.

The most basic set of modules that are created with starport are `auth`, `bank`, `staking`, `params` and `supply`. If you would also add the `wasm` module, this would allow you to upload Webassembly smart contracts on your blockchain application. The `evm` module enables the Ethereum EVM to be used on your blockchain. Each module comes with a clear documentation and codebase. If you wanted changes to a specific module, you can fork the module and change what suits your blockchain application better.

## Summary

- Starport lets you create, develop, and build a blockchain.
- Starport and Cosmos are written in Go
- Today, Cosmos SDK has a unique position worldwide as one of the most successful blockchains.
- Developers can use different Cosmos Modules to customize their blockchain. 
