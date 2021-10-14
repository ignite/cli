---
description: Loan module using Starport
order: x
---

# Creating a Loan Blockchain in Go

`Loan` is a blockchain built using Cosmos SDK and created with [Starport](https://github.com/tendermint/starport)

In this tutorial you will learn how to create, approve and repay loan request. Also, liquidate unpaid loan.

### Borrower:
A borrower will post a loan request with information such as - loan amount, fees, collateral and deadline.
Borrower will repay the loan transfer amount and fee to lender or risk the loosing of collateral.

### Lender:
A lender can approve loan request from borrower. Approving the loan transfers the loan amount to the the borrower. If the borrower is unable to pay the loan, the lender can liquidate the loan which transfers the collateral to the lender.


## Creating a loan blockchain using starport

```bash
starport scaffold chain github.com/cosmonaut/loan --no-module
```

This command creates a Cosmos SDK blockchain called loan in a `loan` directory. The source code inside the `loan` directory contains a fully functional ready-to-use blockchain.


```bash
cd loan
```

```bash
starport scaffold module loan --dep bank
```

<Some information to be added.>



```bash
starport s list loan amount fee collateral deadline state borrower lender --no-message
```

This commands creates CRUD opertaion for loan along with 


```bash
starport s message request-loan amount fee collateral deadline
```

`Request-loan` is a message which request for loan after passing in parameters like amount, fee, collateral and deadline. You also need borrower, which is passed as signer while creating the tx on Blockchain.


```bash
starport s message approve-loan id:uint
```

`Approve-loan` is a message used by lender which needs only 1 parameter: id. We also pass type - that is unsigned integer (uint) to save conversion time from string to uint


```bash
starport s message repay-loan id:uint
```

Repay-loan is a message used by borrower to return the money which needs only 1 parameter: id. We also pass type - that is unsigned integer (uint) to save conversion time from string to uint


```bash
starport s message liquidate-loan id:uint
```

`Liquidate-loan` is a message used by lender to liquidate the loan in case of loan not payed by borrower


```bash
starport s message cancel-loan id:uint
```

`Cancel-loan` is a message used by borrower to cancel loan request after making request and submitting collateral


### Now start adding the following code to `keeper` to handle each function.


Add following code to `keeper/msg_server_request_loan.go`

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

The first step is to deconstruct the loan message into loan types. We start filling in the value in types like Amount, Fee, Collateral, etc from  messages.

The second step is to make state transitions. You need to transfer collateral from the borrower to the module account for which we get borrower's address.


The third step is to convert collateral. ParseCoinsNormalized will parse out coins and normalize it. 

The fourth step is to use functionality from the module bankkeeper to send coins. 

The last step is to append loan. Starport has generated a functionality to append loan which can be found under `keeper/loan.go`

