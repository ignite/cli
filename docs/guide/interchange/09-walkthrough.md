---
order: 10
---

# Setup and Commands

The software is ready. Now you can go through: create an order book, buy a token, and sell a token.

In this chapter you will learn how to

* Setup two blockchains to run in parallel
* Create the relayer configuration
* Run the relayer to connect two blockchains
* Use the command line interface of the software you just built in this tutorial

Get started.

## Start the Blockchains

Change your terminal directory to the `interchain` blockchain that you created in this tutorial.

You have created `mars.yml` and `venus.yml` configuration files. Use these configuration files to start two blockchains on the same machine.

```bash
# start networks
starport chain serve -c mars.yml -r
starport chain serve -c venus.yml -r
```

## Setup the relayer

If you have already started tests with the relayer, make sure to clear your old data on a new start with:

```bash
rm ~/.starport/relayer/config.yml
```

If you have not yet, create a new starport account.

```bash
starport account create alice
```

Configure the relayer. Use `Mars` as source blockchain and `Venus` as target blockchain.
For `Mars` you will use the default ports, while for `Venus` the port details are to be found in the `venus.yml` file.

```bash
# relayer configuration
starport relayer configure -a \
--source-rpc "http://0.0.0.0:26657" \
--source-faucet "http://0.0.0.0:4500" \
--source-port "dex" \
--source-version "dex-1" \
--source-gasprice "0.0000025stake" \
--source-prefix "cosmos" \
--source-account "alice" \
--source-gaslimit 300000 \
--target-rpc "http://0.0.0.0:26659" \
--target-faucet "http://0.0.0.0:4501" \
--target-port "dex" \
--target-version "dex-1" \
--target-gasprice "0.0000025stake" \
--target-prefix "cosmos" \
--target-account "alice" \
--target-gaslimit 300000
```

Now connect the two configured blockchains with the relayer.

```bash
# relayer connection
starport relayer connect
```

## Commands Overview

Commands that you are use in this walkthrough are:

```bash
# To get account balances during the tutorial:
# For mars
interchanged q bank balances [address]

# For venus
interchanged q bank balances [address] --node tcp://localhost:26659

# To show the order book
# For mars
interchanged q dex list-sell-order-book

# For venus
interchanged q dex list-buy-order-book --node tcp://localhost:26659
```

## Create a new Order Book

Create an order book for a new pair of tokens to the exchange. The source blockchain is `Mars`. The target blockchain is `Venus`.

The order book is to sell marscoin and buy venuscoin.

```bash
# create the pair
interchanged tx dex send-create-pair dex channel-0 marscoin venuscoin --from alice --chain-id mars --home ~/.mars
```

Display current order books available.

On Mars

```bash
# show the orderbooks
interchanged q dex list-sell-order-book
```

On Venus

```bash
# show the orderbooks
interchanged q dex list-buy-order-book --node tcp://localhost:26659
```

## Create a sell order

Command to create a `packet` with a sell order:

`interchanged tx dex send-sell-order [src-port] [src-channel] [amount-denom] [amount] [priceDenom] [price]`

This sell order sells 10 `marscoin` token for 15 `venuscoin` token.

```bash
# Create and send the sell order
interchanged tx dex send-sell-order dex channel-0 marscoin 10 venuscoin 15 --from alice --chain-id mars --home ~/.mars
```

## Create a buy order

Command to create a `packet` with a buy order:

`interchanged tx dex send-buy-order [src-port] [src-channel] [amount-denom] [amount] [price-denom] [price]`

This sell order buys 10 `marscoin` token for 5 `venuscoin` token.

```bash
# Create and send the buy order
interchanged tx dex send-buy-order dex channel-0 marscoin 10 venuscoin 5 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

## Cancel Buy or Sell Order

Command to create a `packet` with a cancel buy and sell order:

`interchanged tx dex cancel-sell-order [port] [channel] [amount-denom] [price-denom] [order-id]`

`interchanged tx dex cancel-buy-order [port] [channel] [amount-denom] [price-denom] [order-id]`

```bash
# Sell order
interchanged tx dex cancel-sell-order dex channel-0 marscoin venuscoin 0 --from alice --chain-id mars --home ~/.mars

