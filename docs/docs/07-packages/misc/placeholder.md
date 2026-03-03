---
sidebar_position: 39
title: Placeholder (placeholder)
slug: /packages/placeholder
---

# Placeholder (placeholder)

The `placeholder` package contains reusable utilities used by Ignite CLI internals.

For full API details, see the
[`placeholder` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/placeholder).

## Key APIs

- `type MissingPlaceholdersError struct{ ... }`
- `type Option func(*Tracer)`
- `type Replacer interface{ ... }`
- `type Tracer struct{ ... }`
- `type ValidationMiscError struct{ ... }`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/placeholder"
```
