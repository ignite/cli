---
order: 2
---

# Messages

Messages are a great place to start when building a module because they define the actions that your application can make.

## Message Type

Messages trigger state transitions. Messages (`Msg`) are wrapped in transactions (`Tx`) that clients submit to the network. The Cosmos SDK wraps and unwraps `Msg` from `Tx`, which means, as an app developer, you only have to define messages. `Msgs` must satisfy the following interface:

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

The `Msg` type extends `proto.Message` and contains five methods. Three of them are legacy methods (`Type`, `Route` and `GetSignBytes`).

`ValidateBasic` is called early in the processing of the message in order to discard obviously invalid messages. `ValidateBasic` should only include checks that do not require access to the state, for example, that the `amount` of tokens is a positive value.

`GetSigners` Return the list of signers. The SDK will make sure that each message contained in a transaction is signed by all the signers listed in the list returned by this method.

## `Handlers`

`Handlers` define the action that needs to be taken (which stores need to get updated, how, and under what conditions) when a given `Msg` is received.

In this module you have three types of `Msgs` that users can send to interact with the application state: `SetName`, `BuyName` and `DeleteName`. They will each have an associated `Handler`.

## Scaffolding Messages

The `nameservice` module will have three messages:

* `MsgBuyName`
* `MsgSetName`
* `MsgDeleteName`

### `MsgBuyName`

Use the `starport scaffold message` command to scaffold a new Cosmos SDK message for your module. The command accepts message name as the first argument and a list of fields. By default, a message is scaffolded in a module with a name that matches the name of the project, in our case `nameservice` (this behaviour can be overwritten by using a flag).

```
starport scaffold message buy-name name bid
```

The command has created and modified several files.

* `proto/nameservice/tx.proto`: `MsgBuyName` and `MsgBuyNameResponse` proto messages are added nameservice a `BuyName` RPC is registered in the `Msg` service.
* `x/nameservice/types/message_buy_name.go`: methods are defined to satisfy `Msg` interface.
* `x/nameservice/handler.go`: `MsgBuyName` message is registered in the module message handler.
* `x/nameservice/keeper/msg_server_buy_name.go`: `BuyName` keeper method is defined
* `x/nameservice/client/cli/tx_buy_name.go`: CLI command added to brodcast a transaction with a message.nameservice `x/nameservice/client/cli/tx.go`: CLI command is registered.
* `x/nameservice/types/codec.go`: codecs are registered.

In `x/nameservice/types/message_buy_name.go` you can notice that the message follows the `Msg` interface. The message `struct` contains all the necessary information when buying a name: `Name`, `Bid`, and `Creator` (which was added automatically).

### Add The MsgSetName Message

Set name message needs to contain the following fields:

* Name
* Value - the value that the name resolves to

```
starport scaffold message set-name name value
```

As you're using the same `starport scaffold message` command the set of modified and created files are the same.

### Add The MsgDelete Message

```
starport scaffold message delete-name name
```

Delete name message needs only a `name` as an argument.
