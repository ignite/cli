# Kredi tasfiyesi

`LiquidateLoan` yöntemi, borçlunun krediyi belirtilen son tarihe kadar geri ödeyememesi durumunda borç verenin borçluya ait teminatı satmasına olanak tanıyan bir işlevdir. Bu süreç "tasfiye" olarak bilinir ve genellikle borçlunun geri ödeme yükümlülüklerini yerine getirememesi durumunda borç verenin zararlarını telafi etmesinin bir yolu olarak gerçekleştirilir.

Tasfiye işlemi sırasında, borçlu tarafından kredinin teminatı olarak rehin verilen teminat tokenleri borçlunun hesabından borç verenin hesabına aktarılır. Bu transfer borç veren tarafından başlatılır ve tipik olarak borçlunun üzerinde anlaşılan son tarihe kadar krediyi geri ödeyememesi durumunda tetiklenir. Teminat transfer edildikten sonra, borç veren zararlarını telafi etmek ve ödenmemiş krediyi tazmin etmek için teminatı satabilir.

### Keeper metot

x/loan/keeper/msg\_server\_liquidate\_loan.go

```
package keeper

import (
    "context"
    "fmt"
    "strconv"

    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

    "loan/x/loan/types"
)

func (k msgServer) LiquidateLoan(goCtx context.Context, msg *types.MsgLiquidateLoan) (*types.MsgLiquidateLoanResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    loan, found := k.GetLoan(ctx, msg.Id)
    if !found {
        return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", msg.Id)
    }
    if loan.Lender != msg.Creator {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Cannot liquidate: not the lender")
    }
    if loan.State != "approved" {
        return nil, sdkerrors.Wrapf(types.ErrWrongLoanState, "%v", loan.State)
    }
    lender, _ := sdk.AccAddressFromBech32(loan.Lender)
    collateral, _ := sdk.ParseCoinsNormalized(loan.Collateral)
    deadline, err := strconv.ParseInt(loan.Deadline, 10, 64)
    if err != nil {
        panic(err)
    }
    if ctx.BlockHeight() < deadline {
        return nil, sdkerrors.Wrap(types.ErrDeadline, "Cannot liquidate before deadline")
    }
    err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, lender, collateral)
    if err != nil {
        return nil, err
    }
    loan.State = "liquidated"
    k.SetLoan(ctx, loan)
    return &types.MsgLiquidateLoanResponse{}, nil
}
```

`LiquidateLoan`, girdi olarak bir context ve bir `types.MsgLiquidateLoan` mesajı alır ve çıktı olarak bir `types.MsgLiquidateLoanResponse` mesajı ve bir hata döndürür.

Fonksiyon ilk olarak `GetLoan` yöntemini ve girdi mesajının `Id` alanını kullanarak bir kredi alır. Kredi bulunamazsa, `sdkerrors.Wrap` işlevini ve `sdkerrors.ErrKeyNotFound` hata kodunu kullanarak bir hata döndürür.

Ardından, işlev giriş mesajının `Creator` alanının kredinin `Lender` alanıyla aynı olup olmadığını kontrol eder. Aynı değillerse, `sdkerrors.Wrap` işlevini ve `sdkerrors.ErrUnauthorized` hata kodunu kullanarak bir hata döndürür.

Fonksiyon daha sonra kredinin Durum alanının "onaylandı" değerine eşit olup olmadığını kontrol eder. Değilse, `sdkerrors.Wrapf` işlevini ve `types.ErrWrongLoanState` hata kodunu kullanarak bir hata döndürür.

Fonksiyon daha sonra `sdk.AccAddressFromBech32` fonksiyonunu kullanarak kredinin Lender alanını bir adrese ve `sdk.ParseCoinsNormalized` fonksiyonunu kullanarak `Collateral` alanını coin'e dönüştürür. Ayrıca `strconv.ParseInt` fonksiyonunu kullanarak `Deadline` alanını bir tamsayıya dönüştürür. Bu fonksiyon bir hata döndürürse panikler.

Son olarak, fonksiyon mevcut blok yüksekliğinin son tarihten büyük veya eşit olup olmadığını kontrol eder. Değilse, `sdkerrors.Wrap` işlevini ve `types.ErrDeadline` hata kodunu kullanarak bir hata döndürür. Tüm kontroller başarılı olursa, işlev teminatı modül hesabından borç verenin hesabına aktarmak için `bankKeeper.SendCoinsFromModuleToAccount` yöntemini kullanır ve kredinin `State` alanını "liquidated" olarak günceller. Ardından `SetLoan` yöntemini kullanarak güncellenmiş krediyi saklar ve hata içermeyen bir `types.MsgLiquidateLoanResponse` mesajı döndürür.

### Özel bir hata kaydedin

x/loan/types/errors.go

```
package types

import (
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
    ErrWrongLoanState = sdkerrors.Register(ModuleName, 2, "wrong loan state")
    ErrDeadline = sdkerrors.Register(ModuleName, 3, "deadline")
)
```
