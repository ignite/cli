package plugins

import (
	"context"

	chaincfg "github.com/tendermint/starport/starport/chainconfig"
)

// Caches .so files for rebuilding
func (m *Manager) cache(ctx context.Context, cfg chaincfg.Config) error {
	return nil
}
