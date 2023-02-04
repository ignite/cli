# Request a loan

Bir kullanıcı kredi talep ettiğinde çağrılacak olan `RequestLoan` keeper yöntemini uygulayın. `RequestLoan`, sağlanan verilerle yeni bir kredi oluşturur, teminatı borçlunun hesabından bir modül hesabına gönderir ve krediyi blok zincirinin deposuna ekler.

x/loan/keeper/msg\_server\_request\_loan.go

```
package keeper

import (
    "context"

    sdk "github.com/cosmos/cosmos-sdk/types"

    "loan/x/loan/types"
)

func (k msgServer) RequestLoan(goCtx context.Context, msg *types.MsgRequestLoan) (*types.MsgRequestLoanResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    var loan = types.Loan{
        Amount:     msg.Amount,
        Fee:        msg.Fee,
        Collateral: msg.Collateral,
        Deadline:   msg.Deadline,
        State:      "requested",
        Borrower:   msg.Creator,
    }
    borrower, err := sdk.AccAddressFromBech32(msg.Creator)
    if err != nil {
        panic(err)
    }
    collateral, err := sdk.ParseCoinsNormalized(loan.Collateral)
    if err != nil {
        panic(err)
    }
    sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, borrower, types.ModuleName, collateral)
    if sdkError != nil {
        return nil, sdkError
    }
    k.AppendLoan(ctx, loan)
    return &types.MsgRequestLoanResponse{}, nil
}
```

Fonksiyon iki bağımsız değişken alır: bir `context.Context` nesnesi ve bir `types.MsgRequestLoan` yapısına bir işaretçi. Bir `types.MsgRequestLoanResponse` yapısına bir işaretçi ve bir `error` nesnesi döndürür.

Fonksiyonun yaptığı ilk şey, `types.MsgRequestLoan` struct girişindeki verilerle yeni bir `types.Loan` struct oluşturmaktır. `types.Loan` yapısının `State` alanını "requested" olarak ayarlar.

Ardından, fonksiyon `types.MsgRequestLoan` yapısının `msg.Creator` alanından borçlunun adresini alır. Daha sonra `sdk.ParseCoinsNormalized` fonksiyonunu kullanarak `loan.Collateral` alanını (bir dize olan) `sdk.Coins`'e ayrıştırır.

Fonksiyon daha sonra `k.bankKeeper.SendCoinsFromAccountToModule` fonksiyonunu kullanarak teminatı borçlunun hesabından bir modül hesabına gönderir. Son olarak, `k.AppendLoan` fonksiyonunu kullanarak yeni krediyi bir tutucuya ekler. Fonksiyon bir `types.MsgRequestLoanResponse` yapısı ve her şey yolunda giderse `nil` hatası döndürür.

### Temel mesaj doğrulama

Bir kredi oluşturulduğunda, belirli bir mesaj girişi doğrulaması gereklidir. Son kullanıcının imkansız girdileri denemesi durumunda hata mesajları atmak istersiniz.

x/loan/types/message\_request\_loan.go

```
package types

import (
    "strconv"

    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (msg *MsgRequestLoan) ValidateBasic() error {
    _, err := sdk.AccAddressFromBech32(msg.Creator)
    if err != nil {
        return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
    }
    amount, _ := sdk.ParseCoinsNormalized(msg.Amount)
    if !amount.IsValid() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "amount is not a valid Coins object")
    }
    if amount.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "amount is empty")
    }
    fee, _ := sdk.ParseCoinsNormalized(msg.Fee)
    if !fee.IsValid() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "fee is not a valid Coins object")
    }
    deadline, err := strconv.ParseInt(msg.Deadline, 10, 64)
    if err != nil {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "deadline is not an integer")
    }
    if deadline <= 0 {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "deadline should be a positive integer")
    }
    collateral, _ := sdk.ParseCoinsNormalized(msg.Collateral)
    if !collateral.IsValid() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collateral is not a valid Coins object")
    }
    if collateral.Empty() {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collateral is empty")
    }
    return nil
}
```
