# Kredi geri ödemesi

`RepayLoan` yöntemi, bir kredinin geri ödemesinin gerçekleştirilmesinden sorumludur. Bu, üzerinde anlaşmaya varılan ücretlerle birlikte ödünç alınan fonların borçludan borç verene aktarılmasını içerir. Buna ek olarak, kredi sözleşmesinin bir parçası olarak sağlanan teminat emanet hesabından serbest bırakılacak ve borçluya iade edilecektir.

`RepayLoan` yönteminin yalnızca belirli koşullar altında çağrılabileceğini unutmamak önemlidir. İlk olarak, işlem kredi borçlusu tarafından imzalanmalıdır. Bu, geri ödeme işlemini yalnızca borçlunun başlatabilmesini sağlar. İkinci olarak, kredinin onaylanmış durumda olması gerekir. Bu, kredinin onay aldığı ve geri ödenmeye hazır olduğu anlamına gelir.

`RepayLoan` yöntemini uygulamak için, geri ödeme işlemine geçmeden önce bu koşulların karşılandığından emin olmalıyız. Gerekli kontroller yapıldıktan sonra, metot fon transferini ve teminatın emanet hesabından serbest bırakılmasını gerçekleştirebilir.

### Keeper yöntemi

x/loan/keeper/msg\_server\_repay\_loan.go

```
package keeper

import (
    "context"
    "fmt"

    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

    "loan/x/loan/types"
)

func (k msgServer) RepayLoan(goCtx context.Context, msg *types.MsgRepayLoan) (*types.MsgRepayLoanResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    loan, found := k.GetLoan(ctx, msg.Id)
    if !found {
        return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", msg.Id)
    }
    if loan.State != "approved" {
        return nil, sdkerrors.Wrapf(types.ErrWrongLoanState, "%v", loan.State)
    }
    lender, _ := sdk.AccAddressFromBech32(loan.Lender)
    borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
    if msg.Creator != loan.Borrower {
        return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Cannot repay: not the borrower")
    }
    amount, _ := sdk.ParseCoinsNormalized(loan.Amount)
    fee, _ := sdk.ParseCoinsNormalized(loan.Fee)
    collateral, _ := sdk.ParseCoinsNormalized(loan.Collateral)
    err := k.bankKeeper.SendCoins(ctx, borrower, lender, amount)
    if err != nil {
        return nil, err
    }
    err = k.bankKeeper.SendCoins(ctx, borrower, lender, fee)
    if err != nil {
        return nil, err
    }
    err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, borrower, collateral)
    if err != nil {
        return nil, err
    }
    loan.State = "repayed"
    k.SetLoan(ctx, loan)
    return &types.MsgRepayLoanResponse{}, nil
}
```

`RepayLoan` iki bağımsız değişken alır: bir context ve bir `types.MsgRepayLoan` türüne bir işaretçi. Bir `types.MsgRepayLoanResponse` türüne bir işaretçi ve bir `error` döndürür.

Yöntem ilk olarak, sağlanan kredi kimliğini `k.GetLoan` yöntemine aktararak depodan bir kredi alır. Kredi bulunamazsa, yöntem `sdkerrors.ErrKeyNotFound` hatasına sarılmış bir hata döndürür.

Yöntem daha sonra kredinin durumunun "approved" olup olmadığını kontrol eder. Değilse, yöntem `types.ErrWrongLoanState` hatasına sarılmış bir hata döndürür.

Ardından, yöntem `sdk.AccAddressFromBech32` işlevini kullanarak kredi yapısında saklanan borç veren ve borç alan adreslerini `sdk.AccAddress` türlerine dönüştürür. Ardından, `msg.Creator` alanını kredi yapısında saklanan borç alan adresiyle karşılaştırarak işlemin krediyi alan tarafından imzalanıp imzalanmadığını kontrol eder. Bunlar eşleşmezse, yöntem `sdkerrors.ErrUnauthorized` hatasına sarılmış bir hata döndürür.

Yöntem daha sonra `sdk.ParseCoinsNormalized` işlevini kullanarak kredi yapısında saklanan kredi tutarını, ücreti ve `teminatı sdk.Coins` olarak ayrıştırır. Ardından, kredi tutarını ve ücreti borçludan borç verene aktarmak için `k.bankKeeper.SendCoins` işlevini kullanır. Daha sonra teminatı emanet hesabından borç alana aktarmak için `k.bankKeeper.SendCoinsFromModuleToAccount` fonksiyonunu kullanır.

Son olarak, yöntem kredinin durumunu "repayed" olarak günceller ve `k.SetLoan` yöntemini kullanarak güncellenmiş krediyi depoda saklar. Yöntem, geri ödeme işleminin başarılı olduğunu belirtmek için bir `types.MsgRepayLoanResponse` ve `nil` hatası döndürür.
