---
order: 4
---

# Create the Order Book

In this chapter you implement the code for the order book, for buy orders and sell orders.

In this chapter you will create a `order_book.go` file with the implementation for the order book. 
The order book will allow to publish buy or sell orders. The order book for a certain pair of token has to be registered first. After registering the order book for a pair of token, you can add sell orders and buy orders.
You will create the `sell_order_book.go` file with the implementation of a sell order book. Sell order books contain sell orders that contain the data of the token denomination and a price you offer to sell a token for.
You will create the `buy_order_book.go` file with the implementation of a buy order book. Buy order books contain buy orders that contain the data of the token denomination and a price you offer to buy a token for. Buy orders and sell orders will live on different blockchain apps.
When a buy and a sell order match, the exchange will be executed.

## Add The Order Book

The protobuffer definition defines the data that an order book has. 
Add the `OrderBook` and `Order` messages to the `order.proto` file.

```proto
// proto/ibcdex/order.proto
syntax = "proto3";
package username.interchange.ibcdex;

option go_package = "github.com/username/interchange/x/ibcdex/types";

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

Create a new file `order_book.go` in the `ibcdex` module `types` directory.
In this file, you will define the logic to create a new order book. 
This is the common logic between sell and buy order books.

Create a `order_book.go` file and add the code:

```go
// x/ibcdex/types/order_book.go
package types

import (
	"errors"
	"sort"
)

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

// checkAmountAndPrice checks correct amount or price
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

func NewOrderBook() OrderBook {
	return OrderBook{
		IdCount: 0,
	}
}

// GetOrder gets the order from an index
func (book OrderBook) GetOrder(index int) (order Order, err error) {
	if index >= len(book.Orders) {
		return order, ErrOrderNotFound
	}

	return *book.Orders[index], nil
}

// GetNextOrderID gets the ID of the next order to append
func (book OrderBook) GetNextOrderID() int32 {
	return book.IdCount
}

// GetOrderFromID gets an order from the book from its id
func (book OrderBook) GetOrderFromID(id int32) (Order, error) {
	for _, order := range book.Orders {
		if order.Id == id {
			return *order, nil
		}
	}
	return Order{}, ErrOrderNotFound
}

// SetOrder gets the order from an index
func (book *OrderBook) SetOrder(index int, order Order) error {
	if index >= len(book.Orders) {
		return ErrOrderNotFound
	}

	book.Orders[index] = &order

	return nil
}

// IncrementNextOrderID updates the ID count for orders
func (book *OrderBook) IncrementNextOrderID() {
	// Even numbers to have different ID than buy orders
	book.IdCount++
}

// RemoveOrderFromID removes an order from the book and keep it ordered
func (book *OrderBook) RemoveOrderFromID(id int32) error {
	for i, order := range book.Orders {
		if order.Id == id {
			book.Orders = append(book.Orders[:i], book.Orders[i+1:]...)
			return nil
		}
	}
	return ErrOrderNotFound
}

// AppendOrder initializes and appends a new order in a book from order information
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

// insertOrder inserts the order in the book with the provided order
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

## Add The Sellorder

Modify the `sellOrderBook.proto` file to add the order book into the sell order book.
The proto definition for the `SellOrderBook` should look like follows:

```proto
// proto/ibcdex/sellOrderBook.proto
syntax = "proto3";
package username.interchange.ibcdex;

option go_package = "github.com/username/interchange/x/ibcdex/types";

import "ibcdex/order.proto"; // <--

message SellOrderBook {
  string creator = 1;
  string index = 2;
  string amountDenom = 3;
  string priceDenom = 4;
  OrderBook book = 5; // <--
}
```

For the code of the sell order book, create a `sell_order_book.go` file in the `types` directory and add the following code:

```go
// x/ibcdex/types/sell_order_book.go
package types

// NewSellOrderBook creates a new sell order book
func NewSellOrderBook(AmountDenom string, PriceDenom string) SellOrderBook {
	book := NewOrderBook()
	return SellOrderBook{
		AmountDenom: AmountDenom,
		PriceDenom: PriceDenom,
		Book: &book,
	}
}

// AppendOrder appends an order in sell order book
func (s *SellOrderBook) AppendOrder(creator string, amount int32, price int32) (int32, error) {
	return s.Book.appendOrder(creator, amount, price, Decreasing)
}

// LiquidateFromBuyOrder liquidates the first sell order of the book from the buy order
// if no match is found, return false for match
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

// FillBuyOrder try to fill the buy order with the order book and returns all the side effects
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

## Add The Buyorder

Modify the `buyOrderBook.proto` file to add the order book into the buy order book.
The proto definition for the `BuyOrderBook` should look like follows:

```proto
// proto/ibcdex/buyOrderBook.proto
syntax = "proto3";
package username.interchange.ibcdex;

