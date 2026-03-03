---
sidebar_position: 49
title: Xfilepath (xfilepath)
slug: /packages/xfilepath
---

# Xfilepath (xfilepath)

The `xfilepath` package defines functions to define path retrievers that support error.

For full API details, see the
[`xfilepath` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/xfilepath).

## Key APIs

- `func IsDir(path string) bool`
- `func MustAbs(path string) (string, error)`
- `func MustInvoke(p PathRetriever) string`
- `func RelativePath(appPath string) (string, error)`
- `type PathRetriever func() (path string, err error)`
- `type PathsRetriever func() (path []string, err error)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/xfilepath"
```
