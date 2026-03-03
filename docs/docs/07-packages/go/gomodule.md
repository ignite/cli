---
sidebar_position: 9
title: Go Module Analysis (gomodule)
slug: /packages/gomodule
---

# Go Module Analysis (gomodule)

The `gomodule` package parses `go.mod` files, resolves dependencies (including replacements), and locates module directories.

For full API details, see the
[`gomodule` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/gomodule).

## When to use

- Inspect project dependencies programmatically.
- Resolve replacement rules from `go.mod` before dependency lookup.
- Locate module paths for codegen and analysis workflows.

## Key APIs

- `ParseAt(path string) (*modfile.File, error)`
- `ResolveDependencies(f *modfile.File, includeIndirect bool) ([]Version, error)`
- `FilterVersions(dependencies []Version, paths ...string) []Version`
- `SplitPath(path string) (string, string)`
- `JoinPath(path, version string) string`

## Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/ignite/cli/v29/ignite/pkg/gomodule"
)

func main() {
	mod, err := gomodule.ParseAt(".")
	if err != nil {
		log.Fatal(err)
	}

	deps, err := gomodule.ResolveDependencies(mod, false)
	if err != nil {
		log.Fatal(err)
	}

	for _, dep := range gomodule.FilterVersions(deps, "github.com/cosmos/cosmos-sdk") {
		fmt.Printf("%s %s\n", dep.Path, dep.Version)
	}
}
```
