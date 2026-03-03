---
sidebar_position: 48
title: Xexec (xexec)
slug: /packages/xexec
---

# Xexec (xexec)

The `xexec` package contains reusable utilities used by Ignite CLI internals.

For full API details, see the
[`xexec` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/xexec).

## Key APIs

- `func IsCommandAvailable(name string) bool`
- `func IsExec(binaryPath string) (bool, error)`
- `func ResolveAbsPath(filePath string) (path string, err error)`
- `func TryResolveAbsPath(filePath string) string`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/xexec"
```
