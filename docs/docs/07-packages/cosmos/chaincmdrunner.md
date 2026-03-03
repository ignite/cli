---
sidebar_position: 4
title: Chain Command Runner (chaincmd/runner)
slug: /packages/chaincmdrunner
---

# Chain Command Runner (chaincmd/runner)

The `chaincmdrunner` package wraps chain binary commands into typed, higher-level operations (accounts, genesis setup, tx queries, node control).

For full API details, see the
[`chaincmdrunner` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/chaincmd/runner).

## When to use

- Execute chain lifecycle commands without manually assembling CLI arguments.
- Manage accounts and genesis setup from automation/test flows.
- Query transaction events using typed selectors instead of raw command output parsing.

## Key APIs

- `New(ctx context.Context, chainCmd chaincmd.ChainCmd, options ...Option) (Runner, error)`
- `(Runner) Init(ctx context.Context, moniker string, args ...string) error`
- `(Runner) Start(ctx context.Context, args ...string) error`
- `(Runner) AddAccount(ctx context.Context, name, mnemonic, coinType, accountNumber, addressIndex string) (Account, error)`
- `(Runner) AddGenesisAccount(ctx context.Context, address, coins string) error`
- `(Runner) QueryTxByEvents(ctx context.Context, selectors ...EventSelector) ([]Event, error)`
- `(Runner) WaitTx(ctx context.Context, txHash string, retryDelay time.Duration, maxRetry int) error`

## Common Tasks

- Build a `Runner` from a configured `chaincmd.ChainCmd` and then call `Init`/`Start` for local node workflows.
- Use `AddAccount`, `ListAccounts`, and `ShowAccount` to manage keyring state in scripted flows.
- Query and filter tx events with `NewEventSelector` plus `QueryTxByEvents`.

## Basic import

```go
import chaincmdrunner "github.com/ignite/cli/v29/ignite/pkg/chaincmd/runner"
```
