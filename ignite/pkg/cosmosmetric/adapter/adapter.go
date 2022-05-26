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

	// SetupSchema creates or updates the data backend schema required to save the metrics.
	// This method must be called at least once to create the initial database schema.
	// Calling it when a schema already exists updates the existing schema to the latest
	// version if the current one is older.
	SetupSchema(context.Context) error

	// GetLatestHeight returns the height of the latest block known by the data backend.
	GetLatestHeight(context.Context) (int64, error)
}
