---
sidebar_position: 11
title: Cosmos SDK Versions (cosmosver)
slug: /packages/cosmosver
---

# Cosmos SDK Versions (cosmosver)

The `cosmosver` package parses, detects, and compares Cosmos SDK versions.
It is used to apply version-dependent behavior in Ignite.

For full API details, see the
[`cosmosver` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosver).

## Example: Detect and compare SDK version

```go
package main

import (
	"fmt"
	"log"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosver"
)

func main() {
	version, err := cosmosver.Detect(".")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("detected:", version)

	if version.GTE(cosmosver.StargateFiftyVersion) {
		fmt.Println("SDK is v0.50.0 or newer")
	}
}
```
