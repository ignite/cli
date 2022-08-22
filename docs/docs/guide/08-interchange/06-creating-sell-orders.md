---
sidebar_position: 6
description: Implement logic to create sell orders.
---

# Create Sell Orders

In this chapter, you implement the logic for creating sell orders.

The packet proto file for a sell order is already generated. Add the seller information:

```protobuf
// proto/dex/packet.proto

message SellOrderPacketData {
  // ...
  string seller = 5;
}
```

Now, use Ignite CLI to build the proto files for the `send-sell-order` command. You used this command in a previous chapter. 

```bash
ignite generate proto-go --yes
```

## Message Handling in SendSellOrder

Sell orders are created using the `send-sell-order` command. This command creates a transaction with a `SendSellOrder` message that triggers the `SendSellOrder` keeper method.

The `SendSellOrder` command:

* Checks that an order book for a specified denom pair exists.
* Safely burns or locks token.
  * If the token is an IBC token, burn the token.
  * If the token is a native token, lock the token.
* Saves the voucher that is received on the target chain to later resolve a denom.
* Transmits an IBC packet to the target chain.

```go
// x/dex/keeper/msg_server_sell_order.go

package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"

	"interchange/x/dex/types"
)

func (k msgServer) SendSellOrder(goCtx context.Context, msg *types.MsgSendSellOrder) (*types.MsgSendSellOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// If an order book doesn't exist, throw an error
	pairIndex := types.OrderBookIndex(msg.Port, msg.ChannelID, msg.AmountDenom, msg.PriceDenom)
	_, found := k.GetSellOrderBook(ctx, pairIndex)
	if !found {
		return &types.MsgSendSellOrderResponse{}, errors.New("the pair doesn't exist")
	}

	// Get sender's address
	sender, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return &types.MsgSendSellOrderResponse{}, err
	}

	// Use SafeBurn to ensure no new native tokens are minted
	if err := k.SafeBurn(ctx, msg.Port, msg.ChannelID, sender, msg.AmountDenom, msg.Amount); err != nil {
		return &types.MsgSendSellOrderResponse{}, err
	}

	// Save the voucher received on the other chain, to have the ability to resolve it into the original denom
	k.SaveVoucherDenom(ctx, msg.Port, msg.ChannelID, msg.AmountDenom)

	var packet types.SellOrderPacketData
	packet.Seller = msg.Creator
	packet.AmountDenom = msg.AmountDenom
	packet.Amount = msg.Amount
	packet.PriceDenom = msg.PriceDenom
	packet.Price = msg.Price

	// Transmit the packet
	err = k.TransmitSellOrderPacket(ctx, packet, msg.Port, msg.ChannelID, clienttypes.ZeroHeight(), msg.TimeoutTimestamp)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendSellOrderResponse{}, nil
}
```

## On Receiving a Sell Order

When a "sell order" packet is received on the target chain, you want the module to:

* Update the sell order book
* Distribute sold token to the buyer
* Send the sell order to chain A after the fill attempt

```go
// x/dex/keeper/sell_order.go

func (k Keeper) OnRecvSellOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SellOrderPacketData) (packetAck types.SellOrderPacketAck, err error) {
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	pairIndex := types.OrderBookIndex(packet.SourcePort, packet.SourceChannel, data.AmountDenom, data.PriceDenom)
	book, found := k.GetBuyOrderBook(ctx, pairIndex)
	if !found {
		return packetAck, errors.New("the pair doesn't exist")
	}

	// Fill sell order
	remaining, liquidated, gain, _ := book.FillSellOrder(types.Order{
		Amount: data.Amount,
		Price:  data.Price,
	})

	// Return remaining amount and gains
	packetAck.RemainingAmount = remaining.Amount
	packetAck.Gain = gain

	// Before distributing sales, we resolve the denom
	// First we check if the denom received comes from this chain originally
	finalAmountDenom, saved := k.OriginalDenom(ctx, packet.DestinationPort, packet.DestinationChannel, data.AmountDenom)
	if !saved {
		// If it was not from this chain we use voucher as denom
		finalAmountDenom = VoucherDenom(packet.SourcePort, packet.SourceChannel, data.AmountDenom)
	}

	// Dispatch liquidated buy orders
	for _, liquidation := range liquidated {
		liquidation := liquidation
		addr, err := sdk.AccAddressFromBech32(liquidation.Creator)
		if err != nil {
			return packetAck, err
		}

		if err := k.SafeMint(ctx, packet.DestinationPort, packet.DestinationChannel, addr, finalAmountDenom, liquidation.Amount); err != nil {
			return packetAck, err
		}
	}

	// Save the new order book
	k.SetBuyOrderBook(ctx, book)

	return packetAck, nil
}
```

