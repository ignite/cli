package cosmosmetric

import (
	"context"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosmetric/adapter"
	"golang.org/x/sync/errgroup"
)

// ClientCollecter defines the interface for Cosmos clients that support collection of transactions.
type ClientCollecter interface {
	CollectTXs(ctx context.Context, fromHeight int64, tc chan<- []cosmosclient.TX) error
}

// NewCollector creates a new transactions collector.
func NewCollector(db adapter.Saver, client ClientCollecter) Collector {
	return Collector{db, client}
}

// Collector defines a type to collect and save block transactions in a data backend.
type Collector struct {
	db     adapter.Saver
	client ClientCollecter
}

// Collect gathers transactions for all blocks starting from a specific height.
// Each group of block transactions is saved sequentially after being collected.
func (c Collector) Collect(ctx context.Context, fromHeight int64) error {
	tc := make(chan []cosmosclient.TX)
	wg, ctx := errgroup.WithContext(ctx)

	// Start collecting block transactions.
	// The transactions channel is closed by the client when all transactions
	// are collected or when an error occurs during the collection.
	wg.Go(func() error {
		return c.client.CollectTXs(ctx, fromHeight, tc)
	})

	// The transactions for each block are saved in "bulks" so they are not
	// kept in memory. Also they are saved sequentially to avoid block height
	// gaps that can occur if a group of transactions from a previous block
	// fail to be saved.
	for txs := range tc {
		if err := c.db.Save(ctx, txs); err != nil {
			return err
		}
	}

	// Any collection error is returned after the successfully
	// collected transactions are saved.
	return wg.Wait()
}
