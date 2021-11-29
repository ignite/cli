package pluginsrpc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

// List will list the plugins in a given state.
func (m *Manager) List(ctx context.Context, state PluginState) ([]string, error) {
	if state == Undefined {
		return []string{}, errors.New("undefined state")
	}

	var confPluginNames []string
	for _, cfgPlugin := range m.Config.Plugins {
		confPluginNames = append(confPluginNames, cfgPlugin.Name)
	}

	var pluginIds []string
	switch state {
	case Configured:
		return confPluginNames, nil
	case Downloaded:
		for _, cfgPlugin := range m.Config.Plugins {
			pluginId := getPluginId(cfgPlugin)
			dst, err := formatPluginHome(m.ChainId, pluginId)
			if err != nil {
				return []string{}, err
			}

			file, err := os.Stat(dst)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					continue
				}

				return []string{}, err
			}

			for _, confPluginName := range confPluginNames {
				if file.Name() == confPluginName {
					pluginIds = append(pluginIds, file.Name())
				}
			}
		}
	case Built:
		outputDir, err := formatPluginHome(m.ChainId, "output")
		if err != nil {
			return []string{}, err
		}

		for _, cfgPlugin := range m.Config.Plugins {
			pluginId := getPluginId(cfgPlugin)
			builtCommands, err := listFilesMatch(outputDir, fmt.Sprintf("%s*_cmd", pluginId))
			if err != nil {
				return []string{}, err
			}

			builtHooks, err := listFilesMatch(outputDir, "*_hook")
			if err != nil {
				return []string{}, err
			}

			for _, builtCommand := range builtCommands {
				splitFileName := strings.Split(builtCommand.Name(), "_")
				if pluginId == splitFileName[0] {
					pluginIds = append(pluginIds, pluginId)
				}
			}

			for _, builtHook := range builtHooks {
				splitFileName := strings.Split(builtHook.Name(), "_")
				if pluginId == splitFileName[0] {
					pluginIds = append(pluginIds, pluginId)
				}
			}
		}
	}

	return pluginIds, nil
}
