---
sidebar_position: 1
title: Blockchain Client (cosmosclient)
slug: /packages/cosmosclient
---

# Blockchain Client (cosmosclient)

The `cosmosclient` package provides a high-level client for querying Cosmos SDK chains and building/signing/broadcasting transactions.

For full API details, see the
[`cosmosclient` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosclient).

## When to use

- Connect Ignite tooling to a running node for status and block queries.
- Build and broadcast SDK messages with shared gas/fees/keyring settings.
- Wait for transaction inclusion and inspect block transactions/events.

## Key APIs

- `New(ctx context.Context, options ...Option) (Client, error)`
- `WithNodeAddress(addr string) Option`
- `WithHome(path string) Option`
- `WithKeyringBackend(backend cosmosaccount.KeyringBackend) Option`
- `WithGas(gas string) Option`
- `WithGasPrices(gasPrices string) Option`
- `(Client) BroadcastTx(ctx, account, msgs...) (Response, error)`
- `(Client) WaitForTx(ctx context.Context, hash string) (*ctypes.ResultTx, error)`
- `(Client) Status(ctx context.Context) (*ctypes.ResultStatus, error)`
- `(Client) LatestBlockHeight(ctx context.Context) (int64, error)`

## Common Tasks

- Initialize one `Client` instance with node and keyring options, then reuse it across operations.
- Call `CreateTxWithOptions` or `BroadcastTx` depending on whether you need fine-grained tx overrides.
- Use `WaitForTx`, `WaitForNextBlock`, or `WaitForBlockHeight` for deterministic flows in tests/automation.

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cosmosclient"
```
