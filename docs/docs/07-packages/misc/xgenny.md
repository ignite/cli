---
sidebar_position: 50
title: Xgenny (xgenny)
slug: /packages/xgenny
---

# Xgenny (xgenny)

The `xgenny` package contains reusable utilities used by Ignite CLI internals.

For full API details, see the
[`xgenny` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/xgenny).

## Key APIs

- `func Transformer(ctx *plush.Context) genny.Transformer`
- `type ApplyOption func(r *applyOptions)`
- `type OverwriteCallback func(_, _, duplicated []string) error`
- `type Runner struct{ ... }`
- `type SourceModification struct{ ... }`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/xgenny"
```
