---
sidebar_position: 10
title: Go Module Paths (gomodulepath)
slug: /packages/gomodulepath
---

# Go Module Paths (gomodulepath)

The `gomodulepath` package validates and parses Go module path formats, and can discover a module path from a local project directory.

For full API details, see the
[`gomodulepath` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/gomodulepath).

## When to use

- Validate user-provided module names before scaffolding.
- Extract normalized app paths from full module paths.
- Discover module metadata from a working directory.

## Key APIs

- `Parse(rawpath string) (Path, error)`
- `ParseAt(path string) (Path, error)`
- `Find(path string) (parsed Path, appPath string, err error)`
- `ExtractAppPath(path string) string`

## Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
)

func main() {
	p, err := gomodulepath.Parse("github.com/acme/my-chain")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("root:", p.Root)
	fmt.Println("package:", p.Package)
}
```
