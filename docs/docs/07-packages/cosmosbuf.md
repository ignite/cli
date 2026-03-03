---
sidebar_position: 14
title: Buf Integration (cosmosbuf)
slug: /packages/cosmosbuf
---

# Buf Integration (cosmosbuf)

The `cosmosbuf` package wraps Buf CLI workflows used by Ignite, including:
- `buf generate`,
- `buf export`,
- `buf config migrate`,
- `buf dep update`.

For full API details, see the
[`cosmosbuf` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosbuf).

## Example: Run code generation with Buf

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
	ctx := context.Background()

	storage, err := cache.NewStorage(filepath.Join(os.TempDir(), "ignite-cache.db"))
	if err != nil {
		log.Fatal(err)
	}

	buf, err := cosmosbuf.New(storage, "github.com/acme/my-chain")
	if err != nil {
		log.Fatal(err)
	}

	err = buf.Generate(
		ctx,
		"./proto",
		"./tmp/gen",
		"./proto/buf.gen.gogo.yaml",
		cosmosbuf.ExcludeFiles("**/query.proto"),
		cosmosbuf.IncludeImports(),
	)
	if err != nil {
		log.Fatal(err)
	}
}
```
