# An introduction to Starport

Let us dive into what is special about starport, what we can achieve with it and which other technologies play well with it.

## What is Starport?

Starport is created to make "creating blockchain more comfortable". 

Starport uses the Tendermint Consensus engine and Cosmos SDK to create a blockchain application in Golang. This blockchain has a Proof-of-Stake system with validators (https://en.longhash.com/news/how-cosmos-governance-works-and-how-you-can-become-a-validator) that can be defined in the genesis block (). 

With a few commands or clicks, you can create a blockchain, launch it, server it on the cloud and have a GUI ready to start testing your application.

Boostrapping blockchains was initially the job of scaffold, which helps to create a blockchain application. Starport takes it to the next level and also creates a User Interface with Vue that gives a headstart to interacting with the blockchain more comfortable than with only the CLI tool. 

The CLI will still be there for you, just like scaffold, your blockchain comes with the whole Command Line Interface that enables managing keys, creating validators, sending tokens and let's you add the commands necessary to build your project.

## Projects using Tendermint / Cosmos SDK

There are many projects already showcasing that the Tendermint BFT Consensus and the Cosmos SDK enables a variety of usecases and with a little effort can become part of the most robust blockchains that exist.

The technology stack is being used by 

- Cosmos (Main IBC Hub and "Rolemodel" of the Cosmos SDK)
- Binance Chain (DEX and utility Token)
- IRIS (IBC Hub and developer oriented)
- Kava (DeFi and Stable Coins)
- Aragon (DAO catalyst)
- CosmWasm (ETH EVM machine)
...

The basis of the consensus protocol, network management and blockchain initialisation has been proven to be reliable and competitive. From there, only your imagination can stop you from what can be built with the technology stack.

## Modules

What empowers building blockchains and your own application on top of it are the Cosmos modules. Modules are developed by companies and communities around the Cosmos SDK. Some of the modules are crucial to build a basic blockchain while others enable new features for a blockchain.

Each of the running blockchains most of the time have at least one module developed on their own to enable the features it specialises on. These standards solve the clearity of a project and the understanding of all the different parts that create a specific blockchain application.

The most basic set of modules that are created with starport are `auth`, `bank`, `staking`, `params` and `supply`. If you would add a module like `wasm`, this would enable the Ethereum EVM to be used on your blockchain. Each module comes with a clear documentation and codebase. If you wanted changes to a specific module, you can fork the module and change what suits your blockchain application better.


## Summary

- Starport lets you create, develop, host in the cloud and manage the initiation of the blockchain.
- Starport and Cosmos are written in Go
- Today, Cosmos SDK has a unique position worldwide of the most successful blockchains being built with it.
- A combination of different modules on the Cosmos SDK will create different technology opportunities of a blockchain.


[▶️ Next - Documentation](../../01_introduction/02_documentation_specification/02_documentation_specification.md)  