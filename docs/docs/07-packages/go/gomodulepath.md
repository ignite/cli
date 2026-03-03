---
sidebar_position: 10
title: Go Module Paths (gomodulepath)
slug: /packages/gomodulepath
---

# Go Module Paths (gomodulepath)

The `gomodulepath` package parses and validates Go module paths and helps discover
an app module path from the local filesystem.

For full API details, see the
[`gomodulepath` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/gomodulepath).

## Example: Parse and discover module paths

```go
package main

import (
	"fmt"
	"log"

	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
)

func main() {
	parsed, err := gomodulepath.Parse("github.com/acme/my-chain")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("root=%s package=%s\n", parsed.Root, parsed.Package)

	found, appPath, err := gomodulepath.Find(".")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("module=%s path=%s\n", found.RawPath, appPath)
}
```
