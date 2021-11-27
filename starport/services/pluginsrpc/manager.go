package pluginsrpc

import (
	"context"

	"github.com/spf13/cobra"
	chaincfg "github.com/tendermint/starport/starport/chainconfig"
)

// Plugin manager manages the plugins a given scaffolded blockchain.
type Manager struct {
	ChainId string `yaml:"chain_id"`
	Config  chaincfg.Config

	cmdPlugins  []ExtractedCommandModule
	hookPlugins []ExtractedHookModule
}

func NewManager(chainId string, cfg chaincfg.Config) Manager {
	return Manager{
		ChainId: chainId,
		Config:  cfg,
	}
}

// RunAll runs through all plugin processing steps provided by manager.
func (m *Manager) RunAll(ctx context.Context, rootCommand *cobra.Command, args []string) (bool, error) {
	if len(m.Config.Plugins) == 0 {
		return false, nil
	}

	if err := m.PullBuild(ctx); err != nil {
		return false, err
	}

	// Inject plugins
	return m.InjectPlugins(ctx, rootCommand, args)
}

// PullBuild both pulls and builds plugins specified in config.yml file.
func (m *Manager) PullBuild(ctx context.Context) error {
	if len(m.Config.Plugins) == 0 {
		return nil
	}

	// Check for change in config contents since last
	// Don't check for remote package changes, as theoretically we want it
	// to be up to the user to reload the plugins.
	// configChanged, err := PluginsChanged(cfg, m.ChainId)
	// if err != nil {
	// 	return err
	// }

	// if !configChanged {
	// 	return nil
	// }

	// could also use the pulled git repo to fetch for potential changes, which is
	// definitely possible

	if err := m.Pull(ctx); err != nil {
		return err
	}

	if err := m.Build(ctx); err != nil {
		return err
	}

	return nil
}
