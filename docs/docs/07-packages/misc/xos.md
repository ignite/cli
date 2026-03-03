---
sidebar_position: 55
title: Xos (xos)
slug: /packages/xos
---

# Xos (xos)

The `xos` package provides helpers around `JSONFile`, `CopyFile`, and `CopyFolder`.

For full API details, see the
[`xos` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/xos).

## Key APIs

- `const JSONFile = "json" ...`
- `func CopyFile(srcPath, dstPath string) error`
- `func CopyFolder(srcPath, dstPath string) error`
- `func FileExists(filename string) bool`
- `func FindFiles(directory string, options ...FindFileOptions) ([]string, error)`
- `func RemoveAllUnderHome(path string) error`
- `func Rename(oldPath, newPath string) error`
- `func ValidateFolderCopy(srcPath, dstPath string, exclude ...string) ([]string, error)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/xos"
```
