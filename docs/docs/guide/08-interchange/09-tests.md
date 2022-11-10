---
sidebar_position: 9
description: Add test files.
---

# Write Test Files

To test your application, add the test files to your code.

After you add the test files, change into the `interchange` directory with your terminal, then run:

```bash
go test -timeout 30s ./x/dex/types
```

## Order Book Tests

Create a new `x/dex/types/order_book_test.go` file in the `types` directory.

Add the following testsuite:

```go
// x/dex/types/order_book_test.go

package types_test

import (
	"math/rand"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"interchange/x/dex/types"
)

func GenString(n int) string {
	alpha := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	buf := make([]rune, n)
	for i := range buf {
		buf[i] = alpha[rand.Intn(len(alpha))]
	}

	return string(buf)
}

func GenAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr).String()
}

func GenAmount() int32 {
	return int32(rand.Intn(int(types.MaxAmount)) + 1)
}

func GenPrice() int32 {
	return int32(rand.Intn(int(types.MaxPrice)) + 1)
}

func GenPair() (string, string) {
	return GenString(10), GenString(10)
}

func GenOrder() (string, int32, int32) {
	return GenLocalAccount(), GenAmount(), GenPrice()
}

func GenLocalAccount() string {
	return GenAddress()
}

func MockAccount(str string) string {
	return str
}

func OrderListToOrderBook(list []types.Order) types.OrderBook {
	listCopy := make([]*types.Order, len(list))
	for i, order := range list {
		order := order
		listCopy[i] = &order
	}

	return types.OrderBook{
		IdCount: 0,
		Orders:  listCopy,
	}
}

func TestRemoveOrderFromID(t *testing.T) {
	inputList := []types.Order{
		{Id: 3, Creator: MockAccount("3"), Amount: 2, Price: 10},
		{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
		{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
		{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
	}

	book := OrderListToOrderBook(inputList)
	expectedList := []types.Order{
		{Id: 3, Creator: MockAccount("3"), Amount: 2, Price: 10},
		{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
		{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
	}
	expectedBook := OrderListToOrderBook(expectedList)
	err := book.RemoveOrderFromID(2)
	require.NoError(t, err)
	require.Equal(t, expectedBook, book)

	book = OrderListToOrderBook(inputList)
	expectedList = []types.Order{
		{Id: 3, Creator: MockAccount("3"), Amount: 2, Price: 10},
		{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
		{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
	}
	expectedBook = OrderListToOrderBook(expectedList)
	err = book.RemoveOrderFromID(0)
	require.NoError(t, err)
	require.Equal(t, expectedBook, book)

	book = OrderListToOrderBook(inputList)
	expectedList = []types.Order{
		{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
		{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
		{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
	}
	expectedBook = OrderListToOrderBook(expectedList)
	err = book.RemoveOrderFromID(3)
	require.NoError(t, err)
	require.Equal(t, expectedBook, book)

	book = OrderListToOrderBook(inputList)
	err = book.RemoveOrderFromID(4)
	require.ErrorIs(t, err, types.ErrOrderNotFound)
}
```

## Buy Order Tests

Create a new `x/dex/types/buy_order_book_test.go` file in the `types` directory to add the tests for the Buy Order Book:

