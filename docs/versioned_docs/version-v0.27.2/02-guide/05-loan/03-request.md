# Request a loan

Implement `RequestLoan` keeper method that will be called whenever a user
requests a loan. `RequestLoan` creates a new loan with the provided data, sends
the collateral from the borrower's account to a module account, and adds the
loan to the blockchain's store.

## Keeper method

```go title="x/loan/keeper/msg_server_request_loan.go"
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

The function takes in two arguments: a `context.Context` object and a pointer to
a `types.MsgRequestLoan` struct. It returns a pointer to a
`types.MsgRequestLoanResponse` struct and an `error` object.

The first thing the function does is create a new `types.Loan` struct with the
data from the input `types.MsgRequestLoan` struct. It sets the `State` field of
`the types.Loan` struct to "requested".

Next, the function gets the borrower's address from the `msg.Creator` field of
the input `types.MsgRequestLoan` struct. It then parses the `loan.Collateral`
field (which is a string) into `sdk.Coins` using the `sdk.ParseCoinsNormalized`
function.

The function then sends the collateral from the borrower's account to a module
account using the `k.bankKeeper.SendCoinsFromAccountToModule` function. Finally,
it adds the new loan to a keeper using the `k.AppendLoan` function. The function
returns a `types.MsgRequestLoanResponse` struct and a `nil` error if all goes
well.

## Basic message validation

When a loan is created, a certain message input validation is required. You want
to throw error messages in case the end user tries impossible inputs.

```go title="x/loan/types/message_request_loan.go"
package types

import (
	// highlight-next-line
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (msg *MsgRequestLoan) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	// highlight-start
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
	// highlight-end
	return nil
}
```
