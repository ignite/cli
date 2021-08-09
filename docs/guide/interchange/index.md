---
order: 1
parent:
  title: "Advanced Module: Interchange"
---

# Introduction 

The Interchain Exchange is a module to create buy and sell orders between blockchains.

In this tutorial you will learn how to create a Cosmos SDK module that can create order pairs, buy and sell orders. You will be able to create order books, buy and sell orders across blockchains, which enables to swap tokens from one blockchain to another.

The code in this tutorial is purely written for a tutorial and only for educational purpose. It is not intended to be used in production.

If you want to see the end result, please refer to the [example implementation](https://github.com/tendermint/interchange).

**You will learn how to:**

- Create a blockchain with Starport
- Create a Cosmos SDK IBC module
- Create an order book that hosts buy and sell orders with a module
- Send IBC packets from one blockchain to another
- Deal with timeouts and acknowledgements of IBC packets

## How the module works

You will learn how to build an exchange that works with two or more blockchains. The module is called `ibcdex`.

The module allows to open an exchange order book between a pair of token from one blockchain and a token on another blockchain. The blockchains are required to have the `ibcdex` module available.

Tokens can be bought or sold with Limit Orders on a simple order book, there is no notion of Liquidity Pool or AMM.

The market is unidirectional: the token sold on the source chain cannot be bought back as it is, and the token bought from the target chain cannot be sold back using the same pair. If a token on a source chain is sold, it can only be bought back by creating a new pair on the order book. This is due to the nature of IBC, creating a `voucher` token on the target blockchain. In this tutorial you will learn the difference of a native blockchain token and a `voucher` token that is minted on another blockchain. You will learn how to create a second order book pair in order to receive the native token back.

In the next chapter you will learn details about the design of the interblockchain exchange.