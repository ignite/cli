# Satış Siparişleri Oluşturun

Bu bölümde, satış emirleri oluşturma mantığını uygulayacaksınız.

Bir satış emri için paket proto dosyası zaten oluşturulmuştur. Satıcı bilgilerini ekleyin:

```
// proto/dex/packet.proto

message SellOrderPacketData {
  // ...
  string seller = 5;
}
```

Şimdi, send-sell-order komutu için proto dosyalarını oluşturmak üzere Ignite CLI kullanın. Bu komutu bir önceki bölümde kullanmıştınız.

```
ignite generate proto-go --yes
```

### SendSellOrder'da Mesaj İşleme

Satış emirleri send-sell-order komutu kullanılarak oluşturulur. Bu komut, SendSellOrder keeper yöntemini tetikleyen bir SendSellOrder mesajı içeren bir işlem oluşturur.

SendSellOrder komutu:

* Belirtilen bir denom çifti için bir sipariş defterinin var olup olmadığını kontrol eder.
* Tokeni güvenli bir şekilde yakar veya kilitler.

1. Token bir IBC token ise, tokenı yakar.
2. Eğer token yerel bir token ise, tokenı kilitler.

* Daha sonra bir denomu çözümlemek için hedef zincirde alınan kuponu kaydeder.
* Hedef zincire bir IBC paketi iletir.

```
// x/dex/keeper/msg_server_sell_order.go

package keeper

import (
    "context"
    "errors"

    sdk "github.com/cosmos/cosmos-sdk/types"
    clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"

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

### Bir Satış Emri Alındığında

Hedef zincirde bir "satış emri" paketi alındığında, modülün şunları yapmasını istersiniz:

* Satış emri defterini güncellemek
* Satılan jetonu alıcıya dağıtın
* Doldurma girişiminden sonra satış emrini A zincirine gönderin

```
// x/dex/keeper/sell_order.go

package keeper

// ...

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

#### FillSellOrder İşlevini Uygulayın

FillSellOrder fonksiyonu, alış emrini emir defteri ile doldurmaya çalışır ve tüm yan etkileri döndürür:

```
// x/dex/types/buy_order_book.go

package types

// ...

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

#### LiquidateFromSellOrder İşlevini Uygulayın

LiquidateFromSellOrder fonksiyonu, alış emrinden defterin ilk satış emrini likide eder. Eşleşme bulunamazsa, eşleşme için false döndürür:

```
// x/dex/types/buy_order_book.go

package types

// ...

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

        // Remove the highest bid if it has been entirely liquidated
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

#### Satış Emri Paketleri için OnAcknowledgement İşlevini Uygulayın

Bir IBC paketi hedef zincirde işlendikten sonra, kaynak zincire bir onay döndürülür ve OnAcknowledgementSellOrderPacket işlevi tarafından işlenir.

Kaynak zincirdeki dex modülü:

* Kalan satış emrini satış emri defterinde saklar.
* Satılan tokenları alıcılara dağıtır.
* Satılan miktarın fiyatını satıcıya dağıtır.
* Hata durumunda, yanmış tokenları basar.

```
// x/dex/keeper/sell_order.go

package keeper

// ...

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

```
// x/dex/types/sell_order_book.go

package types

// ...

func (s *SellOrderBook) AppendOrder(creator string, amount int32, price int32) (int32, error) {
    return s.Book.appendOrder(creator, amount, price, Decreasing)
}
```

#### Satış Emri Paketinin OnTimeout İşlevini Ekleme

Bir zaman aşımı meydana gelirse, yerel belirteci geri nane:

```
// x/dex/keeper/sell_order.go

package keeper

// ...

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

### Özet

Harika, satış emri mantığını tamamladınız.

Çalışmanızın durumunu kaydetmek için tekrar bir git commit'i yapmanın tam zamanı:

```
git add .
git commit -m "Add Sell Orders"
```
