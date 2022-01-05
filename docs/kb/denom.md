---
order: 13
description: Denoms on Starport and the Cosmos SDK
---

# Denom

Denom stands for `denomination`.

A denom is the name of a token that can be used for all purposes with Starport and in the Cosmos SDK.

In Starport, the `config.yml` file that is generated in your blockchain folder describes the configuration for your blockchain. In this file, you can set an arbitrary number of denoms before starting your blockchain.

Mostly, example denoms take the format of `token` or `stake`.

## How a Denom Is Used

Assets in the Cosmos SDK are represented by using a `Coins` type that consists of an amount and a denom, where the amount can be any arbitrarily large or small value. The account-based model in the Cosmos SDK has two types of primary accounts -- basic accounts and module accounts. All account types have a set of balances that are composed of `Coins`. The `x/bank` module keeps track of all balances for all accounts and also keeps track of the total supply of balances in an application.

For a balance amount, the Cosmos SDK assumes a static and fixed unit of denomination, regardless of the denomination itself. Clients and apps built on top of a Cosmos SDK-based chain can choose to define and use arbitrary units of denomination to provide a richer user experience. However, by the time a transaction (`tx`) or operation reaches the Cosmos SDK state machine, the amount is treated as a single unit. For example, for the Cosmos Hub (Gaia), clients assume 1 ATOM = 10^6 uatom. All tx and operations in the Cosmos SDK work off of units of 10^6.

## Denoms with IBC

The most used feature of the Inter-Blockchain Communication protocol (IBC) is to send tokens from one blockchain to another. When sending a token from a source chain to a target blockchain, a token `voucher` is generated on the target blockchain.

The voucher token denom is represented with syntax naming convention that starts with `ibc/`.

This naming convention allows apps to identify IBC tokens on a blockchain and deal with them as approriate voucher tokens.

With IBC tokens enabled, you can consider native token in a blockchain with the reference `voucher` token from another blockchain. The tokens are differentiated by the `denom` name of the token.

See [Understand IBC Denoms with Gaia](https://tutorials.cosmos.network/understanding-ibc-denoms/) to learn more about the format and use of a voucher token.
