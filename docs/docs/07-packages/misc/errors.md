---
sidebar_position: 29
title: Errors (errors)
slug: /packages/errors
---

# Errors (errors)

The `errors` package centralizes error creation, wrapping, inspection, and stack enrichment with a consistent API for Ignite code.

For full API details, see the
[`errors` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/errors).

## When to use

- Wrap low-level errors with contextual messages while preserving root causes.
- Check error categories with `Is`/`As` across package boundaries.
- Aggregate multiple failures into a single returned error.

## Key APIs

- `func New(msg string) error`
- `func Errorf(format string, args ...any) error`
- `func Wrap(err error, msg string) error`
- `func Wrapf(err error, format string, args ...any) error`
- `func WithStack(err error) error`
- `func Is(err, reference error) bool`
- `func As(err error, target any) bool`
- `func Unwrap(err error) error`
- `func Join(errs ...error) error`

## Common Tasks

- Return `Wrap(err, "context")` from boundaries where more diagnostic context is needed.
- Use `Is` for sentinel checks and `As` for typed error extraction.
- Use `Join` when multiple independent steps can fail in the same operation.

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/errors"
```
