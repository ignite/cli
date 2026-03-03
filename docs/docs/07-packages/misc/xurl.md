---
sidebar_position: 59
title: Xurl (xurl)
slug: /packages/xurl
---

# Xurl (xurl)

The `xurl` package normalizes and validates URL/address values across HTTP, HTTPS, WS, and TCP forms.

For full API details, see the
[`xurl` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/xurl).

## When to use

- Accept flexible user input and normalize it to valid URLs.
- Enforce protocol-specific address formats.
- Auto-fill missing ports for HTTP endpoints.

## Key APIs

- `HTTP(s string) (string, error)`
- `HTTPS(s string) (string, error)`
- `WS(s string) (string, error)`
- `TCP(s string) (string, error)`
- `HTTPEnsurePort(s string) string`

## Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/ignite/cli/v29/ignite/pkg/xurl"
)

func main() {
	addr := xurl.HTTPEnsurePort("localhost")
	httpURL, err := xurl.HTTP(addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(httpURL)
}
```