option go_package = "github.com/username/interchange/x/ibcdex/types";

import "ibcdex/order.proto"; // <--

message BuyOrderBook {
  string creator = 1;
  string index = 2;
  string amountDenom = 3;
  string priceDenom = 4;
  OrderBook book = 5; // <--
}
```

For the buy order book, create a `buy_order_book.go` file in the `types` directory and add the following code:

```go
 // x/ibcdex/types/buy_order_book.go
package types

// NewBuyOrderBook creates a new buy order book
func NewBuyOrderBook(AmountDenom string, PriceDenom string) BuyOrderBook {
	book := NewOrderBook()
	return BuyOrderBook{
		AmountDenom: AmountDenom,
		PriceDenom: PriceDenom,
		Book: &book,
	}
}

 // AppendOrder appends an order in buy order book
 func (b *BuyOrderBook) AppendOrder(creator string, amount int32, price int32) (int32, error) {
	 return b.Book.appendOrder(creator, amount, price, Increasing)
 }

 // LiquidateFromSellOrder liquidates the first buy order of the book from the sell order
 // if no match is found, return false for match
 func (b *BuyOrderBook) LiquidateFromSellOrder(order Order) (
	 remainingSellOrder Order,
	 liquidatedBuyOrder Order,
	 gain int32,
	 match bool,
	 filled bool,
 ) {
	 remainingSellOrder = order

	 // No match if no order
	 orderCount := len(b.Book.Orders)
	 if orderCount == 0 {
		 return order, liquidatedBuyOrder, gain, false, false
	 }

	 // Check if match
	 highestBid := b.Book.Orders[orderCount-1]
	 if order.Price > highestBid.Price {
		 return order, liquidatedBuyOrder, gain, false, false
	 }

	 liquidatedBuyOrder = *highestBid

	 // Check if sell order can be entirely filled
	 if highestBid.Amount >= order.Amount {
		 remainingSellOrder.Amount = 0
		 liquidatedBuyOrder.Amount = order.Amount
		 gain = order.Amount * highestBid.Price

		 // Remove highest bid if it has been entirely liquidated
		 highestBid.Amount -= order.Amount
		 if highestBid.Amount == 0 {
			 b.Book.Orders = b.Book.Orders[:orderCount-1]
		 } else {
			 b.Book.Orders[orderCount-1] = highestBid
		 }
		 return remainingSellOrder, liquidatedBuyOrder, gain, true, true
	 }

	 // Not entirely filled
	 gain = highestBid.Amount * highestBid.Price
	 b.Book.Orders = b.Book.Orders[:orderCount-1]
	 remainingSellOrder.Amount -= highestBid.Amount

	 return remainingSellOrder, liquidatedBuyOrder, gain, true, false
 }


 // FillSellOrder try to fill the sell order with the order book and returns all the side effects
func (b *BuyOrderBook) FillSellOrder(order Order) (
	remainingSellOrder Order,
	liquidated []Order,
	gain int32,
	filled bool,
) {
	var liquidatedList []Order
	totalGain := int32(0)
	remainingSellOrder = order

	// Liquidate as long as there is match
	for {
		var match bool
		var liquidation Order
		remainingSellOrder, liquidation, gain, match, filled = b.LiquidateFromSellOrder(
			remainingSellOrder,
		)
		if !match {
			break
		}

		// Update gains
		totalGain += gain

		// Update liquidated
		liquidatedList = append(liquidatedList, liquidation)

		if filled {
			break
		}
	}

	return remainingSellOrder, liquidatedList, totalGain, filled
}
```

This finishes your code for the order book module with buy and sell orders.
In the next chapters, you will make them IBC compatible. 
You will have to implement how IBC packets are handled that are sent over a blockchain.
These packets will be received and acknowledged by the recipient blockchain.