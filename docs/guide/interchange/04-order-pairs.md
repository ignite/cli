---
order: 5
---

# Create Order Pairs

In the last chapter you have created the logic for the order book with buy and sell orders.
In this chapter you will build the software to make the order book IBC compatible.
This will require you to create the index for an order book, with their according messages on receive and on acknowledgement of the IBC packet.

## Define the Order Book Indexes

You defined the order books as indexed types in the store but you have to specify how this index is defined with the following code:

```go
// x/ibcdex/types/keys.go
import "fmt"

//...

func OrderBookIndex(
	portID string,
	channelID string,
	sourceDenom string,
	targetDenom string,
) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		portID,
		channelID,
		sourceDenom,
		targetDenom,
	)
}

```

## Check for Existing Pairs

A pair of token always has one order book that everyone can access in the app.
When an order book pair already exists, it should throw an error.

```go
// x/ibcdex/keeper/msg_server_createPair.go
import "errors"

//...

func (k msgServer) SendCreatePair(goCtx context.Context, msg *types.MsgSendCreatePair) (*types.MsgSendCreatePairResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Cannot create the pair if it already exists
	pairIndex := types.OrderBookIndex(msg.Port, msg.ChannelID, msg.SourceDenom, msg.TargetDenom)
	_, found := k.GetSellOrderBook(ctx, pairIndex)
	if found {
		return &types.MsgSendCreatePairResponse{}, errors.New("the pair already exist")
	}

	// ...
}
```

## Create the OnRecv Function

When a packet with an order book creation is received, the validity of the transaction should be check with the `ValidateBasic()` function.
If the pair does not exist yet, it can be added to the keeper.

```go
// x/ibcdex/keeper/createPair.go
func (k Keeper) OnRecvCreatePairPacket(ctx sdk.Context, packet channeltypes.Packet, data types.CreatePairPacketData) (packetAck types.CreatePairPacketAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	// Check if the buy order book exists
	pairIndex := types.OrderBookIndex(packet.SourcePort, packet.SourceChannel, data.SourceDenom, data.TargetDenom)
	_, found := k.GetBuyOrderBook(ctx, pairIndex)
	if found {
		return packetAck, errors.New("the pair already exist")
	}

	// Set the buy order book
	book := types.NewBuyOrderBook(data.SourceDenom, data.TargetDenom)
	book.Index = pairIndex
	k.SetBuyOrderBook(ctx, book)

	return packetAck, nil
}
```

## Create the OnAcknowledgement Function

When a packet sent with IBC is valid and received, it must be acknowledged.
When the acknowledgement is successful, add the sell order book to the database.

```go
// x/ibcdex/keeper/createPair.go
func (k Keeper) OnAcknowledgementCreatePairPacket(ctx sdk.Context, packet channeltypes.Packet, data types.CreatePairPacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.CreatePairPacketAck

		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		// Set the sell order book
		pairIndex := types.OrderBookIndex(packet.SourcePort, packet.SourceChannel, data.SourceDenom, data.TargetDenom)
		book := types.NewSellOrderBook(data.SourceDenom, data.TargetDenom)
		book.Index = pairIndex
		k.SetSellOrderBook(ctx, book)

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}
```

## No Consequences OnTimeout

A timeout will have no consequences in our scenario. The order book will not be acknowledged 