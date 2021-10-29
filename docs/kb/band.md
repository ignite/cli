---
order: 12
description: IBC oracle integration with BandChain
---

# BandChain Oracle Scaffold

BandChain’s Oracle module is a communication module built-in compliance with IBC protocol which can query data points of various types from BandChain. In addition, other chains can ask our Oracle module for real-time information according to their needs.
BandChain has multiples scripts deployed into the network, and we can request any data using the script id.

## IBC Module Packet Scaffold

BandChain oracle queries can be scaffolded only in IBC modules.

To scaffold an oracle:

```bash
starport scaffold band [queryName] --module [moduleName]
```

### Acknowledgement

The BandChain oracle will return the ack messages with the request's id, and we save the last request id for future queries.

## Files and Directories

When you scaffold a BandChain oracle, the following files and directories are created and modified:

- `proto`: oracle request and response data.
- `x/module_name/keeper`: IBC hooks, gRPC message server.
- `x/module_name/types`: message types, IBC events.
- `x/module_name/client/cli`: CLI command to broadcast a transaction containing a message with a packet.
- `x/module_name/oracle.go`: BandChain oracle packet handlers.

## BandChain Oracle Scaffold Example

The following command scaffolds the IBC-enabled oracle. by default, the starport scaffold oracle for [coin rates](https://laozi-testnet4.cosmoscan.io/oracle-script/37#bridge) request and result.

```shell
$ starport scaffold chain github.com/cosmonaut/oracle --no-module && cd oracle 
$ starport scaffold module consuming --ibc
$ starport s band coinRates --module consuming
```

Note: BandChain module uses version "bandchain-1". Make sure to update the `keys.go` file accordingly.

`x/oracle/types/keys.go`

```go
const Version = "bandchain-1"
```

After scaffold and change the data, run the chain:
```shell
$ starport chain server
```

In another tab, configure and run the starport relayer.
```shell
$ starport relayer configure -a \                                              
--source-rpc "http://rpc-laozi-testnet4.bandchain.org:80" \
--source-faucet "https://laozi-testnet4.bandchain.org/faucet" \
--source-port "oracle" \
--source-gasprice "0uband" \
--source-gaslimit 5000000 \
--source-prefix "band" \
--source-version "bandchain-1" \
--target-rpc "http://localhost:26657" \
--target-faucet "http://localhost:4500" \
--target-port "consuming" \
--target-gasprice "0.0stake" \
--target-gaslimit 300000 \
--target-prefix "cosmos"  \
--target-version "bandchain-1"
            
$ starport relayer connect
```

Open one more terminal tab to make a request transaction, passing the script id.
```shell
# Coin Rates (script 37 into the testnet)
$ oracled tx consuming coin-rates-data 37 4 3 --channel channel-0 --symbols "BTC,ETH,XRP,BCH" --multiplier 1000000 --fee-limit 30uband --prepare-gas 600000 --execute-gas 600000 --from alice --chain-id oracle
```

You can check the last request id returned by ack.
```shell
$ oracled query consuming last-coin-rates-id
request_id: "101276"
```

Furthermore, check the data by request id receive the data packet.
```shell
$ oracled query consuming coin-rates-result 101276
```

### Multiple oracles

You can scaffold multiples oracles by module. After scaffold, you must change the `Calldata` and `Result` parameters into the proto file `moduleName.proto` and adapt the request into the  `cli/client/tx_module_name.go` file. Let's create an example to return the [gold price](https://laozi-testnet4.cosmoscan.io/oracle-script/33#bridge):

```shell
$ starport s band goldPrice --module consuming
```

`proto/gold_price.proto`:
```protobuf
syntax = "proto3";
package test.oracle.consuming;

option go_package = "github.com/test/oracle/x/consuming/types";

message GoldPriceCallData {
  uint64 multiplier = 2;
}

message GoldPriceResult {
  uint64 price = 1;
}
```

`x/cli/client/tx_gold_price.go`:
```go
package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/test/x/consuming/types"
)

// CmdRequestGoldPriceData creates and broadcast a GoldPrice request transaction
func CmdRequestGoldPriceData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gold-price-data [oracle-script-id] [requested-validator-count] [sufficient-validator-count]",
		Short: "Make a new GoldPrice query request via an existing BandChain oracle script",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// retrieve the oracle script id.
			uint64OracleScriptID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			oracleScriptID := types.OracleScriptID(uint64OracleScriptID)

			// retrieve the requested validator count.
			askCount, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			// retrieve the sufficient(minimum) validator count.
			minCount, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			channel, err := cmd.Flags().GetString(flagChannel)
			if err != nil {
				return err
			}

			// retrieve the multiplier for the symbols' price.
			multiplier, err := cmd.Flags().GetUint64(flagMultiplier)
			if err != nil {
				return err
			}

			calldata := &types.GoldPriceCallData{
				Multiplier: multiplier,
			}

			// retrieve the amount of coins allowed to be paid for oracle request fee from the pool account.
			coinStr, err := cmd.Flags().GetString(flagFeeLimit)
			if err != nil {
				return err
			}
			feeLimit, err := sdk.ParseCoinsNormalized(coinStr)
			if err != nil {
				return err
			}

			// retrieve the amount of gas allowed for the prepare step of the oracle script.
			prepareGas, err := cmd.Flags().GetUint64(flagPrepareGas)
			if err != nil {
				return err
			}

			// retrieve the amount of gas allowed for the execute step of the oracle script.
			executeGas, err := cmd.Flags().GetUint64(flagExecuteGas)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgGoldPriceData(
				clientCtx.GetFromAddress().String(),
				oracleScriptID,
				channel,
				calldata,
				askCount,
				minCount,
				feeLimit,
				prepareGas,
				executeGas,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagChannel, "", "The channel id")
	cmd.MarkFlagRequired(flagChannel)
	cmd.Flags().Uint64(flagMultiplier, 1000000, "Multiplier used in calling the oracle script")
	cmd.Flags().String(flagFeeLimit, "", "the maximum tokens that will be paid to all data source providers")
	cmd.Flags().Uint64(flagPrepareGas, 200000, "Prepare gas used in fee counting for prepare request")
	cmd.Flags().Uint64(flagExecuteGas, 200000, "Execute gas used in fee counting for execute request")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
```

Make the request transaction.

```shell
# Gold Price (script 33 into the testnet)
$ oracled tx consuming gold-price-data 33 4 3 --channel channel-0 --multiplier 1000000 --fee-limit 30uband --prepare-gas 600000 --execute-gas 600000 --from alice --chain-id oracle
```

Check the last request id returned by ack and the package data.

```shell
$ oracled query consuming last-gold-price-id
request_id: "101290"

$ oracled query consuming gold-price-result 101290
```
