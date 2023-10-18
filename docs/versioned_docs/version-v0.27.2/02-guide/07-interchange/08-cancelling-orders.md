---
sidebar_position: 8
description: Enable cancelling of buy and sell orders.
---

# Cancelling Orders

You have implemented order books, buy and sell orders. In this chapter, you enable cancelling of buy and sell orders.

## Cancel a Sell Order

To cancel a sell order, you have to get the ID of the specific sell order. Then you can use the function
`RemoveOrderFromID` to remove the specific order from the order book and update the keeper accordingly.

Move to the keeper directory and edit the `x/dex/keeper/msg_server_cancel_sell_order.go` file:

```go
// x/dex/keeper/msg_server_cancel_sell_order.go

package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"interchange/x/dex/types"
)

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

### Implement the GetOrderFromID Function

The `GetOrderFromID` function gets an order from the book from its ID.

Add this function to the `x/dex/types/order_book.go` function in the `types` directory:

```go
// x/dex/types/order_book.go

func (book OrderBook) GetOrderFromID(id int32) (Order, error) {
	for _, order := range book.Orders {
		if order.Id == id {
			return *order, nil
		}
	}

	return Order{}, ErrOrderNotFound
}
```

### Implement the RemoveOrderFromID Function

The `RemoveOrderFromID` function removes an order from the book and keeps it ordered:

```go
// x/dex/types/order_book.go

package types

// ...

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

To cancel a buy order, you have to get the ID of the specific buy order. Then you can use the function
`RemoveOrderFromID` to remove the specific order from the order book and update the keeper accordingly:

```go
// x/dex/keeper/msg_server_cancel_buy_order.go

package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"interchange/x/dex/types"
)

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

## Summary

You have completed implementing the functions that are required for the `dex` module. In this chapter, you have
implemented the design for cancelling specific buy or sell orders.

To test if your Ignite CLI blockchain builds correctly, use the `chain build` command:

```bash
ignite chain build
```

Again, it is a good time (a great time!) to add your state to the local GitHub repository:

```bash
git add .
git commit -m "Add Cancelling Orders"
```

Finally, it's now time to write test files.
