---
sidebar_position: 2
description: Add messages to define actions for the nameservice module.
---

# Messages for the Nameservice Module

Messages are a great place to start when building a Cosmos SDK module because they define the actions that your app can make. Remember that the nameservice app lets users buy a name, set a value for a name to resolve to, and delete a name that belongs to them.

With this design in mind for the `nameservice` module, it's time to create these messages to define the actions. End users can send these messages to interact with the application state:

- `BuyName`
- `SetName`
- `DeleteName`

## Message Type

Messages trigger state transitions. Messages (`Msg`) are wrapped in transactions (`Tx`) that clients submit to the network. Because the Cosmos SDK wraps and unwraps messages from transactions, as an app developer, you only have to define messages.

Messages must satisfy the following interface:

```go
// Transactions messages must fulfill the Msg
type Msg interface {
	proto.Message

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() error

	// Signers returns the addrs of signers that must sign.
	// CONTRACT: All signatures must be present to be valid
	GetSigners() []AccAddress

	// Legacy methods
	Type() string
	Route() string
	GetSignBytes() []byte
}
```

The `Msg` type extends `proto.Message` and contains these methods along with the legacy methods (`Type`, `Route`, and `GetSignBytes`):

- `ValidateBasic`

  - Called early in the processing of the message to discard obviously invalid messages.
	- Includes only checks that do not require access to the state. For example, check that the `amount` of tokens is a positive value.

- `GetSigners`

  - Returns the list of signers.
  - The Cosmos SDK ensures that each message contained in a transaction is signed by all the signers in the list that is returned by this method.

## Handlers

Handlers define the action that needs to be taken. Each message has an associated handler.

For example, handlers define which stores to update, how to update the stores, and under what conditions to act when a given message is received.

## Scaffolding Messages

Now, you are ready to implement these Cosmos SDK messages to achieve the desired functionality for your nameservice app:

- `MsgBuyName`
	Allow accounts to buy a name and become its owner. When an end user buys a name, they are required to pay the previous owner of the name a price higher than the price the previous owner paid for it. If a name does not have a previous owner yet, the end user must burn a `MinPrice` amount.
- `MsgSetName`
	Allow name owners to set a value for a given name.
- `MsgDeleteName`
	Allow name owners to delete names that belong to them.

Use the `ignite scaffold message` command to scaffold new messages for your module.

- The [`ignite scaffold message`](../../../references/cli#ignite-scaffold-message) command accepts the message name as the first argument and a list of fields for the message. 
- By default, a message is scaffolded in a module with a name that matches the name of the project, in this case `nameservice`.

### Add the MsgBuyName Message

To create the `MsgBuyName` message for the nameservice module:

```bash
ignite scaffold message buy-name name bid
```

where:

- buy-name is the message name
- name defines the name that the user can buy, sell, and delete
- bid is the price the user bids to buy a name

The `ignite scaffold message buy-name name bid` command creates and modifies several files:

```
modify proto/nameservice/tx.proto
modify x/nameservice/client/cli/tx.go
create x/nameservice/client/cli/tx_buy_name.go
create x/nameservice/keeper/msg_server_buy_name.go
modify x/nameservice/types/codec.go
create x/nameservice/types/message_buy_name.go
```

These are the changes for each one of these files:

- `proto/nameservice/tx.proto`
    - Adds `MsgBuyName` and `MsgBuyNameResponse` proto messages.
    - Registers `BuyName` rpc in the `Msg` service.

    Open the `tx.proto` file to view the changes:

    ```protobuf
    syntax = "proto3";

    package nameservice.nameservice;

    // this line is used by starport scaffolding # proto/tx/import

    option go_package = "nameservice/x/nameservice/types";

    // Msg defines the Msg service.
    service Msg {
      // this line is used by starport scaffolding # proto/tx/rpc
      rpc BuyName(MsgBuyName) returns (MsgBuyNameResponse);
    }

    // this line is used by starport scaffolding # proto/tx/message
    message MsgBuyName {
      string creator = 1;
      string name = 2;
      string bid = 3;
    }

    message MsgBuyNameResponse {
    }
    ```

- `x/nameservice/client/cli/tx.go`

    Registers the CLI command.

- `x/nameservice/types/message_buy_name.go`

    Defines methods to satisfy the `Msg` interface.

- `x/nameservice/keeper/msg_server_buy_name.go`

    Defines the `BuyName` keeper method. You can notice that the message follows the `Msg` interface. The message `struct` contains all the  information required when buying a name: `Name`, `Bid`, and `Creator`. This struct was added automatically.

- `x/nameservice/client/cli/tx_buy_name.go`

  	Adds the CLI command to broadcast a transaction with a message.

- `x/nameservice/types/codec.go`

    Registers the codecs.


### Add The MsgSetName Message

To create the `MsgSetName` for the nameservice module:

```bash
ignite scaffold message set-name name value
```

where:

- set-name is the message name
- name is the name the user sets
- value is the literal value that the name resolves to

This `ignite scaffold message` command modifies and creates the same set of files as the `MsgBuyName` message.

### Add The MsgDeleteName Message

You need a message so that an end user can delete a name that belongs to them.

To create the `MsgDeleteName` for the nameservice module:

```bash
ignite scaffold message delete-name name
```

where:

- delete-name is the message name
- name is message name to delete

## Results

Congratulations, you've defined messages that trigger state transitions. Now it's time to implement types and methods that operate on the state.
