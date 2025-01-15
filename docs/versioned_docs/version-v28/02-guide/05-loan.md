# DeFi Loan

## Introduction

Decentralized finance (DeFi) is a rapidly growing sector that is transforming the way we think about financial instruments and provides an array of inventive financial products and services. These include lending, borrowing, spot trading, margin trading, and flash loans, all of which are available to anyone possessing an internet connection.

A DeFi loan represents a financial contract where the borrower is granted a certain asset, like currency or digital tokens.
In return, the borrower agrees to pay an additional fee and repay the loan within a set period of time.
To secure a loan, the borrower provides collateral that the lender can claim in the event of default.

## You Will Learn

In this tutorial you will learn how to:

- **Scaffold a DeFi Module:**  Learn how to use Ignite CLI to scaffold the basic structure of a DeFi module tailored for loan services.
- **Implement Loan Transactions:** Walk through coding the logic for initiating, managing, and closing loans.
- **Create Custom Tokens:** Understand how to create and manage custom tokens within your DeFi ecosystem, vital for lending and borrowing mechanisms.
- **Integrate Interest Rate Models:** Dive into implementing interest rate models to calculate loan interests dynamically.
- **Ensure Security and Compliance:** Focus on security, ensure your DeFi module is resistant to common vulnerabilities by validating inputs.
- **Test and Debug:** Learn effective strategies for testing your DeFi module and debugging issues that arise during development.

## Setup and Scaffold

1. **Create a New Blockchain:**

```bash
ignite scaffold chain loan --no-module && cd loan
```

Notice the `--no-module` flag, in the next step we make sure the `bank` dependency is included with scaffolding the module.

2. **Create a Module:**
   
Create a new "loan" module that is based on the standard Cosmos SDK `bank` module.

```bash
ignite scaffold module loan --dep bank
```

3. **Define the loan Module:**

The "list" scaffolding command is used to generate files that implement the logic for storing and interacting with data stored as a list in the blockchain state.

```bash
ignite scaffold list loan amount fee collateral deadline state borrower lender --no-message
```

4. **Scaffold the Messages:**

Scaffold the code for handling the messages for requesting, approving, repaying, liquidating, and cancelling loans.

- Handling Loan Requests

```bash
ignite scaffold message request-loan amount fee collateral deadline
```

- Approving and Canceling Loans 

```bash
ignite scaffold message approve-loan id:uint
```

```bash
ignite scaffold message cancel-loan id:uint
```

- Repaying and Liquidating Loans

```bash
ignite scaffold message repay-loan id:uint
```

```bash
ignite scaffold message liquidate-loan id:uint
```

## Additional Features

- **Extend the BankKeeper:**

Ignite takes care of adding the `bank` keeper, but you still need to tell the loan module which bank methods you will be using. You will be using three methods: `SendCoins`, `SendCoinsFromAccountToModule`, and `SendCoinsFromModuleToAccount`. 
Remove the `SpendableCoins` function from the `BankKeeper`.

Add these to the `Bankkeeper` interface.

```go title="x/loan/types/expected_keepers.go"
package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	// SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
```

- **Implement basic Validation to ensure proper loan requests:**

When a loan is created, a certain message input validation is required. You want to throw error messages in case the end user tries impossible inputs.

The `ValidateBasic` function plays a crucial role in maintaining the security and compliance of loan input parameters. By implementing comprehensive input validations, you enhance the security of your application. It's important to rigorously verify all user inputs to ensure they align with the established standards and rules of your platform.

```go title="x/loan/types/message_request_loan.go"
import (
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (msg *MsgRequestLoan) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	amount, _ := sdk.ParseCoinsNormalized(msg.Amount)
	if !amount.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount is not a valid Coins object")
	}
	if amount.Empty() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount is empty")
	}
	fee, _ := sdk.ParseCoinsNormalized(msg.Fee)
	if !fee.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "fee is not a valid Coins object")
	}
	deadline, err := strconv.ParseInt(msg.Deadline, 10, 64)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "deadline is not an integer")
	}
	if deadline <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "deadline should be a positive integer")
	}
	collateral, _ := sdk.ParseCoinsNormalized(msg.Collateral)
	if !collateral.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "collateral is not a valid Coins object")
	}
	if collateral.Empty() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "collateral is empty")
	}
	return nil
}
```

## Using the Platform

1. **As a Borrower:**

Implement `RequestLoan` keeper method that will be called whenever a user requests a loan. `RequestLoan` creates a new loan; Set terms like amount, fee, collateral, and repayment deadline. The collateral from the borrower's account is sent to a module account, and adds the loan to the blockchain's store.

Replace your scaffolded templates with the following code.

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

As a borrower, you have the option to cancel a loan you have created if you no longer want to proceed with it. However, this action is only possible if the loan's current status is marked as "requested".

