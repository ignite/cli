---
sidebar_position: 18
title: Cache (cache)
slug: /packages/cache
---

# Cache (cache)

The `cache` package provides a typed, namespaced key-value storage layer backed by BoltDB.

For full API details, see the
[`cache` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cache).

## When to use

- Persist small state between command runs.
- Store typed values by namespace and key.
- Version cache entries per release/process.

## Key APIs

- `NewStorage(path string, options ...StorageOption) (Storage, error)`
- `New[T any](storage Storage, namespace string) Cache[T]`
- `(Cache[T]) Put(key string, value T) error`
- `(Cache[T]) Get(key string) (T, error)`
- `Key(keyParts ...string) string`

## Example

```go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
)

type BuildInfo struct {
	Version string
}

func main() {
	storage, err := cache.NewStorage(filepath.Join(os.TempDir(), "ignite-cache.db"))
	if err != nil {
		log.Fatal(err)
	}

	c := cache.New[BuildInfo](storage, "build-info")
	if err := c.Put(cache.Key("latest"), BuildInfo{Version: "v1.0.0"}); err != nil {
		log.Fatal(err)
	}

	info, err := c.Get("latest")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(info.Version)
}
```
