---
sidebar_position: 24
title: Ctxticker (ctxticker)
slug: /packages/ctxticker
---

# Ctxticker (ctxticker)

The `ctxticker` package provides helpers around `Do` and `DoNow`.

For full API details, see the
[`ctxticker` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/ctxticker).

## Key APIs

- `func Do(ctx context.Context, d time.Duration, fn func() error) error`
- `func DoNow(ctx context.Context, d time.Duration, fn func() error) error`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/ctxticker"
```
