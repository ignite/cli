---
sidebar_position: 14
title: Buf Integration (cosmosbuf)
slug: /packages/cosmosbuf
---

# Buf Integration (cosmosbuf)

The `cosmosbuf` package provides helpers around `CMDBuf`, `ErrInvalidCommand`, and `Version`.

For full API details, see the
[`cosmosbuf` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosbuf).

## Key APIs

- `const CMDBuf ...`
- `var ErrInvalidCommand ...`
- `func Version(ctx context.Context) (string, error)`
- `type Buf struct{ ... }`
- `type Command string`
- `type GenOption func(*genOptions)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
```
