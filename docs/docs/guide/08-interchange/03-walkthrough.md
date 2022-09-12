---
sidebar_position: 3
description: Walkthrough of commands to use the interchain exchange module.
---

# Use the Interchain Exchange

In this chapter, you learn details about the order book and commands to:

- Create an exchange order book for a token pair between two chains
- Send sell orders on source chain
- Send buy orders on target chain
- Cancel sell or buy orders

The next chapter contains the code for the implementation. 

## Order Book

To use the exchange, start by creating an order book for a pair of tokens:

```bash
# Create pair broadcasted to the source blockchain
# interchanged tx dex send-create-pair [src-port] [src-channel] [sourceDenom] [targetDenom]
interchanged tx dex send-create-pair dex channel-0 marscoin venuscoin
```

Define a pair of token with two denominations: 

- Source denom (in this example, `marscoin`)
- Target denom (`venuscoin`)

Creating an order book affects state on the source blockchain to which the transaction was broadcast and the target blockchain. 

On the source blockchain, the `send-create-pair` command creates an empty sell order book:

```yaml
# Created a sell order book on the source blockchain
SellOrderBook:
- amountDenom: marscoin
  creator: ""
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 0
  orders: []
  priceDenom: venuscoin
```

On the target blockchain, the same `send-createPair` command creates a buy order book:

```yaml
# Created a buy order book on the target blockchain
BuyOrderBook:
- amountDenom: marscoin
  creator: ""
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 1
  orders: []
  priceDenom: venuscoin
```

To make an exchange possible, the `createPair` transaction sends an IBC packet to the target chain. 

- When the target chain receives a packet, the target chain creates a buy order book and sends an acknowledgement back to the source chain. 
- When the source chain receives an acknowledgement, the source chain creates a sell order book. 

Sending an IBC packet requires a user to specify a port and a channel through which a packet is transferred.

## Sell Order

After an order book is created, the next step is to create a sell order:

```bash
# Sell order broadcasted to the source blockchain
# interchanged tx dex send-sell-order [src-port] [src-channel] [amountDenom] [amount] [priceDenom] [price]
interchanged tx dex send-sell-order dex channel-0 marscoin 10 venuscoin 15
```

The `send-sellOrder` command broadcasts a message that locks token on the source blockchain and creates a sell order on the source blockchain:

```yaml
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

A buy order has the same arguments, the amount of token to be purchased and a price:

```bash
# Buy order broadcasted to the target blockchain
# interchanged tx dex send-buy-order [src-port] [src-channel] [amountDenom] [amount] [priceDenom] [price]`
interchanged tx dex send-buy-order dex channel-0 marscoin 10 venuscoin 5
```

The `send-buy-order` command locks token on the target blockchain:

```yaml
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

## Perform an Exchange with a Sell Order

You now have two orders open for marscoin: 

- A sell order on the source chain (for 10marscoin at 15venuscoin)
- A buy order on the target chain (for 5marscoin at 5venuscoin) 

Now, perform an exchange by sending a sell order to the source chain:

```bash
# Sell order broadcasted to the source chain
interchanged tx dex send-sell-order dex channel-0 marscoin 5 venuscoin 3
```

The sell order (for 5marscoin at 3venuscoin) is filled on the target chain by the buy order. 

The amount of the buy order on the target chain is decreased by 5marscoin:

```yaml
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

The sender of the filled sell order exchanged 5marscoin for 25 venuscoin vouchers. 

25 vouchers is a product of the amount of the sell order (5marscoin) and price of the buy order (5venuscoin):

```yaml
# Source blockchain
balances:
- amount: "25" # increased from 0
  denom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833 # venuscoin voucher
- amount: "985" # decreased from 990
  denom: marscoin
```

The counterparty (the sender of the buy marscoin order) receives 5 marscoin vouchers:

```yaml
# Target blockchain
balances:
- amount: "5" # increased from 0
  denom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33 # marscoin voucher
```

The venuscoin balance hasn't changed because the correct amount of venuscoin (50) was locked at the creation of the buy order during the previous step.

## Perform an Exchange with a Buy Order

Now, send an order to buy 5marscoin for 15venuscoin:

```bash
# Buy order broadcasted to the target chain
interchanged tx dex send-buy-order dex channel-0 marscoin 5 venuscoin 15
```

A buy order is immediately filled on the source chain and the sell order creator receives 75 venuscoin vouchers. 

The sell order amount is decreased by the amount of the filled buy order (by 5marscoin):

```yaml
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

The creator of the buy order receives 5 marscoin vouchers for 75 venuscoin (5marscoin * 15venuscoin):

```yaml
# Target blockchain
balances:
- amount: "10" # increased from 5
  denom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33 # marscoin vouchers
- amount: "875" # decreased from 950
  denom: venuscoin
```

## Complete Exchange with a Partially Filled Sell Order

Send an order to sell 10marscoin for 3venuscoin:

```bash
# Source blockchain
interchanged tx dex send-sell-order dex channel-0 marscoin 10 venuscoin 3
```

The sell amount is 10marscoin, but the opened buy order amount is only 5marscoin. The buy order gets filled completely and removed from the order book. The author of the previously created buy order receives 10 marscoin vouchers from the exchange:

```yaml
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

The author of the sell order successfuly exchanged 5 marscoin and received 25 venuscoin vouchers. The other 5marscoin created a sell order:

```yaml
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

Create an order to buy 10 marscoin for 5 venuscoin:

```bash
# Target blockchain
interchanged tx dex send-buy-order dex channel-0 marscoin 10 venuscoin 5
```

The buy order is partially filled for 5marscoin. An existing sell order for 5 marscoin (with a price of 3 venuscoin) on the source chain is completely filled and is removed from the order book. The author of the closed sell order receives 15 venuscoin vouchers (product of 5marscoin and 3venuscoin):

```yaml
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

The author of the buy order receives 5 marscoin vouchers which locks 50 venuscoin of their token. The 5marscoin amount that is not filled by the sell order creates a buy order on the target chain:

```yaml
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

## Cancel an Order

After these exchanges, you still have two orders open: 

- A sell order on the source chain (5marscoin for 15venuscoin)
- A buy order on the target chain (5marscoin for 5venuscoin)

To cancel a sell order:

```bash
# Source blockchain
interchanged tx dex cancel-sell-order dex channel-0 marscoin venuscoin 2
```

The balance of marscoin is increased:

```yaml
# Source blockchain
balances:
- amount: "980" # increased from 975
  denom: marscoin
```

The sell order book on the source blockchain is now empty.

To cancel a buy order:

```bash
# Target blockchain
interchanged tx dex cancel-buy-order dex channel-0 marscoin venuscoin 5
```

The amount of venuscoin is increased:

```yaml
# Target blockchain
balances:
- amount: "850" # increased from 825
  denom: venuscoin
```

The buy order book on the target blokchain is now empty.

This walkthrough of the interchain exchange showed you how to:

- Create an exchange order book for a token pair between two chains
- Send sell orders on source chain
- Send buy orders on target chain
- Cancel sell or buy orders

