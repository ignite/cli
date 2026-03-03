---
sidebar_position: 22
title: Cliui (cliui)
slug: /packages/cliui
---

# Cliui (cliui)

The `cliui` package contains reusable utilities used by Ignite CLI internals.

For full API details, see the
[`cliui` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cliui).

## Key APIs

- `var ErrAbort = errors.New("aborted or not confirmed")`
- `type Option func(s *Session)`
- `type Session struct{ ... }`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cliui"
```
