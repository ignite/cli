---
sidebar_position: 25
title: Debugger (debugger)
slug: /packages/debugger
---

# Debugger (debugger)

The `debugger` package contains reusable utilities used by Ignite CLI internals.

For full API details, see the
[`debugger` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/debugger).

## Key APIs

- `const DefaultAddress = "127.0.0.1:30500" ...`
- `func Run(ctx context.Context, binaryPath string, options ...Option) error`
- `func Start(ctx context.Context, binaryPath string, options ...Option) (err error)`
- `type Option func(*debuggerOptions)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/debugger"
```
