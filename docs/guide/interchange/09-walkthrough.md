---
order: 10
---

# Walkthrough

The software is ready. Now you can go through: create an order book, buy a token, and sell a token.

If you have already started tests with the relayer, make sure to clear your old data on a new start with:

```bash
rm ~/.starport/relayer/config.yml
```

## Start the Blockchains

Start the networks. For input, use the `mars.yml` and `venus.yml` for each blockchain.

```bash
# start networks
starport serve -c mars.yml -r
starport serve -c venus.yml -r
```

## Configure the relayer

Start the relayer, use `Mars` as the source blockchain and `Venus` as the target blockchain.
The ports for `Mars` are default ports. For `Venus`, the port details are in the `venus.yml` file.

```bash
# relayer configuration
starport relayer configure -a \
--source-rpc "http://0.0.0.0:26657" \
--source-faucet "http://0.0.0.0:4500" \
--source-port "ibcdex" \
--source-version "ibcdex-1" \
--source-gasprice "0.0000025stake" \
--source-prefix "cosmos" \
--target-rpc "http://0.0.0.0:26659" \
--target-faucet "http://0.0.0.0:4501" \
--target-port "ibcdex" \
--target-version "ibcdex-1" \
--target-gasprice "0.0000025stake" \
--target-prefix "cosmos"
```
## Connect the relayer

After the relayer is configured with the two blockchains, start running the relayer and connect the two blockchains.

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
interchanged q ibcdex list-sellOrderBook

# For venus
interchanged q ibcdex list-buyOrderBook --node tcp://localhost:26659
```

## Create an Order Book

Create an order book for a new pair of tokens to the exchange. 
The source blockchain is `Mars`. The  target blockchain is `Venus`.

The order book is to sell `mcx` and buy `vcx` token.

```bash
# create the pair
interchanged tx ibcdex send-createPair ibcdex channel-0 mcx vcx --from alice --chain-id mars --home ~/.mars
```

Display current order books available:

```bash
# show the orderbooks
interchanged q ibcdex list-sellOrderBook
interchanged q ibcdex list-buyOrderBook --node tcp://localhost:26659
```

## Create a sell order

Command to create a `packet` with a sell order:

`interchanged tx ibcdex send-sellOrder [src-port] [src-channel] [amountDenom] [amount] [priceDenom] [price]`

This sell order sells 10 `mcx` token for 15 `vcx` token.

```bash
# Create and send the sell order
interchanged tx ibcdex send-sellOrder ibcdex channel-0 mcx 10 vcx 15 --from alice --chain-id mars --home ~/.mars
```

## Create a buy order

Command to create a `packet` with a buy order:

`interchanged tx ibcdex send-buyOrder [src-port] [src-channel] [amountDenom] [amount] [priceDenom] [price]`

This sell order buys 10 `mcx` token for 5 `vcx` token.

```bash
# Create and send the buy order
interchanged tx ibcdex send-buyOrder ibcdex channel-0 mcx 10 vcx 5 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

## Cancel Buy or Sell Order

Command to create a `packet` with a cancel buy and sell order:

`interchanged tx ibcdex cancelSellOrder [port] [channel] [amountDenom] [priceDenom] [orderID]`

`interchanged tx ibcdex cancelBuyOrder [port] [channel] [amountDenom] [priceDenom] [orderID]`

```bash
# Sell order
interchanged tx ibcdex cancelSellOrder ibcdex channel-0 mcx vcx 0 --from alice --chain-id mars --home ~/.mars

# Buy order
interchanged tx ibcdex cancelBuyOrder ibcdex channel-0 mcx vcx 0 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

## Exchange Tokens

Send a sell order for selling 10 `mcx` token for 15 `vcx` token.

```bash
# Sell order
interchanged tx ibcdex send-sellOrder ibcdex channel-0 mcx 10 vcx 15 --from alice --chain-id mars --home ~/.mars
```

Check the sell order book:

```bash
# Sell order book
interchanged q ibcdex list-sellOrderBook
```

Result
```bash
SellOrderBook:
- amountDenom: mcx
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 2
  orders:
  - amount: 10
    creator: cosmos1hfdsvk3rl3a7rfgxvnetnnfvpvjxtrjmwrym3f
    id: 0
    price: 15
  priceDenom: vcx
pagination:
  next_key: null
  total: "1"
```

Send a sell order for buying 10 `mcx` token for 5 `vcx` token.

```bash
# Buy order
interchanged tx ibcdex send-buyOrder ibcdex channel-0 mcx 10 vcx 5 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

