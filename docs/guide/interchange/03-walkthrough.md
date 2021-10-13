---
order: 3
---

# Using the exchange

## Order Book

In this chapter you will build learn details about the order book and the containing commands.
The next chapter contains the code for the implementation. First, learn what you are going to implement.
To use the exchange start from creating a order book for a pair of tokens:

```bash
interchanged tx dex send-create-pair [src-port] [src-channel] [sourceDenom] [targetDenom]

# Create pair broadcasted to the source blockchain
interchanged tx dex send-create-pair dex channel-0 marscoin venuscoin
```

A pair of token is defined by two denominations: source denom (in this example, `marscoin`) and target denom (`venuscoin`). Creating an orderbook affects state on the source blockchain (to which the transaction was broadcasted) and the target blockchain. On the source blockchain `send-createPair` creates an empty sell order book and on the target blockchain a buy order book is created.

```yml
# Created a sell order book on the source blockchain
SellOrderBook:
- amountDenom: marscoin
  creator: ""
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 0
  orders: []
  priceDenom: venuscoin
```

```yml
# Created a buy order book on the target blockchain
BuyOrderBook:
- amountDenom: marscoin
  creator: ""
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 1
  orders: []
  priceDenom: venuscoin
```

To make it possible `createPair` first sends an IBC packet to the target chain. Upon receiving a packet the target chain creates a buy order book and sends back an acknowledgement to the source chain. Upon receiving an acknowledgement, the source chain creates a sell order book. Sending an IBC packet requires a user specify a port and a channel through which a packet will be transferred.

## Sell Order

Once an order book is created, the next step is to create a sell order:

```bash
interchanged tx dex send-sell-order [src-port] [src-channel] [amountDenom] [amount] [priceDenom] [price]

# Sell order broadcasted to the source blockchain
interchanged tx dex send-sell-order dex channel-0 marscoin 10 venuscoin 15
```

The `send-sellOrder` command broadcasts a message that locks tokens on the source blockchain and creates a new sell order on the source blockchain.

```yml
# Source blockchain
balances:
- amount: "990" # decreased from 1000
  denom: marscoin
SellOrderBook:
- amountDenom: marscoin
  creator: ""
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 2
  orders: # a new sell order is created
  - amount: 10
    creator: cosmos1v3p3j7c64c4ls32pcjct333e8vqe45gwwa289q
    id: 0
    price: 15
  priceDenom: venuscoin
```

## Buy Order

A buy order has the same set of arguments: amount of tokens to be purchased and a price.

```bash
`interchanged tx dex send-buy-order [src-port] [src-channel] [amountDenom] [amount] [priceDenom] [price]`

# Buy order broadcasted to the target blockchain
interchanged tx dex send-buy-order dex channel-0 marscoin 10 venuscoin 5
```

The `send-buyOrder` command locks tokens on the target blockchain and creates a buy order book on the target blockchain.

```yml
# Target blockchain
balances:
- amount: "950" # decreased from 1000
  denom: venuscoin
BuyOrderBook:
- amountDenom: marscoin
  creator: ""
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 3
  orders: # a new buy order is created
  - amount: 10
    creator: cosmos1qlrz3peenc6s3xjv9k97e8ef72nk3qn3a0xax2
    id: 1
    price: 5
  priceDenom: venuscoin
```

## Performing an Exchange with a Sell Order

We now have two orders open for marscoin: a sell order on the source chain (for 10marscoin at 15venuscoin) and a buy order on the target chain (for 5marscoin at 5venuscoin). Let's perform an exchange by sending a sell order to the source chain.

```bash
# Sell order broadcasted to the source chain
interchanged tx dex send-sell-order dex channel-0 marscoin 5 venuscoin 3
```

The sell order (for 5marscoin at 3venuscoin) was filled on the target chain by the buy order. The amount of the buy order on the target chain has decreased by 5marscoin.

```yml
# Target blockchain
BuyOrderBook:
- amountDenom: marscoin
  creator: ""
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 5
  orders:
  - amount: 5 # decreased from 10
    creator: cosmos1qlrz3peenc6s3xjv9k97e8ef72nk3qn3a0xax2
    id: 3
    price: 5
  priceDenom: venuscoin
```

The sender of the filled sell order exchanged 5marscoin for 25 venuscoin vouchers. 25 vouchers is a product of the amount of the sell order (5marscoin) and price of the buy order (5venuscoin).

```yml
# Source blockchain
balances:
- amount: "25" # increased from 0
  denom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833 # venuscoin voucher
- amount: "985" # decreased from 990
  denom: marscoin
```

The counterparty (sender of the buy marscoin order) received 5 marscoin vouchers. venuscoin balance hasn't changed, because the correct amount of venuscoin (50) were locked at the creation of the buy order during the previous step.

