---
sidebar_position: 6
description: Loan blockchain using Ignite CLI
title: "Advanced Module: DeFi Loan"
---

# DeFi loan module

As a rapidly growing industry in the blockchain ecosystem, (decentralized finance) DeFi is spurring innovation and revolution in spending, sending, locking, and loaning cryptocurrency tokens.

One of the many goals of blockchain is to make financial instruments available to everyone. A loan in blockchain DeFi can be used in combination with lending, borrowing, spot trading, margin trading, and flash loans.

With DeFi, end users can quickly and easily access loans without having to submit their passports or background checks like in the traditional banking system.

In this tutorial, you learn about a basic loan system as you use Ignite CLI to build a loan module.

**You will learn how to**

* Scaffold a blockchain
* Scaffold a Cosmos SDK loan module
* Scaffold a list for loan objects
* Create messages in the loan module to interact with the loan object
* Interact with other Cosmos SDK modules
* Use an escrow module account
* Add application messages for a loan system
  * Request loan
  * Approve loan
  * Repay loan
  * Liquidate loan
  * Cancel loan

**Note:** The code in this tutorial is written specifically for this learning experience and is intended only for educational purposes. This tutorial code is not intended to be used in production.

## Module design

A loan consists of:

* An `id`
* The `amount` that is being lent
* A `fee` as cost for the loan
* The borrowing party provides a `collateral` to request a loan
* A loan has a `deadline` for repayment, after which the loan can be liquidated
* A loan has a `state` that describes the status as:

	* requested
	* approved
	* paid
	* cancelled
	* liquidated

The two accounts involved in the loan are:

* `borrower`
* `lender`

### The borrower

A borrower posts a loan request with loan information such as:

* `amount`
* `fee`
* `collateral`
* `deadline`

The borrower must repay the loan amount and the loan fee to the lender by the deadline risk losing the collateral.

### The lender

A lender can approve a loan request from a borrower.

- After the lender approves the loan, the loan amount is transferred to the borrower.
- If the borrower is unable to pay the loan, the lender can liquidate the loan.
- Loan liquidation transfers the collateral and the fees to the lender.

## Scaffold the blockchain

Use Ignite CLI to scaffold a fully functional Cosmos SDK blockchain app named `loan`:

```bash
ignite scaffold chain loan --no-module
```

The `--no-module` flag prevents scaffolding a default module. Don't worry, you will add the loan module later.

Change into the newly created `loan` directory:

```bash
cd loan
```

## Scaffold the module

Scaffold the module to create a new `loan` module. Following the Cosmos SDK convention, all modules are scaffolded inside the `x` directory:

```bash
ignite scaffold module loan --dep bank
```

Use the `--dep` flag to specify that this module depends on and is going to interact with the Cosmos SDK `bank` module.

## Scaffold a list

Use the `ignite scaffold list` command to scaffold code necessary to store loans in an array-like data structure:

```bash
ignite scaffold list loan amount fee collateral deadline state borrower lender --no-message
```

Use the `--no-message` flag to disable CRUD messages in the scaffold.

The data you store in an array-like data structure are the loans, with these parameters that are defined in the `Loan` message in `proto/loan/loan.proto`:

```protobuf
message Loan {
  uint64 id = 1;
  string amount = 2;
  string fee = 3;
  string collateral = 4;
  string deadline = 5;
  string state = 6;
  string borrower = 7;
  string lender = 8;  
}
```

Later, you define the messages to interact with the loan list.

Now it is time to use messages to interact with the loan module. But first, make sure to store your current state in a git commit:

```bash
git add .
git commit -m "Scaffold loan module and loan list"
```

## Scaffold the messages

In order to create a loan app, you need the following messages:

* Request loan
* Approve loan
* Repay loan
* Liquidate loan
* Cancel loan

You can use the `ignite scaffold message` command to create each of the messages.

You define the details of each message when you scaffold them.

Create the messages one at a time with the according application logic.

### Request loan message

For a loan, the initial message handles the transaction when a username requests a loan.

The username wants a certain `amount` and is willing to pay `fees` as well as give `collateral`. The `deadline` marks the time when the loan has to be repaid.

The first message is the `request-loan` message that  requires these input parameters:

* `amount`
* `fee`
* `collateral`
* `deadline`

```bash
ignite scaffold message request-loan amount fee collateral deadline
```

For the sake of simplicity, define every parameter as a string.