Check the buy order book:

```bash
# Buy order book
interchanged q ibcdex list-buyOrderBook --node tcp://localhost:26659
```

Result
```bash
BuyOrderBook:
- amountDenom: mcx
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 3
  orders:
  - amount: 10
    creator: cosmos1uhgh5hk8hdmqg8dl8kwfx2c6jxg6ymq7x8h8fz
    id: 1
    price: 5
  priceDenom: vcx
pagination:
  next_key: null
  total: "1"
```

## Complete Exchange with a Sell Order

Now, send a sell order packet that fills a previously created buy order.

```bash
# Perform a sell order that get completely filled
interchanged tx ibcdex send-sellOrder ibcdex channel-0 mcx 5 vcx 3 --from alice --chain-id mars --home ~/.mars
```

Result
```bash
SellOrderBook:
- amountDenom: mcx
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 2
  orders:
  - amount: 10
    creator: cosmos1hfdsvk3rl3a7rfgxvnetnnfvpvjxtrjmwrym3f
    id: 0
    price: 15
  priceDenom: vcx
BuyOrderBook:
- amountDenom: mcx
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 3
  orders:
  - amount: 5
    creator: cosmos1uhgh5hk8hdmqg8dl8kwfx2c6jxg6ymq7x8h8fz
    id: 1
    price: 5
  priceDenom: vcx
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
  denom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833
- amount: "985"
  denom: mcx
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
interchanged tx ibcdex send-buyOrder ibcdex channel-0 mcx 5 vcx 15 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

Result
```bash
SellOrderBook:
- amountDenom: mcx
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 2
  orders:
  - amount: 5
    creator: cosmos1hfdsvk3rl3a7rfgxvnetnnfvpvjxtrjmwrym3f
    id: 0
    price: 15
  priceDenom: vcx
BuyOrderBook:
- amountDenom: mcx
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 3
  orders:
  - amount: 5
    creator: cosmos1uhgh5hk8hdmqg8dl8kwfx2c6jxg6ymq7x8h8fz
    id: 1
    price: 5
  priceDenom: vcx
```

## Compare Balances

```bash
interchanged q bank balances cosmos1hfdsvk3rl3a7rfgxvnetnnfvpvjxtrjmwrym3f
```

Result
```bash
balances:
- amount: "100"
  denom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833
- amount: "985"
  denom: mcx
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
  denom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "875"
  denom: vcx
```

## Complete Exchange with a Partially Filled Sell Order

```bash
# Sell order that gets partially filled
interchanged tx ibcdex send-sellOrder ibcdex channel-0 mcx 10 vcx 3 --from alice --chain-id mars --home ~/.mars
```

Result
```bash
SellOrderBook:
- amountDenom: mcx
  index: ibcdex-channel-0-mcx-vcx
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
  priceDenom: vcx
BuyOrderBook:
- amountDenom: mcx
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 3
  orders: []
  priceDenom: vcx
```

## Compare Balances

```bash
interchanged q bank balances cosmos1hfdsvk3rl3a7rfgxvnetnnfvpvjxtrjmwrym3f
```

Result
```bash
balances:
- amount: "125"
  denom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833
- amount: "975"
  denom: mcx
- amount: "0"
  denom: stake
- amount: "1000"
  denom: token
```

## Complete Exchange with a Partially Filled Buy Order

```bash
# Buy order that gets partially filled
interchanged tx ibcdex send-buyOrder ibcdex channel-0 mcx 10 vcx 5 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

Result
```
SellOrderBook:
- amountDenom: mcx
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 4
  orders:
  - amount: 5
    creator: cosmos1hfdsvk3rl3a7rfgxvnetnnfvpvjxtrjmwrym3f
    id: 0
    price: 15
  priceDenom: vcx
BuyOrderBook:
- amountDenom: mcx
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 5
  orders:
  - amount: 5
    creator: cosmos1uhgh5hk8hdmqg8dl8kwfx2c6jxg6ymq7x8h8fz
    id: 3
    price: 5
  priceDenom: vcx
```

## Compare Balances

```bash
interchanged q bank balances cosmos1uhgh5hk8hdmqg8dl8kwfx2c6jxg6ymq7x8h8fz --node tcp://localhost:26659
```

Result
```bash
balances:
- amount: "20"
  denom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "825"
  denom: vcx
```

## Cancel a Sell Order

```bash
# Cancel a sell order
interchanged tx ibcdex cancelSellOrder ibcdex channel-0 mcx vcx 0 --from alice --chain-id mars --home ~/.mars
```