```yml
# Target blockchain
balances:
- amount: "5" # increased from 0
  denom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33 # marscoin voucher
```

## Performing an Exchange with a Buy Order

An order is sent to buy 5marscoin for 15venuscoin.

```bash
# Buy order broadcasted to the target chain
interchanged tx dex send-buy-order dex channel-0 marscoin 5 venuscoin 15
```

A buy order is immediately filled on the source chain and sell order creator recived 75 venuscoin vouchers. The sell order amount is decreased by the amount of the filled buy order (by 5marscoin).

```yml
# Source blockchain
balances:
- amount: "100" # increased from 25
  denom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833 # venuscoin voucher
SellOrderBook:
- amountDenom: marscoin
  creator: ""
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 4
  orders:
  - amount: 5 # decreased from 10
    creator: cosmos1v3p3j7c64c4ls32pcjct333e8vqe45gwwa289q
    id: 2
    price: 15
  priceDenom: venuscoin
```

Creator of the buy order received 5 marscoin vouchers for 75 venuscoin (5marscoin * 15venuscoin).

```yml
# Target blockchain
balances:
- amount: "10" # increased from 5
  denom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33 # marscoin vouchers
- amount: "875" # decreased from 950
  denom: venuscoin
```

## Complete Exchange with a Partially Filled Sell Order

An order is sent to sell 10marscoin for 3venuscoin.

```bash
# Source blockchain
interchanged tx dex send-sell-order dex channel-0 marscoin 10 venuscoin 3
```

The sell amount is 10marscoin, but the opened buy order amount is only 5marscoin. The buy order gets filled completely and removed from the order book. The author of the previously created buy order recived 10 marscoin vouchers from the exchange.

```yml
# Target blockchain
balances:
- amount: "15" # increased from 5
  denom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33 # marscoin voucher
BuyOrderBook:
- amountDenom: marscoin
  creator: ""
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 5
  orders: [] # buy order with amount 5marscoin has been closed
  priceDenom: venuscoin
```

The author of the sell order successfuly exchanged 5 marscoin and recived 25 venuscoin vouchers. The other 5marscoin created a sell order.

```yml
# Source blockchain
balances:
- amount: "125" # increased from 100
  denom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833 # venuscoin vouchers
- amount: "975" # decreased from 985
  denom: marscoin
- amountDenom: marscoin
SellOrderBook:
  creator: ""
  index: dex-channel-0-marscoin-venuscoin
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

An order is created to buy 10 marscoin for 5 venuscoin.

```bash
# Target blockchain
interchanged tx dex send-buy-order dex channel-0 marscoin 10 venuscoin 5
```

The buy order is partially filled for 5marscoin. An existing sell order for 5 marscoin (with a price of 3 venuscoin) on the source chain has been completely filled and removed from the order book. The author of the closed sell order recived 15 venuscoin vouchers (product of 5marscoin and 3venuscoin).

```yml
# Source blockchain
balances:
- amount: "140" # increased from 125
  denom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833 # venuscoin vouchers
SellOrderBook:
- amountDenom: marscoin
  creator: ""
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 6
  orders:
  - amount: 5 # order hasn't changed
    creator: cosmos1v3p3j7c64c4ls32pcjct333e8vqe45gwwa289q
    id: 2
    price: 15
  # a sell order for 5 marscoin has been closed
  priceDenom: venuscoin
```

Author of buy order recieves 5 marscoin vouchers and 50 venuscoin of their tokens get locked. The 5marscoin amount not filled by the sell order creates a buy order on the target chain.

```yml
# Target blockchain
balances:
- amount: "20" # increased from 15
  denom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33 # marscoin vouchers
- amount: "825" # decreased from 875
  denom: venuscoin
BuyOrderBook:
- amountDenom: marscoin
  creator: ""
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 7
  orders:
  - amount: 5 # new buy order is created
    creator: cosmos1qlrz3peenc6s3xjv9k97e8ef72nk3qn3a0xax2
    id: 5
    price: 5
  priceDenom: venuscoin
```

## Cancelling Orders

After the exchanges we have two orders open: sell order on the source chain (5marscoin for 15venuscoin) and a buy order on the target chain (5marscoin for 5venuscoin).

Cancelling a sell order:

```bash
# Source blockchain
interchanged tx dex cancel-sell-order dex channel-0 marscoin venuscoin 2
```

```yml
# Source blockchain
balances:
- amount: "980" # increased from 975
  denom: marscoin
```

The sell order book on the source blokchain is now empty.

Cancelling a buy order:

```bash
# Target blockchain
interchanged tx dex cancel-buy-order dex channel-0 marscoin venuscoin 5
```

```yml
# Target blockchain
balances:
- amount: "850" # increased from 825
  denom: venuscoin
```

The buy order book on the target blokchain is now empty.
