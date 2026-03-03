---
sidebar_position: 35
title: Localfs (localfs)
slug: /packages/localfs
---

# Localfs (localfs)

The `localfs` package provides helpers around `MkdirAllReset`, `Save`, and `SaveBytesTemp`.

For full API details, see the
[`localfs` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/localfs).

## Key APIs

- `func MkdirAllReset(path string, perm fs.FileMode) error`
- `func Save(f fs.FS, path string) error`
- `func SaveBytesTemp(data []byte, prefix string, perm os.FileMode) (path string, cleanup func(), err error)`
- `func SaveTemp(f fs.FS) (path string, cleanup func(), err error)`
- `func Search(path, pattern string) ([]string, error)`
- `func Watch(ctx context.Context, paths []string, options ...WatcherOption) error`
- `type WatcherOption func(*watcher)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/localfs"
```
