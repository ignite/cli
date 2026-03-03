---
sidebar_position: 13
title: Cosmos Source Analysis (cosmosanalysis)
slug: /packages/cosmosanalysis
---

# Cosmos Source Analysis (cosmosanalysis)

The `cosmosanalysis` package provides static analysis helpers for Cosmos SDK-based projects, especially for app structure and interface/embed discovery.

For full API details, see the
[`cosmosanalysis` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis).

## When to use

- Validate that a directory is a Cosmos chain project before running codegen.
- Locate key app files and embedded types in Cosmos app sources.
- Detect interface implementations across module files.

## Key APIs

- `IsChainPath(path string) error`
- `FindAppFilePath(chainRoot string) (path string, err error)`
- `ValidateGoMod(module *modfile.File) error`
- `FindImplementation(modulePath string, interfaceList []string) (found []string, err error)`
- `DeepFindImplementation(modulePath string, interfaceList []string) (found []string, err error)`
- `FindEmbed(modulePath string, targetEmbeddedTypes []string) (found []string, err error)`
- `FindEmbedInFile(n ast.Node, targetEmbeddedTypes []string) (found []string)`

## Common Tasks

- Call `IsChainPath` early to fail fast on unsupported project layouts.
- Use `FindAppFilePath` before AST transformations that require the chain app entrypoint.
- Use `FindImplementation`/`DeepFindImplementation` to verify generated modules are wired as expected.

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis"
```
