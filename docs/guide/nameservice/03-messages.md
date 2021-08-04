# Messages

Messages are a great place to start when building a module because they define the actions that your application can make.

## `Msgs`

`Msgs` trigger state transitions. `Msgs` are wrapped in [`Txs`](https://github.com/cosmos/cosmos-sdk/blob/master/types/tx_msg.go#L34-L41) that clients submit to the network. The Cosmos SDK wraps and unwraps `Msgs` from `Txs`, which means, as an app developer, you only have to define `Msgs`. `Msgs` must satisfy the following interface:

```go
// Transactions messages must fulfill the Msg
type Msg interface {
	// Return the message type.
	// Must be alphanumeric or empty.
	Type() string

	// Returns a human-readable string for the message, intended for utilization
	// within tags
	Route() string

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() Error

	// Get the canonical byte representation of the Msg.
	GetSignBytes() []byte

	// Signers returns the addrs of signers that must sign.
	// CONTRACT: All signatures must be present to be valid.
	// CONTRACT: Returns addrs in some deterministic order.
	GetSigners() []AccAddress
}
```

`GetSigners` defines whose signature is required on a `Tx` in order for it to be valid.

`GetSignBytes` defines how the `Msg` gets encoded for signing. In most cases this means marshal to sorted JSON.

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

In `x/nameservice/types/message_buy_name.go` you can notice that the message follows the `sdk.Msg` interface. The message `struct` contains all the necessary information when buying a name: `Name`, `Bid`, and `Creator` (which was added automatically).


### `MsgSetName`

Set name message needs to contain the following fields:

* Name
* Value - the value that the name resolves to

```
starport scaffold message set-name name value
```

As you're using the same `starport scaffold message` command the set of modified and created files are the same.

### `MsgDelete`

```
starport scaffold message delete-name name
```

Delete name message needs only a `name` as an argument.