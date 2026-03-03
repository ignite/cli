---
sidebar_position: 6
title: Command Runner (cmdrunner)
slug: /packages/cmdrunner
---

# Command Runner (cmdrunner)

The `cmdrunner` package is a lightweight command execution layer with support for:
- per-command hooks,
- default stdio/workdir configuration,
- sequential or parallel step execution.

For full API details, see the
[`cmdrunner` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cmdrunner).

## Example: Run command steps

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
)

func main() {
	ctx := context.Background()

	runner := cmdrunner.New(
		cmdrunner.DefaultStdout(os.Stdout),
		cmdrunner.DefaultStderr(os.Stderr),
		cmdrunner.DefaultWorkdir("."),
	)

	err := runner.Run(
		ctx,
		step.New(step.Exec("go", "version")),
		step.New(step.Exec("go", "env", "GOMOD")),
	)
	if err != nil {
		log.Fatal(err)
	}
}
```
