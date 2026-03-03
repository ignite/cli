---
sidebar_position: 8
title: Go Command Helpers (gocmd)
slug: /packages/gocmd
---

# Go Command Helpers (gocmd)

The `gocmd` package provides helpers around `CommandInstall`, `Build`, and `BuildPath`.

For full API details, see the
[`gocmd` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/gocmd).

## Key APIs

- `const CommandInstall = "install" ...`
- `func Build(ctx context.Context, out, path string, flags []string, options ...exec.Option) error`
- `func BuildPath(ctx context.Context, output, binary, path string, flags []string, ...) error`
- `func BuildTarget(goos, goarch string) string`
- `func Env(name string) (string, error)`
- `func Fmt(ctx context.Context, path string, options ...exec.Option) error`
- `func Get(ctx context.Context, path string, pkgs []string, options ...exec.Option) error`
- `func GoImports(ctx context.Context, path string) error`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/gocmd"
```
