---
sidebar_position: 31
title: Goanalysis (goanalysis)
slug: /packages/goanalysis
---

# Goanalysis (goanalysis)

The `goanalysis` package provides a toolset for statically analysing Go applications.

For full API details, see the
[`goanalysis` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/goanalysis).

## Key APIs

- `var ErrMultipleMainPackagesFound = errors.New("multiple main packages found")`
- `func AddOrRemoveTools(f *modfile.File, writer io.Writer, importsToAdd, importsToRemove []string) error`
- `func DiscoverMain(path string) (pkgPaths []string, err error)`
- `func DiscoverOneMain(path string) (pkgPath string, err error)`
- `func FindBlankImports(node *ast.File) []string`
- `func FormatImports(f *ast.File) map[string]string`
- `func FuncVarExists(f *ast.File, goImport, methodSignature string) bool`
- `func HasAnyStructFieldsInPkg(pkgPath, structName string, fields []string) (bool, error)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/goanalysis"
```
