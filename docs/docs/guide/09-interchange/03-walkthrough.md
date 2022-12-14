---
sidebar_position: 3
description: Walkthrough of commands to use the interchain exchange module.
---

# Use the Interchain Exchange

In this chapter, you will learn about the exchange and how it will function once
it is implemented. This will give you a better understanding of what you will be
building in the coming chapters.

To achieve this, we will perform the following tasks:

* Start two local blockchains
* Set up an IBC relayer between the two chains
* Create an exchange order book for a token pair on the two chains
* Submit sell orders on the source chain
* Submit buy orders on the Venus chain
* Cancel sell or buy orders

Starting the two local blockchains and setting up the IBC relayer will allow us
to create an exchange order book between the two chains. This order book will
allow us to submit sell and buy orders, as well as cancel any orders that we no
longer want to maintain.

It is important to note that the commands in this chapter will only work
properly if you have completed all the following chapters in this tutorial. By
the end of this chapter, you should have a good understanding of how the
exchange will operate.

## Start blockchain nodes

To start using the interchain exchange, you will need to start two separate
blockchains. This can be done by running the `ignite chain serve` command,
followed by the `-c` flag and the path to the configuration file for each
blockchain. For example, to start the `mars` blockchain, you would run:

```
ignite chain serve -c mars.yml
```

To start the `venus` blockchain, you would run a similar command, but with the
path to the `venus.yml` configuration file:

```
ignite chain serve -c venus.yml
```

Once both blockchains are running, you can proceed with configuring the relayer
to enable interchain exchange between the two chains.

## Relayer

To set up an IBC relayer between two chains, you can use the `ignite relayer
configure` command. This command allows you to specify the source and target
chains, along with their respective RPC endpoints, faucet URLs, port numbers,
versions, gas prices, and gas limits.

```
ignite relayer configure -a --source-rpc "http://0.0.0.0:26657" --source-faucet "http://0.0.0.0:4500" --source-port "dex" --source-version "dex-1" --source-gasprice "0.0000025stake" --source-prefix "cosmos" --source-gaslimit 300000 --target-rpc "http://0.0.0.0:26659" --target-faucet "http://0.0.0.0:4501" --target-port "dex" --target-version "dex-1" --target-gasprice "0.0000025stake" --target-prefix "cosmos" --target-gaslimit 300000
```

To create a connection between the two chains, you can use the ignite relayer
connect command. This command will establish a connection between the source and
target chains, allowing you to transfer data and assets between them

```
ignite relayer connect
```

Now that we have two separate blockchain networks up and running, and a relayer
connection established to facilitate communication between them, we are ready to
begin using the interchain exchange binary to interact with these networks. This
will allow us to create order books and buy/sell orders, enabling us to trade
assets between the two chains.

## Order Book

To use the exchange, start by creating an order book for a pair of tokens:

```
interchanged tx dex send-create-pair dex channel-0 marscoin venuscoin --from alice --chain-id mars --home ~/.mars
```

Define a pair of token with two denominations:

- Source denom (in this example, `marscoin`)
- Target denom (`venuscoin`)

Creating an order book affects state on the source blockchain to which the
transaction was broadcast and the target blockchain.

On the source blockchain, the `send-create-pair` command creates an empty sell
order book.

```
interchanged q dex list-sell-order-book
```

```yml
sellOrderBook:
- amountDenom: marscoin
  book:
    idCount: 0
    orders: []
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

On the target blockchain, the same `send-createPair` command creates a buy order
book:

```
interchanged q dex list-buy-order-book --node tcp://localhost:26659
```

```yml
buyOrderBook:
- amountDenom: marscoin
  book:
    idCount: 0
    orders: []
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

To make an exchange possible, the `create-pair` transaction sends an IBC packet
to the target chain.

- When the target chain receives a packet, the target chain creates a buy order
  book and sends an acknowledgement back to the source chain.
- When the source chain receives an acknowledgement, the source chain creates a
  sell order book.

Sending an IBC packet requires a user to specify a port and a channel through
which a packet is transferred.

## Sell Order

After an order book is created, the next step is to create a sell order:

```
interchanged tx dex send-sell-order dex channel-0 marscoin 10 venuscoin 15  --from alice --chain-id mars --home ~/.mars
```

The `send-sell-order` command broadcasts a message that locks token on the
source blockchain and creates a sell order on the source blockchain.

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.mars)
```

```yml
balances:
- amount: "990"  # decreased from 1000
  denom: marscoin
- amount: "1000"
  denom: token
```

```
interchanged q dex list-sell-order-book
```

```yml
sellOrderBook:
- amountDenom: marscoin
  book:
    idCount: 1
    orders: # a new sell order is created
    - amount: 10
      creator: cosmos14ntyzr6d2dx4ppds9tvenx53fn0xl5jcakrtm4
      id: 0
      price: 15
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

