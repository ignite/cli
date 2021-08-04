---
order: 9
---

# Cancel orders

You have implemented order books, buy and sell orders. 
In this chapter you will enable cancelling buy and sell orders. 

The function `RemoveOrderFromID` will be used to remove the buy or sell order from the order book.

## Cancel the Sell Order

To cancel a sell order, you have to get the ID of the specific sell order.
Then you can use the function `RemoveOrderFromID` to remove the specific order from the order book and update the keeper accordingly.

```go
// x/ibcdex/keeper/msg_server_cancelSellOrder.go
import "errors"

func (k msgServer) CancelSellOrder(goCtx context.Context, msg *types.MsgCancelSellOrder) (*types.MsgCancelSellOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Retrieve the book
	pairIndex := types.OrderBookIndex(msg.Port, msg.Channel, msg.AmountDenom, msg.PriceDenom)
	book, found := k.GetSellOrderBook(ctx, pairIndex)
	if !found {
		return &types.MsgCancelSellOrderResponse{}, errors.New("the pair doesn't exist")
	}

	// Check order creator
	order, err := book.GetOrderFromID(msg.OrderID)
	if err != nil {
		return &types.MsgCancelSellOrderResponse{}, err
	}
	if order.Creator != msg.Creator {
		return &types.MsgCancelSellOrderResponse{}, errors.New("canceller must be creator")
	}

	// Remove order
	if err := book.RemoveOrderFromID(msg.OrderID); err != nil {
		return &types.MsgCancelSellOrderResponse{}, err
	}
	k.SetSellOrderBook(ctx, book)

    // Refund seller with remaining amount
	seller, err := sdk.AccAddressFromBech32(order.Creator)
	if err != nil {
		return &types.MsgCancelSellOrderResponse{}, err
	}
	if err := k.SafeMint(
		ctx,
		msg.Port,
		msg.Channel,
		seller,
		msg.AmountDenom,
		order.Amount,
	); err != nil {
		return &types.MsgCancelSellOrderResponse{}, err
	}

	return &types.MsgCancelSellOrderResponse{}, nil
}
```

## Cancel the Buy Order

To cancel a buy order, you have to get the ID of the specific buy order.
Then you can use the function `RemoveOrderFromID` to remove the specific order from the order book and update the keeper accordingly.

```go
// x/ibcdex/keeper/msg_server_cancelBuyOrder.go
import "errors"

func (k msgServer) CancelBuyOrder(goCtx context.Context, msg *types.MsgCancelBuyOrder) (*types.MsgCancelBuyOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Retrieve the book
	pairIndex := types.OrderBookIndex(msg.Port, msg.Channel, msg.AmountDenom, msg.PriceDenom)
	book, found := k.GetBuyOrderBook(ctx, pairIndex)
	if !found {
		return &types.MsgCancelBuyOrderResponse{}, errors.New("the pair doesn't exist")
	}

	// Check order creator
	order, err := book.GetOrderFromID(msg.OrderID)
	if err != nil {
		return &types.MsgCancelBuyOrderResponse{}, err
	}
	if order.Creator != msg.Creator {
		return &types.MsgCancelBuyOrderResponse{}, errors.New("canceller must be creator")
	}

	// Remove order
	if err := book.RemoveOrderFromID(msg.OrderID); err != nil {
		return &types.MsgCancelBuyOrderResponse{}, err
	}
	k.SetBuyOrderBook(ctx, book)

  // Refund buyer with remaining price amount
	buyer, err := sdk.AccAddressFromBech32(order.Creator)
	if err != nil {
		return &types.MsgCancelBuyOrderResponse{}, err
	}
	if err := k.SafeMint(
		ctx,
		msg.Port,
		msg.Channel,
		buyer,
		msg.PriceDenom,
		order.Amount*order.Price,
	); err != nil {
		return &types.MsgCancelBuyOrderResponse{}, err
	}

	return &types.MsgCancelBuyOrderResponse{}, nil
}
```

That finishes all necessary functions needed for the `ibcdex` module. In this chapter you have implemented the design for cancelling specific buy or sell orders.
In the next chapter, you will be able to use and interact with your `ibcdex` module. You will be using the command line to create order books, buy and sell orders on the `Mars` and `Venus` blockchain.