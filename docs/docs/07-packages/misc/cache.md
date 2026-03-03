---
sidebar_position: 18
title: Cache (cache)
slug: /packages/cache
---

# Cache (cache)

The `cache` package provides helpers around `ErrorNotFound`, `Key`, and `Cache`.

For full API details, see the
[`cache` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cache).

## Key APIs

- `var ErrorNotFound = errors.New("no value was found with the provided key")`
- `func Key(keyParts ...string) string`
- `type Cache[T any] struct{ ... }`
- `type Storage struct{ ... }`
- `type StorageOption func(*Storage)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cache"
```
