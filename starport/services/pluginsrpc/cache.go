package pluginsrpc

import (
	"context"

	plugintypes "github.com/tendermint/starport/starport/services/pluginsrpc/types"
)

// WE CAN CACHE THE ENTIRE EXTRACTEDHOOKMODULE LMFAO
func (m *Manager) Cache(ctx context.Context) error {
	return nil
}

type CachedCommand struct {
	ModuleName string
	PluginDir  string
	plugintypes.Command
}

type CachedHook struct {
	ModuleName string
	PluginDir  string
	plugintypes.Hook
}

func (m *Manager) RetrieveCached() ([]ExtractedCommandModule, error) {
	// cacheHome := formatPluginHome(m.ChainId, "cached")

	// cachedFiles := listFilesMatch()

	return []ExtractedCommandModule{}, nil
}
