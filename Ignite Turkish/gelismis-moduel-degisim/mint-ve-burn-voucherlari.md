# Mint ve Burn Voucher'ları

Bu bölümde, voucher'lar hakkında bilgi edineceksiniz. `Dex` modülü uygulaması voucher'ları basar ve bir blockchain'den yerel tokenları kilitler ve kilidini açar.

Bu `dex` modülü uygulamasından öğrenilecek çok şey var:

* `Bank` saklayıcısı ile çalışır ve sunduğu çeşitli yöntemleri kullanırsınız.
* Başka bir modülle etkileşime giriyor ve tokenları kilitlemek için modül hesabını kullanıyorsunuz.

Bu uygulama size modül hesaplarıyla çeşitli etkileşimleri nasıl kullanacağınızı veya token basmayı, kilitlemeyi veya yakmayı öğretebilir.

### Burn Voucher'lar veya Lock Token'lar için SafeBurn Fonksiyonunu Oluşturun

`SafeBurn` işlevi, IBC kuponuysa (`ibc/` ön ekine sahipse) jetonları yakar ve zincire özgü ise jetonları kilitler.

Yeni bir `x/dex/keeper/mint.go` dosyası oluşturun:

```
// x/dex/keeper/mint.go

package keeper

import (
    "fmt"
    "strings"

    sdkmath "cosmossdk.io/math"
    sdk "github.com/cosmos/cosmos-sdk/types"
    ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

    "interchange/x/dex/types"
)

// isIBCToken checks if the token came from the IBC module
// Each IBC token starts with an ibc/ denom, the check is rather simple
func isIBCToken(denom string) bool {
    return strings.HasPrefix(denom, "ibc/")
}

func (k Keeper) SafeBurn(ctx sdk.Context, port string, channel string, sender sdk.AccAddress, denom string, amount int32) error {
    if isIBCToken(denom) {
        // Burn the tokens
        if err := k.BurnTokens(ctx, sender, sdk.NewCoin(denom, sdkmath.NewInt(int64(amount)))); err != nil {
            return err
        }
    } else {
        // Lock the tokens
        if err := k.LockTokens(ctx, port, channel, sender, sdk.NewCoin(denom, sdkmath.NewInt(int64(amount)))); err != nil {
            return err
        }
    }

    return nil
}
```

Token başka bir blok zincirinden IBC tokeni olarak geliyorsa, yakma yöntemi aslında bu IBC tokenlerini bir zincirde yakar ve diğer zincirde kilidini açar. Yerel tokenlar kilitli kalır.

Şimdi, `BurnTokens` keeper yöntemini önceki işlevde kullanıldığı gibi uygulayın. `BankKeeper` bunun için kullanışlı bir işleve sahiptir:

```
// x/dex/keeper/mint.go

package keeper

// ...

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

`LockTokens` keeper yöntemini uygulayın.

Yerel bir zincirden token kilitlemek için, yerel token'ı Escrow Address'e gönderebilirsiniz:

```
// x/dex/keeper/mint.go

package keeper

// ...

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

`BurnTokens` ve `LockTokens`, `bank` modülünün `SendCoinsFromAccountToModule`, `BurnCoins` ve `SendCoins` keeper yöntemlerini kullanır.

Bu işlevleri `dex` modülünden kullanmaya başlamak için, önce bunları `x/dex/types/expected_keepers.go` dosyasındaki `BankKeeper` arayüzüne ekleyin.

```
// x/dex/types/expected_keepers.go

package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
    //...
    SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
    BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
    SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
```

### SaveVoucherDenom

`SaveVoucherDenom` işlevi, daha sonra geri dönüştürebilmek için voucher denom'unu kaydeder.

Yeni bir `x/dex/keeper/denom.go` dosyası oluşturun:

```
// x/dex/keeper/denom.go

package keeper

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

    "interchange/x/dex/types"
)

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

Son olarak, uygulanacak son fonksiyon, port ID ve kanal ID'sinden denomun voucher'ını döndüren `VoucherDenom` fonksiyonudur:

```
// x/dex/keeper/denom.go

package keeper

// ...

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

#### Bir OriginalDenom İşlevi Uygulayın

`OriginalDenom` fonksiyonu voucher'ın orijinal denom'unu geri döndürür.

Sağlanan port ID'si ve kanal ID'si voucher'ın orijini değilse False döndürülür:

```
// x/dex/keeper/denom.go

package keeper

// ...

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

#### Bir SafeMint İşlevi Uygulama

Bir token bir IBC token ise (ibc/ ön ekine sahipse), SafeMint işlevi MintTokens ile IBC token basar. Aksi takdirde, UnlockTokens ile yerel token kilidini açar.

x/dex/keeper/mint.go dosyasına geri dönün ve aşağıdaki kodu ekleyin:

```
// x/dex/keeper/mint.go

package keeper

// ...

func (k Keeper) SafeMint(ctx sdk.Context, port string, channel string, receiver sdk.AccAddress, denom string, amount int32) error {
    if isIBCToken(denom) {
        // Mint IBC tokens
        if err := k.MintTokens(ctx, receiver, sdk.NewCoin(denom, sdkmath.NewInt(int64(amount)))); err != nil {
            return err
        }
    } else {
        // Unlock native tokens
        if err := k.UnlockTokens(
            ctx,
            port,
            channel,
            receiver,
            sdk.NewCoin(denom, sdkmath.NewInt(int64(amount))),
        ); err != nil {
            return err
        }
    }

    return nil
}
```

#### Bir MintTokens İşlevi Uygulama

MintCoins için bankKeeper işlevini tekrar kullanabilirsiniz. Bu tokenlar daha sonra alıcı hesabına gönderilecektir:

```
// x/dex/keeper/mint.go

package keeper

// ...

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

Son olarak, yerel blockchain'e geri gönderildikten sonra token kilidini açma işlevini ekleyin:

```
// x/dex/keeper/mint.go

package keeper

// ...

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

`MintTokens` fonksiyonu `bank` modülünden iki keeper metodu kullanır: `MintCoins` ve `SendCoinsFromModuleToAccount` . Bu yöntemleri içe aktarmak için imzalarını `x/dex/types/expected_keepers.go` dosyasındaki `BankKeeper` arayüzüne ekleyin:

```
// x/dex/types/expected_keepers.go

package types

// ...

type BankKeeper interface {
    // ...
    MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
    SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}
```

### Özet

Mint ve burn voucher mantığını bitirdiniz.

Çalışmanızın durumunu kaydetmek için başka bir git commit'i yapmanın tam zamanı:

```
git add .
git commit -m "Add Mint and Burn Voucher"
```

Bir sonraki bölümde, satış emirleri oluşturmayı inceleyeceksiniz.
