---
sidebar_position: 5
title: Token Faucet (cosmosfaucet)
slug: /packages/cosmosfaucet
---

# Token Faucet (cosmosfaucet)

The `cosmosfaucet` package provides a local faucet service and client helpers to fund Cosmos accounts during development and tests.

For full API details, see the
[`cosmosfaucet` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosfaucet).

## When to use

- Automatically fund accounts in local/devnet environments.
- Expose a faucet HTTP endpoint backed by a chain key.
- Request funds from an existing faucet endpoint from automation code.

## Key APIs

- `New(ctx context.Context, ccr chaincmdrunner.Runner, options ...Option) (Faucet, error)`
- `TryRetrieve(ctx context.Context, chainID, rpcAddress, faucetAddress, accountAddress string) (string, error)`
- `OpenAPI(apiAddress string) Option`
- `Coin(amount, maxAmount sdkmath.Int, denom string) Option`
- `FeeAmount(amount sdkmath.Int, denom string) Option`
- `RefreshWindow(refreshWindow time.Duration) Option`
- `NewTransferRequest(accountAddress string, coins []string) TransferRequest`

## Common Tasks

- Construct a `Faucet` with chain runner + options, then expose transfer endpoints for local users.
- Use `TryRetrieve` in tests before broadcasting txs to ensure accounts have spendable balance.
- Tune coin amount, max amount, and refresh window to limit faucet abuse.

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cosmosfaucet"
```
