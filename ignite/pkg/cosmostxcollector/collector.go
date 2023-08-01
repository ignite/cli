package cosmostxcollector

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/cosmostxcollector/adapter"
)

// TXsCollector defines the interface for Cosmos clients that support collection of transactions.
//
//go:generate mockery --name TXsCollector --filename txs_collector.go --with-expecter
type TXsCollector interface {
	CollectTXs(ctx context.Context, fromHeight int64, tc chan<- []cosmosclient.TX) error
}

// New creates a new Cosmos transaction collector.
func New(db adapter.Saver, client TXsCollector) Collector {
	return Collector{db, client}
}

// Collector defines a type to collect and save Cosmos transactions in a data backend.
type Collector struct {
	db     adapter.Saver
	client TXsCollector
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
	// kept in memory. Also, they are saved sequentially to avoid block height
	// gaps that can occur if a group of transactions from a previous block
	// fail to be saved.
	wg.Go(func() error {
		for txs := range tc {
			if err := c.db.Save(ctx, txs); err != nil {
				return err
			}
		}

		return nil
	})

	return wg.Wait()
}
