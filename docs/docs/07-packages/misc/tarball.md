---
sidebar_position: 44
title: Tarball (tarball)
slug: /packages/tarball
---

# Tarball (tarball)

The `tarball` package provides helpers around `ErrGzipFileNotFound` and `ExtractFile`.

For full API details, see the
[`tarball` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/tarball).

## Key APIs

- `var ErrGzipFileNotFound = errors.New("file not found in the gzip") ...`
- `func ExtractFile(reader io.Reader, out io.Writer, fileName string) (string, error)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/tarball"
```
