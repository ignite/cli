---
description: Loan blockchain using Starport
order: 6
title: "Advanced module: DeFi Loan"
---

# Loan Module

As a repid growing industry in the blockchain ecosystem, DeFi (Decentralized Finance) is the term for innovation and revolution in spending, sending, locking and loaning cryptocurrency tokens.

One of the many traits of blockchain is to make financial instruments avilable to anyone, via an open protocol. A loan is used in combination with lending, borrowing, spot trading, margin trading and flash loans.

With DeFi, anyone can quick and easy take loans without having to submit their passports or background checks like in the traditional banking system.

In this tutorial you will learn about a basic loan system, built as a module on Starport.

You will learn how to:

* Scaffold a blockchain
* Scaffold a module
* Scaffold a list for loan objects
* Create messages in your module to interact with the loan object
* Interact with other modules in your module
* Use an escrow module account
* Add application messages for a Loan system
  * Request Loan
  * Approve Loan
  * Repay Loan
  * Liquidate Loan
  * Cancel Loan

Warning: This module is purely for learning purposes. It is not tested in production.

## Module Design

A loan consists of an `id`, the `amount` that is being lend out, a `fee` as cost for the loan.
The borrowing party will provide a `collateral` to request a loan. A loan has a `deadline` for when it is supposed to be due and can be liquidated.
Furthermore the `state` of the loan describes if it is in suggestion, approved, payed back, or liquidated.
This is the resulting data schema for the Loan module.

```proto
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

The two accounts that get involved in the loan will be the `borrower` and the `lender`.

### The Borrower

A borrower will post a loan request with information such as - loan `amount`, `fee`, `collateral` and `deadline`.
The borrower will have to repay the loan transfer amount and fee to the lender or the account risks loosing the collateral.

### The Lender

A lender can approve a loan request from a borrower. Approving the loan transfers the loan amount to the the borrower. If the borrower is unable to pay the loan, the lender can liquidate the loan which transfers the collateral and the fees to the lender.

## Scaffold the Blockchain

```bash
starport scaffold chain github.com/cosmonaut/loan --no-module
```

This scaffolds your basic loan blockchain. You will add the Loan module in the next chapter.

Change into the newly created loan directory.

```bash
cd loan
```

## Scaffold the Module

Scaffolding the module will create a new `loan` module inside the `x` directory.

```bash
starport scaffold module loan --dep bank
```

Use the `--dep` flag to specify that this module depends on and is going to interact with the `bank` module.

## Scaffold a List

The [scaffold list](https://docs.starport.com/cli/#starport-scaffold-list) command creates data stored as an array. 
In such a list, you want to store previous mentioned `Loan` proto message.

```bash
starport scaffold list loan amount fee collateral deadline state borrower lender --no-message
```

Use the `--no-message` flag to disable CRUD messages in the scaffold.

The data you store in an array are the loans, with these parameters.

See the `Loan` message in proto/loan/loan.proto

```proto
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

You will define the messages to interact with the loan list in the coming chapters.

Now it is time to interact with the Loan with messages.
But before, make sure to store your current state in a git commit.

```bash
git add .
git commit -m "Scaffold loan module and loan list"
```

## Scaffold the Messages

In order to create a loan app, you will need the following messages:

* Request Loan
* Approve Loan
* Repay Loan
* Liquidate Loan
* Cancel Loan

You can use the `starport scaffold message` command to create each of the messages.
You will learn the details of each of the messages and how to scaffold them in this chapter.

Create the messages one after the other with the according business logic.

### Request Loan

For a loan the initial message to start is with a cosmonaut requesting a loan.
The cosmonaut wants a certain `amount` and is willing to pay `fees` as well as give a `collateral`. The `deadline` marks the time when the loan has to be repayed.

The first message is the `request-loan` message. It needs the input parameters `amount`, `fee`, `collateral` and `deadline`.

```bash
starport scaffold message request-loan amount fee collateral deadline
```

For sake of simplicity every parameter will remain a string.

The `request-loan` message should create a new loan object and lock the tokens to be spent as fee and collateral into an escrow account.
This has to be described in the modules keeper directory `x/loan/keeper/msg_server_request_loan.go`

