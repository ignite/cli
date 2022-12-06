---
sidebar_position: 10
description: IBC oracle integration with BandChain
---

# BandChain oracle

The BandChain oracle communication module has built-in compliance using IBC protocol that can query data points of various types from BandChain.

Other chains can query this oracle module for real-time information.

BandChain has multiple scripts deployed into the network. You can request any data using the script id.

## High-level overview

Steps to scaffold an IBC BandChain query oracle to request real-time data from BandChain scripts in a specific IBC-enabled Cosmos SDK module.

## IBC module packet scaffold

BandChain oracle queries can be scaffolded only in IBC modules.

The basic syntax to scaffold a band oracle module is:

```bash
ignite scaffold band [queryName] --module [moduleName]
```

Customize your band oracle with flags:

- --module string - name of the new IBC Module to add the packets to
- --path string - path of the app, default is the current directory (`"."`)
- --signer string - signer label, default is `creator`

### Acknowledgement messages

The BandChain oracle returns the ack messages with the request id. The last request id is saved for future queries.

## Files and directories

When you scaffold a BandChain oracle module, the following files and directories are created and modified:

- `proto`: oracle request and response data.
- `x/module_name/keeper`: IBC hooks, gRPC message server.
- `x/module_name/types`: message types, IBC events.
- `x/module_name/client/cli`: CLI command to broadcast a transaction containing a message with a packet.
- `x/module_name/oracle.go`: BandChain oracle packet handlers.

## Scaffold a BandChain oracle chain

First, scaffold a chain but don't scaffold a default module:

```bash
ignite scaffold chain oracle --no-module 
```

Next, change to the new `oracle` directory and scaffold an IBC-enabled module named `consuming`:

```bash
cd oracle 
ignite scaffold module consuming --ibc
```

Finally, scaffold a BandChain query oracle that can request real-time data:

```bash
ignite scaffold band coinRates --module consuming
```

So far, you have scaffolded:

- A new `oracle` chain without a default module
- A new IBC-enabled `consuming` module
- A new `coinRates` BandChain query oracle

Now it's time to change the data.

## Update version

The output of the `ignite scaffold band coinRates --module consuming` command prompts you to update the `keys.go` file.

In the `x/consuming/types/keys.go` file, update the `Version` variable in the `const` block to the required version that the IBC module supports:

```go
const (
	// ...

	// Version defines the current version the IBC module supports
	Version = "bandchain-1"

	// ...
)
```

## Start your chain in development

To run the chain from the `oracle` directory:

```bash
ignite chain serve
```

Keep this terminal window open.

## Configure and connect the Ignite CLI relayer

If you previously used the Ignite CLI relayer, it is a good idea to remove existing relayer and Ignite CLI configurations:

1. Stop your blockchains.
2. Delete previous configuration files:

    ```bash
    rm -rf ~/.ignite/relayer
    ```

3. Restart your blockchains.

In another terminal tab, configure the [Ignite CLI relayer](./08-relayer.md):

```bash
ignite relayer configure -a \
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
```

When prompted, press Enter to accept the default source and target accounts.

The command output confirms the relayer is successfully configured:

```
? Source Account default
? Target Account default

🔐  Account on "source" is default(band1dscvlx0mhpys9fazuk7ej9z4cq7qknzn09pjpq)

 |· received coins from a faucet
 |· (balance: 10000000uband)

🔐  Account on "target" is default(cosmos1dscvlx0mhpys9fazuk7ej9z4cq7qknznk2pseg)

 |· received coins from a faucet
 |· (balance: 100000stake,5token)

⛓  Configured chains: band-laozi-testnet4-oracle
```

Connect the relayer:

```bash
ignite relayer connect
```

You can see the paths of the `oracle` port on the testnet and the `consuming` port on your local oracle module in the relayer connection status that is output to the terminal:

```
------
Paths
------

band-laozi-testnet4-oracle:
    band-laozi-testnet4 > (port: oracle)    (channel: channel-405)
    oracle              > (port: consuming) (channel: channel-0)

------
Listening and relaying packets between chains...
------
```

Leave this terminal tab open so you can monitor the relayer.

## Make a request transaction

In another terminal tab, use the `oracled` binary to make a request transaction. Because BandChain has multiple scripts already deployed into the network, you can request any data using the BandChain script id. In this case, use script 37 for Coin Rates:

```bash
# Coin Rates (script 37 into the testnet)
oracled tx consuming coin-rates-data 37 4 3 --channel channel-0 --symbols "BTC,ETH,XRP,BCH" --multiplier 1000000 --fee-limit 30uband --prepare-gas 600000 --execute-gas 600000 --from alice --chain-id oracle
```

You can check the last request id that was returned by ack:

```bash
oracled query consuming last-coin-rates-id
# output: request_id: "101276"
```

Now you can check the data by request id to receive the data packet:

```bash
oracled query consuming coin-rates-result 101276
```

### Multiple oracles

You can scaffold multiples oracles by module. After scaffold, you must change the `Calldata` and `Result` parameters in the proto file `moduleName.proto` and then adapt the request in the  `cli/client/tx_module_name.go` file.

To create an example for the [gold price](https://laozi-testnet6.cosmoscan.io/oracle-script/33#bridge) bridge:

```bash
ignite scaffold band goldPrice --module consuming
```

In the `proto/consuming/gold_price.proto` file:

```protobuf
syntax = "proto3";
package oracle.consuming;

option go_package = "oracle/x/consuming/types";

message GoldPriceCallData {
  uint64 multiplier = 2;
}

message GoldPriceResult {
  uint64 price = 1;
}
```

In the `x/consuming/cli/client/tx_gold_price.go` file:

```go
package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"oracle/x/consuming/types"
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

Make the request transaction:

```bash
# Gold Price (script 33 into the testnet)
oracled tx consuming gold-price-data 33 4 3 --channel channel-0 --multiplier 1000000 --fee-limit 30uband --prepare-gas 600000 --execute-gas 600000 --from alice --chain-id oracle
```

Check the last request id that was returned by ack:

```bash
oracled query consuming last-gold-price-id
# output: request_id: "101290"
```

Request the package data:

```bash
oracled query consuming gold-price-result 101290
```
