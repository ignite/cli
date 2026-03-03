---
sidebar_position: 29
title: Errors (errors)
slug: /packages/errors
---

# Errors (errors)

provides helpers for error creation, avoiding using different.

For full API details, see the
[`errors` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/errors).

## Key APIs

- `func As(err error, target any) bool`
- `func Errorf(format string, args ...any) error`
- `func Is(err, reference error) bool`
- `func Join(errs ...error) error`
- `func New(msg string) error`
- `func Unwrap(err error) error`
- `func WithStack(err error) error`
- `func Wrap(err error, msg string) error`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/errors"
```
