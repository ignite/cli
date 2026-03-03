---
sidebar_position: 26
title: Dircache (dircache)
slug: /packages/dircache
---

# Dircache (dircache)

The `dircache` package contains reusable utilities used by Ignite CLI internals.

For full API details, see the
[`dircache` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/dircache).

## Key APIs

- `var ErrCacheNotFound = errors.New("cache not found")`
- `func ClearCache() error`
- `type Cache struct{ ... }`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/dircache"
```
