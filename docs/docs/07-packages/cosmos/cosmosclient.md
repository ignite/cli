---
sidebar_position: 1
title: Blockchain Client (cosmosclient)
slug: /packages/cosmosclient
---

# Blockchain Client (cosmosclient)

The `cosmosclient` package is a Go client for Cosmos SDK chains. It provides helpers
to query chain data and to create, sign, and broadcast transactions.

For full API details, see the
[`cosmosclient` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosclient).

## Example: Transfer between two cosmoshub accounts

This example sends `1000uatom` from `alice` to `bob` on Cosmos Hub.

It assumes:
- You already imported both accounts into a local keyring directory.
- The `alice` account has enough `uatom` to pay the amount and fees.

```go
package main

import (
	"context"
	"fmt"
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosclient"
)

const (
	nodeAddress     = "https://rpc.cosmos.directory:443/cosmoshub"
	keyringDir      = "./keyring-test"
	fromAccountName = "alice"
	toAccountName   = "bob"
	amountToSend    = "1000uatom"
)

func main() {
	ctx := context.Background()

	// Create a Cosmos Hub client.
	client, err := cosmosclient.New(
		ctx,
		cosmosclient.WithNodeAddress(nodeAddress),
		cosmosclient.WithBech32Prefix(cosmosaccount.AccountPrefixCosmos),
		cosmosclient.WithKeyringBackend(cosmosaccount.KeyringTest),
		cosmosclient.WithKeyringDir(keyringDir),
		cosmosclient.WithGas(cosmosclient.GasAuto),
		cosmosclient.WithGasPrices("0.025uatom"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fromAccount, err := client.Account(fromAccountName)
	if err != nil {
		log.Fatal(err)
	}

	toAccount, err := client.Account(toAccountName)
	if err != nil {
		log.Fatal(err)
	}

	toAddress, err := toAccount.Address(cosmosaccount.AccountPrefixCosmos)
	if err != nil {
		log.Fatal(err)
	}

	amount, err := sdk.ParseCoinsNormalized(amountToSend)
	if err != nil {
		log.Fatal(err)
	}

	// Build and broadcast a bank send transaction.
	txService, err := client.BankSendTx(ctx, fromAccount, toAddress, amount)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := txService.Broadcast(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("transaction hash: %s\n", resp.TxHash)
}
```

To import accounts into the test keyring directory, you can use:

```bash
ignite account import alice --keyring-dir ./keyring-test --keyring-backend test
ignite account import bob --keyring-dir ./keyring-test --keyring-backend test
```
