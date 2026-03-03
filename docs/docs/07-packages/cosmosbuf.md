---
sidebar_position: 14
title: Buf Integration (cosmosbuf)
slug: /packages/cosmosbuf
---

# Buf Integration (cosmosbuf)

The `cosmosbuf` package wraps Buf workflows (`generate`, `export`, `format`, `migrate`, `dep update`) used by Ignite's protobuf pipelines.

For full API details, see the
[`cosmosbuf` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosbuf).

## When to use

- Trigger Buf code generation from Go services.
- Keep Buf invocation flags and error handling consistent.
- Reuse cache-aware generation behavior.

## Key APIs

- `New(cacheStorage cache.Storage, goModPath string) (Buf, error)`
- `(Buf) Generate(ctx, protoPath, output, template, options...)`
- `(Buf) Format(ctx, path)`
- `(Buf) Export(ctx, protoDir, output)`
- `Version(ctx context.Context) (string, error)`

## Example

```go
package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosbuf"
)

func main() {
	storage, err := cache.NewStorage(filepath.Join(os.TempDir(), "ignite-cache.db"))
	if err != nil {
		log.Fatal(err)
	}

	buf, err := cosmosbuf.New(storage, "github.com/acme/my-chain")
	if err != nil {
		log.Fatal(err)
	}

	if err := buf.Format(context.Background(), "./proto"); err != nil {
		log.Fatal(err)
	}
}
```
