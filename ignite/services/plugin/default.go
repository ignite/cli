package plugin

import (
	"errors"
)

var defaultPlugins = []DefaultPlugin{
	{
		Use:     "network",
		Short:   "Launch a blockchain in production",
		Aliases: []string{},
		Path:    "github.com/ignite/cli-plugin-network@v0.1.1",
	},
}

// DefaultPlugin defines a default Ignite plugin.
type DefaultPlugin struct {
	Use     string
	Short   string
	Aliases []string
	Path    string
}

// GetDefaultPlugins returns the list of default Ignite plugins.
func GetDefaultPlugins() []DefaultPlugin {
	return defaultPlugins
}

// GetDefaultNetworkPlugin returns the default network plugin.
func GetDefaultNetworkPlugin() (DefaultPlugin, error) {
	for _, p := range defaultPlugins {
		if p.Use == "network" {
			return p, nil
		}
	}

	return DefaultPlugin{}, errors.New("default network plugin not found")
}
