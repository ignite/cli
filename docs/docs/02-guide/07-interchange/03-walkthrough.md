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
* Submit sell orders on the Mars chain
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

Next, let's set up an IBC relayer between two chains. If you have used a relayer
in the past, reset the relayer configuration directory:

```
rm -rf ~/.ignite/relayer
```

Now you can use the `ignite relayer configure` command. This command allows you
to specify the source and target chains, along with their respective RPC
endpoints, faucet URLs, port numbers, versions, gas prices, and gas limits.

```
ignite relayer configure -a --source-rpc "http://0.0.0.0:26657" --source-faucet "http://0.0.0.0:4500" --source-port "dex" --source-version "dex-1" --source-gasprice "0.0000025stake" --source-prefix "cosmos" --source-gaslimit 300000 --target-rpc "http://0.0.0.0:26659" --target-faucet "http://0.0.0.0:4501" --target-port "dex" --target-version "dex-1" --target-gasprice "0.0000025stake" --target-prefix "cosmos" --target-gaslimit 300000
```

To create a connection between the two chains, you can use the ignite relayer
connect command. This command will establish a connection between the source and
target chains, allowing you to transfer data and assets between them.

```
ignite relayer connect
```

Now that we have two separate blockchain networks up and running, and a relayer
connection established to facilitate communication between them, we are ready to
begin using the interchain exchange binary to interact with these networks. This
will allow us to create order books and buy/sell orders, enabling us to trade
assets between the two chains.

## Order Book

To create an order book for a pair of tokens, you can use the following command:

```
interchanged tx dex send-create-pair dex channel-0 marscoin venuscoin --from alice --chain-id mars --home ~/.mars
```

This command will create an order book for the pair of tokens `marscoin` and
`venuscoin`. The command will be executed by the user `alice` on the Mars
blockchain. The `--home` parameter specifies the location of the configuration
directory for the Mars blockchain.

Creating an order book affects state on the Mars blockchain to which the
transaction was broadcast and the Venus blockchain.

On the Mars blockchain, the `send-create-pair` command creates an empty sell
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

On the Venus blockchain, the same `send-createPair` command creates a buy order
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

In the `create-pair` command on the Mars blockchain, an IBC packet is sent to
the Venus chain. This packet contains information that is used to create a buy
order book on the Venus chain.

When the Venus chain receives the IBC packet, it processes the information
contained in the packet and creates a buy order book. The Venus chain then sends
an acknowledgement back to the Mars chain to confirm that the buy order book has
been successfully created.

Upon receiving the acknowledgement from the Venus chain, the Mars chain creates
a sell order book. This sell order book is associated with the buy order book on
the Venus chain, allowing users to trade assets between the two chains.

## Sell Order

After creating an order book, the next step is to create a sell order. This can
be done using the `send-sell-order` command, which is used to broadcast a
transaction with a message that locks a specified amount of tokens and creates a
sell order on the Mars blockchain.

```
interchanged tx dex send-sell-order dex channel-0 marscoin 10 venuscoin 15  --from alice --chain-id mars --home ~/.mars
```

In the example provided, the `send-sell-order` command is used to create a sell
order for 10 `marscoin` token and 15 `venuscoin` token. This sell order will be
added to the order book on the Mars blockchain.

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

After creating a sell order, the next step in the trading process is typically
to create a buy order. This can be done using the `send-buy-order` command,
which is used to lock a specified amount of tokens and create a buy order on the
Venus blockchain

