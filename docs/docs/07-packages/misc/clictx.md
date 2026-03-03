---
sidebar_position: 20
title: Clictx (clictx)
slug: /packages/clictx
---

# Clictx (clictx)

The `clictx` package provides helpers around `Do` and `From`.

For full API details, see the
[`clictx` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/clictx).

## Key APIs

- `func Do(ctx context.Context, fn func() error) error`
- `func From(ctx context.Context) context.Context`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/clictx"
```
