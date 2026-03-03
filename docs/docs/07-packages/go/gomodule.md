---
sidebar_position: 9
title: Go Module Analysis (gomodule)
slug: /packages/gomodule
---

# Go Module Analysis (gomodule)

The `gomodule` package provides helpers around `ErrGoModNotFound`, `JoinPath`, and `LocatePath`.

For full API details, see the
[`gomodule` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/gomodule).

## Key APIs

- `var ErrGoModNotFound = errors.New("go.mod not found") ...`
- `func JoinPath(path, version string) string`
- `func LocatePath(ctx context.Context, cacheStorage cache.Storage, src string, pkg Version) (path string, err error)`
- `func ParseAt(path string) (*modfile.File, error)`
- `func SplitPath(path string) (string, string)`
- `type Module struct{ ... }`
- `type Version = module.Version`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/gomodule"
```
