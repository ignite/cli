---
sidebar_position: 2
title: Account Registry (cosmosaccount)
slug: /packages/cosmosaccount
---

# Account Registry (cosmosaccount)

The `cosmosaccount` package manages Cosmos keyring accounts (create/import/export/list/delete) with configurable backend and Bech32 settings.

For full API details, see the
[`cosmosaccount` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosaccount).

## When to use

- Manage CLI account keys in Ignite services and commands.
- Switch between `test`, `os`, and `memory` keyring backends.
- Resolve addresses/public keys from named keyring entries.

## Key APIs

- `New(options ...Option) (Registry, error)`
- `NewInMemory(options ...Option) (Registry, error)`
- `WithKeyringBackend(backend KeyringBackend) Option`
- `WithHome(path string) Option`
- `(Registry) Create(name string) (Account, mnemonic string, err error)`
- `(Registry) Import(name, secret, passphrase string) (Account, error)`
- `(Registry) Export(name, passphrase string) (key string, err error)`
- `(Registry) GetByName(name string) (Account, error)`
- `(Registry) List() ([]Account, error)`
- `(Account) Address(accPrefix string) (string, error)`

## Common Tasks

- Instantiate one `Registry` with backend/home options and reuse it for all key operations.
- Call `EnsureDefaultAccount` in setup paths that require a predictable signer account.
- Resolve addresses with `Account.Address(prefix)` when your app uses non-default Bech32 prefixes.

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cosmosaccount"
```
