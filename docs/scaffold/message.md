---
order: 5
description: Generate and modify message files and directories.
---

# Message Scaffold

Cosmos SDK messages modify the state of a blockchain. Messages are bundled into transactions, broadcasted transactions are bundled into blocks, and blocks make a blockchain.

While `type` constructs a message for each CRUD action and implements CRUD logic, message scaffolding constructs a single message without logic.

## Construct Cosmos SDK Messages

Use the `starport scaffold message` command to scaffold Cosmos SDK messages. See [Cosmos SDK messages](https://docs.cosmos.network/v0.42/building-modules/messages-and-queries.html).

```
starport scaffold message [name] [field1] [field2] ... [flags]
```

## Files and Directories

The following files and directories are created and modified by scaffolding:

- `proto`: the message type.
- `x/module_name/keeper`: the gRPC message server.
- `x/module_name/types`: message type definitions and keys
- `x/module_name/client/cli`: the CLI used for broadcasting a transaction with the message.

## Describe Messages


All flags are optional.

`--desc`

The description of the CLI command that broadcasts a transaction with a message.

`--response`

Comma-separated list (no spaces) of fields that describe the response fields of the message.

## Create Message Example

This example command creates the `cancelSellOrder` message with two fields

```
starport scaffold message cancelSellOrder port channel amountDenom priceDenom orderID:int --desc "Cancel a sell order" --response id,amount --module ibcdex
```
