---
order: 2
description: Add messages to define actions for the nameservice module. 
---

# Messages for the Nameservice Module

Messages are a great place to start when building a Cosmos SDK module because they define the actions that your app can make.

## Message Type

Messages trigger state transitions. Messages (`Msg`) are wrapped in transactions (`Tx`) that clients submit to the network. The Cosmos SDK wraps and unwraps `Msg` from `Tx`, which means, as an app developer, you only have to define messages. 

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

The `Msg` type extends `proto.Message` and contains three legacy methods (`Type`, `Route`, and `GetSignBytes`) and these methods:

- `ValidateBasic` is called early in the processing of the message to discard obviously invalid messages. `ValidateBasic` includes only checks that do not require access to the state. For example, check that the `amount` of tokens is a positive value.

- `GetSigners` returns the list of signers. The Cosmos SDK ensures that each message contained in a transaction is signed by all the signers in the list that is returned by this method.

## Handlers

Handlers define the action that needs to be taken. For example, which stores to update, how to update, and under what conditions to act when a given `Msg` is received.

In your `nameservice` module, three types of `Msgs` can be sent to interact with the application state:

- `BuyName` 
- `SetName`
- `DeleteName`

Each message has an associated `Handler`.

## Scaffolding Messages

You must implement these three messages to achieve the desired functionality for your nameservice app:

- `MsgBuyName`: Allow accounts to buy a name and become its owner. When an end user buys a name, they are required to pay the previous owner of the name a price higher than the price the previous owner paid for it. If a name does not have a previous owner yet, the end user must burn a `MinPrice` amount.
- `MsgSetName`: Allow name owners to set a value for a given name.
- `MsgDeleteName`: Allow name owners to delete names that belong to them.

Use the `starport scaffold message` command to scaffold a new Cosmos SDK message for your module. 

The [starport scaffold message](https://docs.starport.network/cli/#starport-scaffold-message) command accepts message name as the first argument and a list of fields. By default, a message is scaffolded in a module with a name that matches the name of the project, in this case `nameservice`. 

### Add The MsgBuyName Message

To create the `MsgBuyName` for the nameservice module:

```bash
starport scaffold message buy-name name bid
```

where:

- buy-name is the message name
- name is the first field <!-- name of what? -->
- bid is second field <!-- the price of the bid for a name? let's say a bit more here, what else can we say here? I wish our CLI reference had examples -->

The `starport scaffold message buy-name name bid` command creates and modifies several files and outputs the changes:

```bash
modify proto/nameservice/tx.proto
modify x/nameservice/client/cli/tx.go
create x/nameservice/client/cli/tx_buy_name.go
modify x/nameservice/handler.go
create x/nameservice/keeper/msg_server_buy_name.go
modify x/nameservice/types/codec.go
create x/nameservice/types/message_buy_name.go
```

- `proto/nameservice/tx.proto`
    - Adds `MsgBuyName` and `MsgBuyNameResponse` proto messages.
    - Registers `BuyName` rpc in the `Msg` service.

    You can these changes in each file. For example:
    ```go
    syntax = "proto3";
    package cosmonaut.nameservice.nameservice;

    // this line is used by starport scaffolding # proto/tx/import

    option go_package = "github.com/cosmonaut/nameservice/x/nameservice/types";

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

- `x/nameservice/handler.go`
    Registers the `MsgBuyName` message in the module message handler.

- `x/nameservice/keeper/msg_server_buy_name.go`
	Defines the `BuyName` keeper method. You can notice that the message follows the `Msg` interface. The message `struct` contains all the necessary information when buying a name: `Name`, `Bid`, and `Creator`. This struct was added automatically.

- `x/nameservice/client/cli/tx_buy_name.go`
  	Adds the CLI command to broadcast a transaction with a message. 

- `x/nameservice/types/codec.go`
	Registers the codecs.



### Add The MsgSetName Message

To create the `MsgSetName` for the nameservice module:


```bash
starport scaffold message set-name name value
```

where:

- set-name is the message name
- name is the first field <!-- name of what? -->
- value is the value that the name resolves to


This is the same `starport scaffold message` command, so the set of modified and created files are the same. 

### Add The MsgDelete Message

To create the `MsgDeleteName` for the nameservice module:


```bash
starport scaffold message delete-name name
```

where:

- delete-name is the message name
- name is message name to delete

## Results

Congratulations, you've defined messages that trigger state transitions. Now it's time to implement types and methods that operate on the state.