The `request-loan` message creates a new loan object and locks the tokens to be spent as fee and collateral into an escrow account. Describe these conditions in the module keeper `x/loan/keeper/msg_server_request_loan.go`:

```go
package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"loan/x/loan/types"
)

func (k msgServer) RequestLoan(goCtx context.Context, msg *types.MsgRequestLoan) (*types.MsgRequestLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Create a new Loan with the following user input
	var loan = types.Loan{
		Amount:     msg.Amount,
		Fee:        msg.Fee,
		Collateral: msg.Collateral,
		Deadline:   msg.Deadline,
		State:      "requested",
		Borrower:   msg.Creator,
	}

	// TODO: collateral has to be more than the amount (+fee?)

	// moduleAcc := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
	// Get the borrower address
	borrower, _ := sdk.AccAddressFromBech32(msg.Creator)

	// Get the collateral as sdk.Coins
	collateral, err := sdk.ParseCoinsNormalized(loan.Collateral)
	if err != nil {
		panic(err)
	}

	// Use the module account as escrow account
	sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, borrower, types.ModuleName, collateral)
	if sdkError != nil {
		return nil, sdkError
	}

	// Add the loan to the keeper
	k.AppendLoan(
		ctx,
		loan,
	)

	return &types.MsgRequestLoanResponse{}, nil
}
```

Since this function is using the `bankKeeper` with the function `SendCoinsFromAccountToModule`, you must add the `SendCoinsFromAccountToModule` function to `x/loan/types/expected_keepers.go` like this:

```go
package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}
```

### Validate the input

When a loan is created, a certain message input validation is required. You want to throw error messages in case the end user tries impossible inputs.

You can describe message validation errors in the modules `types` directory.

Add the following code to the `ValidateBasic()` function in the `x/loan/types/message_request_loan.go` file:

```go
func (msg *MsgRequestLoan) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
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

Congratulations, you have created the `request-loan` message.

## Run and test your first message

You can run the chain and test your first message.

Start the blockchain:

```bash
ignite chain serve
```

Add your first loan:

```bash
loand tx loan request-loan 100token 2token 200token 500 --from alice
```

Query your loan:

```bash
loand query loan list-loan
```

The loan shows in the list:

```bash
Loan:
- amount: 100token
  borrower: cosmos17mnrhwchwc8trg4w09s0gvvfsvt58ejtsykkm6
  collateral: 200token
  deadline: "500"
  fee: 2token
  id: "0"
  lender: ""
  state: requested
```

You can stop the blockchain again with CTRL+C.

### Save iterative changes

This is a good time to add your advancements to git:

```bash
git add .
git commit -m "Add request-loan message"
```

### Approve loan message

After a loan request has been published, another account can approve the loan and agree to the terms of the borrower.

The message `approve-loan` has one parameter, the `id`.
Specify the type of `id` as `uint`. By default, ids are stored as `uint`.

```bash
ignite scaffold message approve-loan id:uint
```

This message must be available for all loan types that are in `"requested"` status.

The loan approval sends the requested coins for the loan to the borrower and sets the loan state to `"approved"`.

Modify the `x/loan/keeper/msg_server_approve_loan.go` to implement this logic:

```go
package keeper

import (
	"context"
	"fmt"

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

	// TODO: for some reason the error doesn't get printed to the terminal
	if loan.State != "requested" {
		return nil, sdkerrors.Wrapf(types.ErrWrongLoanState, "%v", loan.State)
	}

	lender, _ := sdk.AccAddressFromBech32(msg.Creator)
	borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
	amount, err := sdk.ParseCoinsNormalized(loan.Amount)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrWrongLoanState, "Cannot parse coins in loan amount")
	}

	k.bankKeeper.SendCoins(ctx, lender, borrower, amount)

	loan.Lender = msg.Creator
	loan.State = "approved"

	k.SetLoan(ctx, loan)

	return &types.MsgApproveLoanResponse{}, nil
}
```

This module uses the `SendCoins` function of `bankKeeper`. Add this `SendCoins` function to the `x/loan/types/expected_keepers.go` file:

```go
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface {
	// Methods imported from bank should be defined here
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}
```

Now, define the `ErrWrongLoanState` new error type by adding it to the errors definitions in `x/loan/types/errors.go`:

```go
package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/loan module sentinel errors
var (
	ErrWrongLoanState = sdkerrors.Register(ModuleName, 2, "wrong loan state")
)
```

Start the blockchain and use the two commands you already have available:

```bash
ignite chain serve -r
```

Use the `-r` flag to reset the blockchain state and start with a new database.

Now, request a loan from `bob`:

```bash
loand tx loan request-loan 100token 2token 200token 500 --from bob -y
```

Query your loan request:

```bash
loand query loan list-loan
```

Approve the loan:

```bash
loand tx loan approve-loan 0 --from alice -y
```

This approve loan transaction sends the balances according to the loan request.

Check for the loan list again to verify that the loan state is now `approved`.

```bash
Loan:
- amount: 100token
  borrower: cosmos1sx8k358xw5pulv7acjhm6klvn3tukk2r2a74gg
  collateral: 200token
  deadline: "500"
  fee: 2token
  id: "0"
  lender: cosmos1qxm2dtupmr8pp20m0t7tmjq6gq2z8j3d6ltr9d
  state: approved
