---
order: 10
---

# Buy Name

## MsgBuyName

Now it is time to define the `Msg` for buying names and add it to the `./x/nameservice/types/MsgBuyName.go` file. This code is very similar to `SetName`. We can replace the file `MsgCreateWhois.go`, as these two files are similar in nature, and we won't be using `MsgCreateWhois`.

```
mv x/nameservice/types/MsgCreateWhois.go x/nameservice/types/MsgBuyName.go
```

<<< @/nameservice/nameservice/x/nameservice/types/MsgBuyName.go

Next, in the `./x/nameservice/handler.go` file, add the `MsgBuyName` handler to the module router:

```go
// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type()))
		}
	}
}
```

Finally, update the `BuyName` `handler` function which performs the state transitions triggered by the message. Keep in mind that at this point the message has had its `ValidateBasic` function run so there has been some input verification. However, `ValidateBasic` cannot query application state. Validation logic that is dependent on network state (e.g. account balances) should be performed in the `handler` function.

<<< @/nameservice/nameservice/x/nameservice/handlerMsgBuyName.go

First check to make sure that the bid is higher than the current price. Then, check to see whether the name already has an owner. If it does, the former owner will receive the money from the `Buyer`.

If there is no owner, your `nameservice` module "burns" (i.e. sends to an unrecoverable address) the coins from the `Buyer`.

If either `SubtractCoins` or `SendCoins` returns a non-nil error, the handler throws an error, reverting the state transition. Otherwise, using the getters and setters defined on the `Keeper` earlier, the handler sets the buyer to the new owner and sets the new price to be the current bid.

> _*NOTE*_: This handler uses functions from the `coinKeeper` to perform currency operations. If your application is performing currency operations you may want to take a look at the [godocs for this module](https://godoc.org/github.com/cosmos/cosmos-sdk/x/bank#BaseKeeper) to see what functions it exposes.

### Great, now owners can `BuyName`s! But what if they don't want the name any longer? Your module needs a way for users to delete names! Let us define define the `DeleteName` message.
