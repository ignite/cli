# Liquidate loan

The `LiquidateLoan` method is a function that allows the lender to sell off the
collateral belonging to the borrower in the event that the borrower has failed
to repay the loan by the specified deadline. This process is known as
"liquidation" and is typically carried out as a way for the lender to recoup
their losses in the event that the borrower is unable to fulfill their repayment
obligations.

During the liquidation process, the collateral tokens that have been pledged by
the borrower as security for the loan are transferred from the borrower's
account to the lender's account. This transfer is initiated by the lender and is
typically triggered when the borrower fails to repay the loan by the agreed upon
deadline. Once the collateral has been transferred, the lender can then sell it
off in order to recoup their losses and compensate for the unpaid loan.

## Keeper method

```go title="x/loan/keeper/msg_server_liquidate_loan.go"
package keeper

import (
	"context"
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

`LiquidateLoan` takes in a context and a `types.MsgLiquidateLoan` message as input and returns a types.MsgLiquidateLoanResponse message and an error as output.

The function first retrieves a loan using the `GetLoan` method and the `Id` field of the input message. If the loan is not found, it returns an error using the `sdkerrors.Wrap` function and the `sdkerrors.ErrKeyNotFound` error code.

Next, the function checks that the `Creator` field of the input message is the same as the `Lender` field of the loan. If they are not the same, it returns an error using the `sdkerrors.Wrap` function and the `sdkerrors.ErrUnauthorized` error code.

The function then checks that the State field of the loan is equal to "approved". If it is not, it returns an error using the `sdkerrors.Wrapf` function and the `types.ErrWrongLoanState` error code.

The function then converts the Lender field of the loan to an address using the `sdk.AccAddressFromBech32` function and the `Collateral` field to coins using the `sdk.ParseCoinsNormalized` function. It also converts the `Deadline` field to an integer using the `strconv.ParseInt` function. If this function returns an error, it panics.

Finally, the function checks that the current block height is greater than or equal to the deadline. If it is not, it returns an error using the `sdkerrors.Wrap` function and the `types.ErrDeadline` error code. If all checks pass, the function uses the `bankKeeper.SendCoinsFromModuleToAccount` method to transfer the collateral from the module account to the lender's account and updates the `State` field of the loan to `"liquidated"`. It then stores the updated loan using the `SetLoan` method and returns a `types.MsgLiquidateLoanResponse` message with no error.

## Register a custom error

```go title="x/loan/types/errors.go"
package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrWrongLoanState = sdkerrors.Register(ModuleName, 2, "wrong loan state")
	// highlight-next-line
	ErrDeadline = sdkerrors.Register(ModuleName, 3, "deadline")
)
```