```go
package keeper

import (
	"context"

	"github.com/cosmonaut/loan/x/loan/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RequestLoan(goCtx context.Context, msg *types.MsgRequestLoan) (*types.MsgRequestLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
    // Create a new Loan with the according user input
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

Since this function is using the `bankKeeper` with a function `SendCoinsFromAccountToModule`, this function needs to be added to the `x/loan/types/expected_keepers.go`. This is how it should be added.

```go
package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}
```

When creating a loan you will want to require a certain input validation and throw error messages in case the user tries impossible inputs.

You can describe message validation errors in the modules `types` directory.

Add the following code to the `/x/loan/types/message_request_loan.go` function `ValidateBasic()`

```go
func (msg *MsgRequestLoan) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)

	amount, _ := sdk.ParseCoinsNormalized(msg.Amount)
	fee, _ := sdk.ParseCoinsNormalized(msg.Fee)
	collateral, _ := sdk.ParseCoinsNormalized(msg.Collateral)

	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if !amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "amount is not a valid Coins object")
	}
	if amount.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "amount is empty")
	}
	if !fee.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "fee is not a valid Coins object")
	}
	if fee.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "fee is empty")
	}
	if !collateral.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collateral is not a valid Coins object")
	}
	if collateral.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collateral is empty")
	}
	return nil
}
```

This concludes the `request-loan` message.

You can run the chain and test your first message.

Start the blockchain

```bash
starport chain serve
```

Add your first loan

```bash
loand tx loan request-loan 100token 2token 200token 500 --from alice -y
```

Query your loan

```bash
loand query loan list-loan
```

You should see the first Loan in the list

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

You can stop the blockchain again with ctrl + c.

A good time to add your advancements to git.

```bash
git add .
git commit -m "Add request-loan message"
```

### Approve Loan

After a loan request has been published, another account can approve the loan and agree to the terms of the borrower.
The message `approve-loan` has one parameter, the `id`.
Specify the type of `id` as `uint`, by default IDs are stored as `uint`.

```bash
starport scaffold message approve-loan id:uint
```

This message should be available for all loan types that are in the status "requested".

It would send the requested coins for the loan to the borrower and set the loan state to "approved".

Modify the `x/loan/keeper/msg_server_approve_loan.go` to implement this logic.

```go
package keeper

import (
	"context"
	"fmt"

	"github.com/cosmonaut/loan/x/loan/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) ApproveLoan(goCtx context.Context, msg *types.MsgApproveLoan) (*types.MsgApproveLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	loan, found := k.GetLoan(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// TODO: for some reason the error doesn't get printed to the terminal
	if loan.State != "requested" {
		return nil, sdkerrors.Wrapf(types.ErrWrongLoanState, "%v", loan.State)
	}

	lender, _ := sdk.AccAddressFromBech32(msg.Creator)
	borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
	amount, _ := sdk.ParseCoinsNormalized(loan.Amount)

	k.bankKeeper.SendCoins(ctx, lender, borrower, amount)

	loan.Lender = msg.Creator
	loan.State = "approved"

	k.SetLoan(ctx, loan)

	return &types.MsgApproveLoanResponse{}, nil
}
```

This module uses the bankKeepers SendCoins function. Add this to the `x/loan/types/expected_keepers.go` accordingly

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

There is also introduced a new error type `ErrWrongLoanState`.
Add this to the errors definitions in `x/loan/types/errors.go`

```go
package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/loan module sentinel errors
var (
	ErrWrongLoanState = sdkerrors.Register(ModuleName, 1, "wrong loan state")
)
```

Start the blockchain and use the two commands you already have available.

```bash
starport chain serve -r
```

Use the `-r` flag to reset the blockchain state and start with a new database.

```bash
loand tx loan request-loan 100token 2token 200token 500 --from bob -y
```

Query your loan request

```bash
loand query loan list-loan
```

Approve the loan

```bash
loand tx loan approve-loan 0 --from alice -y
```

This should send the balances according to the loan request.
CHeck for the loan list again. You should see the state now approved.

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

You can query for alices balances to see the loan in effect.
Take the lender address from above, this is alice address.

```bash
loand query bank balances <alice_address>
```

In case everything works as expected, this is a good time to save the state with a git commit.

```bash
git add .
git commit -m "Add approve loan message"
```

### Repay Loan

After the loan has been approved, the cosmonaut must be able to repay an approved loan.
Scaffold the message `repay-loan` that is used by a borrower to return tokens borrowed from the lender.

```bash
starport scaffold message repay-loan id:uint
```

Repaying a loan requires the loan to be in the "approved" status.

The coins as described in the loan are collected and sent from the borrower to the lender, as well as the agreed fees.
The collateral will be released from the escrow module account.

This logic is defined in the `x/loan/keeper/msg_server_repay_loan.go`.

```go
package keeper

