---
sidebar_position: 0
title: Indexer (cosmostxcollector)
slug: /packages/cosmostxcollector
---

# Indexer (cosmostxcollector)

The `cosmostxcollector` package streams transactions from a Cosmos client and persists them through a storage adapter.

For full API details, see the
[`cosmostxcollector` Go package documentation](https://pkg.go.dev/github.com/ignite/cli/v29/ignite/pkg/cosmostxcollector).

## When to use

- Build lightweight indexers that ingest transactions block by block.
- Persist transaction history in SQL or custom storage backends.
- Reuse transaction collection logic from `cosmosclient` without duplicating polling code.

## Key APIs

- `New(db adapter.Saver, client TXsCollector) Collector`
- `(Collector) Collect(ctx context.Context, fromHeight int64) error`
- `type TXsCollector interface{ CollectTXs(...) }`

## Common Tasks

- Implement `adapter.Saver` for your storage layer and pass it to `New`.
- Start collection at a chosen block height with `Collect`.
- Use a custom `TXsCollector` implementation for test doubles or alternate chain clients.

## Basic import

```go
import "github.com/ignite/cli/v29/ignite/pkg/cosmostxcollector"
```
