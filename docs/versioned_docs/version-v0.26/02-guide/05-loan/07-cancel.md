# Cancel a loan

As a borrower, you have the option to cancel a loan you have created if you no
longer want to proceed with it. However, this action is only possible if the
loan's current status is marked as "requested".

If you decide to cancel the loan, the collateral tokens that were being held as
security for the loan will be transferred back to your account from the module
account. This means that you will regain possession of the collateral tokens you
had originally put up for the loan.

```go title="x/loan/keeper/msg_server_cancel_loan.go"
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

`CancelLoan` takes in two arguments: a `context.Context` named `goCtx` and a
pointer to a `types.MsgCancelLoan` named `msg`. It returns a pointer to a
`types.MsgCancelLoanResponse` and an error.

The function begins by using the `sdk.UnwrapSDKContext` method to get the
`sdk.Context` from the `context.Context` object. It then uses the `GetLoan`
method of the `msgServer` type to retrieve a loan identified by the `Id` field
of the `msg` argument. If the loan is not found, the function returns an error
using the `sdk.ErrKeyNotFound` error wrapped with the `sdk.Wrap` method.

Next, the function checks if the `Creator` field of the msg argument is the same
as the `Borrower` field of the loan. If they are not the same, the function
returns an error using the `sdk.ErrUnauthorized` error wrapped with the
`sdk.Wrap` method.

The function then checks if the `State` field of the loan is equal to the string
`"requested"`. If it is not, the function returns an error using the
types.`ErrWrongLoanState` error wrapped with the `sdk.Wrapf` method.

If the loan has the correct state and the creator of the message is the borrower
of the loan, the function proceeds to send the collateral coins held in the
`Collateral` field of the loan back to the borrower's account using the
`SendCoinsFromModuleToAccount` method of the `bankKeeper`. The function then
updates the State field of the loan to the string "cancelled" and sets the
updated loan using the `SetLoan` method. Finally, the function returns a
`types.MsgCancelLoanResponse` object and a nil error.