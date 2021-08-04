---
order: 7
---

# Create the Sell Order IBC packet

In this chapter you want to modify the IBC logic to create sell orders on the IBC exchange.
A sell order must be submitted to an existing order book.
When you are dealing with a native token, these tokens will get locked until the IBC packets get reversed.
When you are dealing with an IBC token, these will get burned and you receive back the native token.

## Modify the Proto Definition

The packet proto file for a sell order is already generated. Add the seller information.

```proto
// proto/ibcdex/packet.proto
message SellOrderPacketData {
  string amountDenom = 1;
  int32 amount = 2;
  string priceDenom = 3;
  int32 price = 4;
  string seller = 5;  // <--
}
```

## About the IBC Packet

The IBC packet has four different stages you need to consider: 
1. Before transmitting the packet
2. On Receipt of a packet
3. On Acknowledgment of a packet
4. On Timeout of a packet

### Pre-transmit

Before a sell order will be submitted, make sure it contains the following logic:

- Check if the pair exists on the order book
- If the token is an IBC token, burn the tokens
- If the token is a native token, lock the tokens
- Save the voucher received on the target chain to later resolve a denom

```go
// x/ibcdex/keeper/msg_server_sellOrder.go
import "errors"

func (k msgServer) SendSellOrder(goCtx context.Context, msg *types.MsgSendSourceSellOrder) (*types.MsgSendSourceSellOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Cannot send a order if the order book pair doesn't exist
	pairIndex := types.OrderBookIndex(msg.Port, msg.ChannelID, msg.AmountDenom, msg.PriceDenom)
	_, found := k.GetSellOrderBook(ctx, pairIndex)
	if !found {
		return &types.MsgSendSellOrderResponse{}, errors.New("the pair doesn't exist")
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return &types.MsgSendSellOrderResponse{}, err
	}

	// Use SafeBurn to ensure no new native tokens are minted
	if err := k.SafeBurn(
		ctx,
		msg.Port,
		msg.ChannelID,
		sender,
		msg.AmountDenom,
		msg.Amount,
	); err != nil {
		return &types.MsgSendSellOrderResponse{}, err
	}

	// Save the voucher received on the other chain, to have the ability to resolve it into the original denom
	k.SaveVoucherDenom(ctx, msg.Port, msg.ChannelID, msg.AmountDenom)

	// Construct the packet
	var packet types.SellOrderPacketData

	packet.Seller = msg.Sender  // <- Manually specify the seller here
	packet.AmountDenom = msg.AmountDenom
	packet.Amount = msg.Amount
	packet.PriceDenom = msg.PriceDenom
	packet.Price = msg.Price

	// Transmit the packet
	err = k.TransmitSellOrderPacket(
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

	return &types.MsgSendSellOrderResponse{}, nil
}
```

## Create the OnRecv Function

- Update the buy order book
- Distribute sold token to the buyer
- Send to chain A the sell order after the fill attempt

```go
// x/ibcdex/keeper/sellOrder.go
// OnRecvSellOrderPacket processes packet reception
func (k Keeper) OnRecvSellOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SellOrderPacketData) (packetAck types.SellOrderPacketAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	// Check if the buy order book exists
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

		if err := k.SafeMint(
			ctx,
			packet.DestinationPort,
			packet.DestinationChannel,
			addr,
			finalAmountDenom,
			liquidation.Amount,
		); err != nil {
			return packetAck, err
		}
	}

	// Save the new order book
	k.SetBuyOrderBook(ctx, book)

	return packetAck, nil
}
```

## Create the OnAcknowledgement Function

- Chain `Mars` will store the remaining sell order in the sell order book and will distribute sold `MCX` to the buyers and will distribute to the seller the price of the amount sold
- On error we mint back the burned tokens

```go
// x/ibcdex/keeper/sellOrder.go
func (k Keeper) OnAcknowledgementSellOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SellOrderPacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		// In case of error we mint back the native token
		receiver, err := sdk.AccAddressFromBech32(data.Seller)
		if err != nil {
			return err
		}

		if err := k.SafeMint(
			ctx,
			packet.SourcePort,
			packet.SourceChannel,
			receiver,
			data.AmountDenom,
			data.Amount,
		); err != nil {
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
			_, err := book.AppendOrder(
				data.Seller,
				packetAck.RemainingAmount,
				data.Price,
			)
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
			if err := k.SafeMint(
				ctx,
				packet.SourcePort,
				packet.SourceChannel,
				receiver,
				finalPriceDenom,
				packetAck.Gain,
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

## Create the OnTimeout Function

If a timeout occurs, we mint back the native token.

```go
// x/ibcdex/keeper/sellOrder.go
func (k Keeper) OnTimeoutSellOrderPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SellOrderPacketData) error {
	// In case of error we mint back the native token
	receiver, err := sdk.AccAddressFromBech32(data.Seller)
	if err != nil {
		return err
	}

	if err := k.SafeMint(
		ctx,
		packet.SourcePort,
		packet.SourceChannel,
		receiver,
		data.AmountDenom,
		data.Amount,
	); err != nil {
		return err
	}

	return nil
}
```