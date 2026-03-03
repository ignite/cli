---
sidebar_position: 8
title: Go Command Helpers (gocmd)
slug: /packages/gocmd
---

# Go Command Helpers (gocmd)

The `gocmd` package wraps common `go` tool invocations (`build`, `test`, `mod tidy`, `list`, and more) with consistent execution hooks.

For full API details, see the
[`gocmd` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/gocmd).

## When to use

- Run Go toolchain commands from Ignite services without manually assembling command lines.
- Apply common options (workdir/stdout/stderr) through shared execution abstractions.
- Build binaries for specific targets and manage module commands.

## Key APIs

- `Fmt(ctx, path, options...)`
- `ModTidy(ctx, path, options...)`
- `Build(ctx, out, path, flags, options...)`
- `Test(ctx, path, flags, options...)`
- `List(ctx, path, flags, options...)`

## Example

```go
package main

import (
	"context"
	"log"

	"github.com/ignite/cli/v29/ignite/pkg/gocmd"
)

func main() {
	ctx := context.Background()

	if err := gocmd.ModTidy(ctx, "."); err != nil {
		log.Fatal(err)
	}

	if err := gocmd.Test(ctx, ".", []string{"./..."}); err != nil {
		log.Fatal(err)
	}
}
```
