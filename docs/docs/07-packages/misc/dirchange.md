---
sidebar_position: 27
title: Dirchange (dirchange)
slug: /packages/dirchange
---

# Dirchange (dirchange)

The `dirchange` package provides helpers around `ErrNoFile`, `ChecksumFromPaths`, and `HasDirChecksumChanged`.

For full API details, see the
[`dirchange` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/dirchange).

## Key APIs

- `var ErrNoFile = errors.New("no file in specified paths")`
- `func ChecksumFromPaths(workdir string, paths ...string) ([]byte, error)`
- `func HasDirChecksumChanged(checksumCache cache.Cache[[]byte], cacheKey string, workdir string, ...) (bool, error)`
- `func SaveDirChecksum(checksumCache cache.Cache[[]byte], cacheKey string, workdir string, ...) error`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/dirchange"
```
