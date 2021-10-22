---
order: 13
description: Denoms on Starport and the Cosmos SDK
---

# Denom

## What is the denom

Denom stands for `denomination`.

A denom is the name of a token that can be used for all purposes with Starport and in the Cosmos SDK.

You can set an arbitrary number of denoms before starting your blockchain in the config.yml file at the root of the repository.

Most example denoms take the format of `token` or `stake`.

## How is the denom used

Assets in the Cosmos SDK are represented via a `Coins` type that consists of an amount and a denom, where the amount can be any arbitrarily large or small value. In addition, the Cosmos SDK uses an account-based model where there are two types of primary accounts -- basic accounts and module accounts. All account types have a set of balances that are composed of Coins. The x/bank module keeps track of all balances for all accounts and also keeps track of the total supply of balances in an application.

With regards to a balance amount, the Cosmos SDK assumes a static and fixed unit of denomination, regardless of the denomination itself. In other words, clients and apps built atop a Cosmos-SDK-based chain may choose to define and use arbitrary units of denomination to provide a richer UX, however, by the time a tx or operation reaches the Cosmos SDK state machine, the amount is treated as a single unit. For example, for the Cosmos Hub (Gaia), clients assume 1 ATOM = 10^6 uatom, and so all txs and operations in the Cosmos SDK work off of units of 10^6.

## Denoms with IBC

The most used feature of IBC is to send tokens from one blockchain to another. When sending a token to another blockchain, a token `voucher` is generated on the other (target) blockchain.

The voucher token denom is represented by the format of starting with `ibc/`.

This allows apps to identify IBC tokens on a blockchain and deal with them as approriate vourcher token.

With IBC tokens enabled, you can consider native token in a blockchain with the reference `voucher` token from another blockchain. They can be differentiated by the `denom` name of the token.

[Learn how to parse IBC denoms and how they are encoded](https://tutorials.cosmos.network/understanding-ibc-denoms/).
