package plugins

// import (
// 	"context"
// 	"fmt"

// 	"github.com/spf13/cobra"
// 	chaincfg "github.com/tendermint/starport/starport/chainconfig"
// )

// // NOTE: Go plugins were added in 1.8, so using them is questionable
// // unless we can build the plugins with a greater go version and somehow
// // grab the exposed functions with that build.

// // Plugin manager
// type Manager struct {
// 	ChainId string `yaml:"chain_id"`
// 	Config  chaincfg.Config

// 	cmdPlugins  []CmdPlugin
// 	hookPlugins []HookPlugin
// }

// func NewManager(chainId string, cfg chaincfg.Config) Manager {
// 	return Manager{
// 		ChainId: chainId,
// 		Config:  cfg,
// 	}
// }

// // Run will go through all of the steps:
// // - Downloading
// // - Building
// // - Cache .so files
// // - Execution and Injection
// func (m *Manager) RunAll(ctx context.Context, rootCommand *cobra.Command) error {
// 	if len(m.Config.Plugins) == 0 {
// 		return nil
// 	}

// 	if err := m.PullBuild(ctx); err != nil {
// 		return err
// 	}

// 	// Extraction
// 	if err := m.extractPlugins(ctx, rootCommand); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (m *Manager) PullBuild(ctx context.Context) error {
// 	if len(m.Config.Plugins) == 0 {
// 		return nil
// 	}

// 	// Check for change in config contents since last
// 	// Don't check for remote package changes, as theoretically we want it
// 	// to be up to the user to reload the plugins.
// 	// configChanged, err := PluginsChanged(cfg, m.ChainId)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// if !configChanged {
// 	// 	return nil
// 	// }

// 	fmt.Println("pulling")
// 	// Pull
// 	if err := m.pull(ctx); err != nil {
// 		return err
// 	}

// 	fmt.Println("building")
// 	// Build
// 	if err := m.build(ctx); err != nil {
// 		return err
// 	}

// 	return nil
// }
