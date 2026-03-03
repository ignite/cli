---
sidebar_position: 51
title: Git Helpers (xgit)
slug: /packages/xgit
---

# Git Helpers (xgit)

The `xgit` package wraps common git repository operations used by Ignite scaffolding and upgrade workflows.

For full API details, see the
[`xgit` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/xgit).

## When to use

- Clone template repositories during setup flows.
- Verify local repository state before mutating operations.
- Resolve repository metadata (origin URL, initialized state).

## Key APIs

- `Clone(ctx context.Context, urlRef, dir string) error`
- `IsRepository(path string) (bool, error)`
- `AreChangesCommitted(dir string) (bool, error)`
- `RepositoryURL(path string) (string, error)`
- `InitAndCommit(path string) error`

## Example

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ignite/cli/v29/ignite/pkg/xgit"
)

func main() {
	_ = context.Background()

	isRepo, err := xgit.IsRepository(".")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("is repository:", isRepo)
}
```
