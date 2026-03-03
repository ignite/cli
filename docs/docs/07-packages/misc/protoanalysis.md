---
sidebar_position: 12
title: Proto Analysis (protoanalysis)
slug: /packages/protoanalysis
---

# Proto Analysis (protoanalysis)

The `protoanalysis` package provides a toolset for analyzing proto files and packages.

For full API details, see the
[`protoanalysis` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/protoanalysis).

## Key APIs

- `var ErrImportNotFound = errors.New("proto import not found")`
- `func HasMessages(ctx context.Context, path string, names ...string) error`
- `func IsImported(path string, dependencies ...string) error`
- `type Cache struct{ ... }`
- `type File struct{ ... }`
- `type Files []File`
- `type HTTPRule struct{ ... }`
- `type Message struct{ ... }`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/protoanalysis"
```
