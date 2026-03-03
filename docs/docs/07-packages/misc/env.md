---
sidebar_position: 28
title: Env (env)
slug: /packages/env
---

# Env (env)

The `env` package contains reusable utilities used by Ignite CLI internals.

For full API details, see the
[`env` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/env).

## Key APIs

- `const DebugEnvVar = "IGNT_DEBUG" ...`
- `func ConfigDir() xfilepath.PathRetriever`
- `func IsDebug() bool`
- `func SetConfigDir(dir string)`
- `func SetDebug()`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/env"
```
