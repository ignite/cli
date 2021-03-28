# Message Scaffold

Cosmos SDK messages modify the state of a blockchain. Messages are bundled into transactions, broadcasted transactions are bundled into blocks, and blocks make a blockchain.

While `type` constructs a message for each CRUD action and implements CRUD logic, message scaffolding constructs a single message without logic.

You can use the `starport message` command to scaffold Cosmos SDK messages. See [Cosmos SDK messages](https://docs.cosmos.network/v0.42/building-modules/messages-and-queries.html).

```
starport message [name] [field1] [field2] ... [flags]
```

The following files and directories are created and modified by scaffolding:

- **proto** 
     The message type.
- **x/module_name/keeper** 

    The gRPC message server.
- **x/module_name/types** 

    Message type definitions and keys
- **`x/module_name/client/cli** 

    The CLI used for broadcasting a transaction with the message.

All flags are optional.

- **--desc** 

    The description of the CLI command that broadcasts a transaction with a message.

- **--response** 

Comma-separated list (no spaces) of fields that describe the response fields of the message. 

Message scaffolding is similar to type scaffolding, except `type` scaffolds a message for each CRUD action and implements CRUD logic, where `message` scaffolds a single message without any logic.

## Example

```
starport message cancelSellOrder port channel amountDenom priceDenom orderID:int --desc "Cancel a sell order" --response id,amount --module ibcdex
```
