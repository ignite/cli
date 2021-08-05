---
order: 2
---

# App Init

##  Initialize the Blockchain

In this chapter you create the basic blockchain module for the interchain exchange app. You scaffold the blockchain, the module, the transaction, the IBC packets and messages. In the later chapters you  integrate more code into each of the transaction handlers.

## Create the Blockchain

Scaffold a new blockchain called `interchange`

```bash
starport app github.com/username/interchange
cd interchange
```

A new directory named interchange is created. This directory contains a working blockchain app.
Create a new IBC module next.

## Create the Module

Scaffold a module inside your blockchain named `ibcdex` with IBC capabilities.
The ibcdex module contains the logic for creating and maintaining order books and routing them through IBC to the second blockchain.

```bash
starport module create ibcdex --ibc --ordering unordered
```

## Create CRUD logic for Buy and Sell Order Books

To scaffold two types with create, read, update and delete (CRUD) actions use the Starport `type` command.
The following commands create `sellOrderBook` and `buyOrderBook` types. 

```bash
starport type sellOrderBook amountDenom priceDenom --indexed --no-message --module ibcdex
starport type buyOrderBook amountDenom priceDenom --indexed --no-message --module ibcdex
```

The values are: 

- `amountDenom` represents which token will be sold and in which quantity
- `priceDenom` the token selling price 

The flag `--indexed` flag creates an "indexed type". Without this flag, a type is implemented like a list with new items appended. Indexed types act like key-value stores.

The `--module ibcdex` flag specifies that the type should be scaffolded in the `ibcdex` module.

## Create the IBC Packets

Create three packets for IBC:

- an order book pair `createPair` 
- a sell order `sellOrder` 
- a buy order `buyOrder`

```bash
starport packet createPair sourceDenom targetDenom --module ibcdex
starport packet sellOrder amountDenom amount:int priceDenom price:int --ack remainingAmount:int,gain:int --module ibcdex
starport packet buyOrder amountDenom amount:int priceDenom price:int --ack remainingAmount:int,purchase:int --module ibcdex
```

The optional `--ack` flag defines field names and types of the acknowledgment returned after the packet has been received by the target chain. Value of `--ack` is a comma-separated (no spaces) list of names with optional types appended after a colon.

## Cancel messages

Cancelling orders is done locally in the network, there is no packet to send.
Use the `message` command to create a message to cancel a sell or buy order.

```go
starport message cancelSellOrder port channel amountDenom priceDenom orderID:int --desc "Cancel a sell order" --module ibcdex
starport message cancelBuyOrder port channel amountDenom priceDenom orderID:int --desc "Cancel a buy order" --module ibcdex
```

The optional `--desc` flag lets you define a description of the CLI command that is used to broadcast a transaction with the message.

## Trace the Denom

The token denoms must have the same behavior as described in the `ibc-transfer` module:

- An external token received from a chain has a unique `denom`, reffered to as `voucher`.
- When a token is sent to a blockchain and then sent back and received, the chain can resolve the voucher and convert it back to the original token denomination.

`Voucher` tokens are represented as hashes, therefore you must store which original denomination is related to a voucher, you can do this with an indexed type.

For a `voucher` you store: the source port ID, source channel ID and the original denom

```go
starport type denomTrace port channel origin --indexed --no-message --module ibcdex
```

## Create the Configuration

Add two config files `mars.yml` and `venus.yml` to test two blockchain networks with specific token for each.
Add the config files in the `interchange` folder.
The native denoms for Mars are `mcx`, also known as `marscoin`, and for Venus `vcx`, also known as `venuscoin`.

Create the `mars.yml` file with your content:

```yaml
# mars.yml
accounts:
  - name: alice
    coins: ["1000token", "100000000stake", "1000mcx"]
  - name: bob
    coins: ["500token", "1000mcx", "100000000stake"]
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
    coins: ["1000token", "1000000000stake", "1000vcx"]
  - name: bob
    coins: ["500token", "1000vcx", "100000000stake"]
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
  grpc: ":9091"
  api: ":1318"
  frontend: ":8081"
  dev-ui: ":12346"
genesis:
  chain_id: "venus"
init:
  home: "$HOME/.venus"
```

Implement the code for the order book in the next chapter.
