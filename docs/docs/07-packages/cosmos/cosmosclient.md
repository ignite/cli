---
sidebar_position: 1
title: Blockchain Client (cosmosclient)
slug: /packages/cosmosclient
---

# Blockchain Client (cosmosclient)

The `cosmosclient` package provides a standalone client to connect to Cosmos SDK.

For full API details, see the
[`cosmosclient` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosclient).

## Key APIs

- `const GasAuto = "auto" ...`
- `var FaucetTransferEnsureDuration = time.Second * 40 ...`
- `var WithAddressPrefix = WithBech32Prefix`
- `type BroadcastOption func(*broadcastConfig)`
- `type Client struct{ ... }`
- `type ConsensusInfo struct{ ... }`
- `type FaucetClient interface{ ... }`
- `type Gasometer interface{ ... }`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cosmosclient"
```