Result:

```
SellOrderBook:
- amountDenom: mcx
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 4
  orders: []
  priceDenom: vcx
```

## Cancel a Buy Order

```bash
# Cancel a buy order
interchanged tx ibcdex cancelBuyOrder ibcdex channel-0 mcx vcx 0 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

Result:

```
BuyOrderBook:
- amountDenom: mcx
  index: ibcdex-channel-0-mcx-vcx
  orderIDTrack: 5
  orders: []
  priceDenom: vcx
```

## Exchange a Token and Return to Original Token

After the exchange of a token from blockchain `Mars` to `Venus`, you end up with a voucher token on the `Venus` blockchain. The voucher token was minted into existence on the `Venus` blockchain while locked up on the `Mars` blockchain. When the process is reversed, the token vouchers on the `Venus` is burned and the original token on `Mars` can get unlocked.

Practice this exercise with the following commands.

First, create the order pair:

```bash
interchanged tx ibcdex send-createPair ibcdex channel-0 mcx vcx --from alice --chain-id mars --home ~/.mars
```

Create the sell order:

```bash
interchanged tx ibcdex send-sellOrder ibcdex channel-0 mcx 10 vcx 15 --from alice --chain-id mars --home ~/.mars
```

Create the matching buy order:

```bash
interchanged tx ibcdex send-buyOrder ibcdex channel-0 mcx 10 vcx 15 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

Get the balances:

```bash
interchanged q bank balances cosmos1d745lvvgnrcuzggeze0du30vsdg3y0fmq7p5ct
```

Result:

```bash
balances:
- amount: "150"
  denom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833
- amount: "990"
  denom: mcx
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
  denom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "850"
  denom: vcx
```

Check the denom tracing and the path of the token on `Mars`:

```bash
 # See on each chain the saved trace denom
interchanged q ibcdex list-denomTrace
```

Result:

```bash
DenomTrace:
- channel: channel-0
  index: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33
  origin: mcx
  port: ibcdex
```

See the path of the token voucher received on `Venus`:

```bash
interchanged q ibcdex list-denomTrace --node tcp://localhost:26659
```

Result:

```bash
DenomTrace:
- channel: channel-0
  index: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833
  origin: vcx
  port: ibcde
```

As explained earlier, the process cannot be reversed on the same order book. 

To reverse the exchange path, create a new order book pair:

```bash
# Create a pair in the opposite way
interchanged tx ibcdex send-createPair ibcdex channel-0 ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833 ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33 --from alice --chain-id mars --home ~/.mars
```

The order book is now created on both blockchains. Check for the SellOrderbook on `Mars` and Buyorderbook on `Venus` respectively.

```bash
# On Mars:
SellOrderBook:
- amountDenom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833
  index: ibcdex-channel-0-ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833-ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33
  orderIDTrack: 0
  orders: []
  priceDenom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33
```

```bash
# On Venus
BuyOrderBook:
- amountDenom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833
  index: ibcdex-channel-0-ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833-ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33
  orderIDTrack: 1
  orders: []
  priceDenom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33
```

Exchange the token from the voucher back to the original token denomination with the following commands.

Create the sell order on `Mars`:

```bash
# Exchange tokens back
interchanged tx ibcdex send-sellOrder ibcdex channel-0 ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833 10 ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33 1 --from alice --chain-id mars --home ~/.mars
```

Create the buy order on `Venus`:

```bash
interchanged tx ibcdex send-buyOrder ibcdex channel-0 ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833 5 ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33 1 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

See balances on `Mars`:

```bash
interchanged q bank balances cosmos1d745lvvgnrcuzggeze0du30vsdg3y0fmq7p5ct
```

Result:

```bash
balances:
- amount: "140"
  denom: ibc/50D70B7748FB8AA69F09114EC9E5615C39E07381FE80E628A1AF63A6F5C79833
- amount: "995"
  denom: mcx
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
  denom: ibc/99678A10AF684E33E88959727F2455AE42CCC64CD76ECFA9691E1B5A32342D33
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "855"
  denom: vcx
```

Congratulations, you successfully built and used the `ibcdex` module.
- You created order books across blockchain token pairs. 
- You successfully created buy orders and sell orders.
- You matched orders in full and in part.
- You created an order book for an exchange from Mars to Venus, and another order book from Venus to Mars with the voucher token.
- You successfully made an exchange from one blockchain to the other and returned the original token.

In the next chapter you learn how to add tests to your blockchain module to make sure the logic you are expecting is actually executed. 
