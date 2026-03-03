---
sidebar_position: 58
title: Xtime (xtime)
slug: /packages/xtime
---

# Xtime (xtime)

The `xtime` package provides helpers around `FormatUnix`, `FormatUnixInt`, and `NowAfter`.

For full API details, see the
[`xtime` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/xtime).

## Key APIs

- `func FormatUnix(date time.Time) string`
- `func FormatUnixInt(unix int64) string`
- `func NowAfter(unix time.Duration) string`
- `func Seconds(seconds int64) time.Duration`
- `type Clock interface{ ... }`
- `type ClockMock struct{ ... }`
- `type ClockSystem struct{}`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/xtime"
```
