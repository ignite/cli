---
sidebar_position: 32
title: Goenv (goenv)
slug: /packages/goenv
---

# Goenv (goenv)

The `goenv` package provides helpers around Go-related environment values (`GOPATH`, `GOBIN`, module cache, and tool PATH configuration).

For full API details, see the
[`goenv` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/goenv).

## When to use

- Resolve Go binary and cache paths for tool installation.
- Build stable PATH values before running Go tools.
- Keep environment handling centralized.

## Key APIs

- `Bin() string`
- `GoPath() string`
- `GoModCache() string`
- `Path() string`
- `ConfigurePath() error`

## Example

```go
package main

import (
	"fmt"

	"github.com/ignite/cli/v29/ignite/pkg/goenv"
)

func main() {
	fmt.Println("GOPATH:", goenv.GoPath())
	fmt.Println("GOBIN:", goenv.Bin())
	fmt.Println("PATH:", goenv.Path())
}
```
