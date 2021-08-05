---
order: 6
---

# Creating Sell Orders

In this chapter you will implement the logic for creating sell orders.

The packet proto file for a sell order is already generated. Add the seller information.

```proto
// proto/ibcdex/packet.proto
message SellOrderPacketData {
  // ...
  string seller = 5;
}
```

## `SendSellOrder` Message Handling

Sell orders are created using `send-sellOrder`. This command creates a transaction with a `SendSellOrder` message, which triggers the `SendSellOrder` keeper method.

`SendSellOrder` should:

* Check that an order book for specified denom pair exists
* Safely burn or lock tokens
  * If the token is an IBC token, burn the tokens
  * If the token is a native token, lock the tokens
* Save the voucher received on the target chain to later resolve a denom
* Transmit an IBC packet to the target chain

```go
// x/ibcdex/keeper/msg_server_sellOrder.go
import "errors"

func (k msgServer) SendSellOrder(goCtx context.Context, msg *types.MsgSendSellOrder) (*types.MsgSendSourceSellOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// If an order book doesn't exist, throw an error
	pairIndex := types.OrderBookIndex(msg.Port, msg.ChannelID, msg.AmountDenom, msg.PriceDenom)
	_, found := k.GetSellOrderBook(ctx, pairIndex)
	if !found {
		return &types.MsgSendSellOrderResponse{}, errors.New("the pair doesn't exist")
	}
	// Get sender's address
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
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
	packet.Seller = msg.Sender
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

`SendSellOrder` depends on two new keeper methods: `SafeBurn` and `SaveVoucherDenom`.

### `SafeBurn`

`SafeBurn` burns tokens if they are IBC vouchers (have an `ibc/` prefix) and locks tokens if they are native to the chain.

```go
// x/ibcdex/keeper/mint.go
package keeper

import (
  "fmt"
  sdk "github.com/cosmos/cosmos-sdk/types"
  ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
  "github.com/username/interchange/x/ibcdex/types"
  "strings"
)

// isIBCToken checks if the token came from the IBC module
func isIBCToken(denom string) bool {
  return strings.HasPrefix(denom, "ibc/")
}

func (k Keeper) SafeBurn(ctx sdk.Context, port string, channel string, sender sdk.AccAddress, denom string, amount int32) error {
  if isIBCToken(denom) {
    // Burn the tokens
    if err := k.BurnTokens(ctx, sender, sdk.NewCoin(denom, sdk.NewInt(int64(amount)))); err != nil {
      return err
    }
  } else {
    // Lock the tokens
    if err := k.LockTokens(ctx, port, channel, sender, sdk.NewCoin(denom, sdk.NewInt(int64(amount)))); err != nil {
      return err
    }
  }
  return nil
}
```

Implement the `BurnTokens` keeper method.

```go
// x/ibcdex/keeper/mint.go
func (k Keeper) BurnTokens(ctx sdk.Context, sender sdk.AccAddress, tokens sdk.Coin) error {
  // transfer the coins to the module account and burn them
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(tokens)); err != nil {
		return err
	}
  if err := k.bankKeeper.BurnCoins(
    ctx, types.ModuleName, sdk.NewCoins(tokens),
  ); err != nil {
    // NOTE: should not happen as the module account was
    // retrieved on the step above and it has enough balance
    // to burn.
    panic(fmt.Sprintf("cannot burn coins after a successful send to a module account: %v", err))
  }

  return nil
}
```

Implement the `LockTokens` keeper method.

```go
// x/ibcdex/keeper/mint.go
func (k Keeper) LockTokens(ctx sdk.Context, sourcePort string, sourceChannel string, sender sdk.AccAddress, tokens sdk.Coin) error {
  // create the escrow address for the tokens
  escrowAddress := ibctransfertypes.GetEscrowAddress(sourcePort, sourceChannel)
  // escrow source tokens. It fails if balance insufficient
  if err := k.bankKeeper.SendCoins(
    ctx, sender, escrowAddress, sdk.NewCoins(tokens),
  ); err != nil {
    return err
  }
  return nil
}
```

`BurnTokens` and `LockTokens` use `SendCoinsFromAccountToModule`, `BurnCoins`, and `SendCoins` keeper methods of the `bank` module. To start using these function from the `ibcdex` module, first add them to the `BankKeeper` interface.

```go
// x/ibcdex/types/expected_keeper.go
package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
  SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
  BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
  SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
