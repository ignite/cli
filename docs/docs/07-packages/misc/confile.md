---
sidebar_position: 23
title: Confile (confile)
slug: /packages/confile
---

# Confile (confile)

is helper to load and overwrite configuration files.

For full API details, see the
[`confile` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/confile).

## Key APIs

- `var DefaultJSONEncodingCreator = &JSONEncodingCreator{}`
- `var DefaultTOMLEncodingCreator = &TOMLEncodingCreator{}`
- `var DefaultYAMLEncodingCreator = &YAMLEncodingCreator{}`
- `type ConfigFile struct{ ... }`
- `type Decoder interface{ ... }`
- `type EncodeDecoder interface{ ... }`
- `type Encoder interface{ ... }`
- `type Encoding struct{ ... }`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/confile"
```