# Buy order
interchanged tx dex cancel-buy-order dex channel-0 marscoin venuscoin 0 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

## Exchange Tokens

Send a sell order for selling 10 `marscoin` token for 15 `venuscoin` token.

```bash
# Sell order
interchanged tx dex send-sell-order dex channel-0 marscoin 10 venuscoin 15 --from alice --chain-id mars --home ~/.mars
```

Check the sell order book:

```bash
# Sell order book
interchanged q dex list-sell-order-book
```

Result
```bash
sellOrderBook:
- amountDenom: marscoin
  book:
    idCount: 2
    orders:
    - amount: 10
      creator: cosmos1kqsatce0yaa04xzuxk6sje72herxvr6yefmdqp
      id: 1
      price: 15
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

Send a sell order for buying 10 `marscoin` token for 5 `venuscoin` token.

```bash
# Buy order
interchanged tx dex send-buy-order dex channel-0 marscoin 10 venuscoin 5 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

Wait for the relayer to pick up your transaction

```bash
Relay 1 packets from mars => venus
Relay 1 acks from venus => mars
```

Check the buy order book:

```bash
# Buy order book
interchanged q dex list-buy-order-book --node tcp://localhost:26659
```

Result
```bash
BuyOrderBook:
- amountDenom: marscoin
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 3
  orders:
  - amount: 10
    creator: cosmos1uhgh5hk8hdmqg8dl8kwfx2c6jxg6ymq7x8h8fz
    id: 1
    price: 5
  priceDenom: venuscoin
pagination:
  next_key: null
  total: "1"
```

## Complete Exchange with a Sell Order

Now, send a sell order packet that fills a previously created buy order.

```bash
# Perform a sell order that get completely filled
interchanged tx dex send-sell-order dex channel-0 marscoin 5 venuscoin 3 --from alice --chain-id mars --home ~/.mars
```

Result
```bash
SellOrderBook:
- amountDenom: marscoin
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 2
  orders:
  - amount: 10
    creator: cosmos1hfdsvk3rl3a7rfgxvnetnnfvpvjxtrjmwrym3f
    id: 0
    price: 15
  priceDenom: venuscoin
BuyOrderBook:
- amountDenom: marscoin
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 3
  orders:
  - amount: 5
    creator: cosmos1uhgh5hk8hdmqg8dl8kwfx2c6jxg6ymq7x8h8fz
    id: 1
    price: 5
  priceDenom: venuscoin
```

## Compare Balances

After the exchange has been executed, compare the new balances and see the new `IBC` token voucher in your balance.

```bash
# Get balance
interchanged q bank balances cosmos1hfdsvk3rl3a7rfgxvnetnnfvpvjxtrjmwrym3f
```

Result
```bash
balances:
- amount: "25"
  denom: ibc/50D70B7748FB
- amount: "985"
  denom: marscoin
- amount: "0"
  denom: stake
- amount: "1000"
  denom: token
pagination:
  next_key: null
  total: "0"
```

## Complete Exchange with a Buy Order

```bash
# Filled buy order
interchanged tx dex send-buy-order dex channel-0 marscoin 5 venuscoin 15 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

Result
```bash
SellOrderBook:
- amountDenom: marscoin
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 2
  orders:
  - amount: 5
    creator: cosmos1hfdsvk3rl3a7rfgxvnetnnfvpvjxtrjmwrym3f
    id: 0
    price: 15
  priceDenom: venuscoin
BuyOrderBook:
- amountDenom: marscoin
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 3
  orders:
  - amount: 5
    creator: cosmos1uhgh5hk8hdmqg8dl8kwfx2c6jxg6ymq7x8h8fz
    id: 1
    price: 5
  priceDenom: venuscoin
```

## Compare Balances

```bash
interchanged q bank balances cosmos1hfdsvk3rl3a7rfgxvnetnnfvpvjxtrjmwrym3f
```

Result
```bash
balances:
- amount: "100"
  denom: ibc/50D70B7748FB
- amount: "985"
  denom: marscoin
