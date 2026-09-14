# errors

Gno port of Go's `errors` package.

## Exported API

| Function | Description |
|----------|-------------|
| `New(text string) error` | Create a simple text error. |
| `Unwrap(err error) error` | Return the result of `err.Unwrap()`, or nil. |
| `Is(err, target error) bool` | Report whether any error in the chain matches target. |
| `Join(errs ...error) error` | Combine multiple errors; nil values are discarded. |

## Differences from Go

`errors.As` is omitted because it requires `reflect`, which Gno does not
support. Error types that want custom matching should implement
`Is(error) bool` instead.

Comparability of error values is determined at runtime via recover rather
than via `reflect.Type.Comparable`.
