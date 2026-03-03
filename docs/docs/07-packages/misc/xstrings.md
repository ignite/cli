---
sidebar_position: 57
title: Xstrings (xstrings)
slug: /packages/xstrings
---

# Xstrings (xstrings)

The `xstrings` package contains reusable utilities used by Ignite CLI internals.

For full API details, see the
[`xstrings` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/xstrings).

## Key APIs

- `func AllOrSomeFilter(list, filterList []string) []string`
- `func FormatUsername(s string) string`
- `func List(n int, do func(i int) string) []string`
- `func NoDash(s string) string`
- `func NoNumberPrefix(s string) string`
- `func StringBetween(s, start, end string) string`
- `func Title(s string) string`
- `func ToUpperFirst(s string) string`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/xstrings"
```
