# Sipariş Defterlerinin Uygulanması

Bu bölümde, sipariş defterleri oluşturmak için mantığı uygulayacaksınız.

Cosmos SDK'da durum, bir anahtar-değer deposunda saklanır. Her sipariş defteri, dört değerden oluşan benzersiz bir anahtar altında saklanır:

* Port ID
* Channel ID
* Source denom
* Target denom

Örneğin, marscoin ve venuscoin için bir emir defteri `dex-channel-4-marscoin-venuscoin` altında saklanabilir.

İlk olarak, bir emir defteri deposu anahtarı döndüren bir fonksiyon tanımlayın:

```
// x/dex/types/keys.go
package types

import "fmt"

// ...
func OrderBookIndex(portID string, channelID string, sourceDenom string, targetDenom string) string {
    return fmt.Sprintf("%s-%s-%s-%s", portID, channelID, sourceDenom, targetDenom)
}
```

Send-create-pair komutu sipariş defterleri oluşturmak için kullanılır. Bu komut:

* `SendCreatePair` türünde bir mesajla bir işlem oluşturur ve yayınlar.
* Mesaj `dex` modülüne yönlendirilir.
* Son olarak, bir `SendCreatePair` keeper yöntemi çağrılır.

Aşağıdakileri yapmak için send-create-pair komutuna ihtiyacınız vardır:

* Kaynak zincirde `SendCreatePair` mesajı işlenirken:

1. Verilen denom çiftine sahip bir emir defterinin henüz mevcut olup olmadığını kontrol edin.
2. Bağlantı noktası, kanal, kaynak denomlar ve hedef denomlar hakkında bilgi içeren bir IBC paketi iletin.

* Paket hedef zincirde alındıktan sonra:

1. Verilen denom çiftine sahip bir emir defterinin hedef zincirde henüz mevcut olup olmadığını kontrol edin.
2. Alış emirleri için yeni bir emir defteri oluşturun.
3. Kaynak zincire bir IBC onayını geri iletin.

* Kaynak zincirde onay alındıktan sonra:

1. Satış emirleri için yeni bir emir defteri oluşturun.

### `SendCreatePair'`de Mesaj İşleme

`SendCreatePair` fonksiyonu IBC paket iskelesi sırasında oluşturulmuştur. Fonksiyon bir IBC paketi oluşturur, bu paketi kaynak ve hedef denomlarla doldurur ve bu paketi IBC üzerinden iletir.

Şimdi, belirli bir denom çifti için mevcut bir sipariş defteri olup olmadığını kontrol etme mantığını ekleyin:

```
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

### Bir IBC Paketinin Yaşam Döngüsü

Başarılı bir iletim sırasında bir IBC paketi şu aşamalardan geçer:

1. Kaynak zincirinde paket iletiminden önce mesaj işleme
2. Hedef zincirde bir paketin alınması
3. Kaynak zincirindeki bir paketin onaylanması
4. Kaynak zincirindeki bir paketin zaman aşımı

Aşağıdaki bölümde, `OnRecvCreatePairPacket` işlevinde paket alım mantığını ve `OnAcknowledgementCreatePairPacket` işlevinde paket onay mantığını uygulayın.

Zaman Aşımı işlevini boş bırakın.

### Bir IBC paketi alın

Protokol tampon tanımı, bir sipariş defterinin içerdiği verileri tanımlar.

`OrderBook` ve `Order` mesajlarını `order.proto` dosyasına ekleyin.

İlk olarak, Go kod dosyalarını oluşturmak için proto tampon dosyalarını ekleyin. Bu dosyaları uygulamanızın amacına uygun olarak değiştirebilirsiniz.

`proto/interchange/dex` dizininde yeni bir `order.proto` dosyası oluşturun ve içeriği ekleyin:

```
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

`buy_order_book.proto` dosyasını, sipariş defterinde bir satın alma emri oluşturmak için alanlara sahip olacak şekilde değiştirin. İçe aktarmayı da eklemeyi unutmayın.

**İpucu**: İçe aktarmayı da eklemeyi unutmayın.

```
// proto/interchange/dex/buy_order_book.proto

// ...

import "interchange/dex/order.proto";

message BuyOrderBook {
  // ...
  OrderBook book = 4;
}
```

