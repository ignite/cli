package pluginsrpc

import (
	"context"
	"errors"
	"os"

	"github.com/spf13/cobra"
	chaincfg "github.com/tendermint/starport/starport/chainconfig"
)

var (
	requiredDirs = []string{"output", "cached"}
)

// Plugin manager manages the plugins a given scaffolded blockchain.
type Manager struct {
	ChainId string `yaml:"chain_id"`
	Config  chaincfg.Config

	cmdPlugins  []ExtractedCommandModule
	hookPlugins []ExtractedHookModule
}

func NewManager(chainId string, cfg chaincfg.Config) (Manager, error) {
	// Pull cmdplugins and hookplugins from cache
	m := Manager{
		ChainId: chainId,
		Config:  cfg,
	}

	for _, requiredDir := range requiredDirs {
		dirHome, err := formatPluginHome(m.ChainId, requiredDir)
		if err != nil {
			return Manager{}, err
		}

		_, err = os.Stat(dirHome)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				err := os.Mkdir(dirHome, 0755)
				if err != nil {
					return Manager{}, err
				}
			}
		}
	}

	return m, nil
}

// RunAll runs through all plugin processing steps provided by manager.
func (m *Manager) RunAll(ctx context.Context, rootCommand *cobra.Command, args []string) (bool, error) {
	if len(m.Config.Plugins) == 0 {
		return false, nil
	}

	if err := m.PullBuild(ctx, false); err != nil {
		return false, err
	}

	// Inject plugins
	return m.InjectPlugins(ctx, rootCommand, args)
}

// PullBuild both pulls and builds plugins specified in config.yml file.
func (m *Manager) PullBuild(ctx context.Context, force bool) error {
	if len(m.Config.Plugins) == 0 {
		return nil
	}

	// If the call does not force a rerun, check for a change in config
	if !force {
		// Check for change in config contents since last
		// Don't check for remote package changes, as theoretically we want it
		// to be up to the user to reload the plugins.
		configChanged, err := pluginsChanged(m.Config, m.ChainId)
		if err != nil {
			return err
		}

		if !configChanged {
			return nil
		}

		return nil
	}

	if err := m.Pull(ctx); err != nil {
		return err
	}

	if err := m.Build(ctx); err != nil {
		return err
	}

	return nil
}
