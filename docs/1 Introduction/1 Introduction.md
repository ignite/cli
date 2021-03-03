# An introduction to Starport

Starport is a developer-friendly interface to the Cosmos SDK.

## What is Starport?

The Starport tool is the easiest way to create blockchains using the Tendermint Consensus engine and the Cosmos SDK. <!-- @fadeev does it matter that we use Go? --> Starport is written in the Go programming language. The new blockchain has a Proof-of-Stake (PoS) system with validators that can be defined in the genesis block.

With just a few commands, you can use Starport to:

1. Create a blockchain
2. Launch a blockchain
3. Serve your blockchain on the cloud
4. Create a client-side app you can interact with in web browser

<!-- do we need to mention Vue.js? -->

 Starport creates a user interface with Vue.js. This user interface provides a good starting point for developers by creating a browser-based client-side app for your blockchain. The scaffolded app created with Starport includes a command line interface that lets you manage keys, create validators, send tokens.

Note: Starport replaces the `scaffold` program that was previously used to create a blockchain app.

## Projects using Tendermint and Cosmos SDK

Many projects already showcase the Tendermint BFT Consensus Engine and the Cosmos SDK, including:

- [Cosmos](https://github.com/cosmos/gaia) The main IBC Hub and role model of the Cosmos SDK

- [Binance Chain](https://github.com/binance-chain) DEX and utility token

- [Crypto.com Chain](https://github.com/crypto-com/chain-main) Payments, DeFi, and utility token

- [IRIS](https://github.com/irisnet) IBC Hub and developer-oriented

- [Kava](https://github.com/Kava-Labs/kava) DeFi and Stable Coins

- [Aragon](https://docs.chain.aragon.org/) DAO catalyst

- [CosmWasm](https://cosmwasm.com/) Smart contracts using WASM

- [Ethermint](https://ethermint.zone/) Ethereum virtual machine

For the full list, explore the [Cosmos Network](https://cosmonauts.world/) to discover a wide variety of apps, blockchains, wallets, and explorers that are built in the Cosmos ecosystem.

## Cosmos SDK Modules

Cosmos SDK modules are the foundational building blocks for building a blockchain. The Cosmos SDK offers a variety of native modules to make a blockchain work. Modules can be created and shared with anyone. You can use Starport to add modules. The code edits and additions are managed. <!-- link to the architecture modules topic and list high-level benefits -->

## Summary

- Starport lets you create, develop, and build a blockchain.

- Starport and Cosmos SDK are written in Go.

- Today, Cosmos SDK has a unique position worldwide as one of the most successful blockchains.

- Developers can use different Cosmos SDK modules to customize their blockchain.
