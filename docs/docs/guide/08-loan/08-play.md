# Play

Add `10000foocoin` to Alice's account. These tokens will be used as a loan collateral.

```yml
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

## Repaying a loan

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

```
loand tx loan approve-loan 0 --from bob
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
  # highlight-start
  lender: cosmos1qfzpxfhsu2qfy2exkukuanrkzrrexh9yeg2pr4
  state: approved
  # highlight-end
```

```
loand q bank balances $(loand keys show alice -a)
```

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

```yml
balances:
- amount: "100000000"
  denom: stake
  # highlight-start
- amount: "9000"
  denom: token
  # highlight-end
```

```
loand tx loan repay-loan 0 --from alice
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
  lender: cosmos1qfzpxfhsu2qfy2exkukuanrkzrrexh9yeg2pr4
  # highlight-next-line
  state: repayed
```

```
loand q bank balances $(loand keys show alice -a)
```

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

```
loand tx loan request-loan 1000token 100token 1000foocoin 20 --from alice
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

```
loand tx loan approve-loan 0 --from bob
```

```
loand tx loan liquidate-loan 1 --from bob
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
q bank balances $(loand keys show alice -a)
```

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