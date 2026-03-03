---
sidebar_position: 2
title: Account Registry (cosmosaccount)
slug: /packages/cosmosaccount
---

# Account Registry (cosmosaccount)

The `cosmosaccount` package provides helpers around `KeyringServiceName`, `CoinTypeCosmos`, and `ErrAccountExists`.

For full API details, see the
[`cosmosaccount` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosaccount).

## Key APIs

- `const KeyringServiceName = "ignite" ...`
- `const CoinTypeCosmos = sdktypes.CoinType ...`
- `var ErrAccountExists = errors.New("account already exists")`
- `var KeyringHome = os.ExpandEnv("$HOME/.ignite/accounts")`
- `type Account struct{ ... }`
- `type AccountDoesNotExistError struct{ ... }`
- `type KeyringBackend string`
- `type Option func(*Registry)`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cosmosaccount"
```
