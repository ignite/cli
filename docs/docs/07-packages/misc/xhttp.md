---
sidebar_position: 52
title: Xhttp (xhttp)
slug: /packages/xhttp
---

# Xhttp (xhttp)

The `xhttp` package contains reusable utilities used by Ignite CLI internals.

For full API details, see the
[`xhttp` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/xhttp).

## Key APIs

- `const ShutdownTimeout = time.Minute`
- `func ResponseJSON(w http.ResponseWriter, status int, data interface{}) error`
- `func Serve(ctx context.Context, s *http.Server) error`
- `type ErrorResponse struct{ ... }`
- `type ErrorResponseBody struct{ ... }`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/xhttp"
```
