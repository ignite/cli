---
order: 3
---

# Using the exchange

## Order Book

Using the exchange starts from creating a order book for a pair of tokens:

```
interchanged tx ibcdex send-createPair [src-port] [src-channel] [sourceDenom] [targetDenom]

# Create pair broadcasted to the source blockchain
interchanged tx ibcdex send-createPair ibcdex channel-0 mcx vcx
```

A pair of token is defined by two denominations: source denom (in this example, `mcx`) and target denom (`vcx`). Creating an orderbook affects state on the source blockchain (to which the transaction was broadcasted) and the target blockchain. On the source blockchain `send-createPair` creates an empty sell order book and on the target blockchain a buy order book is created.

```yml
# Created a sell order book on the source blockchain
SellOrderBook:
- amountDenom: mcx
  creator: ""
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 0
  orders: []
  priceDenom: vcx
```

```yml
# Created a buy order book on the target blockchain
BuyOrderBook:
- amountDenom: mcx
  creator: ""
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 1
  orders: []
  priceDenom: vcx
```

To make it possible `createPair` first sends an IBC packet to the target chain. Upon receiving a packet the target chain creates a buy order book and sends back an acknowledgement to the source chain. Upon receiving an acknowledgement, the source chain creates a sell order book. Sending an IBC packet requires a user specify a port and a channel through which a packet will be transferred.

## Sell Order

Once an order book is created, the next step is to create a sell order:

```
interchanged tx ibcdex send-sellOrder [src-port] [src-channel] [amountDenom] [amount] [priceDenom] [price]

# Sell order broadcasted to the source blockchain
interchanged tx ibcdex send-sellOrder ibcdex channel-0 mcx 10 vcx 15
```

The `send-sellOrder` command broadcasts a message that locks tokens on the source blockchain and creates a new sell order on the source blockchain.

```yml
# Source blockchain
balances:
- amount: "990" # decreased from 1000
  denom: mcx
SellOrderBook:
- amountDenom: mcx
  creator: ""
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 2
  orders: # a new sell order is created
  - amount: 10
    creator: cosmos1v3p3j7c64c4ls32pcjct333e8vqe45gwwa289q
    id: 0
    price: 15
  priceDenom: vcx
```

## Buy Order

A buy order has the same set of arguments: amount of tokens to be purchased and a price.

```
`interchanged tx ibcdex send-buyOrder [src-port] [src-channel] [amountDenom] [amount] [priceDenom] [price]`

# Buy order broadcasted to the target blockchain
interchanged tx ibcdex send-buyOrder ibcdex channel-0 mcx 10 vcx 5
```

The `send-buyOrder` command locks tokens on the target blockchain and creates a buy order book on the target blockchain.

```yml
# Target blockchain
balances:
- amount: "950" # decreased from 1000
  denom: vcx
BuyOrderBook:
- amountDenom: mcx
  creator: ""
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 3
  orders: # a new buy order is created
  - amount: 10
    creator: cosmos1qlrz3peenc6s3xjv9k97e8ef72nk3qn3a0xax2
    id: 1
    price: 5
  priceDenom: vcx
```

## Performing an Exchange with a Sell Order

We now have two orders open for MCX: a sell order on the source chain (for 10mcx at 15vcx) and a buy order on the target chain (for 5mcx at 5vcx). Let's perform an exchange by sending a sell order to the source chain.

```
# Sell order broadcasted to the source chain
interchanged tx ibcdex send-sellOrder ibcdex channel-0 mcx 5 vcx 3
```

The sell order (for 5mcx at 3vcx) was filled on the target chain by the buy order. The amount of the buy order on the target chain has decreased by 5mcx.

```yml
# Target blockchain
BuyOrderBook:
- amountDenom: mcx
  creator: ""
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 5
  orders:
  - amount: 5 # decreased from 10
    creator: cosmos1qlrz3peenc6s3xjv9k97e8ef72nk3qn3a0xax2
    id: 3
    price: 5
  priceDenom: vcx
```

The sender of the filled sell order exchanged 5mcx for 25 vcx vouchers. 25 vouchers is a product of the amount of the sell order (5mcx) and price of the buy order (5vcx).

```yml
# Source blockchain
balances:
- amount: "25" # increased from 0
  denom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833 # vcx voucher
- amount: "985" # decreased from 990
  denom: mcx
```

The counterparty (sender of the buy mcx order) received 5 mcx vouchers. vcx balance hasn't changed, because the correct amount of vcx (50) were locked at the creation of the buy order during the previous step.

```yml
# Target blockchain
balances:
- amount: "5" # increased from 0
  denom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33 # mcx voucher
```

## Performing an Exchange with a Buy Order

An order is sent to buy 5mcx for 15vcx.

