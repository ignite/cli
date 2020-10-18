---
order: 7
---

# Msgs and Handlers

Now that you have the `Keeper` setup, it is time to update the `Msgs` and `Handlers` so that users can buy names and set values for them.

## `Msgs`

`Msgs` trigger state transitions. `Msgs` are wrapped in [`Txs`](https://github.com/cosmos/cosmos-sdk/blob/master/types/tx_msg.go#L34-L41) that clients submit to the network. The Cosmos SDK wraps and unwraps `Msgs` from `Txs`, which means, as an app developer, you only have to define `Msgs`. `Msgs` must satisfy the following interface (we'll implement all of these in the next section):

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

## `Handlers`

`Handlers` define the action that needs to be taken (which stores need to get updated, how, and under what conditions) when a given `Msg` is received.

In this module you have three types of `Msgs` that users can send to interact with the application state: [`SetName`](set-name.md), [`BuyName`](./buy-name.md) and [`DeleteName`](./delete-name.md). They will each have an associated `Handler`.

We can see that a few files have already been scaffolded by the `type` command, and we can modify these files to fit our needs for messages and handlers.

Now that you have a better understanding of `Msgs` and `Handlers`, you can start building your first message: `SetName`
