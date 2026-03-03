---
sidebar_position: 30
title: Events (events)
slug: /packages/events
---

# Events (events)

provides functionalities for packages to log their states as.

For full API details, see the
[`events` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/events).

## Key APIs

- `const DefaultBufferSize = 50`
- `const GroupError = "error"`
- `type Bus struct{ ... }`
- `type BusOption func(*Bus)`
- `type Event struct{ ... }`
- `type Option func(*Event)`
- `type ProgressIndication uint8`
- `type Provider interface{ ... }`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/events"
```
