---
order: 0
parent:
  title: "Sending Tokens: Nameservice"
  order: 4
  description: Build a blockchain app that lets users buy names and set a value these names resolve to. 
---

# Nameservice Tutorial

The nameservice tutorial provides step-by-step instructions to build a blockchain app for a nameservice. This nameservice app sends tokens between participants so that end users can buy names and set a value to the names. 
<!-- rewrite intro for clarity 
- what problem does this nameservice solve?
- why should user care?
-->

This tutorial builds on knowlege and skills developed in the earlier tutorials in the Starport Developer Guide. Before you start this building your nameservice app, we recommend that you complete these foundational tutorials:

- [Install Starport](../install.md)
- [Hello, World](../hello.md)
- [Module Basics](../blog.md)

This tutorial helps you build a blockchain for a nameservice app with these steps:

* Create a blockchain without a default module
* Create a Cosmos SDK nameservice module with a dependency on another module
* Create CRUD (create, read, update, and delete) actions for a type stored as a map
* Declare functions of the bank module to be available to the nameservice module
* Implement keeper functions that implement the logic

The goal of this tutorial is to build a functional nameservice app and a mapping of strings to other strings (`map[string]string`).
<!-- rewrite for clarity
 This is similar to [Namecoin](https://namecoin.org/), [ENS](https://ens.domains/), or [Handshake](https://handshake.org/), which all model the traditional DNS systems (`map[domain]zonefile`). App end users can buy unused names, sell their names, or trade their names.
-->

## Prerequisites 

- A supported version of [Starport](https://docs.starport.network/). To install Starport, see [Install Starport](../starport/install.md). 
* A text editor like [Visual Studio Code](https://code.visualstudio.com/download) or [Atom](https://atom.io/). 
* A web browser like [Chrome](https://www.google.com/chrome/) or [Firefox](https://www.mozilla.org/en-US/firefox/new/).
- Familiarity with [Cosmos SDK modules](https://docs.cosmos.network/master/building-modules/intro.html) 

## Nameservice App Goals

The goal of the app you are building is to let users buy names and to set a value these names resolve to. The owner of a given name is the current highest bidder. 

First, see how these simple requirements translate to app design. 

### Core Concepts 
<!-- needs so much love 
- A blockchain app is a [replicated deterministic state machine](https://en.wikipedia.org/wiki/State_machine_replication). 
- As a developer, you have to define the state machine with a starting state and messages that trigger state transitions. 
-->

- [Starport](https://docs.starport.network/) is built on top of Cosmos SDK and accelerates chain development by scaffolding everything you need. 
- The [Cosmos SDK](https://github.com/cosmos/cosmos-sdk/) modular framework allows developers like you to create custom blockchains that can natively interact with other blockchains. 
- [Tendermint](https://docs.tendermint.com/master/introduction/what-is-tendermint.html) software securely and consistently replicates an app on many machines. The Tendermint app-agnostic engine handles the networking and consensus layers of your blockchain. 

## Cosmos SDK Modules 

In a Cosmos SDK blockchain, application-specific logic is implemented in separate modules. Modules keep code easy to understand and reuse. Each module contains its own message and transaction processor, while the Cosmos SDK is responsible for routing each message to its respective module.

Your nameservice app requires the following Cosmos SDK modules:

- [auth](https://docs.cosmos.network/master/modules/auth/): Specifies the base transaction and account types for an application. For your nameservice app, it defines accounts and fees and gives access to these functionalities to the rest of your app.
- [bank](https://docs.cosmos.network/master/modules/bank/): Enables the app to create and manage tokens and token balances.
- [distribution](https://docs.cosmos.network/master/modules/distribution/): This module give a functional way to passively distribute rewards between validators and delegators.
- [slashing](https://docs.cosmos.network/master/modules/slashing/): Enables punishing misbehavior of validators when evidence of validator fraud is reported.
- [staking](https://docs.cosmos.network/master/modules/staking/): Enables the app to have validators that users can delegate to.
- nameservice: This module does not exist yet! You will build this module to handle the core logic for your new `nameservice` app. The `nameservice` module is the main piece of software you develop to build your app.

Now, take a look at the two main parts of your app: the state and the message types.

## Application State

The state represents your app at a given moment. The state of your nameservice app defines how many tokens are in each account, who the account owners are, the price of each name, and to what value each name resolves to.

The state of tokens and accounts is defined by the `auth` and `bank` modules so you can direct your focus instead on defining the part of the state that relates specifically to your `nameservice` module.

The Cosmos SDK comes with a large set of stores to persist the state of applications. By default, the main store of Cosmos SDK apps is a multistore (a store of stores). You can add any number of key-value stores [KVStores in Go](https://pkg.go.dev/github.com/cosmos/cosmos-sdk/types#KVStore) to the multistore. 

For your nameservice app, use one store to map a `name` key to its respective `whois` value that holds a name's value, owner, and price.

## Messages

In the Cosmos SDK, [messages](https://docs.cosmos.network/master/building-modules/messages-and-queries.html#messages) are objects that are contained in transactions to trigger state transitions. Each Cosmos SDK module defines a list of messages and how to handle them. 

You must create [messages for the nameservice module](02-messages.md) that support this functionality:

- When a transaction that is included in a block reaches a Tendermint node, the transaction is passed to the application using the [ABCI](https://docs.cosmos.network/master/intro/sdk-app-architecture.html#abci) interface between Tendermint and your app. 
- The transaction is decoded to get the message. 
- The message is then routed to the appropriate module and handled according to the logic defined in the corresponding `Handler`. 
- If the state needs to be updated, the `Handler` calls the `Keeper` to perform the update. 

You will learn more about these concepts in the next steps of this tutorial.

Now that you have an idea of how your app functions from a high-level perspective, it is time to start implementing it.
