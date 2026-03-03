---
sidebar_position: 34
title: Jsonfile (jsonfile)
slug: /packages/jsonfile
---

# Jsonfile (jsonfile)

The `jsonfile` package contains reusable utilities used by Ignite CLI internals.

For full API details, see the
[`jsonfile` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/jsonfile).

## Key APIs

- `var ErrFieldNotFound = errors.New("JSON field not found") ...`
- `type JSONFile struct{ ... }`
- `type ReadWriteSeeker interface{ ... }`
- `type UpdateFileOption func(map[string][]byte)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/jsonfile"
```
