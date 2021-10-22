---
description: Loan blockchain using Starport
order: 6
title: "Advanced module: Loan"
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
* Add application logic for a Loan system
  * Request
  * Approve
  * Repay
  * Liquidate
  * Cancel Loan

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

Use the `--dep` flag to specify that this module depends on the `bank` module.

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

While scaffolded as strings, for easier convertion later, you would want `amount`, `fee` and `collateral` as a Coins object.

Modify the `loan.proto` to import cosmos coin proto and modify the types to be Coins. While the deadline should be a uint64.

```proto
syntax = "proto3";
package cosmonaut.loan.loan;

option go_package = "github.com/cosmonaut/loan/x/loan/types";
import "gogoproto/gogo.proto";
import "cosmos/base/v1beta1/coin.proto";

message Loan {
  uint64 id = 1;
  repeated cosmos.base.v1beta1.Coin amount = 2 [(gogoproto.nullable) = false];
  repeated cosmos.base.v1beta1.Coin fee = 3 [(gogoproto.nullable) = false];
  repeated cosmos.base.v1beta1.Coin collateral = 4 [(gogoproto.nullable) = false];
  uint64 deadline = 5;
  string state = 6;
  string borrower = 7;
  string lender = 8;
}
```

After this change, you should run a proto build with Starport.

```bash
starport generate proto-go
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

On requesting the loan, the `amount`, `fee` and `collateral` parameters should be described as `Coins`.
First, add the imports to the `proto/loan/tx.proto` file 

```proto
import "cosmos/base/v1beta1/coin.proto";
import "gogoproto/gogo.proto";
```

and use the following `MsgRequestLoan` message, using Coin as input for the three parameters.

```proto
message MsgRequestLoan {
  string creator = 1;
  repeated cosmos.base.v1beta1.Coin amount = 2 [(gogoproto.nullable) = false];
  repeated cosmos.base.v1beta1.Coin fee = 3 [(gogoproto.nullable) = false];
  repeated cosmos.base.v1beta1.Coin collateral = 4 [(gogoproto.nullable) = false];
  uint64 deadline = 5;
}
```

After changing a proto file, it is recommended to run a proto build with Starport.

```bash
starport generate proto-go
```

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

  // Add the input into a loan type
	loan := types.Loan{
		Amount:     msg.Amount,
		Fee:        msg.Fee,
		Collateral: msg.Collateral,
		Deadline:   msg.Deadline,
		State:      "requested",
		Borrower:   msg.Creator,
	}

	// TODO: collateral has to be more than the amount (+fee?)
  // Get borrowers Account
	borrower, _ := sdk.AccAddressFromBech32(msg.Creator)
  // Send the Collateral into a module escrow account
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, borrower, types.ModuleName, msg.Collateral)
	if err != nil {
		return nil, err
	}
  // add the created loan to the keeper
	loan.Id = k.AppendLoan(ctx, loan)

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
Add the following code to the `/x/loan/types/message_request_loan.go` function `func (msg *MsgRequestLoan) ValidateBasic()`

```go
func (msg *MsgRequestLoan) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if !sdk.Coins(msg.Amount).IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "amount is not a valid Coins object")
	}
	if sdk.Coins(msg.Amount).Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "amount is empty")
	}
	if !sdk.Coins(msg.Fee).IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "fee is not a valid Coins object")
	}
	if sdk.Coins(msg.Fee).Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "fee is empty")
	}
	if !sdk.Coins(msg.Collateral).IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collateral is not a valid Coins object")
	}
	if sdk.Coins(msg.Collateral).Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collateral is empty")
	}
	return nil
}
```

This concludes the `request-loan` message.

A good time to add your advancements to git.

```bash
git add .
git commit -m "Add request-loan message"
```

### Approve Loan

After a loan request has been published, another account an approve the loan and agree to the terms of the borrower.
`approve-loan` contains 1 field, the `id`.
Specify the type of `id` as `uint`, by default IDs are stored as `uint`.

```bash
starport scaffold message approve-loan id:uint
```

### Repay Loan

After the loan has been approved, the cosmonaut must be able to repay an approved loan.
Scaffold the message `repay-loan` that is used by a borrower to return tokens borrowed from the lender. 

```bash
starport s message repay-loan id:uint
```

### Liquidate Loan

```bash
starport s message liquidate-loan id:uint
```

`liquidate-loan` is a message used by lender to liquidate the loan in case of loan not payed by borrower.

### Cancel Loan

```bash
starport s message cancel-loan id:uint
```

`cancel-loan` is a message used by a borrower to cancel a loan request after making request and submitting collateral.


## Start adding the following code to `keeper` to handle each function.


### Add following code to `keeper/msg_server_request_loan.go`

```go
// Add import:
import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: Handling the message
	var loan = types.Loan{
		Amount:     msg.Amount,
		Fee:        msg.Fee,
		Collateral: msg.Collateral,
		Deadline:   msg.Deadline,
		State:      "requested",
		Borrower:   msg.Creator,
	}

	
	borrower, _ := sdk.AccAddressFromBech32(msg.Creator)

	collateral, err := sdk.ParseCoinsNormalized(loan.Collateral)
	if err != nil {
		panic(err)
	}

	sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, borrower, types.ModuleName, collateral)
	if sdkError != nil {
		return nil, sdkError
	}

	k.AppendLoan(
		ctx,
		loan,
	)
