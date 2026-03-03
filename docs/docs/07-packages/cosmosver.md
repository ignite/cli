---
sidebar_position: 11
title: Cosmos SDK Versions (cosmosver)
slug: /packages/cosmosver
---

# Cosmos SDK Versions (cosmosver)

The `cosmosver` package parses, compares, and detects Cosmos SDK versions used by a chain project.

For full API details, see the
[`cosmosver` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosver).

## When to use

- Detect the Cosmos SDK version from a project before scaffolding or migrations.
- Compare versions to enable/disable version-specific features.
- Access Ignite's known SDK version set and latest supported baseline.

## Key APIs

- `Detect(appPath string) (version Version, err error)`
- `Parse(version string) (v Version, err error)`
- `var Versions = []Version{ ... }`
- `var Latest = Versions[len(Versions)-1]`
- `(Version) Is(version Version) bool`
- `(Version) LT(version Version) bool`
- `(Version) LTE(version Version) bool`
- `(Version) GTE(version Version) bool`

## Common Tasks

- Use `Detect` against a chain root to gate generation paths by SDK version.
- Parse user-provided versions with `Parse` before comparisons.
- Branch behavior with `LT`/`GTE` checks against well-known constants.

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cosmosver"
```
