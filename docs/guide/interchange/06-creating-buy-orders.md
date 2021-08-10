---
order: 6
---

# Creating Buy Orders

In this chapter you will be implementing creation of buy orders. The logic is very similar to the previous chapter.

## Modify the Proto Definition

Add the buyer to the proto file definition

```proto
// proto/ibcdex/packet.proto
message BuyOrderPacketData {
  // ...
  string buyer = 5;
}
```

## IBC Message Handling in SendBuyOrder

* Check if the pair exists on the order book
* If the token is an IBC token, burn the tokens
* If the token is a native token, lock the tokens
* Save the voucher received on the target chain to later resolve a denom

```go
// x/ibcdex/keeper/msg_server_buy_order.go
import (
  //...
  "errors"
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
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
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
	packet.Buyer = msg.Sender
  // Transmit an IBC packet...
	return &types.MsgSendBuyOrderResponse{}, nil
}
```

## `OnRecv`

- Update the buy order book
- Distribute sold token to the buyer
- Send to chain A the sell order after the fill attempt

```go
// x/ibcdex/keeper/buy_order.go
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
		Price: data.Price,
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

### `FillBuyOrder`

`FillBuyOrder` try to fill the buy order with the order book and returns all the side effects.

```go
// x/ibcdex/types/sell_order_book.go
func (s *SellOrderBook) FillBuyOrder(order Order) (remainingBuyOrder Order, liquidated []Order, purchase int32, filled bool) {
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

#### `LiquidateFromBuyOrder`

`LiquidateFromBuyOrder` liquidates the first sell order of the book from the buy order. If no match is found, return false for match

```go
// x/ibcdex/types/sell_order_book.go
func (s *SellOrderBook) LiquidateFromBuyOrder(order Order) (remainingBuyOrder Order, liquidatedSellOrder Order, purchase int32, match bool, filled bool) {
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

## `OnAcknowledgement`

- Chain `Mars` will store the remaining sell order in the sell order book and will distribute sold `MCX` to the buyers and will distribute to the seller the price of the amount sold
- On error we mint back the burned tokens

```go
// x/ibcdex/keeper/buy_order.go
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

```go
// x/ibcdex/types/buy_order_book.go
func (b *BuyOrderBook) AppendOrder(creator string, amount int32, price int32) (int32, error) {
  return b.Book.appendOrder(creator, amount, price, Increasing)
}
```

## `OnTimeout`

If a timeout occurs, we mint back the native token.

```go
// x/ibcdex/keeper/buy_order.go
func (k Keeper) OnTimeoutBuyOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SellOrderPacketData) error {
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