```

Next, in the `keeper` directory, specify the bank so that you can access it in your module. 

```go
// x/ibcdex/keeper/keeper.go
type (
  Keeper struct {
    // ...
    bankKeeper types.BankKeeper
  }
)

func NewKeeper(
  // ...
  bankKeeper types.BankKeeper,
) *Keeper {
  return &Keeper{
    // ...
    bankKeeper: bankKeeper,
  }
}
```

Lastly, the `app.go` file that describes which modules are used in the blockchain application, add the bank keeper to the `ibcdexKeeper`

```go
// app/app.go
app.ibcdexKeeper = *ibcdexkeeper.NewKeeper(
  // ...
  app.BankKeeper,
)
```

The `ibcdex` module will need to mint and burn token using the `bank` account. The use this feature, the module must have a _module account_. To enable the _module account_ declare this permission in the _module account permissions_ structure of the auth module.

```go
// app/app.go
maccPerms = map[string][]string{
    // ...
    ibcdextypes.ModuleName: {authtypes.Minter, authtypes.Burner},
}
```

### `SaveVoucherDenom`

`SaveVoucherDenom` saves the voucher denom to be able to convert it back later.

```go
// x/ibcdex/keeper/denom.go
func (k Keeper) SaveVoucherDenom(ctx sdk.Context, port string, channel string, denom string) {
	voucher := VoucherDenom(port, channel, denom)

	// Store the origin denom
	_, saved := k.GetDenomTrace(ctx, voucher)
	if !saved {
		k.SetDenomTrace(ctx, types.DenomTrace{
			Index:   voucher,
			Port:    port,
			Channel: channel,
			Origin:  denom,
		})
	}
}
```

Finally, last function we need to implement is `VoucherDenom`. `VoucherDenom` returns the voucher of the denom from the port ID and channel ID.

```go
// x/ibcdex/keeper/denom.go
import (
  sdk "github.com/cosmos/cosmos-sdk/types"
  ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
  "github.com/username/interchange/x/ibcdex/types"
)

func VoucherDenom(port string, channel string, denom string) string {
  // since SendPacket did not prefix the denomination, we must prefix denomination here
  sourcePrefix := ibctransfertypes.GetDenomPrefix(port, channel)
  // NOTE: sourcePrefix contains the trailing "/"
  prefixedDenom := sourcePrefix + denom
  // construct the denomination trace from the full raw denomination
  denomTrace := ibctransfertypes.ParseDenomTrace(prefixedDenom)
  voucher := denomTrace.IBCDenom()
  return voucher[:16]
}
```

## `OnRecv`

When a "sell order" packet is received on the target chain, the module should  ????

- Update the sell order book
- Distribute sold token to the buyer
- Send to chain A the sell order after the fill attempt

```go
// x/ibcdex/keeper/sellOrder.go
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

### `FillSellOrder`

`FillSellOrder` try to fill the sell order with the order book and returns all the side effects.

