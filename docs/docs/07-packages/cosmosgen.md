---
sidebar_position: 15
title: Code Generation (cosmosgen)
slug: /packages/cosmosgen
---

# Code Generation (cosmosgen)

The `cosmosgen` package orchestrates multi-target code generation from protobuf sources, including Go code, TS clients, composables, and OpenAPI output.

For full API details, see the
[`cosmosgen` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosgen).

## When to use

- Run full generation pipelines from application services.
- Configure selective outputs (Go only, TS only, OpenAPI only, etc.).
- Check tool availability and maintain buf-related configuration.

## Key APIs

- `Generate(ctx, cacheStorage, appPath, protoDir, goModPath, frontendPath, options...)`
- `WithGoGeneration()`
- `WithTSClientGeneration(out, tsClientRootPath, useCache)`
- `WithOpenAPIGeneration(out, excludeList)`
- `DepTools() []string`

## Example

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
	storage, err := cache.NewStorage(filepath.Join(os.TempDir(), "ignite-cache.db"))
	if err != nil {
		log.Fatal(err)
	}

	err = cosmosgen.Generate(
		context.Background(),
		storage,
		".",
		"proto",
		"github.com/acme/my-chain",
		"./web",
		cosmosgen.WithGoGeneration(),
		cosmosgen.WithOpenAPIGeneration("./api/openapi.yml", nil),
	)
	if err != nil {
		log.Fatal(err)
	}
}
```
