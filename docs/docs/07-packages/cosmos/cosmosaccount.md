---
sidebar_position: 2
title: Account Registry (cosmosaccount)
slug: /packages/cosmosaccount
---

# Account Registry (cosmosaccount)

The `cosmosaccount` package manages blockchain accounts using Cosmos SDK keyring backends.
It supports creating, importing, exporting, listing, and deleting accounts.

For full API details, see the
[`cosmosaccount` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosaccount).

## Example: Create and query accounts

This example creates an account in a local test keyring and then retrieves it by address.

```go
package main

import (
	"fmt"
	"log"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosaccount"
)

func main() {
	registry, err := cosmosaccount.New(
		cosmosaccount.WithHome("./keyring-test"),
		cosmosaccount.WithKeyringBackend(cosmosaccount.KeyringTest),
		cosmosaccount.WithBech32Prefix(cosmosaccount.AccountPrefixCosmos),
	)
	if err != nil {
		log.Fatal(err)
	}

	account, mnemonic, err := registry.Create("alice")
	if err != nil {
		log.Fatal(err)
	}

	// Store this mnemonic securely. Anyone with it can control the account.
	fmt.Printf("alice mnemonic: %s\n", mnemonic)

	address, err := account.Address(cosmosaccount.AccountPrefixCosmos)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("alice address: %s\n", address)

	loaded, err := registry.GetByAddress(address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("loaded account name: %s\n", loaded.Name)

	accounts, err := registry.List()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("accounts in keyring: %d\n", len(accounts))
}
```
