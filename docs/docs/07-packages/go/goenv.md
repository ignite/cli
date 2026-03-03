---
sidebar_position: 32
title: Goenv (goenv)
slug: /packages/goenv
---

# Goenv (goenv)

defines env variables known by Go and some utilities around it.

For full API details, see the
[`goenv` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/goenv).

## Key APIs

- `const GOBIN = "GOBIN" ...`
- `func Bin() string`
- `func ConfigurePath() error`
- `func GoModCache() string`
- `func GoPath() string`
- `func Path() string`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/goenv"
```