```go title="x/loan/keeper/msg_server_cancel_loan.go"
package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"loan/x/loan/types"
)

func (k msgServer) CancelLoan(goCtx context.Context, msg *types.MsgCancelLoan) (*types.MsgCancelLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	loan, found := k.GetLoan(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", msg.Id)
	}
	if loan.Borrower != msg.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Cannot cancel: not the borrower")
	}
	if loan.State != "requested" {
		return nil, errorsmod.Wrapf(types.ErrWrongLoanState, "%v", loan.State)
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

2. **As a Lender:**

Approve loan requests and liquidate loans if borrowers fail to repay.

```go title="x/loan/keeper/msg_server_approve_loan.go"
package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"loan/x/loan/types"
)

func (k msgServer) ApproveLoan(goCtx context.Context, msg *types.MsgApproveLoan) (*types.MsgApproveLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	loan, found := k.GetLoan(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", msg.Id)
	}
	if loan.State != "requested" {
		return nil, errorsmod.Wrapf(types.ErrWrongLoanState, "%v", loan.State)
	}
	lender, _ := sdk.AccAddressFromBech32(msg.Creator)
	borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
	amount, err := sdk.ParseCoinsNormalized(loan.Amount)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrWrongLoanState, "Cannot parse coins in loan amount")
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

```go title="x/loan/keeper/msg_server_repay_loan.go"
package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"loan/x/loan/types"
)

func (k msgServer) RepayLoan(goCtx context.Context, msg *types.MsgRepayLoan) (*types.MsgRepayLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	loan, found := k.GetLoan(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", msg.Id)
	}
	if loan.State != "approved" {
		return nil, errorsmod.Wrapf(types.ErrWrongLoanState, "%v", loan.State)
	}
	lender, _ := sdk.AccAddressFromBech32(loan.Lender)
	borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
	if msg.Creator != loan.Borrower {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Cannot repay: not the borrower")
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

```go title="x/loan/keeper/msg_server_liquidate_loan.go"
package keeper

import (
	"context"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"loan/x/loan/types"
)

func (k msgServer) LiquidateLoan(goCtx context.Context, msg *types.MsgLiquidateLoan) (*types.MsgLiquidateLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	loan, found := k.GetLoan(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "key %d doesn't exist", msg.Id)
	}
	if loan.Lender != msg.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Cannot liquidate: not the lender")
	}
	if loan.State != "approved" {
		return nil, errorsmod.Wrapf(types.ErrWrongLoanState, "%v", loan.State)
	}
	lender, _ := sdk.AccAddressFromBech32(loan.Lender)
	collateral, _ := sdk.ParseCoinsNormalized(loan.Collateral)
	deadline, err := strconv.ParseInt(loan.Deadline, 10, 64)
	if err != nil {
		panic(err)
	}
	if ctx.BlockHeight() < deadline {
		return nil, errorsmod.Wrap(types.ErrDeadline, "Cannot liquidate before deadline")
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

Add the custom errors `ErrWrongLoanState` and `ErrDeadline`:

```go title="x/loan/types/errors.go"
package types

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrWrongLoanState = sdkerrors.Register(ModuleName, 2, "wrong loan state")
	ErrDeadline       = sdkerrors.Register(ModuleName, 3, "deadline")
)
```

## Testing the Application

- **Add Test Tokens:**

Configure config.yml to add tokens (e.g., 10000foocoin) to test accounts.

```bash title="config.yml"
version: 1
accounts:
  - name: alice
    coins:
      - 20000token
      - 10000foocoin
      - 200000000stake
  - name: bob
    coins:
      - 10000token
      - 100000000stake
client:
  openapi:
    path: docs/static/openapi.yml
faucet:
  name: bob
  coins:
    - 5token
    - 100000stake
validators:
  - name: alice
    bonded: 100000000stake
```

- **Start the Blockchain:**

```bash
ignite chain serve
```

If everything works successful, you should see the `Blockchain is running` message in the Terminal.

- **Request a loan:**

In a new terminal window, request a loan of `1000token` with `100token` as a fee and `1000foocoin` as a collateral from Alice's account. The deadline is set to `500` blocks:

```bash
loand tx loan request-loan 1000token 100token 1000foocoin 500 --from alice --chain-id loan
```

- **Approve the loan:**

```bash
loand tx loan approve-loan 0 --from bob --chain-id loan
```

- **Repay a loan:**

```bash
loand tx loan repay-loan 0 --from alice --chain-id loan
```

- **Liquidate a loan:**

```bash
loand tx loan request-loan 1000token 100token 1000foocoin 20 --from alice --chain-id loan -y
loand tx loan approve-loan 1 --from bob --chain-id loan -y
loand tx loan liquidate-loan 1 --from bob --chain-id loan -y
```

At any state in the process, use `q list loan` to see the active state of all loans.

```bash
loand q loan list-loan
```

Query the blockchain for balances to confirm they have changed according to your transactions.

```bash
loand q bank balances $(loand keys show alice -a)
```

## Conclusion

This tutorial outlines the process of setting up a decentralized platform for digital asset loans using blockchain technology. By following these steps, you can create a DeFi platform that allows users to engage in secure and transparent lending and borrowing activities.
