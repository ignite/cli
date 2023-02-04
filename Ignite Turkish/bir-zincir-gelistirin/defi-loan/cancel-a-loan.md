# Cancel a loan

Bir borçlu olarak, artık devam etmek istemiyorsanız, oluşturduğunuz bir krediyi iptal etme seçeneğiniz vardır. Ancak, bu işlem yalnızca kredinin mevcut durumu "talep edildi" olarak işaretlenmişse mümkündür.

Krediyi iptal etmeye karar verirseniz, kredi için teminat olarak tutulan teminat jetonları modül hesabından hesabınıza geri aktarılacaktır. Bu, kredi için başlangıçta koyduğunuz teminat jetonlarına yeniden sahip olacağınız anlamına gelir.

x/loan/keeper/msg\_server\_cancel\_loan.go

```
package keeper

import (
    "context"

    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

    "loan/x/loan/types"
)

func (k msgServer) CancelLoan(goCtx context.Context, msg *types.MsgCancelLoan) (*types.MsgCancelLoanResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    loan, found := k.GetLoan(ctx, msg.Id)
    if !found {
        return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", msg.Id)
    }
    if loan.Borrower != msg.Creator {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Cannot cancel: not the borrower")
    }
    if loan.State != "requested" {
        return nil, sdkerrors.Wrapf(types.ErrWrongLoanState, "%v", loan.State)
    }
    borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
    collateral, _ := sdk.ParseCoinsNormalized(loan.Collateral)
    err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, borrower, collateral)
    if err != nil {
        return nil, err
    }
    loan.State = "cancelled"
    k.SetLoan(ctx, loan)
    return &types.MsgCancelLoanResponse{}, nil
}
```

`CancelLoan` iki argüman alır: `goCtx` adında bir `context.Context` ve `msg` adında bir `types.MsgCancelLoan` işaretçisi. Bir `types.MsgCancelLoanResponse` işaretçisi ve bir hata döndürür.

Fonksiyon, `context.Context` nesnesinden `sdk.Context` öğesini almak için `sdk.UnwrapSDKContext` yöntemini kullanarak başlar. Daha sonra `msgServer` türünün `GetLoan` yöntemini kullanarak msg bağımsız değişkeninin Id alanı tarafından tanımlanan bir krediyi alır. Kredi bulunamazsa, işlev `sdk.Wrap` yöntemiyle sarılmış s`dk.ErrKeyNotFound` hatasını kullanarak bir hata döndürür.

Ardından, fonksiyon msg bağımsız değişkeninin `Creator` alanının kredinin `Borrower` alanıyla aynı olup olmadığını kontrol eder. Aynı değillerse, fonksiyon `sdk.Wrap` yöntemiyle sarılmış `sdk.ErrUnauthorized` hatasını kullanarak bir hata döndürür.

Fonksiyon daha sonra kredinin `State` alanının "`requested`" dizesine eşit olup olmadığını kontrol eder. Değilse, işlev `sdk.Wrapf` yöntemiyle sarılmış `types.ErrWrongLoanState` hatasını kullanarak bir hata döndürür.

Kredi doğru duruma sahipse ve mesajı oluşturan kişi krediyi alan kişiyse, fonksiyon kredinin `Collateral` alanında tutulan teminat paralarını `bankKeeper`'ın `SendCoinsFromModuleToAccount` yöntemini kullanarak krediyi alan kişinin hesabına geri göndermeye devam eder. Fonksiyon daha sonra kredinin `State` alanını "cancelled" dizesine günceller ve `SetLoan` yöntemini kullanarak güncellenmiş krediyi ayarlar. Son olarak, fonksiyon bir `types.MsgCancelLoanResponse` nesnesi ve `nil` hatası döndürür.
