---
order: 2
---

# App Design

In this chapter you will learn how the interchain exchange module is designed.
The module has order books, buy- and sell orders.
First, an order book for a pair of token has to be created.
After an order book exists, you can create buy and sell orders for this pair of token.

The module will make use of the Interblockchain Communication Standard [IBC](https://github.com/cosmos/ics/blob/master/ibc/2_IBC_ARCHITECTURE.md). With use of the IBC, the module can create order books for tokens to have multiple blockchains interact and exchange their tokens. 
You will be able to create an order book pair with one token from one blockchain and another token from another blockchain. We will call the module you create in this tutorial `ibcdex`.
Both blockchains will need to have the `ibcdex` module installed and running.

When a user exchanges a token with the `ibcdex`, you receive a `voucher` of that token on the other blockchain.
This is similar to how a `ibc-transfer` is constructed.
Since a blockchain module does not have the rights to mint new token of a blockchain into existence, the token on the target chain would be locked up and the buyer would receive a `voucher` of that token.
This process can be reversed when the `voucher` get burned again to unlock the original token. This will be explained throghout the tutorial in more detail.

## Assumption

An order book can be created for the exchange of any tokens between any pair of chains. The requirement is to have the `ibcdex` module available. There can only be one order book for a pair of token at the same time.
<!-- There is no condition to check for open channels between two chains. -->
A specific chain cannot mint new of its native token. 
<!-- The module is trustless, there is no condition to check when opening a channel between two chains. Any pair of tokens can be exchanged between any pair of chains. -->

This module is inspired by the [`ibc-transfer`](https://github.com/cosmos/cosmos-sdk/tree/v0.42.1/x/ibc/applications/transfer) module and will have some similarities, like the `voucher` creation. It will be more complex but it will display how to create:

- Several types of packets to send
- Several types of acknowledgments to treat
- Some more complex logic on how to treat a packet on receipt, on timeout and more

## Overview

Assume you have two blockchains: `Venus` and `Mars`. The native token on Venus is called `vcx`, the token on Mars is `mcx`. 
When exchanging a token from Mars to Venus, on the Venus blockchain you would end up with an IBC `voucher` token with a denom that looks like `ibc/B5CB286...A7B21307F `. The long string of characters after `ibc/` is a denom trace hash of a token transferred through IBC. Using the blockchain's API you can get a denom trace from that hash. Denom trace consists of a `base_denom` and a `path`. In our example, `base_denom` will be `mcx` and the `path` will contain pairs of ports and channels through which the token has been transferred. For a single-hop transfer `path` will look like `transfer/channel-0`. Learn more about token paths in [ICS 20](https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer).
This token `ibc/Venus/mcx` cannot be sold back using the same order book. If you want to "reverse" the exchange and receive back the Mars token, a new order book `ibc/Venux/mcx` to `mcx` needs to be created.

## Order books

As a typical exchange, a new pair implies the creation of an order book with orders to sell `MCX` or orders to buy `VCX`. Here, you have two chains and this data-structure must be split between `Mars` and `Venus`.

Users from chain `Mars` will sell `MCX` and users from chain `Venus` will buy `MCX`. Therefore, we represent all orders to sell `MCX` on chain `Mars` and all the orders to buy `MCX` on chain `Venus`.

In this example blockchain `Mars` holds the sell oders and blockchain `Venus` holds the buy orders.

## Exchanging tokens back

Like `ibc-transfer` each blockchain keep a trace of the token voucher created on the other blockchain.

If a blockchain `Mars` sells `MCX` to `Venus` and `ibc/Venus/mcx` is minted on `Venus` then, if `ibc/Venus/mcx` is sold back on `Mars` the token unlocked and received will be `MCX`.

## Features

The features supported by the module are:

- Creating an exchange order book for a token pair between two chains
- Send sell orders on source chain
- Send buy orders on target chain
- Cancel sell or buy orders
