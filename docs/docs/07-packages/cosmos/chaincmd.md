---
sidebar_position: 7
title: Chain Command Builder (chaincmd)
slug: /packages/chaincmd
---

# Chain Command Builder (chaincmd)

The `chaincmd` package provides helpers around `SimulationCommand`, `BankSendOption`, and `ChainCmd`.

For full API details, see the
[`chaincmd` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/chaincmd).

## Key APIs

- `func SimulationCommand(appPath string, simName string, options ...SimappOption) step.Option`
- `type BankSendOption func([]string) []string`
- `type ChainCmd struct{ ... }`
- `type GentxOption func([]string) []string`
- `type InPlaceOption func([]string) []string`
- `type KeyringBackend string`
- `type MultiNodeOption func([]string) []string`
- `type Option func(*ChainCmd)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/chaincmd"
```
