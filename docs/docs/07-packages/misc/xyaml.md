---
sidebar_position: 60
title: Xyaml (xyaml)
slug: /packages/xyaml
---

# Xyaml (xyaml)

The `xyaml` package contains reusable utilities used by Ignite CLI internals.

For full API details, see the
[`xyaml` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/xyaml).

## Key APIs

- `func Marshal(ctx context.Context, obj interface{}, paths ...string) (string, error)`
- `type Map map[string]interface{}`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/xyaml"
```
