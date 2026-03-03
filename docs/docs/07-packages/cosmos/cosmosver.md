---
sidebar_position: 11
title: Cosmos SDK Versions (cosmosver)
slug: /packages/cosmosver
---

# Cosmos SDK Versions (cosmosver)

The `cosmosver` package provides helpers around `StargateFortyVersion`, `Versions`, and `CosmosSDKRepoName`.

For full API details, see the
[`cosmosver` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosver).

## Key APIs

- `var StargateFortyVersion = newVersion("0.40.0") ...`
- `var Versions = []Version{ ... } ...`
- `var CosmosSDKRepoName = "cosmos-sdk" ...`
- `type Version struct{ ... }`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cosmosver"
```
