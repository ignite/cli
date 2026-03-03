---
sidebar_position: 54
title: Xnet (xnet)
slug: /packages/xnet
---

# Xnet (xnet)

The `xnet` package provides helpers around `AnyIPv4Address`, `IncreasePort`, and `IncreasePortBy`.

For full API details, see the
[`xnet` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/xnet).

## Key APIs

- `func AnyIPv4Address(port int) string`
- `func IncreasePort(addr string) (string, error)`
- `func IncreasePortBy(addr string, inc uint64) (string, error)`
- `func LocalhostIPv4Address(port int) string`
- `func MustIncreasePortBy(addr string, inc uint64) string`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/xnet"
```
