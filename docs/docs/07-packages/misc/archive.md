---
sidebar_position: 16
title: Archive (archive)
slug: /packages/archive
---

# Archive (archive)

The `archive` package contains reusable utilities used by Ignite CLI internals.

For full API details, see the
[`archive` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/archive).

## Key APIs

- `func CreateArchive(dir string, buf io.Writer) error`
- `func ExtractArchive(outDir string, gzipStream io.Reader) error`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/archive"
```
