---
sidebar_position: 2
description: Create the blockchain for the interchain exchange app.
---

# App Init

## Initialize the Blockchain

In this chapter, you create the basic blockchain module for the interchain exchange app. You scaffold the blockchain, the module, the transaction, the IBC packets, and messages. In later chapters, you integrate more code into each of the transaction handlers.

## Create the Blockchain

Scaffold a new blockchain called `interchange`:

```bash
ignite scaffold chain interchange --no-module
```

A new directory named `interchange` is created. 

Change into this directory where you can scaffold modules, types, and maps:

```bash
cd interchange
```

The `interchange` directory contains a working blockchain app.

A local GitHub repository has been created for you with the initial scaffold.

Next, create a new IBC module.

## Create the dex Module

Scaffold a module inside your blockchain named `dex` with IBC capabilities.

The dex module contains the logic to create and maintain order books and route them through IBC to the second blockchain.

```bash
ignite scaffold module dex --ibc --ordering unordered --dep bank
```

## Create CRUD logic for Buy and Sell Order Books

Scaffold two types with create, read, update, and delete (CRUD) actions. 

Run the following Ignite CLI `type` commands to create `sellOrderBook` and `buyOrderBook` types:

```bash
ignite scaffold map sell-order-book amountDenom priceDenom --no-message --module dex
ignite scaffold map buy-order-book amountDenom priceDenom --no-message --module dex
```

The values are:

- `amountDenom`: the token to be sold and in which quantity
- `priceDenom`: the token selling price

The `--no-message` flag specifies to skip the message creation. Custom messages will be created in the next steps.

The `--module dex` flag specifies to scaffold the type in the `dex` module.

## Create the IBC Packets

Create three packets for IBC:

- An order book pair `createPair`
- A sell order `sellOrder`
- A buy order `buyOrder`

```bash
ignite scaffold packet create-pair sourceDenom targetDenom --module dex
ignite scaffold packet sell-order amountDenom amount:int priceDenom price:int --ack remainingAmount:int,gain:int --module dex
ignite scaffold packet buy-order amountDenom amount:int priceDenom price:int --ack remainingAmount:int,purchase:int --module dex
```

The optional `--ack` flag defines field names and types of the acknowledgment returned after the packet has been received by the target chain. The value of the `--ack` flag is a comma-separated list of names (no spaces). Append optional types after a colon (`:`).

## Cancel messages

Cancelling orders is done locally in the network, there is no packet to send.

Use the `message` command to create a message to cancel a sell or buy order:

```bash
ignite scaffold message cancel-sell-order port channel amountDenom priceDenom orderID:int --desc "Cancel a sell order" --module dex
ignite scaffold message cancel-buy-order port channel amountDenom priceDenom orderID:int --desc "Cancel a buy order" --module dex
```

Use the optional `--desc` flag to define a description of the CLI command that is used to broadcast a transaction with the message.

## Trace the Denom

The token denoms must have the same behavior as described in the `ibc-transfer` module:

- An external token received from a chain has a unique `denom`, reffered to as `voucher`.
- When a token is sent to a blockchain and then sent back and received, the chain can resolve the voucher and convert it back to the original token denomination.

`Voucher` tokens are represented as hashes, therefore you must store which original denomination is related to a voucher. You can do this with an indexed type.

For a `voucher` you store, define the source port ID, source channel ID, and the original denom:

```bash
ignite scaffold map denom-trace port channel origin --no-message --module dex
```

## Create the Configuration for Two Blockchains

Add two config files `mars.yml` and `venus.yml` to test two blockchain networks with specific token for each.

Add the config files in the `interchange` folder.

The native denoms for Mars are `marscoin`, and for Venus `venuscoin`.

Create the `mars.yml` file with your content:

```yaml
# mars.yml
accounts:
  - name: alice
    coins: ["1000token", "100000000stake", "1000marscoin"]
  - name: bob
    coins: ["500token", "1000marscoin", "100000000stake"]
validator:
  name: alice
  staked: "100000000stake"
faucet:
  name: bob
  coins: ["5token", "100000stake"]
genesis:
  chain_id: "mars"
init:
  home: "$HOME/.mars"
```

Create the `venus.yml` file with your content:

```yaml
# venus.yml
accounts:
  - name: alice
    coins: ["1000token", "1000000000stake", "1000venuscoin"]
  - name: bob
    coins: ["500token", "1000venuscoin", "100000000stake"]
validator:
  name: alice
  staked: "100000000stake"
faucet:
  host: ":4501"
  name: bob
  coins: ["5token", "100000stake"]
host:
  rpc: ":26659"
  p2p: ":26658"
  prof: ":6061"
  grpc: ":9092"
  grpc-web: ":9093"
  api: ":1318"
genesis:
  chain_id: "venus"
init:
  home: "$HOME/.venus"
```

In the `venus.yml` file, you can see the specific `host` parameter that you can use to change the ports for various running services (rpc, p2p, prof, grpc, api, frontend, and dev-ui). This `host` parameter can be used later so you can run two blockchains in parallel and prevent conflicts when the chains are using the same ports.

You can also use the `host` parameter to use specific ports for any of the services.

After scaffolding, now is a good time to make a commit to the local GitHub repository that was created for you.

```bash
git add .
git commit -m "Scaffold module, maps, packages and messages for the dex"
```

Implement the code for the order book in the next chapter.