- amount: "0"
  denom: stake
- amount: "1000"
  denom: token
```

```bash
interchanged q bank balances cosmos1uhgh5hk8hdmqg8dl8kwfx2c6jxg6ymq7x8h8fz --node tcp://localhost:26659
```

Result
```
balances:
- amount: "10"
  denom: ibc/99678A10AF68
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "875"
  denom: venuscoin
```

## Complete Exchange with a Partially Filled Sell Order

```bash
# Sell order that gets partially filled
interchanged tx dex send-sell-order dex channel-0 marscoin 10 venuscoin 3 --from alice --chain-id mars --home ~/.mars
```

Result
```bash
SellOrderBook:
- amountDenom: marscoin
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 4
  orders:
  - amount: 5
    creator: cosmos1hfdsvk3rl3a7rfgxvnetnnfvpvjxtrjmwrym3f
    id: 0
    price: 15
  - amount: 5
    creator: cosmos1hfdsvk3rl3a7rfgxvnetnnfvpvjxtrjmwrym3f
    id: 2
    price: 3
  priceDenom: venuscoin
BuyOrderBook:
- amountDenom: marscoin
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 3
  orders: []
  priceDenom: venuscoin
```

## Compare Balances

```bash
interchanged q bank balances cosmos1hfdsvk3rl3a7rfgxvnetnnfvpvjxtrjmwrym3f
```

Result

```bash
balances:
- amount: "125"
  denom: ibc/50D70B7748FB
- amount: "975"
  denom: marscoin
- amount: "0"
  denom: stake
- amount: "1000"
  denom: token
```

## Complete Exchange with a Partially Filled Buy Order

```bash
# Buy order that gets partially filled
interchanged tx dex send-buy-order dex channel-0 marscoin 10 venuscoin 5 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

Result

```bash
SellOrderBook:
- amountDenom: marscoin
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 4
  orders:
  - amount: 5
    creator: cosmos1hfdsvk3rl3a7rfgxvnetnnfvpvjxtrjmwrym3f
    id: 0
    price: 15
  priceDenom: venuscoin
BuyOrderBook:
- amountDenom: marscoin
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 5
  orders:
  - amount: 5
    creator: cosmos1uhgh5hk8hdmqg8dl8kwfx2c6jxg6ymq7x8h8fz
    id: 3
    price: 5
  priceDenom: venuscoin
```

## Compare Balances

```bash
interchanged q bank balances cosmos1uhgh5hk8hdmqg8dl8kwfx2c6jxg6ymq7x8h8fz --node tcp://localhost:26659
```

Result

```bash
balances:
- amount: "20"
  denom: ibc/99678A10AF68
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "825"
  denom: venuscoin
```

## Cancel a Sell Order

```bash
# Cancel a sell order
interchanged tx dex cancel-sell-order dex channel-0 marscoin venuscoin 0 --from alice --chain-id mars --home ~/.mars
```

Result:

```bash
SellOrderBook:
- amountDenom: marscoin
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 4
  orders: []
  priceDenom: venuscoin
```

## Cancel a Buy Order

```bash
# Cancel a buy order
interchanged tx dex cancel-buy-order dex channel-0 marscoin venuscoin 0 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

Result:

```bash
BuyOrderBook:
- amountDenom: marscoin
  index: dex-channel-0-marscoin-venuscoin
  orderIDTrack: 5
  orders: []
  priceDenom: venuscoin
```

## Exchange a Token and Return to Original Token

After the exchange of a token from blockchain `Mars` to `Venus`, you end up with a voucher token on the `Venus` blockchain. The voucher token was minted into existence on the `Venus` blockchain while locked up on the `Mars` blockchain. When the process is reversed, the token vouchers on the `Venus` is burned and the original token on `Mars` can get unlocked.

Practice this exercise with the following commands.

First, create the order pair:

```bash
interchanged tx dex send-create-pair dex channel-0 marscoin venuscoin --from alice --chain-id mars --home ~/.mars
```

Create the sell order:

```bash
interchanged tx dex send-sell-order dex channel-0 marscoin 10 venuscoin 15 --from alice --chain-id mars --home ~/.mars
```

Create the matching buy order:

```bash
interchanged tx dex send-buy-order dex channel-0 marscoin 10 venuscoin 15 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

