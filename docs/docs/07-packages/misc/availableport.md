---
sidebar_position: 17
title: Availableport (availableport)
slug: /packages/availableport
---

# Availableport (availableport)

The `availableport` package contains reusable utilities used by Ignite CLI internals.

For full API details, see the
[`availableport` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/availableport).

## Key APIs

- `func Find(n uint, options ...Options) (ports []uint, err error)`
- `type Options func(o *availablePortOptions)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/availableport"
```
