---
sidebar_position: 6
title: Command Runner (cmdrunner)
slug: /packages/cmdrunner
---

# Command Runner (cmdrunner)

The `cmdrunner` package provides helpers around `Env`, `Executor`, and `Option`.

For full API details, see the
[`cmdrunner` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cmdrunner).

## Key APIs

- `func Env(key, val string) string`
- `type Executor interface{ ... }`
- `type Option func(*Runner)`
- `type Runner struct{ ... }`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cmdrunner"
```
