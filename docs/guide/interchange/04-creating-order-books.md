---
order: 4
---
# Creating Order Books

In this chapter you will implement the logic for creating order books.

In Cosmos SDK the state is stored in a key-value store. Each order book will be stored under a unique key composed of four values: port ID, channel ID, source denom and target denom. For example, an order book for `mcx` and `vcx` could be stored under `ibcdex-channel-4-mcx-vcx`. Define a function that returns an order book store key.

```go
// x/ibcdex/types/keys.go
import "fmt"

//...
func OrderBookIndex( portID string, channelID string, sourceDenom string, targetDenom string, ) string {
  return fmt.Sprintf("%s-%s-%s-%s", portID, channelID, sourceDenom, targetDenom, )
}
```

`send-create-pair` is used to create order books. This command creates and broadcasts a transaction with a message of type `SendCreatePair`. The message gets routed to the `ibcdex` module, processed by the message handler in `x/ibcdex/handler.go` and finally a `SendCreatePair` keeper method is called.

You need `send-create-pair` to do the following:

* When processing `SendCreatePair` message on the source chain
  * Check that an order book with the given pair of denoms does not yet exist
  * Transmit an IBC packet with information about port, channel, source and target denoms
* Upon receiving the packet on the target chain
  * Check that an order book with the given pair of denoms does not yet exist on the target chain
  * Create a new order book for buy orders
  * Transmit an IBC acknowledgement back to the source chain
* Upon receiving the acknowledgement on the source chain
  * Create a new order book for sell orders

## `SendCreatePair` Message Handling

`SendCreatePair` function was created during the IBC packet scaffolding. Currently, it creates an IBC packet, populates it with source and target denoms and transmits this packet over IBC. Add the logic to check for an existing order book for a particular pair of denoms.

```go
// x/ibcdex/keeper/msg_server_create_pair.go
import (
  "errors"
  //...
)

func (k msgServer) SendCreatePair(goCtx context.Context, msg *types.MsgSendCreatePair) (*types.MsgSendCreatePairResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Get an order book index
	pairIndex := types.OrderBookIndex(msg.Port, msg.ChannelID, msg.SourceDenom, msg.TargetDenom)
	// If an order book is found, return an error
	_, found := k.GetSellOrderBook(ctx, pairIndex)
	if found {
		return &types.MsgSendCreatePairResponse{}, errors.New("the pair already exist")
	}

	// Construct the packet
	var packet types.CreatePairPacketData

	packet.SourceDenom = msg.SourceDenom
	packet.TargetDenom = msg.TargetDenom

	// Transmit the packet
	err := k.TransmitCreatePairPacket(
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

	return &types.MsgSendCreatePairResponse{}, nil
}
```

## Lifecycle of an IBC Packet

During a successful transmission, an IBC packet goes through 4 stages:

1. Message processing before packet transmission (on the source cahin)
2. Reception of a packet (on the target chain)
3. Acknowledgment of a packet (on the source chain)
4. Timeout of a packet (on the source chain)

In the following section you'll be implementing packet reception logic in the `OnRecvCreatePairPacket` function and packet acknowledgement logic in the `OnAcknowledgementCreatePairPacket` function. Timeout function will be left empty.

## `OnRecv`

On the target chain when an IBC packet is recieved, the module should check whether a book already exists, if not, create a new buy order book for specified denoms.

```go
// x/ibcdex/keeper/create_pair.go
func (k Keeper) OnRecvCreatePairPacket(ctx sdk.Context, packet channeltypes.Packet, data types.CreatePairPacketData) (packetAck types.CreatePairPacketAck, err error) {
  // ...
  // Get an order book index
  pairIndex := types.OrderBookIndex(packet.SourcePort, packet.SourceChannel, data.SourceDenom, data.TargetDenom)
  // If an order book is found, return an error
  _, found := k.GetBuyOrderBook(ctx, pairIndex)
  if found {
    return packetAck, errors.New("the pair already exist")
  }
  // Create a new buy order book for source and target denoms
  book := types.NewBuyOrderBook(data.SourceDenom, data.TargetDenom)
  // Assign order book index
  book.Index = pairIndex
  // Save the order book to the store
  k.SetBuyOrderBook(ctx, book)
  return packetAck, nil
}
```

Define `NewBuyOrderBook` creates a new buy order book.

```go
// x/ibcdex/types/buy_order_book.go
package types

func NewBuyOrderBook(AmountDenom string, PriceDenom string) BuyOrderBook {
	book := NewOrderBook()
	return BuyOrderBook{
		AmountDenom: AmountDenom,
		PriceDenom: PriceDenom,
		Book: &book,
	}
}
```

Modify the `buy_order_book.proto` file to have the fields for creating a buy order on the order book.

```proto
// proto/ibcdex/buy_order_book.proto
import "ibcdex/order.proto";

message BuyOrderBook {
  // ...
  OrderBook book = 5;
}
```

```go
// x/ibcdex/types/order_book.go
package types

func NewOrderBook() OrderBook {
	return OrderBook{
		IdCount: 0,
	}
}
```

The protocol buffer definition defines the data that an order book has. Add the `OrderBook` and `Order` messages to the `order.proto` file.

```proto
// proto/ibcdex/order.proto
syntax = "proto3";
package cosmonaut.interchange.ibcdex;

option go_package = "github.com/cosmonaut/interchange/x/ibcdex/types";

message OrderBook {
  int32 idCount = 1;
  repeated Order orders = 2;
}

message Order {
  int32 id = 1;
  string creator = 2;
  int32 amount = 3;
  int32 price = 4;
}
```

## `OnAcknowledgement`

On the source chain when an IBC acknowledgement is recieved, the module should check whether a book already exists, if not, create a new sell order book for specified denoms.

```go
// x/ibcdex/keeper/create_pair.go
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

`NewSellOrderBook` creates a new sell order book.

```go
// x/ibcdex/types/sell_order_book.go
package types

func NewSellOrderBook(AmountDenom string, PriceDenom string) SellOrderBook {
	book := NewOrderBook()
	return SellOrderBook{
		AmountDenom: AmountDenom,
		PriceDenom: PriceDenom,
		Book: &book,
	}
}
```

Modify the `sell_order_book.proto` file to add the order book into the buy order book. The proto definition for the `SellOrderBook` should look like follows:

```proto
// proto/ibcdex/sell_order_book.proto
// ...
import "ibcdex/order.proto";

message SellOrderBook {
  // ...
  OrderBook book = 6;
}
```

In this chapter we implemented the logic behind `send-create-pair` command that upon recieving of an IBC packet on the target chain creates a buy order book and upon recieving of an IBC acknowledgement on the source chain creates a sell order book.