```
# Buy order broadcasted to the target chain
interchanged tx ibcdex send-buyOrder ibcdex channel-0 mcx 5 vcx 15
```

A buy order is immediately filled on the source chain and sell order creator recived 75 vcx vouchers. The sell order amount is decreased by the amount of the filled buy order (by 5mcx).

```yml
# Source blockchain
balances:
- amount: "100" # increased from 25
  denom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833 # vcx voucher
SellOrderBook:
- amountDenom: mcx
  creator: ""
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 4
  orders:
  - amount: 5 # decreased from 10
    creator: cosmos1v3p3j7c64c4ls32pcjct333e8vqe45gwwa289q
    id: 2
    price: 15
  priceDenom: vcx
```

Creator of the buy order received 5 mcx vouchers for 75 vcx (5mcx * 15vcx).

```yml
# Target blockchain
balances:
- amount: "10" # increased from 5
  denom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33 # mcx vouchers
- amount: "875" # decreased from 950
  denom: vcx
```

## Complete Exchange with a Partially Filled Sell Order

An order is sent to sell 10mcx for 3vcx.

```
# Source blockchain
interchanged tx ibcdex send-sellOrder ibcdex channel-0 mcx 10 vcx 3
```

The sell amount is 10mcx, but the opened buy order amount is only 5mcx. The buy order gets filled completely and removed from the order book. The author of the previously created buy order recived 10 mcx vouchers from the exchange.

```yml
# Target blockchain
balances:
- amount: "15" # increased from 5
  denom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33 # mcx voucher
BuyOrderBook:
- amountDenom: mcx
  creator: ""
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 5
  orders: [] # buy order with amount 5mcx has been closed
  priceDenom: vcx
```

The author of the sell order successfuly exchanged 5 mcx and recived 25 vcx vouchers. The other 5mcx created a sell order.

```yml
# Source blockchain
balances:
- amount: "125" # increased from 100
  denom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833 # vcx vouchers
- amount: "975" # decreased from 985
  denom: mcx
- amountDenom: mcx
SellOrderBook:
  creator: ""
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 6
  orders:
  - amount: 5 # hasn't changed
    creator: cosmos1v3p3j7c64c4ls32pcjct333e8vqe45gwwa289q
    id: 2
    price: 15
  - amount: 5 # new order is created
    creator: cosmos1v3p3j7c64c4ls32pcjct333e8vqe45gwwa289q
    id: 4
    price: 3
```

## Complete Exchange with a Partially Filled Buy Order

An order is created to buy 10 mcx for 5 vcx.

```
# Target blockchain
interchanged tx ibcdex send-buyOrder ibcdex channel-0 mcx 10 vcx 5
```

The buy order is partially filled for 5mcx. An existing sell order for 5 mcx (with a price of 3 vcx) on the source chain has been completely filled and removed from the order book. The author of the closed sell order recived 15 vcx vouchers (product of 5mcx and 3vcx).

```yml
# Source blockchain
balances:
- amount: "140" # increased from 125
  denom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833 # vcx vouchers
SellOrderBook:
- amountDenom: mcx
  creator: ""
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 6
  orders:
  - amount: 5 # order hasn't changed
    creator: cosmos1v3p3j7c64c4ls32pcjct333e8vqe45gwwa289q
    id: 2
    price: 15
  # a sell order for 5 mcx has been closed
  priceDenom: vcx
```

Author of buy order recieves 5 mcx vouchers and 50 vcx of their tokens get locked. The 5mcx amount not filled by the sell order creates a buy order on the target chain.

```yml
# Target blockchain
balances:
- amount: "20" # increased from 15
  denom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33 # mcx vouchers
- amount: "825" # decreased from 875
  denom: vcx
BuyOrderBook:
- amountDenom: mcx
  creator: ""
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 7
  orders:
  - amount: 5 # new buy order is created
    creator: cosmos1qlrz3peenc6s3xjv9k97e8ef72nk3qn3a0xax2
    id: 5
    price: 5
  priceDenom: vcx
```

## Cancelling Orders

After the exchanges we have two orders open: sell order on the source chain (5mcx for 15vcx) and a buy order on the target chain (5mcx for 5vcx).

Cancelling a sell order:

```
# Source blockchain
interchanged tx ibcdex cancelSellOrder ibcdex channel-0 mcx vcx 2
```

```yml
# Source blockchain
balances:
- amount: "980" # increased from 975
  denom: mcx
```

Sell order book on the source blokchain is now empty.

Cancelling a buy order:

```
# Target blockchain
interchanged tx ibcdex cancelBuyOrder ibcdex channel-0 mcx vcx 5
```

```yml
# Target blockchain
balances:
- amount: "850" # increased from 825
  denom: vcx
```

Buy order book on the target blokchain is now empty.
