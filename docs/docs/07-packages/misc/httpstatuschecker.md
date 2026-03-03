---
sidebar_position: 33
title: Httpstatuschecker (httpstatuschecker)
slug: /packages/httpstatuschecker
---

# Httpstatuschecker (httpstatuschecker)

The `httpstatuschecker` package is a tool check health of http pages.

For full API details, see the
[`httpstatuschecker` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/httpstatuschecker).

## Key APIs

- `func Check(ctx context.Context, addr string, options ...Option) (isAvailable bool, err error)`
- `type Option func(*checker)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/httpstatuschecker"
```