```go
// x/dex/types/buy_order_book_test.go

package types_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"interchange/x/dex/types"
)

func OrderListToBuyOrderBook(list []types.Order) types.BuyOrderBook {
	listCopy := make([]*types.Order, len(list))
	for i, order := range list {
		order := order
		listCopy[i] = &order
	}

	book := types.BuyOrderBook{
		AmountDenom: "foo",
		PriceDenom:  "bar",
		Book: &types.OrderBook{
			IdCount: 0,
			Orders:  listCopy,
		},
	}
	return book
}

func TestAppendOrder(t *testing.T) {
	buyBook := types.NewBuyOrderBook(GenPair())

	// Prevent zero amount
	seller, amount, price := GenOrder()
	_, err := buyBook.AppendOrder(seller, 0, price)
	require.ErrorIs(t, err, types.ErrZeroAmount)

	// Prevent big amount
	_, err = buyBook.AppendOrder(seller, types.MaxAmount+1, price)
	require.ErrorIs(t, err, types.ErrMaxAmount)

	// Prevent zero price
	_, err = buyBook.AppendOrder(seller, amount, 0)
	require.ErrorIs(t, err, types.ErrZeroPrice)

	// Prevent big price
	_, err = buyBook.AppendOrder(seller, amount, types.MaxPrice+1)
	require.ErrorIs(t, err, types.ErrMaxPrice)

	// Can append buy orders
	for i := 0; i < 20; i++ {
		// Append a new order
		creator, amount, price := GenOrder()
		newOrder := types.Order{
			Id:      buyBook.Book.IdCount,
			Creator: creator,
			Amount:  amount,
			Price:   price,
		}
		orderID, err := buyBook.AppendOrder(creator, amount, price)

		// Checks
		require.NoError(t, err)
		require.Contains(t, buyBook.Book.Orders, &newOrder)
		require.Equal(t, newOrder.Id, orderID)
	}

	require.Len(t, buyBook.Book.Orders, 20)
	require.True(t, sort.SliceIsSorted(buyBook.Book.Orders, func(i, j int) bool {
		return buyBook.Book.Orders[i].Price < buyBook.Book.Orders[j].Price
	}))
}

type liquidateSellRes struct {
	Book       []types.Order
	Remaining  types.Order
	Liquidated types.Order
	Gain       int32
	Match      bool
	Filled     bool
}

func simulateLiquidateFromSellOrder(
	t *testing.T,
	inputList []types.Order,
	inputOrder types.Order,
	expected liquidateSellRes,
) {
	book := OrderListToBuyOrderBook(inputList)
	expectedBook := OrderListToBuyOrderBook(expected.Book)

	require.True(t, sort.SliceIsSorted(book.Book.Orders, func(i, j int) bool {
		return book.Book.Orders[i].Price < book.Book.Orders[j].Price
	}))
	require.True(t, sort.SliceIsSorted(expectedBook.Book.Orders, func(i, j int) bool {
		return expectedBook.Book.Orders[i].Price < expectedBook.Book.Orders[j].Price
	}))

	remaining, liquidated, gain, match, filled := book.LiquidateFromSellOrder(inputOrder)

	require.Equal(t, expectedBook, book)
	require.Equal(t, expected.Remaining, remaining)
	require.Equal(t, expected.Liquidated, liquidated)
	require.Equal(t, expected.Gain, gain)
	require.Equal(t, expected.Match, match)
	require.Equal(t, expected.Filled, filled)
}

func TestLiquidateFromSellOrder(t *testing.T) {
	// No match for empty book
	inputOrder := types.Order{Id: 10, Creator: MockAccount("1"), Amount: 100, Price: 30}
	book := OrderListToBuyOrderBook([]types.Order{})
	_, _, _, match, _ := book.LiquidateFromSellOrder(inputOrder)
	require.False(t, match)

	// Buy book
	inputBook := []types.Order{
		{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
		{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
		{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
	}

	// Test no match if highest bid too low (25 < 30)
	book = OrderListToBuyOrderBook(inputBook)
	_, _, _, match, _ = book.LiquidateFromSellOrder(inputOrder)
	require.False(t, match)

	// Entirely filled (30 < 50)
	inputOrder = types.Order{Id: 10, Creator: MockAccount("1"), Amount: 30, Price: 22}
	expected := liquidateSellRes{
		Book: []types.Order{
			{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
			{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
			{Id: 0, Creator: MockAccount("0"), Amount: 20, Price: 25},
		},
		Remaining:  types.Order{Id: 10, Creator: MockAccount("1"), Amount: 0, Price: 22},
		Liquidated: types.Order{Id: 0, Creator: MockAccount("0"), Amount: 30, Price: 25},
		Gain:       int32(30 * 25),
		Match:      true,
		Filled:     true,
	}
	simulateLiquidateFromSellOrder(t, inputBook, inputOrder, expected)

	// Entirely filled and liquidated ( 50 = 50)
	inputOrder = types.Order{Id: 10, Creator: MockAccount("1"), Amount: 50, Price: 15}
	expected = liquidateSellRes{
		Book: []types.Order{
			{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
			{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
		},
		Remaining:  types.Order{Id: 10, Creator: MockAccount("1"), Amount: 0, Price: 15},
		Liquidated: types.Order{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
		Gain:       int32(50 * 25),
		Match:      true,
		Filled:     true,
	}
	simulateLiquidateFromSellOrder(t, inputBook, inputOrder, expected)

	// Not filled and entirely liquidated (60 > 50)
	inputOrder = types.Order{Id: 10, Creator: MockAccount("1"), Amount: 60, Price: 10}
	expected = liquidateSellRes{
		Book: []types.Order{
			{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
			{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
		},
		Remaining:  types.Order{Id: 10, Creator: MockAccount("1"), Amount: 10, Price: 10},
		Liquidated: types.Order{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
		Gain:       int32(50 * 25),
		Match:      true,
		Filled:     false,
	}
	simulateLiquidateFromSellOrder(t, inputBook, inputOrder, expected)
}

type fillSellRes struct {
	Book       []types.Order
	Remaining  types.Order
	Liquidated []types.Order
	Gain       int32
	Filled     bool
}

func simulateFillSellOrder(
	t *testing.T,
	inputList []types.Order,
	inputOrder types.Order,
	expected fillSellRes,
) {
	book := OrderListToBuyOrderBook(inputList)
	expectedBook := OrderListToBuyOrderBook(expected.Book)

	require.True(t, sort.SliceIsSorted(book.Book.Orders, func(i, j int) bool {
		return book.Book.Orders[i].Price < book.Book.Orders[j].Price
	}))
	require.True(t, sort.SliceIsSorted(expectedBook.Book.Orders, func(i, j int) bool {
		return expectedBook.Book.Orders[i].Price < expectedBook.Book.Orders[j].Price
	}))

	remaining, liquidated, gain, filled := book.FillSellOrder(inputOrder)

	require.Equal(t, expectedBook, book)
	require.Equal(t, expected.Remaining, remaining)
	require.Equal(t, expected.Liquidated, liquidated)
	require.Equal(t, expected.Gain, gain)
	require.Equal(t, expected.Filled, filled)
}

func TestFillSellOrder(t *testing.T) {
	var inputBook []types.Order

	// Empty book
	inputOrder := types.Order{Id: 10, Creator: MockAccount("1"), Amount: 30, Price: 30}
	expected := fillSellRes{
		Book:       []types.Order{},
		Remaining:  inputOrder,
		Liquidated: []types.Order(nil),
		Gain:       int32(0),
		Filled:     false,
	}
	simulateFillSellOrder(t, inputBook, inputOrder, expected)

	// No match
	inputBook = []types.Order{
		{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
		{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
		{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
	}
	expected = fillSellRes{
		Book:       inputBook,
		Remaining:  inputOrder,
		Liquidated: []types.Order(nil),
		Gain:       int32(0),
		Filled:     false,
	}
	simulateFillSellOrder(t, inputBook, inputOrder, expected)

	// First order liquidated, not filled
	inputOrder = types.Order{Id: 10, Creator: MockAccount("1"), Amount: 60, Price: 22}
	expected = fillSellRes{
		Book: []types.Order{
			{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
			{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
		},
		Remaining: types.Order{Id: 10, Creator: MockAccount("1"), Amount: 10, Price: 22},
		Liquidated: []types.Order{
			{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
		},
		Gain:   int32(50 * 25),
		Filled: false,
	}
	simulateFillSellOrder(t, inputBook, inputOrder, expected)

	// Filled with two order
	inputOrder = types.Order{Id: 10, Creator: MockAccount("1"), Amount: 60, Price: 18}
	expected = fillSellRes{
		Book: []types.Order{
			{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
			{Id: 1, Creator: MockAccount("1"), Amount: 190, Price: 20},
		},
		Remaining: types.Order{Id: 10, Creator: MockAccount("1"), Amount: 0, Price: 18},
		Liquidated: []types.Order{
			{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
			{Id: 1, Creator: MockAccount("1"), Amount: 10, Price: 20},
		},
		Gain:   int32(50*25 + 10*20),
		Filled: true,
	}
	simulateFillSellOrder(t, inputBook, inputOrder, expected)

	// Not filled, buy order book liquidated
	inputOrder = types.Order{Id: 10, Creator: MockAccount("1"), Amount: 300, Price: 10}
	expected = fillSellRes{
		Book:      []types.Order{},
		Remaining: types.Order{Id: 10, Creator: MockAccount("1"), Amount: 20, Price: 10},
		Liquidated: []types.Order{
			{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
			{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
			{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
		},
		Gain:   int32(50*25 + 200*20 + 30*15),
		Filled: false,
	}
	simulateFillSellOrder(t, inputBook, inputOrder, expected)
}
```

