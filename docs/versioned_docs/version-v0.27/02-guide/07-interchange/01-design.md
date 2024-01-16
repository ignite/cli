---
sidebar_position: 1
description: Learn about the interchain exchange module design.
---

# App Design

In this chapter, you learn how the interchain exchange module is designed. The module has order books, buy orders, and
sell orders.

- First, create an order book for a pair of token.
- After an order book exists, you can create buy and sell orders for this pair of token.

The module uses the Inter-Blockchain Communication
protocol [IBC](https://github.com/cosmos/ibc/blob/old/ibc/2_IBC_ARCHITECTURE.md).
By using IBC, the module can create order books so that multiple blockchains can interact and exchange their token.

You create an order book pair with a token from one blockchain and another token from another blockchain. In this
tutorial, call the module you create the `dex` module.

> When a user exchanges a token with the `dex` module, a `voucher` of that token is received on the other blockchain.
> This voucher is similar to how an `ibc-transfer` is constructed. Since a blockchain module does not have the rights
> to mint new token of a blockchain into existence, the token on the target chain is locked up, and the buyer receives
> a `voucher` of that token.

This process can be reversed when the `voucher` gets burned to unlock the original token. This exchange process is
explained in more detail throughout the tutorial.

## Assumption of the Design

An order book can be created for the exchange of any tokens between any pair of chains.

- Both blockchains require the `dex` module to be installed and running.
- There can only be one order book for a pair of token at the same time.

<!-- There is no condition to check for open channels between two chains. -->

A specific chain cannot mint new coins of its native token.

<!-- The module is trustless, there is no condition to check when opening a channel between two chains. 
Any pair of tokens can be exchanged between any pair of chains. -->

This module is inspired by the [`ibc transfer`](https://github.com/cosmos/ibc-go/tree/main/modules/apps/transfer)
module on the Cosmos SDK. The `dex` module you create in this tutorial has similarities, like the `voucher` creation.

However, the new `dex` module you are creating is more complex because it supports creation of:

- Several types of packets to send
- Several types of acknowledgments to treat
- More complex logic on how to treat a packet on receipt, on timeout, and more

## Interchain Exchange Overview

Assume you have two blockchains: Venus and Mars.

- The native token on Venus is `venuscoin`.
- The native token on Mars is `marscoin`.

When a token is exchanged from Mars to Venus:

- The Venus blockchain has an IBC `voucher` token with a denom that looks like `ibc/B5CB286...A7B21307F`.
- The long string of characters after `ibc/` is a denom trace hash of a token that was transferred using IBC.

Using the blockchain's API you can get a denom trace from that hash. The denom trace consists of a `base_denom` and a
`path`. In our example:

- The `base_denom` is `marscoin`.
- The `path` contains pairs of ports and channels through which the token has been transferred.

For a single-hop transfer, the `path` is identified by `transfer/channel-0`.

Learn more about token paths
in [ICS 20 Fungible Token Transfer](https://github.com/cosmos/ibc/tree/main/spec/app/ics-020-fungible-token-transfer).

**Note:** This token `ibc/Venus/marscoin` cannot be sold back using the same order book. If you want to "reverse" the
exchange and receive the Mars token back, you must create and use a new order book for the `ibc/Venus/marscoin` to
`marscoin` transfer.

## The Design of the Order Books

As a typical exchange, a new pair implies the creation of an order book with orders to sell `marscoin` or orders to buy
`venuscoin`. Here, you have two chains and this data structure must be split between Mars and Venus.

- Users from chain Mars sell `marscoin`.
- Users from chain Venus buy `marscoin`.

Therefore, we represent:

- All orders to sell `marscoin` on chain Mars.
- All orders to buy `marscoin` on chain Venus.

In this example, blockchain Mars holds the sell orders and blockchain Venus holds the buy orders.

## Exchanging Tokens Back

Like `ibc-transfer`, each blockchain keeps a trace of the token voucher that was created on the other blockchain.

If blockchain Mars sells `marscoin` to chain Venus and `ibc/Venus/marscoin` is minted on Venus then, if
`ibc/Venus/marscoin` is sold back to Mars, the token is unlocked and the token that is received is `marscoin`.

## Features

The features supported by the interchain exchange module are:

- Create an exchange order book for a token pair between two chains
- Send sell orders on source chain
- Send buy orders on target chain
- Cancel sell or buy orders
