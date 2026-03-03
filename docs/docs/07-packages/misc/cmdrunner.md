---
sidebar_position: 6
title: Command Runner (cmdrunner)
slug: /packages/cmdrunner
---

# Command Runner (cmdrunner)

The `cmdrunner` package executes command steps with configurable stdio, hooks, environment, and optional parallelism.

For full API details, see the
[`cmdrunner` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cmdrunner).

## When to use

- Execute multiple shell-like steps with common defaults.
- Attach pre/post hooks and inline stdin payloads.
- Reuse one runner setup across flows.

## Key APIs

- `New(options ...Option) *Runner`
- `DefaultStdout/DefaultStderr/DefaultStdin/DefaultWorkdir`
- `RunParallel()`
- `(Runner) Run(ctx context.Context, steps ...*step.Step) error`

## Example

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
	r := cmdrunner.New(
		cmdrunner.DefaultStdout(os.Stdout),
		cmdrunner.DefaultStderr(os.Stderr),
	)

	err := r.Run(
		context.Background(),
		step.New(step.Exec("go", "version")),
		step.New(step.Exec("go", "env", "GOMOD")),
	)
	if err != nil {
		log.Fatal(err)
	}
}
```