## Buy order 

A buy order has the same arguments, the amount of token to be purchased and a
price:

```
interchanged tx dex send-buy-order dex channel-0 marscoin 10 venuscoin 5 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```
The `send-buy-order` command locks token on the target blockchain.

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.venus) --node tcp://localhost:26659
```

```yml
balances:
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "950" # decreased from 1000
  denom: venuscoin
```

```
interchanged q dex list-buy-order-book --node tcp://localhost:26659
```

```yml
buyOrderBook:
- amountDenom: marscoin
  book:
    idCount: 1
    orders: # a new buy order is created
    - amount: 10
      creator: cosmos1mrrttwtdcp47pl4hq6sar3mwqpmtc7pcl9e6ss
      id: 0
      price: 5
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

## Perform an Exchange with a Sell Order

You now have two orders open for marscoin:

- A sell order on the source chain (for 10marscoin at 15venuscoin)
- A buy order on the target chain (for 5marscoin at 5venuscoin)

Now, perform an exchange by sending a sell order to the source chain:

```
interchanged tx dex send-sell-order dex channel-0 marscoin 5 venuscoin 3 --from alice --home ~/.mars
```

The sell order (for 5marscoin at 3venuscoin) is filled on the target chain by
the buy order.

The amount of the buy order on the target chain is decreased by 5marscoin.

```
interchanged q dex list-buy-order-book --node tcp://localhost:26659
```

```yml
buyOrderBook:
- amountDenom: marscoin
  book:
    idCount: 1
    orders:
    - amount: 5 # decreased from 10
      creator: cosmos1mrrttwtdcp47pl4hq6sar3mwqpmtc7pcl9e6ss
      id: 0
      price: 5
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

The sender of the filled sell order exchanged 5marscoin for 25 venuscoin
vouchers.

25 vouchers is a product of the amount of the sell order (5marscoin) and price
of the buy order (5venuscoin):

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.mars)
```

```yml
balances:
- amount: "25" # increased from 0
  denom: ibc/BB38C24E9877
- amount: "985" # decreased from 990
  denom: marscoin
- amount: "1000"
  denom: token
```

The counterparty (the sender of the buy marscoin order) receives 5 marscoin
vouchers:

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.venus) --node tcp://localhost:26659
```

```yml
balances:
- amount: "5" # increased from 0
  denom: ibc/745B473BFE24 # marscoin voucher
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "950"
  denom: venuscoin
```

The venuscoin balance hasn't changed because the correct amount of venuscoin
(50) was locked at the creation of the buy order during the previous step.


## Perform an Exchange with a Buy Order

Now, send an order to buy 5marscoin for 15venuscoin:

```
interchanged tx dex send-buy-order dex channel-0 marscoin 5 venuscoin 15 --from alice --home ~/.venus --node tcp://localhost:26659
```

A buy order is immediately filled on the source chain and the sell order creator
receives 75 venuscoin vouchers.

The sell order amount is decreased by the amount of the filled buy order (by
5marscoin):

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.mars)
```

```yml
balances:
- amount: "100" # increased from 25
  denom: ibc/BB38C24E9877 # venuscoin voucher
- amount: "985"
  denom: marscoin
- amount: "1000"
  denom: token
```

```
interchanged q dex list-sell-order-book
```

```yml
sellOrderBook:
- amountDenom: marscoin
  book:
    idCount: 1
    orders:
    - amount: 5 # decreased from 10
      creator: cosmos14ntyzr6d2dx4ppds9tvenx53fn0xl5jcakrtm4
      id: 0
      price: 15
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

The creator of the buy order receives 5 marscoin vouchers for 75 venuscoin
(5marscoin * 15venuscoin):

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.venus) --node tcp://localhost:26659
```

```yml
balances:
- amount: "10" # increased from 5
  denom: ibc/745B473BFE24 # marscoin vouchers
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "875" # decreased from 950
  denom: venuscoin
```

## Complete Exchange with a Partially Filled Sell Order

Send an order to sell 10marscoin for 3venuscoin:

```
interchanged tx dex send-sell-order dex channel-0 marscoin 10 venuscoin 3 --from alice --home ~/.mars
```

The sell amount is 10marscoin, but the opened buy order amount is only
5marscoin. The buy order gets filled completely and removed from the order book.
The author of the previously created buy order receives 10 marscoin vouchers
from the exchange:

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.venus) --node tcp://localhost:26659
```

```yml
balances:
- amount: "15" # increased from 5
  denom: ibc/745B473BFE24 # marscoin voucher
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "875"
  denom: venuscoin
