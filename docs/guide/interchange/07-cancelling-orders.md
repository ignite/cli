---
order: 7
---

# Cancelling Orders

You have implemented order books, buy and sell orders.  In this chapter you will enable cancelling buy and sell orders. 

## Cancel a Sell Order

To cancel a sell order, you have to get the ID of the specific sell order. Then you can use the function `RemoveOrderFromID` to remove the specific order from the order book and update the keeper accordingly.

```go
// x/ibcdex/keeper/msg_server_cancel_sell_order.go
import "errors"

func (k msgServer) CancelSellOrder(goCtx context.Context, msg *types.MsgCancelSellOrder) (*types.MsgCancelSellOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Retrieve the book
	pairIndex := types.OrderBookIndex(msg.Port, msg.Channel, msg.AmountDenom, msg.PriceDenom)
	s, found := k.GetSellOrderBook(ctx, pairIndex)
	if !found {
		return &types.MsgCancelSellOrderResponse{}, errors.New("the pair doesn't exist")
	}
	// Check order creator
	order, err := s.Book.GetOrderFromID(msg.OrderID)
	if err != nil {
		return &types.MsgCancelSellOrderResponse{}, err
	}
	if order.Creator != msg.Creator {
		return &types.MsgCancelSellOrderResponse{}, errors.New("canceller must be creator")
	}
	// Remove order
	if err := s.Book.RemoveOrderFromID(msg.OrderID); err != nil {
		return &types.MsgCancelSellOrderResponse{}, err
	}
	k.SetSellOrderBook(ctx, s)
	// Refund seller with remaining amount
	seller, err := sdk.AccAddressFromBech32(order.Creator)
	if err != nil {
		return &types.MsgCancelSellOrderResponse{}, err
	}
	if err := k.SafeMint(ctx, msg.Port, msg.Channel, seller, msg.AmountDenom, order.Amount); err != nil {
		return &types.MsgCancelSellOrderResponse{}, err
	}
	return &types.MsgCancelSellOrderResponse{}, nil
}
```

`GetOrderFromID` gets an order from the book from its ID.

```go
// x/ibcdex/types/order_book.go
func (book OrderBook) GetOrderFromID(id int32) (Order, error) {
	for _, order := range book.Orders {
		if order.Id == id {
			return *order, nil
		}
	}
	return Order{}, ErrOrderNotFound
}
```

`RemoveOrderFromID` removes an order from the book and keep it ordered.

```go
// x/ibcdex/types/order_book.go
func (book *OrderBook) RemoveOrderFromID(id int32) error {
	for i, order := range book.Orders {
		if order.Id == id {
			book.Orders = append(book.Orders[:i], book.Orders[i+1:]...)
			return nil
		}
	}
	return ErrOrderNotFound
}
```

## Cancel a Buy Order

To cancel a buy order, you have to get the ID of the specific buy order. Then you can use the function `RemoveOrderFromID` to remove the specific order from the order book and update the keeper accordingly.

```go
// x/ibcdex/keeper/msg_server_cancel_buy_order.go
import "errors"

func (k msgServer) CancelBuyOrder(goCtx context.Context, msg *types.MsgCancelBuyOrder) (*types.MsgCancelBuyOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Retrieve the book
	pairIndex := types.OrderBookIndex(msg.Port, msg.Channel, msg.AmountDenom, msg.PriceDenom)
	b, found := k.GetBuyOrderBook(ctx, pairIndex)
	if !found {
		return &types.MsgCancelBuyOrderResponse{}, errors.New("the pair doesn't exist")
	}
	// Check order creator
	order, err := b.Book.GetOrderFromID(msg.OrderID)
	if err != nil {
		return &types.MsgCancelBuyOrderResponse{}, err
	}
	if order.Creator != msg.Creator {
		return &types.MsgCancelBuyOrderResponse{}, errors.New("canceller must be creator")
	}
	// Remove order
	if err := b.Book.RemoveOrderFromID(msg.OrderID); err != nil {
		return &types.MsgCancelBuyOrderResponse{}, err
	}
	k.SetBuyOrderBook(ctx, b)
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
