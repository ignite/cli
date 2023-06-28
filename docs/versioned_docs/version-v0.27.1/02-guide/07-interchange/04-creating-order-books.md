---
sidebar_position: 4
description: Implement logic to create order books.
---

# Implement the Order Books

In this chapter, you implement the logic to create order books.

In the Cosmos SDK, the state is stored in a key-value store. Each order book is stored under a unique key that is
composed of four values:

- Port ID
- Channel ID
- Source denom
- Target denom

For example, an order book for marscoin and venuscoin could be stored under `dex-channel-4-marscoin-venuscoin`.

First, define a function that returns an order book store key:

```go
// x/dex/types/keys.go
package types

import "fmt"

// ...
func OrderBookIndex(portID string, channelID string, sourceDenom string, targetDenom string) string {
	return fmt.Sprintf("%s-%s-%s-%s", portID, channelID, sourceDenom, targetDenom)
}
```

The `send-create-pair` command is used to create order books. This command:

- Creates and broadcasts a transaction with a message of type `SendCreatePair`.
- The message gets routed to the `dex` module.
- Finally, a `SendCreatePair` keeper method is called.

You need the `send-create-pair` command to do the following:

- When processing `SendCreatePair` message on the source chain:
    - Check that an order book with the given pair of denoms does not yet exist.
    - Transmit an IBC packet with information about port, channel, source denoms, and target denoms.
- After the packet is received on the target chain:
    - Check that an order book with the given pair of denoms does not yet exist on the target chain.
    - Create a new order book for buy orders.
    - Transmit an IBC acknowledgement back to the source chain.
- After the acknowledgement is received on the source chain:
    - Create a new order book for sell orders.

## Message Handling in SendCreatePair

The `SendCreatePair` function was created during the IBC packet scaffolding. The function creates an IBC packet,
populates it with source and target denoms, and transmits this packet over IBC.

Now, add the logic to check for an existing order book for a particular pair of denoms:

```go
// x/dex/keeper/msg_server_create_pair.go

package keeper

import (
	"errors"
	// ...
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
	_, err := k.TransmitCreatePairPacket(
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

During a successful transmission, an IBC packet goes through these stages:

1. Message processing before packet transmission on the source chain
2. Reception of a packet on the target chain
3. Acknowledgment of a packet on the source chain
4. Timeout of a packet on the source chain

In the following section, implement the packet reception logic in the `OnRecvCreatePairPacket` function and the packet
acknowledgement logic in the `OnAcknowledgementCreatePairPacket` function.

Leave the Timeout function empty.

## Receive an IBC packet

The protocol buffer definition defines the data that an order book contains.

Add the `OrderBook` and `Order` messages to the `order.proto` file.

First, add the proto buffer files to build the Go code files. You can modify these files for the purpose of your app.

Create a new `order.proto` file in the `proto/interchange/dex` directory and add the content:

```protobuf
// proto/interchange/dex/order.proto

syntax = "proto3";

package interchange.dex;

option go_package = "interchange/x/dex/types";

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

Modify the `buy_order_book.proto` file to have the fields for creating a buy order on the order book.
Don't forget to add the import as well.

**Tip:** Don't forget to add the import as well.

```protobuf
// proto/interchange/dex/buy_order_book.proto

// ...

import "interchange/dex/order.proto";

message BuyOrderBook {
  // ...
  OrderBook book = 4;
}
```

Modify the `sell_order_book.proto` file to add the order book into the buy order book.

The proto definition for the `SellOrderBook` looks like:

```protobuf
// proto/interchange/dex/sell_order_book.proto

// ...
import "interchange/dex/order.proto";

message SellOrderBook {
  // ...
  OrderBook book = 4;
}
```

Now, use Ignite CLI to build the proto files for the `send-create-pair` command:

```bash
ignite generate proto-go --yes
```

Start enhancing the functions for the IBC packets.

Create a new file `x/dex/types/order_book.go`.

Add the new order book function to the corresponding Go file:

```go
// x/dex/types/order_book.go

package types

func NewOrderBook() OrderBook {
	return OrderBook{
		IdCount: 0,
	}
}
```

To create a new buy order book type, define `NewBuyOrderBook` in a new file `x/dex/types/buy_order_book.go` :

```go
// x/dex/types/buy_order_book.go

package types

func NewBuyOrderBook(AmountDenom string, PriceDenom string) BuyOrderBook {
	book := NewOrderBook()
	return BuyOrderBook{
		AmountDenom: AmountDenom,
		PriceDenom:  PriceDenom,
		Book:        &book,
	}
}
```

When an IBC packet is received on the target chain, the module must check whether a book already exists. If not, then
create a buy order book for the specified denoms.

