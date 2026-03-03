---
sidebar_position: 31
title: Goanalysis (goanalysis)
slug: /packages/goanalysis
---

# Goanalysis (goanalysis)

The `goanalysis` package provides static analysis helpers for Go source code. It is used in Ignite to inspect imports, discover binaries, and apply targeted source rewrites.

For full API details, see the
[`goanalysis` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/goanalysis).

## When to use

- Discover main packages in a project before build/install flows.
- Inspect and normalize import usage in parsed AST files.
- Check or patch specific code patterns in generated sources.

## Key APIs

- `DiscoverMain(path string) ([]string, error)`
- `DiscoverOneMain(path string) (string, error)`
- `FindBlankImports(node *ast.File) []string`
- `FormatImports(f *ast.File) map[string]string`
- `ReplaceCode(pkgPath, oldFunctionName, newFunction string) error`

## Example

```go
package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"

	"github.com/ignite/cli/v29/ignite/pkg/goanalysis"
)

func main() {
	mainPkg, err := goanalysis.DiscoverOneMain(".")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("main package:", mainPkg)

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "main.go", nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	blank := goanalysis.FindBlankImports(file)
	fmt.Println("blank imports:", blank)
}
```
