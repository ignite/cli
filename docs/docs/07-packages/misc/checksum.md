---
sidebar_position: 19
title: Checksum (checksum)
slug: /packages/checksum
---

# Checksum (checksum)

The `checksum` package provides helpers around `Binary`, `Strings`, and `Sum`.

For full API details, see the
[`checksum` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/checksum).

## Key APIs

- `func Binary(binaryName string) (string, error)`
- `func Strings(inputs ...string) string`
- `func Sum(dirPath, outPath string) error`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/checksum"
```
