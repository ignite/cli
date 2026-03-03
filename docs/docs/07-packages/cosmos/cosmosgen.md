---
sidebar_position: 15
title: Code Generation (cosmosgen)
slug: /packages/cosmosgen
---

# Code Generation (cosmosgen)

The `cosmosgen` package orchestrates code generation from protobuf definitions for:
- Go protobuf code,
- TypeScript client code,
- OpenAPI specifications,
- optional frontend composables/templates.

For full API details, see the
[`cosmosgen` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosgen).

## Example: Run generators

```go
package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosgen"
)

func main() {
	ctx := context.Background()

	storage, err := cache.NewStorage(filepath.Join(os.TempDir(), "ignite-cache.db"))
	if err != nil {
		log.Fatal(err)
	}

	err = cosmosgen.Generate(
		ctx,
		storage,
		".",
		"proto",
		"github.com/acme/my-chain",
		"./web",
		cosmosgen.WithGoGeneration(),
		cosmosgen.WithTSClientGeneration(
			cosmosgen.TypescriptModulePath("./ts-client"),
			"./ts-client",
			true,
		),
		cosmosgen.WithOpenAPIGeneration("./api/openapi.yml", nil),
		cosmosgen.UpdateBufModule(),
	)
	if err != nil {
		log.Fatal(err)
	}
}
```
