# Importing methods from the Bank keeper

Bir önceki adımda `--dep bank` komutunu kullanarak `ignite scaffold module` ile `loan` modülünü oluşturdunuz. Bu komut yeni bir modül oluşturdu ve `bank` keeper'ını kredi modülüne ekledi, bu da bankanın keeper yöntemlerini kredinin keeper yöntemlerine eklemenize ve kullanmanıza olanak tanır.

`--dep bank` tarafından yapılan değişiklikleri görmek için şu dosyaları inceleyin: `x/loan/keeper/keeper.go` ve `x/loan/module.go`.

Ignite `bank` keeper'ını eklemeyi halleder, ancak yine de `loan` modülüne hangi `bank` yöntemlerini kullanacağınızı söylemeniz gerekir. Üç yöntem kullanacaksınız: `SendCoins`, `SendCoinsFromAccountToModule` ve `SendCoinsFromModuleToAccount`. Bunu `BankKeeper` arayüzüne yöntem imzaları ekleyerek yapabilirsiniz:

x/loan/types/expected\_keepers.go

```
package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface {
    SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
    SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
    SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
    SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}
```
