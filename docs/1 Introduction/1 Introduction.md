# An introduction to Starport

Let us dive into Starport, what we can achieve with it, and which other technologies play well with it.

## What is Starport?

Starport is a tool that makes it easier to create blockchains.

Starport uses the Tendermint Consensus engine and the Cosmos SDK to create a blockchain application in the Go programming language. This blockchain has a Proof-of-Stake system with validators (https://en.longhash.com/news/how-cosmos-governance-works-and-how-you-can-become-a-validator) that can be defined in the genesis block.

With just a few command lines, you can create a blockchain, launch it, serve it on the cloud and have a GUI ready to start testing your application.

Bootstrapping blockchains was initially the job of the `scaffold` program, which was used to create a blockchain application. Starport takes it to the next level and also creates a user interface with Vue.js, which provides a good starting point for developers creating a browser-based client-side application for your blockchain.

The scaffolded application still includes a command line interface that lets you manage keys, create validators, send tokens.

## Projects using Tendermint / Cosmos SDK

There are many projects already showcasing that the Tendermint BFT Consensus Engine and the Cosmos SDK.

The following projects are using the technology:

- [Cosmos](https://github.com/cosmos/gaia) (Main IBC Hub and "Rolemodel" of the Cosmos SDK)
- [Binance Chain](https://github.com/binance-chain) (DEX and utility token)
- [Crypto.com Chain](https://github.com/crypto-com/chain-main) (Payments, DeFi, and utility token)
- [IRIS](https://github.com/irisnet) (IBC Hub and developer oriented)
- [Kava](https://github.com/Kava-Labs/kava) (DeFi and Stable Coins)
- [Aragon](https://docs.chain.aragon.org/) (DAO catalyst)
- [CosmWasm](https://cosmwasm.com/) (smart contracts using WASM)
- [Ethermint](https://ethermint.zone/) (Ethereum virtual machine)

[See the full list here](https://cosmonauts.world/).

## Modules

Cosmos modules are the foundational building blocks for building a blockchain. Modules can be created and shared with anyone. Each module plugs into the Cosmos SDK. Some of the modules are core building blocks for creating blockchains, while other modules enable new features.

Many of the live blockchains use multiple Cosmos modules. The foundational modules for starport are: `auth`, `bank`, `staking`, `params` and `supply`. We also recommend adding the `wasm` or the `evm` module, this allows you to deploy Web Assembly smart contracts to your blockchain. The `evm` module enables the Ethereum EVM to be used in your blockchain. Each module comes with clear documentation and codebase. If you wanted to make changes to a specific module, you can fork the module and change what suits your use case better.

## Summary

- Starport lets you create, develop, and build a blockchain.
- Starport and Cosmos are written in Go.
- Today, Cosmos SDK has a unique position worldwide as one of the most successful blockchains.
- Developers can use different Cosmos SDK modules to customize their blockchain.
