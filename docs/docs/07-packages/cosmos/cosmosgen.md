---
sidebar_position: 15
title: Code Generation (cosmosgen)
slug: /packages/cosmosgen
---

# Code Generation (cosmosgen)

The `cosmosgen` package provides helpers around `ErrBufConfig`, `DepTools`, and `Generate`.

For full API details, see the
[`cosmosgen` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosgen).

## Key APIs

- `var ErrBufConfig = errors.New("invalid Buf config") ...`
- `func DepTools() []string`
- `func Generate(ctx context.Context, cacheStorage cache.Storage, ...) error`
- `func MissingTools(f *modfile.File) (missingTools []string)`
- `func UnusedTools(f *modfile.File) (unusedTools []string)`
- `func Vue(path string) error`
- `type ModulePathFunc func(module.Module) string`
- `type Option func(*generateOptions)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cosmosgen"
```
