---
sidebar_position: 7
description: Implement the buy order logic.
---

# Creating Buy Orders

In this chapter, you implement the creation of buy orders. The logic is very similar to the sell order logic you
implemented in the previous chapter.

## Modify the Proto Definition

Add the buyer to the proto file definition:

```protobuf
// proto/interchange/dex/packet.proto

message BuyOrderPacketData {
  // ...
  string buyer = 5;
}
```

Now, use Ignite CLI to build the proto files for the `send-buy-order` command. You used this command in previous
chapters.

```bash
ignite generate proto-go --yes
```

## IBC Message Handling in SendBuyOrder

* Check if the pair exists on the order book
* If the token is an IBC token, burn the tokens
* If the token is a native token, lock the tokens
* Save the voucher received on the target chain to later resolve a denom

```go
// x/dex/keeper/msg_server_buy_order.go

package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"interchange/x/dex/types"
)

func (k msgServer) SendBuyOrder(goCtx context.Context, msg *types.MsgSendBuyOrder) (*types.MsgSendBuyOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Cannot send a order if the pair doesn't exist
	pairIndex := types.OrderBookIndex(msg.Port, msg.ChannelID, msg.AmountDenom, msg.PriceDenom)
	_, found := k.GetBuyOrderBook(ctx, pairIndex)
	if !found {
		return &types.MsgSendBuyOrderResponse{}, errors.New("the pair doesn't exist")
	}

	// Lock the token to send
	sender, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return &types.MsgSendBuyOrderResponse{}, err
	}

	// Use SafeBurn to ensure no new native tokens are minted
	if err := k.SafeBurn(ctx, msg.Port, msg.ChannelID, sender, msg.PriceDenom, msg.Amount*msg.Price); err != nil {
		return &types.MsgSendBuyOrderResponse{}, err
	}

	// Save the voucher received on the other chain, to have the ability to resolve it into the original denom
	k.SaveVoucherDenom(ctx, msg.Port, msg.ChannelID, msg.PriceDenom)

	// Construct the packet
	var packet types.BuyOrderPacketData

	packet.Buyer = msg.Creator
	packet.AmountDenom = msg.AmountDenom
	packet.Amount = msg.Amount
	packet.PriceDenom = msg.PriceDenom
	packet.Price = msg.Price

	// Transmit the packet
	err = k.TransmitBuyOrderPacket(
		ctx,
		packet,
		msg.Port,
		msg.ChannelID,
		clienttypes.ZeroHeight(),
		msg.TimeoutTimestamp,
	)
	if err != nil {
		return nil, err
	}

	// Transmit an IBC packet...
	return &types.MsgSendBuyOrderResponse{}, nil
}
```

## On Receiving a Buy Order

* Update the buy order book
* Distribute sold token to the buyer
* Send the sell order to chain A after the fill attempt

```go
// x/dex/keeper/buy_order.go

package keeper

// ...

func (k Keeper) OnRecvBuyOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.BuyOrderPacketData) (packetAck types.BuyOrderPacketAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	// Check if the sell order book exists
	pairIndex := types.OrderBookIndex(packet.SourcePort, packet.SourceChannel, data.AmountDenom, data.PriceDenom)
	book, found := k.GetSellOrderBook(ctx, pairIndex)
	if !found {
		return packetAck, errors.New("the pair doesn't exist")
	}

	// Fill buy order
	remaining, liquidated, purchase, _ := book.FillBuyOrder(types.Order{
		Amount: data.Amount,
		Price:  data.Price,
	})

	// Return remaining amount and gains
	packetAck.RemainingAmount = remaining.Amount
	packetAck.Purchase = purchase

	// Before distributing gains, we resolve the denom
	// First we check if the denom received comes from this chain originally
	finalPriceDenom, saved := k.OriginalDenom(ctx, packet.DestinationPort, packet.DestinationChannel, data.PriceDenom)
	if !saved {
		// If it was not from this chain we use voucher as denom
		finalPriceDenom = VoucherDenom(packet.SourcePort, packet.SourceChannel, data.PriceDenom)
	}

	// Dispatch liquidated buy order
	for _, liquidation := range liquidated {
		liquidation := liquidation
		addr, err := sdk.AccAddressFromBech32(liquidation.Creator)
		if err != nil {
			return packetAck, err
		}

		if err := k.SafeMint(
			ctx,
			packet.DestinationPort,
			packet.DestinationChannel,
			addr,
			finalPriceDenom,
			liquidation.Amount*liquidation.Price,
		); err != nil {
			return packetAck, err
		}
	}

	// Save the new order book
	k.SetSellOrderBook(ctx, book)

	return packetAck, nil
}
```

### Implement a FillBuyOrder Function

The `FillBuyOrder` function tries to fill the sell order with the order book and returns all the side effects:

```go
// x/dex/types/sell_order_book.go

package types

// ...

func (s *SellOrderBook) FillBuyOrder(order Order) (
	remainingBuyOrder Order,
	liquidated []Order,
	purchase int32,
	filled bool,
) {
	var liquidatedList []Order
	totalPurchase := int32(0)
	remainingBuyOrder = order

	// Liquidate as long as there is match
	for {
		var match bool
		var liquidation Order
		remainingBuyOrder, liquidation, purchase, match, filled = s.LiquidateFromBuyOrder(
			remainingBuyOrder,
		)
		if !match {
			break
		}

		// Update gains
		totalPurchase += purchase

		// Update liquidated
		liquidatedList = append(liquidatedList, liquidation)

		if filled {
			break
		}
	}

	return remainingBuyOrder, liquidatedList, totalPurchase, filled
}
```

