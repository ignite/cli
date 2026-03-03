---
sidebar_position: 5
title: Token Faucet (cosmosfaucet)
slug: /packages/cosmosfaucet
---

# Token Faucet (cosmosfaucet)

The `cosmosfaucet` package provides:
- A faucet service (`http.Handler`) that sends tokens from a faucet account.
- A client to request tokens from faucet endpoints.

For full API details, see the
[`cosmosfaucet` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosfaucet).

## Example: Start a faucet server

```go
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	sdkmath "cosmossdk.io/math"

	"github.com/ignite/cli/v29/ignite/pkg/chaincmd"
	chaincmdrunner "github.com/ignite/cli/v29/ignite/pkg/chaincmd/runner"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosfaucet"
)

func main() {
	ctx := context.Background()

	cmd := chaincmd.New(
		"simd",
		chaincmd.WithHome("./.simapp"),
		chaincmd.WithKeyringBackend(chaincmd.KeyringBackendTest),
	)

	runner, err := chaincmdrunner.New(ctx, cmd)
	if err != nil {
		log.Fatal(err)
	}

	faucet, err := cosmosfaucet.New(
		ctx,
		runner,
		cosmosfaucet.Account("faucet", "", "", "", ""),
		cosmosfaucet.Coin(sdkmath.NewInt(1000000), sdkmath.NewInt(100000000), "stake"),
		cosmosfaucet.RefreshWindow(24*time.Hour),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":4500", faucet))
}
```

## Example: Request tokens from a faucet

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosfaucet"
)

func main() {
	ctx := context.Background()
	client := cosmosfaucet.NewClient("http://localhost:4500")

	resp, err := client.Transfer(
		ctx,
		cosmosfaucet.NewTransferRequest("cosmos1youraddresshere", []string{"1000000stake"}),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx hash: %s\n", resp.Hash)
}
```
