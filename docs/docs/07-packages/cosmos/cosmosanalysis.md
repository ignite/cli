---
sidebar_position: 13
title: Cosmos Source Analysis (cosmosanalysis)
slug: /packages/cosmosanalysis
---

# Cosmos Source Analysis (cosmosanalysis)

The `cosmosanalysis` package provides a toolset for statically analysing Cosmos SDK's.

For full API details, see the
[`cosmosanalysis` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis).

## Key APIs

- `var AppEmbeddedTypes = []string{ ... }`
- `func DeepFindImplementation(modulePath string, interfaceList []string) (found []string, err error)`
- `func FindAppFilePath(chainRoot string) (path string, err error)`
- `func FindEmbed(modulePath string, targetEmbeddedTypes []string) (found []string, err error)`
- `func FindEmbedInFile(n ast.Node, targetEmbeddedTypes []string) (found []string)`
- `func FindImplementation(modulePath string, interfaceList []string) (found []string, err error)`
- `func FindImplementationInFile(n ast.Node, interfaceList []string) (found []string)`
- `func IsChainPath(path string) error`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis"
```
