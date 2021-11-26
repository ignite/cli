package pluginsrpc

import (
	"context"
	"log"
	"os/exec"
	"path"

	"github.com/hashicorp/go-plugin"
	"github.com/lukerhoads/plugintypes"
	"github.com/spf13/cobra"
	chaincfg "github.com/tendermint/starport/starport/chainconfig"
)

func (m *Manager) extractPlugins(ctx context.Context, rootCmd *cobra.Command) error {
	pluginHome, err := formatPluginHome(m.ChainId, "")
	if err != nil {
		return err
	}

	for _, cfgPlugin := range m.Config.Plugins {
		pluginId := getPluginId(cfgPlugin)
		pluginDir := path.Join(pluginHome, pluginId)
		outputDir := path.Join(pluginHome, "output")

		cmdPlugins, err := extractCommandPlugins(ctx, pluginDir, outputDir, rootCmd, m.Config)
		if err != nil {
			return err
		}

		hookPlugins, err := extractHookPlugins(ctx, pluginDir, outputDir, rootCmd, m.Config)
		if err != nil {
			return err
		}

		m.cmdPlugins = append(m.cmdPlugins, cmdPlugins...)
		m.hookPlugins = append(m.hookPlugins, hookPlugins...)
	}

	return nil
}

func extractCommandPlugins(
	ctx context.Context,
	pluginDir string,
	outputDir string,
	parentCommand *cobra.Command,
	cfg chaincfg.Config,
) ([]ExtractedCommandModule, error) {
	pluginFiles, err := listFiles(outputDir, `*_cmd`)
	if err != nil {
		return nil, err
	}

	if len(pluginFiles) == 0 {
		return []ExtractedCommandModule{}, nil
	}

	var extractedCommandModules []ExtractedCommandModule
	for _, pluginFile := range pluginFiles {
		pluginDir := path.Join(outputDir, pluginFile.Name())
		PluginMap := BasePluginMap

		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: HandshakeConfig,
			Plugins:         PluginMap,
			Cmd:             exec.Command(pluginDir),
		})

		rpcClient, err := client.Client()
		if err != nil {
			return []ExtractedCommandModule{}, err
		}

		raw, err := rpcClient.Dispense("command_map")
		if err != nil {
			return []ExtractedCommandModule{}, err
		}

		cmdMapper := raw.(plugintypes.CommandMapper)

		// Edit pluginMap off of that, and then make execute functions for each thing
		for _, loadedModule := range cmdMapper.Commands() {
			NewPluginMap := map[string]plugin.Plugin{
				loadedModule: &plugintypes.CommandModulePlugin{},
			}

			client2 := plugin.NewClient(&plugin.ClientConfig{
				HandshakeConfig: HandshakeConfig,
				Plugins:         NewPluginMap,
				Cmd:             exec.Command(pluginDir),
			})

			rpcClient2, err := client2.Client()
			if err != nil {
				return []ExtractedCommandModule{}, err
			}

			raw2, err := rpcClient2.Dispense(loadedModule)
			if err != nil {
				return []ExtractedCommandModule{}, err
			}

			cmdModule := raw2.(plugintypes.CommandModule)
			extractedCommandModules = append(extractedCommandModules, ExtractedCommandModule{
				ParentCommand: []string{},
				Name:          cmdModule.GetName(),
				Usage:         cmdModule.GetUsage(),
				ShortDesc:     cmdModule.GetShortDesc(),
				LongDesc:      cmdModule.GetLongDesc(),
				NumArgs:       cmdModule.GetNumArgs(),
				Exec: func(cmd *cobra.Command, args []string) error {
					client := plugin.NewClient(&plugin.ClientConfig{
						HandshakeConfig: HandshakeConfig,
						Plugins:         NewPluginMap,
						Cmd:             exec.Command(pluginDir),
					})

					rpcClient, err := client.Client()
					if err != nil {
						return err
					}

					raw, err := rpcClient.Dispense(loadedModule)
					if err != nil {
						return err
					}

					cmdModuleExec := raw.(plugintypes.CommandModule)
					return cmdModuleExec.Exec(cmd, args)
				},
			})

			log.Println("success!", cmdModule)
			client2.Kill()
		}

		client.Kill()
	}

	return extractedCommandModules, nil
}

func extractHookPlugins(
	ctx context.Context,
	pluginDir string,
	outputDir string,
	parentCommand *cobra.Command,
	cfg chaincfg.Config,
) ([]HookModule, error) {
	pluginFiles, err := listFiles(outputDir, `.*_hook`)
	if err != nil {
		return nil, err
	}

	if len(pluginFiles) == 0 {
		return []HookModule{}, nil
	}

	var hookPlugins []HookModule
	for _, pluginFile := range pluginFiles {
		pluginDir := path.Join(outputDir, pluginFile.Name())
		PluginMap := BasePluginMap

		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: HandshakeConfig,
			Plugins:         PluginMap,
			Cmd:             exec.Command(pluginDir),
		})
		defer client.Kill()

		rpcClient, err := client.Client()
		if err != nil {
			return []HookModule{}, err
		}

		raw, err := rpcClient.Dispense("hook")
		if err != nil {
			return []HookModule{}, err
		}

		cmdModule := raw.(HookModule)
		hookPlugins = append(hookPlugins, cmdModule)
	}

	return hookPlugins, nil
}