import (
	"context"
	"fmt"

	"github.com/cosmonaut/loan/x/loan/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	amount, _ := sdk.ParseCoinsNormalized(loan.Amount)
	fee, _ := sdk.ParseCoinsNormalized(loan.Fee)
	collateral, _ := sdk.ParseCoinsNormalized(loan.Collateral)

	k.bankKeeper.SendCoins(ctx, borrower, lender, amount)
	k.bankKeeper.SendCoins(ctx, borrower, lender, fee)
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, borrower, collateral)

	loan.State = "repayed"

	k.SetLoan(ctx, loan)

	return &types.MsgRepayLoanResponse{}, nil
}
```

After the coins have been successfully exchanged, the state of the loan will be set to `repayed`.

Releasing tokens with the `bankKeepers` `SendCoinsFromModuleToAccount` function, you will need to add this to the `x/loan/types/expected_keepers.go`

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

Start the blockchain and use the two commands you already have available.

```bash
starport chain serve -r
```

Use the `-r` flag to reset the blockchain state and start with a new database.

```bash
loand tx loan request-loan 100token 2token 200token 500 --from bob -y
```

Query your loan request

```bash
loand query loan list-loan
```

Approve the loan

```bash
loand tx loan approve-loan 0 --from alice -y
```

You can query for alices balances to see the loan in effect.
Take the lender address from above, this is alice address.

```bash
loand query bank balances <alice_address>
```

Now repay the loan

```bash
loand tx loan repay-loan 0 --from bob -y
```

The loan should now be status `repayed`

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

And alice balance reflect the repayed amount plus fees

```bash
loand query bank balances <alice_address>
```

Good job!
Update your git with the changes you made.

```bash
git add .
git commit -m "Add repay-loan message"
```

### Liquidate Loan

A lender can liquidate a loan when the borrower does not pay the tokens back after the passed deadline. The message to `liquidate-loan` refers to the loan `id`.

```bash
starport scaffold message liquidate-loan id:uint
```

The `liquidate-loan` message should be able to be executed by the `lender`.
The status of the loan has to be `approved`. The `deadline` has to be met.

When these properties are valid, the collateral shall be liquidated from the `borrower`.

Add this to the `keeper` in `x/loan/keeper/msg_server_liquidate_loan.go`

```go
package keeper

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmonaut/loan/x/loan/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

Add the new error `ErrDeadline` to the error messages in `x/loan/types/errors.go`

```go
package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/loan module sentinel errors
var (
	ErrWrongLoanState = sdkerrors.Register(ModuleName, 1, "wrong loan state")
	ErrDeadline       = sdkerrors.Register(ModuleName, 2, "deadline")
)
```

These are all changes necessary to the `liquidate-loan` message.

You can test the liquidation message now.

```bash
starport chain serve -r
```

Set the deadline for the loan request to 1 block.

```bash
loand tx loan request-loan 100token 2token 200token 1 --from bob -y
```

Query your loan request

```bash
loand query loan list-loan
```

Approve the loan

```bash
loand tx loan approve-loan 0 --from alice -y
```

You can query for alices balances to see the loan in effect.
Take the lender address from above, this is alice address.

```bash
loand query bank balances <alice_address>
```

Now repay the loan

```bash
loand tx loan liquidate-loan 0 --from alice -y
```

The loan should now be status `liquidated`

```bash
loand query loan list-loan
```

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

And alice balance reflect the repayed amount plus fees

```bash
loand query bank balances <alice_address>
```

Add the changes to your git.

```bash
git add .
git commit -m "Add liquidate-loan message"
```

### Cancel Loan

After a loan request has been made and not been approved, the `borrower` should be able to cancel a loan request. Scaffold the message for `cancel-loan`.

```bash
starport s message cancel-loan id:uint
```

Only the `borrower` should be able to cancel a loan request.
The state of the request must be `requested`.
Then the collateral coins can be released from escrow and the status set to `cancelled`.

Add this to the `keeper` in `x/loan/keeper/msg_server_cancel_loan.go`.

```go
package keeper

import (
	"context"
	"fmt"

	"github.com/cosmonaut/loan/x/loan/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

```bash
starport chain serve -r
```

```bash
loand tx loan request-loan 100token 2token 200token 100 --from bob -y
```

Query your loan request

```bash
loand query loan list-loan
```

```bash
loand tx loan cancel-loan 0 --from bob -y
```

Query your loan request

```bash
loand query loan list-loan
```

Then the collateral coins can be released from escrow and the status set to `cancelled`.

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

Congratulations. This concludes the loan module.

Consider adding this to your git commit and maybe publish it on a public repository for others to see your accomplisments.

```bash
git add .
git commit -m "Add cancel-loan message"
```

You have learned how to

* Scaffold a blockchain
* Scaffold a module
* Scaffold a list for loan objects
* Create messages in your module to interact with the loan object
* Interact with other modules in your module
* Use an escrow module account
* Add application messages for a Loan system
  * Request Loan
  * Approve Loan
  * Repay Loan
  * Liquidate Loan
  * Cancel Loan
