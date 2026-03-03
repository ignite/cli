---
sidebar_position: 5
title: Token Faucet (cosmosfaucet)
slug: /packages/cosmosfaucet
---

# Token Faucet (cosmosfaucet)

The `cosmosfaucet` package is a faucet to request tokens for sdk accounts.

For full API details, see the
[`cosmosfaucet` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmosfaucet).

## Key APIs

- `const DefaultAccountName = "faucet" ...`
- `func TryRetrieve(ctx context.Context, chainID, rpcAddress, faucetAddress, accountAddress string) (string, error)`
- `type ErrTransferRequest struct{ ... }`
- `type Faucet struct{ ... }`
- `type FaucetInfoResponse struct{ ... }`
- `type HTTPClient struct{ ... }`
- `type Option func(*Faucet)`
- `type TransferRequest struct{ ... }`

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cosmosfaucet"
```