```go
// x/ibcdex/types/buy_order_book.go
func (b *BuyOrderBook) FillSellOrder(order Order) (remainingSellOrder Order, liquidated []Order, gain int32, filled bool) {
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

#### `LiquidateFromSellOrder`

`LiquidateFromSellOrder` liquidates the first buy order of the book from the sell order if no match is found, return false for match.

```go
// x/ibcdex/types/buy_order_book.go
func (b *BuyOrderBook) LiquidateFromSellOrder(order Order) ( remainingSellOrder Order, liquidatedBuyOrder Order, gain int32, match bool, filled bool) {
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
```

### `OriginalDenom`

`OriginalDenom` returns back the original denom of the voucher. False is returned if the port ID and channel ID provided are not the origins of the voucher

```go
// x/ibcdex/keeper/denom.go
func (k Keeper) OriginalDenom(ctx sdk.Context, port string, channel string, voucher string) (string, bool) {
	trace, exist := k.GetDenomTrace(ctx, voucher)
	if exist {
		// Check if original port and channel
		if trace.Port == port && trace.Channel == channel {
			return trace.Origin, true
		}
	}
	// Not the original chain
	return "", false
}
```

If a token is an IBC token (has an `ibc/` prefix) `SafeMint` mints IBC tokens with `MintTokens`, otherwise, it unlocks native tokens with `UnlockTokens`.

### `SafeMint`

```go
// x/ibcdex/keeper/mint.go
func (k Keeper) SafeMint(ctx sdk.Context, port string, channel string, receiver sdk.AccAddress, denom string, amount int32) error {
	if isIBCToken(denom) {
		// Mint IBC tokens
		if err := k.MintTokens(ctx, receiver, sdk.NewCoin(denom, sdk.NewInt(int64(amount)))); err != nil {
			return err
		}
	} else {
		// Unlock native tokens
		if err := k.UnlockTokens(
			ctx,
			port,
			channel,
			receiver,
			sdk.NewCoin(denom, sdk.NewInt(int64(amount))),
		); err != nil {
			return err
		}
	}
	return nil
}
```

#### `MintTokens`

```go
// x/ibcdex/keeper/mint.go
func (k Keeper) MintTokens(ctx sdk.Context, receiver sdk.AccAddress, tokens sdk.Coin) error {
	// mint new tokens if the source of the transfer is the same chain
	if err := k.bankKeeper.MintCoins(
		ctx, types.ModuleName, sdk.NewCoins(tokens),
	); err != nil {
		return err
	}
	// send to receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, receiver, sdk.NewCoins(tokens),
	); err != nil {
		panic(fmt.Sprintf("unable to send coins from module to account despite previously minting coins to module account: %v", err))
	}
	return nil
}
```

`MintTokens` uses two keeper methods from the `bank` module: `MintCoins` and `SendCoinsFromModuleToAccount`. Import them by adding their signatures to the `BankKeeper` interface.

```go
// x/ibcdex/types/expected_keeper.go
type BankKeeper interface {
  // ...
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}
```

```go
// x/ibcdex/keeper/mint.go
func (k Keeper) UnlockTokens(ctx sdk.Context, sourcePort string, sourceChannel string, receiver sdk.AccAddress, tokens sdk.Coin) error {
	// create the escrow address for the tokens
	escrowAddress := ibctransfertypes.GetEscrowAddress(sourcePort, sourceChannel)
	// escrow source tokens. It fails if balance insufficient
	if err := k.bankKeeper.SendCoins(
		ctx, escrowAddress, receiver, sdk.NewCoins(tokens),
	); err != nil {
		return err
	}
	return nil
}
```

## `OnAcknowledgement`

Once an IBC packet is processed on the target chain, an acknowledgement is returned to the source chain and processed in `OnAcknowledgementSellOrderPacket`. The module on the source chain will store the remaining sell order in the sell order book and will distribute sold tokens to the buyers and will distribute to the seller the price of the amount sold. On error the module mints the burned tokens.

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
// x/ibcdex/types/sell_order_book.go
func (s *SellOrderBook) AppendOrder(creator string, amount int32, price int32) (int32, error) {
	return s.Book.appendOrder(creator, amount, price, Decreasing)
}
```

### `appendOrder` Implementation

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
```

`AppendOrder` initializes and appends a new order to an order book from the order information.

```go
// x/ibcdex/types/order_book.go
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

#### `checkAmountAndPrice`

`checkAmountAndPrice` checks correct amount or price.

```go
// x/ibcdex/types/order_book.go
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

#### `GetNextOrderID`

`GetNextOrderID` gets the ID of the next order to append

```go
// x/ibcdex/types/order_book.go
func (book OrderBook) GetNextOrderID() int32 {
	return book.IdCount
}
```

#### `IncrementNextOrderID`

`IncrementNextOrderID` updates the ID count for orders

```go
// x/ibcdex/types/order_book.go
func (book *OrderBook) IncrementNextOrderID() {
	// Even numbers to have different ID than buy orders
	book.IdCount++
}
```

#### `insertOrder`

`insertOrder` inserts the order in the book with the provided order

```go
// x/ibcdex/types/order_book.go
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

## `OnTimeout`

If a timeout occurs, we mint back the native token.

```go
// x/ibcdex/keeper/sellOrder.go
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