---
sidebar_position: 42
title: Safeconverter (safeconverter)
slug: /packages/safeconverter
---

# Safeconverter (safeconverter)

The `safeconverter` package provides helpers around `ToInt`, `ToInt64`, and `SafeToConvertToInt`.

For full API details, see the
[`safeconverter` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/safeconverter).

## Key APIs

- `func ToInt[T SafeToConvertToInt](x T) int`
- `func ToInt64[T SafeToConvertToInt](x T) int64`
- `type SafeToConvertToInt interface{ ... }`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/safeconverter"
```
