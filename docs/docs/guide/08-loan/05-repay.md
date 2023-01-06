# Repay a loan

The `RepayLoan` method is responsible for handling the repayment of a loan. This
involves transferring the borrowed funds, along with any agreed upon fees, from
the borrower to the lender. In addition, the collateral that was provided as
part of the loan agreement will be released from the escrow account and returned
to the borrower.

It is important to note that the `RepayLoan` method can only be called under
certain conditions. Firstly, the transaction must be signed by the borrower of
the loan. This ensures that only the borrower has the ability to initiate the
repayment process. Secondly, the loan must be in an approved status. This means
that the loan has received approval and is ready to be repaid.

To implement the `RepayLoan` method, we must ensure that these conditions are
met before proceeding with the repayment process. Once the necessary checks have
been performed, the method can then handle the transfer of funds and the release
of the collateral from the escrow account.

## Keeper method

```go title="x/loan/keeper/msg_server_repay_loan.go"
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

`RepayLoan` takes in two arguments: a context and a pointer to a
`types.MsgRepayLoan` type. It returns a pointer to a
`types.MsgRepayLoanResponse` type and an `error`.

The method first retrieves a loan from storage by passing the provided loan ID
to the `k.GetLoan` method. If the loan cannot be found, the method returns an
error wrapped in a `sdkerrors.ErrKeyNotFound` error.

The method then checks that the state of the loan is "approved". If it is not,
the method returns an error wrapped in a `types.ErrWrongLoanState` error.

Next, the method converts the lender and borrower addresses stored in the loan
struct to `sdk.AccAddress` types using the `sdk.AccAddressFromBech32` function.
It then checks that the transaction is signed by the borrower of the loan by
comparing the `msg.Creator` field to the borrower address stored in the loan
struct. If these do not match, the method returns an error wrapped in a
`sdkerrors.ErrUnauthorized` error.

The method then parses the loan amount, fee, and collateral stored in the loan
struct as `sdk.Coins` using the `sdk.ParseCoinsNormalized` function. It then
uses the `k.bankKeeper.SendCoins` function to transfer the loan amount and fee
from the borrower to the lender. It then uses the
`k.bankKeeper.SendCoinsFromModuleToAccount` function to transfer the collateral
from the escrow account to the borrower.

Finally, the method updates the state of the loan to "repayed" and stores the
updated loan in storage using the `k.SetLoan` method. The method returns a
`types.MsgRepayLoanResponse` and a `nil` error to indicate that the repayment
process was successful.