---
sidebar_position: 7
title: Chain Command Builder (chaincmd)
slug: /packages/chaincmd
---

# Chain Command Builder (chaincmd)

The `chaincmd` package builds command definitions for Cosmos chain binaries (`simd`, `gaiad`, etc.).
It does not execute commands directly; it builds `step.Option` values that can be executed by a runner.

For full API details, see the
[`chaincmd` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/chaincmd).

## Example: Build an `init` command

```go
package main

import (
	"fmt"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/chaincmd"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
)

func main() {
	cmd := chaincmd.New(
		"simd",
		chaincmd.WithHome("./.simapp"),
		chaincmd.WithChainID("demo-1"),
		chaincmd.WithKeyringBackend(chaincmd.KeyringBackendTest),
	)

	initStep := step.New(cmd.InitCommand("validator"))

	fmt.Println("binary:", initStep.Exec.Command)
	fmt.Println("args:", strings.Join(initStep.Exec.Args, " "))
}
```
