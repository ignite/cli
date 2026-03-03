---
sidebar_position: 27
title: Directory Change Detection (dirchange)
slug: /packages/dirchange
---

# Directory Change Detection (dirchange)

The `dirchange` package computes and compares directory checksums so callers can skip expensive work when inputs have not changed.

For full API details, see the
[`dirchange` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/dirchange).

## When to use

- Short-circuit generation/build steps when source directories are unchanged.
- Persist checksum state between command runs.
- Detect content drift in selected file/folder sets.

## Key APIs

- `ChecksumFromPaths(workdir string, paths ...string) ([]byte, error)`
- `HasDirChecksumChanged(checksumCache cache.Cache[[]byte], cacheKey string, workdir string, paths ...string) (bool, error)`
- `SaveDirChecksum(checksumCache cache.Cache[[]byte], cacheKey string, workdir string, paths ...string) error`

## Example

```go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/dirchange"
)

func main() {
	storage, err := cache.NewStorage(filepath.Join(os.TempDir(), "ignite-cache.db"))
	if err != nil {
		log.Fatal(err)
	}
	c := cache.New[[]byte](storage, "dir-checksum")

	changed, err := dirchange.HasDirChecksumChanged(c, "proto", ".", "proto")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("changed:", changed)
}
```
