---
sidebar_position: 4
title: Chain Command Runner (chaincmd/runner)
slug: /packages/chaincmdrunner
---

# Chain Command Runner (chaincmd/runner)

The `chaincmd/runner` package wraps chain daemon CLI commands with a higher-level Go API.
It is useful when you need to automate workflows like querying node status, managing accounts,
or sending tokens through a chain binary.

For full API details, see the
[`chaincmd/runner` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/chaincmd/runner).

## Example: Query node status and list keyring accounts

This example uses `gaiad` as the chain binary.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ignite/cli/v29/ignite/pkg/chaincmd"
	chaincmdrunner "github.com/ignite/cli/v29/ignite/pkg/chaincmd/runner"
)

func main() {
	ctx := context.Background()

	cmd := chaincmd.New(
		"gaiad",
		chaincmd.WithHome(os.ExpandEnv("$HOME/.gaia")),
		chaincmd.WithNodeAddress("https://rpc.cosmos.directory:443/cosmoshub"),
		chaincmd.WithKeyringBackend(chaincmd.KeyringBackendTest),
	)

	runner, err := chaincmdrunner.New(ctx, cmd)
	if err != nil {
		log.Fatal(err)
	}

	status, err := runner.Status(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("chain id: %s\n", status.ChainID)

	accounts, err := runner.ListAccounts(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, account := range accounts {
		fmt.Printf("%s: %s\n", account.Name, account.Address)
	}
}
```
