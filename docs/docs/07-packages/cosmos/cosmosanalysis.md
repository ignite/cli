---
sidebar_position: 13
title: Cosmos Source Analysis (cosmosanalysis)
slug: /packages/cosmosanalysis
---

# Cosmos Source Analysis (cosmosanalysis)

The `cosmosanalysis` package provides static analysis helpers for Cosmos SDK chains.
It can validate chain structure, locate app files, and inspect module wiring.

For full API details, see the
[`cosmosanalysis` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis).

## Example: Validate chain and inspect app modules

```go
package main

import (
	"fmt"
	"log"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis"
	appanalysis "github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis/app"
)

func main() {
	chainRoot := "."

	if err := cosmosanalysis.IsChainPath(chainRoot); err != nil {
		log.Fatal(err)
	}

	appFilePath, err := cosmosanalysis.FindAppFilePath(chainRoot)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("app file:", appFilePath)

	modules, err := appanalysis.FindRegisteredModules(chainRoot)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("registered modules: %d\n", len(modules))
}
```