### Implement a LiquidateFromBuyOrder Function

The `LiquidateFromBuyOrder` function liquidates the first buy order of the book from the sell order. If no match is
found, return false for match:

```go
// x/dex/types/sell_order_book.go

package types

// ...

func (s *SellOrderBook) LiquidateFromBuyOrder(order Order) (
	remainingBuyOrder Order,
	liquidatedSellOrder Order,
	purchase int32,
	match bool,
	filled bool,
) {
	remainingBuyOrder = order

	// No match if no order
	orderCount := len(s.Book.Orders)
	if orderCount == 0 {
		return order, liquidatedSellOrder, purchase, false, false
	}

	// Check if match
	lowestAsk := s.Book.Orders[orderCount-1]
	if order.Price < lowestAsk.Price {
		return order, liquidatedSellOrder, purchase, false, false
	}

	liquidatedSellOrder = *lowestAsk

	// Check if buy order can be entirely filled
	if lowestAsk.Amount >= order.Amount {
		remainingBuyOrder.Amount = 0
		liquidatedSellOrder.Amount = order.Amount
		purchase = order.Amount

		// Remove lowest ask if it has been entirely liquidated
		lowestAsk.Amount -= order.Amount
		if lowestAsk.Amount == 0 {
			s.Book.Orders = s.Book.Orders[:orderCount-1]
		} else {
			s.Book.Orders[orderCount-1] = lowestAsk
		}

		return remainingBuyOrder, liquidatedSellOrder, purchase, true, true
	}

	// Not entirely filled
	purchase = lowestAsk.Amount
	s.Book.Orders = s.Book.Orders[:orderCount-1]
	remainingBuyOrder.Amount -= lowestAsk.Amount

	return remainingBuyOrder, liquidatedSellOrder, purchase, true, false
}
```

## Receiving a Buy Order Acknowledgment

After a buy order acknowledgement is received, chain `Mars`:

* Stores the remaining sell order in the sell order book.
* Distributes sold `marscoin` to the buyers.
* Distributes to the seller the price of the amount sold.
* On error, mints back the burned tokens.

```go
// x/dex/keeper/buy_order.go

package keeper

// ...

func (k Keeper) OnAcknowledgementBuyOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.BuyOrderPacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		// In case of error we mint back the native token
		receiver, err := sdk.AccAddressFromBech32(data.Buyer)
		if err != nil {
			return err
		}

		if err := k.SafeMint(
			ctx,
			packet.SourcePort,
			packet.SourceChannel,
			receiver,
			data.PriceDenom,
			data.Amount*data.Price,
		); err != nil {
			return err
		}

		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.BuyOrderPacketAck

		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		// Get the sell order book
		pairIndex := types.OrderBookIndex(packet.SourcePort, packet.SourceChannel, data.AmountDenom, data.PriceDenom)
		book, found := k.GetBuyOrderBook(ctx, pairIndex)
		if !found {
			panic("buy order book must exist")
		}

		// Append the remaining amount of the order
		if packetAck.RemainingAmount > 0 {
			_, err := book.AppendOrder(
				data.Buyer,
				packetAck.RemainingAmount,
				data.Price,
			)
			if err != nil {
				return err
			}

			// Save the new order book
			k.SetBuyOrderBook(ctx, book)
		}

		// Mint the purchase
		if packetAck.Purchase > 0 {
			receiver, err := sdk.AccAddressFromBech32(data.Buyer)
			if err != nil {
				return err
			}

			finalAmountDenom, saved := k.OriginalDenom(ctx, packet.SourcePort, packet.SourceChannel, data.AmountDenom)
			if !saved {
				// If it was not from this chain we use voucher as denom
				finalAmountDenom = VoucherDenom(packet.DestinationPort, packet.DestinationChannel, data.AmountDenom)
			}

			if err := k.SafeMint(
				ctx,
				packet.SourcePort,
				packet.SourceChannel,
				receiver,
				finalAmountDenom,
				packetAck.Purchase,
			); err != nil {
				return err
			}
		}

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}
```

`AppendOrder` appends an order in the buy order book.
Add the following function to the `x/dex/types/buy_order_book.go` file in the `types` directory.

```go
// x/dex/types/buy_order_book.go

package types

// ...

func (b *BuyOrderBook) AppendOrder(creator string, amount int32, price int32) (int32, error) {
	return b.Book.appendOrder(creator, amount, price, Increasing)
}
```

## OnTimeout of a Buy Order Packet

If a timeout occurs, mint back the native token:

```go
// x/dex/keeper/buy_order.go

package keeper

// ...

func (k Keeper) OnTimeoutBuyOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.BuyOrderPacketData) error {
	// In case of error we mint back the native token
	receiver, err := sdk.AccAddressFromBech32(data.Buyer)
	if err != nil {
		return err
	}

	if err := k.SafeMint(
		ctx,
		packet.SourcePort,
		packet.SourceChannel,
		receiver,
		data.PriceDenom,
		data.Amount*data.Price,
	); err != nil {
		return err
	}

	return nil
}
```

## Summary

Congratulations, you implemented the buy order logic.

Again, it's a good time to save your current state to your local GitHub repository:

```bash
git add .
git commit -m "Add Buy Orders"
```
