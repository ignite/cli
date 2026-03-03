---
sidebar_position: 12
title: Proto Analysis (protoanalysis)
slug: /packages/protoanalysis
---

# Proto Analysis (protoanalysis)

The `protoanalysis` package parses `.proto` files and builds high-level metadata for:
- packages,
- messages and fields,
- services and RPC signatures,
- HTTP rules.

For full API details, see the
[`protoanalysis` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/protoanalysis).

## Example: Parse proto packages

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ignite/cli/v29/ignite/pkg/protoanalysis"
)

func main() {
	pkgs, err := protoanalysis.Parse(context.Background(), protoanalysis.NewCache(), "./proto")
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		fmt.Printf("package: %s (%d files)\n", pkg.Name, len(pkg.Files))
	}
}
```
