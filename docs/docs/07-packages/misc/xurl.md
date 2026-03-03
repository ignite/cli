---
sidebar_position: 59
title: Xurl (xurl)
slug: /packages/xurl
---

# Xurl (xurl)

The `xurl` package provides helpers around `Address`, `HTTP`, and `HTTPEnsurePort`.

For full API details, see the
[`xurl` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/xurl).

## Key APIs

- `func Address(address string) string`
- `func HTTP(s string) (string, error)`
- `func HTTPEnsurePort(s string) string`
- `func HTTPS(s string) (string, error)`
- `func IsHTTP(address string) bool`
- `func MightHTTPS(s string) (string, error)`
- `func TCP(s string) (string, error)`
- `func WS(s string) (string, error)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/xurl"
```
