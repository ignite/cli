---
description: Loan blockchain using Starport
order: x
---

# Creating a Loan Blockchain in Go

`Loan` is a blockchain built using Cosmos SDK and created with [Starport](https://github.com/tendermint/starport)

In this tutorial you will learn how to create, approve and repay loan request. Also, liquidate unpaid loan.

### Borrower:
A borrower will post a loan request with information such as - loan amount, fees, collateral and deadline.
Borrower will repay the loan transfer amount and fee to lender or risk loosing the collateral.

### Lender:
A lender can approve loan request from borrower. Approving the loan transfers the loan amount to the the borrower. If the borrower is unable to pay the loan, the lender can liquidate the loan which transfers the collateral to the lender.


## Creating a loan blockchain using starport

```bash
starport scaffold chain github.com/cosmonaut/loan --no-module
```

This command creates a Cosmos SDK blockchain called loan in a `loan` directory. The source code inside the `loan` directory contains a fully functional ready-to-use blockchain. Use `--no-module` to disable scaffolding an empty module: you'll scaffold a module in the next step.


```bash
cd loan
```

```bash
starport scaffold module loan --dep bank
```

Scaffold a new module called `loan`. Use the `--dep` flag to specify that this module depends on the standard `bank` module.


```bash
starport s list loan amount fee collateral deadline state borrower lender --no-message
```

This command creates CRUD operation for `loan` stored as a list. `--no-message` disables CRUD interaction messages in scaffolding.


```bash
starport s message request-loan amount fee collateral deadline
```

`request-loan` is a message which lets a borrower to request a loan. This message accepts 4 parameters: amount, fee, collateral and deadline.


```bash
starport s message approve-loan id:uint
```

`approve-loan` is a message used by a lender to approve an already requested loan. The message contains 1 field: `id`. Specify the type of `id` as `uint`, because by default IDs are stored as `uint`s, 


```bash
starport s message repay-loan id:uint
```

`repay-loan` is a message used by a borrower to return tokens borrowed from the lender. 


```bash
starport s message liquidate-loan id:uint
```

`liquidate-loan` is a message used by lender to liquidate the loan in case of loan not payed by borrower.


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
