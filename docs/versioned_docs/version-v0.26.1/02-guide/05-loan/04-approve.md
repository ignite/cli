# Approve a loan

After a loan request has been made, it is possible for another account to
approve the loan and accept the terms proposed by the borrower. This process
involves the transfer of the requested funds from the lender to the borrower.

To be eligible for approval, a loan must have a status of "requested." This
means that the borrower has made a request for a loan and is waiting for a
lender to agree to the terms and provide the funds. Once a lender has decided to
approve the loan, they can initiate the transfer of the funds to the borrower.

Upon loan approval, the status of the loan is changed to "approved." This
signifies that the funds have been successfully transferred and that the loan
agreement is now in effect.

## Keeper method

```go title="x/loan/keeper/msg_server_approve_loan.go"
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

`ApproveLoan` takes a context and a message of type `*types.MsgApproveLoan` as
input, and returns a pointer to a `types.MsgApproveLoanResponse` and an `error`.

The function first retrieves a loan object by calling `k.GetLoan(ctx, msg.Id)`,
where `ctx` is a context object, `k` is the `msgServer` object, `GetLoan` is a
method on `k`, and `msg.Id` is a field of the msg object passed as an argument.
If the loan is not found, it returns `nil` and an error wrapped with
`sdkerrors.ErrKeyNotFound`.

Next, the function checks if the loan's state is `"requested"`. If it is not, it
returns `nil` and an error wrapped with `types.ErrWrongLoanState`.

If the loan's state is `"requested"`, the function parses the addresses of the
lender and borrower from bech32 strings, and then parses the `amount` of the
loan from a string. If there is an error parsing the coins in the loan amount,
it returns `nil` and an error wrapped with `types.ErrWrongLoanState`.

Otherwise, the function calls the `SendCoins` method on the `k.bankKeeper`
object, passing it the context, the lender and borrower addresses, and the
amount of the loan. It then updates the lender field of the loan object and sets
its state to `"approved"`. Finally, it stores the updated loan object by calling
`k.SetLoan(ctx, loan)`.

At the end, the function returns a `types.MsgApproveLoanResponse` object and
`nil` for the error.

## Register a custom error

To register the custom error `ErrWrongLoanState` that is used in the
`ApproveLoan` function, modify the "errors.go" file:

```go title="x/loan/types/errors.go"
package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrWrongLoanState = sdkerrors.Register(ModuleName, 2, "wrong loan state")
)
```