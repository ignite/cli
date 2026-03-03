---
sidebar_position: 9
title: Go Module Analysis (gomodule)
slug: /packages/gomodule
---

# Go Module Analysis (gomodule)

The `gomodule` package parses `go.mod` files, resolves dependencies (including replacements),
and locates modules on disk.

For full API details, see the
[`gomodule` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/gomodule).

## Example: Resolve direct dependencies

```go
package main

import (
	"fmt"
	"log"

	"github.com/ignite/cli/v29/ignite/pkg/gomodule"
)

func main() {
	modFile, err := gomodule.ParseAt(".")
	if err != nil {
		log.Fatal(err)
	}

	deps, err := gomodule.ResolveDependencies(modFile, false)
	if err != nil {
		log.Fatal(err)
	}

	cosmosSDKDeps := gomodule.FilterVersions(deps, "github.com/cosmos/cosmos-sdk")
	for _, dep := range cosmosSDKDeps {
		fmt.Printf("%s %s\n", dep.Path, dep.Version)
	}
}
```
