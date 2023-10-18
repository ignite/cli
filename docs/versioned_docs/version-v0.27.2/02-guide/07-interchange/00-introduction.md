---
sidebar_position: 0
slug: /guide/interchange
---

# Introduction

The Interchain Exchange is a module to create buy and sell orders between blockchains.

In this tutorial, you learn how to create a Cosmos SDK module that can create order pairs, buy orders, and sell orders.
You create order books and buy and sell orders across blockchains, which in turn enables you to swap token from one
blockchain to another.

**Note:** The code in this tutorial is written specifically for this tutorial and is intended only for educational
purposes. This tutorial code is not intended to be used in production.

If you want to see the end result, see the example implementation in
the [interchange repo](https://github.com/tendermint/interchange).

**You will learn how to:**

- Create a blockchain with Ignite CLI
- Create a Cosmos SDK IBC module
- Create an order book that hosts buy and sell orders with a module
- Send IBC packets from one blockchain to another
- Deal with timeouts and acknowledgements of IBC packets

## How the Interchange Exchange Module Works

To build an exchange that works with two or more blockchains, follow the steps in this tutorial to create a Cosmos SDK
module called `dex`.

The new `dex` module allows you to open an exchange order book for a pair of token: a token from one blockchain and a token
on another blockchain. The blockchains are required to have the `dex` module available.

Token can be bought or sold with limit orders on a simple order book. In this tutorial, there is no notion of a
liquidity pool or automated market maker (AMM).

The market is unidirectional:

- The token sold on the source chain cannot be bought back as it is
- The token bought from the target chain cannot be sold back using the same pair.

If a token on a source chain is sold, it can only be bought back by creating a new pair on the order book.
This workflow is due to the nature of the Inter-Blockchain Communication protocol (IBC) which creates a `voucher`
token on the target blockchain. There is a difference of a native blockchain token and a `voucher` token that is minted
on another blockchain. You must create a second order book pair in order to receive the native token back.

In the next chapter, you learn details about the design of the interblockchain exchange.