pagination:
  next_key: null
  total: "0"
```

You can query for alice's balance to see the loan in effect. Take the lender address from above, this is alice address:

```bash
loand query bank balances <alice_address>
```

In case everything works as expected, this is a good time to save the state with a git commit:

```bash
git add .
git commit -m "Add approve loan message"
```

### Repay Loan Message

After the loan has been approved, the username must be able to repay an approved loan.

Scaffold the message `repay-loan` that a borrower uses to return tokens that were borrowed from the lender:

```bash
ignite scaffold message repay-loan id:uint
```

Repaying a loan requires that the loan is in `"approved"` status.

The coins as described in the loan are collected and sent from the borrower to the lender, along with the agreed fees.

The `collateral` is released from the escrow module account.

Only the `borrower` can repay the loan.

This loan repayment logic is defined in `x/loan/keeper/msg_server_repay_loan.go`:

```go
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
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
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
		return nil, sdkerrors.Wrap(types.ErrWrongLoanState, "Cannot send coins")
	}
	err = k.bankKeeper.SendCoins(ctx, borrower, lender, fee)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrWrongLoanState, "Cannot send coins")
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, borrower, collateral)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrWrongLoanState, "Cannot send coins")
	}

	loan.State = "repayed"

	k.SetLoan(ctx, loan)

	return &types.MsgRepayLoanResponse{}, nil
}
```

After the coins have been successfully exchanged, the state of the loan is set to `repayed`.

To release tokens with the `SendCoinsFromModuleToAccount` function of `bankKeepers`, you need to add the `SendCoinsFromModuleToAccount` function to the `x/loan/types/expected_keepers.go`:

```go
package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface {
	// Methods imported from bank should be defined here
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}
```

Start the blockchain and use the two commands you already have available:

```bash
ignite chain serve -r
```

Use the `-r` flag to reset the blockchain state and start with a new database:

```bash
loand tx loan request-loan 100token 2token 200token 500 --from bob -y
```

Query your loan request:

```bash
loand query loan list-loan
```

Approve the loan:

```bash
loand tx loan approve-loan 0 --from alice -y
```

You can query for alice's balance to see the loan in effect.

Take the lender address from above, this is alice address:

```bash
loand query bank balances <alice_address>
```

Now repay the loan:

```bash
loand tx loan repay-loan 0 --from bob -y
```

The loan status is now `repayed`:

```bash
Loan:
- amount: 100token
  borrower: cosmos1200nsqsxcyxtllfgal5x8qhqwj8km64ft0eu2d
  collateral: 200token
  deadline: "500"
  fee: 2token
  id: "0"
  lender: cosmos194pn6vly2nlald3zjqcxfnvasa0xt7ect6h6qk
  state: repayed
```

The alice balance reflects the repayed amount plus fees:

```bash
loand query bank balances <alice_address>
```

Good job!

Update your git with the changes you made:

```bash
git add .
git commit -m "Add repay-loan message"
```

### Liquidate Loan Message

After the deadline is passed, a lender can liquidate a loan when the borrower does not repay the tokens. The message to `liquidate-loan` refers to the loan `id`:

```bash
ignite scaffold message liquidate-loan id:uint
```

* The `liquidate-loan` message must be able to be executed by the `lender`.
* The status of the loan must be `approved`.
* The `deadline` block height must have passed.

When these properties are valid, the collateral shall be liquidated from the `borrower`.

Add this liquidate loan logic to the `keeper` in `x/loan/keeper/msg_server_liquidate_loan.go`:

```go
package keeper

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"loan/x/loan/types"
)

