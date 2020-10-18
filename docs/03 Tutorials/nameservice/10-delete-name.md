---
order: 11
---

# Delete Name

## MsgDeleteName


Now it is time to update the `Msg` for deleting names and add it to the `./x/nameservice/types/MsgDeleteName.go` file. 

<<< @/nameservice/nameservice/x/nameservice/types/MsgDeleteName.go


Finally, define the `DeleteName` `handler` function which performs the state transitions triggered by the message. Keep in mind that at this point the message has had its `ValidateBasic` function run so there has been some input verification. However, `ValidateBasic` cannot query application state. Validation logic that is dependent on network state (e.g. account balances) should be performed in the `handler` function.

<<< @/nameservice/nameservice/x/nameservice/handlerMsgDeleteName.go

Afterwards, we'll follow the same steps as earlier and add the `MsgDeleteName` handler to the module router in `./x/nameservice/handler.go`:

```go
// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		case MsgDeleteName:
			return handleMsgDeleteName(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type()))
		}
	}
}
```

### Now that you have your `Msgs` and `Handlers` defined it's time to learn about making the data from these transactions [available for querying.