```

The functionality of this module is to allow people to make loan request.

The first step is to deconstruct the loan message into loan types. Start filling in the value in types like Amount, Fee, Collateral, etc from  messages.

The second step is to make state transitions. You need to transfer collateral from the borrower to the module account for which we get borrower's address.

The third step is to convert collateral. `ParseCoinsNormalized` will parse out coins and normalize it. 

The fourth step is to use functionality from the module bankkeeper to send coins. 

The last step is to append loan. Starport has generated a functionality to append loan which can be found under `keeper/loan.go`


### Add following code to `keeper/msg_server_approve_loan.go`

```go
// Add import:
import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO

loan, found := k.GetLoan(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

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
```

The functionality of this module is to allow lender to approve loan request.

The first step is to get loan using the keeper function `GetLoan` before it can be approved.

The second step is to make sure only loans that are requested are approved and not already approved loans.

The third step is to populate values of lender, borrower and amount.

The fourth step is to send coins and change the state to `approved`

The last step is to set loan. Starport has generated a functionality to set loan which can be found under `keeper/loan.go`


### Add following code to `keeper/msg_server_repay_loan.go`

```go
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
```

The functionality of this module is to allow the borrower to repay loan.

The first step is to get loan using the keeper function `GetLoan` before it can be repayed.

The second step is to make sure only loans that are approved are repayed and not the pending loans.

The third step is to populate values of lender, borrower, amount, fee and collateral.

The fourth step is to send coins (loan amount and fees) to borrower.

The fifth step is to send the collateral amount to the borrower after the loan amount is repayed.

The last step is to change the state to `repayed` and set loan. Starport has generated a functionality to set loan which can be found under `keeper/loan.go`


### Add following code to `keeper/msg_server_liquidate_loan.go`

```go
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
```

The functionality of this module is to allow the lender to liquidate loan if unpaid past deadline.

The first step is to get loan using the keeper function `GetLoan` before it can be repayed.

The second step is to make sure only loans that are approved are liquidated and not the pending loans.

The third step is to populate values of lender and collateral.

The fourth step is to get loan deadline and compare it with block height. If its past the block height the collateral can be liquidated.

The fifth step is to send the collateral amount to the lender after the collateral is liquidated.

The last step is to change the state to `liquidated` and set loan. Starport has generated a functionality to set loan which can be found under `keeper/loan.go`


### Add following code to `keeper/msg_server_cancel_loan.go`

```go
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
```

The functionality of this module is to allow the borrower to cancel the loan request.

The first step is to check if the loan exist.

The second step is to make sure the borrower can cancel only its loan.

The third step is to check state of loan which should be requested and not approved or liquidated.

The fourth step is to fetch values of borrower and collateral. Then send collateral back to borrower.

The last step is to change the state to `cancelled` and set loan. Starport has generated a functionality to set loan which can be found under `keeper/loan.go`



## Running the Blockchain

Run your loan blockchain `starport chain serve`


### Request loan

```bash
loand tx loan request-loan [amount] [fee] [collateral] [deadline] [flags]
```

```markdown
loand tx loan request-loan 100token 2token 200token 500 --from alice -y
```

Where:  
--from is the name or address of private key with which to sign
-y is to skip tx broadcasting prompt confirmation

You should see an output similar to:

```bash
code: 0
codespace: ""
data: 0A250A232F636F736D6F6E6175742E6C6F616E2E6C6F616E2E4D7367526571756573744C6F616E
gas_used: "57234"
gas_wanted: "200000"
height: "442"
info: ""
logs:
- events:
  - attributes:
    - key: receiver
      value: cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l
    - key: amount
      value: 200token
    type: coin_received
  - attributes:
    - key: spender
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: amount
      value: 200token
    type: coin_spent
  - attributes:
    - key: action
      value: RequestLoan
    - key: sender
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    type: message
  - attributes:
    - key: recipient
      value: cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l
    - key: sender
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: amount
      value: 200token
    type: transfer
  log: ""
  msg_index: 0
raw_log: '[{"events":[{"type":"coin_received","attributes":[{"key":"receiver","value":"cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l"},{"key":"amount","value":"200token"}]},{"type":"coin_spent","attributes":[{"key":"spender","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"amount","value":"200token"}]},{"type":"message","attributes":[{"key":"action","value":"RequestLoan"},{"key":"sender","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l"},{"key":"sender","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"amount","value":"200token"}]}]}]'
timestamp: ""
tx: null
txhash: E2F12B96991FD15ECA93E373C66056D41DCE1B1C0DD33A09177F36D5F5566D94
```

This can also be checked using `query` loan function.

```bash
loand query loan list-loan
```

This returns a list of all loans. 

You should see an output similar to:

```bash
Loan:
- amount: 100token
  borrower: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
  collateral: 200token
  deadline: "500"
  fee: 2token
  id: "0"
  lender: ""
  state: requested
- amount: 100token
  borrower: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
  collateral: 200token
  deadline: "500"
  fee: 2token
  id: "1"
  lender: ""
  state: requested
```


### Approve loan

```bash
loand tx loan approve-loan [id] [flags]
```

```markdown
loand tx loan approve-loan 0 --from alice -y
```

You should see an output similar to:

```bash
code: 0
codespace: ""
data: 0A250A232F636F736D6F6E6175742E6C6F616E2E6C6F616E2E4D7367417070726F76654C6F616E
gas_used: "55050"
gas_wanted: "200000"
height: "828"
info: ""
logs:
- events:
  - attributes:
    - key: receiver
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: amount
      value: 100token
    type: coin_received
  - attributes:
    - key: spender
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: amount
      value: 100token
    type: coin_spent
  - attributes:
    - key: action
      value: ApproveLoan
    - key: sender
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    type: message
  - attributes:
    - key: recipient
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: sender
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: amount
      value: 100token
    type: transfer
  log: ""
  msg_index: 0
raw_log: '[{"events":[{"type":"coin_received","attributes":[{"key":"receiver","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"amount","value":"100token"}]},{"type":"coin_spent","attributes":[{"key":"spender","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"amount","value":"100token"}]},{"type":"message","attributes":[{"key":"action","value":"ApproveLoan"},{"key":"sender","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"sender","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"amount","value":"100token"}]}]}]'
timestamp: ""
tx: null
txhash: F1B52A2BB721529C244A2AAAFA77554D773B3D75D274EEEBA4680EB94840408E
```


Check the state of the loan using the following command:

```bash
loand query loan show-loan 0
```

This returns the loan requested by id.

You should see an output similar to:

```bash
Loan:
  amount: 100token
  borrower: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
  collateral: 200token
  deadline: "500"
  fee: 2token
  id: "0"
  lender: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
  state: approved
```

Note: The state has changed from `requested` to `approved`


### Repay loan

```bash
loand tx loan repay-loan [id] [flags]
```

```markdown
loand tx loan repay-loan 0 --from alice -y
```

You should see an output similar to:

```bash
code: 0
codespace: ""
data: 0A230A212F636F736D6F6E6175742E6C6F616E2E6C6F616E2E4D736752657061794C6F616E
gas_used: "74693"
gas_wanted: "200000"
height: "1167"
info: ""
logs:
- events:
  - attributes:
    - key: receiver
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: amount
      value: 100token
    - key: receiver
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: amount
      value: 2token
    - key: receiver
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: amount
      value: 200token
    type: coin_received
  - attributes:
    - key: spender
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: amount
      value: 100token
    - key: spender
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: amount
      value: 2token
    - key: spender
      value: cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l
    - key: amount
      value: 200token
    type: coin_spent
  - attributes:
    - key: action
      value: RepayLoan
    - key: sender
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: sender
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: sender
      value: cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l
    type: message
  - attributes:
    - key: recipient
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: sender
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: amount
      value: 100token
    - key: recipient
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: sender
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: amount
      value: 2token
    - key: recipient
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: sender
      value: cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l
    - key: amount
      value: 200token
    type: transfer
  log: ""
  msg_index: 0
raw_log: '[{"events":[{"type":"coin_received","attributes":[{"key":"receiver","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"amount","value":"100token"},{"key":"receiver","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"amount","value":"2token"},{"key":"receiver","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"amount","value":"200token"}]},{"type":"coin_spent","attributes":[{"key":"spender","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"amount","value":"100token"},{"key":"spender","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"amount","value":"2token"},{"key":"spender","value":"cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l"},{"key":"amount","value":"200token"}]},{"type":"message","attributes":[{"key":"action","value":"RepayLoan"},{"key":"sender","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"sender","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"sender","value":"cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"sender","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"amount","value":"100token"},{"key":"recipient","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"sender","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"amount","value":"2token"},{"key":"recipient","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"sender","value":"cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l"},{"key":"amount","value":"200token"}]}]}]'
timestamp: ""
tx: null
txhash: F84F0E7DE78BD9BBD34B0BCC538F83AC74574EA7FFD158F7AB720529FC1F989B
```

Check the state of the loan using the following command:

```bash
loand query loan show-loan 0
```

You should see an output similar to:

```bash
Loan:
  amount: 100token
  borrower: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
  collateral: 200token
  deadline: "500"
  fee: 2token
  id: "0"
  lender: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
  state: repayed
```

Note: The state has changed from `approved` to `repayed`


### Cancel loan

```bash
loand tx loan cancel-loan [id] [flags]
```

```markdown
loand tx loan cancel-loan 1 --from alice -y
```

You should see an output similar to:

```bash
code: 0
codespace: ""
data: 0A240A222F636F736D6F6E6175742E6C6F616E2E6C6F616E2E4D736743616E63656C4C6F616E
gas_used: "53569"
gas_wanted: "200000"
height: "1707"
info: ""
logs:
- events:
  - attributes:
    - key: receiver
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: amount
      value: 200token
    type: coin_received
  - attributes:
    - key: spender
      value: cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l
    - key: amount
      value: 200token
    type: coin_spent
  - attributes:
    - key: action
      value: CancelLoan
    - key: sender
      value: cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l
    type: message
  - attributes:
    - key: recipient
      value: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
    - key: sender
      value: cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l
    - key: amount
      value: 200token
    type: transfer
  log: ""
  msg_index: 0
raw_log: '[{"events":[{"type":"coin_received","attributes":[{"key":"receiver","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"amount","value":"200token"}]},{"type":"coin_spent","attributes":[{"key":"spender","value":"cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l"},{"key":"amount","value":"200token"}]},{"type":"message","attributes":[{"key":"action","value":"CancelLoan"},{"key":"sender","value":"cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l"}]},{"type":"transfer","attributes":[{"key":"recipient","value":"cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0"},{"key":"sender","value":"cosmos1gu4m79yj8ch8em7c22vzt3qparg69ymm75qf6l"},{"key":"amount","value":"200token"}]}]}]'
timestamp: ""
tx: null
txhash: 8AE8A3A9F502ECB6A3747B445FA8BB63FFBFFC4A1EF15DA9E678D08B8EC03913
```

Check the state of the loan using the following command:

```bash
loand query loan show-loan 1
```

You should see an output similar to:

```bash
Loan:
  amount: 100token
  borrower: cosmos1ulk2f49lhljvldqw09queq82pphsuv759t32t0
  collateral: 200token
  deadline: "500"
  fee: 2token
  id: "1"
  lender: ""
  state: cancelled
```

Note: The state has changed from `approved` to `cancelled`


Congratulations, you have just created a `loan blockchain` using starport.