### Implement a FillBuyOrder Function

The `FillBuyOrder` function tries to fill the sell order with the order book and returns all the side effects:

```go
// x/dex/types/sell_order_book.go

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

The `LiquidateFromBuyOrder` function liquidates the first buy order of the book from the sell order. If no match is found, return false for match:

```go
// x/dex/types/sell_order_book.go

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

### Implement the OnAcknowledgement Function for Sell Order Packets

After an IBC packet is processed on the target chain, an acknowledgement is returned to the source chain and processed by the `OnAcknowledgementSellOrderPacket` function. 

The dex module on the source chain:

- Stores the remaining sell order in the sell order book.
- Distributes sold tokens to the buyers.
- Distributes the price of the amount sold to the seller. 
- On error, mints the burned tokens.

```go
// x/dex/keeper/sell_order.go

func (k Keeper) OnAcknowledgementSellOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SellOrderPacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		// In case of error we mint back the native token
		receiver, err := sdk.AccAddressFromBech32(data.Seller)
		if err != nil {
			return err
		}

		if err := k.SafeMint(ctx, packet.SourcePort, packet.SourceChannel, receiver, data.AmountDenom, data.Amount); err != nil {
			return err
		}

		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.SellOrderPacketAck
		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		// Get the sell order book
		pairIndex := types.OrderBookIndex(packet.SourcePort, packet.SourceChannel, data.AmountDenom, data.PriceDenom)
		book, found := k.GetSellOrderBook(ctx, pairIndex)
		if !found {
			panic("sell order book must exist")
		}

		// Append the remaining amount of the order
		if packetAck.RemainingAmount > 0 {
			_, err := book.AppendOrder(data.Seller, packetAck.RemainingAmount, data.Price)
			if err != nil {
				return err
			}

			// Save the new order book
			k.SetSellOrderBook(ctx, book)
		}

		// Mint the gains
		if packetAck.Gain > 0 {
			receiver, err := sdk.AccAddressFromBech32(data.Seller)
			if err != nil {
				return err
			}

			finalPriceDenom, saved := k.OriginalDenom(ctx, packet.SourcePort, packet.SourceChannel, data.PriceDenom)
			if !saved {
				// If it was not from this chain we use voucher as denom
				finalPriceDenom = VoucherDenom(packet.DestinationPort, packet.DestinationChannel, data.PriceDenom)
			}

			if err := k.SafeMint(ctx, packet.SourcePort, packet.SourceChannel, receiver, finalPriceDenom, packetAck.Gain); err != nil {
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

```go
// x/dex/types/sell_order_book.go

func (s *SellOrderBook) AppendOrder(creator string, amount int32, price int32) (int32, error) {
	return s.Book.appendOrder(creator, amount, price, Decreasing)
}
```

### Add the OnTimeout of a Sell Order Packet Function

If a timeout occurs, mint back the native token:

```go
// x/dex/keeper/sell_order.go

func (k Keeper) OnTimeoutSellOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SellOrderPacketData) error {
	// In case of error we mint back the native token
	receiver, err := sdk.AccAddressFromBech32(data.Seller)
	if err != nil {
		return err
	}

	if err := k.SafeMint(ctx, packet.SourcePort, packet.SourceChannel, receiver, data.AmountDenom, data.Amount); err != nil {
		return err
	}

	return nil
}
```

## Summary

Great, you have completed the sell order logic.

It is a good time to make another git commit again to save the state of your work:

```bash
git add .
git commit -m "Add Sell Orders"
```
