---
sidebar_position: 7
title: Chain Command Builder (chaincmd)
slug: /packages/chaincmd
---

# Chain Command Builder (chaincmd)

The `chaincmd` package builds `step.Option` command definitions for Cosmos SDK daemon binaries (`simd`, `gaiad`, and others). It does not execute commands directly.

For full API details, see the
[`chaincmd` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/chaincmd).

## When to use

- Build consistent daemon command lines from typed options.
- Reuse command composition across services and tests.
- Keep chain binary-specific flags centralized.

## Key APIs

- `New(appCmd string, options ...Option) ChainCmd`
- `WithHome(home string) Option`
- `WithChainID(chainID string) Option`
- `InitCommand(moniker string, options ...string) step.Option`
- `BankSendCommand(fromAddress, toAddress, amount string, options ...BankSendOption) step.Option`

## Example

```go
package main

import (
	"fmt"

	"github.com/ignite/cli/v29/ignite/pkg/chaincmd"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
)

func main() {
	cmd := chaincmd.New(
		"simd",
		chaincmd.WithHome("./.simapp"),
		chaincmd.WithChainID("demo-1"),
	)

	initStep := step.New(cmd.InitCommand("validator"))
	fmt.Println(initStep.Exec.Command)
	fmt.Println(initStep.Exec.Args)
}
```
