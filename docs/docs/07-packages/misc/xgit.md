---
sidebar_position: 51
title: Xgit (xgit)
slug: /packages/xgit
---

# Xgit (xgit)

The `xgit` package provides helpers around `AreChangesCommitted`, `Clone`, and `InitAndCommit`.

For full API details, see the
[`xgit` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/xgit).

## Key APIs

- `func AreChangesCommitted(dir string) (bool, error)`
- `func Clone(ctx context.Context, urlRef, dir string) error`
- `func InitAndCommit(path string) error`
- `func IsRepository(path string) (bool, error)`
- `func RepositoryURL(path string) (string, error)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/xgit"
```
