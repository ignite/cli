---
sidebar_position: 46
title: Xast (xast)
slug: /packages/xast
---

# Xast (xast)

The `xast` package contains reusable utilities used by Ignite CLI internals.

For full API details, see the
[`xast` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/xast).

## Key APIs

- `var AppendFuncCodeAtLine = AppendFuncAtLine`
- `var ErrStop = errors.New("ast stop")`
- `func AppendFunction(fileContent string, function string) (modifiedContent string, err error)`
- `func AppendImports(fileContent string, imports ...ImportOptions) (string, error)`
- `func InsertGlobal(fileContent string, globalType GlobalType, globals ...GlobalOptions) (modifiedContent string, err error)`
- `func Inspect(n ast.Node, f func(n ast.Node) error) (err error)`
- `func ModifyCaller(content, callerExpr string, modifiers func([]string) ([]string, error)) (string, error)`
- `func ModifyFunction(content string, funcName string, functions ...FunctionOptions) (string, error)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/xast"
```
