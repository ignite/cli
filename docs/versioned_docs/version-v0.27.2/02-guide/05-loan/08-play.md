# Play

Add `10000foocoin` to Alice's account. These tokens will be used as a loan
collateral.

```yml title="config.yml"
version: 1
accounts:
  - name: alice
    coins:
      - 20000token
      # highlight-next-line
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

Start a blockchain node:

```
ignite chain serve
```

## Repaying a loan

Request a loan of `1000token` with `100token` as a fee and `1000foocoin` as a
collateral from Alice's account. The deadline is set to `500` blocks:

```
loand tx loan request-loan 1000token 100token 1000foocoin 500 --from alice
```

```
loand q loan list-loan
```

```yml
Loan:
- amount: 1000token
  borrower: cosmos153dk8qh56v4yg6e4uzrvvqjueu6d36fptlr2kw
  collateral: 1000foocoin
  deadline: "500"
  fee: 100token
  id: "0"
  lender: ""
  state: requested
```

Please be aware that the addresses displayed in your terminal window (such as those in the `borrower` field) will not match the ones provided in this tutorial. This is because Ignite generates new accounts each time a chain is started, unless an account has a mnemonic specified in the `config.yml` file.

Approve the loan from Bob's account:

```
loand tx loan approve-loan 0 --from bob
```

```
loand q loan list-loan         
```

The `lender` field has been updated to Bob's address and the `state` field has
been updated to `approved`:

```yml        
Loan:
- amount: 1000token
  borrower: cosmos153dk8qh56v4yg6e4uzrvvqjueu6d36fptlr2kw
  collateral: 1000foocoin
  deadline: "500"
  fee: 100token
  id: "0"
  # highlight-start
  lender: cosmos1qfzpxfhsu2qfy2exkukuanrkzrrexh9yeg2pr4
  state: approved
  # highlight-end
```

```
loand q bank balances $(loand keys show alice -a)
```

The `foocoin` balance has been updated to `9000`, because `1000foocoin` has been
transferred as collateral to the module account. The `token` balance has been
updated to `21000`, because `1000token` has been transferred from Bob's account
to Alice's account as a loan:

```yml
balances:
  # highlight-start
- amount: "9000"
  denom: foocoin
  # highlight-end
- amount: "100000000"
  denom: stake
  # highlight-start
- amount: "21000"
  denom: token
  # highlight-end
```

```
loand q bank balances $(loand keys show bob -a)  
```

The `token` balance has been updated to `9000`, because `1000token` has been
transferred from Bob's account to Alice's account as a loan:

```yml
balances:
- amount: "100000000"
  denom: stake
  # highlight-start
- amount: "9000"
  denom: token
  # highlight-end
```

Repay the loan from Alice's account:

```
loand tx loan repay-loan 0 --from alice
```

```
loand q loan list-loan
```

The `state` field has been updated to `repayed`:

```yml
Loan:
- amount: 1000token
  borrower: cosmos153dk8qh56v4yg6e4uzrvvqjueu6d36fptlr2kw
  collateral: 1000foocoin
  deadline: "500"
  fee: 100token
  id: "0"
  lender: cosmos1qfzpxfhsu2qfy2exkukuanrkzrrexh9yeg2pr4
  # highlight-next-line
  state: repayed
```

```
loand q bank balances $(loand keys show alice -a)
```

The `foocoin` balance has been updated to `10000`, because `1000foocoin` has
been transferred from the module account to Alice's account. The `token` balance
has been updated to `19900`, because `1000token` has been transferred from
Alice's account to Bob's account as a repayment and `100token` has been
transferred from Alice's account to Bob's account as a fee:

```yml
balances:
  # highlight-start
- amount: "10000"
  denom: foocoin
  # highlight-end
- amount: "100000000"
  denom: stake
  # highlight-start
- amount: "19900"
  denom: token
  # highlight-end
```

```
loand q bank balances $(loand keys show bob -a)  
```

The `token` balance has been updated to `10100`, because `1000token` has been
transferred from Alice's account to Bob's account as a repayment and `100token`
has been transferred from Alice's account to Bob's account as a fee:

```yml
balances:
- amount: "100000000"
  denom: stake
  # highlight-start
- amount: "10100"
  denom: token
  # highlight-end
```

## Liquidating a loan

Request a loan of `1000token` with `100token` as a fee and `1000foocoin` as a
collateral from Alice's account. The deadline is set to `20` blocks. The
deadline is set to a very small value, so that the loan can be quickly
liquidated in the next step: 

```
loand tx loan request-loan 1000token 100token 1000foocoin 20 --from alice
```

```
loand q loan list-loan
```

A loan has been added to the list:

```yml
Loan:
- amount: 1000token
  borrower: cosmos153dk8qh56v4yg6e4uzrvvqjueu6d36fptlr2kw
  collateral: 1000foocoin
  deadline: "500"
  fee: 100token
  id: "0"
  lender: cosmos1qfzpxfhsu2qfy2exkukuanrkzrrexh9yeg2pr4
  state: repayed
  # highlight-start
- amount: 1000token
  borrower: cosmos153dk8qh56v4yg6e4uzrvvqjueu6d36fptlr2kw
  collateral: 1000foocoin
  deadline: "20"
  fee: 100token
  id: "1"
  lender: ""
  state: requested
  # highlight-end
```

Approve the loan from Bob's account:

```
loand tx loan approve-loan 1 --from bob
```

Liquidate the loan from Bob's account:

```
loand tx loan liquidate-loan 1 --from bob
```

```
loand q loan list-loan
```

The loan has been liquidated:

```yml
Loan:
- amount: 1000token
  borrower: cosmos153dk8qh56v4yg6e4uzrvvqjueu6d36fptlr2kw
  collateral: 1000foocoin
  deadline: "500"
  fee: 100token
  id: "0"
  lender: cosmos1qfzpxfhsu2qfy2exkukuanrkzrrexh9yeg2pr4
  state: repayed
  # highlight-start
- amount: 1000token
  borrower: cosmos153dk8qh56v4yg6e4uzrvvqjueu6d36fptlr2kw
  collateral: 1000foocoin
  deadline: "20"
  fee: 100token
  id: "1"
  lender: cosmos1qfzpxfhsu2qfy2exkukuanrkzrrexh9yeg2pr4
  state: liquidated
  # highlight-end
```

```
loand q bank balances $(loand keys show alice -a)
```

The `foocoin` balance has been updated to `9000`, because `1000foocoin` has been
transferred from Alice's account to the module account as a collateral. Alice
has lost her collateral, but she has kept the loan amount:

```yml
balances:
  # highlight-start
- amount: "9000"
  denom: foocoin
  # highlight-end
- amount: "100000000"
  denom: stake
  # highlight-start
- amount: "20900"
  denom: token
  # highlight-end
```

```
loand q bank balances $(loand keys show bob -a)  
```

The `foocoin` balance has been updated to `1000`, because `1000foocoin` has been
transferred from the module account to Bob's account as a collateral. Bob has
gained the collateral, but he has lost the loan amount:

```yml
balances:
  # highlight-start
- amount: "1000"
  denom: foocoin
  # highlight-end
- amount: "100000000"
  denom: stake
- amount: "9100"
  denom: token
```
