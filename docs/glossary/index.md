---
title: Glossary
order: 0
parent:
  order: 0
  title: Developer Tutorials
description: Learn more about the language of Starport.
---

Learn more about the language of Starport. Terms are defined in the context of blockchain networks.

# Glossary 

<!-- add navigation A to Z -->

## A

### account

A pair of private key/public key/addresses.

### application specific blockchain

A special purpose blockchain that is customized to operate a single application.

### ATOM

The native staking currency of the Cosmos Hub. 

## B



## C

### campaign

A project to start a mainnet chain network. A campaign can be a sequence of testnets that precede and result in a mainnet launch.

### chain 

The project of starting a new blockchain network. When you create a chain, you create a project to start a new blockchain network.

### chain information

Data about a blockchain that validators require to perform a decentralized start of a blockchain network, including chain ID, source URL, and hash.

### coordinator 

An entity that initiates and manages chain launches and campaigns.

### Cosmos 

A network of interoperable blockchains.

### Cosmos Hub

The first of interconnected blockchains that comprise the Cosmos Network.

### Cosmos Network

A decentralized network of independent, scalable, and interoperable blockchains.

### Cosmos SDK

A modular framework that simplifies the process of building secure blockchain applications.

### CRUD 

The four functions that are considered necessary to implement a persistent storage application: create, read, update and delete. Starport scaffold commands create and modify the source code files to add CRUD functionality.

### current reward height

The current block height tracked for a chain for reward distribution. This value is increased when a monitoring packet is received.

## D

### delegators

A safeguard against validators that exhibit bad behavior. 

### distribution round

A round of reward distribution for a chain. During a distribution round, all the reward pools are iterated to determine reward for each validator.

## E



## F

### faucet

An application that dispenses cryptocurrency for use on test networks. Used by developers for testing before launching a chain to mainnet. Tokens that are dispensed by a test faucet cannot be exchanged for mainnet equivalents. 

### finality

A transaction is final after the block that contains it is validated. After finality is reached, the transaction is immutable.

### full node 

A node on the blockchain network that stores and verifies the entire state of a blockchain. 

## G

### Gaia

The name of the Cosmos SDK application for the Cosmos Hub.

### genesis block

The initial block of data in the history of a blockchain network.

## H



## I

### IBC 

A protocol that allows blockchains to talk to each other. The backbone of the Cosmos ecosystem, the Inter-Blockchain Communication protocol (IBC) handles reliable transport across different sovereign blockchains. This end-to-end, connection-oriented, stateful protocol provides reliable, ordered, and authenticated communication between heterogeneous blockchains. 

### IBC module

The standard for the interaction between two blockchains in the Cosmos SDK. The IBC module defines how packets and messages are constructed to be interpreted by the sending and the receiving blockchain.

### IBC relayer

The IBC relayer lets you connect between sets of IBC-enabled chains. A built-in relayer in Starport connects blockchains that run on your local computer to blockchains that run on remote computers. The Starport relayer uses the TypeScript relayer.

### incentivized block

A block that is eligible for a reward for the validator that signs it.

## J



## K



## L

### launch information

The set of information relied on to deterministically launch the network of a chain, including the information required to generate the genesis file and the list of persistent peers for the network.

### last block height

The height of the last block that was monitored in a monitoring round.

### last reward height

The block height of the last incentivized block in a testnet for a reward pool.

### light node 

A program that processes only block headers and a small subset of transactions. 

## M

### mainnet

Short for main network, the original and functional blockchain where actual transactions take place in the distributed ledger.

### module

A subset of the state-machine. A Cosmos SDK app is usually composed of an aggregation of modules.

### monitoring packet 

An atomic IBC packet that is sent from a testnet to SPN that contains data from the monitoring period.

### monitoring period

A number of blocks specific to the testnet to monitor validator activities. 

### monitoring round

The monitored validator activities during a round. A round is a variable number of blocks. A round ends when the validator set of the testnet changes or when the monitoring period ends.

### minting 

The process of creating new units of cryptocurrency tokens.

### mnemonic

A sequence of words that is used as seed to derive private keys. The mnemonic is at the core of each wallet.

## N

### node

Any computer that is connected to a blockchain network is referred to as a node.

## O

### oracle

A communication module that has built-in compliance with IBC protocol to query data points of various types.

## P

### pool reward

The percentage of reward for a given reward pool that is distributed during a distribution round. The percentage depends on the current reward height and the last reward height of the reward pool.

## Q



## R

### refund per pool

The percentage of reward that is refunded to the reward provider during a distribution round. The percentage for each reward pool depends on the number of unsigned blocks.

### relaying tip

The reward offered for relaying a monitoring packet.

### relaying tip pool

A pool that contains relaying tips. The pool determines how the tips are distributed by defining the number of tips per block. Distribution instructions cannot be modified, but the tips contained in the pool can be refunded to the provider. A testnet can have any number of associated pools. 

### request 

Solicit a change in the launch information of the chain. To be effective, a request must be approved by a coordinator. 

### request pool

The list of requests that are waiting to be approved or rejected by the chain coordinator.

### reward

Tokens that exist on SPN are sent into an escrow for direct payment to validators to incentivize and reward them to participate in a blockchain network.

### reward pool 

A pool that is associated with a testnet that contains rewards to be distributed. The pool determines how to distribute the rewards. For example, how much reward per block, block range for distribution. The distribution instructions for distribution cannot be modified but the rewards contained in the pool can be refunded to the provider. A testnet can have any number of associated pools.

## S

### scaffold

A code generation function. Starport comes with a number of scaffolding commands that are designed to make development easier by creating all of the code that's required to start working on a particular task. For example, to build a fully functional Cosmos SDK blockchain foundation, use the `starport scaffold chain` command.

### scalability 

The ability of a network to continue functioning when the number of actors increases to infinity.

### share

A percentage of the total supply. Each coin defined in the total supply has an associated share. If an entity holds 100 of the 1000 token shares, that entity receives 10% of the supply of the token.

### signature counts

A structure that contains the number of signatures with their associated consensus address in a monitoring round. This structure also stores the total number of signatures and the number of unsigned blocks by validator to determine rewards for validators.

### sovereign blockchain app

An app whose governance system has full authority over the blockchain on which it runs. This governance includes having its own independent validator set.

### start reward height

The block height of the first incentivized block in a testnet for a reward pool.

### state machine 

A program that holds a state and modifies the state when it receives inputs. For example, a list of accounts and balances for a cryptocurrency. Each block in a blockchain represents a change to the state. 

## T

### testnet

A blockchain network dedicated to testing the software.

### token

A digital asset that is usually stored and secured by a blockchain.

### transactions 

Signed messages that trigger state transitions.

### transaction fees 

Transactions in a blockchain network can include a fee paid with tokens in order for the transaction to be processed.

## U



## V

### validator, Cosmos SDK 

A special full-node that takes part in the consensus algorithm to collectively add blocks to the blockchain.

### validator, SPN

An entity that participates in a chain launch or a campaign with an intent to become a validator. 

### validator reward per pool

The percentage of reward that is distributed to a validator  during a distribution round for a given reward pool. The percentage depends depending on the number of signatures associated with the validator's address.

### voucher

A tokenized share that can be transferred. 

## W



## X



## Y



## Z



