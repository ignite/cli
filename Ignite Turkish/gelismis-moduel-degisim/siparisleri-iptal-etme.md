# Siparişleri İptal Etme

Emir defterlerini, alış ve satış emirlerini uyguladınız. Bu bölümde, alış ve satış emirlerinin iptal edilmesini etkinleştireceksiniz.

### Satış Emrini İptal Etme

Bir satış emrini iptal etmek için, belirli satış emrinin kimliğini almanız gerekir. Daha sonra `RemoveOrderFromID` fonksiyonunu kullanarak belirli bir emri emir defterinden kaldırabilir ve saklayıcıyı buna göre güncelleyebilirsiniz.

Keeper dizinine gidin ve `x/dex/keeper/msg_server_cancel_sell_order.go` dosyasını düzenleyin:

```
// x/dex/keeper/msg_server_cancel_sell_order.go

package keeper

import (
    "context"
    "errors"

    sdk "github.com/cosmos/cosmos-sdk/types"

    "interchange/x/dex/types"
)

func (k msgServer) CancelSellOrder(goCtx context.Context, msg *types.MsgCancelSellOrder) (*types.MsgCancelSellOrderResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Retrieve the book
    pairIndex := types.OrderBookIndex(msg.Port, msg.Channel, msg.AmountDenom, msg.PriceDenom)
    s, found := k.GetSellOrderBook(ctx, pairIndex)
    if !found {
        return &types.MsgCancelSellOrderResponse{}, errors.New("the pair doesn't exist")
    }

    // Check order creator
    order, err := s.Book.GetOrderFromID(msg.OrderID)
    if err != nil {
        return &types.MsgCancelSellOrderResponse{}, err
    }

    if order.Creator != msg.Creator {
        return &types.MsgCancelSellOrderResponse{}, errors.New("canceller must be creator")
    }

    // Remove order
    if err := s.Book.RemoveOrderFromID(msg.OrderID); err != nil {
        return &types.MsgCancelSellOrderResponse{}, err
    }

    k.SetSellOrderBook(ctx, s)

    // Refund seller with remaining amount
    seller, err := sdk.AccAddressFromBech32(order.Creator)
    if err != nil {
        return &types.MsgCancelSellOrderResponse{}, err
    }

    if err := k.SafeMint(ctx, msg.Port, msg.Channel, seller, msg.AmountDenom, order.Amount); err != nil {
        return &types.MsgCancelSellOrderResponse{}, err
    }

    return &types.MsgCancelSellOrderResponse{}, nil
}
```

`GetOrderFromID` İşlevini Uygulayın

`GetOrderFromID` fonksiyonu kitabın ID'sinden bir sipariş alır.

Bu fonksiyonu `types` dizinindeki `x/dex/types/order_book.go` fonksiyonuna ekleyin:

```
// x/dex/types/order_book.go

func (book OrderBook) GetOrderFromID(id int32) (Order, error) {
    for _, order := range book.Orders {
        if order.Id == id {
            return *order, nil
        }
    }

    return Order{}, ErrOrderNotFound
}
```

#### RemoveOrderFromID İşlevini Uygulayın

RemoveOrderFromID işlevi bir siparişi defterden kaldırır ve sıralı olarak tutar:

```
// x/dex/types/order_book.go

package types

// ...

func (book *OrderBook) RemoveOrderFromID(id int32) error {
    for i, order := range book.Orders {
        if order.Id == id {
            book.Orders = append(book.Orders[:i], book.Orders[i+1:]...)
            return nil
        }
    }

    return ErrOrderNotFound
}
```

### Satın Alma Emrini İptal Etme

Bir satın alma emrini iptal etmek için, belirli satın alma emrinin kimliğini almanız gerekir. Ardından, belirli bir emri emir defterinden kaldırmak ve keeper'ı buna göre güncellemek için `RemoveOrderFromID` işlevini kullanabilirsiniz:

```
// x/dex/keeper/msg_server_cancel_buy_order.go

package keeper

import (
    "context"
    "errors"

    sdk "github.com/cosmos/cosmos-sdk/types"

    "interchange/x/dex/types"
)

func (k msgServer) CancelBuyOrder(goCtx context.Context, msg *types.MsgCancelBuyOrder) (*types.MsgCancelBuyOrderResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Retrieve the book
    pairIndex := types.OrderBookIndex(msg.Port, msg.Channel, msg.AmountDenom, msg.PriceDenom)
    b, found := k.GetBuyOrderBook(ctx, pairIndex)
    if !found {
        return &types.MsgCancelBuyOrderResponse{}, errors.New("the pair doesn't exist")
    }

    // Check order creator
    order, err := b.Book.GetOrderFromID(msg.OrderID)
    if err != nil {
        return &types.MsgCancelBuyOrderResponse{}, err
    }

    if order.Creator != msg.Creator {
        return &types.MsgCancelBuyOrderResponse{}, errors.New("canceller must be creator")
    }

    // Remove order
    if err := b.Book.RemoveOrderFromID(msg.OrderID); err != nil {
        return &types.MsgCancelBuyOrderResponse{}, err
    }

    k.SetBuyOrderBook(ctx, b)

    // Refund buyer with remaining price amount
    buyer, err := sdk.AccAddressFromBech32(order.Creator)
    if err != nil {
        return &types.MsgCancelBuyOrderResponse{}, err
    }

    if err := k.SafeMint(
        ctx,
        msg.Port,
        msg.Channel,
        buyer,
        msg.PriceDenom,
        order.Amount*order.Price,
    ); err != nil {
        return &types.MsgCancelBuyOrderResponse{}, err
    }

    return &types.MsgCancelBuyOrderResponse{}, nil
}
```

### Özet

dex modülü için gerekli olan fonksiyonları uygulamayı tamamladınız. Bu bölümde, belirli alım veya satım emirlerini iptal etme tasarımını uyguladınız.

Ignite CLI blockchain'inizin doğru şekilde oluşturulup oluşturulmadığını test etmek için chain build komutunu kullanın:

```
ignite chain build
```

Yine, durumunuzu yerel GitHub deposuna eklemek için iyi bir zaman (harika bir zaman!):

```
git add .
git commit -m "Add Cancelling Orders"
```

Son olarak, şimdi test dosyalarını yazma zamanı.