```
interchanged tx dex send-buy-order dex channel-0 marscoin 10 venuscoin 5 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

In the example provided, the `send-buy-order` command is used to create a buy
order for 10 `marscoin` token and 5 `venuscoin` token. This buy order will be
added to the order book on the Venus blockchain.

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

You currently have two open orders for `marscoin`:

* A sell order on the Mars chain, where you are offering to sell 10 `marscoin`
  for 15 `venuscoin`.
* A buy order on the Venus chain, where you are willing to buy 5 `marscoin` for
  5 `venuscoin`.

To perform an exchange, you can send a sell order to the Mars chain using the
following command:

```
interchanged tx dex send-sell-order dex channel-0 marscoin 5 venuscoin 3 --from alice --home ~/.mars
```

This sell order, offering to sell 5 `marscoin` for 3 `venuscoin`, will be filled
on the Venus chain by the existing buy order. This will result in the amount of
the buy order on the Venus chain being reduced by 5 `marscoin`.

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

The sender of the filled sell order traded 5 `marscoin` for 25 `venuscoin`
tokens. This means that the amount of the sell order (5 `marscoin`) was
multiplied by the price of the buy order (5 `venuscoin`) to determine the value
of the exchange. In this case, the value of the exchange was 25 `venuscoin`
vouchers.

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

The counterparty, or the sender of the buy `marscoin` order, will receive 5
`marscoin` as a result of the exchange.

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

The `venuscoin` balance has remained unchanged because the appropriate amount of
`venuscoin` (50) was already locked at the time the buy order was created in the
previous step.


## Perform an Exchange with a Buy Order

To perform an exchange with a buy order, send a transaction to the decentralized
exchange to buy 5 `marscoin` for 15 `venuscoin`. This is done by running the
following command:

```
interchanged tx dex send-buy-order dex channel-0 marscoin 5 venuscoin 15 --from alice --home ~/.venus --node tcp://localhost:26659
```

This buy order will be immediately filled on the Mars chain, and the creator of
the sell order will receive 75 `venuscoin` vouchers as payment.

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

The amount of the sell order will be decreased by the amount of the filled buy
order, so in this case it will be decreased by 5 `marscoin`.

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

To complete the exchange with a partially filled sell order, send a transaction
to the decentralized exchange to sell 10 `marscoin` for 3 `venuscoin`. This is
done by running the following command:

```
interchanged tx dex send-sell-order dex channel-0 marscoin 10 venuscoin 3 --from alice --home ~/.mars
```

In this scenario, the sell amount is 10 `marscoin`, but there is an existing buy
order for only 5 `marscoin`. The buy order will be filled completely and removed
from the order book. The author of the previously created buy order will receive
10 `marscoin` vouchers from the exchange.

To check the balances, she can run the following command:

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

To complete the exchange with a partially filled buy order, send a transaction
to the decentralized exchange to buy 10 `marscoin` for 5 `venuscoin`. This is
done by running the following command:

```
interchanged tx dex send-buy-order dex channel-0 marscoin 10 venuscoin 5 --from alice --home ~/.venus --node tcp://localhost:26659
```

In this scenario, the buy order is only partially filled for 5 `marscoin`. There
is an existing sell order for 5 `marscoin` (with a price of 3 `venuscoin`) on
the Mars chain, which is completely filled and removed from the order book. The
author of the closed sell order will receive 15 `venuscoin` vouchers as payment,
which is the product of 5 `marscoin` and 3 `venuscoin`.

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

In this scenario, the author of the buy order will receive 5 `marscoin` vouchers
as payment, which locks up 50 `venuscoin` of their token. The remaining 5
`marscoin` that is not filled by the sell order will create a new buy order on
the Venus chain. This means that the author of the buy order is still interested
in purchasing 5 `marscoin`, and is willing to pay the specified price for it.
The new buy order will remain on the order book until it is filled by another
sell order, or it is cancelled by the buyer.

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

After the exchanges described, there are still two open orders: a sell order on
the Mars chain (5 `marscoin` for 15 `venuscoin`), and a buy order on the Venus
chain (5 `marscoin` for 5 `venuscoin`).

To cancel an order on a blockchain, you can use the `cancel-sell-order` or
`cancel-buy-order` command, depending on the type of order you want to cancel.
The command takes several arguments, including the `channel-id` of the IBC
connection, the `amount-denom` and `price-denom` of the order, and the
`order-id` of the order you want to cancel.

To cancel a sell order on the Mars chain, you would run the following command:

```
interchanged tx dex cancel-sell-order dex channel-0 marscoin venuscoin 0 --from alice --home ~/.mars
```

This will cancel the sell order and remove it from the order book. The balance
of Alice's `marscoin` will be increased by the amount of the cancelled sell
order.

To check Alice's balances, including her updated `marscoin` balance, run the
following command:

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.mars) 
```

This will return a list of Alice's balances, including her updated `marscoin`
balance.

```yml
balances:
- amount: "140"
  denom: ibc/BB38C24E9877
- amount: "980" # increased from 975
  denom: marscoin
- amount: "1000"
  denom: token
```

After the sell order on the Mars chain has been cancelled, the sell order book
on that blockchain will be empty. This means that there are no longer any active
sell orders on the Mars chain, and anyone interested in purchasing `marscoin`
will need to create a new buy order. The sell order book will remain empty until
a new sell order is created and added to it.

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

To cancel a buy order on the `Venus` chain, you can run the following command:

```
interchanged tx dex cancel-buy-order dex channel-0 marscoin venuscoin 1 --from alice --home ~/.venus --node tcp://localhost:26659
```

This will cancel the buy order and remove it from the order book. The balance of
Alice's `venuscoin` will be increased by the amount of the cancelled buy order.

To check Alice's balances, including her updated `venuscoin` balance, you can
run the following command:

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.venus) --node tcp://localhost:26659
```

The amount of `venuscoin` is increased:

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

This will return a list of Alice's balances, including her updated `venuscoin`
balance.

After canceling a buy order, the buy order book on the Venus blockchain will be
empty. This means that there are no longer any active buy orders on the chain,
and anyone interested in selling `marscoin` will need to create a new sell
order. The buy order book will remain empty until a new buy order is created and
added to it.

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

In this walkthrough, we demonstrated how to set up an interchain exchange for
trading tokens between two different blockchain networks. This involved creating
an exchange order book for a specific token pair and establishing a fixed
exchange rate between the two.

Once the exchange was set up, users could send sell orders on the Mars chain and
buy orders on the Venus chain. This allowed them to offer their tokens for sale
or purchase tokens from the exchange. In addition, users could also cancel their
orders if needed. 