---
sidebar_position: 3
title: Chain Registry Types (chainregistry)
slug: /packages/chainregistry
---

# Chain Registry Types (chainregistry)

The `chainregistry` package defines strongly-typed Go structs for Cosmos chain-registry data (`chain.json` and `assetlist.json`).

For full API details, see the
[`chainregistry` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/chainregistry).

## When to use

- Parse chain-registry JSON into typed values.
- Build tooling that reads chain metadata (APIs, fees, staking tokens, assets).
- Validate or transform registry documents before writing them back.

## Key APIs

- `type Chain struct{ ... }`
- `type APIs struct{ ... }`
- `type APIProvider struct{ ... }`
- `type AssetList struct{ ... }`
- `type Asset struct{ ... }`
- `type Fees struct{ ... }`
- `type Staking struct{ ... }`
- `type Codebase struct{ ... }`
- `type ChainStatus string`
- `type ChainType string`

## Common Tasks

- Decode `chain.json` data into a `Chain` value and inspect RPC/REST metadata.
- Decode `assetlist.json` into `AssetList` to access denom units and logo URIs.
- Use enum-like types (`ChainStatus`, `NetworkType`, `ChainType`) to keep metadata checks explicit.

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/chainregistry"
```