Get the balances:

```bash
interchanged q bank balances cosmos1d745lvvgnrcuzggeze0du30vsdg3y0fmq7p5ct
```

Result:

```bash
balances:
- amount: "150"
  denom: ibc/50D70B7748FB
- amount: "990"
  denom: marscoin
- amount: "0"
  denom: stake
- amount: "1000"
  denom: token
```

See the balance on `Venus` blockchain:

```bash
interchanged q bank balances cosmos14r4pkeat7v6r5n5msr4a33c8lptqrryqau3zrg  --node tcp://localhost:26659
```

Result:

```bash
balances:
- amount: "10"
  denom: ibc/99678A10AF68
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "850"
  denom: venuscoin
```

Check the denom tracing and the path of the token on `Mars`:

```bash
 # See on each chain the saved trace denom
interchanged q dex list-denom-trace
```

Result:

```bash
DenomTrace:
- channel: channel-0
  index: ibc/99678A10AF68
  origin: marscoin
  port: dex
```

See the path of the token voucher received on `Venus`:

```bash
interchanged q dex list-denom-trace --node tcp://localhost:26659
```

Result:

```bash
DenomTrace:
- channel: channel-0
  index: ibc/50D70B7748FB
  origin: venuscoin
  port: ibcde
```

As explained earlier, the process cannot be reversed on the same order book. 

To reverse the exchange path, create a new order book pair:

```bash
# Create a pair in the opposite way
interchanged tx dex send-create-pair dex channel-0 ibc/50D70B7748FB ibc/99678A10AF68 --from alice --chain-id mars --home ~/.mars
```

The order book is now created on both blockchains. Check for the SellOrderBook on `Mars` and BuyOrderBook on `Venus` respectively.

```bash
# On Mars:
SellOrderBook:
- amountDenom: ibc/50D70B7748FB
  index: dex-channel-0-ibc/50D70B7748FB-ibc/99678A10AF68
  orderIDTrack: 0
  orders: []
  priceDenom: ibc/99678A10AF68
```

```bash
# On Venus
BuyOrderBook:
- amountDenom: ibc/50D70B7748FB
  index: dex-channel-0-ibc/50D70B7748FB-ibc/99678A10AF68
  orderIDTrack: 1
  orders: []
  priceDenom: ibc/99678A10AF68
```

Exchange the token from the voucher back to the original token denomination with the following commands.

Create the sell order on `Mars`:

```bash
# Exchange tokens back
interchanged tx dex send-sell-order dex channel-0 ibc/50D70B7748FB 10 ibc/99678A10AF68 1 --from alice --chain-id mars --home ~/.mars
```

Create the buy order on `Venus`:

```bash
interchanged tx dex send-buy-order dex channel-0 ibc/50D70B7748FB 5 ibc/99678A10AF68 1 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

See balances on `Mars`:

```bash
interchanged q bank balances cosmos1d745lvvgnrcuzggeze0du30vsdg3y0fmq7p5ct
```

Result:

```bash
balances:
- amount: "140"
  denom: ibc/50D70B7748FB
- amount: "995"
  denom: marscoin
- amount: "0"
  denom: stake
- amount: "1000"
  denom: token
```

See balances on `Venus`:

```bash
interchanged q bank balances cosmos14r4pkeat7v6r5n5msr4a33c8lptqrryqau3zrg  --node tcp://localhost:26659
```

Result:

```bash
balances:
- amount: "5"
  denom: ibc/99678A10AF68
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "855"
  denom: venuscoin
```

Congratulations, you successfully built and used the `dex` module.

* You created order books across blockchain token pairs.
* You successfully created buy orders and sell orders.
* You matched orders in full and in part.
* You created an order book for an exchange from Mars to Venus, and another order book from Venus to Mars with the voucher token.
* You successfully made an exchange from one blockchain to the other and returned the original token.

In the next chapter you learn how to add tests to your blockchain module to make sure the logic you are expecting is actually executed.
