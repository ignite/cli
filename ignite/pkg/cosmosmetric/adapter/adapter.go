package adapter

import (
	"context"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
)

// Saver is the interface that wraps the transactions save method.
type Saver interface {
	// Save a list of transactions into a data backend.
	Save(context.Context, []cosmosclient.TX) error
}

// Adapter defines the interface for data backend adaptors.
type Adapter interface {
	Saver

	// GetType returns the adapter type.
	GetType() string

	// GetLatestHeight returns the height of the latest block known by the data backend.
	GetLatestHeight(context.Context) (int64, error)
}