func (k msgServer) LiquidateLoan(goCtx context.Context, msg *types.MsgLiquidateLoan) (*types.MsgLiquidateLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	loan, found := k.GetLoan(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
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

	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, lender, collateral)

	loan.State = "liquidated"

	k.SetLoan(ctx, loan)

	return &types.MsgLiquidateLoanResponse{}, nil
}
```

Add the new error `ErrDeadline` to the error messages in `x/loan/types/errors.go`:

```go
package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/loan module sentinel errors
var (
	ErrWrongLoanState = sdkerrors.Register(ModuleName, 2, "wrong loan state")
	ErrDeadline       = sdkerrors.Register(ModuleName, 3, "deadline")
)
```

These changes are required for the `liquidate-loan` message.

### Test liquidation message

You can test the liquidation message now. Start your chain and reset the state of the app:

```bash
ignite chain serve -r
```

Set the deadline for the loan request to 1 block:

```bash
loand tx loan request-loan 100token 2token 200token 1 --from bob -y
```

Query your loan request:

```bash
loand query loan list-loan
```

Approve the loan:

```bash
loand tx loan approve-loan 0 --from alice -y
```

You can query for alice's balances to see the loan in effect.

Take the lender address from above, this is alice address.

```bash
loand query bank balances <alice_address>
```

Now, liquidate the loan:

```bash
loand tx loan liquidate-loan 0 --from alice -y
```

Query the loan:

```bash
loand query loan list-loan
```

The loan status is now `liquidated`:

```bash
Loan:
- amount: 100token
  borrower: cosmos1lp4ghp4mmsdgpf2fm22f0qtqmnjeh3gr9h3cau
  collateral: 200token
  deadline: "1"
  fee: 2token
  id: "0"
  lender: cosmos1w6pfj52jp809pyp2a2h573cta23rc0zsulpafm
  state: liquidated
```

And alice balance reflects the repayed amount plus fees:

```bash
loand query bank balances <alice_address>
```

Add the changes to your local repository:

```bash
git add .
git commit -m "Add liquidate-loan message"
```

### Cancel loan message

After a loan request has been made and not been approved, the `borrower` must be able to cancel a loan request.

Scaffold the message for `cancel-loan`:

```bash
ignite s message cancel-loan id:uint
```

* Only the `borrower` can cancel a loan request.
* The state of the request must be `requested`.
* Then the collateral coins can be released from escrow and the status set to `cancelled`.

Add this functionality to the `keeper` in `x/loan/keeper/msg_server_cancel_loan.go`:

```go
package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"loan/x/loan/types"
)

func (k msgServer) CancelLoan(goCtx context.Context, msg *types.MsgCancelLoan) (*types.MsgCancelLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	loan, found := k.GetLoan(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	if loan.Borrower != msg.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Cannot cancel: not the borrower")
	}

	if loan.State != "requested" {
		return nil, sdkerrors.Wrapf(types.ErrWrongLoanState, "%v", loan.State)
	}

	borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
	collateral, _ := sdk.ParseCoinsNormalized(loan.Collateral)
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, borrower, collateral)

	loan.State = "cancelled"

	k.SetLoan(ctx, loan)

	return &types.MsgCancelLoanResponse{}, nil
}
```

### Test cancelling a loan

Test the changes for cancelling a loan request:

```bash
ignite chain serve -r
```

```bash
loand tx loan request-loan 100token 2token 200token 100 --from bob -y
```

Query your loan request:

```bash
loand query loan list-loan
```

```bash
loand tx loan cancel-loan 0 --from bob -y
```

Query your loan request:

```bash
loand query loan list-loan
```

Now the collateral coins can be released from escrow and the status set to `cancelled`.

```bash
- amount: 100token
  borrower: cosmos1lp4ghp4mmsdgpf2fm22f0qtqmnjeh3gr9h3cau
  collateral: 200token
  deadline: "100"
  fee: 2token
  id: "2"
  lender: ""
  state: cancelled
```

Consider again updating your local repository with a git commit. After you test and use your loan module, consider publishing your code to a public repository for others to see your accomplishments.

```bash
git add .
git commit -m "Add cancel-loan message"
```

## Complete

Congratulations. You have completed the loan module tutorial.

You executed commands and updated files to:

* Scaffold a blockchain
* Scaffold a module
* Scaffold a list for loan objects
* Create messages in your module to interact with the loan object
* Interact with other modules in your module
* Use an escrow module account
* Add application messages for a loan system
  * Request Loan
  * Approve Loan
  * Repay Loan
  * Liquidate Loan
  * Cancel Loan
