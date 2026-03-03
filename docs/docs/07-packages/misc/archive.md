---
sidebar_position: 16
title: Archive (archive)
slug: /packages/archive
---

# Archive (archive)

The `archive` package provides helpers around `CreateArchive` and `ExtractArchive`.

For full API details, see the
[`archive` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/archive).

## Key APIs

- `func CreateArchive(dir string, buf io.Writer) error`
- `func ExtractArchive(outDir string, gzipStream io.Reader) error`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/archive"
```
