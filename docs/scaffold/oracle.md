---
order: 7
description: IBC oracle integration with Bandchain
---

# Bandchain Oracle Scaffold

BandChain’s Oracle module is a communication module built-in compliance with IBC protocol which can query data points of various types from BandChain. In addition, other chains can ask our Oracle module for real-time information according to their needs.
Bandchain has multiples scripts deployed into the network, and we can request any data using the script id.

## IBC Module Packet Scaffold

Bandchain oracles can be scaffolded only in IBC modules.

To scaffold an oracle:

```
starport scaffold band [oracleName] --module [module_name]
```

### Acknowledgement

The Bandchain oracle will return the ack messages with the id of the request, and we save the last request id for future queries.

## Files and Directories

When you scaffold a Bandchain oracle, the following files and directories are created and modified:

- `proto`: oracle request and response data.
- `x/module_name/keeper`: IBC hooks, gRPC message server.
- `x/module_name/types`: message types, IBC events.
- `x/module_name/client/cli`: CLI command to broadcast a transaction containing a message with a packet.
- `x/module_name/oracle.go`: Bandchain oracle packet handlers.

## Bandchain Oracle Scaffold Example

The following command scaffolds the IBC-enabled oracle. by default, the starport scaffold oracle for `coinRates` request and result.

```shell
$ starport scaffold chain github.com/test/ibcoracle && cd ibcoracle 
$ starport scaffold module consuming --ibc
$ starport s band coinRates --module consuming
```

Also, you can scaffold multiples oracles by module. After scaffold, you must change the `Calldata` and `Result` objects into the proto file `moduleName.proto` and adapt the request into the  `cli/client/tx_module_name.go` file.

```shell
$ starport s band goldPrice --module consuming
```

After scaffold and change the data, configure and run the starport relayer.
```shell
$ starport relayer configure -a \
--source-rpc "http://rpc-laozi-testnet2.bandchain.org:26657" \
--source-faucet "https://laozi-testnet2.bandchain.org/faucet/request" \
--source-port "oracle" \
--source-gasprice "0uband" \
--source-prefix "band" \
--source-version "bandchain-1" \
--target-rpc "http://localhost:26657" \
--target-faucet "http://localhost:4500" \
--target-port "consuming" \
--target-gasprice "0.0stake" \
--target-prefix "cosmos"  \
--target-version "ibcoracle-1"

$ starport relayer connect
```

And make a request transaction, passing the script id.
```shell
# Coin Rates (script 37 into the testnet)
$ ibcoracled tx consuming coin-rates-data 37 4 3 --channel channel-0 --symbols "BTC,ETH,XRP,BCH" --multiplier 1000000 --fee-limit 30uband --request-key "random_string" --prepare-gas 600000 --execute-gas 600000 --from alice --chain-id ibcoracle
# Gold Price (script 33 into the testnet)
$ ibcoracled tx consuming gold-price-data 33 4 3 --channel channel-0 --multiplier 1000000 --fee-limit 30uband --request-key "random_string" --prepare-gas 600000 --execute-gas 600000 --from alice --chain-id ibcoracle
```

You can check the last request id returned by ack.
````shell
$ ibcoracled query consuming last-coin-rates-price-id
request_id: "101276"

$ ibcoracled query consuming last-gold-price-id
request_id: "101290"
````

Furthermore, check the data by request id receive the data packet.
```shell
$ bandchaind query consuming coin-rates-result 101276
$ bandchaind query consuming gold-price-result 101290
```