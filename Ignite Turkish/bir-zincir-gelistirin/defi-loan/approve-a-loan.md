# Approve a loan

Bir kredi talebi yapıldıktan sonra, başka bir hesabın krediyi onaylaması ve borçlu tarafından önerilen şartları kabul etmesi mümkündür. Bu süreç, talep edilen fonların borç verenden borç alana aktarılmasını içerir.

Bir kredinin onaylanabilmesi için "talep edildi" statüsünde olması gerekir. Bu, borçlunun bir kredi talebinde bulunduğu ve bir borç verenin şartları kabul etmesini ve fonları sağlamasını beklediği anlamına gelir. Bir borç veren krediyi onaylamaya karar verdikten sonra, fonların borçluya transferini başlatabilir.

Kredi onaylandıktan sonra, kredinin durumu "onaylandı" olarak değiştirilir. Bu, fonların başarıyla transfer edildiğini ve kredi sözleşmesinin artık yürürlükte olduğunu gösterir.

### Keeper Metot

x/loan/keeper/msg\_server\_approve\_loan.go

```
package keeper

import (
    "context"

    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

    "loan/x/loan/types"
)

func (k msgServer) ApproveLoan(goCtx context.Context, msg *types.MsgApproveLoan) (*types.MsgApproveLoanResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    loan, found := k.GetLoan(ctx, msg.Id)
    if !found {
        return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", msg.Id)
    }
    if loan.State != "requested" {
        return nil, sdkerrors.Wrapf(types.ErrWrongLoanState, "%v", loan.State)
    }
    lender, _ := sdk.AccAddressFromBech32(msg.Creator)
    borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
    amount, err := sdk.ParseCoinsNormalized(loan.Amount)
    if err != nil {
        return nil, sdkerrors.Wrap(types.ErrWrongLoanState, "Cannot parse coins in loan amount")
    }
    err = k.bankKeeper.SendCoins(ctx, lender, borrower, amount)
    if err != nil {
        return nil, err
    }
    loan.Lender = msg.Creator
    loan.State = "approved"
    k.SetLoan(ctx, loan)
    return &types.MsgApproveLoanResponse{}, nil
}
```

`ApproveLoan`, girdi olarak bir context ve `*types.MsgApproveLoan` türünde bir mesaj alır ve `types.MsgApproveLoanResponse` türünde bir işaretçi ve bir `error` döndürür.

Fonksiyon ilk olarak `k.GetLoan(ctx, msg.Id)` çağrısını yaparak bir kredi nesnesi alır; burada `ctx` bir context nesnesidir, `k` `msgServer` nesnesidir, `GetLoan` `k` üzerinde bir yöntemdir ve `msg.Id` argüman olarak aktarılan msg nesnesinin bir alanıdır. Kredi bulunamazsa, `nil` ve `sdkerrors.ErrKeyNotFound` ile sarılmış bir hata döndürür.

Ardından, işlev kredinin durumunun "`requested`" olup olmadığını kontrol eder. Değilse, `nil` ve `types.ErrWrongLoanState` ile sarılmış bir hata döndürür.

Kredinin durumu "`requested`" ise, fonksiyon borç veren ve borç alanın adreslerini bech32 dizelerinden ayrıştırır ve ardından kredi `amount`'unu bir dizeden ayrıştırır. Kredi tutarındaki coin'lerin ayrıştırılmasında bir hata varsa, `nil` ve `types.ErrWrongLoanState` ile sarılmış bir hata döndürür.

Aksi takdirde, fonksiyon `k.bankKeeper` nesnesi üzerinde `SendCoins` yöntemini çağırır ve ona context'i, borç veren ve borç alan adreslerini ve kredi miktarını iletir. Ardından kredi nesnesinin borç veren alanını günceller ve durumunu "`approved`" olarak ayarlar. Son olarak, `k.SetLoan(ctx, loan)` çağrısını yaparak güncellenmiş kredi nesnesini saklar.

Sonunda, fonksiyon bir `types.MsgApproveLoanResponse` nesnesi ve hata için `nil` döndürür.

### Özel bir hata kaydetme

&#x20;`ApproveLoan` fonksiyonunda kullanılan `ErrWrongLoanState` özel hatasını kaydetmek için "`errors.go`" dosyasını değiştirin:

x/loan/types/errors.go

```
package types

import (
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
    ErrWrongLoanState = sdkerrors.Register(ModuleName, 2, "wrong loan state")
)
```