```go
// x/dex/keeper/create_pair.go

package keeper

// ...

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

## Receive an IBC Acknowledgement

When an IBC acknowledgement is received on the source chain, the module must check whether a book already exists. If
not,
create a sell order book for the specified denoms.

Create a new file `x/dex/types/sell_order_book.go`.
Insert the `NewSellOrderBook` function which creates a new sell order book.

```go
// x/dex/types/sell_order_book.go

package types

func NewSellOrderBook(AmountDenom string, PriceDenom string) SellOrderBook {
	book := NewOrderBook()
	return SellOrderBook{
		AmountDenom: AmountDenom,
		PriceDenom:  PriceDenom,
		Book:        &book,
	}
}
```

Modify the Acknowledgement function in the `x/dex/keeper/create_pair.go` file:

```go
// x/dex/keeper/create_pair.go

package keeper

// ...

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

In this section, you implemented the logic behind the new `send-create-pair` command:

- When an IBC packet is received on the target chain, `send-create-pair` command creates a buy order book.
- When an IBC acknowledgement is received on the source chain, the `send-create-pair` command creates a sell order book.

### Implement the appendOrder Function to Add Orders to the Order Book

```go
// x/dex/types/order_book.go

package types

import (
	"errors"
	"sort"
)

func NewOrderBook() OrderBook {
	return OrderBook{
		IdCount: 0,
	}
}

const (
	MaxAmount = int32(100000)
	MaxPrice  = int32(100000)
)

type Ordering int

const (
	Increasing Ordering = iota
	Decreasing
)

var (
	ErrMaxAmount     = errors.New("max amount reached")
	ErrMaxPrice      = errors.New("max price reached")
	ErrZeroAmount    = errors.New("amount is zero")
	ErrZeroPrice     = errors.New("price is zero")
	ErrOrderNotFound = errors.New("order not found")
)
```

The `AppendOrder` function initializes and appends a new order to an order book from the order information:

```go
// x/dex/types/order_book.go

func (book *OrderBook) appendOrder(creator string, amount int32, price int32, ordering Ordering) (int32, error) {
	if err := checkAmountAndPrice(amount, price); err != nil {
		return 0, err
	}

	// Initialize the order
	var order Order
	order.Id = book.GetNextOrderID()
	order.Creator = creator
	order.Amount = amount
	order.Price = price

	// Increment ID tracker
	book.IncrementNextOrderID()

	// Insert the order
	book.insertOrder(order, ordering)
	return order.Id, nil
}
```

#### Implement the checkAmountAndPrice Function For an Order

The `checkAmountAndPrice` function checks for the correct amount or price:

```go
// x/dex/types/order_book.go

func checkAmountAndPrice(amount int32, price int32) error {
	if amount == int32(0) {
		return ErrZeroAmount
	}
	if amount > MaxAmount {
		return ErrMaxAmount
	}

	if price == int32(0) {
		return ErrZeroPrice
	}
	if price > MaxPrice {
		return ErrMaxPrice
	}

	return nil
}
```

#### Implement the GetNextOrderID Function

The `GetNextOrderID` function gets the ID of the next order to append:

```go
// x/dex/types/order_book.go

func (book OrderBook) GetNextOrderID() int32 {
	return book.IdCount
}
```

#### Implement the IncrementNextOrderID Function

The `IncrementNextOrderID` function updates the ID count for orders:

```go
// x/dex/types/order_book.go

func (book *OrderBook) IncrementNextOrderID() {
	// Even numbers to have different ID than buy orders
	book.IdCount++
}
```

#### Implement the insertOrder Function

The `insertOrder` function inserts the order in the book with the provided order:

```go
// x/dex/types/order_book.go

func (book *OrderBook) insertOrder(order Order, ordering Ordering) {
	if len(book.Orders) > 0 {
		var i int

		// get the index of the new order depending on the provided ordering
		if ordering == Increasing {
			i = sort.Search(len(book.Orders), func(i int) bool { return book.Orders[i].Price > order.Price })
		} else {
			i = sort.Search(len(book.Orders), func(i int) bool { return book.Orders[i].Price < order.Price })
		}

		// insert order
		orders := append(book.Orders, &order)
		copy(orders[i+1:], orders[i:])
		orders[i] = &order
		book.Orders = orders
	} else {
		book.Orders = append(book.Orders, &order)
	}
}
```

This completes the order book setup.

Now is a good time to save the state of your implementation.
Because your project is in a local repository, you can use git. Saving your current state lets you jump back and forth
in case you introduce errors or need a break.

```bash
git add .
git commit -m "Create Order Books"
```

In the next chapter, you learn how to deal with vouchers by minting and burning vouchers and locking and unlocking
native blockchain token in your app.
