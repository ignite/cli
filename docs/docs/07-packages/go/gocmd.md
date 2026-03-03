---
sidebar_position: 8
title: Go Command Helpers (gocmd)
slug: /packages/gocmd
---

# Go Command Helpers (gocmd)

The `gocmd` package wraps common `go` tool invocations (`mod tidy`, `build`, `test`, `list`, etc.)
and integrates with Ignite's command execution layer.

For full API details, see the
[`gocmd` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/gocmd).

## Example: List packages and run tests

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ignite/cli/v29/ignite/pkg/gocmd"
)

func main() {
	ctx := context.Background()

	pkgs, err := gocmd.List(ctx, ".", []string{"./..."})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("found %d packages\n", len(pkgs))

	if err := gocmd.Test(ctx, ".", []string{"./..."}); err != nil {
		log.Fatal(err)
	}
}
```