Sipariş defterini satın alma sipariş defterine eklemek için `sell_order_book.proto` dosyasını değiştirin.

`SellOrderBook` için proto tanımı şöyle görünür:

```
// proto/interchange/dex/sell_order_book.proto

// ...
import "interchange/dex/order.proto";

message SellOrderBook {
  // ...
  OrderBook book = 4;
}
```

Şimdi, `send-create-pair` komutu için proto dosyalarını oluşturmak üzere Ignite CLI kullanın:

```
ignite generate proto-go --yes
```

IBC paketleri için işlevleri geliştirmeye başlayın.

Yeni bir `x/dex/types/order_book.go` dosyası oluşturun.

Yeni sipariş defteri fonksiyonunu ilgili Go dosyasına ekleyin:

```
// x/dex/types/order_book.go

package types

func NewOrderBook() OrderBook {
    return OrderBook{
        IdCount: 0,
    }
}
```

Yeni bir alış emri defteri türü oluşturmak için, `x/dex/types/buy_order_book.go` dosyasında `NewBuyOrderBook`'u tanımlayın:

```
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

Hedef zincirde bir IBC paketi alındığında, modül zaten bir defterin mevcut olup olmadığını kontrol etmelidir. Eğer yoksa, belirtilen denomlar için bir alış emri defteri oluşturmalıdır.

```
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

### Bir IBC Onayı Alın

Kaynak zincirinde bir IBC onayı alındığında, modül bir defterin zaten mevcut olup olmadığını kontrol etmelidir. Eğer yoksa, belirtilen denomlar için bir satış emri defteri oluşturmalıdır.

Yeni bir `x/dex/types/sell_order_book.go` dosyası oluşturun. Yeni bir satış emri defteri oluşturan `NewSellOrderBook` fonksiyonunu ekleyin.

```
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

`x/dex/keeper/create_pair.go` dosyasındaki `Acknowledgement` işlevini değiştirin:

```
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

Bu bölümde, yeni `send-create-pair` komutunun arkasındaki mantığı uyguladınız:

Hedef zincirde bir IBC paketi alındığında, `send-create-pair` komutu bir alış emri defteri oluşturur.

Kaynak zincirde bir IBC onayı alındığında, `send-create-pair` komutu bir satış emri defteri oluşturur.

#### Emir Defterine Emir Eklemek için appendOrder Fonksiyonunu Uygulayın

```
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

`AppendOrder` fonksiyonu, emir bilgilerinden yeni bir emri başlatır ve emir defterine ekler:

```
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

Bir Sipariş İçin `checkAmountAndPrice` İşlevini Uygulama

`checkAmountAndPrice` fonksiyonu doğru miktar veya fiyat olup olmadığını kontrol eder:

```
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

#### `GetNextOrderID` İşlevini Uygulayın

`GetNextOrderID` fonksiyonu eklenecek bir sonraki siparişin ID'sini alır:

```
// x/dex/types/order_book.go

func (book OrderBook) GetNextOrderID() int32 {
    return book.IdCount
}
```

#### IncrementNextOrderID İşlevini Uygulama

IncrementNextOrderID fonksiyonu, siparişler için ID sayısını günceller:

```
// x/dex/types/order_book.go

func (book *OrderBook) IncrementNextOrderID() {
    // Even numbers to have different ID than buy orders
    book.IdCount++
}
```

#### `insertOrder` İşlevini Uygulama

`insertOrder` işlevi, siparişi sağlanan siparişle birlikte deftere ekler:

```
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

Bu, sipariş defteri kurulumunu tamamlar.

Şimdi uygulamanızın durumunu kaydetmek için iyi bir zaman. Projeniz yerel bir depoda olduğu için git kullanabilirsiniz. Mevcut durumunuzu kaydetmek, hata yapmanız veya ara vermeniz gerektiğinde ileri geri atlamanıza olanak tanır.

```
git add .
git commit -m "Create Order Books"
```

Bir sonraki bölümde, kupon basarak ve yakarak ve uygulamanızda yerel blok zinciri tokenini kilitleyerek ve kilidini açarak kuponlarla nasıl başa çıkacağınızı öğreneceksiniz.