```

```
interchanged q dex list-buy-order-book --node tcp://localhost:26659
```

```yml
buyOrderBook:
- amountDenom: marscoin
  book:
    idCount: 1
    orders: [] # buy order with amount 5marscoin has been closed
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.mars)
```

The author of the sell order successfuly exchanged 5 marscoin and received 25
venuscoin vouchers. The other 5marscoin created a sell order:

```yml
balances:
- amount: "125" # increased from 100
  denom: ibc/BB38C24E9877 # venuscoin vouchers
- amount: "975" # decreased from 985
  denom: marscoin
- amount: "1000"
  denom: token
```

```
interchanged q dex list-sell-order-book
```

```yml
sellOrderBook:
- amountDenom: marscoin
  book:
    idCount: 2
    orders:
    - amount: 5 # hasn't changed
      creator: cosmos14ntyzr6d2dx4ppds9tvenx53fn0xl5jcakrtm4
      id: 0
      price: 15
    - amount: 5 # new order is created
      creator: cosmos14ntyzr6d2dx4ppds9tvenx53fn0xl5jcakrtm4
      id: 1
      price: 3
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

## Complete Exchange with a Partially Filled Buy Order

Create an order to buy 10 marscoin for 5 venuscoin:

```
interchanged tx dex send-buy-order dex channel-0 marscoin 10 venuscoin 5 --from alice --home ~/.venus --node tcp://localhost:26659
```

The buy order is partially filled for 5marscoin. An existing sell order for 5
marscoin (with a price of 3 venuscoin) on the source chain is completely filled
and is removed from the order book. The author of the closed sell order receives
15 venuscoin vouchers (product of 5marscoin and 3venuscoin):

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.mars)
```

```yml
balances:
- amount: "140" # increased from 125
  denom: ibc/BB38C24E9877 # venuscoin vouchers
- amount: "975"
  denom: marscoin
- amount: "1000"
  denom: token
```

```
interchanged q dex list-sell-order-book
```

```yml
sellOrderBook:
- amountDenom: marscoin
  book:
    idCount: 2
    orders:
    - amount: 5 # order hasn't changed
      creator: cosmos14ntyzr6d2dx4ppds9tvenx53fn0xl5jcakrtm4
      id: 0
      price: 15
    # a sell order for 5 marscoin has been closed
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

The author of the buy order receives 5 marscoin vouchers which locks 50
venuscoin of their token. The 5marscoin amount that is not filled by the sell
order creates a buy order on the target chain:

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.venus) --node tcp://localhost:26659
```

```yml
balances:
- amount: "20" # increased from 15
  denom: ibc/745B473BFE24 # marscoin vouchers
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "825" # decreased from 875
  denom: venuscoin
```

```
interchanged q dex list-buy-order-book --node tcp://localhost:26659
```

```yml
buyOrderBook:
- amountDenom: marscoin
  book:
    idCount: 2
    orders:
    - amount: 5 # new buy order is created
      creator: cosmos1mrrttwtdcp47pl4hq6sar3mwqpmtc7pcl9e6ss
      id: 1
      price: 5
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

## Cancel an Order

After these exchanges, you still have two orders open:

- A sell order on the source chain (5marscoin for 15venuscoin)
- A buy order on the target chain (5marscoin for 5venuscoin)

To cancel a sell order:

```
interchanged tx dex cancel-sell-order dex channel-0 marscoin venuscoin 0 --from alice --home ~/.mars
```

The balance of marscoin is increased:

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.mars) 
```

```yml
balances:
- amount: "140"
  denom: ibc/BB38C24E9877
- amount: "980" # increased from 975
  denom: marscoin
- amount: "1000"
  denom: token
```

The sell order book on the source blockchain is now empty.

```
interchanged q dex list-sell-order-book
```

```yml
sellOrderBook:
- amountDenom: marscoin
  book:
    idCount: 2
    orders: []
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

To cancel a buy order:

```
interchanged tx dex cancel-buy-order dex channel-0 marscoin venuscoin 1 --from alice --home ~/.venus --node tcp://localhost:26659
```

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.venus) --node tcp://localhost:26659
```

The amount of venuscoin is increased:

```yml
balances:
- amount: "20"
  denom: ibc/745B473BFE24
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "850" # increased from 825
  denom: venuscoin
```

The buy order book on the target blockchain is now empty.

```
interchanged q dex list-buy-order-book --node tcp://localhost:26659
```

```yml
buyOrderBook:
- amountDenom: marscoin
  book:
    idCount: 2
    orders: []
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

This walkthrough of the interchain exchange showed you how to:

- Create an exchange order book for a token pair between two chains
- Send sell orders on source chain
- Send buy orders on target chain
- Cancel sell or buy orders