## Sell Order Tests

Create a new testsuite for Sell Orders in a new file `x/dex/types/sell_order_book_test.go`:

```go
// x/dex/types/sell_order_book_test.go

package types_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"interchange/x/dex/types"
)

func OrderListToSellOrderBook(list []types.Order) types.SellOrderBook {
	listCopy := make([]*types.Order, len(list))
	for i, order := range list {
		order := order
		listCopy[i] = &order
	}

	book := types.SellOrderBook{
		AmountDenom: "foo",
		PriceDenom:  "bar",
		Book: &types.OrderBook{
			IdCount: 0,
			Orders:  listCopy,
		},
	}
	return book
}

func TestSellOrderBook_AppendOrder(t *testing.T) {
	sellBook := types.NewSellOrderBook(GenPair())

	// Prevent zero amount
	seller, amount, price := GenOrder()
	_, err := sellBook.AppendOrder(seller, 0, price)
	require.ErrorIs(t, err, types.ErrZeroAmount)

	// Prevent big amount
	_, err = sellBook.AppendOrder(seller, types.MaxAmount+1, price)
	require.ErrorIs(t, err, types.ErrMaxAmount)

	// Prevent zero price
	_, err = sellBook.AppendOrder(seller, amount, 0)
	require.ErrorIs(t, err, types.ErrZeroPrice)

	// Prevent big price
	_, err = sellBook.AppendOrder(seller, amount, types.MaxPrice+1)
	require.ErrorIs(t, err, types.ErrMaxPrice)

	// Can append sell orders
	for i := 0; i < 20; i++ {
		// Append a new order
		creator, amount, price := GenOrder()
		newOrder := types.Order{
			Id:      sellBook.Book.IdCount,
			Creator: creator,
			Amount:  amount,
			Price:   price,
		}
		orderID, err := sellBook.AppendOrder(creator, amount, price)

		// Checks
		require.NoError(t, err)
		require.Contains(t, sellBook.Book.Orders, &newOrder)
		require.Equal(t, newOrder.Id, orderID)
	}
	require.Len(t, sellBook.Book.Orders, 20)
	require.True(t, sort.SliceIsSorted(sellBook.Book.Orders, func(i, j int) bool {
		return sellBook.Book.Orders[i].Price > sellBook.Book.Orders[j].Price
	}))
}

type liquidateBuyRes struct {
	Book       []types.Order
	Remaining  types.Order
	Liquidated types.Order
	Purchase   int32
	Match      bool
	Filled     bool
}

func simulateLiquidateFromBuyOrder(
	t *testing.T,
	inputList []types.Order,
	inputOrder types.Order,
	expected liquidateBuyRes,
) {
	book := OrderListToSellOrderBook(inputList)
	expectedBook := OrderListToSellOrderBook(expected.Book)
	require.True(t, sort.SliceIsSorted(book.Book.Orders, func(i, j int) bool {
		return book.Book.Orders[i].Price > book.Book.Orders[j].Price
	}))
	require.True(t, sort.SliceIsSorted(expectedBook.Book.Orders, func(i, j int) bool {
		return expectedBook.Book.Orders[i].Price > expectedBook.Book.Orders[j].Price
	}))

	remaining, liquidated, purchase, match, filled := book.LiquidateFromBuyOrder(inputOrder)

	require.Equal(t, expectedBook, book)
	require.Equal(t, expected.Remaining, remaining)
	require.Equal(t, expected.Liquidated, liquidated)
	require.Equal(t, expected.Purchase, purchase)
	require.Equal(t, expected.Match, match)
	require.Equal(t, expected.Filled, filled)
}

func TestLiquidateFromBuyOrder(t *testing.T) {
	// No match for empty book
	inputOrder := types.Order{Id: 10, Creator: MockAccount("1"), Amount: 100, Price: 10}
	book := OrderListToSellOrderBook([]types.Order{})
	_, _, _, match, _ := book.LiquidateFromBuyOrder(inputOrder)
	require.False(t, match)

	// Sell book
	inputBook := []types.Order{
		{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
		{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
		{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
	}

	// Test no match if lowest ask too high (25 < 30)
	book = OrderListToSellOrderBook(inputBook)
	_, _, _, match, _ = book.LiquidateFromBuyOrder(inputOrder)
	require.False(t, match)

	// Entirely filled (30 > 15)
	inputOrder = types.Order{Id: 10, Creator: MockAccount("1"), Amount: 20, Price: 30}
	expected := liquidateBuyRes{
		Book: []types.Order{
			{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
			{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
			{Id: 2, Creator: MockAccount("2"), Amount: 10, Price: 15},
		},
		Remaining:  types.Order{Id: 10, Creator: MockAccount("1"), Amount: 0, Price: 30},
		Liquidated: types.Order{Id: 2, Creator: MockAccount("2"), Amount: 20, Price: 15},
		Purchase:   int32(20),
		Match:      true,
		Filled:     true,
	}
	simulateLiquidateFromBuyOrder(t, inputBook, inputOrder, expected)

	// Entirely filled (30 = 30)
	inputOrder = types.Order{Id: 10, Creator: MockAccount("1"), Amount: 30, Price: 30}
	expected = liquidateBuyRes{
		Book: []types.Order{
			{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
			{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
		},
		Remaining:  types.Order{Id: 10, Creator: MockAccount("1"), Amount: 0, Price: 30},
		Liquidated: types.Order{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
		Purchase:   int32(30),
		Match:      true,
		Filled:     true,
	}
	simulateLiquidateFromBuyOrder(t, inputBook, inputOrder, expected)

	// Not filled and entirely liquidated (60 > 30)
	inputOrder = types.Order{Id: 10, Creator: MockAccount("1"), Amount: 60, Price: 30}
	expected = liquidateBuyRes{
		Book: []types.Order{
			{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
			{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
		},
		Remaining:  types.Order{Id: 10, Creator: MockAccount("1"), Amount: 30, Price: 30},
		Liquidated: types.Order{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
		Purchase:   int32(30),
		Match:      true,
		Filled:     false,
	}
	simulateLiquidateFromBuyOrder(t, inputBook, inputOrder, expected)
}

type fillBuyRes struct {
	Book       []types.Order
	Remaining  types.Order
	Liquidated []types.Order
	Purchase   int32
	Filled     bool
}

func simulateFillBuyOrder(
	t *testing.T,
	inputList []types.Order,
	inputOrder types.Order,
	expected fillBuyRes,
) {
	book := OrderListToSellOrderBook(inputList)
	expectedBook := OrderListToSellOrderBook(expected.Book)

	require.True(t, sort.SliceIsSorted(book.Book.Orders, func(i, j int) bool {
		return book.Book.Orders[i].Price > book.Book.Orders[j].Price
	}))
	require.True(t, sort.SliceIsSorted(expectedBook.Book.Orders, func(i, j int) bool {
		return expectedBook.Book.Orders[i].Price > expectedBook.Book.Orders[j].Price
	}))

	remaining, liquidated, purchase, filled := book.FillBuyOrder(inputOrder)

	require.Equal(t, expectedBook, book)
	require.Equal(t, expected.Remaining, remaining)
	require.Equal(t, expected.Liquidated, liquidated)
	require.Equal(t, expected.Purchase, purchase)
	require.Equal(t, expected.Filled, filled)
}

func TestFillBuyOrder(t *testing.T) {
	var inputBook []types.Order

	// Empty book
	inputOrder := types.Order{Id: 10, Creator: MockAccount("1"), Amount: 30, Price: 10}
	expected := fillBuyRes{
		Book:       []types.Order{},
		Remaining:  inputOrder,
		Liquidated: []types.Order(nil),
		Purchase:   int32(0),
		Filled:     false,
	}
	simulateFillBuyOrder(t, inputBook, inputOrder, expected)

	// No match
	inputBook = []types.Order{
		{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
		{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
		{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
	}
	expected = fillBuyRes{
		Book:       inputBook,
		Remaining:  inputOrder,
		Liquidated: []types.Order(nil),
		Purchase:   int32(0),
		Filled:     false,
	}
	simulateFillBuyOrder(t, inputBook, inputOrder, expected)

	// First order liquidated, not filled
	inputOrder = types.Order{Id: 10, Creator: MockAccount("1"), Amount: 60, Price: 18}
	expected = fillBuyRes{
		Book: []types.Order{
			{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
			{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
		},
		Remaining: types.Order{Id: 10, Creator: MockAccount("1"), Amount: 30, Price: 18},
		Liquidated: []types.Order{
			{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
		},
		Purchase: int32(30),
		Filled:   false,
	}
	simulateFillBuyOrder(t, inputBook, inputOrder, expected)

	// Filled with two order
	inputOrder = types.Order{Id: 10, Creator: MockAccount("1"), Amount: 60, Price: 22}
	expected = fillBuyRes{
		Book: []types.Order{
			{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
			{Id: 1, Creator: MockAccount("1"), Amount: 170, Price: 20},
		},
		Remaining: types.Order{Id: 10, Creator: MockAccount("1"), Amount: 0, Price: 22},
		Liquidated: []types.Order{
			{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
			{Id: 1, Creator: MockAccount("1"), Amount: 30, Price: 20},
		},
		Purchase: int32(30 + 30),
		Filled:   true,
	}
	simulateFillBuyOrder(t, inputBook, inputOrder, expected)

	// Not filled, sell order book liquidated
	inputOrder = types.Order{Id: 10, Creator: MockAccount("1"), Amount: 300, Price: 30}
	expected = fillBuyRes{
		Book:      []types.Order{},
		Remaining: types.Order{Id: 10, Creator: MockAccount("1"), Amount: 20, Price: 30},
		Liquidated: []types.Order{
			{Id: 2, Creator: MockAccount("2"), Amount: 30, Price: 15},
			{Id: 1, Creator: MockAccount("1"), Amount: 200, Price: 20},
			{Id: 0, Creator: MockAccount("0"), Amount: 50, Price: 25},
		},
		Purchase: int32(30 + 200 + 50),
		Filled:   false,
	}
	simulateFillBuyOrder(t, inputBook, inputOrder, expected)
}
```

## Successful Test Output

When the tests are successful, your output is:

```
ok      interchange/x/dex/types       0.550s
```
