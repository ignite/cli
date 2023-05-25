---
sidebar_position: 0
title: cosmostxcollector
slug: /packages/cosmostxcollector
---

# cosmostxcollector

The package implements support for collecting transactions and events from Cosmos blockchains
into a data backend and it also adds support for querying the collected data.

## Transaction and event data collecting

Transactions and events can be collected using the `cosmostxcollector.Collector` type. This
type uses a `cosmosclient.Client` instance to fetch the data from each block and a data backend
adapter to save the data.

### Data backend adapters

Data backend adapters are used to query and save the collected data into different types of data
backends and must implement the `cosmostxcollector.adapter.Adapter` interface.

An adapter for PostgreSQL is already implemented in `cosmostxcollector.adapter.postgres.Adapter`.
This is the one used in the examples.

### Example: Data collection

The data collection example assumes that there is a PostgreSQL database running in the local
environment containing an empty database named "cosmos".

The required database tables will be created automatically by the collector the first time it is run.

When the application is run it will fetch all the transactions and events starting from one of the
recent blocks until the current block height and populate the database:

```go
package main

import (
	"context"
	"log"

	"github.com/ignite/cli/ignite/pkg/clictx"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
)

const (
	// Name of a local PostgreSQL database
	dbName = "cosmos"

	// Cosmos RPC address
	rpcAddr = "https://rpc.cosmos.network:443"
)

func collect(ctx context.Context, db postgres.Adapter) error {
	// Make sure that the data backend schema is up to date
	if err := db.Init(ctx); err != nil {
		return err
	}

	// Init the Cosmos client
	client, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(rpcAddr))
	if err != nil {
		return err
	}

	// Get the latest block height
	latestHeight, err := client.LatestBlockHeight(ctx)
	if err != nil {
		return err
	}

	// Collect transactions and events starting from a block height.
	// The collector stops at the latest height available at the time of the call.
	collector := cosmostxcollector.New(db, client)
	if err := collector.Collect(ctx, latestHeight-50); err != nil {
		return err
	}

	return nil
}

func main() {
	ctx := clictx.From(context.Background())

	// Init an adapter for a local PostgreSQL database running with the default values
	params := map[string]string{"sslmode": "disable"}
	db, err := postgres.NewAdapter(dbName, postgres.WithParams(params))
	if err != nil {
		log.Fatal(err)
	}

	if err := collect(ctx, db); err != nil {
		log.Fatal(err)
	}
}
```

## Queries

Collected data can be queried through the data backend adapters using event queries or
cursor-based queries.

Queries support sorting, paging and filtering by using different options during creation.
The cursor-based ones also support the selection of specific fields or properties and also
passing arguments in cases where the query is a function.

By default no sorting, filtering nor paging is applied to the queries.

### Event queries

The event queries return events and their attributes as `[]cosmostxcollector.query.Event`.

### Example: Query events

The example reads transfer events from Cosmos' bank module and paginates the results.

```go
import (
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/query"
)

func queryBankTransferEvents(ctx context.Context, db postgres.Adapter) ([]query.Event, error) {
	// Create an event query that returns events of type "transfer"
	qry := query.NewEventQuery(
		query.WithFilters(
			// Filter transfer events from Cosmos' bank module
			postgres.FilterByEventType(banktypes.EventTypeTransfer),
		),
		query.WithPageSize(10),
		query.AtPage(1),
	)

	// Execute the query
	return db.QueryEvents(ctx, qry)
}
```

### Cursor-based queries

This type of queries is meant to be used in contexts where the Event queries are not
useful.

Cursor-based queries can query a single "entity" which can be a table, view or function
in relational databases or a collection or function in non relational data backends.

The result of these types of queries is a cursor that implements the `cosmostxcollector.query.Cursor`
interface.

### Example: Query events using cursors

```go
import (
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter/postgres"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/query"
)

func queryBankTransferEventIDs(ctx context.Context, db postgres.Adapter) (ids []int64, err error) {
	// Create a query that returns the IDs for events of type "transfer"
	qry := query.New(
		"event",
		query.Fields("id"),
		query.WithFilters(
			// Filter transfer events from Cosmos' bank module
			postgres.NewFilter("type", banktypes.EventTypeTransfer),
		),
		query.WithPageSize(10),
		query.AtPage(1),
		query.SortByFields(query.SortOrderAsc, "id"),
	)

	// Execute the query
	cr, err := db.Query(ctx, qry)
	if err != nil {
		return nil, err
	}

	// Read the results
	for cr.Next() {
		var eventID int64

		if err := cr.Scan(&eventID); err != nil {
			return nil, err
		}

		ids = append(ids, eventID)
	}

	return ids, nil
}
```
